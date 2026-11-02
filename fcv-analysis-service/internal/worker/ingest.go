package worker

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"strings"
	"time"

	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourusername/fair-contribution-view/backend/internal/ai"
	"github.com/yourusername/fair-contribution-view/backend/internal/github"
	"github.com/yourusername/fair-contribution-view/backend/internal/ingestion"
	"github.com/yourusername/fair-contribution-view/backend/internal/processor"
	"github.com/yourusername/fair-contribution-view/backend/internal/queue"
	"github.com/yourusername/fair-contribution-view/backend/internal/reports"
	"github.com/yourusername/fair-contribution-view/backend/internal/scoring"
)

type IngestWorker struct {
	DB          *pgxpool.Pool
	GitHub      *github.Client
	AI          *ai.Client
	CommitLimit int
}

func (w *IngestWorker) Handle(ctx context.Context, task *asynq.Task) error {
	payload, err := queue.DecodeIngest(task.Payload())
	if err != nil {
		return fmt.Errorf("decode ingestion task: %w", err)
	}

	_, err = w.DB.Exec(ctx,
		`UPDATE analyses SET status='processing', started_at=$2,
		 progress='{"stage":"fetching_commits"}'::jsonb WHERE id=$1`,
		payload.AnalysisID, time.Now().UTC())
	if err != nil {
		return err
	}

	repo, err := w.GitHub.GetRepository(ctx, payload.Owner, payload.Repo)
	if err != nil {
		return err
	}
	commits, err := w.GitHub.ListCommits(ctx, payload.Owner, payload.Repo,
		repo.DefaultBranch, w.CommitLimit)
	if err != nil {
		return err
	}

	for i, item := range commits {
		var userID *int64
		var login, email string
		displayName := item.Commit.Author.Name
		if item.Author != nil {
			userID = &item.Author.ID
			login = item.Author.Login
			displayName = item.Author.Login
		}
		email = item.Commit.Author.Email
		contributorID, resolveErr := ingestion.ResolveContributor(ctx, w.DB, payload.AnalysisID, userID, login, displayName, email)
		if resolveErr != nil {
			return fmt.Errorf("resolve contributor for commit %s: %w", item.SHA, resolveErr)
		}

		_, err = w.DB.Exec(ctx, `
			INSERT INTO commits
			(analysis_id, sha, contributor_id, message, author_name,
			 author_email, authored_at, html_url)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
			ON CONFLICT (analysis_id, sha) DO NOTHING
		`, payload.AnalysisID, item.SHA, contributorID, item.Commit.Message,
			item.Commit.Author.Name, item.Commit.Author.Email,
			item.Commit.Author.Date, item.HTMLURL)
		if err != nil {
			return err
		}

		_, err = w.DB.Exec(ctx, `
			UPDATE analyses SET progress=$2::jsonb WHERE id=$1
		`, payload.AnalysisID, fmt.Sprintf(
			`{"stage":"fetching_commits","completed":%d,"total":%d}`,
			i+1, len(commits)))
		if err != nil {
			return err
		}
	}

	// Classify commit messages through the Python service in batches. If it is
	// unavailable, fall back to deterministic Go heuristics and continue ingestion.
	for start := 0; start < len(commits); start += 100 {
		end := start + 100
		if end > len(commits) {
			end = len(commits)
		}
		batch := make([]ai.Commit, 0, end-start)
		for _, item := range commits[start:end] {
			batch = append(batch, ai.Commit{SHA: item.SHA, Message: item.Commit.Message})
		}
		var results []ai.Result
		if w.AI != nil {
			results, err = w.AI.Classify(ctx, batch)
		}
		if err == nil && len(results) == len(batch) {
			for _, result := range results {
				_, err = w.DB.Exec(ctx, `UPDATE commits SET category=$3,classification_confidence=$4,classification_explanation=$5 WHERE analysis_id=$1 AND sha=$2`, payload.AnalysisID, result.SHA, result.Category, result.Confidence, result.Explanation)
				if err != nil {
					return err
				}
			}
		} else {
			for _, item := range commits[start:end] {
				category := processor.ClassifyCommit(item.Commit.Message, nil)
				_, err = w.DB.Exec(ctx, `UPDATE commits SET category=$3,classification_confidence=$4,classification_explanation=$5 WHERE analysis_id=$1 AND sha=$2`, payload.AnalysisID, item.SHA, string(category), 0.45, "Go heuristic fallback; no diff files were analyzed.")
				if err != nil {
					return err
				}
			}
		}
	}

	// ------------------------------------
	// STEP 2: FETCH PULL REQUESTS
	// ------------------------------------

	pullRequests, err := w.GitHub.ListPullRequests(
		ctx,
		payload.Owner,
		payload.Repo,
		50,
	)
	if err != nil {
		return fmt.Errorf("fetch pull requests: %w", err)
	}

	for _, listedPR := range pullRequests {
		// Fetch detailed PR statistics.
		pr, err := w.GitHub.GetPullRequest(
			ctx,
			payload.Owner,
			payload.Repo,
			listedPR.Number,
		)
		if err != nil {
			return err
		}

		if pr.User.Login != "" {
			if err := w.ensureLoginContributor(ctx, payload.AnalysisID, pr.User.Login); err != nil {
				return err
			}
		}
		_, err = w.DB.Exec(ctx, `
			INSERT INTO pull_requests (
				analysis_id, pr_number, title, state,
				author_login, created_at, updated_at,
				merged_at, html_url, additions, deletions,
				changed_files, comments_count
			)
			VALUES (
				$1,$2,$3,$4,$5,$6,$7,
				$8,$9,$10,$11,$12,$13
			)
			ON CONFLICT (analysis_id, pr_number)
			DO UPDATE SET
				state=EXCLUDED.state,
				merged_at=EXCLUDED.merged_at,
				additions=EXCLUDED.additions,
				deletions=EXCLUDED.deletions,
				changed_files=EXCLUDED.changed_files
		`,
			payload.AnalysisID,
			pr.Number,
			pr.Title,
			pr.State,
			pr.User.Login,
			pr.CreatedAt,
			pr.UpdatedAt,
			pr.MergedAt,
			pr.HTMLURL,
			pr.Additions,
			pr.Deletions,
			pr.ChangedFiles,
			pr.Comments,
		)
		if err != nil {
			return fmt.Errorf("save pull request %d: %w", pr.Number, err)
		}

		// ------------------------------------
		// STEP 3: FETCH PR REVIEWS
		// ------------------------------------

		reviews, err := w.GitHub.ListReviews(
			ctx,
			payload.Owner,
			payload.Repo,
			pr.Number,
		)
		if err != nil {
			return err
		}

		for _, review := range reviews {
			if review.User.Login != "" {
				if err := w.ensureLoginContributor(ctx, payload.AnalysisID, review.User.Login); err != nil {
					return err
				}
			}
			_, err = w.DB.Exec(ctx, `
				INSERT INTO reviews (
					analysis_id, review_id, pr_number,
					reviewer_login, review_state,
					body, submitted_at, html_url
				)
				VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
				ON CONFLICT (analysis_id, review_id)
				DO UPDATE SET
					review_state=EXCLUDED.review_state,
					body=EXCLUDED.body
			`,
				payload.AnalysisID,
				review.ID,
				pr.Number,
				review.User.Login,
				review.State,
				review.Body,
				review.SubmittedAt,
				review.HTMLURL,
			)

			if err != nil {
				return fmt.Errorf("save review %d: %w", review.ID, err)
			}
		}
	}

	// ------------------------------------
	// STEP 4: FETCH ISSUES
	// ------------------------------------

	issues, err := w.GitHub.ListIssues(
		ctx,
		payload.Owner,
		payload.Repo,
		100,
	)
	if err != nil {
		return fmt.Errorf("fetch issues: %w", err)
	}

	for _, issue := range issues {
		if issue.User.Login != "" {
			if err := w.ensureLoginContributor(ctx, payload.AnalysisID, issue.User.Login); err != nil {
				return err
			}
		}
		_, err = w.DB.Exec(ctx, `
			INSERT INTO issues (
				analysis_id, issue_number, title,
				state, author_login, created_at,
				updated_at, closed_at, html_url,
				comments_count
			)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
			ON CONFLICT (analysis_id, issue_number)
			DO UPDATE SET
				state=EXCLUDED.state,
				updated_at=EXCLUDED.updated_at,
				closed_at=EXCLUDED.closed_at,
				comments_count=EXCLUDED.comments_count
		`,
			payload.AnalysisID,
			issue.Number,
			issue.Title,
			issue.State,
			issue.User.Login,
			issue.CreatedAt,
			issue.UpdatedAt,
			issue.ClosedAt,
			issue.HTMLURL,
			issue.Comments,
		)

		if err != nil {
			return fmt.Errorf("save issue %d: %w", issue.Number, err)
		}
	}

	if err := w.calculateMetrics(ctx, payload.AnalysisID); err != nil {
		return fmt.Errorf("calculate contributor metrics: %w", err)
	}
	if _, err := reports.Generate(ctx, w.DB, payload.AnalysisID); err != nil {
		return fmt.Errorf("generate report: %w", err)
	}
	_, err = w.DB.Exec(ctx, `UPDATE analyses SET status='completed', completed_at=$2, progress=$3::jsonb, error_message=NULL WHERE id=$1`, payload.AnalysisID, time.Now().UTC(), `{"stage":"completed"}`)
	if err != nil {
		return err
	}
	return nil
}

// calculateMetrics creates a transparent activity-proxy score for each contributor.
func (w *IngestWorker) calculateMetrics(ctx context.Context, analysisID string) error {
	rows, err := w.DB.Query(ctx, `
      SELECT c.id::text, c.github_login, c.display_name,
        (SELECT count(*) FROM commits x WHERE x.analysis_id=c.analysis_id AND x.contributor_id=c.id)::int,
        (SELECT count(*) FROM pull_requests p WHERE p.analysis_id=c.analysis_id AND lower(coalesce(p.author_login,''))=lower(coalesce(c.github_login,'')))::int,
        (SELECT count(*) FROM reviews r WHERE r.analysis_id=c.analysis_id AND lower(coalesce(r.reviewer_login,''))=lower(coalesce(c.github_login,'')))::int,
        (SELECT count(*) FROM issues i WHERE i.analysis_id=c.analysis_id AND lower(coalesce(i.author_login,''))=lower(coalesce(c.github_login,'')))::int
      FROM contributors c WHERE c.analysis_id=$1`, analysisID)
	if err != nil {
		return err
	}
	defer rows.Close()
	type item struct {
		id                            string
		commits, prs, reviews, issues int
	}
	all := make([]item, 0)
	maxActivity := 0
	for rows.Next() {
		var x item
		var login, name *string
		if err := rows.Scan(&x.id, &login, &name, &x.commits, &x.prs, &x.reviews, &x.issues); err != nil {
			return err
		}
		raw := x.commits + 2*x.prs + 2*x.reviews + x.issues
		if raw > maxActivity {
			maxActivity = raw
		}
		all = append(all, x)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	for _, x := range all {
		raw := float64(x.commits + 2*x.prs + 2*x.reviews + x.issues)
		code := scoring.Normalize(float64(x.commits), float64(maxActivity))
		collab := scoring.Normalize(float64(x.reviews+x.prs), float64(maxActivity))
		breadth := scoring.Normalize(float64(x.prs+x.reviews+x.issues), float64(maxActivity))
		consistency := scoring.Normalize(float64(x.commits+x.prs+x.reviews+x.issues), float64(maxActivity))
		confidence := 0.25
		if raw >= 5 {
			confidence = 0.55
		}
		if raw >= 15 {
			confidence = 0.75
		}
		result := scoring.Calculate(scoring.Components{Code: code, Collaboration: collab, Breadth: breadth, Consistency: consistency}, confidence)
		_, err := w.DB.Exec(ctx, `INSERT INTO contributor_metrics(analysis_id,contributor_id,commit_count,pr_count,review_count,issue_count,activity_score,score,confidence,score_version) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) ON CONFLICT(analysis_id,contributor_id) DO UPDATE SET commit_count=EXCLUDED.commit_count,pr_count=EXCLUDED.pr_count,review_count=EXCLUDED.review_count,issue_count=EXCLUDED.issue_count,activity_score=EXCLUDED.activity_score,score=EXCLUDED.score,confidence=EXCLUDED.confidence,score_version=EXCLUDED.score_version`, analysisID, x.id, x.commits, x.prs, x.reviews, x.issues, raw, result.Score, result.Confidence, result.Version)
		if err != nil {
			return err
		}
	}
	return nil
}

func (w *IngestWorker) ensureLoginContributor(ctx context.Context, analysisID, login string) error {
	var existing string
	err := w.DB.QueryRow(ctx, `SELECT id FROM contributors WHERE analysis_id=$1 AND lower(github_login)=lower($2) LIMIT 1`, analysisID, login).Scan(&existing)
	if err == nil {
		return nil
	}
	if err != pgx.ErrNoRows {
		return err
	}
	_, err = w.DB.Exec(ctx, `INSERT INTO contributors(analysis_id,github_login,display_name,identity_key) VALUES($1,$2,$2,$3) ON CONFLICT(analysis_id,identity_key) DO NOTHING`, analysisID, login, "login:"+strings.ToLower(login))
	return err
}

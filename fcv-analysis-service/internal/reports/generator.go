package reports

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Contributor struct {
	Login        string  `json:"login"`
	DisplayName  string  `json:"display_name"`
	Commits      int     `json:"commits"`
	PullRequests int     `json:"pull_requests"`
	Reviews      int     `json:"reviews"`
	Issues       int     `json:"issues"`
	Score        float64 `json:"score"`
	Confidence   float64 `json:"confidence"`
}
type Report struct {
	AnalysisID     string        `json:"analysis_id"`
	GeneratedAt    time.Time     `json:"generated_at"`
	ScoreVersion   string        `json:"score_version"`
	Interpretation string        `json:"interpretation"`
	Limitations    []string      `json:"limitations"`
	Contributors   []Contributor `json:"contributors"`
}

func Generate(ctx context.Context, db *pgxpool.Pool, analysisID string) (Report, error) {
	r := Report{AnalysisID: analysisID, GeneratedAt: time.Now().UTC(), ScoreVersion: "fcv-1.0", Interpretation: "Scores are explainable activity proxies, not measurements of human effort, skill, or individual worth.", Limitations: []string{"Private or unavailable activity may be missing.", "Commit counts and code volume do not establish effort or impact.", "Classification uses deterministic heuristics unless an external model is explicitly configured."}, Contributors: []Contributor{}}
	rows, err := db.Query(ctx, `SELECT c.github_login,c.display_name, m.commit_count,m.pr_count,m.review_count,m.issue_count,m.score,m.confidence FROM contributor_metrics m JOIN contributors c ON c.id=m.contributor_id WHERE m.analysis_id=$1 ORDER BY m.score DESC`, analysisID)
	if err != nil {
		return r, err
	}
	defer rows.Close()
	for rows.Next() {
		var c Contributor
		if err := rows.Scan(&c.Login, &c.DisplayName, &c.Commits, &c.PullRequests, &c.Reviews, &c.Issues, &c.Score, &c.Confidence); err != nil {
			return r, err
		}
		r.Contributors = append(r.Contributors, c)
	}
	if err := rows.Err(); err != nil {
		return r, err
	}
	b, err := json.Marshal(r)
	if err != nil {
		return r, err
	}
	_, err = db.Exec(ctx, `INSERT INTO analysis_reports(analysis_id,report) VALUES($1,$2::jsonb) ON CONFLICT(analysis_id) DO UPDATE SET report=EXCLUDED.report,generated_at=now()`, analysisID, string(b))
	return r, err
}

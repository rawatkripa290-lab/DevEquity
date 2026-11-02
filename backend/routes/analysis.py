from fastapi import APIRouter, HTTPException
from pydantic import BaseModel
from urllib.parse import urlparse

from services.github_services import (
    get_repository_details,
    get_repository_commits,
    get_commit_statistics,
    get_repository_pull_requests,
    get_pull_request_reviews,
    get_repository_issues,
    get_commit_files
)

from services.metrices import (
    calculate_contributors,
    calculate_code_changes,
    calculate_fair_scores,
    calculate_pull_request_metrices,
    calculate_review_metrices,
    calculate_issue_metrics,
    classify_commit,
    calculate_commit_category_metrics,
    calculate_consistency_metrics,
     calculate_weighted_fair_scores
)
router = APIRouter()
class RepositoryRequest(BaseModel):
    repo_url: str

@router.post("/analyze")
def analyze_repository(request: RepositoryRequest):
    repo_url = request.repo_url.strip().rstrip("/")
    parsed_url = urlparse(repo_url)

    if (
        parsed_url.scheme != "https"
        or parsed_url.netloc.lower() not in ["github.com", "www.github.com"]
    ):
        raise HTTPException(
            status_code=400,
            detail="Please enter a valid GitHub repository link"
        )

    parts = parsed_url.path.strip("/").split("/")

    if len(parts) != 2 or not all(parts):
        raise HTTPException(
            status_code=400,
            detail="URL must contain owner and repository name"
        )

    owner = parts[0]
    repo = parts[1].removesuffix(".git")

    details = get_repository_details(owner, repo)

    if details is None:
        raise HTTPException(
            status_code=404,
            detail="Repository not found or inaccessible"
        )

    commits = get_repository_commits(owner, repo)

    if commits is None:
        raise HTTPException(
            status_code=500,
            detail="Unable to fetch repository commits"
        )

    # 1. Commit metrics
    contributors = calculate_contributors(commits)

    code_changes = calculate_code_changes(
        commits,
        owner,
        repo,
        get_commit_files
    )

    # 2. Pull request metrics
    pull_requests = get_repository_pull_requests(owner, repo)

    if pull_requests is None:
        raise HTTPException(
            status_code=500,
            detail="Unable to fetch pull requests"
        )

    pr_metrics = calculate_pull_request_metrices(pull_requests)

    # 3. Review metrics
    review_data = []

    for pr in pull_requests[:10]:
        reviews = get_pull_request_reviews(
            owner,
            repo,
            pr["number"]
        )

        if reviews is not None:
            review_data.append({
                "pr_author": pr["author"],
                "reviews": reviews
            })

    review_metrics = calculate_review_metrices(review_data)

    # 4. Issue metrics
    issues = get_repository_issues(owner, repo)

    if issues is None:
        raise HTTPException(
            status_code=500,
            detail="Unable to fetch repository issues"
        )

    issue_metrics = calculate_issue_metrics(issues)

    # 5. Commit category metrics
    classified_commits = []

    for commit in commits:
        classified_commits.append({
            "sha": commit["sha"],
            "author": commit["author"],
            "message": commit["message"],
            "date": commit["date"],
            "category": classify_commit(commit["message"])
        })

    commit_category_metrics = calculate_commit_category_metrics(
        classified_commits
    )

    # 6. Consistency metrics
    consistency_metrics = calculate_consistency_metrics(commits)

    # 7. Calculate weighted fair scores
    fair_scores = calculate_weighted_fair_scores(
        contributors,
        code_changes,
        review_metrics,
        pr_metrics,
        issue_metrics,
        commit_category_metrics,
        consistency_metrics
    )

    # 8. Create lookup maps
    code_map = {x["contributor"]: x for x in code_changes}
    pr_map = {x["contributor"]: x for x in pr_metrics}
    review_map = {x["contributor"]: x for x in review_metrics}
    issue_map = {x["contributor"]: x for x in issue_metrics}
    category_map = {
        x["contributor"]: x for x in commit_category_metrics
    }
    consistency_map = {
        x["contributor"]: x for x in consistency_metrics
    }
    score_map = {x["contributor"]: x for x in fair_scores}

    # 9. Combine all metrics for each contributor
    enriched_contributors = []

    for contributor in contributors:
        name = contributor["contributor"]

        code = code_map.get(name, {})
        prs = pr_map.get(name, {})
        reviews = review_map.get(name, {})
        issue = issue_map.get(name, {})
        categories = category_map.get(name, {})
        consistency = consistency_map.get(name, {})
        score = score_map.get(name, {})

        enriched_contributors.append({
            "contributor": name,

            "commits": contributor["total_commits"],
            "total_commits": contributor["total_commits"],

            "additions": code.get("additions", 0),
            "deletions": code.get("deletions", 0),
            "total_changes": code.get("total_changes", 0),

            "total_prs": prs.get("total_prs", 0),
            "merged_prs": prs.get("merged_prs", 0),
            "closed_prs": prs.get("closed_prs", 0),
            "open_prs": prs.get("open_prs", 0),

            "total_reviews": reviews.get("total_reviews", 0),
            "approved_reviews": reviews.get("approved_reviews", 0),
            "changes_requested": reviews.get("changes_requested", 0),

            "total_issues": issue.get("total_issues", 0),
            "open_issues": issue.get("open_issues", 0),
            "closed_issues": issue.get("closed_issues", 0),

            "feature": categories.get("feature", 0),
            "bugfix": categories.get("bugfix", 0),
            "refactor": categories.get("refactor", 0),
            "docs": categories.get("docs", 0),
            "tests": categories.get("tests", 0),

            "active_days": consistency.get("active_days", 0),

            "score": score.get("fair_score", 0),
            "fair_score": score.get("fair_score", 0)
        })

    # 10. Return complete response
    return {
        "message": "Repository analyzed successfully",
        "data": details,
        "total_commits_fetched": len(commits),
        "contributors": enriched_contributors,
        "code_changes": code_changes,
        "commits": commits,
        "fair_scores": fair_scores,
        "pull_request_metrics": pr_metrics,
        "review_metrics": review_metrics,
        "issue_metrics": issue_metrics,
        "classified_commits": classified_commits,
        "commit_category_metrics": commit_category_metrics,
        "consistency_metrics": consistency_metrics
    }

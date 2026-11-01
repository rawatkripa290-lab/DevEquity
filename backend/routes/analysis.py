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

    contributors = calculate_contributors(commits)

    code_changes = calculate_code_changes(
    commits,
    owner,
    repo,
    get_commit_files)
    pull_requests = get_repository_pull_requests(owner, repo)
    if pull_requests is None:
        raise HTTPException(
        status_code=500,
        detail="Unable to fetch pull requests"
    )
    pr_metrics = calculate_pull_request_metrices(pull_requests)
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
        issues = get_repository_issues(owner, repo)
        if issues is None:
            raise HTTPException(
        status_code=500,
        detail="Unable to fetch repository issues"
    )
        issue_metrics = calculate_issue_metrics(issues)
        classified_commits = []
        for commit in commits:
            classified_commits.append({
        "sha": commit["sha"],
        "author": commit["author"],
        "message": commit["message"],
        "date": commit["date"],
        "category": classify_commit(commit["message"])
    })
        commit_category_metrics = calculate_commit_category_metrics(classified_commits)
        consistency_metrics = calculate_consistency_metrics(commits)
        fair_scores = calculate_weighted_fair_scores(
            contributors,
            code_changes,
            review_metrics,
            pr_metrics,
            issue_metrics,
            commit_category_metrics,
            consistency_metrics
        )
        return {
        "message": "Repository analyzed successfully",
        "data": details,
        "total_commits_fetched": len(commits),
        "contributors": contributors,
        "code_changes": code_changes,
        "commits": commits,
        "fair_scores":fair_scores,
        "pull_request_metrics": pr_metrics,
        "review_metrics": review_metrics,
        "issue_metrics": issue_metrics,
        "classified_commits": classified_commits,
        "commit_category_metrics": commit_category_metrics,
        "consistency_metrics": consistency_metrics,
        "fair_scores": fair_scores
    }
import requests
import os
from dotenv import load_dotenv

BASE_DIR = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))

load_dotenv(os.path.join(BASE_DIR, ".env"))

GITHUB_TOKEN = os.getenv("GITHUB_TOKEN")

HEADERS = {
    "Accept": "application/vnd.github+json"
}

if GITHUB_TOKEN:
    HEADERS["Authorization"] = f"Bearer {GITHUB_TOKEN}"
print("Token loaded:", bool(GITHUB_TOKEN))
print("Authorization header added:", "Authorization" in HEADERS)

def get_repository_details(owner, repo):
    url = f"https://api.github.com/repos/{owner}/{repo}"

    response = requests.get(
        url,
        headers=HEADERS,
        timeout=15
    )

    print("GitHub Status Code:", response.status_code)

    if response.status_code == 200:
        data = response.json()

        return {
            "repository": data["name"],
            "owner": data["owner"]["login"],
            "stars": data["stargazers_count"],
            "forks": data["forks_count"],
            "open_issues": data["open_issues_count"],
            "language": data["language"]
        }

    print("GitHub Error:", response.text)
    return None


def get_repository_commits(owner, repo):
    url = f"https://api.github.com/repos/{owner}/{repo}/commits"

    params = {
        "per_page": 100,
        "page": 1
    }

    response = requests.get(
        url,
        params=params,
        headers=HEADERS,
        timeout=15
    )

    if response.status_code == 200:
        commits = response.json()
        result = []

        for commit in commits:
            result.append({
                "sha": commit["sha"],
                "author": (
                    commit["author"]["login"]
                    if commit.get("author")
                    else commit["commit"]["author"]["name"]
                ),
                "message": commit["commit"]["message"],
                "date": commit["commit"]["author"]["date"]
            })

        return result

    print("GitHub Commit Error:", response.text)
    return None


def get_commit_statistics(owner, repo, sha):
    url = f"https://api.github.com/repos/{owner}/{repo}/commits/{sha}"

    response = requests.get(
        url,
        headers=HEADERS,
        timeout=15
    )

    if response.status_code == 200:
        data = response.json()
        stats = data["stats"]

        return {
            "additions": stats["additions"],
            "deletions": stats["deletions"],
            "total_changes": stats["total"]
        }

    print("GitHub Statistics Error:", response.text)
    return None
def get_repository_pull_requests(owner, repo):
    url = f"https://api.github.com/repos/{owner}/{repo}/pulls"

    params = {
        "state": "all",
        "per_page": 100,
        "page": 1
    }

    response = requests.get(
        url,
        headers=HEADERS,
        params=params,
        timeout=15
    )

    if response.status_code == 200:
        pull_requests = response.json()
        result = []

        for pr in pull_requests:
            result.append({
                "number": pr["number"],
                "author": (
                    pr["user"]["login"]
                    if pr.get("user")
                    else None
                ),
                "title": pr["title"],
                "state": pr["state"],
                "merged": pr["merged_at"] is not None,
                "created_at": pr["created_at"],
                "closed_at": pr["closed_at"]
            })

        return result

    print("Pull request Error:", response.text)
    return None
def get_pull_request_reviews(owner,repo,pr_number):
    url=(
    f"https://api.github.com/repos/"
        f"{owner}/{repo}/pulls/{pr_number}/reviews")
    response = requests.get(
        url,
        headers=HEADERS,
        timeout=15
    )
    if response.status_code == 200:
        return response.json()
    print("Review Error:" , response.text)
    return None

def get_repository_issues(owner, repo):

    url = f"https://api.github.com/repos/{owner}/{repo}/issues"

    params = {
        "state": "all",
        "per_page": 100,
        "page": 1
    }

    response = requests.get(
        url,
        headers=HEADERS,
        params=params,
        timeout=15
    )

    if response.status_code == 200:

        issues = response.json()
        result = []

        for issue in issues:

            # GitHub's issues API also returns pull requests.
            # We exclude them because PRs are already analyzed separately.
            if "pull_request" in issue:
                continue

            result.append({
                "number": issue["number"],
                "author": (
                    issue["user"]["login"]
                    if issue.get("user")
                    else None
                ),
                "title": issue["title"],
                "state": issue["state"],
                "created_at": issue["created_at"],
                "closed_at": issue["closed_at"]
            })

        return result

    print("Issue Error:", response.text)
    return None
def get_commit_files(owner, repo, sha):
    url = f"https://api.github.com/repos/{owner}/{repo}/commits/{sha}"

    response = requests.get(
        url,
        headers=HEADERS,
        timeout=15
    )

    if response.status_code == 200:
        data = response.json()
        files = []

        for file in data.get("files", []):
            files.append({
                "filename": file["filename"],
                "additions": file["additions"],
                "deletions": file["deletions"],
                "changes": file["changes"],
                "status": file["status"]
            })

        return files

    print("Commit Files Error:", response.text)
    return None
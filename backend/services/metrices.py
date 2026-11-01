from datetime import datetime
def calculate_contributors(commits):
    contributors = {}
    for commit in commits:
        author = commit["author"]

        if author not in contributors:
            contributors[author] = 0

        contributors[author] += 1

    result = []

    for author, count in contributors.items():
        result.append({
            "contributor": author,
            "total_commits": count
        })

    return sorted(
        result,
        key=lambda s: s["total_commits"],
        reverse=True
    )


def calculate_code_changes(commits, owner, repo, get_commit_files):
    contributors = {}

    for commit in commits:
        author = commit["author"]
        sha = commit["sha"]

        if author not in contributors:
            contributors[author] = {
                "additions": 0,
                "deletions": 0,
                "total_changes": 0
            }

        files = get_commit_files(owner, repo, sha)

        if files is None:
            continue

        meaningful_files = filter_commit_files(files)

        for file in meaningful_files:
            contributors[author]["additions"] += file["additions"]
            contributors[author]["deletions"] += file["deletions"]
            contributors[author]["total_changes"] += file["changes"]

    result = []

    for author, stats in contributors.items():
        result.append({
            "contributor": author,
            **stats
        })

    return result

def calculate_fair_scores(contributors, code_changes):

    code_map = {
        item["contributor"]: item
        for item in code_changes
    }

    max_commits = max(
        (item["total_commits"] for item in contributors),
        default=0
    )

    max_additions = max(
        (item["additions"] for item in code_changes),
        default=0
    )

    max_deletions = max(
        (item["deletions"] for item in code_changes),
        default=0
    )

    result = []

    for contributor in contributors:

        author = contributor["contributor"]
        commits = contributor["total_commits"]

        changes = code_map.get(author, {})

        additions = changes.get("additions", 0)
        deletions = changes.get("deletions", 0)

        # Calculate individual scores
        commit_score = (
            commits / max_commits * 100
            if max_commits > 0 else 0
        )

        addition_score = (
            additions / max_additions * 100
            if max_additions > 0 else 0
        )

        deletion_score = (
            deletions / max_deletions * 100
            if max_deletions > 0 else 0
        )

        fair_score = (
            commit_score * 0.40
            + addition_score * 0.35
            + deletion_score * 0.25
        )

        result.append({
            "contributor": author,
            "total_commits": commits,
            "additions": additions,
            "deletions": deletions,
            "commit_score": round(commit_score, 2),
            "addition_score": round(addition_score, 2),
            "deletion_score": round(deletion_score, 2),
            "fair_score": round(fair_score, 2)
        })

    return sorted(
        result,
        key=lambda item: item["fair_score"],
        reverse=True
    )

def calculate_pull_request_metrices(pull_requests):
    contributors = {}

    for pr in pull_requests:

        author = pr.get("author")

        if not author:
            continue

        if author not in contributors:
            contributors[author] = {
                "total_prs": 0,
                "merged_prs": 0,
                "closed_prs": 0,
                "open_prs": 0
            }

        contributors[author]["total_prs"] += 1

        if pr["merged"]:
            contributors[author]["merged_prs"] += 1

        elif pr["state"] == "closed":
            contributors[author]["closed_prs"] += 1

        else:
            contributors[author]["open_prs"] += 1

    result = []

    for author, stats in contributors.items():
        result.append({
            "contributor": author,
            **stats
        })

    return sorted(
        result,
        key=lambda item: item["merged_prs"],
        reverse=True
    )

def calculate_review_metrices(review_data):

    contributors = {}

    for item in review_data:

        pr_author = item["pr_author"]
        reviews = item["reviews"]

        for review in reviews:

            reviewer = (
                review["user"]["login"]
                if review.get("user")
                else None
            )

            if not reviewer or reviewer == pr_author:
                continue

            if reviewer not in contributors:
                contributors[reviewer] = {
                    "total_reviews": 0,
                    "approved_reviews": 0,
                    "changes_requested": 0,
                    "commented_reviews": 0
                }

            state = review["state"]

            contributors[reviewer]["total_reviews"] += 1

            if state == "APPROVED":
                contributors[reviewer]["approved_reviews"] += 1

            elif state == "CHANGES_REQUESTED":
                contributors[reviewer]["changes_requested"] += 1

            elif state == "COMMENTED":
                contributors[reviewer]["commented_reviews"] += 1

    result = []

    for reviewer, stats in contributors.items():
        result.append({
            "contributor": reviewer,
            **stats
        })

    return sorted(
        result,
        key=lambda item: item["total_reviews"],
        reverse=True
    )

def calculate_issue_metrics(issues):

    contributors = {}

    for issue in issues:

        author = issue.get("author")

        if not author:
            continue

        if author not in contributors:
            contributors[author] = {
                "total_issues": 0,
                "open_issues": 0,
                "closed_issues": 0
            }

        contributors[author]["total_issues"] += 1

        if issue["state"] == "open":
            contributors[author]["open_issues"] += 1

        elif issue["state"] == "closed":
            contributors[author]["closed_issues"] += 1

    result = []

    for author, stats in contributors.items():
        result.append({
            "contributor": author,
            **stats
        })

    return sorted(
        result,
        key=lambda item: item["total_issues"],
        reverse=True
    )
def classify_commit(message):
    message = message.lower()

    if any(word in message for word in ["feat", "feature", "add"]):
        return "feature"

    elif any(word in message for word in ["fix", "bug", "error", "issue"]):
        return "bugfix"

    elif any(word in message for word in ["refactor", "improve", "optimize"]):
        return "refactor"

    elif any(word in message for word in ["docs", "documentation", "readme"]):
        return "docs"

    elif any(word in message for word in ["test", "tests", "testing"]):
        return "tests"

    else:
        return "chore"
    
def calculate_commit_category_metrics(classified_commits):
    contributors = {}

    for commit in classified_commits:
        author = commit["author"]
        category = commit["category"]

        if author not in contributors:
            contributors[author] = {
                "feature": 0,
                "bugfix": 0,
                "refactor": 0,
                "docs": 0,
                "tests": 0,
                "chore": 0
            }

        contributors[author][category] += 1

    result = []

    for author, stats in contributors.items():
        result.append({
            "contributor": author,
            **stats
        })

    return result
def is_noise_file(filename):
    noise_files = [
        "package-lock.json",
        "yarn.lock",
        "pnpm-lock.yaml",
        "poetry.lock",
        "Pipfile.lock",
        "composer.lock"
    ]

    if filename in noise_files:
        return True

    if filename.endswith((".min.js", ".min.css", ".map", ".pyc")):
        return True

    if filename.startswith(("dist/", "build/", "__pycache__/")):
        return True

    return False
def filter_commit_files(files):
    meaningful_files = []

    for file in files:
        filename = file["filename"]

        if not is_noise_file(filename):
            meaningful_files.append(file)

    return meaningful_files
def calculate_consistency_metrics(commits):
    contributors = {}

    for commit in commits:
        author = commit["author"]

        if author not in contributors:
            contributors[author] = set()

        commit_date = datetime.fromisoformat(
            commit["date"].replace("Z", "+00:00")
        ).date()

        contributors[author].add(commit_date)

    result = []

    for author, active_days in contributors.items():
        result.append({
            "contributor": author,
            "active_days": len(active_days)
        })

    return result
def calculate_weighted_fair_scores(
    contributors,
    code_changes,
    review_metrics,
    pr_metrics,
    issue_metrics,
    category_metrics,
    consistency_metrics
):
    def normalize(value, maximum):
        if maximum == 0:
            return 0
        return (value / maximum) * 100

    code_map = {x["contributor"]: x for x in code_changes}
    review_map = {x["contributor"]: x for x in review_metrics}
    pr_map = {x["contributor"]: x for x in pr_metrics}
    issue_map = {x["contributor"]: x for x in issue_metrics}
    category_map = {x["contributor"]: x for x in category_metrics}
    consistency_map = {
        x["contributor"]: x for x in consistency_metrics
    }

    all_names = {x["contributor"] for x in contributors}

    max_changes = max(
        (x["total_changes"] for x in code_changes), default=0
    )
    max_reviews = max(
        (x["total_reviews"] for x in review_metrics), default=0
    )
    max_activity = max(
        (
            pr_map.get(name, {}).get("total_prs", 0)
            + issue_map.get(name, {}).get("total_issues", 0)
            for name in all_names
        ),
        default=0
    )
    max_docs_tests = max(
        (
            category_map.get(name, {}).get("docs", 0)
            + category_map.get(name, {}).get("tests", 0)
            for name in all_names
        ),
        default=0
    )
    max_active_days = max(
        (
            consistency_map.get(name, {}).get("active_days", 0)
            for name in all_names
        ),
        default=0
    )

    result = []

    for name in all_names:
        code = code_map.get(name, {})
        reviews = review_map.get(name, {})
        prs = pr_map.get(name, {})
        issues = issue_map.get(name, {})
        categories = category_map.get(name, {})
        consistency = consistency_map.get(name, {})

        code_score = normalize(
            code.get("total_changes", 0), max_changes
        )
        review_score = normalize(
            reviews.get("total_reviews", 0), max_reviews
        )
        activity = (
            prs.get("total_prs", 0)
            + issues.get("total_issues", 0)
        )
        activity_score = normalize(activity, max_activity)

        docs_tests = categories.get("docs", 0) + categories.get("tests", 0)
        docs_tests_score = normalize(docs_tests, max_docs_tests)

        consistency_score = normalize(
            consistency.get("active_days", 0), max_active_days
        )

        fair_score = (
            code_score * 0.40
            + review_score * 0.25
            + activity_score * 0.15
            + docs_tests_score * 0.10
            + consistency_score * 0.10
        )

        result.append({
            "contributor": name,
            "code_score": round(code_score, 2),
            "review_score": round(review_score, 2),
            "activity_score": round(activity_score, 2),
            "docs_tests_score": round(docs_tests_score, 2),
            "consistency_score": round(consistency_score, 2),
            "fair_score": round(fair_score, 2)
        })

    return sorted(
        result,
        key=lambda x: x["fair_score"],
        reverse=True
    )
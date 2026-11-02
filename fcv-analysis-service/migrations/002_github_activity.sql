
-- Pull requests
CREATE TABLE IF NOT EXISTS pull_requests (
    analysis_id UUID NOT NULL REFERENCES analyses(id) ON DELETE CASCADE,
    pr_number BIGINT NOT NULL,
    title TEXT NOT NULL DEFAULT '',
    state TEXT NOT NULL,
    author_login TEXT,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    merged_at TIMESTAMPTZ,
    html_url TEXT,
    additions INTEGER NOT NULL DEFAULT 0,
    deletions INTEGER NOT NULL DEFAULT 0,
    changed_files INTEGER NOT NULL DEFAULT 0,
    comments_count INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (analysis_id, pr_number)
);

-- Reviews
CREATE TABLE IF NOT EXISTS reviews (
    analysis_id UUID NOT NULL REFERENCES analyses(id) ON DELETE CASCADE,
    review_id BIGINT NOT NULL,
    pr_number BIGINT NOT NULL,
    reviewer_login TEXT,
    review_state TEXT,
    body TEXT,
    submitted_at TIMESTAMPTZ,
    html_url TEXT,
    PRIMARY KEY (analysis_id, review_id),
    FOREIGN KEY (analysis_id, pr_number)
        REFERENCES pull_requests(analysis_id, pr_number)
        ON DELETE CASCADE
);

-- Issues
CREATE TABLE IF NOT EXISTS issues (
    analysis_id UUID NOT NULL REFERENCES analyses(id) ON DELETE CASCADE,
    issue_number BIGINT NOT NULL,
    title TEXT NOT NULL DEFAULT '',
    state TEXT NOT NULL,
    author_login TEXT,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    closed_at TIMESTAMPTZ,
    html_url TEXT,
    comments_count INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (analysis_id, issue_number)
);

CREATE INDEX IF NOT EXISTS idx_pr_author
ON pull_requests(analysis_id, author_login);

CREATE INDEX IF NOT EXISTS idx_review_reviewer
ON reviews(analysis_id, reviewer_login);

CREATE INDEX IF NOT EXISTS idx_issue_author
ON issues(analysis_id, author_login);

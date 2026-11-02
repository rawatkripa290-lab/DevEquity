CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS repositories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    github_repo_id BIGINT UNIQUE NOT NULL,
    owner TEXT NOT NULL,
    name TEXT NOT NULL,
    full_name TEXT NOT NULL,
    default_branch TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS analyses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    repository_id UUID NOT NULL REFERENCES repositories(id),
    status TEXT NOT NULL CHECK (status IN ('queued','processing','completed','failed')),
    progress JSONB NOT NULL DEFAULT '{}'::jsonb,
    error_message TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS contributors (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    analysis_id UUID NOT NULL REFERENCES analyses(id) ON DELETE CASCADE,
    github_user_id BIGINT,
    github_login TEXT,
    display_name TEXT NOT NULL,
    UNIQUE (analysis_id, github_user_id)
);

-- Composite key allows the same Git SHA to appear in separate analysis snapshots.
CREATE TABLE IF NOT EXISTS commits (
    analysis_id UUID NOT NULL REFERENCES analyses(id) ON DELETE CASCADE,
    sha TEXT NOT NULL,
    contributor_id UUID REFERENCES contributors(id),
    message TEXT NOT NULL DEFAULT '',
    author_name TEXT,
    author_email TEXT,
    authored_at TIMESTAMPTZ,
    html_url TEXT,
    additions INTEGER NOT NULL DEFAULT 0,
    deletions INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (analysis_id, sha)
);

CREATE INDEX IF NOT EXISTS idx_analyses_status ON analyses(status);
CREATE INDEX IF NOT EXISTS idx_commits_contributor ON commits(analysis_id, contributor_id);

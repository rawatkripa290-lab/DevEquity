ALTER TABLE commits ADD COLUMN IF NOT EXISTS category TEXT NOT NULL DEFAULT 'unknown';
ALTER TABLE commits ADD COLUMN IF NOT EXISTS classification_confidence DOUBLE PRECISION NOT NULL DEFAULT 0;
ALTER TABLE commits ADD COLUMN IF NOT EXISTS classification_explanation TEXT NOT NULL DEFAULT 'Not classified';

ALTER TABLE contributors ADD COLUMN IF NOT EXISTS identity_key TEXT;
UPDATE contributors SET identity_key = CASE WHEN github_user_id IS NOT NULL THEN 'gh:' || github_user_id::text ELSE 'legacy:' || id::text END WHERE identity_key IS NULL;
CREATE UNIQUE INDEX IF NOT EXISTS uq_contributors_analysis_identity ON contributors(analysis_id, identity_key);

CREATE TABLE IF NOT EXISTS contributor_metrics (
  analysis_id UUID NOT NULL REFERENCES analyses(id) ON DELETE CASCADE,
  contributor_id UUID NOT NULL REFERENCES contributors(id) ON DELETE CASCADE,
  commit_count INTEGER NOT NULL DEFAULT 0,
  pr_count INTEGER NOT NULL DEFAULT 0,
  review_count INTEGER NOT NULL DEFAULT 0,
  issue_count INTEGER NOT NULL DEFAULT 0,
  additions INTEGER NOT NULL DEFAULT 0,
  deletions INTEGER NOT NULL DEFAULT 0,
  feature_count INTEGER NOT NULL DEFAULT 0,
  bugfix_count INTEGER NOT NULL DEFAULT 0,
  refactor_count INTEGER NOT NULL DEFAULT 0,
  test_count INTEGER NOT NULL DEFAULT 0,
  docs_count INTEGER NOT NULL DEFAULT 0,
  chore_count INTEGER NOT NULL DEFAULT 0,
  activity_score DOUBLE PRECISION NOT NULL DEFAULT 0,
  score DOUBLE PRECISION NOT NULL DEFAULT 0,
  confidence DOUBLE PRECISION NOT NULL DEFAULT 0,
  score_version TEXT NOT NULL DEFAULT 'fcv-1.0',
  PRIMARY KEY (analysis_id, contributor_id)
);

CREATE TABLE IF NOT EXISTS analysis_reports (
  analysis_id UUID PRIMARY KEY REFERENCES analyses(id) ON DELETE CASCADE,
  report JSONB NOT NULL,
  generated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

package ingestion

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ResolveContributor returns one contributor row per GitHub user. For commits
// without a GitHub account, it falls back to a normalized email identity.
func ResolveContributor(ctx context.Context, db *pgxpool.Pool, analysisID string, userID *int64, login, name, email string) (string, error) {
	display := name
	if login != "" {
		display = login
	}
	if display == "" {
		display = "Unknown contributor"
	}

	var id string
	if userID != nil && *userID > 0 {
		err := db.QueryRow(ctx, `SELECT id FROM contributors WHERE analysis_id=$1 AND github_user_id=$2`, analysisID, *userID).Scan(&id)
		if err == nil {
			return id, nil
		}
		if err != pgx.ErrNoRows {
			return "", err
		}
		err = db.QueryRow(ctx, `INSERT INTO contributors (analysis_id,github_user_id,github_login,display_name,identity_key) VALUES ($1,$2,$3,$4,$5) RETURNING id`, analysisID, *userID, nullString(login), display, fmt.Sprintf("gh:%d", *userID)).Scan(&id)
		return id, err
	}

	key := "unknown"
	if strings.TrimSpace(email) != "" {
		key = "email:" + strings.ToLower(strings.TrimSpace(email))
	} else if strings.TrimSpace(name) != "" {
		key = "name:" + strings.ToLower(strings.TrimSpace(name))
	}
	err := db.QueryRow(ctx, `SELECT id FROM contributors WHERE analysis_id=$1 AND identity_key=$2`, analysisID, key).Scan(&id)
	if err == nil {
		return id, nil
	}
	if err != pgx.ErrNoRows {
		return "", err
	}
	err = db.QueryRow(ctx, `INSERT INTO contributors (analysis_id,github_login,display_name,identity_key) VALUES ($1,$2,$3,$4) ON CONFLICT (analysis_id,identity_key) DO UPDATE SET display_name=EXCLUDED.display_name RETURNING id`, analysisID, nullString(login), display, key).Scan(&id)
	return id, err
}

func nullString(s string) any {
	if s == "" {
		return nil
	}
	return s
}

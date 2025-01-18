package setting

import (
	"context"
	"database/sql"
	"post-backend/internal/custom"
)

type SettingRepository interface {
	FindBy(ctx context.Context, tx *sql.Tx, by any) (Setting, error)
	FindAll(ctx context.Context, tx *sql.Tx) ([]Setting, error)
}

type SettingRepositoryImpl struct {
}

func NewSettingRepository() SettingRepository {
	return &SettingRepositoryImpl{}
}

func (s *SettingRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) ([]Setting, error) {
	sql := "select id, key, value, created_at, updated_at from setting"

	rows, err := tx.QueryContext(ctx, sql)
	if err != nil {
		return []Setting{}, err
	}

	var settings []Setting
	for rows.Next() {
		var setting Setting
		if err := rows.Scan(&setting.Id, &setting.Key, &setting.Value, &setting.CreatedAt, &setting.UpdatedAt); err != nil {
			return []Setting{}, err
		}

		settings = append(settings, setting)
	}

	return settings, nil
}

func (s *SettingRepositoryImpl) FindBy(ctx context.Context, tx *sql.Tx, by any) (Setting, error) {
	sql := ""
	switch by.(type) {
	case int:
		sql = "select id, key, value, created_at, updated_at from setting where id = ?"
	case string:
		sql = "select id, key, value, created_at, updated_at from setting where key = ?"
	}

	if sql == "" {
		return Setting{}, custom.ErrInternal
	}

	row, err := tx.QueryContext(ctx, sql, by)
	if err != nil {
		return Setting{}, err
	}

	var setting Setting
	if row.Next() {
		if err := row.Scan(&setting.Id, &setting.Key, &setting.Value, &setting.CreatedAt, &setting.UpdatedAt); err != nil {
			return Setting{}, err
		}

		return setting, nil
	}

	return Setting{}, custom.ErrNotFound
}

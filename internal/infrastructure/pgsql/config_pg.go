package pgsql

import (
	"database/sql"
	"fmt"
	"go-arch/internal/entity"
	"go-arch/internal/repository"

	"github.com/jmoiron/sqlx"
)

type configRepo struct {
	db *sqlx.DB
}

func NewConfigRepoPg(db *sqlx.DB) repository.ConfigRepository {
	return &configRepo{db: db}
}

func (r *configRepo) FindAll() ([]*entity.Config, error) {
	query := `
		SELECT id, key, value, created_at, updated_at
		FROM configs
		ORDER BY id ASC
	`
	var configs []*entity.Config
	err := r.db.Select(&configs, query)
	if err != nil {
		return nil, fmt.Errorf("failed fetching configs: %w", err)
	}

	if configs == nil {
		configs = []*entity.Config{}
	}
	return configs, nil
}

func (r *configRepo) FindByKey(key string) (*entity.Config, error) {
	query := `
		SELECT id, key, value, created_at, updated_at
		FROM configs
		WHERE key = $1
	`
	var item entity.Config
	err := r.db.Get(&item, query, key)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed fetching config by key '%s': %w", key, err)
	}
	return &item, nil
}

func (r *configRepo) GetMap() (map[string]string, error) {
	configs, err := r.FindAll()
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for _, cfg := range configs {
		result[cfg.Key] = cfg.Value
	}
	return result, nil
}

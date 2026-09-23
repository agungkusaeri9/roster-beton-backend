package repository

import "go-arch/internal/entity"

type ConfigRepository interface {
	FindAll() ([]*entity.Config, error)
	FindByKey(key string) (*entity.Config, error)
	GetMap() (map[string]string, error)
}

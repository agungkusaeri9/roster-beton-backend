package usecase

import (
	"go-arch/internal/entity"
	"go-arch/internal/repository"
)

type ConfigUsecase interface {
	GetAll() ([]*entity.Config, error)
	GetByKey(key string) (*entity.Config, error)
	GetMap() (map[string]string, error)
}

type configUsecase struct {
	repo repository.ConfigRepository
}

func NewConfigUsecase(repo repository.ConfigRepository) ConfigUsecase {
	return &configUsecase{repo: repo}
}

func (u *configUsecase) GetAll() ([]*entity.Config, error) {
	return u.repo.FindAll()
}

func (u *configUsecase) GetByKey(key string) (*entity.Config, error) {
	return u.repo.FindByKey(key)
}

func (u *configUsecase) GetMap() (map[string]string, error) {
	return u.repo.GetMap()
}

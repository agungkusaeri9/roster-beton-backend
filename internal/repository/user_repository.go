package repository

import "go-arch/internal/entity"

type UserRepository interface {
	All() ([]*entity.User, error)
	GetByUsername(username string) (*entity.User, error)
	Find(id string) (*entity.User, error)
	Create(user *entity.User) (*entity.User, error)
	Update(user *entity.User) (*entity.User, error)
	Delete(id string) error
}

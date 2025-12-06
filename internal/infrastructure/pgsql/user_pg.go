package pgsql

import (
	"go-arch/internal/entity"
	"go-arch/internal/repository"

	"github.com/jmoiron/sqlx"
)

type userRepo struct {
	db *sqlx.DB
}

func NewUserRepoPg(db *sqlx.DB) repository.UserRepository {
	return &userRepo{
		db: db,
	}
}

func (r *userRepo) All() ([]*entity.User, error) {
	var users []*entity.User
	err := r.db.Select(&users, "SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepo) Find(id string) (*entity.User, error) {
	var user entity.User
	err := r.db.Get(&user, "SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) GetByUsername(username string) (*entity.User, error) {
	var user entity.User
	err := r.db.Get(&user, "SELECT * FROM users WHERE username = $1", username)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) Create(user *entity.User) (*entity.User, error) {
	// Set default role if not provided
	if user.Role == "" {
		user.Role = "user"
	}
	
	err := r.db.Get(user, `
		INSERT INTO users (name, username, password, role)
		VALUES ($1, $2, $3, $4)
		RETURNING *`,
		user.Name, user.Username, user.Password, user.Role,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepo) Update(user *entity.User) (*entity.User, error) {
	err := r.db.Get(user, `
		UPDATE users
		SET name = $1, username = $2, password = $3, role = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $5
		RETURNING *`,
		user.Name, user.Username, user.Password, user.Role, user.ID,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepo) Delete(id string) error {
	_, err := r.db.Exec("DELETE FROM users WHERE id = $1", id)
	return err
}

package usecase

import (
	"go-arch/internal/entity"
	"go-arch/internal/errors"
	"go-arch/internal/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase interface {
	Login(username string, password string) (string, error)
	Register(name string, username string, password string) (string, error)
}

type authUsecase struct {
	repo      repository.UserRepository
	jwtSecret string
}

func NewAuthUseCase(repo repository.UserRepository, jwtSecret string) AuthUsecase {
	return &authUsecase{repo: repo, jwtSecret: jwtSecret}
}

func (uc *authUsecase) Login(username string, password string) (string, error) {
	user, err := uc.repo.GetByUsername(username)
	if err != nil {
		return "", errors.NewBadRequest("Username not found")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.NewBadRequest("Password not match")
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(uc.jwtSecret))
	if err != nil {
		return "", errors.NewBadRequest("Failed to signing token")
	}
	return signed, nil
}

func (uc *authUsecase) Register(name string, username string, password string) (string, error) {
	findUser, _ := uc.repo.GetByUsername(username)
	if findUser != nil {
		return "", errors.NewBadRequest("User already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.NewBadRequest("failed to has a password")
	}

	user, err := uc.repo.Create(&entity.User{Name: name, Username: username, Password: string(hashedPassword)})
	if err != nil {
		return "", err
	}
	return user.ID, nil
}

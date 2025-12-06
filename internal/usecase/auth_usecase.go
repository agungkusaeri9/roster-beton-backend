package usecase

import (
	"go-arch/internal/entity"
	"go-arch/internal/errors"
	"go-arch/internal/repository"
	"regexp"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type LoginResult struct {
	Token     string
	ExpiredAt time.Time
	User      LoginUserInfo
}

type LoginUserInfo struct {
	Username string
	Name     string
}

type AuthUsecase interface {
	Login(username string, password string) (*LoginResult, error)
	Register(name string, username string, password string) (string, error)
}

type authUsecase struct {
	repo              repository.UserRepository
	jwtSecret         string
	tokenExpiryDuration time.Duration
}

func NewAuthUseCase(repo repository.UserRepository, jwtSecret string, tokenExpiryDuration time.Duration) AuthUsecase {
	return &authUsecase{
		repo:              repo,
		jwtSecret:         jwtSecret,
		tokenExpiryDuration: tokenExpiryDuration,
	}
}

func (uc *authUsecase) Login(username string, password string) (*LoginResult, error) {
	// Validate input
	if err := uc.validateLoginInput(username, password); err != nil {
		return nil, err
	}

	user, err := uc.repo.GetByUsername(username)
	if err != nil {
		return nil, errors.NewBadRequest("Username not found")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.NewBadRequest("Password not match")
	}

	expiredAt := time.Now().Add(uc.tokenExpiryDuration)
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     expiredAt.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(uc.jwtSecret))
	if err != nil {
		return nil, errors.NewInternal("Failed to signing token")
	}

	return &LoginResult{
		Token:     signed,
		ExpiredAt: expiredAt,
		User: LoginUserInfo{
			Username: user.Username,
			Name:     user.Name,
		},
	}, nil
}

func (uc *authUsecase) Register(name string, username string, password string) (string, error) {
	// Validate input
	if err := uc.validateRegisterInput(name, username, password); err != nil {
		return "", err
	}

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

// validateRegisterInput validates register request input
func (uc *authUsecase) validateRegisterInput(name, username, password string) error {
	var validationErrors []string

	// Validate name
	if strings.TrimSpace(name) == "" {
		validationErrors = append(validationErrors, "name is required")
	} else if len(name) < 2 {
		validationErrors = append(validationErrors, "name must be at least 2 characters")
	} else if len(name) > 100 {
		validationErrors = append(validationErrors, "name must not exceed 100 characters")
	}

	// Validate username
	if strings.TrimSpace(username) == "" {
		validationErrors = append(validationErrors, "username is required")
	} else if len(username) < 3 {
		validationErrors = append(validationErrors, "username must be at least 3 characters")
	} else if len(username) > 150 {
		validationErrors = append(validationErrors, "username must not exceed 150 characters")
	} else if !isValidUsername(username) {
		validationErrors = append(validationErrors, "username can only contain letters, numbers, underscores, and hyphens")
	}

	// Validate password
	if strings.TrimSpace(password) == "" {
		validationErrors = append(validationErrors, "password is required")
	} else if len(password) < 6 {
		validationErrors = append(validationErrors, "password must be at least 6 characters")
	} else if len(password) > 255 {
		validationErrors = append(validationErrors, "password must not exceed 255 characters")
	}

	if len(validationErrors) > 0 {
		return errors.NewBadRequest(strings.Join(validationErrors, "; "))
	}

	return nil
}

// validateLoginInput validates login request input
func (uc *authUsecase) validateLoginInput(username, password string) error {
	var validationErrors []string

	// Validate username
	if strings.TrimSpace(username) == "" {
		validationErrors = append(validationErrors, "username is required")
	}

	// Validate password
	if strings.TrimSpace(password) == "" {
		validationErrors = append(validationErrors, "password is required")
	}

	if len(validationErrors) > 0 {
		return errors.NewBadRequest(strings.Join(validationErrors, "; "))
	}

	return nil
}

// isValidUsername checks if username contains only allowed characters
// Allowed: letters (a-z, A-Z), numbers (0-9), underscore (_), hyphen (-)
func isValidUsername(username string) bool {
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, username)
	return matched
}

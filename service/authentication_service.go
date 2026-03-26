package service

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/snck/book-keeper-api/model"
	"github.com/snck/book-keeper-api/repository"
	"golang.org/x/crypto/bcrypt"
)

type MyCustomClaims struct {
	ID       uuid.UUID `json:"id"`
	UserName string    `json:"user_name"`
	jwt.RegisteredClaims
}

type AuthenticationService struct {
	repository *repository.AuthenticationRepository
}

func NewAuthenticationService(repository *repository.AuthenticationRepository) *AuthenticationService {
	return &AuthenticationService{repository: repository}
}

func (s *AuthenticationService) AddNewUser(userName string, password string) (*model.User, error) {

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return nil, err
	}

	newUser := model.User{
		UserName:     userName,
		PasswordHash: string(passwordHash),
	}

	return s.repository.AddNewUser(newUser)
}

func (s *AuthenticationService) GetUser(userName string) (*model.User, error) {
	user, err := s.repository.GetUserByUserName(userName)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, nil
	}

	return user, nil
}

func (s *AuthenticationService) IsPasswordValid(passwordHash string, password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
		return false
	} else {
		return true
	}
}

func (s *AuthenticationService) GenerateToken(user model.User) (string, error) {
	key := os.Getenv("KEY")
	t := jwt.NewWithClaims(jwt.SigningMethodHS256,
		MyCustomClaims{
			ID:       user.ID,
			UserName: user.UserName,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			},
		})

	return t.SignedString([]byte(key))
}

func (s *AuthenticationService) ParseAndValidateToken(tokenStr string) (*MyCustomClaims, error) {

	claims := &MyCustomClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
		key := os.Getenv("KEY")
		return []byte(key), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, nil
	}

	return claims, nil
}

func (s *AuthenticationService) InvalidateToken(id uuid.UUID, token string, expiredAt time.Time) error {
	return s.repository.AddTokenToBlocklist(id, token, expiredAt)
}

func (s *AuthenticationService) IsTokenExistInBlocklist(token string) (bool, error) {
	return s.repository.IsTokenExistInBlocklist(token)
}

package service

import (
	"github.com/snck/book-keeper-api/model"
	"github.com/snck/book-keeper-api/repository"
	"golang.org/x/crypto/bcrypt"
)

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

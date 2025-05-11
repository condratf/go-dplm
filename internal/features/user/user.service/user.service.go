package userservice

import (
	"github.com/condratf/go-musthave-diploma-tpl/internal/errors_custom"
	"golang.org/x/crypto/bcrypt"
)

func NewUserService(
	userRepo userRepository,
) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) RegisterUser(login, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	err = s.userRepo.CreateUser(login, string(hashedPassword))
	if err != nil {
		return err
	}
	return nil
}

func (s *userService) LoginUser(login, password string) error {
	hashedPassword, err := s.userRepo.GetUserPassword(login)
	if err != nil {
		return err
	}
	if bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) != nil {
		return errors_custom.ErrInvalidAuth
	}

	return nil
}

package service

import (
	"errors"
	"log"

	"github.com/alexey-shedrin/avito-test-task/internal/model/dto/request"
	"github.com/alexey-shedrin/avito-test-task/internal/model/dto/response"
	"github.com/alexey-shedrin/avito-test-task/internal/model/entity"
	"github.com/alexey-shedrin/avito-test-task/internal/utils/token"
	"github.com/google/uuid"
)

var InvalidCredentials = errors.New("invalid credentials")

type UserRepository interface {
	Create(user *entity.User) error
	GetByEmail(email string) (*entity.User, error)
}

type UserService struct {
	userRepo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{
		userRepo: repo,
	}
}

func (s *UserService) DummyLogin(req *request.DummyLogin) (*response.DummyLogin, error) {
	log.SetPrefix("service.DummyLogin")
	jwt, err := token.GenerateJWT(req.Role)
	if err != nil {
		log.Printf("error: %v", err)

		return nil, err
	}

	return &response.DummyLogin{
		Token: jwt,
	}, nil
}

func (s *UserService) Register(req *request.Register) (*entity.User, error) {
	log.SetPrefix("service.Register")
	user := entity.User{
		Id:       uuid.New(),
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
	}
	if err := user.HashPassword(); err != nil {
		log.Printf("error: %v", err)

		return nil, err
	}

	if err := s.userRepo.Create(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserService) Login(req *request.Login) (*response.Login, error) {
	log.SetPrefix("service.Login")

	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	if !user.CheckPassword(req.Password) {
		return nil, InvalidCredentials
	}

	jwt, err := token.GenerateJWT(user.Role)
	if err != nil {
		log.Printf("error: %v", err)

		return nil, err
	}

	return &response.Login{
		Token: jwt,
	}, nil
}

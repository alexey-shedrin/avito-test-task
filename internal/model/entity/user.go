package entity

import (
	"github.com/alexey-shedrin/avito-test-task/internal/model/dto/response"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	EmployeeRole  = "employee"
	ModeratorRole = "moderator"
)

type User struct {
	Id       uuid.UUID
	Email    string
	Password string
	Role     string
}

func (u *User) ToResponse() *response.User {
	return &response.User{
		Id:    u.Id,
		Email: u.Email,
		Role:  u.Role,
	}
}

func (u *User) HashPassword() error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(bytes)

	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

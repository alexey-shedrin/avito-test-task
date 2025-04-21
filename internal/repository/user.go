package repository

import (
	"database/sql"
	"errors"
	"log"

	"github.com/alexey-shedrin/avito-test-task/internal/database"
	"github.com/alexey-shedrin/avito-test-task/internal/model/entity"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}
func (r *UserRepository) Create(user *entity.User) error {
	query := `INSERT INTO users (id, email, password, user_role) VALUES ($1, $2, $3, $4)`

	if _, err := r.db.Exec(query, user.Id, user.Email, user.Password, user.Role); err != nil {
		if database.IsUniqueViolation(err) {
			return ErrUserAlreadyExists
		}

		log.Printf("error: %v", err)

		return err
	}

	return nil
}

func (r *UserRepository) GetByEmail(email string) (*entity.User, error) {
	query := `SELECT id, email, password, user_role FROM users WHERE email = $1`

	var user entity.User
	if err := r.db.QueryRow(query, email).Scan(&user.Id, &user.Email, &user.Password, &user.Role); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		log.Printf("error: %v", err)

		return nil, err
	}

	return &user, nil
}

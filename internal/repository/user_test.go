package repository_test

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/alexey-shedrin/avito-test-task/internal/repository"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alexey-shedrin/avito-test-task/internal/model/entity"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repository.NewUserRepository(db)

	userID := uuid.New()

	user := &entity.User{
		Id:       userID,
		Email:    "test@example.com",
		Password: "hashedpassword",
		Role:     "user",
	}

	testCases := []struct {
		name          string
		mockSetup     func()
		expectedError error
	}{
		{
			name: "Success",
			mockSetup: func() {
				mock.ExpectExec("INSERT INTO users").
					WithArgs(user.Id, user.Email, user.Password, user.Role).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectedError: nil,
		},
		{
			name: "UserAlreadyExists",
			mockSetup: func() {
				pqErr := &pq.Error{Code: "23505"}
				mock.ExpectExec("INSERT INTO users").
					WithArgs(user.Id, user.Email, user.Password, user.Role).
					WillReturnError(pqErr)
			},
			expectedError: repository.ErrUserAlreadyExists,
		},
		{
			name: "OtherError",
			mockSetup: func() {
				someError := errors.New("some error")
				mock.ExpectExec("INSERT INTO users").
					WithArgs(user.Id, user.Email, user.Password, user.Role).
					WillReturnError(someError)
			},
			expectedError: errors.New("some error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()

			err = repo.Create(user)

			if tc.expectedError != nil {
				if tc.name == "OtherError" {
					require.Error(t, err)
					require.Equal(t, tc.expectedError.Error(), err.Error())
				} else {
					require.ErrorIs(t, err, tc.expectedError)
				}
			} else {
				require.NoError(t, err)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repository.NewUserRepository(db)

	email := "test@example.com"
	userID := uuid.New()
	expectedUser := &entity.User{
		Id:       userID,
		Email:    email,
		Password: "hashedpassword",
		Role:     "user",
	}

	testCases := []struct {
		name          string
		mockSetup     func()
		expectedUser  *entity.User
		expectedError error
	}{
		{
			name: "Success",
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"id", "email", "password", "user_role"}).
					AddRow(expectedUser.Id, expectedUser.Email, expectedUser.Password, expectedUser.Role)
				mock.ExpectQuery("SELECT id, email, password, user_role FROM users WHERE email = \\$1").
					WithArgs(email).
					WillReturnRows(rows)
			},
			expectedUser:  expectedUser,
			expectedError: nil,
		},
		{
			name: "UserNotFound",
			mockSetup: func() {
				mock.ExpectQuery("SELECT id, email, password, user_role FROM users WHERE email = \\$1").
					WithArgs(email).
					WillReturnError(sql.ErrNoRows)
			},
			expectedUser:  nil,
			expectedError: repository.ErrUserNotFound,
		},
		{
			name: "OtherError",
			mockSetup: func() {
				someError := errors.New("some error")
				mock.ExpectQuery("SELECT id, email, password, user_role FROM users WHERE email = \\$1").
					WithArgs(email).
					WillReturnError(someError)
			},
			expectedUser:  nil,
			expectedError: errors.New("some error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()

			user, err := repo.GetByEmail(email)

			if tc.expectedError != nil {
				if tc.name == "OtherError" {
					require.Error(t, err)
					require.Equal(t, tc.expectedError.Error(), err.Error())
				} else {
					require.ErrorIs(t, err, tc.expectedError)
				}
				require.Nil(t, user)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedUser, user)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

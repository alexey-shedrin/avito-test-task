package repository_test

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alexey-shedrin/avito-test-task/internal/model/entity"
	"github.com/alexey-shedrin/avito-test-task/internal/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type PVZRepositoryTestSuite struct {
	suite.Suite
	db   *sql.DB
	mock sqlmock.Sqlmock
	repo *repository.PVZRepository
}

func (s *PVZRepositoryTestSuite) SetupTest() {
	var err error
	s.db, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)
	s.repo = repository.NewPVZRepository(s.db)
}

func (s *PVZRepositoryTestSuite) TearDownTest() {
	s.db.Close()
}

func TestPVZRepositorySuite(t *testing.T) {
	suite.Run(t, new(PVZRepositoryTestSuite))
}

func (s *PVZRepositoryTestSuite) TestCreatePvz_Success() {
	pvz := &entity.Pvz{
		City: "Moscow",
	}

	s.mock.ExpectExec("INSERT INTO pvz \\(id, registration_date, city\\) VALUES \\(\\$1, \\$2, \\$3\\)").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), pvz.City).
		WillReturnResult(sqlmock.NewResult(1, 1))

	result, err := s.repo.CreatePvz(pvz)

	require.NoError(s.T(), err)
	require.NotNil(s.T(), result)
	require.Equal(s.T(), pvz.City, result.City)
	require.NotEqual(s.T(), uuid.UUID{}, result.Id)
	require.NotZero(s.T(), result.RegistrationDate)
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func (s *PVZRepositoryTestSuite) TestCreatePvz_DBError() {
	pvz := &entity.Pvz{
		City: "Moscow",
	}
	dbErr := errors.New("database error")

	s.mock.ExpectExec("INSERT INTO pvz \\(id, registration_date, city\\) VALUES \\(\\$1, \\$2, \\$3\\)").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), pvz.City).
		WillReturnError(dbErr)

	result, err := s.repo.CreatePvz(pvz)

	require.Error(s.T(), err)
	require.Equal(s.T(), dbErr, err)
	require.Nil(s.T(), result)
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

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

type ReceptionRepositoryTestSuite struct {
	suite.Suite
	db   *sql.DB
	mock sqlmock.Sqlmock
	repo *repository.ReceptionRepository
}

func (s *ReceptionRepositoryTestSuite) SetupTest() {
	var err error
	s.db, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)
	s.repo = repository.NewReceptionRepository(s.db)
}

func (s *ReceptionRepositoryTestSuite) TearDownTest() {
	s.db.Close()
}

func TestReceptionRepositorySuite(t *testing.T) {
	suite.Run(t, new(ReceptionRepositoryTestSuite))
}

func (s *ReceptionRepositoryTestSuite) TestGetOpenedReceptionId_Success() {
	pvzID := uuid.New()
	expectedID := uuid.New()

	s.mock.ExpectQuery("SELECT id FROM reception WHERE pvz_id = \\$1 AND status = 'in_progress'").
		WithArgs(pvzID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedID))

	id, err := s.repo.GetOpenedReceptionId(pvzID)

	require.NoError(s.T(), err)
	require.Equal(s.T(), expectedID, id)
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func (s *ReceptionRepositoryTestSuite) TestGetOpenedReceptionId_NotFound() {
	pvzID := uuid.New()

	s.mock.ExpectQuery("SELECT id FROM reception WHERE pvz_id = \\$1 AND status = 'in_progress'").
		WithArgs(pvzID).
		WillReturnError(sql.ErrNoRows)

	id, err := s.repo.GetOpenedReceptionId(pvzID)

	require.NoError(s.T(), err)
	require.Equal(s.T(), uuid.UUID{}, id)
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func (s *ReceptionRepositoryTestSuite) TestGetOpenedReceptionId_DBError() {
	pvzID := uuid.New()
	dbErr := errors.New("database error")

	s.mock.ExpectQuery("SELECT id FROM reception WHERE pvz_id = \\$1 AND status = 'in_progress'").
		WithArgs(pvzID).
		WillReturnError(dbErr)

	id, err := s.repo.GetOpenedReceptionId(pvzID)

	require.Error(s.T(), err)
	require.Equal(s.T(), dbErr, err)
	require.Equal(s.T(), uuid.UUID{}, id)
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func (s *ReceptionRepositoryTestSuite) TestCreateReception_Success() {
	reception := &entity.Reception{
		PvzId: uuid.New(),
	}

	s.mock.ExpectExec("INSERT INTO reception \\(id, reception_datetime, pvz_id, status\\) VALUES \\(\\$1, \\$2, \\$3, \\$4\\)").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), reception.PvzId, "in_progress").
		WillReturnResult(sqlmock.NewResult(1, 1))

	result, err := s.repo.CreateReception(reception)

	require.NoError(s.T(), err)
	require.NotNil(s.T(), result)
	require.Equal(s.T(), "in_progress", result.Status)
	require.NotEqual(s.T(), uuid.UUID{}, result.Id)
	require.NotZero(s.T(), result.DateTime)
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func (s *ReceptionRepositoryTestSuite) TestCreateReception_DBError() {
	reception := &entity.Reception{
		PvzId: uuid.New(),
	}
	dbErr := errors.New("database error")

	s.mock.ExpectExec("INSERT INTO reception \\(id, reception_datetime, pvz_id, status\\) VALUES \\(\\$1, \\$2, \\$3, \\$4\\)").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), reception.PvzId, "in_progress").
		WillReturnError(dbErr)

	result, err := s.repo.CreateReception(reception)

	require.Error(s.T(), err)
	require.Equal(s.T(), dbErr, err)
	require.Nil(s.T(), result)
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func (s *ReceptionRepositoryTestSuite) TestCreateProduct_Success() {
	product := &entity.Product{
		Type:        "smartphone",
		ReceptionId: uuid.New(),
	}

	s.mock.ExpectExec("INSERT INTO product \\(id, product_type, acceptance_datetime, reception_id\\) VALUES \\(\\$1, \\$2, \\$3, \\$4\\)").
		WithArgs(sqlmock.AnyArg(), product.Type, sqlmock.AnyArg(), product.ReceptionId).
		WillReturnResult(sqlmock.NewResult(1, 1))

	result, err := s.repo.CreateProduct(product)

	require.NoError(s.T(), err)
	require.NotNil(s.T(), result)
	require.Equal(s.T(), product.Type, result.Type)
	require.Equal(s.T(), product.ReceptionId, result.ReceptionId)
	require.NotEqual(s.T(), uuid.UUID{}, result.Id)
	require.NotZero(s.T(), result.DateTime)
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func (s *ReceptionRepositoryTestSuite) TestCreateProduct_DBError() {
	product := &entity.Product{
		Type:        "smartphone",
		ReceptionId: uuid.New(),
	}
	dbErr := errors.New("database error")

	s.mock.ExpectExec("INSERT INTO product \\(id, product_type, acceptance_datetime, reception_id\\) VALUES \\(\\$1, \\$2, \\$3, \\$4\\)").
		WithArgs(sqlmock.AnyArg(), product.Type, sqlmock.AnyArg(), product.ReceptionId).
		WillReturnError(dbErr)

	result, err := s.repo.CreateProduct(product)

	require.Error(s.T(), err)
	require.Equal(s.T(), dbErr, err)
	require.Nil(s.T(), result)
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func (s *ReceptionRepositoryTestSuite) TestDeleteLastProduct_Success() {
	receptionID := uuid.New()

	s.mock.ExpectExec("DELETE FROM product WHERE id = \\(SELECT id FROM product WHERE reception_id = \\$1 ORDER BY acceptance_datetime DESC LIMIT 1\\)").
		WithArgs(receptionID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := s.repo.DeleteLastProduct(receptionID)

	require.NoError(s.T(), err)
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func (s *ReceptionRepositoryTestSuite) TestDeleteLastProduct_NotFound() {
	receptionID := uuid.New()

	s.mock.ExpectExec("DELETE FROM product WHERE id = \\(SELECT id FROM product WHERE reception_id = \\$1 ORDER BY acceptance_datetime DESC LIMIT 1\\)").
		WithArgs(receptionID).
		WillReturnError(sql.ErrNoRows)

	err := s.repo.DeleteLastProduct(receptionID)

	require.Error(s.T(), err)
	require.Equal(s.T(), errors.New("product not found"), err)
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func (s *ReceptionRepositoryTestSuite) TestDeleteLastProduct_DBError() {
	receptionID := uuid.New()
	dbErr := errors.New("database error")

	s.mock.ExpectExec("DELETE FROM product WHERE id = \\(SELECT id FROM product WHERE reception_id = \\$1 ORDER BY acceptance_datetime DESC LIMIT 1\\)").
		WithArgs(receptionID).
		WillReturnError(dbErr)

	err := s.repo.DeleteLastProduct(receptionID)

	require.Error(s.T(), err)
	require.Equal(s.T(), dbErr, err)
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func (s *ReceptionRepositoryTestSuite) TestCloseLastReception_Success() {
	receptionID := uuid.New()

	s.mock.ExpectExec("UPDATE reception SET status = 'closed' WHERE id = \\$1").
		WithArgs(receptionID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	result, err := s.repo.CloseLastReception(receptionID)

	require.NoError(s.T(), err)
	require.NotNil(s.T(), result)
	require.Equal(s.T(), receptionID, result.Id)
	require.Equal(s.T(), "closed", result.Status)
	require.NotZero(s.T(), result.DateTime)
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func (s *ReceptionRepositoryTestSuite) TestCloseLastReception_DBError() {
	receptionID := uuid.New()
	dbErr := errors.New("database error")

	s.mock.ExpectExec("UPDATE reception SET status = 'closed' WHERE id = \\$1").
		WithArgs(receptionID).
		WillReturnError(dbErr)

	result, err := s.repo.CloseLastReception(receptionID)

	require.Error(s.T(), err)
	require.Equal(s.T(), dbErr, err)
	require.Nil(s.T(), result)
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

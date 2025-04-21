package service_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alexey-shedrin/avito-test-task/internal/model/entity"
	"github.com/alexey-shedrin/avito-test-task/internal/repository/mocks"
	"github.com/alexey-shedrin/avito-test-task/internal/service"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestReceptionService_CreateReception(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockReceptionRepository(ctrl)
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		receptionSvc := service.NewReceptionService(mockRepo, db)

		pvzID := uuid.New()
		expectedReception := &entity.Reception{
			PvzId: pvzID,
		}
		returnedReception := &entity.Reception{
			Id:    uuid.New(),
			PvzId: pvzID,
		}

		mock.ExpectBegin()
		mockRepo.EXPECT().GetOpenedReceptionId(pvzID).Return(uuid.Nil, nil)
		mockRepo.EXPECT().CreateReception(expectedReception).Return(returnedReception, nil)
		mock.ExpectCommit()

		result, err := receptionSvc.CreateReception(expectedReception)

		require.NoError(t, err)
		require.Equal(t, returnedReception, result)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Reception already opened", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockReceptionRepository(ctrl)
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		receptionSvc := service.NewReceptionService(mockRepo, db)

		pvzID := uuid.New()
		openedReceptionID := uuid.New()
		reception := &entity.Reception{
			PvzId: pvzID,
		}

		mock.ExpectBegin()
		mockRepo.EXPECT().GetOpenedReceptionId(pvzID).Return(openedReceptionID, nil)
		mock.ExpectRollback()

		result, err := receptionSvc.CreateReception(reception)

		require.Error(t, err)
		require.Equal(t, service.ReceptionAlreadyOpened, err)
		require.Nil(t, result)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Transaction begin error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockReceptionRepository(ctrl)
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		receptionSvc := service.NewReceptionService(mockRepo, db)

		expectedError := errors.New("db error")
		mock.ExpectBegin().WillReturnError(expectedError)

		reception := &entity.Reception{
			PvzId: uuid.New(),
		}

		result, err := receptionSvc.CreateReception(reception)

		require.Error(t, err)
		require.Equal(t, expectedError, err)
		require.Nil(t, result)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GetOpenedReceptionId error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockReceptionRepository(ctrl)
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		receptionSvc := service.NewReceptionService(mockRepo, db)

		expectedError := errors.New("repo error")
		pvzID := uuid.New()
		reception := &entity.Reception{
			PvzId: pvzID,
		}

		mock.ExpectBegin()
		mockRepo.EXPECT().GetOpenedReceptionId(pvzID).Return(uuid.Nil, expectedError)
		mock.ExpectRollback()

		result, err := receptionSvc.CreateReception(reception)

		require.Error(t, err)
		require.Equal(t, expectedError, err)
		require.Nil(t, result)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("CreateReception error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockReceptionRepository(ctrl)
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		receptionSvc := service.NewReceptionService(mockRepo, db)

		expectedError := errors.New("repo error")
		pvzID := uuid.New()
		reception := &entity.Reception{
			PvzId: pvzID,
		}

		mock.ExpectBegin()
		mockRepo.EXPECT().GetOpenedReceptionId(pvzID).Return(uuid.Nil, nil)
		mockRepo.EXPECT().CreateReception(reception).Return(nil, expectedError)
		mock.ExpectRollback()

		result, err := receptionSvc.CreateReception(reception)

		require.Error(t, err)
		require.Equal(t, expectedError, err)
		require.Nil(t, result)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestReceptionService_CreateProduct(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockReceptionRepository(ctrl)
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		receptionSvc := service.NewReceptionService(mockRepo, db)

		pvzID := uuid.New()
		receptionID := uuid.New()
		product := &entity.Product{
			Type: "Test Product",
		}
		expectedProduct := &entity.Product{
			Type:        "Test Product",
			ReceptionId: receptionID,
		}
		returnedProduct := &entity.Product{
			Id:          uuid.New(),
			Type:        "Test Product",
			ReceptionId: receptionID,
		}

		mock.ExpectBegin()
		mockRepo.EXPECT().GetOpenedReceptionId(pvzID).Return(receptionID, nil)
		mockRepo.EXPECT().CreateProduct(expectedProduct).Return(returnedProduct, nil)
		mock.ExpectCommit()

		result, err := receptionSvc.CreateProduct(product, pvzID)

		require.NoError(t, err)
		require.Equal(t, returnedProduct, result)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Reception not opened", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockReceptionRepository(ctrl)
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		receptionSvc := service.NewReceptionService(mockRepo, db)

		pvzID := uuid.New()
		product := &entity.Product{
			Type: "Test Product",
		}

		mock.ExpectBegin()
		mockRepo.EXPECT().GetOpenedReceptionId(pvzID).Return(uuid.Nil, nil)
		mock.ExpectRollback()

		result, err := receptionSvc.CreateProduct(product, pvzID)

		require.Error(t, err)
		require.Equal(t, service.ReceptionNotOpened, err)
		require.Nil(t, result)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Transaction begin error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockReceptionRepository(ctrl)
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		receptionSvc := service.NewReceptionService(mockRepo, db)

		expectedError := errors.New("db error")
		mock.ExpectBegin().WillReturnError(expectedError)

		product := &entity.Product{
			Type: "Test Product",
		}

		result, err := receptionSvc.CreateProduct(product, uuid.New())

		require.Error(t, err)
		require.Equal(t, expectedError, err)
		require.Nil(t, result)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GetOpenedReceptionId error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockReceptionRepository(ctrl)
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		receptionSvc := service.NewReceptionService(mockRepo, db)

		pvzID := uuid.New()
		expectedError := errors.New("repo error")
		product := &entity.Product{
			Type: "Test Product",
		}

		mock.ExpectBegin()
		mockRepo.EXPECT().GetOpenedReceptionId(pvzID).Return(uuid.Nil, expectedError)
		mock.ExpectRollback()

		result, err := receptionSvc.CreateProduct(product, pvzID)

		require.Error(t, err)
		require.Equal(t, expectedError, err)
		require.Nil(t, result)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("CreateProduct error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockReceptionRepository(ctrl)
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		receptionSvc := service.NewReceptionService(mockRepo, db)

		expectedError := errors.New("repo error")
		pvzID := uuid.New()
		receptionID := uuid.New()
		product := &entity.Product{
			Type: "Test Product",
		}
		expectedProduct := &entity.Product{
			Type:        "Test Product",
			ReceptionId: receptionID,
		}

		mock.ExpectBegin()
		mockRepo.EXPECT().GetOpenedReceptionId(pvzID).Return(receptionID, nil)
		mockRepo.EXPECT().CreateProduct(expectedProduct).Return(nil, expectedError)
		mock.ExpectRollback()

		result, err := receptionSvc.CreateProduct(product, pvzID)

		require.Error(t, err)
		require.Equal(t, expectedError, err)
		require.Nil(t, result)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestReceptionService_DeleteLastProduct(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockReceptionRepository(ctrl)
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		receptionSvc := service.NewReceptionService(mockRepo, db)

		pvzID := uuid.New()
		receptionID := uuid.New()

		mock.ExpectBegin()
		mockRepo.EXPECT().GetOpenedReceptionId(pvzID).Return(receptionID, nil)
		mockRepo.EXPECT().DeleteLastProduct(receptionID).Return(nil)
		mock.ExpectCommit()

		err = receptionSvc.DeleteLastProduct(pvzID)

		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Reception not opened", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockReceptionRepository(ctrl)
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		receptionSvc := service.NewReceptionService(mockRepo, db)

		pvzID := uuid.New()

		mock.ExpectBegin()
		mockRepo.EXPECT().GetOpenedReceptionId(pvzID).Return(uuid.Nil, nil)
		mock.ExpectRollback()

		err = receptionSvc.DeleteLastProduct(pvzID)

		require.Error(t, err)
		require.Equal(t, service.ReceptionNotOpened, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Transaction begin error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockReceptionRepository(ctrl)
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		receptionSvc := service.NewReceptionService(mockRepo, db)

		expectedError := errors.New("db error")
		mock.ExpectBegin().WillReturnError(expectedError)

		err = receptionSvc.DeleteLastProduct(uuid.New())

		require.Error(t, err)
		require.Equal(t, expectedError, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GetOpenedReceptionId error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockReceptionRepository(ctrl)
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		receptionSvc := service.NewReceptionService(mockRepo, db)

		pvzID := uuid.New()
		expectedError := errors.New("repo error")

		mock.ExpectBegin()
		mockRepo.EXPECT().GetOpenedReceptionId(pvzID).Return(uuid.Nil, expectedError)
		mock.ExpectRollback()

		err = receptionSvc.DeleteLastProduct(pvzID)

		require.Error(t, err)
		require.Equal(t, expectedError, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DeleteLastProduct error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockReceptionRepository(ctrl)
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		receptionSvc := service.NewReceptionService(mockRepo, db)

		pvzID := uuid.New()
		receptionID := uuid.New()
		expectedError := errors.New("repo error")

		mock.ExpectBegin()
		mockRepo.EXPECT().GetOpenedReceptionId(pvzID).Return(receptionID, nil)
		mockRepo.EXPECT().DeleteLastProduct(receptionID).Return(expectedError)
		mock.ExpectRollback()

		err = receptionSvc.DeleteLastProduct(pvzID)

		require.Error(t, err)
		require.Equal(t, expectedError, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestReceptionService_CloseLastReception(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockReceptionRepository(ctrl)
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		receptionSvc := service.NewReceptionService(mockRepo, db)

		pvzID := uuid.New()
		receptionID := uuid.New()
		returnedReception := &entity.Reception{
			Id:    receptionID,
			PvzId: pvzID,
		}

		mock.ExpectBegin()
		mockRepo.EXPECT().GetOpenedReceptionId(pvzID).Return(receptionID, nil)
		mockRepo.EXPECT().CloseLastReception(receptionID).Return(returnedReception, nil)
		mock.ExpectCommit()

		result, err := receptionSvc.CloseLastReception(pvzID)

		require.NoError(t, err)
		require.Equal(t, returnedReception, result)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Reception already closed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockReceptionRepository(ctrl)
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		receptionSvc := service.NewReceptionService(mockRepo, db)

		pvzID := uuid.New()

		mock.ExpectBegin()
		mockRepo.EXPECT().GetOpenedReceptionId(pvzID).Return(uuid.Nil, nil)
		mock.ExpectRollback()

		result, err := receptionSvc.CloseLastReception(pvzID)

		require.Error(t, err)
		require.Equal(t, service.ReceptionAlreadyClosed, err)
		require.Nil(t, result)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Transaction begin error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockReceptionRepository(ctrl)
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		receptionSvc := service.NewReceptionService(mockRepo, db)

		expectedError := errors.New("db error")
		mock.ExpectBegin().WillReturnError(expectedError)

		result, err := receptionSvc.CloseLastReception(uuid.New())

		require.Error(t, err)
		require.Equal(t, expectedError, err)
		require.Nil(t, result)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GetOpenedReceptionId error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockReceptionRepository(ctrl)
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		receptionSvc := service.NewReceptionService(mockRepo, db)

		pvzID := uuid.New()
		expectedError := errors.New("repo error")

		mock.ExpectBegin()
		mockRepo.EXPECT().GetOpenedReceptionId(pvzID).Return(uuid.Nil, expectedError)
		mock.ExpectRollback()

		result, err := receptionSvc.CloseLastReception(pvzID)

		require.Error(t, err)
		require.Equal(t, expectedError, err)
		require.Nil(t, result)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("CloseLastReception error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockReceptionRepository(ctrl)
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		receptionSvc := service.NewReceptionService(mockRepo, db)

		pvzID := uuid.New()
		receptionID := uuid.New()
		expectedError := errors.New("repo error")

		mock.ExpectBegin()
		mockRepo.EXPECT().GetOpenedReceptionId(pvzID).Return(receptionID, nil)
		mockRepo.EXPECT().CloseLastReception(receptionID).Return(nil, expectedError)
		mock.ExpectRollback()

		result, err := receptionSvc.CloseLastReception(pvzID)

		require.Error(t, err)
		require.Equal(t, expectedError, err)
		require.Nil(t, result)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

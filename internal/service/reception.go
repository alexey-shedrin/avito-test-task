package service

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/alexey-shedrin/avito-test-task/internal/metrics"
	"github.com/alexey-shedrin/avito-test-task/internal/model/entity"
	"github.com/google/uuid"
)

var (
	ReceptionAlreadyOpened = errors.New("reception is already opened")
	ReceptionNotOpened     = errors.New("reception is not opened")
	ReceptionAlreadyClosed = errors.New("reception is already closed")
)

type ReceptionRepository interface {
	GetOpenedReceptionId(pvzID uuid.UUID) (uuid.UUID, error)
	CreateReception(reception *entity.Reception) (*entity.Reception, error)
	CreateProduct(product *entity.Product) (*entity.Product, error)
	DeleteLastProduct(pvzID uuid.UUID) error
	CloseLastReception(receptionId uuid.UUID) (*entity.Reception, error)
}

type ReceptionService struct {
	receptionRepo ReceptionRepository
	db            *sql.DB
}

func NewReceptionService(receptionRepo ReceptionRepository, db *sql.DB) *ReceptionService {
	return &ReceptionService{
		receptionRepo: receptionRepo,
		db:            db,
	}
}

func (s *ReceptionService) CreateReception(reception *entity.Reception) (*entity.Reception, error) {
	log.SetPrefix("ReceptionService.CreateReception")

	tx, err := s.db.BeginTx(context.Background(), nil)
	if err != nil {
		log.Printf("error start transaction: %v", err)

		return nil, err
	}
	defer tx.Rollback()

	id, err := s.receptionRepo.GetOpenedReceptionId(reception.PvzId)
	if err != nil {
		return nil, err
	}

	if id != uuid.Nil {
		return nil, ReceptionAlreadyOpened
	}

	reception, err = s.receptionRepo.CreateReception(reception)
	if err != nil {
		return nil, err
	}

	tx.Commit()

	metrics.CreateReception()

	return reception, nil
}

func (s *ReceptionService) CreateProduct(product *entity.Product, pvzID uuid.UUID) (*entity.Product, error) {
	log.SetPrefix("ReceptionService.CreateProduct")

	tx, err := s.db.BeginTx(context.Background(), nil)
	if err != nil {
		log.Printf("error start transaction: %v", err)

		return nil, err
	}

	defer tx.Rollback()

	id, err := s.receptionRepo.GetOpenedReceptionId(pvzID)
	if err != nil {
		return nil, err
	}

	if id == uuid.Nil {
		return nil, ReceptionNotOpened
	}

	product.ReceptionId = id

	product, err = s.receptionRepo.CreateProduct(product)
	if err != nil {
		return nil, err
	}

	tx.Commit()

	metrics.AddProduct()

	return product, nil
}

func (s *ReceptionService) DeleteLastProduct(pvzID uuid.UUID) error {
	log.SetPrefix("ReceptionService.DeleteLastProduct")

	tx, err := s.db.BeginTx(context.Background(), nil)
	if err != nil {
		log.Printf("error start transaction: %v", err)

		return err
	}

	defer tx.Rollback()

	id, err := s.receptionRepo.GetOpenedReceptionId(pvzID)
	if err != nil {
		return err
	}

	if id == uuid.Nil {
		return ReceptionNotOpened
	}

	if err = s.receptionRepo.DeleteLastProduct(id); err != nil {
		return err
	}

	tx.Commit()

	metrics.DeleteProduct()

	return nil
}

func (s *ReceptionService) CloseLastReception(pvzID uuid.UUID) (*entity.Reception, error) {
	log.SetPrefix("ReceptionService.CloseLastReception")

	tx, err := s.db.BeginTx(context.Background(), nil)
	if err != nil {
		log.Printf("error start transaction: %v", err)

		return nil, err
	}

	defer tx.Rollback()

	id, err := s.receptionRepo.GetOpenedReceptionId(pvzID)
	if err != nil {
		return nil, err
	}

	if id == uuid.Nil {
		return nil, ReceptionAlreadyClosed
	}

	reception, err := s.receptionRepo.CloseLastReception(id)
	if err != nil {
		return nil, err
	}

	tx.Commit()

	return reception, nil
}

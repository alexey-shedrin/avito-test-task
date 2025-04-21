package repository

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/alexey-shedrin/avito-test-task/internal/model/entity"
	"github.com/google/uuid"
)

var errProductNotFound = errors.New("product not found")

type ReceptionRepository struct {
	db *sql.DB
}

func NewReceptionRepository(db *sql.DB) *ReceptionRepository {
	return &ReceptionRepository{
		db: db,
	}
}
func (r *ReceptionRepository) GetOpenedReceptionId(pvzID uuid.UUID) (uuid.UUID, error) {
	log.SetPrefix("repository.CheckOpenedReception")
	query := `SELECT id FROM reception WHERE pvz_id = $1 AND status = 'in_progress'`

	id := uuid.UUID{}
	if err := r.db.QueryRow(query, pvzID).Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uuid.UUID{}, nil
		}

		log.Printf("error: %v", err)

		return uuid.UUID{}, err
	}

	return id, nil
}

func (r *ReceptionRepository) CreateReception(reception *entity.Reception) (*entity.Reception, error) {
	log.SetPrefix("repository.CreateReception")
	query := `INSERT INTO reception (id, reception_datetime, pvz_id, status) VALUES ($1, $2, $3, $4)`

	reception.DateTime = time.Now()
	reception.Status = "in_progress"
	reception.Id = uuid.New()

	if _, err := r.db.Exec(query, reception.Id, reception.DateTime, reception.PvzId, reception.Status); err != nil {
		log.Printf("error: %v", err)

		return nil, err
	}

	return reception, nil
}

func (r *ReceptionRepository) CreateProduct(product *entity.Product) (*entity.Product, error) {
	log.SetPrefix("repository.CreateProduct")
	query := `INSERT INTO product (id, product_type, acceptance_datetime, reception_id) VALUES ($1, $2, $3, $4)`

	product.DateTime = time.Now()
	product.Id = uuid.New()

	if _, err := r.db.Exec(query, product.Id, product.Type, product.DateTime, product.ReceptionId); err != nil {
		log.Printf("error: %v", err)

		return nil, err
	}

	return product, nil
}

func (r *ReceptionRepository) DeleteLastProduct(receptionID uuid.UUID) error {
	log.SetPrefix("repository.DeleteLastProduct")
	query := `DELETE FROM product WHERE id = (SELECT id FROM product WHERE reception_id = $1 ORDER BY acceptance_datetime DESC LIMIT 1)`
	if _, err := r.db.Exec(query, receptionID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errProductNotFound
		}

		log.Printf("error: %v", err)

		return err
	}
	return nil
}

func (r *ReceptionRepository) CloseLastReception(receptionID uuid.UUID) (*entity.Reception, error) {
	log.SetPrefix("repository.CloseLastReception")
	query := `UPDATE reception SET status = 'closed' WHERE id = $1`
	if _, err := r.db.Exec(query, receptionID); err != nil {
		log.Printf("error: %v", err)

		return nil, err
	}

	return &entity.Reception{
		Id:       receptionID,
		Status:   "closed",
		DateTime: time.Now(),
	}, nil
}

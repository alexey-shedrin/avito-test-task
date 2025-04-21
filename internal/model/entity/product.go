package entity

import (
	"time"

	"github.com/alexey-shedrin/avito-test-task/internal/model/dto/response"
	"github.com/google/uuid"
)

type Product struct {
	Id          uuid.UUID
	ReceptionId uuid.UUID
	Type        string
	DateTime    time.Time
}

func (p *Product) ToResponse() *response.Product {
	return &response.Product{
		Id:          p.Id,
		ReceptionId: p.ReceptionId,
		Type:        p.Type,
		DateTime:    p.DateTime,
	}
}

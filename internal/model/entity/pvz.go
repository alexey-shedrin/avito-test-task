package entity

import (
	"time"

	"github.com/alexey-shedrin/avito-test-task/internal/model/dto/response"
	"github.com/google/uuid"
)

type Pvz struct {
	Id               uuid.UUID
	City             string
	RegistrationDate time.Time
}

func (p *Pvz) ToResponse() response.Pvz {
	return response.Pvz{
		Id:               p.Id,
		City:             p.City,
		RegistrationDate: p.RegistrationDate,
	}
}

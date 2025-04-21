package entity

import (
	"time"

	"github.com/alexey-shedrin/avito-test-task/internal/model/dto/response"
	"github.com/google/uuid"
)

type Reception struct {
	Id       uuid.UUID
	PvzId    uuid.UUID
	Status   string
	DateTime time.Time
}

func (r *Reception) ToResponse() response.Reception {
	return response.Reception{
		Id:       r.Id,
		PvzId:    r.PvzId,
		Status:   r.Status,
		DateTime: r.DateTime,
	}
}

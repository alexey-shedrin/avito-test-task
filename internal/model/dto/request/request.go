package request

import (
	"github.com/google/uuid"
)

type DummyLogin struct {
	Role string `json:"role" binding:"required,oneof=employee moderator"`
}

type Register struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required,oneof=client employee moderator"`
}

type Login struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type PvzId struct {
	PvzId uuid.UUID `json:"pvzId" binding:"required"`
}

type Pvz struct {
	City string `json:"city" binding:"required,oneof=Москва Санкт-Петербург Казань"`
}

type Reception struct {
	PvzId uuid.UUID `json:"pvzId" binding:"required"`
}

type CreateProduct struct {
	PvzId uuid.UUID `json:"pvzId" binding:"required"`
	Type  string    `json:"type" binding:"required,oneof=электроника одежда обувь"`
}

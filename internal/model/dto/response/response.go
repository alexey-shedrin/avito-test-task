package response

import (
	"github.com/google/uuid"
	"time"
)

type DummyLogin struct {
	Token string `json:"token"`
}

type User struct {
	Id    uuid.UUID `json:"id"`
	Email string    `json:"email"`
	Role  string    `json:"role"`
}

type Login struct {
	Token string `json:"token"`
}

type Pvz struct {
	Id               uuid.UUID `json:"id"`
	City             string    `json:"city"`
	RegistrationDate time.Time `json:"registrationDate"`
}

type Reception struct {
	Id       uuid.UUID `json:"id"`
	PvzId    uuid.UUID `json:"pvzId"`
	Status   string    `json:"status"`
	DateTime time.Time `json:"dateTime"`
}

type Product struct {
	Id          uuid.UUID `json:"id"`
	ReceptionId uuid.UUID `json:"receptionId"`
	Type        string    `json:"type"`
	DateTime    time.Time `json:"dateTime"`
}

type ReceptionsWithProducts struct {
	Products  []Product `json:"products"`
	Reception Reception `json:"reception"`
}

type PvzInfo struct {
	Pvz        Pvz                      `json:"pvz"`
	Receptions []ReceptionsWithProducts `json:"receptions"`
}

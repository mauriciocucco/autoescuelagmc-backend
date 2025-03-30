package models

import (
	"time"
)

// ServiceType representa los diferentes servicios que ofrece la autoescuela
type ServiceType string

const (
    ServiceIntensive    ServiceType = "clases_intensivas"
    ServiceExamPrep     ServiceType = "preparacion_examen"
    ServiceRegularClass ServiceType = "clases_regulares"
    ServiceOther        ServiceType = "otros"
)

// ContactRequest representa una solicitud de contacto del formulario
type ContactRequest struct {
    ID           int         `json:"id" db:"id"`
    FullName     string      `json:"full_name" db:"full_name"`
    Email        string      `json:"email" db:"email"`
    PhoneNumber  string      `json:"phone_number" db:"phone_number"`
    ServiceType  ServiceType `json:"service_type" db:"service_type"`
    Message      string      `json:"message" db:"message"`
    CreatedAt    time.Time   `json:"created_at" db:"created_at"`
}
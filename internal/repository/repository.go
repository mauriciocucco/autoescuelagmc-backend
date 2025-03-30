package repository

import (
	"context"
	"database/sql"
	"time"

	"autoescuelagmc-backend/internal/models"
)

// ContactRepository maneja las operaciones de base de datos para los contactos
type ContactRepository struct {
    db *sql.DB
}

// NewContactRepository crea un nuevo repositorio de contactos
func NewContactRepository(db *sql.DB) *ContactRepository {
    return &ContactRepository{
        db: db,
    }
}

// SaveContact guarda una solicitud de contacto en la base de datos
func (r *ContactRepository) SaveContact(ctx context.Context, contact models.ContactRequest) (int, error) {
    query := `
        INSERT INTO contact_requests (full_name, email, phone_number, service_type, message, created_at)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id
    `
    
    var id int
    err := r.db.QueryRowContext(
        ctx,
        query,
        contact.FullName,
        contact.Email,
        contact.PhoneNumber,
        contact.ServiceType,
        contact.Message,
        time.Now(),
    ).Scan(&id)
    
    return id, err
}
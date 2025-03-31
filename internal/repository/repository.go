package repository

import (
	"context"
	"database/sql"
	"log"
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
        VALUES (?, ?, ?, ?, ?, ?)
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

func (r *ContactRepository) GetAllContacts(ctx context.Context) ([]models.ContactRequest, error) {
    query := `SELECT id, full_name, email, phone_number, service_type, message, created_at FROM contact_requests ORDER BY created_at DESC`
    
    rows, err := r.db.QueryContext(ctx, query)

    if err != nil {
        log.Printf("ERROR en GetAllContacts: %+v", err)

        return nil, err
    }

    defer rows.Close()
    
    var contacts []models.ContactRequest
    
    for rows.Next() {
        var contact models.ContactRequest
        
        if err := rows.Scan(
            &contact.ID,
            &contact.FullName,
            &contact.Email,
            &contact.PhoneNumber,
            &contact.ServiceType,
            &contact.Message,
            &contact.CreatedAt,
        ); err != nil {
            return nil, err
        }
        
        contacts = append(contacts, contact)
    }
    
    if err := rows.Err(); err != nil {
        return nil, err
    }
    
    return contacts, nil
}
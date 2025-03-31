package database

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

type Config struct {
    // PostgreSQL config (keeping for backward compatibility)
    Host     string
    Port     int
    User     string
    Password string
    DBName   string
    SSLMode  string
    
    // Turso config
    TursoEnabled  bool
    TursoDbPath   string
    TursoPrimaryUrl string
    TursoAuthToken  string
}

func NewDB(cfg Config) (*sql.DB, error) {
    if cfg.TursoEnabled {
        return newTursoDB(cfg)
    }

    return NewPostgresDB(cfg)
}

func newTursoDB(cfg Config) (*sql.DB, error) {
    // Crear la URL de conexión para Turso
    url := fmt.Sprintf("%s?authToken=%s", cfg.TursoPrimaryUrl, cfg.TursoAuthToken)
    
    // Asegurar que la URL comienza con libsql:// o https://
    if !strings.HasPrefix(url, "libsql://") && !strings.HasPrefix(url, "https://") {
        url = "libsql://" + url
    }
    
    // Abrir conexión usando el driver libsql
    db, err := sql.Open("libsql", url)
    
    if err != nil {
        return nil, fmt.Errorf("error opening libsql connection: %w", err)
    }
    
    // Verificar conexión
    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("error connecting to Turso database: %w", err)
    }
    
    return db, nil
}

func NewPostgresDB(cfg Config) (*sql.DB, error) {
    dsn := fmt.Sprintf(
        "host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
        cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
    )
    
    db, err := sql.Open("postgres", dsn)

    if err != nil {
        return nil, fmt.Errorf("error opening database connection: %w", err)
    }
    
    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("error connecting to database: %w", err)
    }
    
    return db, nil
}

func CreateTablesIfNotExist(db *sql.DB) error {
    query := `
        CREATE TABLE IF NOT EXISTS contact_requests (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            full_name TEXT NOT NULL,
            email TEXT NOT NULL,
            phone_number TEXT NOT NULL,
            service_type TEXT NOT NULL,
            message TEXT,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        );
    `
    
    _, err := db.Exec(query)
    
    return err
}
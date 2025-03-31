package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/joho/godotenv"

	"autoescuelagmc-backend/internal/app"
	"autoescuelagmc-backend/internal/database"
	"autoescuelagmc-backend/internal/services"
)

func main() {
    fmt.Println("Iniciando aplicación...")
    // En una implementación real, estas configuraciones vendrían de variables de entorno
    // o de un archivo de configuración

    // Cargar variables de entorno desde .env
    if err := godotenv.Load(); err != nil {
        log.Println("No se pudo cargar el archivo .env:", err)
    }
    
    // Configuración de la base de datos
    dbConfig := database.Config{
        // Set Turso as enabled
        TursoEnabled:   true,
        TursoDbPath:    getEnv("TURSO_DB_PATH", "./data/local.db"),
        TursoPrimaryUrl: getEnv("TURSO_PRIMARY_URL", "https://autoescuelagmc-maurisc.aws-us-east-1.turso.io"),
        TursoAuthToken:  getEnv("TURSO_AUTH_TOKEN", "your-auth-token"),
        
        // Keep PostgreSQL config for backward compatibility
        Host:     getEnv("DB_HOST", "localhost"),
        Port:     getEnvAsInt("DB_PORT", 5432),
        User:     getEnv("DB_USER", "postgres"),
        Password: getEnv("DB_PASSWORD", "postgres"),
        DBName:   getEnv("DB_NAME", "autoescuela"),
        SSLMode:  getEnv("DB_SSL_MODE", "disable"),
    }
    
    // Configuración del email
    emailConfig := services.EmailConfig{
        SMTPHost:     getEnv("SMTP_HOST", "smtp.gmail.com"),
        SMTPPort:     getEnvAsInt("SMTP_PORT", 587),
        SMTPUsername: getEnv("SMTP_USERNAME", "tu_correo@gmail.com"),
        SMTPPassword: getEnv("SMTP_PASSWORD", "tu_contraseña"),
        FromEmail:    getEnv("FROM_EMAIL", "noreply@autoescuelagmc.com"),
        ToEmail:      getEnv("TO_EMAIL", "autoescuelagmc@hotmail.com"),
    }
    
    // Configuración de la aplicación
    appConfig := app.Config{
        ServerPort: getEnv("SERVER_PORT", "8080"),
        DB:         dbConfig,
        Email:      emailConfig,
    }
    
    // Crear e iniciar la aplicación
    application := app.NewApp(appConfig)
    
    if err := application.Start(); err != nil {
        log.Fatalf("Error starting the application: %v", err)
    }
    
    // Configurar canal para señales de sistema operativo
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    
    // Esperar señal de terminación
    <-quit
    log.Println("Shutting down server...")
    
    // Detener la aplicación de forma controlada
    if err := application.Stop(); err != nil {
        log.Fatalf("Error stopping the application: %v", err)
    }
    
    log.Println("Server exited properly")
}

// getEnv obtiene una variable de entorno o devuelve un valor por defecto
func getEnv(key, defaultValue string) string {
    value := os.Getenv(key)

    if value == "" {
        return defaultValue
    }

    return value
}

// getEnvAsInt obtiene una variable de entorno como entero
func getEnvAsInt(key string, defaultValue int) int {
    valueStr := os.Getenv(key)

    if valueStr == "" {
        return defaultValue
    }
    
    value, err := strconv.Atoi(valueStr)
    
    if err != nil {
        return defaultValue
    }
    
    return value
}
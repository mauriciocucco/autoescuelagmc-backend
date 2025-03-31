package app

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"autoescuelagmc-backend/internal/database"
	"autoescuelagmc-backend/internal/handlers"
	"autoescuelagmc-backend/internal/repository"
	"autoescuelagmc-backend/internal/services"
)

// Config contiene la configuración de la aplicación
type Config struct {
    ServerPort  string
    DB          database.Config
    Email       services.EmailConfig
}

// App encapsula la funcionalidad de la aplicación
type App struct {
    config     Config
    router     *gin.Engine
    db         *sql.DB
    httpServer *http.Server
}

// NewApp crea una nueva instancia de la aplicación
func NewApp(config Config) *App {
    return &App{
        config: config,
        router: gin.Default(),
    }
}

// Start inicia la aplicación
func (a *App) Start() error {
    // Configurar base de datos
    db, err := database.NewDB(a.config.DB)
    
    if err != nil {
        return err
    }
    a.db = db
    
    // Crear tablas si no existen
    // if err := database.CreateTablesIfNotExist(a.db); err != nil {
    //     return err
    // }
    
    // Inicializar repositorios
    contactRepo := repository.NewContactRepository(a.db)
    
    // Inicializar servicios
    emailService := services.NewEmailService(a.config.Email)
    
    // Inicializar handlers
    contactHandler := handlers.NewContactHandler(contactRepo, emailService)
    
    // Configurar CORS
    a.router.Use(func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(http.StatusOK)
            return
        }
        
        c.Next()
    })
    
    // Configurar rutas
    api := a.router.Group("/api")
    {
        api.POST("/contact", contactHandler.CreateContact)
        api.GET("/contacts", contactHandler.GetAllContacts)
    }
    
    // Iniciar servidor HTTP
    a.httpServer = &http.Server{
        Addr:    ":" + a.config.ServerPort,
        Handler: a.router,
    }
    
    go func() {
        if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Error starting server: %v", err)
        }
    }()
    
    log.Printf("Server started on port %s", a.config.ServerPort)
    
    return nil
}

// Stop detiene la aplicación
func (a *App) Stop() error {
    // Cerrar servidor HTTP
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    if err := a.httpServer.Shutdown(ctx); err != nil {
        return err
    }
    
    // Cerrar conexión de base de datos
    if a.db != nil {
        if err := a.db.Close(); err != nil {
            return err
        }
    }
    
    return nil
}
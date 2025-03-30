package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"autoescuelagmc-backend/internal/models"
	"autoescuelagmc-backend/internal/repository"
	"autoescuelagmc-backend/internal/services"
)

// ContactHandler maneja las solicitudes HTTP relacionadas con contactos
type ContactHandler struct {
    contactRepo *repository.ContactRepository
    emailService *services.EmailService
}

// NewContactHandler crea un nuevo manejador de contactos
func NewContactHandler(contactRepo *repository.ContactRepository, emailService *services.EmailService) *ContactHandler {
    return &ContactHandler{
        contactRepo: contactRepo,
        emailService: emailService,
    }
}

// CreateContact maneja la creación de nuevas solicitudes de contacto
func (h *ContactHandler) CreateContact(c *gin.Context) {
    var contact models.ContactRequest
    if err := c.ShouldBindJSON(&contact); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de formulario inválidos"})
        return
    }
    
    // Validar datos
    if contact.FullName == "" || contact.Email == "" || contact.PhoneNumber == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Nombre, email y teléfono son obligatorios"})
        return
    }
    
    // Guardar en base de datos
    id, err := h.contactRepo.SaveContact(c.Request.Context(), contact)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar el contacto"})
        return
    }
    
    // Enviar email
    err = h.emailService.SendContactNotification(contact)
    if err != nil {
        // Logear el error pero continuar (no fallar la solicitud)
        // En una implementación real, podrías añadir una cola para reintentar
    }
    
    c.JSON(http.StatusCreated, gin.H{
        "id": id,
        "message": "Solicitud de contacto recibida correctamente",
    })
}
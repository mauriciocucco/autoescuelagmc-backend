package services

import (
	"autoescuelagmc-backend/internal/models"
	"fmt"

	"gopkg.in/mail.v2"
)

// EmailConfig contiene la configuración para el servicio de correo
type EmailConfig struct {
    SMTPHost     string
    SMTPPort     int
    SMTPUsername string
    SMTPPassword string
    FromEmail    string
    ToEmail      string
}

// EmailService maneja el envío de correos electrónicos
type EmailService struct {
    config EmailConfig
}

// NewEmailService crea un nuevo servicio de correo
func NewEmailService(config EmailConfig) *EmailService {
    return &EmailService{
        config: config,
    }
}

// SendContactNotification envía una notificación por correo sobre un nuevo contacto
func (s *EmailService) SendContactNotification(contact models.ContactRequest) error {
    m := mail.NewMessage()
    m.SetHeader("From", s.config.FromEmail)
    m.SetHeader("To", s.config.ToEmail)
    m.SetHeader("Subject", "Nueva solicitud de contacto - Autoescuela GMC")
    
    body := fmt.Sprintf(`
        <h2>Nueva solicitud de contacto</h2>
        <p><strong>Nombre:</strong> %s</p>
        <p><strong>Email:</strong> %s</p>
        <p><strong>Teléfono:</strong> %s</p>
        <p><strong>Servicio de interés:</strong> %s</p>
        <p><strong>Mensaje:</strong></p>
        <p>%s</p>
    `, contact.FullName, contact.Email, contact.PhoneNumber, contact.ServiceType, contact.Message)
    
    m.SetBody("text/html", body)
    
    d := mail.NewDialer(s.config.SMTPHost, s.config.SMTPPort, s.config.SMTPUsername, s.config.SMTPPassword)
    
    return d.DialAndSend(m)
}
// utils/email.go
package utils

import (
	"os"

	"gopkg.in/gomail.v2"
)

type EmailConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

// NewEmailConfig returns a new EmailConfig after environment variables are loaded.
func NewEmailConfig() EmailConfig {
	return EmailConfig{
		Host:     "smtp.gmail.com",
		Port:     587,
		Username: os.Getenv("EMAIL"),
		Password: os.Getenv("APP_PASSWORD"),
	}
}

// SendPasswordResetEmail sends a password reset email to the specified recipient
func SendPasswordResetEmail(recipient, resetLink string) error {
	config := NewEmailConfig()
	// Configure email settings - ideally these should be in environment variables
	host := config.Host
	port := config.Port
	username := config.Username
	password := config.Password // Make sure to set this environment variable

	// Create a new email message
	m := gomail.NewMessage()
	m.SetHeader("From", username)
	m.SetHeader("To", recipient)
	m.SetHeader("Subject", "Mačva Press - Resetovanje Lozinke")

	// HTML body with a button that links to the reset page
	htmlBody := `
	<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
		<img src="https://macva-press.duckdns.org/static/assets/macva-1-300x71.png" alt="Mačva Press Logo" style="width: 200px; height: auto; margin-bottom: 20px;" />
		<h2>Zahtev za resetovanje lozinke</h2>
		<p>Primili smo zahtev za resetovanje lozinke za vaš nalog. Kliknite na dugme ispod da biste resetovali lozinku:</p>
		<p style="margin: 30px 0;">
			<a href="` + resetLink + `" style="background-color: #3B82F6; color: white; padding: 12px 20px; text-decoration: none; border-radius: 5px; font-weight: bold;">Resetuj Lozinku</a>
		</p>
		<p>Ako niste zatražili resetovanje lozinke, molimo vas da ignorišete ovu poruku.</p>
		<p>Ovaj link će isteći za 1 sat.</p>
		<hr style="margin: 30px 0; border: none; border-top: 1px solid #eaeaea;" />
		<p style="font-size: 12px; color: #666;">Mačva Press Tim</p>
	</div>
	`
	m.SetBody("text/html", htmlBody)

	// Set up the email dialer
	d := gomail.NewDialer(host, port, username, password)

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

// SendEmailVerificationEmail sends a verification email to the specified recipient
func SendEmailVerificationEmail(recipient, verificationLink string) error {
	config := NewEmailConfig()
	// Configure email settings using the config
	host := config.Host
	port := config.Port
	username := config.Username
	password := config.Password

	// Create a new email message
	m := gomail.NewMessage()
	m.SetHeader("From", username)
	m.SetHeader("To", recipient)
	m.SetHeader("Subject", "Mačva Press - Verifikacija Email Adrese")

	// HTML body with a button that links to the verification page
	htmlBody := `
	<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
		<img src="https://macva-press.duckdns.org/static/assets/macva-1-300x71.png" alt="Mačva Press Logo" style="width: 200px; height: auto; margin-bottom: 20px;" />
		<h2>Verifikacija email adrese</h2>
		<p>Hvala što ste se registrovali na Mačva Press portal. Da biste aktivirali svoj nalog, kliknite na dugme ispod:</p>
		<p style="margin: 30px 0;">
			<a href="` + verificationLink + `" style="background-color: #3B82F6; color: white; padding: 12px 20px; text-decoration: none; border-radius: 5px; font-weight: bold;">Verifikuj Email</a>
		</p>
		<p>Ako niste kreirali nalog na našem portalu, molimo vas da ignorišete ovu poruku.</p>
		<p>Ovaj link će isteći za 24 sata.</p>
		<hr style="margin: 30px 0; border: none; border-top: 1px solid #eaeaea;" />
		<p style="font-size: 12px; color: #666;">Mačva Press Tim</p>
	</div>
	`
	m.SetBody("text/html", htmlBody)

	// Set up the email dialer
	d := gomail.NewDialer(host, port, username, password)

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

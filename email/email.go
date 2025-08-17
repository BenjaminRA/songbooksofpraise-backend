package email

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strconv"
)

// SendEmail sends an email using Zoho SMTP with TLS encryption
// Supports both port 587 (TLS) and port 465 (SSL)
// For Zoho SMTP configuration:
// - Port 587 with TLS (recommended)
// - Port 465 with SSL (alternative)
// - Requires authentication
func SendEmail(to string, subject string, msg string) error {
	// log.Printf("Starting email send process...")
	// log.Printf("Recipient: %s", to)
	// log.Printf("Subject: %s", subject)

	from := os.Getenv("MAIL_USERNAME")
	password := os.Getenv("MAIL_PASSWORD")
	host := os.Getenv("MAIL_HOST")
	port := os.Getenv("MAIL_PORT")

	// log.Printf("SMTP Configuration:")
	// log.Printf("  Host: %s", host)
	// log.Printf("  Port: %s", port)
	// log.Printf("  Username: %s", from)
	// log.Printf("  Password: %s", maskPassword(password))

	if port == "" {
		port = "587" // Default to 587 for TLS
		log.Printf("No port specified, defaulting to 587")
	}

	// Convert port to int for validation
	portInt, err := strconv.Atoi(port)
	if err != nil {
		// log.Printf("ERROR: Invalid port number: %s", port)
		return fmt.Errorf("invalid port number: %s", port)
	}

	// Format the message with proper headers
	message := fmt.Sprintf("From: Songbooks of Praise <%s>\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s", from, to, subject, msg)
	// log.Printf("Message formatted successfully")

	// SMTP server configuration
	smtpHost := host
	smtpPort := fmt.Sprintf("%s:%d", smtpHost, portInt)
	// log.Printf("Connecting to SMTP server: %s", smtpPort)

	var client *smtp.Client

	if portInt == 465 {
		// Port 465 uses direct TLS/SSL connection
		// log.Printf("Using SSL connection (port 465)...")
		tlsConfig := &tls.Config{
			InsecureSkipVerify: false,
			ServerName:         smtpHost,
		}

		conn, err := tls.Dial("tcp", smtpPort, tlsConfig)
		if err != nil {
			// log.Printf("ERROR: Failed to connect to SMTP server with SSL: %v", err)
			return fmt.Errorf("failed to connect to SMTP server with SSL: %v", err)
		}
		defer conn.Close()
		// log.Printf("SSL connection established successfully")

		client, err = smtp.NewClient(conn, smtpHost)
		if err != nil {
			// log.Printf("ERROR: Failed to create SMTP client: %v", err)
			return fmt.Errorf("failed to create SMTP client: %v", err)
		}
	} else {
		// Port 587 and others use STARTTLS
		log.Printf("Using STARTTLS connection (port %d)...", portInt)

		// Connect without TLS first
		client, err = smtp.Dial(smtpPort)
		if err != nil {
			// log.Printf("ERROR: Failed to connect to SMTP server: %v", err)
			return fmt.Errorf("failed to connect to SMTP server: %v", err)
		}
		// log.Printf("Initial connection established successfully")

		// Start TLS
		tlsConfig := &tls.Config{
			InsecureSkipVerify: false,
			ServerName:         smtpHost,
		}

		if err = client.StartTLS(tlsConfig); err != nil {
			client.Quit()
			// log.Printf("ERROR: Failed to start TLS: %v", err)
			return fmt.Errorf("failed to start TLS: %v", err)
		}
		// log.Printf("STARTTLS established successfully")
	}
	defer client.Quit()
	// log.Printf("SMTP client ready")

	// Authenticate
	// log.Printf("Authenticating with SMTP server...")
	auth := smtp.PlainAuth("", from, password, smtpHost)
	if err = client.Auth(auth); err != nil {
		// log.Printf("ERROR: Failed to authenticate: %v", err)
		return fmt.Errorf("failed to authenticate: %v", err)
	}
	// log.Printf("Authentication successful")

	// Set sender
	// log.Printf("Setting sender: %s", from)
	if err = client.Mail(from); err != nil {
		// log.Printf("ERROR: Failed to set sender: %v", err)
		return fmt.Errorf("failed to set sender: %v", err)
	}
	// log.Printf("Sender set successfully")

	// Set recipient
	// log.Printf("Setting recipient: %s", to)
	if err = client.Rcpt(to); err != nil {
		// log.Printf("ERROR: Failed to set recipient: %v", err)
		return fmt.Errorf("failed to set recipient: %v", err)
	}
	// log.Printf("Recipient set successfully")

	// Send the email body
	// log.Printf("Sending email body...")
	writer, err := client.Data()
	if err != nil {
		// log.Printf("ERROR: Failed to get data writer: %v", err)
		return fmt.Errorf("failed to get data writer: %v", err)
	}

	_, err = writer.Write([]byte(message))
	if err != nil {
		// log.Printf("ERROR: Failed to write message: %v", err)
		return fmt.Errorf("failed to write message: %v", err)
	}

	err = writer.Close()
	if err != nil {
		// log.Printf("ERROR: Failed to close writer: %v", err)
		return fmt.Errorf("failed to close writer: %v", err)
	}

	log.Printf("Email sent successfully to %s", to)
	return nil
}

// maskPassword masks the password for logging purposes
// func maskPassword(password string) string {
// 	if len(password) == 0 {
// 		return "[EMPTY]"
// 	}
// 	if len(password) <= 4 {
// 		return "****"
// 	}
// 	return password[:2] + strings.Repeat("*", len(password)-4) + password[len(password)-2:]
// }

// // Test program to debug email sending
// // Run with: go run debug/test_email.go
package main

// import (
// 	"fmt"
// 	"log"
// 	"os"
// 	"strings"

// 	"github.com/BenjaminRA/himnario-backend/email"
// 	"github.com/joho/godotenv"
// )

// func main() {
// 	// Load environment variables
// 	err := godotenv.Load("../.env")
// 	if err != nil {
// 		log.Fatal("Error loading .env file:", err)
// 	}

// 	// Test email configuration
// 	fmt.Println("Testing email configuration...")
// 	fmt.Printf("MAIL_HOST: %s\n", os.Getenv("MAIL_HOST"))
// 	fmt.Printf("MAIL_PORT: %s\n", os.Getenv("MAIL_PORT"))
// 	fmt.Printf("MAIL_USERNAME: %s\n", os.Getenv("MAIL_USERNAME"))

// 	// Mask password for security
// 	password := os.Getenv("MAIL_PASSWORD")
// 	if len(password) > 4 {
// 		fmt.Printf("MAIL_PASSWORD: %s...%s (length: %d)\n", password[:2], password[len(password)-2:], len(password))
// 	} else {
// 		fmt.Printf("MAIL_PASSWORD: **** (length: %d)\n", len(password))
// 	}

// 	// Check if all required variables are set
// 	if os.Getenv("MAIL_HOST") == "" {
// 		log.Fatal("MAIL_HOST is not set")
// 	}
// 	if os.Getenv("MAIL_USERNAME") == "" {
// 		log.Fatal("MAIL_USERNAME is not set")
// 	}
// 	if os.Getenv("MAIL_PASSWORD") == "" {
// 		log.Fatal("MAIL_PASSWORD is not set")
// 	}

// 	// Test sending an email
// 	testEmail := "benjamin.gra572@gmail.com" // Change this to your test email
// 	subject := "Test Email from Songbooks of Praise"
// 	message := `
// 	<html>
// 	<body>
// 		<h2>Test Email</h2>
// 		<p>This is a test email to verify the SMTP configuration.</p>
// 		<p>If you received this email, the configuration is working correctly!</p>
// 		<br>
// 		<p>Best regards,<br>Songbooks of Praise Team</p>
// 	</body>
// 	</html>
// 	`

// 	fmt.Printf("\nSending test email to: %s\n", testEmail)
// 	fmt.Println("Check the logs below for detailed debugging information...")
// 	fmt.Println("=" + strings.Repeat("=", 60))

// 	err = email.SendEmail(testEmail, subject, message)
// 	if err != nil {
// 		log.Printf("Failed to send email: %v", err)
// 		os.Exit(1)
// 	}

// 	fmt.Println("=" + strings.Repeat("=", 60))
// 	fmt.Println("Email sent successfully! Check your inbox (and spam folder).")
// }

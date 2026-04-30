package support

import (
	"net/http"
	"strings"

	"github.com/BenjaminRA/himnario-backend/email"
	"github.com/gin-gonic/gin"
)

type contactRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Topic   string `json:"topic"`
	Message string `json:"message"`
	Lang    string `json:"lang"`
}

// PostSupportContact handles POST /support/contact.
// It is a public endpoint — no authentication is required.
func PostSupportContact(c *gin.Context) {
	var req contactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body."})
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(req.Email)
	req.Message = strings.TrimSpace(req.Message)

	if req.Name == "" || req.Email == "" || req.Message == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name, email, and message are required."})
		return
	}

	if err := email.SendSupportContactEmail(req.Name, req.Email, req.Topic, req.Message, req.Lang); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message. Please try again later."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Your message has been sent successfully."})
}

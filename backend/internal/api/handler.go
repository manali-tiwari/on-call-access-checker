package api

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/manali-tiwari/on-call-access-checker/backend/internal/auth"
)

const ProductionAccessDuration = 12 * time.Hour

type Handler struct {
	oktaAuth auth.OktaAuthenticator
	awsAuth  auth.AWSAuthenticator
}

type AccessCheckResponse struct {
	VPN            bool     `json:"vpn"`
	Production     bool     `json:"production"`
	ConfigTool     bool     `json:"configTool"`
	CurrentProfile string   `json:"currentProfile"`
	MissingGroups  []string `json:"missingGroups"`
	ValidUntil     string   `json:"validUntil,omitempty"`
	ProfileARN     string   `json:"profileArn,omitempty"`
}

type AccessCheckRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Environment string `json:"environment" binding:"required"`
}

func NewHandler(oktaAuth auth.OktaAuthenticator, awsAuth auth.AWSAuthenticator) *Handler {
	return &Handler{
		oktaAuth: oktaAuth,
		awsAuth:  awsAuth,
	}
}

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	router.POST("/api/check-access", h.checkAccess)
}

func (h *Handler) checkAccess(c *gin.Context) {
	// Log and read the raw request body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	log.Printf("Raw request body: %s", string(body))

	// Restore the body for binding
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	var request AccessCheckRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Printf("Request binding error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Processing request for email: %s, environment: %s", request.Email, request.Environment)

	// Check Okta access
	accessStatus, err := h.oktaAuth.CheckAccess(request.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Check AWS profile
	awsProfile, err := h.awsAuth.GetProfileInfo(request.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := AccessCheckResponse{
		VPN:            accessStatus.VPN,
		Production:     accessStatus.Production,
		ConfigTool:     accessStatus.ConfigTool,
		CurrentProfile: awsProfile.Name,
		ProfileARN:     awsProfile.ARN,
		MissingGroups:  []string{},
	}

	// Collect missing groups
	if !accessStatus.VPN {
		response.MissingGroups = append(response.MissingGroups, "VPN")
	}
	if !accessStatus.Production {
		response.MissingGroups = append(response.MissingGroups, "Production")
	}
	if !accessStatus.ConfigTool {
		response.MissingGroups = append(response.MissingGroups, "Config Tool")
	}

	// Set expiration if production access is granted
	if accessStatus.Production && request.Environment == "Production" {
		response.ValidUntil = time.Now().Add(ProductionAccessDuration).Format(time.RFC3339)
	}

	log.Printf("Returning response: %+v", response)
	c.JSON(http.StatusOK, response)
}

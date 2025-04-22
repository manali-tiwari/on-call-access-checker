package api

import (
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
	email := c.PostForm("email")
	environment := c.DefaultPostForm("environment", "Production")

	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email is required"})
		return
	}

	// Check Okta access
	accessStatus, err := h.oktaAuth.CheckAccess(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Check AWS profile (only if checking production environment)
	var awsProfile *auth.AWSProfile
	if environment == "Production" {
		awsProfile, err = h.awsAuth.GetProfileInfo()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		awsProfile = &auth.AWSProfile{Name: "dev", ARN: ""}
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
	if accessStatus.Production && environment == "Production" {
		response.ValidUntil = time.Now().Add(ProductionAccessDuration).Format(time.RFC3339)
	}

	c.JSON(http.StatusOK, response)
}

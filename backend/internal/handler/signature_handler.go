package handler

import (
	"net/http"

	"github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/internal/models"
	"github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type SignatureHandler struct {
	svc service.SignatureService
	validate *validator.Validate
}

func NewSignatureHandler(svc service.SignatureService, v *validator.Validate) *SignatureHandler {
	return &SignatureHandler{
		svc: svc,
		validate: v,
	}
}

// Request DTOs
type SignContractRequest struct {
	ContractID uuid.UUID `json:"contract_id" validate:"required"`
	VectorData []byte    `json:"vector_data" validate:"required_without=FileURL"`
	FileURL    string    `json:"file_url" validate:"required_without=VectorData,omitempty,url"`
}

// SignContract godoc
// @Summary      Create a signature
// @Description  Adds a user's signature (vector or image) to a pending contract
// @Tags         signatures
// @Accept       json
// @Produce      json
// @Param        signature body SignContractRequest true "Signature Data"
// @Success      201  {object}  models.Signature
// @Failure      400  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Router       /signatures [post]
// @Security     AuthKey
func (h *SignatureHandler) SignContract(c *gin.Context) {
	// 1. Get UserID from Middleware (already a uuid.UUID object)
	uidVal, _ := c.Get("userID")
	userID, ok := uidVal.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal context error"})
		return
	}

	// 2. Bind and Validate
	var req SignContractRequest
	if err := h.bindAndValidate(c, &req); err != nil {
		return
	}

	// 3. Map to Model
	signature := &models.Signature{
		ContractID: req.ContractID,
		UserID:     userID,
		VectorData: req.VectorData,
		FileURL:    req.FileURL,
	}
	
	// 4. Call Service (Service handles the "Status Check" and "Hashing")
	if err := h.svc.SignContract(c.Request.Context(), signature); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, signature)
}

// GetSignature godoc
// @Summary      Get a user's signature
// @Description  Retrieves a specific signature record by ID
// @Tags         signatures
// @Produce      json
// @Param        contract_id path string true "Contract ID"
// @Success      200  {object}  models.Signature
// @Failure      404  {object}  map[string]string
// @Router       /signatures/contract/{contract_id} [get]
// @Security     AuthKey
func (h *SignatureHandler) GetSignature(c *gin.Context) {
	contractID, err := uuid.Parse(c.Param("contract_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Contract ID format"})
		return
	}

	signatures, err := h.svc.GetContractSignatures(c.Request.Context(), contractID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch signatures"})
		return
	}

	c.JSON(http.StatusOK, signatures)
}

// RevokeSignature godoc
// @Summary      Remove a signature
// @Description  Deletes a signature mark if the contract is still pending
// @Tags         signatures
// @Param        id  path  string  true  "Signature ID"
// @Success      200  {object}  map[string]string
// @Router       /signatures/{id} [delete]
// @Security     AuthKey
func (h *SignatureHandler) RevokeSignature(c *gin.Context){
	sigID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Signature ID format"})
		return
	}

	uidVal, _ := c.Get("userID")
	userID, ok := uidVal.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal context error"})
		return
	}

	if err := h.svc.RevokeSignature(c.Request.Context(), sigID, userID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Signature revoked successfully"})
}

// Helper to keep handler logic dry
func (h *SignatureHandler) bindAndValidate(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindJSON(obj); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return err
	}
	if err := h.validate.Struct(obj); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}
	return nil
}
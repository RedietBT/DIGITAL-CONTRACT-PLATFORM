package handler

import (
	"net/http"
	"strings"

	"github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/internal/middleware"
	"github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/internal/models"
	"github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type ContractHandler struct {
	svc      service.ContractService
	validate *validator.Validate
}

func NewContractHandler(svc service.ContractService, v *validator.Validate) *ContractHandler {
	// Register the custom "no_scripts" validation
    v.RegisterValidation("no_scripts", func(fl validator.FieldLevel) bool {
        // Returns false if the string contains a script tag
        return !strings.Contains(strings.ToLower(fl.Field().String()), "<script")
    })
	return &ContractHandler{svc: svc, validate: v}
}

type CreateContractRequest struct {
	Title       string `json:"title" validate:"required,min=3,no_scripts"`
	Description string `json:"description" validate:"required,min=3,max=1000,no_scripts"`
	Content     string `json:"content" validate:"required,no_scripts"` // Initial version text
}

//CreateContract godoc
// @Summary      Create a new contract
// @Description  Creates a contract and initializes Version 1
// @Tags         contracts
// @Accept       json
// @Produce      json
// @Param contract body CreateContractRequest true "Contract"
// @Success      201         {object} models.Contract
// @Failure      400         {object} map[string]string
// @Failure      500         {object} map[string]string
// @Router       /contracts [post]
// @Security     BearerAuth
func (h *ContractHandler) CreateContract(c *gin.Context) {
	// 1. GEt UserID from Middleware context
	uidVal, _ := c.Get(middleware.UserIDKey)
	userID, _ := uuid.Parse(uidVal.(string))


	// 2. Bind JSON to request struct
	var req CreateContractRequest
	if err := h.bindAndValidate(c, &req); err != nil {
		return // h.bindAndValidate handles the error response
	}

	// 3. Map to Model
	contract := &models.Contract{
		OwnerUserID: userID,
		Title: req.Title,
		Description: req.Description,
	}

	// 4. Call Service 
	if err := h.svc.CreateContract(c.Request.Context(), contract, req.Content); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create contract"})
		return
	}

	c.JSON(http.StatusCreated, contract)
}

// bindAndValidate is a helper to keep handler clean
func(h *ContractHandler) bindAndValidate(c *gin.Context, obj interface{}) error {
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

// GetContract godoc
// @Summary      Get a single contract
// @Description  Retrieves contract details including versions and participants
// @Tags         contracts
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Contract ID"
// @Success      200  {object}  models.Contract
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /contracts/{id} [get]
// @Security     BearerAuth
func (h *ContractHandler) GetContract(c *gin.Context) {
	contractsID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	contract, err := h.svc.GetContractDetails(c.Request.Context(), contractsID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve contract"})
		return
	}

	c.JSON(http.StatusOK, contract)
}

// ListContracts godoc
// @Summary      List user contracts
// @Description  Get all contracts owned by the authenticated user
// @Tags         contracts
// @Accept       json
// @Produce      json
// @Success      200  {array}   models.Contract
// @Failure      500  {object}  map[string]string
// @Router       /contracts [get]
// @Security     BearerAuth
func (h *ContractHandler) ListContracts(c *gin.Context) {
	uidVal, _ := c.Get(middleware.UserIDKey)
	userID, _ := uuid.Parse(uidVal.(string))

	contracts, err := h.svc.ListUserContracts(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve contracts"})
		return
	}

	c.JSON(http.StatusOK, contracts)
}

// UpdateContractRequest defines what fields can be modified
type UpdateContractRequest struct {
	Title       string `json:"title" validate:"required,min=3,alphanum_start,no_scripts"`
	Description string `json:"description" validate:"required,min=3,max=1000,no_scripts"`
}

// UpdateContract godoc
// @Summary      Update contract metadata
// @Description  Update title/description of a draft contract
// @Tags         contracts
// @Accept       json
// @Produce      json
// @Param        id        path      string                 true  "Contract ID"
// @Param        contract  body      UpdateContractRequest  true  "Updated Data"
// @Success      200       {object}  map[string]string
// @Router       /contracts/{id} [put]
// @Security     BearerAuth
func(h *ContractHandler) UpdateContract(c *gin.Context) {
	contractID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	uidVal, _ := c.Get(middleware.UserIDKey)
	userID, _ := uuid.Parse(uidVal.(string))

	var req UpdateContractRequest
	if err := h.bindAndValidate(c, &req); err != nil {
		return
	}

	err = h.svc.UpdateContract(c.Request.Context(), contractID, userID, req.Title, req.Description)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Contract updated successfully"})
}

// DeleteContract godoc
// @Summary      Delete a contract
// @Description  Hard delete a contract and its associated data
// @Tags         contracts
// @Produce      json
// @Param        id   path      string  true  "Contract ID"
// @Success      200  {object}  map[string]string
// @Router       /contracts/{id} [delete]
// @Security     BearerAuth
func (h *ContractHandler) DeleteContract(c *gin.Context) {
	contractID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	uidVal, _ := c.Get(middleware.UserIDKey)
	userID, _ := uuid.Parse(uidVal.(string))

	err = h.svc.DeleteContract(c.Request.Context(), contractID, userID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Contract deleted successfully"})
}

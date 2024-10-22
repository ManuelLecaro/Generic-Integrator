package handler

import (
	"context"
	"generic-integration-platform/internal/application/dto"
	"generic-integration-platform/internal/application/services"
	errorDTO "generic-integration-platform/internal/infra/http/dto"

	"net/http"

	"github.com/gin-gonic/gin"
)

// IntegrationHandler manages integration operations.
type IntegrationHandler struct {
	service services.IIntegrationService
}

// NewIntegrationHandler creates a new IntegrationHandler.
func NewIntegrationHandler(s services.IIntegrationService) *IntegrationHandler {
	return &IntegrationHandler{
		service: s,
	}
}

// @Summary List all integrations
// @Description Get a list of all integrations
// @Tags Integrations
// @Produce json
// @Success 200 {array} dto.IntegrationResponseDTO
// @Failure 500 {object} errorDTO.ErrorResponseDTO
// @Router /integrations [get]
func (h *IntegrationHandler) GetIntegrations(c *gin.Context) {
	integrations, err := h.service.ListIntegrations(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorDTO.ErrorResponseDTO{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, integrations)
}

// @Summary Create a new integration
// @Description Create a new integration with the provided details
// @Tags Integrations
// @Accept json
// @Produce json
// @Param integration body dto.IntegrationRequestDTO true "Integration data"
// @Success 201 {object} dto.IntegrationResponseDTO
// @Failure 400 {object} errorDTO.ErrorResponseDTO
// @Failure 500 {object} errorDTO.ErrorResponseDTO
// @Router /integrations [post]
func (h *IntegrationHandler) CreateIntegration(c *gin.Context) {
	var input dto.IntegrationRequestDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errorDTO.ErrorResponseDTO{Message: "Invalid request payload"})
		return
	}

	integration, err := h.service.CreateIntegration(context.Background(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorDTO.ErrorResponseDTO{Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, integration)
}

// @Summary Get a specific integration by ID
// @Description Retrieve details of a specific integration by its ID
// @Tags Integrations
// @Produce json
// @Param id path string true "Integration ID"
// @Success 200 {object} dto.IntegrationResponseDTO
// @Failure 404 {object} errorDTO.ErrorResponseDTO
// @Failure 500 {object} errorDTO.ErrorResponseDTO
// @Router /integrations/{id} [get]
func (h *IntegrationHandler) GetIntegrationDetails(c *gin.Context) {
	id := c.Param("id")

	integration, err := h.service.GetIntegrationByID(context.Background(), id)
	if err != nil {
		if err == services.ErrIntegrationNotFound {
			c.JSON(http.StatusNotFound, errorDTO.ErrorResponseDTO{Message: "Integration not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, errorDTO.ErrorResponseDTO{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, integration)
}

// @Summary Update an existing integration by ID
// @Description Update the details of an existing integration
// @Tags Integrations
// @Accept json
// @Produce json
// @Param id path string true "Integration ID"
// @Param integration body dto.IntegrationRequestDTO true "Updated integration data"
// @Success 200 {object} dto.IntegrationResponseDTO
// @Failure 400 {object} errorDTO.ErrorResponseDTO
// @Failure 404 {object} errorDTO.ErrorResponseDTO
// @Failure 500 {object} errorDTO.ErrorResponseDTO
// @Router /integrations/{id} [put]
func (h *IntegrationHandler) UpdateIntegration(c *gin.Context) {
	id := c.Param("id")

	var input dto.IntegrationRequestDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errorDTO.ErrorResponseDTO{Message: "Invalid request payload"})
		return
	}

	integration, err := h.service.UpdateIntegration(context.Background(), id, input)
	if err != nil {
		if err == services.ErrIntegrationNotFound {
			c.JSON(http.StatusNotFound, errorDTO.ErrorResponseDTO{Message: "Integration not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, errorDTO.ErrorResponseDTO{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, integration)
}

// @Summary Delete an integration by ID
// @Description Remove a specific integration by its ID
// @Tags Integrations
// @Produce json
// @Param id path string true "Integration ID"
// @Success 204
// @Failure 404 {object} errorDTO.ErrorResponseDTO
// @Failure 500 {object} errorDTO.ErrorResponseDTO
// @Router /integrations/{id} [delete]
func (h *IntegrationHandler) DeleteIntegration(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.DeleteIntegration(context.Background(), id); err != nil {
		if err == services.ErrIntegrationNotFound {
			c.JSON(http.StatusNotFound, errorDTO.ErrorResponseDTO{Message: "Integration not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, errorDTO.ErrorResponseDTO{Message: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

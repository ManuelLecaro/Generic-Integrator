package handler

import (
	"context"
	"generic-integration-platform/internal/application/dto"
	"generic-integration-platform/internal/application/services"
	errorDTO "generic-integration-platform/internal/infra/http/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

type FlowHandler struct {
	service services.IFlowService
}

func NewFlowHandler(s services.IFlowService) *FlowHandler {
	return &FlowHandler{service: s}
}

// GetFlows handles the GET request to list all flows.
// @Summary Get all flows
// @Description Retrieve a list of all flows
// @Tags Flows
// @Produce json
// @Success 200 {array} dto.FlowDTO
// @Failure 500 {object} errorDTO.ErrorResponseDTO
// @Router /flows [get]
func (h *FlowHandler) GetFlows(c *gin.Context) {
	flows, err := h.service.ListFlows(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorDTO.ErrorResponseDTO{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, flows)
}

// CreateFlow handles the POST request to create a new flow.
// @Summary Create a new flow
// @Description Create a new flow by providing flow details
// @Tags Flows
// @Accept json
// @Produce json
// @Param flow body dto.FlowDTO true "Flow Data"
// @Success 201 {object} dto.FlowDTO
// @Failure 400 {object} errorDTO.ErrorResponseDTO
// @Failure 500 {object} errorDTO.ErrorResponseDTO
// @Router /flows [post]
func (h *FlowHandler) CreateFlow(c *gin.Context) {
	var flowDTO dto.FlowDTO
	if err := c.ShouldBindJSON(&flowDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	flow, err := h.service.CreateFlow(context.Background(), flowDTO)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, flow)
}

// GetFlowDetails handles the GET request to retrieve a specific flow by ID.
// @Summary Get flow details by ID
// @Description Retrieve details of a specific flow by ID
// @Tags Flows
// @Produce json
// @Param id path string true "Flow ID"
// @Success 200 {object} dto.FlowDTO
// @Failure 404 {object} errorDTO.ErrorResponseDTO
// @Failure 500 {object} errorDTO.ErrorResponseDTO
// @Router /flows/{id} [get]
func (h *FlowHandler) GetFlowDetails(c *gin.Context) {
	id := c.Param("id")
	flow, err := h.service.GetFlowByID(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, flow)
}

// UpdateFlow handles the PUT request to update an existing flow by ID.
// @Summary Update a flow
// @Description Update an existing flow by providing its ID and updated details
// @Tags Flows
// @Accept json
// @Produce json
// @Param id path string true "Flow ID"
// @Param flow body dto.FlowDTO true "Updated Flow Data"
// @Success 200 {object} dto.FlowDTO
// @Failure 400 {object} errorDTO.ErrorResponseDTO
// @Failure 404 {object} errorDTO.ErrorResponseDTO
// @Failure 500 {object} errorDTO.ErrorResponseDTO
// @Router /flows/{id} [put]
func (h *FlowHandler) UpdateFlow(c *gin.Context) {
	id := c.Param("id")
	var flowDTO dto.FlowDTO
	if err := c.ShouldBindJSON(&flowDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	flow, err := h.service.UpdateFlow(context.Background(), id, flowDTO)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, flow)
}

// DeleteFlow handles the DELETE request to remove a specific flow by ID.
// @Summary Delete a flow by ID
// @Description Remove a specific flow by its ID
// @Tags Flows
// @Produce json
// @Param id path string true "Flow ID"
// @Success 204
// @Failure 404 {object} errorDTO.ErrorResponseDTO
// @Failure 500 {object} errorDTO.ErrorResponseDTO
// @Router /flows/{id} [delete]
func (h *FlowHandler) DeleteFlow(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteFlow(context.Background(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// ExecuteFlow handles the POST request to execute a specific flow by ID.
// @Summary Execute a flow by ID
// @Description Execute a specific flow by its ID
// @Tags Flows
// @Produce json
// @Param id path string true "Flow ID"
// @Success 200 {object} dto.FlowDTO
// @Failure 404 {object} errorDTO.ErrorResponseDTO
// @Failure 500 {object} errorDTO.ErrorResponseDTO
// @Router /flows/{id}/execute [post]
func (h *FlowHandler) ExecuteFlow(c *gin.Context) {
	id := c.Param("id")
	result, err := h.service.ExecuteFlow(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

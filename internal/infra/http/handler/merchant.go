package handler

import (
	"agnostic-payment-platform/internal/application/merchant/dto"
	service "agnostic-payment-platform/internal/application/merchant/ports/services"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MerchantHandler struct {
	service service.MerchantService
}

func NewMerchantHandler(s service.MerchantService) *MerchantHandler {
	return &MerchantHandler{
		service: s,
	}
}

// @Summary Register a new merchant
// @Description Register a new merchant with an email, password, and name.
// @Tags Merchants
// @Accept json
// @Produce json
// @Param signup body dto.SignupRequestDTO true "Signup request"
// @Success 201 {object} dto.SuccessResponseDTO
// @Failure 400 {object} dto.ErrorResponseDTO
// @Failure 500 {object} dto.ErrorResponseDTO
// @Router /merchant/signup [post]
func (h *MerchantHandler) Signup(c *gin.Context) {
	var input dto.SignupRequestDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponseDTO{Message: "Invalid request payload"})
		return
	}

	if err := h.service.Register(context.Background(), input.Email, input.Password, input.Name); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponseDTO{Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.SuccessResponseDTO{Message: "Merchant registered successfully"})
}

// @Summary Login as a merchant
// @Description Login a merchant with an email and password.
// @Tags Merchants
// @Accept json
// @Produce json
// @Param login body dto.LoginRequestDTO true "Login request"
// @Success 200 {object} dto.LoginResponseDTO
// @Failure 400 {object} dto.ErrorResponseDTO
// @Failure 401 {object} dto.ErrorResponseDTO
// @Failure 500 {object} dto.ErrorResponseDTO
// @Router /merchant/login [post]
func (h *MerchantHandler) Login(c *gin.Context) {
	var input dto.LoginRequestDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponseDTO{Message: "Invalid request payload"})
		return
	}

	merchant, err := h.service.Login(context.Background(), input.Email, input.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponseDTO{Message: err.Error()})
		return
	}

	response := dto.LoginResponseDTO{
		ID:    merchant.ID,
		Email: merchant.Email,
		Name:  merchant.Name,
	}

	c.JSON(http.StatusOK, response)
}

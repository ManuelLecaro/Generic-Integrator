package handler

import (
	"agnostic-payment-platform/internal/application/payment/dto"
	"agnostic-payment-platform/internal/application/payment/service"
	"agnostic-payment-platform/internal/domain/payment"
	"net/http"

	"github.com/gin-gonic/gin"
)

// PaymentHandler is the struct that holds the dependencies for the payment handlers.
type PaymentHandler struct {
	Service service.PaymentService
}

// NewPaymentHandler initializes a new PaymentHandler with its dependencies.
func NewPaymentHandler(service *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{Service: *service}
}

// ProcessPayment handles the payment processing request.
// @Summary Process a payment
// @Description Processes a payment based on the provided details
// @Tags payments
// @Accept json
// @Produce json
// @Param payment body dto.PaymentRequestDTO true "Payment Request"
// @Success 200 {object} dto.PaymentResponseDTO
// @Failure 400 {object} gin.H "Invalid request payload"
// @Failure 500 {object} gin.H "Unable to process payment"
// @Router /payments [post]
func (h *PaymentHandler) ProcessPayment(c *gin.Context) {
	var input dto.PaymentRequestDTO

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	response, err := h.Service.ProcessPayment(c, &input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to process payment"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetPaymentDetails handles the request to get details of a previous payment by its ID.
// @Summary Get payment details
// @Description Retrieves details of a payment by its ID
// @Tags payments
// @Accept json
// @Produce json
// @Param id path string true "Payment ID"
// @Success 200 {object} dto.PaymentDetailsDTO
// @Failure 400 {object} gin.H "Missing payment ID"
// @Failure 404 {object} gin.H "Payment not found"
// @Failure 500 {object} gin.H "Unable to retrieve payment details"
// @Router /payments/{id} [get]
func (h *PaymentHandler) GetPaymentDetails(c *gin.Context) {
	var input dto.PaymentsByIDDTO

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing payment ID"})
		return
	}

	input.ID = id

	details, err := h.Service.GetPaymentDetails(c, input.ID)
	if err != nil {
		if err == payment.ErrPaymentNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve payment details"})
		return
	}

	c.JSON(http.StatusOK, details)
}

// @Summary Process a refund request
// @Description Handles the refund request by extracting refund details from the request body and processing the refund.
// @Tags payments
// @Accept json
// @Produce json
// @Param refundDTO body dto.RefundDTO true "Refund details"
// @Success 200 {object} dto.PaymentResponseDTO "Refund processed successfully"
// @Failure 400 {object} gin.H "Invalid request payload"
// @Failure 500 {object} gin.H "Unable to process payment"
// @Router /payments/refund [post]
func (h *PaymentHandler) Refund(c *gin.Context) {
	var input dto.RefundDTO

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	response, err := h.Service.ProcessRefund(c, &input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to process payment"})
		return
	}

	c.JSON(http.StatusOK, response)
}

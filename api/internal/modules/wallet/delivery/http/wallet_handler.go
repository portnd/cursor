package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/portnd/the-sentinel-core/internal/modules/wallet/domain"
)

type WalletHandler struct {
	usecase domain.Usecase
}

// NewWalletHandler creates a new wallet handler
func NewWalletHandler(usecase domain.Usecase) *WalletHandler {
	return &WalletHandler{
		usecase: usecase,
	}
}

// GetMyWallet retrieves the authenticated user's wallet
// GET /api/wallets/me
func (h *WalletHandler) GetMyWallet(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "user not authenticated",
		})
		return
	}

	// Convert to uint
	uid, ok := userID.(float64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "invalid user ID format",
		})
		return
	}

	// Get or create wallet
	wallet, err := h.usecase.GetMyWallet(uint(uid))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Wallet retrieved successfully",
		"data":    wallet,
	})
}

// Transfer performs a money transfer
// POST /api/wallets/transfer
func (h *WalletHandler) Transfer(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "user not authenticated",
		})
		return
	}

	// Convert to uint
	uid, ok := userID.(float64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "invalid user ID format",
		})
		return
	}

	// Parse request body
	var req domain.TransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	// Perform transfer
	transaction, err := h.usecase.Transfer(uint(uid), req.ToWalletID, req.Amount)
	if err != nil {
		// Check for business logic errors (400) vs system errors (500)
		statusCode := http.StatusInternalServerError
		if err.Error() == "insufficient funds" || 
		   err.Error()[:len("insufficient funds")] == "insufficient funds" ||
		   err.Error() == "sender wallet not found" ||
		   err.Error() == "receiver wallet not found" ||
		   err.Error() == "cannot transfer to yourself" ||
		   err.Error()[:len("invalid amount")] == "invalid amount" {
			statusCode = http.StatusBadRequest
		}

		c.JSON(statusCode, gin.H{
			"error":   "Transfer Failed",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Transfer successful",
		"data":    transaction,
	})
}

// GetMyTransactions retrieves user's transaction history
// GET /api/wallets/transactions?limit=20
func (h *WalletHandler) GetMyTransactions(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "user not authenticated",
		})
		return
	}

	// Convert to uint
	uid, ok := userID.(float64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "invalid user ID format",
		})
		return
	}

	// Parse query params
	limitStr := c.DefaultQuery("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 20
	}

	// Get transactions
	transactions, err := h.usecase.GetMyTransactions(uint(uid), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Transactions retrieved successfully",
		"data":    transactions,
	})
}

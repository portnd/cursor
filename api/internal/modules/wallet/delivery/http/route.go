package http

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers wallet routes
// authMiddleware is required to protect these endpoints
func RegisterRoutes(router *gin.RouterGroup, handler *WalletHandler, authMiddleware gin.HandlerFunc) {
	// All wallet routes require authentication
	wallet := router.Group("/wallets")
	wallet.Use(authMiddleware)
	{
		wallet.GET("/me", handler.GetMyWallet)
		wallet.POST("/transfer", handler.Transfer)
		wallet.GET("/transactions", handler.GetMyTransactions)
	}
}

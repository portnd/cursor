package domain

import (
	"time"

	"gorm.io/gorm"
)

// Wallet represents a user's wallet with balance
type Wallet struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"uniqueIndex;not null" json:"user_id"` // One wallet per user
	Balance   float64   `gorm:"type:decimal(15,2);default:0;not null" json:"balance"`
	Currency  string    `gorm:"type:varchar(3);default:'THB';not null" json:"currency"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TransactionType represents the type of transaction
type TransactionType string

const (
	TransactionTypeDeposit  TransactionType = "DEPOSIT"
	TransactionTypeWithdraw TransactionType = "WITHDRAW"
	TransactionTypeTransfer TransactionType = "TRANSFER"
)

// TransactionStatus represents the status of transaction
type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "PENDING"
	TransactionStatusCompleted TransactionStatus = "COMPLETED"
	TransactionStatusFailed    TransactionStatus = "FAILED"
)

// Transaction represents a financial transaction
type Transaction struct {
	ID           uint              `gorm:"primaryKey" json:"id"`
	FromWalletID *uint             `gorm:"index" json:"from_wallet_id,omitempty"` // Nullable for deposits
	ToWalletID   *uint             `gorm:"index" json:"to_wallet_id,omitempty"`   // Nullable for withdrawals
	Amount       float64           `gorm:"type:decimal(15,2);not null" json:"amount"`
	Type         TransactionType   `gorm:"type:varchar(20);not null" json:"type"`
	Status       TransactionStatus `gorm:"type:varchar(20);default:'PENDING';not null" json:"status"`
	Description  string            `gorm:"type:text" json:"description,omitempty"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
}

// TransferRequest represents a transfer request DTO
type TransferRequest struct {
	ToWalletID uint    `json:"to_wallet_id" binding:"required"`
	Amount     float64 `json:"amount" binding:"required,gt=0"`
}

// WalletResponse represents wallet response DTO
type WalletResponse struct {
	ID       uint    `json:"id"`
	UserID   uint    `json:"user_id"`
	Balance  float64 `json:"balance"`
	Currency string  `json:"currency"`
}

// TransactionResponse represents transaction response DTO
type TransactionResponse struct {
	ID           uint              `json:"id"`
	FromWalletID *uint             `json:"from_wallet_id,omitempty"`
	ToWalletID   *uint             `json:"to_wallet_id,omitempty"`
	Amount       float64           `json:"amount"`
	Type         TransactionType   `json:"type"`
	Status       TransactionStatus `json:"status"`
	Description  string            `json:"description,omitempty"`
	CreatedAt    time.Time         `json:"created_at"`
}

// Repository interface for wallet persistence
type Repository interface {
	// GetWalletByUserID retrieves wallet by user ID
	// tx is optional, if provided, uses that transaction context
	GetWalletByUserID(userID uint, tx *gorm.DB) (*Wallet, error)

	// GetWalletByID retrieves wallet by ID with optional locking
	// tx is optional, if provided, uses that transaction context
	// forUpdate: if true, applies FOR UPDATE lock (pessimistic locking)
	GetWalletByID(walletID uint, tx *gorm.DB, forUpdate bool) (*Wallet, error)

	// CreateWallet creates a new wallet for user
	CreateWallet(wallet *Wallet, tx *gorm.DB) error

	// UpdateBalance updates wallet balance
	// CRITICAL: This should only be called within a transaction
	UpdateBalance(walletID uint, newBalance float64, tx *gorm.DB) error

	// CreateTransaction logs a transaction
	CreateTransaction(transaction *Transaction, tx *gorm.DB) error

	// GetTransactionsByUserID retrieves transaction history
	GetTransactionsByUserID(userID uint, limit int) ([]Transaction, error)
}

// Usecase interface for wallet business logic
type Usecase interface {
	// GetMyWallet retrieves or creates user's wallet
	GetMyWallet(userID uint) (*WalletResponse, error)

	// Transfer performs a transfer between wallets with ACID guarantees
	// CRITICAL: Uses pessimistic locking to prevent race conditions
	Transfer(fromUserID uint, toWalletID uint, amount float64) (*TransactionResponse, error)

	// GetMyTransactions retrieves user's transaction history
	GetMyTransactions(userID uint, limit int) ([]TransactionResponse, error)
}

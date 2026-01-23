package repository

import (
	"fmt"

	"github.com/komgrip/starter-kit/internal/modules/wallet/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type postgresRepository struct {
	db *gorm.DB
}

// NewPostgresRepository creates a new wallet repository
func NewPostgresRepository(db *gorm.DB) domain.Repository {
	return &postgresRepository{db: db}
}

// GetWalletByUserID retrieves wallet by user ID
func (r *postgresRepository) GetWalletByUserID(userID uint, tx *gorm.DB) (*domain.Wallet, error) {
	db := r.getDB(tx)

	var wallet domain.Wallet
	err := db.Where("user_id = ?", userID).First(&wallet).Error
	if err != nil {
		return nil, err
	}

	return &wallet, nil
}

// GetWalletByID retrieves wallet by ID with optional pessimistic locking
// CRITICAL: forUpdate = true applies SELECT ... FOR UPDATE (row-level lock)
func (r *postgresRepository) GetWalletByID(walletID uint, tx *gorm.DB, forUpdate bool) (*domain.Wallet, error) {
	db := r.getDB(tx)

	var wallet domain.Wallet
	query := db.Where("id = ?", walletID)

	// Apply pessimistic locking if requested
	// This prevents other transactions from reading/modifying this row
	// until current transaction commits or rolls back
	if forUpdate {
		query = query.Clauses(clause.Locking{Strength: "UPDATE"})
	}

	err := query.First(&wallet).Error
	if err != nil {
		return nil, err
	}

	return &wallet, nil
}

// CreateWallet creates a new wallet
func (r *postgresRepository) CreateWallet(wallet *domain.Wallet, tx *gorm.DB) error {
	db := r.getDB(tx)
	return db.Create(wallet).Error
}

// UpdateBalance updates wallet balance
// CRITICAL: This must be called within a transaction with locked wallet
func (r *postgresRepository) UpdateBalance(walletID uint, newBalance float64, tx *gorm.DB) error {
	db := r.getDB(tx)

	// Security check: Prevent negative balance at database level
	if newBalance < 0 {
		return fmt.Errorf("invalid balance: cannot set negative balance")
	}

	result := db.Model(&domain.Wallet{}).
		Where("id = ?", walletID).
		Update("balance", newBalance)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("wallet not found or balance unchanged")
	}

	return nil
}

// CreateTransaction logs a transaction
func (r *postgresRepository) CreateTransaction(transaction *domain.Transaction, tx *gorm.DB) error {
	db := r.getDB(tx)
	return db.Create(transaction).Error
}

// GetTransactionsByUserID retrieves user's transaction history
func (r *postgresRepository) GetTransactionsByUserID(userID uint, limit int) ([]domain.Transaction, error) {
	// First get user's wallet to find wallet ID
	wallet, err := r.GetWalletByUserID(userID, nil)
	if err != nil {
		return nil, err
	}

	var transactions []domain.Transaction
	query := r.db.Where("from_wallet_id = ? OR to_wallet_id = ?", wallet.ID, wallet.ID).
		Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err = query.Find(&transactions).Error
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

// getDB returns the transaction DB if provided, otherwise returns default DB
func (r *postgresRepository) getDB(tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx
	}
	return r.db
}

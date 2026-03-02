package usecase

import (
	"fmt"

	"github.com/portnd/the-sentinel-core/internal/modules/wallet/domain"
	"gorm.io/gorm"
)

type walletUsecase struct {
	repo domain.Repository
	db   *gorm.DB
}

// NewWalletUsecase creates a new wallet usecase
func NewWalletUsecase(repo domain.Repository, db *gorm.DB) domain.Usecase {
	return &walletUsecase{
		repo: repo,
		db:   db,
	}
}

// GetMyWallet retrieves or creates user's wallet
func (u *walletUsecase) GetMyWallet(userID uint) (*domain.WalletResponse, error) {
	// Try to get existing wallet
	wallet, err := u.repo.GetWalletByUserID(userID, nil)

	// If not found, create new wallet with 0 balance
	if err == gorm.ErrRecordNotFound {
		newWallet := &domain.Wallet{
			UserID:   userID,
			Balance:  0,
			Currency: "THB",
		}

		err = u.repo.CreateWallet(newWallet, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create wallet: %w", err)
		}

		wallet = newWallet
	} else if err != nil {
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}

	return &domain.WalletResponse{
		ID:       wallet.ID,
		UserID:   wallet.UserID,
		Balance:  wallet.Balance,
		Currency: wallet.Currency,
	}, nil
}

// Transfer performs a money transfer with ACID guarantees and pessimistic locking
// CRITICAL: This function implements the following security measures:
// 1. Database Transaction (ACID compliance)
// 2. Pessimistic Locking (SELECT ... FOR UPDATE)
// 3. Balance Validation (prevent negative balance)
// 4. Transaction Logging (audit trail)
func (u *walletUsecase) Transfer(fromUserID uint, toWalletID uint, amount float64) (*domain.TransactionResponse, error) {
	// Validation: Amount must be positive
	if amount <= 0 {
		return nil, fmt.Errorf("invalid amount: must be greater than 0")
	}

	var transactionRecord *domain.Transaction

	// Start ACID Transaction
	// All operations within this transaction are atomic
	// If any step fails, entire transaction rolls back
	err := u.db.Transaction(func(tx *gorm.DB) error {
		// Step 1: Get sender's wallet (without lock)
		senderWallet, err := u.repo.GetWalletByUserID(fromUserID, tx)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("sender wallet not found")
			}
			return fmt.Errorf("failed to get sender wallet: %w", err)
		}

		// Step 2: Lock sender's wallet (Pessimistic Locking)
		// This prevents other transactions from modifying this wallet
		// until our transaction commits or rolls back
		senderWallet, err = u.repo.GetWalletByID(senderWallet.ID, tx, true)
		if err != nil {
			return fmt.Errorf("failed to lock sender wallet: %w", err)
		}

		// Step 3: Check sender's balance (Business Logic Validation)
		if senderWallet.Balance < amount {
			return fmt.Errorf("insufficient funds: balance %.2f, required %.2f", senderWallet.Balance, amount)
		}

		// Step 4: Lock receiver's wallet (Pessimistic Locking)
		// CRITICAL: Lock in consistent order (by ID) to prevent deadlocks
		receiverWallet, err := u.repo.GetWalletByID(toWalletID, tx, true)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("receiver wallet not found")
			}
			return fmt.Errorf("failed to lock receiver wallet: %w", err)
		}

		// Security: Prevent self-transfer
		if senderWallet.ID == receiverWallet.ID {
			return fmt.Errorf("cannot transfer to yourself")
		}

		// Step 5: Calculate new balances
		newSenderBalance := senderWallet.Balance - amount
		newReceiverBalance := receiverWallet.Balance + amount

		// Security: Double-check no negative balance
		if newSenderBalance < 0 {
			return fmt.Errorf("operation would result in negative balance")
		}

		// Step 6: Update sender's balance
		err = u.repo.UpdateBalance(senderWallet.ID, newSenderBalance, tx)
		if err != nil {
			return fmt.Errorf("failed to update sender balance: %w", err)
		}

		// Step 7: Update receiver's balance
		err = u.repo.UpdateBalance(receiverWallet.ID, newReceiverBalance, tx)
		if err != nil {
			return fmt.Errorf("failed to update receiver balance: %w", err)
		}

		// Step 8: Log transaction (Audit Trail)
		transactionRecord = &domain.Transaction{
			FromWalletID: &senderWallet.ID,
			ToWalletID:   &receiverWallet.ID,
			Amount:       amount,
			Type:         domain.TransactionTypeTransfer,
			Status:       domain.TransactionStatusCompleted,
			Description:  fmt.Sprintf("Transfer from wallet %d to wallet %d", senderWallet.ID, receiverWallet.ID),
		}

		err = u.repo.CreateTransaction(transactionRecord, tx)
		if err != nil {
			return fmt.Errorf("failed to log transaction: %w", err)
		}

		// If we reach here, all operations succeeded
		// Transaction will commit automatically
		return nil
	})

	if err != nil {
		return nil, err
	}

	// Return transaction response
	return &domain.TransactionResponse{
		ID:           transactionRecord.ID,
		FromWalletID: transactionRecord.FromWalletID,
		ToWalletID:   transactionRecord.ToWalletID,
		Amount:       transactionRecord.Amount,
		Type:         transactionRecord.Type,
		Status:       transactionRecord.Status,
		Description:  transactionRecord.Description,
		CreatedAt:    transactionRecord.CreatedAt,
	}, nil
}

// GetMyTransactions retrieves user's transaction history
func (u *walletUsecase) GetMyTransactions(userID uint, limit int) ([]domain.TransactionResponse, error) {
	// Default limit
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	transactions, err := u.repo.GetTransactionsByUserID(userID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}

	// Convert to response DTOs
	responses := make([]domain.TransactionResponse, len(transactions))
	for i, tx := range transactions {
		responses[i] = domain.TransactionResponse{
			ID:           tx.ID,
			FromWalletID: tx.FromWalletID,
			ToWalletID:   tx.ToWalletID,
			Amount:       tx.Amount,
			Type:         tx.Type,
			Status:       tx.Status,
			Description:  tx.Description,
			CreatedAt:    tx.CreatedAt,
		}
	}

	return responses, nil
}

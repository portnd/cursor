# 💰 Wallet & Transaction Module

## 📋 Overview

The **Wallet Module** is a mission-critical financial system built with **God-Tier Security Standards**. It handles monetary transactions with **ACID compliance**, **Pessimistic Locking**, and **Zero-Tolerance for Data Inconsistency**.

This module is designed to prevent:
- ❌ **Double-Spending Attacks**
- ❌ **Race Conditions** (concurrent transfers)
- ❌ **Negative Balances** (overdrafts)
- ❌ **Transaction Loss** (all operations are atomic)

---

## 🏗️ Architecture (Hexagonal/Clean Architecture)

```
┌─────────────────────────────────────────────────────────────┐
│                      HTTP Request                           │
│                  POST /wallets/transfer                     │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│  🌐 DELIVERY LAYER (HTTP Handler)                          │
│  📁 delivery/http/wallet_handler.go                        │
│                                                             │
│  • Validate Request (DTO binding)                          │
│  • Extract UserID from JWT Context                         │
│  • Call Usecase                                            │
│  • Return JSON Response                                    │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│  💼 USECASE LAYER (Business Logic)                         │
│  📁 usecase/wallet_usecase.go                              │
│                                                             │
│  🔒 CRITICAL: DB Transaction Starts Here                   │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ 1. Lock Sender Wallet (SELECT FOR UPDATE)          │   │
│  │ 2. Validate Balance (prevent negative)             │   │
│  │ 3. Lock Receiver Wallet (SELECT FOR UPDATE)        │   │
│  │ 4. Update Both Balances                            │   │
│  │ 5. Log Transaction Record                          │   │
│  │ 6. Commit (or Rollback on ANY error)               │   │
│  └─────────────────────────────────────────────────────┘   │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│  🗄️ REPOSITORY LAYER (Data Access)                         │
│  📁 repository/postgres_repo.go                            │
│                                                             │
│  • GetWalletByID(id, tx, forUpdate=true)                   │
│  • UpdateBalance(id, newBalance, tx)                       │
│  • CreateTransaction(transaction, tx)                      │
│                                                             │
│  🔐 Uses: tx.Clauses(clause.Locking{Strength: "UPDATE"})   │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│  💾 DATABASE (PostgreSQL with ACID)                        │
│                                                             │
│  SELECT * FROM wallets WHERE id = ? FOR UPDATE;            │
│  UPDATE wallets SET balance = ? WHERE id = ?;              │
│  INSERT INTO transactions (...) VALUES (...);              │
│                                                             │
│  🔒 Row-Level Locking Active (no concurrent access)        │
└─────────────────────────────────────────────────────────────┘
```

---

## 🔑 Key Features

### 1. **Pessimistic Locking** 🔒

**Problem:** What if two users transfer money from the same wallet simultaneously?

```go
// Without Locking (DANGEROUS ❌)
wallet := repo.GetWallet(1) // User A reads: balance = 1000
wallet := repo.GetWallet(1) // User B reads: balance = 1000 (same time!)
repo.Update(1, 800)         // User A transfers 200
repo.Update(1, 600)         // User B transfers 400
// Result: Final balance = 600 (should be 200) ❌ RACE CONDITION!
```

**Solution:** Use `SELECT ... FOR UPDATE` (Pessimistic Locking)

```go
// With Locking (SAFE ✅)
tx.Begin()
wallet := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&wallet, id)
// 🔒 This wallet row is now LOCKED. User B must WAIT until User A finishes.
wallet.Balance -= 200
tx.Save(&wallet)
tx.Commit() // 🔓 Lock released
```

**Code Location:** `repository/postgres_repo.go:50`

```go
func (r *postgresRepository) GetWalletByID(walletID uint, tx *gorm.DB, forUpdate bool) (*domain.Wallet, error) {
	var wallet domain.Wallet
	query := tx.Where("id = ?", walletID)
	
	if forUpdate {
		// 🔒 CRITICAL: This prevents concurrent access
		query = query.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	
	result := query.First(&wallet)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get wallet by ID %d: %w", walletID, result.Error)
	}
	return &wallet, nil
}
```

---

### 2. **Atomic Transactions** ⚛️

**Guarantee:** ALL operations succeed, or NONE succeed. No partial updates.

```go
// Example: Transfer 500 THB from Wallet 1 to Wallet 2
err := db.Transaction(func(tx *gorm.DB) error {
    // Step 1: Lock & Get Sender
    sender, err := repo.GetWalletByID(1, tx, true)
    if err != nil { return err } // ❌ Rollback entire transaction
    
    // Step 2: Validate Balance
    if sender.Balance < 500 {
        return fmt.Errorf("insufficient funds") // ❌ Rollback
    }
    
    // Step 3: Lock & Get Receiver
    receiver, err := repo.GetWalletByID(2, tx, true)
    if err != nil { return err } // ❌ Rollback
    
    // Step 4: Update Balances
    sender.Balance -= 500
    receiver.Balance += 500
    tx.Save(&sender)
    tx.Save(&receiver)
    
    // Step 5: Log Transaction
    transaction := &Transaction{...}
    tx.Create(&transaction)
    
    return nil // ✅ Commit all changes atomically
})

if err != nil {
    // All changes above are ROLLED BACK automatically
    log.Printf("Transfer failed: %v", err)
}
```

**Code Location:** `usecase/wallet_usecase.go:30`

---

### 3. **Balance Validation** ✅

**Prevent Negative Balance:** No user can overdraft their account.

```go
// usecase/wallet_usecase.go:47
if senderWallet.Balance < req.Amount {
    return fmt.Errorf("insufficient funds: balance %.2f, required %.2f", 
        senderWallet.Balance, req.Amount)
}
```

---

### 4. **Transaction Logging** 📝

Every transfer is recorded in the `transactions` table for audit purposes.

```go
transactionRecord := &domain.Transaction{
    FromWalletID: senderWallet.ID,
    ToWalletID:   receiverWallet.ID,
    Amount:       req.Amount,
    Type:         domain.TransactionTypeTransfer,
    Status:       domain.TransactionStatusCompleted,
    Description:  fmt.Sprintf("Transfer from wallet %d to wallet %d", senderWallet.ID, receiverWallet.ID),
    CreatedAt:    time.Now(),
}
txRepo.CreateTransaction(transactionRecord, tx)
```

---

## 🌐 API Reference

### 1. **Get My Wallet**

**Endpoint:** `GET /wallets/me`

**Authentication:** Required (JWT Bearer Token)

**Description:** Retrieves the authenticated user's wallet. If the wallet doesn't exist, it will be **auto-created** with a balance of 0.00 THB.

**Request:**
```bash
curl -X GET http://localhost:8080/wallets/me \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Response (200 OK):**
```json
{
  "message": "Wallet retrieved successfully",
  "data": {
    "id": 1,
    "user_id": 5,
    "balance": 1000.50,
    "currency": "THB",
    "created_at": "2026-01-23T10:30:00Z",
    "updated_at": "2026-01-23T14:20:00Z"
  }
}
```

---

### 2. **Transfer Money**

**Endpoint:** `POST /wallets/transfer`

**Authentication:** Required (JWT Bearer Token)

**Description:** Transfers money from the authenticated user's wallet to another wallet. This operation is **atomic** and uses **pessimistic locking**.

**Request:**
```bash
curl -X POST http://localhost:8080/wallets/transfer \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "to_wallet_id": 2,
    "amount": 250.50
  }'
```

**Response (200 OK):**
```json
{
  "message": "Transfer successful",
  "data": {
    "id": 42,
    "from_wallet_id": 1,
    "to_wallet_id": 2,
    "amount": 250.50,
    "type": "TRANSFER",
    "status": "COMPLETED",
    "description": "Transfer from wallet 1 to wallet 2",
    "created_at": "2026-01-23T14:25:30Z"
  }
}
```

**Error Response (409 Conflict - Insufficient Funds):**
```json
{
  "error": "Transfer Failed",
  "message": "transfer failed: insufficient funds: balance 100.00, required 250.50"
}
```

**Error Response (404 Not Found - Invalid Wallet ID):**
```json
{
  "error": "Transfer Failed",
  "message": "failed to get receiver wallet: record not found"
}
```

---

### 3. **Get Transaction History**

**Endpoint:** `GET /wallets/transactions`

**Authentication:** Required (JWT Bearer Token)

**Description:** Retrieves all transactions related to the authenticated user's wallet (both sent and received).

**Request:**
```bash
curl -X GET http://localhost:8080/wallets/transactions \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Response (200 OK):**
```json
{
  "message": "Transactions retrieved successfully",
  "data": [
    {
      "id": 42,
      "from_wallet_id": 1,
      "to_wallet_id": 2,
      "amount": 250.50,
      "type": "TRANSFER",
      "status": "COMPLETED",
      "description": "Transfer from wallet 1 to wallet 2",
      "created_at": "2026-01-23T14:25:30Z"
    },
    {
      "id": 41,
      "from_wallet_id": 3,
      "to_wallet_id": 1,
      "amount": 500.00,
      "type": "TRANSFER",
      "status": "COMPLETED",
      "description": "Transfer from wallet 3 to wallet 1",
      "created_at": "2026-01-23T13:10:15Z"
    }
  ]
}
```

---

## 🧪 Testing Guide

### **Prerequisites:**
1. Start the services: `make up`
2. Register two test users: `testuser1@komgrip.com` and `testuser2@komgrip.com`
3. Login to get JWT tokens for both users

### **Test Scenario: Transfer Money**

#### **Step 1: Login as User 1**
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "testuser1@komgrip.com",
    "password": "password123"
  }'
```

**Copy the `token` from response:**
```json
{
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {...}
  }
}
```

Set as environment variable:
```bash
export TOKEN1="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

---

#### **Step 2: Get User 1's Wallet**
```bash
curl -X GET http://localhost:8080/wallets/me \
  -H "Authorization: Bearer $TOKEN1"
```

**Expected:** Wallet is auto-created with 0.00 THB balance.

```json
{
  "message": "Wallet retrieved successfully",
  "data": {
    "id": 1,
    "user_id": 1,
    "balance": 0.00,
    "currency": "THB"
  }
}
```

---

#### **Step 3: Manually Add Balance (For Testing Only)**

**Open PostgreSQL Shell:**
```bash
make shell-db
```

**Add 1000 THB to User 1's wallet:**
```sql
UPDATE wallets SET balance = 1000.00 WHERE id = 1;
SELECT * FROM wallets WHERE id = 1;
```

**Exit:**
```bash
\q
```

---

#### **Step 4: Login as User 2 and Get Wallet ID**

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "testuser2@komgrip.com",
    "password": "password123"
  }'

export TOKEN2="<USER_2_JWT_TOKEN>"

curl -X GET http://localhost:8080/wallets/me \
  -H "Authorization: Bearer $TOKEN2"
```

**Note the `id` field (e.g., Wallet ID = 2)**

---

#### **Step 5: Transfer from User 1 to User 2**

```bash
curl -X POST http://localhost:8080/wallets/transfer \
  -H "Authorization: Bearer $TOKEN1" \
  -H "Content-Type: application/json" \
  -d '{
    "to_wallet_id": 2,
    "amount": 250.50
  }'
```

**Expected Response:**
```json
{
  "message": "Transfer successful",
  "data": {
    "id": 1,
    "from_wallet_id": 1,
    "to_wallet_id": 2,
    "amount": 250.50,
    "type": "TRANSFER",
    "status": "COMPLETED"
  }
}
```

---

#### **Step 6: Verify Balances**

**User 1 (Sender):**
```bash
curl -X GET http://localhost:8080/wallets/me \
  -H "Authorization: Bearer $TOKEN1"
```

**Expected:** `balance: 749.50` (1000.00 - 250.50)

**User 2 (Receiver):**
```bash
curl -X GET http://localhost:8080/wallets/me \
  -H "Authorization: Bearer $TOKEN2"
```

**Expected:** `balance: 250.50` (0.00 + 250.50)

---

#### **Step 7: Test Insufficient Funds**

```bash
curl -X POST http://localhost:8080/wallets/transfer \
  -H "Authorization: Bearer $TOKEN1" \
  -H "Content-Type: application/json" \
  -d '{
    "to_wallet_id": 2,
    "amount": 999999.00
  }'
```

**Expected Response (409 Conflict):**
```json
{
  "error": "Transfer Failed",
  "message": "transfer failed: insufficient funds: balance 749.50, required 999999.00"
}
```

✅ **Balance Validation Works!**

---

#### **Step 8: Test Invalid Wallet ID**

```bash
curl -X POST http://localhost:8080/wallets/transfer \
  -H "Authorization: Bearer $TOKEN1" \
  -H "Content-Type: application/json" \
  -d '{
    "to_wallet_id": 99999,
    "amount": 10.00
  }'
```

**Expected Response (500 Internal Server Error):**
```json
{
  "error": "Transfer Failed",
  "message": "failed to get receiver wallet: record not found"
}
```

---

## 🛡️ Security Features Summary

| Feature | Implementation | Prevents |
|---------|----------------|----------|
| **Pessimistic Locking** | `SELECT FOR UPDATE` | Race conditions, double-spending |
| **ACID Transactions** | `db.Transaction()` | Partial updates, data corruption |
| **Balance Validation** | `if balance < amount` | Negative balances, overdrafts |
| **JWT Authentication** | Middleware | Unauthorized access |
| **Transaction Logging** | Audit table | Fraud detection, dispute resolution |
| **Row-Level Locking** | PostgreSQL | Concurrent modification conflicts |

---

## 📂 Module Structure

```
api/internal/modules/wallet/
├── domain/
│   └── entity.go              # Wallet, Transaction structs + Interfaces
├── repository/
│   └── postgres_repo.go       # Data access with Locking
├── usecase/
│   └── wallet_usecase.go      # Business logic with Transactions
├── delivery/
│   └── http/
│       ├── wallet_handler.go  # HTTP handlers
│       └── route.go           # Route registration
└── README.md                  # This file
```

---

## 🚀 Production Considerations

### **1. Rate Limiting**
Add rate limiting to prevent abuse:
```go
// Limit: 10 transfers per minute per user
middleware.RateLimiter(10, time.Minute)
```

### **2. Notification System**
Send notifications for:
- ✅ Successful transfers
- ❌ Failed transfers (insufficient funds)
- 🔔 Large transactions (> 10,000 THB)

### **3. Fraud Detection**
- Monitor unusual patterns (e.g., 100 transfers in 1 minute)
- Implement transaction limits (e.g., max 50,000 THB per day)

### **4. Database Indexing**
Ensure indexes exist:
```sql
CREATE INDEX idx_wallets_user_id ON wallets(user_id);
CREATE INDEX idx_transactions_from_wallet ON transactions(from_wallet_id);
CREATE INDEX idx_transactions_to_wallet ON transactions(to_wallet_id);
```

### **5. Monitoring**
Track metrics:
- Average transaction time
- Failed transaction rate
- Lock wait times
- Balance consistency checks

---

## 🎯 Future Enhancements

- [ ] **Multi-Currency Support** (USD, EUR, etc.)
- [ ] **Deposit/Withdrawal** endpoints
- [ ] **Transaction Fees** (configurable percentage)
- [ ] **Refund Mechanism**
- [ ] **Scheduled Transfers**
- [ ] **Wallet-to-Wallet QR Code**
- [ ] **Real-time WebSocket Updates**

---

## 📚 References

- [PostgreSQL Locking](https://www.postgresql.org/docs/current/explicit-locking.html)
- [GORM Transactions](https://gorm.io/docs/transactions.html)
- [ACID Compliance](https://en.wikipedia.org/wiki/ACID)

---

**Built with 💪 by the Komgrip Team for National-Scale Financial Systems.**

© 2026 Komgrip. All rights reserved.

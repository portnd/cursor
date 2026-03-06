package domain

import (
	"time"

	"gorm.io/gorm"
)

// MonthlyEntry is raw accounting data for one month (accountant input).
// One row per year-month; upsert by (year, month).
type MonthlyEntry struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Year        int            `gorm:"uniqueIndex:idx_finance_year_month;not null" json:"year"`
	Month       int            `gorm:"uniqueIndex:idx_finance_year_month;not null" json:"month"` // 1-12
	Revenue     float64        `gorm:"type:decimal(15,2);default:0;not null" json:"revenue"`
	Expenses    float64        `gorm:"type:decimal(15,2);default:0;not null" json:"expenses"`
	CashBalance float64        `gorm:"type:decimal(15,2);default:0;not null" json:"cash_balance"`
	Note        string         `gorm:"type:text" json:"note,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName overrides table name for GORM
func (MonthlyEntry) TableName() string {
	return "finance_monthly_entries"
}

// CreateOrUpdateRequest is the payload for upserting a monthly entry (accountant form).
type CreateOrUpdateRequest struct {
	Year        int     `json:"year" binding:"required,min=2020,max=2100"`
	Month       int     `json:"month" binding:"required,min=1,max=12"`
	Revenue     float64 `json:"revenue" binding:"gte=0"`
	Expenses    float64 `json:"expenses" binding:"gte=0"`
	CashBalance float64 `json:"cash_balance" binding:"gte=0"`
	Note        string  `json:"note"`
}

// FinanceSummary is computed for the CEO dashboard (runway, burn, revenue).
type FinanceSummary struct {
	CashBalance    float64 `json:"cash_balance"`
	RunwayMonths   float64 `json:"runway_months"`
	BurnRate       float64 `json:"burn_rate"`
	LastMonthMRR   float64 `json:"last_month_mrr"`
	NetNewARR      float64 `json:"net_new_arr"`
	Currency       string  `json:"currency"`
	LastEntryYear  int     `json:"last_entry_year"`
	LastEntryMonth int     `json:"last_entry_month"`
}

// Repository defines persistence for finance entries.
type Repository interface {
	CreateOrUpdate(entry *MonthlyEntry) error
	GetByYearMonth(year, month int) (*MonthlyEntry, error)
	List(limit int) ([]MonthlyEntry, error)
	GetLatestEntry() (*MonthlyEntry, error)
	GetEntriesForSummary(lastNMonths int) ([]MonthlyEntry, error)
}

// Usecase defines finance business logic (entry CRUD + summary computation).
type Usecase interface {
	CreateOrUpdateEntry(req *CreateOrUpdateRequest) (*MonthlyEntry, error)
	ListEntries(limit int) ([]MonthlyEntry, error)
	GetSummary() (*FinanceSummary, error)
}

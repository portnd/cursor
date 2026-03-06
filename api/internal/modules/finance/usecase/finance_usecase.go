package usecase

import (
	"math"

	"github.com/portnd/the-sentinel-core/internal/modules/finance/domain"
)

const (
	defaultSummaryMonths = 12
	currency             = "THB"
)

type financeUsecase struct {
	repo domain.Repository
}

// NewFinanceUsecase creates the finance usecase.
func NewFinanceUsecase(repo domain.Repository) domain.Usecase {
	return &financeUsecase{repo: repo}
}

func (u *financeUsecase) CreateOrUpdateEntry(req *domain.CreateOrUpdateRequest) (*domain.MonthlyEntry, error) {
	entry := &domain.MonthlyEntry{
		Year:        req.Year,
		Month:       req.Month,
		Revenue:     req.Revenue,
		Expenses:    req.Expenses,
		CashBalance: req.CashBalance,
		Note:        req.Note,
	}
	if err := u.repo.CreateOrUpdate(entry); err != nil {
		return nil, err
	}
	return u.repo.GetByYearMonth(req.Year, req.Month)
}

func (u *financeUsecase) ListEntries(limit int) ([]domain.MonthlyEntry, error) {
	return u.repo.List(limit)
}

func (u *financeUsecase) GetSummary() (*domain.FinanceSummary, error) {
	latest, err := u.repo.GetLatestEntry()
	if err != nil {
		return nil, err
	}
	entries, err := u.repo.GetEntriesForSummary(defaultSummaryMonths)
	if err != nil {
		return nil, err
	}

	out := &domain.FinanceSummary{Currency: currency}
	if latest != nil {
		out.CashBalance = latest.CashBalance
		out.LastMonthMRR = latest.Revenue
		out.LastEntryYear = latest.Year
		out.LastEntryMonth = latest.Month
	}

	if len(entries) > 0 {
		var totalExpenses float64
		for _, e := range entries {
			totalExpenses += e.Expenses
		}
		out.BurnRate = totalExpenses / float64(len(entries))
		if out.BurnRate > 0 {
			out.RunwayMonths = out.CashBalance / out.BurnRate
		}
		// Net new ARR: latest revenue - previous month revenue
		if len(entries) >= 2 {
			out.NetNewARR = entries[0].Revenue - entries[1].Revenue
		}
	}

	// Avoid negative runway or NaN
	if math.IsNaN(out.RunwayMonths) || math.IsInf(out.RunwayMonths, 0) || out.RunwayMonths < 0 {
		out.RunwayMonths = 0
	}
	return out, nil
}

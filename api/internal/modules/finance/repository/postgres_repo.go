package repository

import (
	"github.com/portnd/the-sentinel-core/internal/modules/finance/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type postgresRepo struct {
	db *gorm.DB
}

// NewPostgresRepository returns a finance repository.
func NewPostgresRepository(db *gorm.DB) domain.Repository {
	return &postgresRepo{db: db}
}

func (r *postgresRepo) CreateOrUpdate(entry *domain.MonthlyEntry) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "year"}, {Name: "month"}},
		DoUpdates: clause.AssignmentColumns([]string{"revenue", "expenses", "cash_balance", "note", "updated_at"}),
	}).Create(entry).Error
}

func (r *postgresRepo) GetByYearMonth(year, month int) (*domain.MonthlyEntry, error) {
	var e domain.MonthlyEntry
	err := r.db.Where("year = ? AND month = ?", year, month).First(&e).Error
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *postgresRepo) List(limit int) ([]domain.MonthlyEntry, error) {
	if limit <= 0 {
		limit = 24
	}
	if limit > 60 {
		limit = 60
	}
	var list []domain.MonthlyEntry
	err := r.db.Order("year DESC, month DESC").Limit(limit).Find(&list).Error
	return list, err
}

func (r *postgresRepo) GetLatestEntry() (*domain.MonthlyEntry, error) {
	var e domain.MonthlyEntry
	err := r.db.Order("year DESC, month DESC").First(&e).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *postgresRepo) GetEntriesForSummary(lastNMonths int) ([]domain.MonthlyEntry, error) {
	if lastNMonths <= 0 {
		lastNMonths = 12
	}
	var list []domain.MonthlyEntry
	err := r.db.Order("year DESC, month DESC").Limit(lastNMonths).Find(&list).Error
	return list, err
}

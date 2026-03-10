package repository

import (
	"fmt"
	"time"

	"github.com/portnd/the-sentinel-core/internal/modules/pulse/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type postgresRepository struct {
	db *gorm.DB
}

// NewPostgresRepository returns a PulseRepository backed by PostgreSQL.
func NewPostgresRepository(db *gorm.DB) domain.PulseRepository {
	return &postgresRepository{db: db}
}

// ─── Standup ──────────────────────────────────────────────────────────────────

func (r *postgresRepository) SaveStandup(s *domain.DailyStandup) error {
	// Upsert: one standup per user per calendar day.
	result := r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "date"}},
		DoUpdates: clause.AssignmentColumns([]string{"yesterday_summary", "today_task_ids", "blocker"}),
	}).Create(s)
	if result.Error != nil {
		return fmt.Errorf("pulse: save standup: %w", result.Error)
	}
	return nil
}

func (r *postgresRepository) GetStandupByUserAndDate(userID uint, date time.Time) (*domain.DailyStandup, error) {
	var s domain.DailyStandup
	dayStart := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	dayEnd := dayStart.Add(24 * time.Hour)
	err := r.db.
		Where("user_id = ? AND date >= ? AND date < ?", userID, dayStart, dayEnd).
		First(&s).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("pulse: get standup by user+date: %w", err)
	}
	return &s, nil
}

func (r *postgresRepository) GetStandupsByDate(date time.Time) ([]domain.DailyStandup, error) {
	dayStart := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	dayEnd := dayStart.Add(24 * time.Hour)

	var standups []domain.DailyStandup
	err := r.db.
		Where("date >= ? AND date < ?", dayStart, dayEnd).
		Find(&standups).Error
	if err != nil {
		return nil, fmt.Errorf("pulse: get standups by date: %w", err)
	}
	return standups, nil
}

// ─── Cross-module reads ───────────────────────────────────────────────────────

func (r *postgresRepository) GetTimeLogsByDate(date time.Time) ([]domain.TimeLogSummary, error) {
	dayStart := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	dayEnd := dayStart.Add(24 * time.Hour)

	type row struct {
		UserID      uint      `gorm:"column:user_id"`
		Minutes     int       `gorm:"column:minutes"`
		Description string    `gorm:"column:description"`
		LoggedAt    time.Time `gorm:"column:logged_at"`
	}
	var rows []row
	err := r.db.Raw(
		`SELECT user_id, minutes, description, logged_at
		 FROM time_logs
		 WHERE logged_at >= ? AND logged_at < ?
		 ORDER BY logged_at DESC`,
		dayStart, dayEnd,
	).Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("pulse: get time logs by date: %w", err)
	}

	out := make([]domain.TimeLogSummary, len(rows))
	for i, row := range rows {
		out[i] = domain.TimeLogSummary{
			UserID:      row.UserID,
			Minutes:     row.Minutes,
			Description: row.Description,
			LoggedAt:    row.LoggedAt,
		}
	}
	return out, nil
}

func (r *postgresRepository) GetSubmissionsByDate(date time.Time) ([]domain.SubmissionSummary, error) {
	dayStart := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	dayEnd := dayStart.Add(24 * time.Hour)

	type row struct {
		DevID        uint      `gorm:"column:dev_id"`
		ReferenceURL string    `gorm:"column:reference_url"`
		Note         string    `gorm:"column:note"`
		CreatedAt    time.Time `gorm:"column:created_at"`
	}
	var rows []row
	err := r.db.Raw(
		`SELECT dev_id, reference_url, note, created_at
		 FROM submissions
		 WHERE created_at >= ? AND created_at < ?
		 ORDER BY created_at DESC`,
		dayStart, dayEnd,
	).Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("pulse: get submissions by date: %w", err)
	}

	out := make([]domain.SubmissionSummary, len(rows))
	for i, row := range rows {
		out[i] = domain.SubmissionSummary{
			DevID:        row.DevID,
			ReferenceURL: row.ReferenceURL,
			Note:         row.Note,
			CreatedAt:    row.CreatedAt,
		}
	}
	return out, nil
}

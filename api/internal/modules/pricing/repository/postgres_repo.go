package repository

import (
	"fmt"
	"time"

	"github.com/portnd/the-sentinel-core/internal/modules/pricing/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type postgresRepository struct {
	db *gorm.DB
}

// NewPostgresRepository creates a new pricing repository backed by PostgreSQL.
func NewPostgresRepository(db *gorm.DB) domain.Repository {
	return &postgresRepository{db: db}
}

// GetCostConfig returns the singleton company cost configuration row.
func (r *postgresRepository) GetCostConfig() (*domain.CompanyCostConfig, error) {
	var cfg domain.CompanyCostConfig
	result := r.db.First(&cfg)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return &domain.CompanyCostConfig{
				WorkingDaysPerMonth: 22,
				WorkingHoursPerDay:  8,
				OverheadMultiplier:  1.30,
				DefaultProfitMargin: 0.25,
				DefaultRiskBuffer:   0.10,
				Currency:            "THB",
				ExecutiveExpense:    0,
				CompanyExpense:      0,
			}, nil
		}
		return nil, fmt.Errorf("pricing: fetch cost config: %w", result.Error)
	}
	return &cfg, nil
}

// UpsertCostConfig updates or creates the singleton company cost config (id=1).
func (r *postgresRepository) UpsertCostConfig(req *domain.UpdateCostConfigRequest) (*domain.CompanyCostConfig, error) {
	currency := req.Currency
	if currency == "" {
		currency = "THB"
	}
	cfg := domain.CompanyCostConfig{
		ID:                  1,
		WorkingDaysPerMonth: req.WorkingDaysPerMonth,
		WorkingHoursPerDay:  req.WorkingHoursPerDay,
		OverheadMultiplier:  req.OverheadMultiplier,
		DefaultProfitMargin: req.DefaultProfitMargin,
		DefaultRiskBuffer:   req.DefaultRiskBuffer,
		Currency:            currency,
		ExecutiveExpense:    req.ExecutiveExpense,
		CompanyExpense:      req.CompanyExpense,
	}
	if err := r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"working_days_per_month", "working_hours_per_day", "overhead_multiplier", "default_profit_margin", "default_risk_buffer", "currency", "executive_expense", "company_expense", "updated_at"}),
	}).Create(&cfg).Error; err != nil {
		return nil, fmt.Errorf("pricing: upsert cost config: %w", err)
	}
	return r.GetCostConfig()
}

// GetEmployeeSalaries fetches the currently active salary record for each requested user ID.
func (r *postgresRepository) GetEmployeeSalaries(userIDs []uint) ([]domain.EmployeeSalary, error) {
	if len(userIDs) == 0 {
		return nil, nil
	}
	today := time.Now().Truncate(24 * time.Hour)
	var salaries []domain.EmployeeSalary
	err := r.db.
		Where("user_id IN ? AND effective_from <= ? AND (effective_to IS NULL OR effective_to >= ?)", userIDs, today, today).
		Order("user_id, effective_from DESC").
		Find(&salaries).Error
	if err != nil {
		return nil, fmt.Errorf("pricing: fetch employee salaries: %w", err)
	}
	seen := make(map[uint]bool, len(salaries))
	deduplicated := make([]domain.EmployeeSalary, 0, len(userIDs))
	for _, s := range salaries {
		if !seen[s.UserID] {
			seen[s.UserID] = true
			deduplicated = append(deduplicated, s)
		}
	}
	return deduplicated, nil
}

// ListAllSalaries returns all salary records joined with user display info.
func (r *postgresRepository) ListAllSalaries() ([]domain.SalaryWithUser, error) {
	type row struct {
		ID             int64      `gorm:"column:id"`
		UserID         uint       `gorm:"column:user_id"`
		MonthlySalary  float64    `gorm:"column:monthly_salary"`
		Currency       string     `gorm:"column:currency"`
		EffectiveFrom  time.Time  `gorm:"column:effective_from"`
		EffectiveTo    *time.Time `gorm:"column:effective_to"`
		EmploymentType string     `gorm:"column:employment_type"`
		CostPerMinute  float64    `gorm:"column:cost_per_minute"`
		CreatedAt      time.Time  `gorm:"column:created_at"`
		UpdatedAt      time.Time  `gorm:"column:updated_at"`
		UserEmail       string     `gorm:"column:user_email"`
		UserDisplayName string     `gorm:"column:user_display_name"`
		UserRole        string     `gorm:"column:user_role"`
	}
	var rows []row
	err := r.db.Raw(`
		SELECT es.*, u.email AS user_email,
		       COALESCE(NULLIF(u.display_name,''), u.email) AS user_display_name,
		       u.role AS user_role
		FROM employee_salaries es
		JOIN users u ON u.id = es.user_id
		ORDER BY u.role, u.email, es.effective_from DESC
	`).Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("pricing: list salaries: %w", err)
	}
	result := make([]domain.SalaryWithUser, len(rows))
	for i, row := range rows {
		result[i] = domain.SalaryWithUser{
			EmployeeSalary: domain.EmployeeSalary{
				ID:             row.ID,
				UserID:         row.UserID,
				MonthlySalary:  row.MonthlySalary,
				Currency:       row.Currency,
				EffectiveFrom:  row.EffectiveFrom,
				EffectiveTo:    row.EffectiveTo,
				EmploymentType: row.EmploymentType,
				CostPerMinute:  row.CostPerMinute,
				CreatedAt:      row.CreatedAt,
				UpdatedAt:      row.UpdatedAt,
			},
			UserEmail:       row.UserEmail,
			UserDisplayName: row.UserDisplayName,
			UserRole:        row.UserRole,
			SsCost:          domain.SSCost(row.MonthlySalary),
		}
	}
	return result, nil
}

// UpsertSalary inserts or updates an employee salary record.
// Unique constraint is (user_id, effective_from) — on conflict update salary fields.
func (r *postgresRepository) UpsertSalary(req *domain.UpsertSalaryRequest) (*domain.EmployeeSalary, error) {
	currency := req.Currency
	if currency == "" {
		currency = "THB"
	}
	empType := req.EmploymentType
	if empType == "" {
		empType = "FULLTIME"
	}
	effectiveFrom, err := time.Parse("2006-01-02", req.EffectiveFrom)
	if err != nil {
		return nil, fmt.Errorf("pricing: invalid effective_from date: %w", err)
	}
	var effectiveTo *time.Time
	if req.EffectiveTo != "" {
		t, err := time.Parse("2006-01-02", req.EffectiveTo)
		if err != nil {
			return nil, fmt.Errorf("pricing: invalid effective_to date: %w", err)
		}
		effectiveTo = &t
	}

	sal := domain.EmployeeSalary{
		UserID:         req.UserID,
		MonthlySalary:  req.MonthlySalary,
		Currency:       currency,
		EffectiveFrom:  effectiveFrom,
		EffectiveTo:    effectiveTo,
		EmploymentType: empType,
	}
	if err := r.db.Clauses(clause.OnConflict{
		OnConstraint: "uq_employee_salaries_user_effective",
		DoUpdates:    clause.AssignmentColumns([]string{"monthly_salary", "currency", "effective_to", "employment_type", "updated_at"}),
	}).Create(&sal).Error; err != nil {
		return nil, fmt.Errorf("pricing: upsert salary: %w", err)
	}
	// Reload to get the auto-assigned ID.
	var saved domain.EmployeeSalary
	if err := r.db.Where("user_id = ? AND effective_from = ?", req.UserID, effectiveFrom).First(&saved).Error; err != nil {
		return nil, fmt.Errorf("pricing: reload salary after upsert: %w", err)
	}
	return &saved, nil
}

// DeleteSalary removes a salary record by primary key.
func (r *postgresRepository) DeleteSalary(id int64) error {
	if err := r.db.Delete(&domain.EmployeeSalary{}, id).Error; err != nil {
		return fmt.Errorf("pricing: delete salary %d: %w", id, err)
	}
	return nil
}

// GetActiveSalariesByRoles fetches currently active salary records for users whose role
// is in the given list. Used to compute company overhead from MANAGER + SUPPORT payroll.
func (r *postgresRepository) GetActiveSalariesByRoles(roles []string) ([]domain.EmployeeSalary, error) {
	if len(roles) == 0 {
		return nil, nil
	}
	today := time.Now().Truncate(24 * time.Hour)
	var salaries []domain.EmployeeSalary
	err := r.db.
		Joins("JOIN users u ON u.id = employee_salaries.user_id").
		Where("u.role IN ? AND employee_salaries.effective_from <= ? AND (employee_salaries.effective_to IS NULL OR employee_salaries.effective_to >= ?)", roles, today, today).
		Order("employee_salaries.user_id, employee_salaries.effective_from DESC").
		Find(&salaries).Error
	if err != nil {
		return nil, fmt.Errorf("pricing: fetch salaries by roles: %w", err)
	}
	// Deduplicate: keep only the latest active record per user
	seen := make(map[uint]bool)
	result := make([]domain.EmployeeSalary, 0, len(salaries))
	for _, s := range salaries {
		if !seen[s.UserID] {
			seen[s.UserID] = true
			result = append(result, s)
		}
	}
	return result, nil
}

// GetProjectTasks fetches tasks for a project. When taskIDs is non-empty only those tasks
// are returned; when epicIDs is non-empty only tasks belonging to those epics are returned.
// When both slices are empty all project tasks are included.
func (r *postgresRepository) GetProjectTasks(projectID string, taskIDs []string, epicIDs []string) ([]domain.PricingTask, error) {
	type row struct {
		ID        string     `gorm:"column:id"`
		Title     string     `gorm:"column:title"`
		EpicTitle string     `gorm:"column:epic_title"`
		StartDate *time.Time `gorm:"column:start_date"`
		EndDate   *time.Time `gorm:"column:end_date"`
	}

	baseSQL := `
		SELECT t.id::text, t.title, COALESCE(e.title, '') AS epic_title, t.start_date, t.end_date
		FROM tasks t
		LEFT JOIN epics e ON e.id = t.epic_id
		WHERE t.project_id = ?::uuid`

	var (
		rows []row
		err  error
	)

	if len(taskIDs) > 0 {
		err = r.db.Raw(baseSQL+` AND t.id::text IN ? ORDER BY e.sort_order NULLS LAST, t.sort_order`, projectID, taskIDs).Scan(&rows).Error
	} else if len(epicIDs) > 0 {
		err = r.db.Raw(baseSQL+` AND t.epic_id::text IN ? ORDER BY e.sort_order NULLS LAST, t.sort_order`, projectID, epicIDs).Scan(&rows).Error
	} else {
		err = r.db.Raw(baseSQL+` ORDER BY e.sort_order NULLS LAST, t.sort_order`, projectID).Scan(&rows).Error
	}

	if err != nil {
		return nil, fmt.Errorf("pricing: fetch project tasks: %w", err)
	}

	tasks := make([]domain.PricingTask, len(rows))
	for i, row := range rows {
		tasks[i] = domain.PricingTask{
			ID:        row.ID,
			Title:     row.Title,
			EpicTitle: row.EpicTitle,
			StartDate: row.StartDate,
			EndDate:   row.EndDate,
		}
	}
	return tasks, nil
}

// SaveSnapshot upserts a ProjectCostSnapshot.
func (r *postgresRepository) SaveSnapshot(snapshot *domain.ProjectCostSnapshot) error {
	if err := r.db.Save(snapshot).Error; err != nil {
		return fmt.Errorf("pricing: save snapshot: %w", err)
	}
	return nil
}

// GetAllProjectSnapshots returns all projects (from sentinel projects table) with their
// latest cost snapshot joined. When projectIDs is non-empty only those projects are included.
func (r *postgresRepository) GetAllProjectSnapshots(projectIDs []string) ([]domain.ProjectWithSnapshot, error) {
	type row struct {
		ProjectID   string  `gorm:"column:project_id"`
		ProjectName string  `gorm:"column:project_name"`
		Status      string  `gorm:"column:status"`
		// snapshot fields (nullable — project may have no snapshot)
		SnapshotID      *string  `gorm:"column:snapshot_id"`
		Version         *int     `gorm:"column:version"`
		TotalLaborCost  *float64 `gorm:"column:total_labor_cost"`
		TotalExpenses   *float64 `gorm:"column:total_expenses"`
		TotalCost       *float64 `gorm:"column:total_cost"`
		SuggestedPrice  *float64 `gorm:"column:suggested_price"`
		ProfitMargin    *float64 `gorm:"column:profit_margin"`
		RiskBuffer      *float64 `gorm:"column:risk_buffer"`
		EstimatedHours  *float64 `gorm:"column:estimated_hours"`
		EstimatedDays   *float64 `gorm:"column:estimated_days"`
		SnapshotStatus  *string  `gorm:"column:snapshot_status"`
		Notes           *string  `gorm:"column:notes"`
	}

	baseSQL := `
		SELECT
			p.id::text AS project_id,
			p.name AS project_name,
			p.status,
			s.id::text AS snapshot_id,
			s.version,
			s.total_labor_cost,
			s.total_expenses,
			s.total_cost,
			s.suggested_price,
			s.profit_margin,
			s.risk_buffer,
			s.estimated_hours,
			s.estimated_days,
			s.status AS snapshot_status,
			s.notes
		FROM projects p
		LEFT JOIN LATERAL (
			SELECT * FROM project_cost_snapshots
			WHERE project_id = p.id
			ORDER BY version DESC
			LIMIT 1
		) s ON true`

	var rows []row
	var err error
	if len(projectIDs) > 0 {
		err = r.db.Raw(baseSQL+` WHERE p.id::text IN ? ORDER BY p.name`, projectIDs).Scan(&rows).Error
	} else {
		err = r.db.Raw(baseSQL+` ORDER BY p.name`).Scan(&rows).Error
	}
	if err != nil {
		return nil, fmt.Errorf("pricing: get project snapshots: %w", err)
	}

	result := make([]domain.ProjectWithSnapshot, len(rows))
	for i, row := range rows {
		ps := domain.ProjectWithSnapshot{
			ProjectID:   row.ProjectID,
			ProjectName: row.ProjectName,
			Status:      row.Status,
		}
		if row.SnapshotID != nil {
			version := 1
			if row.Version != nil {
				version = *row.Version
			}
			snap := &domain.ProjectCostSnapshot{
				Version: version,
			}
			if row.TotalLaborCost != nil {
				snap.TotalLaborCost = *row.TotalLaborCost
			}
			if row.TotalExpenses != nil {
				snap.TotalExpenses = *row.TotalExpenses
			}
			if row.TotalCost != nil {
				snap.TotalCost = *row.TotalCost
			}
			if row.SuggestedPrice != nil {
				snap.SuggestedPrice = *row.SuggestedPrice
			}
			if row.ProfitMargin != nil {
				snap.ProfitMargin = *row.ProfitMargin
			}
			if row.RiskBuffer != nil {
				snap.RiskBuffer = *row.RiskBuffer
			}
			if row.EstimatedHours != nil {
				snap.EstimatedHours = *row.EstimatedHours
			}
			if row.EstimatedDays != nil {
				snap.EstimatedDays = *row.EstimatedDays
			}
			if row.SnapshotStatus != nil {
				snap.Status = *row.SnapshotStatus
			}
			if row.Notes != nil {
				snap.Notes = *row.Notes
			}
			ps.Snapshot = snap
		}
		result[i] = ps
	}
	return result, nil
}

// GetAllProjectExpenses returns project_expenses filtered by optional project IDs and date range.
func (r *postgresRepository) GetAllProjectExpenses(projectIDs []string, dateFrom, dateTo *time.Time) ([]domain.ProjectExpense, error) {
	query := r.db.Model(&domain.ProjectExpense{})
	if len(projectIDs) > 0 {
		query = query.Where("project_id::text IN ?", projectIDs)
	}
	if dateFrom != nil {
		query = query.Where("incurred_at >= ?", *dateFrom)
	}
	if dateTo != nil {
		query = query.Where("incurred_at <= ?", *dateTo)
	}
	var expenses []domain.ProjectExpense
	if err := query.Order("incurred_at DESC").Find(&expenses).Error; err != nil {
		return nil, fmt.Errorf("pricing: get project expenses: %w", err)
	}
	return expenses, nil
}

// GetProjectsWithExpenseNames returns a map of projectID → projectName for the given project IDs.
// Used to enrich expense rows with human-readable project names.
func (r *postgresRepository) GetProjectsWithExpenseNames(projectIDs []string) (map[string]string, error) {
	if len(projectIDs) == 0 {
		return map[string]string{}, nil
	}
	type row struct {
		ID   string `gorm:"column:id"`
		Name string `gorm:"column:name"`
	}
	var rows []row
	if err := r.db.Raw(`SELECT id::text AS id, name FROM projects WHERE id::text IN ?`, projectIDs).Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("pricing: get project names: %w", err)
	}
	m := make(map[string]string, len(rows))
	for _, r := range rows {
		m[r.ID] = r.Name
	}
	return m, nil
}

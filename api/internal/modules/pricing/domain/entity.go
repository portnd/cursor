package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// Social security (ประกันสังคม) constants for cost calculation.
const (
	SSRate   = 0.05  // Employee contribution rate (5%)
	SSCapTHB = 875.0 // Monthly cap in THB
)

// SSCost returns the monthly social security cost: Min(Salary * 0.05, 875).
func SSCost(monthlySalary float64) float64 {
	if monthlySalary <= 0 {
		return 0
	}
	v := monthlySalary * SSRate
	if v > SSCapTHB {
		return SSCapTHB
	}
	return v
}

// EmployeeSalary maps the employee_salaries table — effective-dated monthly salary per user.
type EmployeeSalary struct {
	ID             int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID         uint       `json:"user_id" gorm:"not null;index"`
	MonthlySalary  float64    `json:"monthly_salary" gorm:"type:decimal(15,2);not null"`
	Currency       string     `json:"currency" gorm:"type:varchar(3);default:'THB'"`
	EffectiveFrom  time.Time  `json:"effective_from" gorm:"type:date;not null"`
	EffectiveTo    *time.Time `json:"effective_to" gorm:"type:date"`
	EmploymentType string     `json:"employment_type" gorm:"type:varchar(20);default:'FULLTIME'"`
	CostPerMinute  float64    `json:"cost_per_minute" gorm:"type:decimal(10,6)"`
	CreatedAt      time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

func (EmployeeSalary) TableName() string { return "employee_salaries" }

// CompanyCostConfig is the singleton row controlling cost calculation settings.
type CompanyCostConfig struct {
	ID                   int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	WorkingDaysPerMonth  int       `json:"working_days_per_month" gorm:"default:22"`
	WorkingHoursPerDay   int       `json:"working_hours_per_day" gorm:"default:8"`
	OverheadMultiplier   float64   `json:"overhead_multiplier" gorm:"type:decimal(5,2);default:1.30"`
	DefaultProfitMargin  float64   `json:"default_profit_margin" gorm:"type:decimal(5,2);default:0.25"`
	DefaultRiskBuffer    float64   `json:"default_risk_buffer" gorm:"type:decimal(5,2);default:0.10"`
	Currency             string    `json:"currency" gorm:"type:varchar(3);default:'THB'"`
	ExecutiveExpense     float64   `json:"executive_expense" gorm:"type:decimal(15,2);default:0"`   // Default monthly executive cost (for quotation)
	CompanyExpense       float64   `json:"company_expense" gorm:"type:decimal(15,2);default:0"`     // Default monthly company overhead (rent, utilities, etc.)
	CreatedAt            time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt            time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (CompanyCostConfig) TableName() string { return "company_cost_configs" }

// ProjectCostSnapshot is a versioned, immutable cost estimate for a project.
type ProjectCostSnapshot struct {
	ID             uuid.UUID      `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ProjectID      uuid.UUID      `json:"project_id" gorm:"type:uuid;not null;index"`
	Version        int            `json:"version" gorm:"default:1"`
	TotalLaborCost float64        `json:"total_labor_cost" gorm:"type:decimal(15,2);default:0"`
	TotalExpenses  float64        `json:"total_expenses" gorm:"type:decimal(15,2);default:0"`
	TotalCost      float64        `json:"total_cost" gorm:"type:decimal(15,2);default:0"`
	SuggestedPrice float64        `json:"suggested_price" gorm:"type:decimal(15,2);default:0"`
	ProfitMargin   float64        `json:"profit_margin" gorm:"type:decimal(5,2);default:0.25"`
	RiskBuffer     float64        `json:"risk_buffer" gorm:"type:decimal(5,2);default:0.10"`
	EstimatedHours float64        `json:"estimated_hours" gorm:"type:decimal(10,2);default:0"`
	EstimatedDays  float64        `json:"estimated_days" gorm:"type:decimal(10,2);default:0"`
	Status         string         `json:"status" gorm:"type:varchar(20);default:'DRAFT'"`
	Breakdown      datatypes.JSON `json:"breakdown" gorm:"type:jsonb;default:'[]'"`
	Notes          string         `json:"notes" gorm:"type:text"`
	ValidUntil     *time.Time     `json:"valid_until" gorm:"type:date"`
	CreatedBy      *uint          `json:"created_by"`
	ApprovedBy     *uint          `json:"approved_by"`
	CreatedAt      time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
}

func (ProjectCostSnapshot) TableName() string { return "project_cost_snapshots" }

// ProjectExpense is a direct out-of-pocket expense for a project.
type ProjectExpense struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ProjectID   uuid.UUID  `json:"project_id" gorm:"type:uuid;not null;index"`
	Category    string     `json:"category" gorm:"type:varchar(50);not null"`
	Description string     `json:"description" gorm:"type:text"`
	Amount      float64    `json:"amount" gorm:"type:decimal(15,2);not null"`
	Currency    string     `json:"currency" gorm:"type:varchar(3);default:'THB'"`
	IncurredAt  time.Time  `json:"incurred_at" gorm:"type:date;not null"`
	CreatedBy   *uint      `json:"created_by"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

func (ProjectExpense) TableName() string { return "project_expenses" }

// ─── DTOs ────────────────────────────────────────────────────────────────────

// QuotationRequest is the input payload for a quotation calculation.
// Company overhead (company_expense, executive_expense, PM/MANAGER/SUPPORT payroll, total SS)
// is taken from config and DB; no manual overhead/pm/exec fields.
type QuotationRequest struct {
	// DevUserIDs lists the user IDs of developers whose salaries are used in the cost model.
	DevUserIDs []uint `json:"dev_user_ids" binding:"required,min=1"`
	// RiskMarginPct is e.g. 0.10 for 10% risk buffer.
	RiskMarginPct float64 `json:"risk_margin_pct" binding:"min=0"`
	// ProfitMarginPct is e.g. 0.25 for 25% profit.
	ProfitMarginPct float64 `json:"profit_margin_pct" binding:"min=0"`
	// TaskIDs optionally restricts cost calculation to specific task IDs.
	// When empty, all project tasks are included.
	TaskIDs []string `json:"task_ids"`
	// EpicIDs optionally restricts cost calculation to tasks belonging to specific epics.
	// Ignored when TaskIDs is provided.
	EpicIDs []string `json:"epic_ids"`
	// CustomerView controls which PDF template to use.
	// When true, the customer-facing template is used: no mandays, no internal cost breakdowns,
	// risk buffer and profit margin are absorbed into the total without disclosure.
	CustomerView bool `json:"customer_view"`
	// ProjectName is the human-readable project name shown in the customer PDF header.
	// Falls back to ProjectID when empty.
	ProjectName string `json:"project_name"`
}

// TaskCostLine is one task row in the quotation breakdown.
// EstDays is derived from task start_date and end_date (calendar days).
type TaskCostLine struct {
	TaskID        string  `json:"task_id"`
	Title         string  `json:"title"`
	EpicTitle     string  `json:"epic_title,omitempty"`
	EstDays       float64 `json:"est_days"`   // Number of days from start_date to end_date
	Mandays       float64 `json:"mandays"`
	CostPerManday float64 `json:"cost_per_manday"`
	Cost          float64 `json:"cost"`
}

// QuotationResponse is the full quotation result returned by CalculateQuotation.
type QuotationResponse struct {
	ProjectID     string         `json:"project_id"`
	Tasks         []TaskCostLine `json:"tasks"`
	Subtotal      float64        `json:"subtotal"`
	RiskAmount    float64        `json:"risk_amount"`
	ProfitAmount  float64        `json:"profit_amount"`
	VAT           float64        `json:"vat"`
	GrandTotal    float64        `json:"grand_total"`
	CostPerManday float64        `json:"cost_per_manday"`
	BillableDays  float64        `json:"billable_days"`
	TotalMandays  float64        `json:"total_mandays"`
	Currency      string         `json:"currency"`
}

// ─── Admin DTOs ──────────────────────────────────────────────────────────────

// UpsertSalaryRequest is the payload for creating or updating an employee salary record.
type UpsertSalaryRequest struct {
	UserID         uint    `json:"user_id" binding:"required"`
	MonthlySalary  float64 `json:"monthly_salary" binding:"required,min=0"`
	Currency       string  `json:"currency"`
	EffectiveFrom  string  `json:"effective_from" binding:"required"` // YYYY-MM-DD
	EffectiveTo    string  `json:"effective_to"`                      // YYYY-MM-DD or empty
	EmploymentType string  `json:"employment_type"`                   // FULLTIME / PARTTIME / CONTRACTOR
}

// UpdateCostConfigRequest is the payload for updating the singleton company cost config.
type UpdateCostConfigRequest struct {
	WorkingDaysPerMonth int     `json:"working_days_per_month" binding:"required,min=1,max=31"`
	WorkingHoursPerDay  int     `json:"working_hours_per_day" binding:"required,min=1,max=24"`
	OverheadMultiplier  float64 `json:"overhead_multiplier" binding:"required,min=1"`
	DefaultProfitMargin float64 `json:"default_profit_margin" binding:"min=0"`
	DefaultRiskBuffer   float64 `json:"default_risk_buffer" binding:"min=0"`
	Currency            string  `json:"currency"`
	ExecutiveExpense    float64 `json:"executive_expense" binding:"min=0"`
	CompanyExpense      float64 `json:"company_expense" binding:"min=0"`
}

// SalaryWithUser enriches an EmployeeSalary with the user's display fields for the admin UI.
// SsCost is computed as Min(monthly_salary * 0.05, 875) for display and costing transparency.
type SalaryWithUser struct {
	EmployeeSalary
	UserEmail       string  `json:"user_email"`
	UserDisplayName string  `json:"user_display_name"`
	UserRole        string  `json:"user_role"`
	SsCost          float64 `json:"ss_cost"` // Min(monthly_salary * 5%, 875 THB)
}

// ─── Cost Analysis Report DTOs ────────────────────────────────────────────────

// CostReportRequest is the input for generating a full Cost Analysis Report.
type CostReportRequest struct {
	// Period filter: "current_month" | "current_quarter" | "ytd" | "all" | "custom"
	Period     string  `json:"period"`
	DateFrom   string  `json:"date_from"`   // YYYY-MM-DD, used when period = "custom"
	DateTo     string  `json:"date_to"`     // YYYY-MM-DD, used when period = "custom"
	ProjectIDs []string `json:"project_ids"` // empty = all projects
}

// HeadcountRow is one employee row in the headcount analysis section.
type HeadcountRow struct {
	UserDisplayName      string  `json:"user_display_name"`
	UserEmail            string  `json:"user_email"`
	UserRole             string  `json:"user_role"`
	EmploymentType       string  `json:"employment_type"`
	MonthlySalary        float64 `json:"monthly_salary"`
	SsCost               float64 `json:"ss_cost"`
	FullyLoadedMonthly   float64 `json:"fully_loaded_monthly"`
	AnnualCost           float64 `json:"annual_cost"`
	EffectiveFrom        string  `json:"effective_from"`
}

// ProjectCostSummary is the cost analysis row for one project in the portfolio section.
type ProjectCostSummary struct {
	ProjectID    string  `json:"project_id"`
	ProjectName  string  `json:"project_name"`
	Status       string  `json:"status"`
	TotalMandays float64 `json:"total_mandays"`
	LaborCost    float64 `json:"labor_cost"`
	DirectExpenses float64 `json:"direct_expenses"`
	RiskBuffer   float64 `json:"risk_buffer"`
	Profit       float64 `json:"profit"`
	GrandTotal   float64 `json:"grand_total"`
	MarginPct    float64 `json:"margin_pct"`
	SnapshotStatus string `json:"snapshot_status"`
	Version      int     `json:"version"`
}

// ExpenseRow is one row in the direct expense breakdown section.
type ExpenseRow struct {
	ProjectName string  `json:"project_name"`
	Category    string  `json:"category"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
	IncurredAt  string  `json:"incurred_at"`
}

// ExpenseByCategoryRow is an aggregated expense total per category.
type ExpenseByCategoryRow struct {
	Category string  `json:"category"`
	Total    float64 `json:"total"`
}

// MonthlyCashFlowRow represents one month of cash outflow (salaries + expenses).
type MonthlyCashFlowRow struct {
	Month         string  `json:"month"`           // e.g. "Jan 2026"
	SalaryCost    float64 `json:"salary_cost"`
	ExpenseCost   float64 `json:"expense_cost"`
	TotalOutflow  float64 `json:"total_outflow"`
	Cumulative    float64 `json:"cumulative"`
}

// SensitivityCell is one cell in the risk×profit sensitivity matrix.
type SensitivityCell struct {
	RiskPct    float64 `json:"risk_pct"`
	ProfitPct  float64 `json:"profit_pct"`
	GrandTotal float64 `json:"grand_total"`
}

// CostReportData is the aggregated data structure for the Company Cost Analysis Report.
// Focus: monthly and annual cost structure — not project-specific.
type CostReportData struct {
	// Metadata
	GeneratedAt time.Time `json:"generated_at"`
	Currency    string    `json:"currency"`

	// ── Section 1: Monthly Cost Summary (KPI tiles) ──────────────────────────
	// Total cash out the door every month
	MonthlyTotalPayroll   float64 `json:"monthly_total_payroll"`    // gross salary all staff
	MonthlyTotalSS        float64 `json:"monthly_total_ss"`         // employer SS contribution
	MonthlyCompanyExpense float64 `json:"monthly_company_expense"`  // rent, utilities, etc.
	MonthlyExecExpense    float64 `json:"monthly_exec_expense"`     // executive costs
	MonthlyTotalOverhead  float64 `json:"monthly_total_overhead"`   // company + exec
	MonthlyBurnRate       float64 `json:"monthly_burn_rate"`        // total monthly cash out
	AnnualBurnRate        float64 `json:"annual_burn_rate"`         // × 12

	// ── Section 2: Cost Model Parameters ────────────────────────────────────
	Config          *CompanyCostConfig `json:"config"`
	CostPerManday   float64            `json:"cost_per_manday"`
	CostPerHour     float64            `json:"cost_per_hour"`
	BillableDays    float64            `json:"billable_days"`
	UtilizationRate float64            `json:"utilization_rate"`

	// Per-dev breakdown (for waterfall chart)
	AvgDevSalary         float64 `json:"avg_dev_salary"`
	AvgDevSS             float64 `json:"avg_dev_ss"`
	TotalPMSalaryPerDev  float64 `json:"total_pm_salary_per_dev"`
	CompanyExpensePerDev float64 `json:"company_expense_per_dev"`
	ExecExpensePerDev    float64 `json:"exec_expense_per_dev"`
	OverheadPerDev       float64 `json:"overhead_per_dev"`
	FullyLoadedMonthly   float64 `json:"fully_loaded_monthly"`

	// ── Section 3: Headcount & Salary Analysis ──────────────────────────────
	Headcount           []HeadcountRow `json:"headcount"`
	DevCount            int            `json:"dev_count"`
	PMCount             int            `json:"pm_count"`
	OtherCount          int            `json:"other_count"`
	TotalMonthlyPayroll float64        `json:"total_monthly_payroll"`
	TotalAnnualPayroll  float64        `json:"total_annual_payroll"`
	TotalMonthlySS      float64        `json:"total_monthly_ss"`
	TotalAnnualSS       float64        `json:"total_annual_ss"`

	// ── Section 4: Monthly Cost Waterfall (12-month projection) ─────────────
	// Fixed monthly cost components (same every month)
	MonthlySalaryPool   float64 `json:"monthly_salary_pool"`   // all salaries
	MonthlySSPool       float64 `json:"monthly_ss_pool"`       // all SS
	MonthlyCashFlow     []MonthlyCashFlowRow `json:"monthly_cash_flow"` // 12-month projection
	PeakMonthlyOutflow  float64 `json:"peak_monthly_outflow"`
	AnnualProjection    float64 `json:"annual_projection"`     // exact 12-month total

	// ── Section 5: Sensitivity Analysis ─────────────────────────────────────
	SensitivityMatrix    []SensitivityCell `json:"sensitivity_matrix"`
	BreakEvenRate        float64           `json:"break_even_rate"` // cost/manday minimum
}

// ─── Ports (interfaces) ───────────────────────────────────────────────────────

// PricingTask is the minimal task data needed for cost calculation.
// EstDays is computed from start_date and end_date in the usecase.
type PricingTask struct {
	ID        string
	Title     string
	EpicTitle string
	StartDate *time.Time
	EndDate   *time.Time
}

// ProjectWithSnapshot aggregates a project with its latest cost snapshot.
type ProjectWithSnapshot struct {
	ProjectID   string
	ProjectName string
	Status      string
	Snapshot    *ProjectCostSnapshot
}

// Repository is the data-access port for the pricing module.
type Repository interface {
	GetCostConfig() (*CompanyCostConfig, error)
	UpsertCostConfig(req *UpdateCostConfigRequest) (*CompanyCostConfig, error)
	GetEmployeeSalaries(userIDs []uint) ([]EmployeeSalary, error)
	ListAllSalaries() ([]SalaryWithUser, error)
	UpsertSalary(req *UpsertSalaryRequest) (*EmployeeSalary, error)
	DeleteSalary(id int64) error
	// GetActiveSalariesByRoles fetches currently active salary records for users with the given roles.
	// Used to sum MANAGER + SUPPORT payroll as part of company overhead.
	GetActiveSalariesByRoles(roles []string) ([]EmployeeSalary, error)
	// GetProjectTasks fetches tasks for a project. When taskIDs is non-empty only those tasks
	// are returned; when epicIDs is non-empty only tasks belonging to those epics are returned.
	// When both are empty all project tasks are returned.
	GetProjectTasks(projectID string, taskIDs []string, epicIDs []string) ([]PricingTask, error)
	SaveSnapshot(snapshot *ProjectCostSnapshot) error
	// Cost report aggregation queries
	GetAllProjectSnapshots(projectIDs []string) ([]ProjectWithSnapshot, error)
	GetAllProjectExpenses(projectIDs []string, dateFrom, dateTo *time.Time) ([]ProjectExpense, error)
	GetProjectsWithExpenseNames(expenseIDs []string) (map[string]string, error)
}

// CompanyMandayRateResponse is the response for the company manday rate calculation.
// It shows the fully loaded cost per manday derived from ALL employee salaries + company overheads.
//
// Company Expense Total = company_expense (config) + executive_expense (config)
//                       + Σ monthly_salary of MANAGER role
//                       + Σ monthly_salary of SUPPORT role
type CompanyMandayRateResponse struct {
	// Salary pool
	TotalMonthlySalaries float64 `json:"total_monthly_salaries"` // sum of all active employee salaries
	TotalMonthlySS       float64 `json:"total_monthly_ss"`        // sum of all SS contributions
	// Overheads from company config
	CompanyExpense   float64 `json:"company_expense"`    // rent, utilities, etc. (from config)
	ExecutiveExpense float64 `json:"executive_expense"`  // executive cost (from config)
	// Non-billable role salaries treated as overhead
	OverheadRoleSalaryTotal float64 `json:"overhead_role_salary_total"` // Σ salary of MANAGER + SUPPORT
	// Combined company expense for manday calculation
	CompanyExpenseTotal float64 `json:"company_expense_total"` // company_expense + executive_expense + overhead_role_salary_total
	// Total monthly burden
	TotalMonthlyBurnRate float64 `json:"total_monthly_burn_rate"` // all salaries + SS + config overheads
	// Headcount
	ActiveHeadcount int `json:"active_headcount"` // number of active employees with salary records
	DevCount        int `json:"dev_count"`         // number of billable engineers (ENGINEER + CHIEF_ENGINEER)
	// Per-dev breakdown (authoritative — matches CostPerManday calculation)
	AvgDevSalary           float64 `json:"avg_dev_salary"`             // average dev salary = totalDevSalary / devCount
	OverheadPerDev         float64 `json:"overhead_per_dev"`           // companyExpenseTotal / devCount
	FullyLoadedMonthlyPerDev float64 `json:"fully_loaded_monthly_per_dev"` // avgDevSalary + overheadPerDev
	// Config
	WorkingDaysPerMonth int     `json:"working_days_per_month"`
	OverheadMultiplier  float64 `json:"overhead_multiplier"`
	BillableDays        float64 `json:"billable_days"` // working_days / overhead_multiplier
	// Rates
	CostPerManday float64 `json:"cost_per_manday"` // fully_loaded_monthly_per_dev / billable_days
	CostPerHour   float64 `json:"cost_per_hour"`   // cost_per_manday / working_hours_per_day
	Currency      string  `json:"currency"`
}

// MADeliveryMilestone represents one payment/delivery milestone in an MA quotation.
type MADeliveryMilestone struct {
	Label     string  `json:"label"`
	TaskCount *int    `json:"task_count"` // nil = end-of-MA milestone (no task delivery)
	Amount    float64 `json:"amount"`
	IsEndOfMA bool    `json:"is_end_of_ma"`
}

// MATaskItem is a task included in the MA scope.
type MATaskItem struct {
	Code  string `json:"code"`
	Title string `json:"title"`
}

// MAQuotationExportRequest is the payload for POST /ma-quotation/export.
type MAQuotationExportRequest struct {
	ProjectName     string                `json:"project_name"`
	QuoteNo         string                `json:"quote_no"`
	IssueDate       string                `json:"issue_date"`
	MAPrice         float64               `json:"ma_price"`
	MADurationYears int                   `json:"ma_duration_years"`
	Tasks           []MATaskItem          `json:"tasks"`
	Milestones      []MADeliveryMilestone `json:"milestones"`
}

// Usecase is the business-logic port for the pricing module.
type Usecase interface {
	CalculateQuotation(projectID string, req *QuotationRequest) (*QuotationResponse, error)
	ExportQuotationPDF(projectID string, req *QuotationRequest) ([]byte, error)
	// MA Quotation PDF (HTML → chromedp)
	ExportMAQuotationPDF(req *MAQuotationExportRequest) ([]byte, error)
	// Cost Analysis Report
	GenerateCostReport(req *CostReportRequest) ([]byte, error)
	// Company Manday Rate — fully loaded rate derived from all salaries + overheads
	GetCompanyMandayRate() (*CompanyMandayRateResponse, error)
	// Admin
	GetCostConfig() (*CompanyCostConfig, error)
	UpdateCostConfig(req *UpdateCostConfigRequest) (*CompanyCostConfig, error)
	ListAllSalaries() ([]SalaryWithUser, error)
	UpsertSalary(req *UpsertSalaryRequest) (*EmployeeSalary, error)
	DeleteSalary(id int64) error
}

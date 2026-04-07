package usecase

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	chromepdf "github.com/portnd/the-sentinel-core/internal/core/pdf"
	"github.com/portnd/the-sentinel-core/internal/modules/pricing/domain"
)

const vatRate = 0.07

type costingUsecase struct {
	repo domain.Repository
}

// NewCostingUsecase constructs the pricing usecase.
func NewCostingUsecase(repo domain.Repository) domain.Usecase {
	return &costingUsecase{repo: repo}
}

// ─── Admin usecase methods ────────────────────────────────────────────────────

func (u *costingUsecase) GetCostConfig() (*domain.CompanyCostConfig, error) {
	return u.repo.GetCostConfig()
}

func (u *costingUsecase) UpdateCostConfig(req *domain.UpdateCostConfigRequest) (*domain.CompanyCostConfig, error) {
	return u.repo.UpsertCostConfig(req)
}

func (u *costingUsecase) ListAllSalaries() ([]domain.SalaryWithUser, error) {
	return u.repo.ListAllSalaries()
}

func (u *costingUsecase) UpsertSalary(req *domain.UpsertSalaryRequest) (*domain.EmployeeSalary, error) {
	return u.repo.UpsertSalary(req)
}

func (u *costingUsecase) DeleteSalary(id int64) error {
	return u.repo.DeleteSalary(id)
}

// GetCompanyMandayRate computes the company-wide fully loaded cost per manday.
//
// CompanyExpenseTotal = company_expense + executive_expense
//                     + Σ monthly_salary (role = MANAGER)
//                     + Σ monthly_salary (role = SUPPORT)
//                     + Total SS/mo (Σ SS of all employees)
//
// These roles (MANAGER, SUPPORT) are non-billable overhead — their time is not
// sold to clients. SS is also treated as company overhead, not per-dev cost.
//
// Formula:
//
//	DevCount      = employees with role = DEV
//	AvgDevSalary  = Σ DEV salary / DevCount
//	OverheadPerDev = CompanyExpenseTotal / DevCount
//	FullyLoadedMonthly = AvgDevSalary + OverheadPerDev
//	BillableDays  = working_days_per_month / overhead_multiplier
//	CostPerManday = FullyLoadedMonthly / BillableDays
func (u *costingUsecase) GetCompanyMandayRate() (*domain.CompanyMandayRateResponse, error) {
	cfg, err := u.repo.GetCostConfig()
	if err != nil {
		return nil, fmt.Errorf("manday rate: get config: %w", err)
	}

	// All employee salaries (for total payroll display)
	allSalaries, err := u.repo.ListAllSalaries()
	if err != nil {
		return nil, fmt.Errorf("manday rate: list salaries: %w", err)
	}

	// Product Owner, MANAGER, SUPPORT salaries counted as overhead, not billable engineer cost
	overheadRoleSalaries, err := u.repo.GetActiveSalariesByRoles([]string{"PRODUCT_OWNER", "MANAGER", "SUPPORT"})
	if err != nil {
		return nil, fmt.Errorf("manday rate: fetch overhead role salaries: %w", err)
	}

	// DEV-only salaries for avg dev cost
	devSalaries, err := u.repo.GetActiveSalariesByRoles([]string{"ENGINEER", "CHIEF_ENGINEER"})
	if err != nil {
		return nil, fmt.Errorf("manday rate: fetch dev salaries: %w", err)
	}

	var totalSalary, totalSS float64
	for _, s := range allSalaries {
		totalSalary += s.MonthlySalary
		totalSS += domain.SSCost(s.MonthlySalary)
	}

	var overheadRoleSalarySum float64
	for _, s := range overheadRoleSalaries {
		overheadRoleSalarySum += s.MonthlySalary
	}

	var totalDevSalary float64
	for _, s := range devSalaries {
		totalDevSalary += s.MonthlySalary
	}

	devCount := len(devSalaries)
	if devCount == 0 {
		devCount = 1
	}

	// SS is company overhead, not per-dev
	companyExpenseTotal := cfg.CompanyExpense + cfg.ExecutiveExpense + overheadRoleSalarySum + totalSS

	billableDays := float64(cfg.WorkingDaysPerMonth) / cfg.OverheadMultiplier
	if billableDays <= 0 {
		billableDays = 1
	}

	avgDevSalary := totalDevSalary / float64(devCount)
	overheadPerDev := companyExpenseTotal / float64(devCount)
	fullyLoadedMonthly := avgDevSalary + overheadPerDev
	costPerManday := fullyLoadedMonthly / billableDays

	costPerHour := costPerManday
	if cfg.WorkingHoursPerDay > 0 {
		costPerHour = costPerManday / float64(cfg.WorkingHoursPerDay)
	}

	// TotalBurnRate = all salaries + SS + config overheads (for reporting)
	totalBurnRate := totalSalary + totalSS + cfg.CompanyExpense + cfg.ExecutiveExpense

	return &domain.CompanyMandayRateResponse{
		TotalMonthlySalaries:         round2(totalSalary),
		TotalMonthlySS:               round2(totalSS),
		CompanyExpense:               cfg.CompanyExpense,
		ExecutiveExpense:             cfg.ExecutiveExpense,
		OverheadRoleSalaryTotal:      round2(overheadRoleSalarySum),
		CompanyExpenseTotal:          round2(companyExpenseTotal),
		TotalMonthlyBurnRate:         round2(totalBurnRate),
		ActiveHeadcount:              len(allSalaries),
		DevCount:                     devCount,
		AvgDevSalary:                 round2(avgDevSalary),
		OverheadPerDev:               round2(overheadPerDev),
		FullyLoadedMonthlyPerDev:     round2(fullyLoadedMonthly),
		WorkingDaysPerMonth:          cfg.WorkingDaysPerMonth,
		OverheadMultiplier:           cfg.OverheadMultiplier,
		BillableDays:                 round2(billableDays),
		CostPerManday:                round2(costPerManday),
		CostPerHour:                  round2(costPerHour),
		Currency:                     cfg.Currency,
	}, nil
}

// CalculateQuotation implements the Fully Loaded Cost model.
func (u *costingUsecase) CalculateQuotation(projectID string, req *domain.QuotationRequest) (*domain.QuotationResponse, error) {
	cfg, err := u.repo.GetCostConfig()
	if err != nil {
		return nil, fmt.Errorf("costing: get config: %w", err)
	}

	salaries, err := u.repo.GetEmployeeSalaries(req.DevUserIDs)
	if err != nil {
		return nil, fmt.Errorf("costing: get salaries: %w", err)
	}

	// Product Owner, MANAGER, SUPPORT salaries are overhead (from DB; no manual input)
	overheadRoleSalaries, err := u.repo.GetActiveSalariesByRoles([]string{"PRODUCT_OWNER", "MANAGER", "SUPPORT"})
	if err != nil {
		return nil, fmt.Errorf("costing: get overhead role salaries: %w", err)
	}

	// All employee salaries — needed to compute total SS as company overhead
	allSalaries, err := u.repo.ListAllSalaries()
	if err != nil {
		return nil, fmt.Errorf("costing: list all salaries: %w", err)
	}

	tasks, err := u.repo.GetProjectTasks(projectID, req.TaskIDs, req.EpicIDs)
	if err != nil {
		return nil, fmt.Errorf("costing: get tasks: %w", err)
	}
	if len(tasks) == 0 {
		return nil, fmt.Errorf("costing: no tasks found for project %s", projectID)
	}

	costPerManday := u.computeCostPerManday(cfg, req, salaries, overheadRoleSalaries, allSalaries)

	// BillableDays = WorkingDaysPerMonth / OverheadMultiplier
	// e.g. 22 / 1.30 ≈ 16.92 billable days
	billableDays := float64(cfg.WorkingDaysPerMonth) / cfg.OverheadMultiplier

	lines := make([]domain.TaskCostLine, 0, len(tasks))
	var laborSubtotal float64
	var totalMandays float64

	for _, t := range tasks {
		mandays := daysFromStartEnd(t.StartDate, t.EndDate)
		cost := round2(mandays * costPerManday)
		totalMandays += mandays
		laborSubtotal += cost
		lines = append(lines, domain.TaskCostLine{
			TaskID:        t.ID,
			Title:         t.Title,
			EpicTitle:     t.EpicTitle,
			EstDays:       round2(mandays),
			Mandays:       round2(mandays),
			CostPerManday: round2(costPerManday),
			Cost:          cost,
		})
	}

	subtotal := round2(laborSubtotal)
	riskAmount := round2(subtotal * req.RiskMarginPct)
	afterRisk := subtotal + riskAmount
	profitAmount := round2(afterRisk * req.ProfitMarginPct)
	afterProfit := afterRisk + profitAmount
	vat := round2(afterProfit * vatRate)
	grandTotal := round2(afterProfit + vat)

	return &domain.QuotationResponse{
		ProjectID:     projectID,
		Tasks:         lines,
		Subtotal:      subtotal,
		RiskAmount:    riskAmount,
		ProfitAmount:  profitAmount,
		VAT:           vat,
		GrandTotal:    grandTotal,
		CostPerManday: round2(costPerManday),
		BillableDays:  round2(billableDays),
		TotalMandays:  round2(totalMandays),
		Currency:      cfg.Currency,
	}, nil
}

// ExportQuotationPDF calculates the quotation and renders it to a PDF via chromedp.
// When req.CustomerView is true the customer-facing template is used (no mandays, no margin disclosure).
func (u *costingUsecase) ExportQuotationPDF(projectID string, req *domain.QuotationRequest) ([]byte, error) {
	result, err := u.CalculateQuotation(projectID, req)
	if err != nil {
		return nil, err
	}

	var html string
	if req.CustomerView {
		html = buildCustomerQuotationHTML(result, req)
	} else {
		html = buildQuotationHTML(result, req)
	}

	ctx, cancel := chromepdf.NewChromedpContext(context.Background())
	defer cancel()

	var buf []byte
	if err := chromedp.Run(ctx, chromepdf.PrintToPDF(html, &buf, false)); err != nil {
		return nil, fmt.Errorf("costing: render PDF: %w", err)
	}
	return buf, nil
}

// ─── Internal helpers ─────────────────────────────────────────────────────────

// computeCostPerManday derives the fully loaded daily cost from config + DB.
//
// CompanyExpenseTotal = cfg.CompanyExpense + cfg.ExecutiveExpense  (from config)
//                     + Σ monthly_salary of Product Owner + MANAGER + SUPPORT (from DB)
//                     + Total SS/mo (Σ SS of ALL employees with salary records)
//
// OverheadPerDev   = CompanyExpenseTotal / len(DevUserIDs)
// FullyLoadedCost  = AvgSalary + OverheadPerDev  (monthly)
// CostPerManday    = FullyLoadedCost / BillableDays
func (u *costingUsecase) computeCostPerManday(cfg *domain.CompanyCostConfig, req *domain.QuotationRequest, salaries []domain.EmployeeSalary, overheadRoleSalaries []domain.EmployeeSalary, allSalaries []domain.SalaryWithUser) float64 {
	numDevs := float64(len(req.DevUserIDs))
	if numDevs == 0 {
		numDevs = 1
	}

	// Salary of the selected devs only
	var totalDevSalary float64
	for _, s := range salaries {
		totalDevSalary += s.MonthlySalary
	}
	avgSalary := totalDevSalary / numDevs

	// Sum MANAGER + SUPPORT salaries as part of overhead
	var overheadRoleSalarySum float64
	for _, s := range overheadRoleSalaries {
		overheadRoleSalarySum += s.MonthlySalary
	}

	// SS of ALL employees goes into company overhead
	var totalSS float64
	for _, s := range allSalaries {
		totalSS += domain.SSCost(s.MonthlySalary)
	}

	// Company overhead: config (company_expense, executive_expense) + Product Owner/MANAGER/SUPPORT payroll + total SS
	companyExpenseTotal := cfg.CompanyExpense + cfg.ExecutiveExpense + overheadRoleSalarySum + totalSS
	overheadPerDev := companyExpenseTotal / numDevs
	fullyLoadedMonthly := avgSalary + overheadPerDev

	billableDays := float64(cfg.WorkingDaysPerMonth) / cfg.OverheadMultiplier
	if billableDays == 0 {
		billableDays = 1
	}
	return fullyLoadedMonthly / billableDays
}

// daysFromStartEnd returns the number of working days (Mon–Fri, excluding Thai public holidays)
// between start and end inclusive. Returns 0 if either date is nil or end < start.
func daysFromStartEnd(start, end *time.Time) float64 {
	if start == nil || end == nil {
		return 0
	}
	s := truncateToDay(*start)
	e := truncateToDay(*end)
	if e.Before(s) {
		return 0
	}
	count := 0
	for d := s; !d.After(e); d = d.AddDate(0, 0, 1) {
		if isWorkingDay(d) {
			count++
		}
	}
	return float64(count)
}

// truncateToDay zeroes out the time component keeping the date.
func truncateToDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// isWorkingDay returns true if the given day is Mon–Fri and not a Thai public holiday.
func isWorkingDay(t time.Time) bool {
	wd := t.Weekday()
	if wd == time.Saturday || wd == time.Sunday {
		return false
	}
	return !isThaiPublicHoliday(t)
}

// isThaiPublicHoliday returns true only for New Year's Day (1 Jan).
// All other weekdays are counted as working days for costing.
func isThaiPublicHoliday(t time.Time) bool {
	m, d := t.Month(), t.Day()
	return m == time.January && d == 1
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}

// ─── Cost Analysis Report ──────────────────────────────────────────────────────

// GenerateCostReport aggregates all financial data and exports a CFO-grade multi-section PDF.
func (u *costingUsecase) GenerateCostReport(req *domain.CostReportRequest) ([]byte, error) {
	data, err := u.buildCostReportData(req)
	if err != nil {
		return nil, fmt.Errorf("costing: build report data: %w", err)
	}

	html := buildCostReportHTML(data)

	ctx, cancel := chromepdf.NewChromedpContext(context.Background())
	defer cancel()

	var buf []byte
	if err := chromedp.Run(ctx, chromepdf.PrintToPDF(html, &buf, false)); err != nil {
		return nil, fmt.Errorf("costing: render cost report PDF: %w", err)
	}
	return buf, nil
}

// buildCostReportData aggregates company cost structure data for monthly/annual analysis.
// It uses only salary records and company_cost_configs — no project data needed.
func (u *costingUsecase) buildCostReportData(req *domain.CostReportRequest) (*domain.CostReportData, error) {
	cfg, err := u.repo.GetCostConfig()
	if err != nil {
		return nil, err
	}

	salaries, err := u.repo.ListAllSalaries()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	data := &domain.CostReportData{
		GeneratedAt: now,
		Currency:    cfg.Currency,
		Config:      cfg,
	}

	// ── Classify employees ────────────────────────────────────────────────────
	var devSalaries, pmSalaries, otherSalaries []domain.SalaryWithUser
	for _, s := range salaries {
		r := strings.ToUpper(s.UserRole)
		switch r {
		case "DEV", "DEVELOPER", "ENGINEER", "CHIEF_ENGINEER":
			devSalaries = append(devSalaries, s)
		case "PM", "PRODUCT_OWNER":
			pmSalaries = append(pmSalaries, s)
		default:
			otherSalaries = append(otherSalaries, s)
		}
	}
	// Fallback: if no explicit engineer role found, non-CEO/Product Owner/EXEC are treated as engineers
	if len(devSalaries) == 0 {
		for _, s := range salaries {
			r := strings.ToUpper(s.UserRole)
			if r != "CEO" && r != "CTO" && r != "PM" && r != "PRODUCT_OWNER" && r != "EXEC" {
				devSalaries = append(devSalaries, s)
			}
		}
	}
	numDevs := float64(len(devSalaries))
	if numDevs == 0 {
		numDevs = 1
	}

	// ── Section 1: Monthly Cost Summary ──────────────────────────────────────
	var totalPayroll, totalSS float64
	for _, s := range salaries {
		totalPayroll += s.MonthlySalary
		totalSS += s.SsCost
	}
	var totalPMSalary float64
	for _, s := range pmSalaries {
		totalPMSalary += s.MonthlySalary
	}

	monthlyCompanyExpense := cfg.CompanyExpense
	monthlyExecExpense    := cfg.ExecutiveExpense
	monthlyTotalOverhead  := monthlyCompanyExpense + monthlyExecExpense
	monthlyBurnRate       := totalPayroll + totalSS + monthlyTotalOverhead
	annualBurnRate        := monthlyBurnRate * 12

	data.MonthlyTotalPayroll   = round2(totalPayroll)
	data.MonthlyTotalSS        = round2(totalSS)
	data.MonthlyCompanyExpense = round2(monthlyCompanyExpense)
	data.MonthlyExecExpense    = round2(monthlyExecExpense)
	data.MonthlyTotalOverhead  = round2(monthlyTotalOverhead)
	data.MonthlyBurnRate       = round2(monthlyBurnRate)
	data.AnnualBurnRate        = round2(annualBurnRate)

	// ── Section 2: Cost Model Parameters ─────────────────────────────────────
	billableDays := float64(cfg.WorkingDaysPerMonth) / cfg.OverheadMultiplier

	var totalDevSalary, totalDevSS float64
	for _, s := range devSalaries {
		totalDevSalary += s.MonthlySalary
		totalDevSS += s.SsCost
	}
	avgDevSalary := totalDevSalary / numDevs
	avgDevSS     := totalDevSS / numDevs
	overheadPerDev := (monthlyCompanyExpense + monthlyExecExpense + totalPMSalary) / numDevs
	fullyLoadedMonthly := avgDevSalary + avgDevSS + overheadPerDev
	costPerManday := 0.0
	if billableDays > 0 {
		costPerManday = fullyLoadedMonthly / billableDays
	}
	costPerHour := 0.0
	if cfg.WorkingHoursPerDay > 0 {
		costPerHour = costPerManday / float64(cfg.WorkingHoursPerDay)
	}

	data.BillableDays         = round2(billableDays)
	data.UtilizationRate      = round2(1.0 / cfg.OverheadMultiplier)
	data.CostPerManday        = round2(costPerManday)
	data.CostPerHour          = round2(costPerHour)
	data.AvgDevSalary         = round2(avgDevSalary)
	data.AvgDevSS             = round2(avgDevSS)
	data.TotalPMSalaryPerDev  = round2(totalPMSalary / numDevs)
	data.CompanyExpensePerDev = round2(monthlyCompanyExpense / numDevs)
	data.ExecExpensePerDev    = round2(monthlyExecExpense / numDevs)
	data.OverheadPerDev       = round2(overheadPerDev)
	data.FullyLoadedMonthly   = round2(fullyLoadedMonthly)

	// ── Section 3: Headcount & Salary Analysis ───────────────────────────────
	var totalAnnualPayroll float64
	headcount := make([]domain.HeadcountRow, 0, len(salaries))
	for _, s := range salaries {
		ss := s.SsCost
		fullyLoaded := s.MonthlySalary + ss
		annual := fullyLoaded * 12
		totalAnnualPayroll += annual
		headcount = append(headcount, domain.HeadcountRow{
			UserDisplayName:    s.UserDisplayName,
			UserEmail:          s.UserEmail,
			UserRole:           s.UserRole,
			EmploymentType:     s.EmploymentType,
			MonthlySalary:      round2(s.MonthlySalary),
			SsCost:             round2(ss),
			FullyLoadedMonthly: round2(fullyLoaded),
			AnnualCost:         round2(annual),
			EffectiveFrom:      s.EffectiveFrom.Format("02 Jan 2006"),
		})
	}
	data.Headcount           = headcount
	data.DevCount            = len(devSalaries)
	data.PMCount             = len(pmSalaries)
	data.OtherCount          = len(otherSalaries)
	data.TotalMonthlyPayroll = round2(totalPayroll)
	data.TotalAnnualPayroll  = round2(totalAnnualPayroll)
	data.TotalMonthlySS      = round2(totalSS)
	data.TotalAnnualSS       = round2(totalSS * 12)

	// ── Section 4: 12-Month Cash Flow Projection ──────────────────────────────
	// Project forward 12 months from the current month using fixed monthly costs.
	data.MonthlySalaryPool = round2(totalPayroll)
	data.MonthlySSPool     = round2(totalSS)

	cashFlow := make([]domain.MonthlyCashFlowRow, 12)
	var cumulative float64
	startMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	for i := 0; i < 12; i++ {
		m := startMonth.AddDate(0, i, 0)
		expCost := monthlyTotalOverhead
		total := monthlyBurnRate
		cumulative += total
		cashFlow[i] = domain.MonthlyCashFlowRow{
			Month:        m.Format("Jan 2006"),
			SalaryCost:   round2(totalPayroll + totalSS),
			ExpenseCost:  round2(expCost),
			TotalOutflow: round2(total),
			Cumulative:   round2(cumulative),
		}
	}
	data.MonthlyCashFlow    = cashFlow
	data.PeakMonthlyOutflow = round2(monthlyBurnRate) // constant, same every month
	data.AnnualProjection   = round2(cumulative)

	// ── Section 5: Sensitivity Analysis ──────────────────────────────────────
	// Base = fully loaded monthly cost × 22 billable days (1 full month of work)
	riskSteps   := []float64{0.05, 0.10, 0.15, 0.20}
	profitSteps := []float64{0.15, 0.20, 0.25, 0.30, 0.35}
	baseLaborCost := fullyLoadedMonthly * numDevs // total dev-team monthly cost

	sensitivityMatrix := make([]domain.SensitivityCell, 0, len(riskSteps)*len(profitSteps))
	for _, rPct := range riskSteps {
		for _, pPct := range profitSteps {
			risk       := baseLaborCost * rPct
			afterRisk  := baseLaborCost + risk
			profit     := afterRisk * pPct
			afterProfit := afterRisk + profit
			vat        := afterProfit * vatRate
			grandTotal := afterProfit + vat
			sensitivityMatrix = append(sensitivityMatrix, domain.SensitivityCell{
				RiskPct:    rPct * 100,
				ProfitPct:  pPct * 100,
				GrandTotal: round2(grandTotal),
			})
		}
	}
	data.SensitivityMatrix = sensitivityMatrix
	if billableDays > 0 {
		data.BreakEvenRate = round2(fullyLoadedMonthly / billableDays)
	}

	return data, nil
}

// buildQuotationHTML produces a self-contained A4 HTML document for the PDF export.
func buildQuotationHTML(r *domain.QuotationResponse, req *domain.QuotationRequest) string {
	now := time.Now().Format("02 January 2006")

	var taskRows strings.Builder
	for i, t := range r.Tasks {
		taskRows.WriteString(fmt.Sprintf(`
		<tr class="%s">
			<td class="num">%d</td>
			<td>%s</td>
			<td>%s</td>
			<td class="right">%.2f</td>
			<td class="right">%s</td>
		</tr>`,
			rowClass(i),
			i+1,
			escapeHTML(t.EpicTitle),
			escapeHTML(t.Title),
			t.Mandays,
			formatTHB(t.Cost),
		))
	}

	totalBeforeVAT := r.Subtotal + r.RiskAmount + r.ProfitAmount

	html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8"/>
<link rel="preconnect" href="https://fonts.googleapis.com" />
<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin />
<link rel="stylesheet" href="https://fonts.googleapis.com/css2?family=Sarabun:ital,wght@0,300;0,400;0,500;0,600;0,700;1,400&display=block" />
<style>
  * { margin: 0; padding: 0; box-sizing: border-box; }
  @page { size: A4; margin: 18mm 15mm; }
  body { font-family: 'Sarabun', sans-serif; font-size: 10pt; color: #1a2035; background: #fff; }

  /* ── Header ── */
  .header { display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 28px; padding-bottom: 20px; border-bottom: 3px solid #1e3a5f; }
  .brand-block { display: flex; flex-direction: column; gap: 4px; }
  .company-name { font-size: 20pt; font-weight: 800; color: #1e3a5f; letter-spacing: 1px; text-transform: uppercase; }
  .company-tagline { font-size: 8.5pt; color: #4a6fa5; font-weight: 500; letter-spacing: 0.5px; text-transform: uppercase; }
  .doc-meta { text-align: right; font-size: 9pt; color: #4b5563; line-height: 1.7; }
  .doc-meta .doc-title { font-size: 12pt; font-weight: 700; color: #1e3a5f; margin-bottom: 6px; text-transform: uppercase; letter-spacing: 0.5px; }
  .doc-meta strong { color: #1a2035; }

  /* ── Section headings ── */
  h2 { font-size: 10pt; font-weight: 700; color: #1e3a5f; margin-bottom: 10px; margin-top: 22px;
       border-bottom: 1.5px solid #c7d8ed; padding-bottom: 5px; text-transform: uppercase; letter-spacing: 0.8px; }

  /* ── Params bar ── */
  .params { display: flex; gap: 0; margin-bottom: 22px; border: 1px solid #c7d8ed; border-radius: 6px; overflow: hidden; }
  .param-item { flex: 1; padding: 10px 14px; border-right: 1px solid #c7d8ed; background: #f0f5fb; }
  .param-item:last-child { border-right: none; }
  .param-item label { display: block; font-size: 7.5pt; font-weight: 700; color: #4a6fa5; text-transform: uppercase; letter-spacing: 0.6px; margin-bottom: 3px; }
  .param-item span { font-size: 11pt; font-weight: 700; color: #1e3a5f; }

  /* ── Table ── */
  table { width: 100%%; border-collapse: collapse; margin-bottom: 22px; font-size: 9pt; }
  thead tr { background: #1e3a5f; }
  th { color: #fff; padding: 8px 10px; text-align: left; font-weight: 600; font-size: 8.5pt; text-transform: uppercase; letter-spacing: 0.4px; }
  th.right, td.right { text-align: right; }
  th.num, td.num { text-align: center; width: 28px; }
  td { padding: 6px 10px; border-bottom: 1px solid #dde8f4; vertical-align: top; color: #1a2035; }
  tr.alt td { background: #f5f9ff; }
  tr:last-child td { border-bottom: none; }
  td.cost { font-weight: 600; color: #1e3a5f; }

  /* ── Summary box ── */
  .summary-wrap { display: flex; justify-content: flex-end; margin-bottom: 24px; }
  .summary-box { width: 360px; border: 1px solid #c7d8ed; border-radius: 6px; overflow: hidden; }
  .summary-row { display: flex; justify-content: space-between; padding: 8px 16px; font-size: 9.5pt; border-bottom: 1px solid #dde8f4; }
  .summary-row:last-child { border-bottom: none; }
  .summary-row:nth-child(odd) { background: #f5f9ff; }
  .summary-row.subtotal-before-vat { background: #e8f0fa; border-top: 1.5px solid #4a6fa5; border-bottom: 1.5px solid #4a6fa5; }
  .summary-row.subtotal-before-vat .summary-label { font-weight: 700; color: #1e3a5f; }
  .summary-row.subtotal-before-vat .summary-amount { font-weight: 700; color: #1e3a5f; }
  .summary-row.grand-total { background: #1e3a5f; border-bottom: none; }
  .summary-row.grand-total .summary-label { color: #fff; font-weight: 700; font-size: 11pt; }
  .summary-row.grand-total .summary-amount { color: #fff; font-weight: 800; font-size: 11pt; }
  .summary-label { color: #374151; }
  .summary-amount { font-weight: 600; font-variant-numeric: tabular-nums; }

  /* ── Footer ── */
  .footer { margin-top: 28px; font-size: 7.5pt; color: #9ca3af; text-align: center; border-top: 1px solid #dde8f4; padding-top: 10px; }
  .footer strong { color: #4a6fa5; }
</style>
</head>
<body>

<div class="header">
  <div class="brand-block">
    <div class="company-name">Komgrip Technologies</div>
    <div class="company-tagline">Software Engineering &amp; Digital Solutions</div>
  </div>
  <div class="doc-meta">
    <div class="doc-title">Project Quotation</div>
    <div><strong>Date:</strong> %s</div>
    <div><strong>Project ID:</strong> %s</div>
    <div><strong>Currency:</strong> %s</div>
  </div>
</div>

<h2>Cost Model Parameters</h2>
<div class="params">
  <div class="param-item"><label>Cost / Manday</label><span>%s</span></div>
  <div class="param-item"><label>Total Mandays</label><span>%.2f days</span></div>
  <div class="param-item"><label>Risk Buffer</label><span>%.0f%%</span></div>
  <div class="param-item"><label>Profit Margin</label><span>%.0f%%</span></div>
</div>

<h2>Task Breakdown</h2>
<table>
  <thead>
    <tr>
      <th class="num">#</th>
      <th>Epic</th>
      <th>Task</th>
      <th class="right">Mandays</th>
      <th class="right">Cost (THB)</th>
    </tr>
  </thead>
  <tbody>
    %s
  </tbody>
</table>

<div class="summary-wrap">
  <div class="summary-box">
    <div class="summary-row"><span class="summary-label">Labor Subtotal</span><span class="summary-amount">%s</span></div>
    <div class="summary-row"><span class="summary-label">Risk Buffer (%.0f%%)</span><span class="summary-amount">+ %s</span></div>
    <div class="summary-row"><span class="summary-label">Profit Margin (%.0f%%)</span><span class="summary-amount">+ %s</span></div>
    <div class="summary-row subtotal-before-vat"><span class="summary-label">Total (before VAT)</span><span class="summary-amount">%s</span></div>
    <div class="summary-row"><span class="summary-label">VAT (7%%)</span><span class="summary-amount">+ %s</span></div>
    <div class="summary-row grand-total"><span class="summary-label">Grand Total</span><span class="summary-amount">%s</span></div>
  </div>
</div>

<div class="footer">
  Generated by <strong>Komgrip Technologies</strong> · Project Cost Engine · %s
</div>
</body>
</html>`,
		now,
		r.ProjectID,
		r.Currency,
		formatTHB(r.CostPerManday),
		r.TotalMandays,
		req.RiskMarginPct*100,
		req.ProfitMarginPct*100,
		taskRows.String(),
		formatTHB(r.Subtotal),
		req.RiskMarginPct*100,
		formatTHB(r.RiskAmount),
		req.ProfitMarginPct*100,
		formatTHB(r.ProfitAmount),
		formatTHB(totalBeforeVAT),
		formatTHB(r.VAT),
		formatTHB(r.GrandTotal),
		now,
	)
	return html
}

// buildCustomerQuotationHTML produces a customer-facing A4 PDF.
// It omits internal cost parameters: no mandays column, no risk/profit margin line items.
// Risk buffer and profit margin are silently absorbed into the grand total.
func buildCustomerQuotationHTML(r *domain.QuotationResponse, req *domain.QuotationRequest) string {
	now := time.Now().Format("02 January 2006")

	projectLabel := req.ProjectName
	if projectLabel == "" {
		projectLabel = r.ProjectID
	}

	var taskRows strings.Builder
	for i, t := range r.Tasks {
		taskRows.WriteString(fmt.Sprintf(`
		<tr class="%s">
			<td class="num">%d</td>
			<td>%s</td>
			<td>%s</td>
			<td class="right cost">%s</td>
		</tr>`,
			rowClass(i),
			i+1,
			escapeHTML(t.EpicTitle),
			escapeHTML(t.Title),
			formatTHB(t.Cost),
		))
	}

	totalBeforeVAT := r.Subtotal + r.RiskAmount + r.ProfitAmount

	html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8"/>
<link rel="preconnect" href="https://fonts.googleapis.com" />
<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin />
<link rel="stylesheet" href="https://fonts.googleapis.com/css2?family=Sarabun:ital,wght@0,300;0,400;0,500;0,600;0,700;1,400&display=block" />
<style>
  * { margin: 0; padding: 0; box-sizing: border-box; }
  @page { size: A4; margin: 18mm 15mm; }
  body { font-family: 'Sarabun', sans-serif; font-size: 10pt; color: #1a2035; background: #fff; }

  /* ── Header ── */
  .header { display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 28px; padding-bottom: 20px; border-bottom: 3px solid #1e3a5f; }
  .brand-block { display: flex; flex-direction: column; gap: 4px; }
  .company-name { font-size: 20pt; font-weight: 800; color: #1e3a5f; letter-spacing: 1px; text-transform: uppercase; }
  .company-tagline { font-size: 8.5pt; color: #4a6fa5; font-weight: 500; letter-spacing: 0.5px; text-transform: uppercase; }
  .doc-meta { text-align: right; font-size: 9pt; color: #4b5563; line-height: 1.7; }
  .doc-meta .doc-title { font-size: 12pt; font-weight: 700; color: #1e3a5f; margin-bottom: 6px; text-transform: uppercase; letter-spacing: 0.5px; }
  .doc-meta strong { color: #1a2035; }

  /* ── Section headings ── */
  h2 { font-size: 10pt; font-weight: 700; color: #1e3a5f; margin-bottom: 10px; margin-top: 22px;
       border-bottom: 1.5px solid #c7d8ed; padding-bottom: 5px; text-transform: uppercase; letter-spacing: 0.8px; }

  /* ── Table ── */
  table { width: 100%%; border-collapse: collapse; margin-bottom: 22px; font-size: 9pt; }
  thead tr { background: #1e3a5f; }
  th { color: #fff; padding: 8px 10px; text-align: left; font-weight: 600; font-size: 8.5pt; text-transform: uppercase; letter-spacing: 0.4px; }
  th.right, td.right { text-align: right; }
  th.num, td.num { text-align: center; width: 28px; }
  td { padding: 6px 10px; border-bottom: 1px solid #dde8f4; vertical-align: top; color: #1a2035; }
  tr.alt td { background: #f5f9ff; }
  tr:last-child td { border-bottom: none; }
  td.cost { font-weight: 600; color: #1e3a5f; }

  /* ── Summary box ── */
  .summary-wrap { display: flex; justify-content: flex-end; margin-bottom: 24px; }
  .summary-box { width: 360px; border: 1px solid #c7d8ed; border-radius: 6px; overflow: hidden; }
  .summary-row { display: flex; justify-content: space-between; padding: 8px 16px; font-size: 9.5pt; border-bottom: 1px solid #dde8f4; }
  .summary-row:last-child { border-bottom: none; }
  .summary-row:nth-child(odd) { background: #f5f9ff; }
  .summary-row.subtotal-before-vat { background: #e8f0fa; border-top: 1.5px solid #4a6fa5; border-bottom: 1.5px solid #4a6fa5; }
  .summary-row.subtotal-before-vat .summary-label { font-weight: 700; color: #1e3a5f; }
  .summary-row.subtotal-before-vat .summary-amount { font-weight: 700; color: #1e3a5f; }
  .summary-row.grand-total { background: #1e3a5f; border-bottom: none; }
  .summary-row.grand-total .summary-label { color: #fff; font-weight: 700; font-size: 11pt; }
  .summary-row.grand-total .summary-amount { color: #fff; font-weight: 800; font-size: 11pt; }
  .summary-label { color: #374151; }
  .summary-amount { font-weight: 600; font-variant-numeric: tabular-nums; }

  /* ── Validity notice ── */
  .validity-note { margin-top: 18px; padding: 10px 14px; background: #f0f5fb; border: 1px solid #c7d8ed;
    border-radius: 6px; font-size: 8.5pt; color: #4b5563; line-height: 1.6; }
  .validity-note strong { color: #1e3a5f; }

  /* ── Footer ── */
  .footer { margin-top: 28px; font-size: 7.5pt; color: #9ca3af; text-align: center; border-top: 1px solid #dde8f4; padding-top: 10px; }
  .footer strong { color: #4a6fa5; }
</style>
</head>
<body>

<div class="header">
  <div class="brand-block">
    <div class="company-name">Komgrip Technologies</div>
    <div class="company-tagline">Software Engineering &amp; Digital Solutions</div>
  </div>
  <div class="doc-meta">
    <div class="doc-title">Project Quotation</div>
    <div><strong>Date:</strong> %s</div>
    <div><strong>Project:</strong> %s</div>
    <div><strong>Currency:</strong> %s</div>
  </div>
</div>

<h2>Scope of Work</h2>
<table>
  <thead>
    <tr>
      <th class="num">#</th>
      <th>Epic</th>
      <th>Task / Deliverable</th>
      <th class="right">Amount (%s)</th>
    </tr>
  </thead>
  <tbody>
    %s
  </tbody>
</table>

<div class="summary-wrap">
  <div class="summary-box">
    <div class="summary-row subtotal-before-vat"><span class="summary-label">Total (before VAT)</span><span class="summary-amount">%s</span></div>
    <div class="summary-row"><span class="summary-label">VAT (7%%)</span><span class="summary-amount">+ %s</span></div>
    <div class="summary-row grand-total"><span class="summary-label">Grand Total</span><span class="summary-amount">%s</span></div>
  </div>
</div>

<div class="validity-note">
  <strong>Validity:</strong> This quotation is valid for 30 days from the date of issue.<br/>
  Prices are quoted in Thai Baht (THB) and inclusive of VAT at 7%%.
</div>

<div class="footer">
  Prepared by <strong>Komgrip Technologies</strong> · %s
</div>
</body>
</html>`,
		now,
		escapeHTML(projectLabel),
		r.Currency,
		r.Currency,
		taskRows.String(),
		formatTHB(totalBeforeVAT),
		formatTHB(r.VAT),
		formatTHB(r.GrandTotal),
		now,
	)
	return html
}

// ─── Cost Analysis Report HTML Builder ───────────────────────────────────────

// buildCostReportHTML produces a 5-section A4 PDF focused on monthly/annual company cost structure.
func buildCostReportHTML(d *domain.CostReportData) string {
	var sb strings.Builder

	thb := func(v float64) string { return "฿" + commaSep(v) }

	// ── SVG Helpers ───────────────────────────────────────────────────────────

	// Horizontal stacked bar: salary pool / ss / overhead
	svgStackedBar := func(salary, ss, overhead float64) string {
		total := salary + ss + overhead
		if total == 0 { return "" }
		w := 460.0
		h := 28.0
		sW := salary / total * w
		ssW := ss / total * w
		ohW := overhead / total * w
		return fmt.Sprintf(
			`<svg width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f" xmlns="http://www.w3.org/2000/svg" style="border-radius:4px;overflow:hidden">
  <rect x="0" y="0" width="%.1f" height="%.0f" fill="#1e3a5f"/>
  <rect x="%.1f" y="0" width="%.1f" height="%.0f" fill="#4a6fa5"/>
  <rect x="%.1f" y="0" width="%.1f" height="%.0f" fill="#93c5fd"/>
</svg>`, w, h, w, h,
			sW, h,
			sW, ssW, h,
			sW+ssW, ohW, h)
	}

	// Vertical bar chart for monthly 12-month projection
	svgMonthlyBars := func(rows []domain.MonthlyCashFlowRow) string {
		if len(rows) == 0 { return "" }
		maxVal := 0.0
		for _, r := range rows {
			if r.TotalOutflow > maxVal { maxVal = r.TotalOutflow }
		}
		if maxVal == 0 { maxVal = 1 }
		chartW, chartH := 520.0, 110.0
		n := len(rows)
		barW := chartW/float64(n) - 4
		var bars, labels, cumLine strings.Builder
		var prev float64
		maxCum := rows[len(rows)-1].Cumulative
		if maxCum == 0 { maxCum = 1 }
		for i, r := range rows {
			x := float64(i)*(chartW/float64(n)) + 2
			// Salary bar
			sh := r.SalaryCost / maxVal * chartH
			bars.WriteString(fmt.Sprintf(`<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" fill="#3b82f6" rx="2"/>`,
				x, chartH-sh, barW, sh))
			// Overhead bar stacked
			oh := r.ExpenseCost / maxVal * chartH
			bars.WriteString(fmt.Sprintf(`<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" fill="#f59e0b" rx="2"/>`,
				x, chartH-sh-oh, barW, oh))
			// Cumulative line
			cy := chartH - (r.Cumulative/maxCum)*chartH*0.8
			cx := x + barW/2
			if i == 0 {
				cumLine.WriteString(fmt.Sprintf("M %.1f %.1f", cx, cy))
			} else {
				cumLine.WriteString(fmt.Sprintf(" L %.1f %.1f", cx, cy))
			}
			_ = prev
			prev = cy
			// Month label every 3rd
			if i%3 == 0 || i == n-1 {
				lbl := r.Month
				if len(lbl) > 8 { lbl = lbl[:8] }
				labels.WriteString(fmt.Sprintf(`<text x="%.1f" y="%.1f" font-size="6.5" text-anchor="middle" fill="#6b7280">%s</text>`,
					cx, chartH+11, lbl))
			}
		}
		return fmt.Sprintf(
			`<svg width="%.0f" height="130" viewBox="0 0 %.0f 130" xmlns="http://www.w3.org/2000/svg">
  %s
  <path d="%s" fill="none" stroke="#1e3a5f" stroke-width="1.5" stroke-dasharray="4,2"/>
  %s
</svg>`, chartW, chartW, bars.String(), cumLine.String(), labels.String())
	}

	// Pie chart for cost composition
	svgPie := func(labels []string, values []float64) string {
		total := 0.0
		for _, v := range values { total += v }
		if total == 0 { return "" }
		colors := []string{"#1e3a5f","#3b82f6","#f59e0b","#10b981","#8b5cf6","#ef4444"}
		cx, cy, r := 55.0, 55.0, 50.0
		var paths, legend strings.Builder
		angle := -math.Pi / 2
		for i, v := range values {
			if v == 0 { continue }
			sweep := v / total * 2 * math.Pi
			x1 := cx + r*math.Cos(angle)
			y1 := cy + r*math.Sin(angle)
			x2 := cx + r*math.Cos(angle+sweep)
			y2 := cy + r*math.Sin(angle+sweep)
			la := 0; if sweep > math.Pi { la = 1 }
			col := colors[i%len(colors)]
			paths.WriteString(fmt.Sprintf(
				`<path d="M %.1f %.1f L %.1f %.1f A %.1f %.1f 0 %d 1 %.1f %.1f Z" fill="%s" stroke="white" stroke-width="1"/>`,
				cx, cy, x1, y1, r, r, la, x2, y2, col))
			lbl := ""
			if i < len(labels) { lbl = labels[i] }
			ly := float64(13 + i*14)
			legend.WriteString(fmt.Sprintf(
				`<rect x="118" y="%.1f" width="9" height="9" fill="%s" rx="1"/>
<text x="132" y="%.1f" font-size="7.5" fill="#374151">%s (%.1f%%)</text>`,
				ly-8, col, ly, escapeHTML(lbl), v/total*100))
			angle += sweep
		}
		return fmt.Sprintf(
			`<svg width="250" height="120" viewBox="0 0 250 120" xmlns="http://www.w3.org/2000/svg">%s%s</svg>`,
			paths.String(), legend.String())
	}

	// ── CSS ───────────────────────────────────────────────────────────────────
	css := `
* { margin:0; padding:0; box-sizing:border-box; }
@page { size:A4; margin:14mm 13mm; }
body { font-family:'Sarabun',sans-serif; font-size:9.5pt; color:#111827; background:#fff; }
body::before {
  content:"CONFIDENTIAL — INTERNAL USE ONLY";
  position:fixed; top:50%; left:50%;
  transform:translate(-50%,-50%) rotate(-35deg);
  font-size:26pt; font-weight:900; color:rgba(0,0,0,0.04);
  white-space:nowrap; pointer-events:none; z-index:0;
}
.hdr { display:flex; justify-content:space-between; align-items:flex-start;
  padding-bottom:12px; border-bottom:3px solid #1e3a5f; margin-bottom:14px; }
.co-name { font-size:17pt; font-weight:900; color:#1e3a5f; text-transform:uppercase; letter-spacing:1px; }
.co-tag  { font-size:8pt; color:#4a6fa5; font-weight:500; text-transform:uppercase; letter-spacing:0.5px; margin-top:3px; }
.rpt-title { font-size:12pt; font-weight:800; color:#1e3a5f; text-transform:uppercase; }
.rpt-meta { text-align:right; font-size:8.5pt; color:#4b5563; line-height:1.7; }
.conf-badge { display:inline-block; font-size:7pt; font-weight:700; color:#dc2626;
  border:1px solid #fca5a5; background:#fef2f2; padding:2px 6px; border-radius:3px;
  text-transform:uppercase; letter-spacing:0.5px; margin-top:4px; }
.sec { font-size:9.5pt; font-weight:800; color:#fff; background:#1e3a5f;
  padding:6px 12px; margin-top:16px; margin-bottom:10px; border-radius:4px;
  text-transform:uppercase; letter-spacing:0.8px; page-break-after:avoid; }
.sec-sub { font-size:8pt; font-weight:500; color:#93c5fd; float:right; text-transform:none; }
.kpi-grid { display:grid; gap:8px; margin-bottom:12px; }
.g2 { grid-template-columns:repeat(2,1fr); }
.g3 { grid-template-columns:repeat(3,1fr); }
.g4 { grid-template-columns:repeat(4,1fr); }
.g5 { grid-template-columns:repeat(5,1fr); }
.tile { border-radius:8px; padding:12px 14px; border:1px solid #e5e7eb; }
.tile.navy { background:#1e3a5f; border-color:#1e3a5f; }
.tile.blue { background:linear-gradient(135deg,#eff6ff,#dbeafe); border-color:#93c5fd; }
.tile.amber{ background:linear-gradient(135deg,#fffbeb,#fef3c7); border-color:#fcd34d; }
.tile.green{ background:linear-gradient(135deg,#f0fdf4,#dcfce7); border-color:#86efac; }
.tile.purple{background:linear-gradient(135deg,#faf5ff,#ede9fe); border-color:#c4b5fd; }
.tile.red  { background:linear-gradient(135deg,#fef2f2,#fee2e2); border-color:#fca5a5; }
.tile.gray { background:#f9fafb; border-color:#e5e7eb; }
.kpi-lbl { font-size:7pt; font-weight:700; color:#4b5563; text-transform:uppercase; letter-spacing:0.5px; margin-bottom:3px; }
.kpi-lbl.lt { color:#93c5fd; }
.kpi-lbl.th { font-size:6.5pt; color:#9ca3af; font-weight:400; text-transform:none; letter-spacing:0; }
.kpi-val { font-size:15pt; font-weight:900; color:#1e3a5f; line-height:1; }
.kpi-val.white { color:#fff; }
.kpi-val.sm { font-size:12pt; }
.kpi-val.xs { font-size:10pt; }
.kpi-sub { font-size:7.5pt; color:#6b7280; margin-top:3px; }
.kpi-sub.lt { color:#93c5fd; }
.card { border:1px solid #dde8f4; border-radius:6px; padding:12px 14px; background:#f8fafc; margin-bottom:10px; }
.card h4 { font-size:8.5pt; font-weight:700; color:#1e3a5f; margin-bottom:8px; text-transform:uppercase; letter-spacing:0.5px; }
.row  { display:flex; justify-content:space-between; align-items:center; padding:3px 0; border-bottom:1px solid #e5e7eb; font-size:8.5pt; }
.row:last-child { border-bottom:none; font-weight:700; color:#1e3a5f; font-size:9pt; }
.row .lbl { color:#4b5563; }
.row .val { font-variant-numeric:tabular-nums; font-weight:600; }
.row .indent { padding-left:12px; color:#6b7280; font-size:8pt; }
table { width:100%%; border-collapse:collapse; margin-bottom:12px; font-size:8.5pt; }
thead tr { background:#1e3a5f; }
th { color:#fff; padding:7px 9px; text-align:left; font-size:7.5pt; font-weight:700; text-transform:uppercase; letter-spacing:0.3px; }
th.r, td.r { text-align:right; }
td { padding:5px 9px; border-bottom:1px solid #e5e7eb; vertical-align:middle; }
tr.alt td { background:#f8fafc; }
tr.grp td { background:#dbeafe; color:#1e40af; font-weight:700; font-size:8pt; padding:4px 9px; }
tr.tot td { background:#1e3a5f; color:#fff; font-weight:700; }
.leg { display:flex; gap:14px; font-size:7.5pt; color:#4b5563; align-items:center; margin-bottom:6px; }
.leg-dot { width:12px; height:12px; border-radius:2px; }
.pgbrk { page-break-before:always; }
.ftr { margin-top:18px; padding-top:7px; border-top:1px solid #e5e7eb;
  font-size:7pt; color:#9ca3af; display:flex; justify-content:space-between; }
.ftr strong { color:#4a6fa5; }
`

	// ── Document Head ─────────────────────────────────────────────────────────
	sb.WriteString(fmt.Sprintf(`<!DOCTYPE html>
<html lang="th">
<head>
<meta charset="UTF-8"/>
<link rel="preconnect" href="https://fonts.googleapis.com"/>
<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin/>
<link rel="stylesheet" href="https://fonts.googleapis.com/css2?family=Sarabun:wght@300;400;500;600;700;800;900&display=block"/>
<style>%s</style>
</head>
<body>`, css))

	// ── Report Header ─────────────────────────────────────────────────────────
	sb.WriteString(fmt.Sprintf(`
<div class="hdr">
  <div>
    <div class="co-name">Komgrip Technologies</div>
    <div class="co-tag">Software Engineering &amp; Digital Solutions</div>
  </div>
  <div class="rpt-meta">
    <div class="rpt-title">Company Cost Analysis</div>
    <div><strong>Generated:</strong> %s (B.E. %d)</div>
    <div><strong>Currency:</strong> %s</div>
    <div><span class="conf-badge">Confidential — Internal Use Only</span></div>
  </div>
</div>`,
		d.GeneratedAt.Format("02 January 2006 15:04"),
		d.GeneratedAt.Year()+543,
		d.Currency,
	))

	// ══════════════════════════════════════════════════════════════════════════
	// SECTION 1 — Monthly Cost Summary
	// ══════════════════════════════════════════════════════════════════════════
	sb.WriteString(`<div class="sec">Section 1 — Monthly Cost Summary <span class="sec-sub">สรุปต้นทุนรายเดือน</span></div>`)

	sb.WriteString(fmt.Sprintf(`
<div class="kpi-grid g5">
  <div class="tile navy">
    <div class="kpi-lbl lt">Monthly Burn Rate</div>
    <div class="kpi-lbl th">อัตราเผาเงินต่อเดือน</div>
    <div class="kpi-val white sm">%s</div>
    <div class="kpi-sub lt">Total cash out / month</div>
  </div>
  <div class="tile blue">
    <div class="kpi-lbl">Gross Payroll</div>
    <div class="kpi-lbl th">เงินเดือนรวม</div>
    <div class="kpi-val sm">%s</div>
    <div class="kpi-sub">%d employees</div>
  </div>
  <div class="tile amber">
    <div class="kpi-lbl">Social Security</div>
    <div class="kpi-lbl th">ประกันสังคมรายเดือน</div>
    <div class="kpi-val sm">%s</div>
    <div class="kpi-sub">employer contribution</div>
  </div>
  <div class="tile purple">
    <div class="kpi-lbl">Company Overhead</div>
    <div class="kpi-lbl th">โสหุ้ยบริษัท</div>
    <div class="kpi-val sm">%s</div>
    <div class="kpi-sub">exec + company</div>
  </div>
  <div class="tile green">
    <div class="kpi-lbl">Annual Burn Rate</div>
    <div class="kpi-lbl th">อัตราเผาเงินต่อปี</div>
    <div class="kpi-val sm">%s</div>
    <div class="kpi-sub">× 12 months</div>
  </div>
</div>`,
		thb(d.MonthlyBurnRate),
		thb(d.MonthlyTotalPayroll), len(d.Headcount),
		thb(d.MonthlyTotalSS),
		thb(d.MonthlyTotalOverhead),
		thb(d.AnnualBurnRate),
	))

	// Monthly cost breakdown card
	sb.WriteString(fmt.Sprintf(`
<div style="display:grid;grid-template-columns:1fr 1fr;gap:12px;margin-bottom:12px">
  <div class="card">
    <h4>Monthly Cost Breakdown / รายละเอียดต้นทุนรายเดือน</h4>
    <div class="row"><span class="lbl">Gross Payroll (เงินเดือนรวม)</span><span class="val">%s</span></div>
    <div class="row"><span class="lbl indent">+ Social Security (ประกันสังคม)</span><span class="val">%s</span></div>
    <div class="row"><span class="lbl indent">+ Company Overhead (ค่าดำเนินการ)</span><span class="val">%s</span></div>
    <div class="row"><span class="lbl indent">+ Executive Expense (ผู้บริหาร)</span><span class="val">%s</span></div>
    <div class="row"><span class="lbl">Monthly Burn Rate</span><span class="val">%s</span></div>
  </div>
  <div class="card">
    <h4>Annual Cost Breakdown / รายละเอียดต้นทุนรายปี</h4>
    <div class="row"><span class="lbl">Annual Payroll (เงินเดือนรายปี)</span><span class="val">%s</span></div>
    <div class="row"><span class="lbl indent">+ Annual Social Security</span><span class="val">%s</span></div>
    <div class="row"><span class="lbl indent">+ Annual Company Overhead (×12)</span><span class="val">%s</span></div>
    <div class="row"><span class="lbl indent">+ Annual Executive Expense (×12)</span><span class="val">%s</span></div>
    <div class="row"><span class="lbl">Annual Burn Rate</span><span class="val">%s</span></div>
  </div>
</div>`,
		thb(d.MonthlyTotalPayroll),
		thb(d.MonthlyTotalSS),
		thb(d.MonthlyCompanyExpense),
		thb(d.MonthlyExecExpense),
		thb(d.MonthlyBurnRate),
		thb(d.TotalMonthlyPayroll),
		thb(d.TotalAnnualSS),
		thb(d.MonthlyCompanyExpense*12),
		thb(d.MonthlyExecExpense*12),
		thb(d.AnnualBurnRate),
	))

	// Cost composition pie chart
	pieLabels := []string{"Payroll", "Social Security", "Company OH", "Executive"}
	pieVals   := []float64{d.MonthlyTotalPayroll, d.MonthlyTotalSS, d.MonthlyCompanyExpense, d.MonthlyExecExpense}
	sb.WriteString(fmt.Sprintf(`
<div style="display:flex;gap:20px;align-items:flex-start;margin-bottom:10px">
  <div>
    <p style="font-size:8pt;font-weight:600;color:#4b5563;margin-bottom:5px">Monthly Cost Composition</p>
    %s
  </div>
  <div style="flex:1">
    <p style="font-size:8pt;font-weight:600;color:#4b5563;margin-bottom:5px">Cost Waterfall per Dev / สัดส่วนต้นทุนต่อ Developer</p>
    %s
    <div class="leg" style="margin-top:5px">
      <div class="leg-item" style="display:flex;align-items:center;gap:4px"><div class="leg-dot" style="background:#1e3a5f"></div><span>Salary %s</span></div>
      <div class="leg-item" style="display:flex;align-items:center;gap:4px"><div class="leg-dot" style="background:#4a6fa5"></div><span>SS %s</span></div>
      <div class="leg-item" style="display:flex;align-items:center;gap:4px"><div class="leg-dot" style="background:#93c5fd"></div><span>Overhead %s</span></div>
    </div>
  </div>
</div>`,
		svgPie(pieLabels, pieVals),
		svgStackedBar(d.AvgDevSalary, d.AvgDevSS, d.OverheadPerDev),
		thb(d.AvgDevSalary), thb(d.AvgDevSS), thb(d.OverheadPerDev),
	))

	// ══════════════════════════════════════════════════════════════════════════
	// SECTION 2 — Billing Rate & Cost Model
	// ══════════════════════════════════════════════════════════════════════════
	sb.WriteString(`<div class="sec pgbrk">Section 2 — Billing Rate &amp; Cost Model <span class="sec-sub">อัตราเรียกเก็บและโมเดลต้นทุน</span></div>`)

	sb.WriteString(fmt.Sprintf(`
<div class="kpi-grid g4" style="margin-bottom:12px">
  <div class="tile navy">
    <div class="kpi-lbl lt">Cost / Manday</div>
    <div class="kpi-lbl th">ต้นทุนต่อแมนเดย์</div>
    <div class="kpi-val white">%s</div>
    <div class="kpi-sub lt">Fully loaded ÷ billable days</div>
  </div>
  <div class="tile blue">
    <div class="kpi-lbl">Cost / Hour</div>
    <div class="kpi-lbl th">ต้นทุนต่อชั่วโมง</div>
    <div class="kpi-val sm">%s</div>
    <div class="kpi-sub">Manday ÷ %d hrs</div>
  </div>
  <div class="tile amber">
    <div class="kpi-lbl">Billable Days / Month</div>
    <div class="kpi-lbl th">วันที่เรียกเก็บได้ต่อเดือน</div>
    <div class="kpi-val sm">%.1f days</div>
    <div class="kpi-sub">%d work days ÷ %.2f×</div>
  </div>
  <div class="tile green">
    <div class="kpi-lbl">Utilisation Rate</div>
    <div class="kpi-lbl th">อัตราการใช้ประโยชน์</div>
    <div class="kpi-val sm">%.1f%%</div>
    <div class="kpi-sub">1 ÷ overhead multiplier</div>
  </div>
</div>`,
		thb(d.CostPerManday),
		thb(d.CostPerHour), d.Config.WorkingHoursPerDay,
		d.BillableDays, d.Config.WorkingDaysPerMonth, d.Config.OverheadMultiplier,
		d.UtilizationRate*100,
	))

	sb.WriteString(fmt.Sprintf(`
<div style="display:grid;grid-template-columns:1fr 1fr;gap:12px">
  <div class="card">
    <h4>Fully Loaded Cost per Dev / ต้นทุนเต็มรูปแบบต่อ Developer</h4>
    <div class="row"><span class="lbl">① Avg Dev Salary</span><span class="val">%s / mo</span></div>
    <div class="row"><span class="lbl">② Avg Social Security</span><span class="val">%s / mo</span></div>
    <div class="row"><span class="lbl" style="color:#6b7280;font-size:8pt">── Overhead per Dev ──</span><span class="val"></span></div>
    <div class="row"><span class="lbl indent">Product Owner salaries ÷ %d devs</span><span class="val">%s</span></div>
    <div class="row"><span class="lbl indent">Company Expense ÷ %d devs</span><span class="val">%s</span></div>
    <div class="row"><span class="lbl indent">Exec Expense ÷ %d devs</span><span class="val">%s</span></div>
    <div class="row"><span class="lbl">③ Total Overhead per Dev</span><span class="val">%s / mo</span></div>
    <div class="row"><span class="lbl">Fully Loaded Monthly (①+②+③)</span><span class="val">%s / mo</span></div>
  </div>
  <div class="card">
    <h4>Config Parameters / พารามิเตอร์การตั้งค่า</h4>
    <div class="row"><span class="lbl">Working Days / Month</span><span class="val">%d days</span></div>
    <div class="row"><span class="lbl">Working Hours / Day</span><span class="val">%d hrs</span></div>
    <div class="row"><span class="lbl">Overhead Multiplier</span><span class="val">%.2f×</span></div>
    <div class="row"><span class="lbl">Default Risk Buffer</span><span class="val">%.0f%%</span></div>
    <div class="row"><span class="lbl">Default Profit Margin</span><span class="val">%.0f%%</span></div>
    <div class="row"><span class="lbl">Monthly Company Expense</span><span class="val">%s</span></div>
    <div class="row"><span class="lbl">Monthly Executive Expense</span><span class="val">%s</span></div>
  </div>
</div>`,
		thb(d.AvgDevSalary),
		thb(d.AvgDevSS),
		d.DevCount, thb(d.TotalPMSalaryPerDev),
		d.DevCount, thb(d.CompanyExpensePerDev),
		d.DevCount, thb(d.ExecExpensePerDev),
		thb(d.OverheadPerDev),
		thb(d.FullyLoadedMonthly),
		d.Config.WorkingDaysPerMonth, d.Config.WorkingHoursPerDay,
		d.Config.OverheadMultiplier,
		d.Config.DefaultRiskBuffer*100, d.Config.DefaultProfitMargin*100,
		thb(d.Config.CompanyExpense), thb(d.Config.ExecutiveExpense),
	))

	// ══════════════════════════════════════════════════════════════════════════
	// SECTION 3 — Headcount & Salary Analysis
	// ══════════════════════════════════════════════════════════════════════════
	sb.WriteString(`<div class="sec pgbrk">Section 3 — Headcount &amp; Salary Analysis <span class="sec-sub">วิเคราะห์อัตรากำลังและเงินเดือน</span></div>`)

	sb.WriteString(fmt.Sprintf(`
<div class="kpi-grid g4" style="margin-bottom:12px">
  <div class="tile blue">
    <div class="kpi-lbl">Total Headcount</div>
    <div class="kpi-lbl th">จำนวนพนักงานทั้งหมด</div>
    <div class="kpi-val">%d</div>
    <div class="kpi-sub">%d ENG · %d PO · %d other</div>
  </div>
  <div class="tile amber">
    <div class="kpi-lbl">Monthly Payroll</div>
    <div class="kpi-lbl th">เงินเดือนรวมต่อเดือน</div>
    <div class="kpi-val sm">%s</div>
    <div class="kpi-sub">gross salaries only</div>
  </div>
  <div class="tile purple">
    <div class="kpi-lbl">Monthly SS Total</div>
    <div class="kpi-lbl th">ประกันสังคมรวมต่อเดือน</div>
    <div class="kpi-val sm">%s</div>
    <div class="kpi-sub">all employees</div>
  </div>
  <div class="tile navy">
    <div class="kpi-lbl lt">Annual Payroll + SS</div>
    <div class="kpi-lbl th">เงินเดือน+ประกันสังคมรายปี</div>
    <div class="kpi-val white sm">%s</div>
    <div class="kpi-sub lt">fully loaded × 12</div>
  </div>
</div>`,
		len(d.Headcount), d.DevCount, d.PMCount, d.OtherCount,
		thb(d.TotalMonthlyPayroll),
		thb(d.TotalMonthlySS),
		thb(d.TotalAnnualPayroll),
	))

	// Headcount table
	sb.WriteString(`<table>
<thead><tr>
  <th>Employee / พนักงาน</th><th>Role</th><th>Type</th>
  <th class="r">Monthly Salary</th><th class="r">SS / mo</th>
  <th class="r">Fully Loaded / mo</th><th class="r">Annual Cost</th>
  <th>Effective From</th>
</tr></thead><tbody>`)

	// Group by role
	grouped := make(map[string][]domain.HeadcountRow)
	order := []string{}
	seen := make(map[string]bool)
	for _, h := range d.Headcount {
		k := strings.ToUpper(h.UserRole)
		grouped[k] = append(grouped[k], h)
		if !seen[k] { order = append(order, k); seen[k] = true }
	}
	for _, role := range order {
		rows := grouped[role]
		sb.WriteString(fmt.Sprintf(`<tr class="grp"><td colspan="8">%s (%d)</td></tr>`, role, len(rows)))
		var gPay, gSS, gFL, gAnn float64
		for i, h := range rows {
			cls := ""; if i%2==1 { cls = ` class="alt"` }
			sb.WriteString(fmt.Sprintf(`<tr%s>
  <td><strong>%s</strong><br/><span style="font-size:7.5pt;color:#6b7280">%s</span></td>
  <td>%s</td><td style="font-size:8pt">%s</td>
  <td class="r" style="font-variant-numeric:tabular-nums">%s</td>
  <td class="r" style="font-variant-numeric:tabular-nums">%s</td>
  <td class="r" style="font-variant-numeric:tabular-nums"><strong>%s</strong></td>
  <td class="r" style="font-variant-numeric:tabular-nums">%s</td>
  <td style="font-size:8pt;color:#6b7280">%s</td>
</tr>`,
				cls,
				escapeHTML(h.UserDisplayName), escapeHTML(h.UserEmail),
				escapeHTML(h.UserRole), escapeHTML(h.EmploymentType),
				thb(h.MonthlySalary), thb(h.SsCost),
				thb(h.FullyLoadedMonthly), thb(h.AnnualCost),
				h.EffectiveFrom,
			))
			gPay += h.MonthlySalary; gSS += h.SsCost
			gFL += h.FullyLoadedMonthly; gAnn += h.AnnualCost
		}
		sb.WriteString(fmt.Sprintf(`<tr class="tot">
  <td colspan="3">%s Subtotal</td>
  <td class="r">%s</td><td class="r">%s</td>
  <td class="r">%s</td><td class="r">%s</td><td></td>
</tr>`, role, thb(gPay), thb(gSS), thb(gFL), thb(gAnn)))
	}
	totalFL := d.TotalMonthlyPayroll + d.TotalMonthlySS
	sb.WriteString(fmt.Sprintf(`<tr class="tot">
  <td colspan="3">GRAND TOTAL</td>
  <td class="r">%s</td><td class="r">%s</td>
  <td class="r">%s</td><td class="r">%s</td><td></td>
</tr>`,
		thb(d.TotalMonthlyPayroll), thb(d.TotalMonthlySS),
		thb(totalFL), thb(d.TotalAnnualPayroll),
	))
	sb.WriteString(`</tbody></table>`)

	// ══════════════════════════════════════════════════════════════════════════
	// SECTION 4 — 12-Month Cash Flow Projection
	// ══════════════════════════════════════════════════════════════════════════
	sb.WriteString(`<div class="sec pgbrk">Section 4 — 12-Month Cash Flow Projection <span class="sec-sub">การฉายภาพกระแสเงินสด 12 เดือน</span></div>`)

	sb.WriteString(fmt.Sprintf(`
<div class="kpi-grid g3" style="margin-bottom:10px">
  <div class="tile amber">
    <div class="kpi-lbl">Monthly Burn Rate</div>
    <div class="kpi-lbl th">ค่าใช้จ่ายต่อเดือน (คงที่)</div>
    <div class="kpi-val sm">%s</div>
    <div class="kpi-sub">Salary+SS+Overhead / month</div>
  </div>
  <div class="tile blue">
    <div class="kpi-lbl">Payroll + SS / Month</div>
    <div class="kpi-lbl th">เงินเดือน+ประกันสังคม</div>
    <div class="kpi-val sm">%s</div>
    <div class="kpi-sub">direct labor cost</div>
  </div>
  <div class="tile navy">
    <div class="kpi-lbl lt">12-Month Total</div>
    <div class="kpi-lbl th">ยอดรวม 12 เดือน</div>
    <div class="kpi-val white sm">%s</div>
    <div class="kpi-sub lt">cumulative projection</div>
  </div>
</div>`,
		thb(d.MonthlyBurnRate),
		thb(d.MonthlySalaryPool+d.MonthlySSPool),
		thb(d.AnnualProjection),
	))

	// Chart
	sb.WriteString(fmt.Sprintf(`
<div style="margin-bottom:8px">
  %s
  <div class="leg" style="margin-top:5px">
    <div style="display:flex;align-items:center;gap:4px"><div class="leg-dot" style="background:#3b82f6"></div>Payroll+SS %s/mo</div>
    <div style="display:flex;align-items:center;gap:4px"><div class="leg-dot" style="background:#f59e0b"></div>Overhead %s/mo</div>
    <div style="display:flex;align-items:center;gap:4px"><div style="width:20px;height:2px;background:#1e3a5f;border-top:2px dashed #1e3a5f"></div>Cumulative</div>
  </div>
</div>`,
		svgMonthlyBars(d.MonthlyCashFlow),
		thb(d.MonthlySalaryPool+d.MonthlySSPool),
		thb(d.MonthlyTotalOverhead),
	))

	// Table
	sb.WriteString(`<table>
<thead><tr>
  <th>Month / เดือน</th>
  <th class="r">Payroll + SS</th>
  <th class="r">Overhead (OH)</th>
  <th class="r">Total Outflow / เดือน</th>
  <th class="r">Cumulative / สะสม</th>
</tr></thead><tbody>`)
	for i, m := range d.MonthlyCashFlow {
		cls := ""; if i%2==1 { cls = ` class="alt"` }
		sb.WriteString(fmt.Sprintf(`<tr%s>
  <td><strong>%s</strong></td>
  <td class="r" style="font-variant-numeric:tabular-nums">%s</td>
  <td class="r" style="font-variant-numeric:tabular-nums">%s</td>
  <td class="r" style="font-variant-numeric:tabular-nums"><strong>%s</strong></td>
  <td class="r" style="font-variant-numeric:tabular-nums;color:#4a6fa5">%s</td>
</tr>`, cls, m.Month,
			thb(m.SalaryCost), thb(m.ExpenseCost),
			thb(m.TotalOutflow), thb(m.Cumulative)))
	}
	if len(d.MonthlyCashFlow) > 0 {
		last := d.MonthlyCashFlow[len(d.MonthlyCashFlow)-1]
		sb.WriteString(fmt.Sprintf(`<tr class="tot">
  <td>12-MONTH TOTAL</td>
  <td class="r">%s</td>
  <td class="r">%s</td>
  <td class="r">%s</td>
  <td class="r">%s</td>
</tr>`,
			thb((d.MonthlySalaryPool+d.MonthlySSPool)*12),
			thb(d.MonthlyTotalOverhead*12),
			thb(d.AnnualProjection),
			thb(last.Cumulative),
		))
	}
	sb.WriteString(`</tbody></table>`)

	// ══════════════════════════════════════════════════════════════════════════
	// SECTION 5 — Sensitivity Analysis
	// ══════════════════════════════════════════════════════════════════════════
	sb.WriteString(`<div class="sec pgbrk">Section 5 — Pricing Sensitivity Analysis <span class="sec-sub">การวิเคราะห์ความอ่อนไหวด้านราคา</span></div>`)

	sb.WriteString(fmt.Sprintf(`
<p style="font-size:8.5pt;color:#374151;margin-bottom:10px">
  Based on fully loaded dev-team monthly cost <strong>%s</strong> (%d devs × %s/dev).
  Matrix shows Grand Total (incl. 7%% VAT) at varying Risk Buffer × Profit Margin.
  Highlighted cell = default config (Risk %.0f%% / Profit %.0f%%).
</p>`,
		thb(d.FullyLoadedMonthly*float64(d.DevCount)),
		d.DevCount, thb(d.FullyLoadedMonthly),
		d.Config.DefaultRiskBuffer*100, d.Config.DefaultProfitMargin*100,
	))

	riskLevels   := []float64{5, 10, 15, 20}
	profitLevels := []float64{15, 20, 25, 30, 35}
	lookup := make(map[string]float64)
	for _, c := range d.SensitivityMatrix {
		lookup[fmt.Sprintf("%.0f_%.0f", c.RiskPct, c.ProfitPct)] = c.GrandTotal
	}

	sb.WriteString(`<table>
<thead><tr><th style="background:#374151">Risk \ Profit</th>`)
	for _, p := range profitLevels {
		sb.WriteString(fmt.Sprintf(`<th class="r" style="background:#374151">%.0f%%</th>`, p))
	}
	sb.WriteString(`</tr></thead><tbody>`)
	for _, risk := range riskLevels {
		sb.WriteString(fmt.Sprintf(`<tr><td style="background:#eff6ff;font-weight:700;color:#1e40af">Risk %.0f%%</td>`, risk))
		for _, profit := range profitLevels {
			val := lookup[fmt.Sprintf("%.0f_%.0f", risk, profit)]
			style := ""
			if risk == d.Config.DefaultRiskBuffer*100 && profit == d.Config.DefaultProfitMargin*100 {
				style = ` style="background:#dbeafe;font-weight:800;color:#1e40af"`
			}
			sb.WriteString(fmt.Sprintf(`<td class="r" style="font-variant-numeric:tabular-nums"%s>%s</td>`, style, thb(val)))
		}
		sb.WriteString(`</tr>`)
	}
	sb.WriteString(`</tbody></table>`)

	// Break-even + billing rate KPIs
	sb.WriteString(fmt.Sprintf(`
<div class="kpi-grid g3" style="margin-top:10px">
  <div class="tile amber">
    <div class="kpi-lbl">Break-even Rate / Manday</div>
    <div class="kpi-lbl th">อัตราคุ้มทุนต่อแมนเดย์</div>
    <div class="kpi-val sm">%s</div>
    <div class="kpi-sub">minimum billing to cover cost</div>
  </div>
  <div class="tile blue">
    <div class="kpi-lbl">Recommended Rate (default)</div>
    <div class="kpi-lbl th">ราคาแนะนำ (ตามค่าเริ่มต้น)</div>
    <div class="kpi-val sm">%s</div>
    <div class="kpi-sub">Risk %.0f%% + Profit %.0f%% + VAT</div>
  </div>
  <div class="tile navy">
    <div class="kpi-lbl lt">Fully Loaded / Manday</div>
    <div class="kpi-lbl th">ต้นทุนเต็มรูปแบบต่อแมนเดย์</div>
    <div class="kpi-val white sm">%s</div>
    <div class="kpi-sub lt">before risk &amp; profit</div>
  </div>
</div>`,
		thb(d.BreakEvenRate),
		thb(lookup[fmt.Sprintf("%.0f_%.0f", d.Config.DefaultRiskBuffer*100, d.Config.DefaultProfitMargin*100)]/d.BillableDays),
		d.Config.DefaultRiskBuffer*100, d.Config.DefaultProfitMargin*100,
		thb(d.CostPerManday),
	))

	// ── Footer ───────────────────────────────────────────────────────────────
	sb.WriteString(fmt.Sprintf(`
<div class="ftr">
  <div>Generated by <strong>Komgrip Technologies</strong> — Cost Analysis Engine v2.1 | %s (B.E. %d)</div>
  <div>CONFIDENTIAL — FOR INTERNAL USE ONLY</div>
</div>
</body></html>`,
		d.GeneratedAt.Format("02 Jan 2006 15:04:05"),
		d.GeneratedAt.Year()+543,
	))

	return sb.String()
}


// shortNum formats large numbers as compact strings (e.g. 150000 → "150K").
func shortNum(v float64) string {
	if v >= 1_000_000 {
		return fmt.Sprintf("%.1fM", v/1_000_000)
	} else if v >= 1000 {
		return fmt.Sprintf("%.0fK", v/1000)
	}
	return fmt.Sprintf("%.0f", v)
}

func rowClass(i int) string {
	if i%2 == 1 {
		return "alt"
	}
	return ""
}

func formatTHB(v float64) string {
	// Simple comma-formatted number with 2 decimal places.
	return fmt.Sprintf("%s", commaSep(v))
}

func commaSep(v float64) string {
	s := fmt.Sprintf("%.2f", v)
	// Insert commas every 3 digits before the decimal point.
	dotIdx := strings.Index(s, ".")
	intPart := s[:dotIdx]
	decPart := s[dotIdx:]
	var b strings.Builder
	for i, c := range intPart {
		pos := len(intPart) - i
		if i > 0 && pos%3 == 0 {
			b.WriteRune(',')
		}
		b.WriteRune(c)
	}
	b.WriteString(decPart)
	return b.String()
}

func escapeHTML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	return s
}

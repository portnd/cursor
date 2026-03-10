import { useAuth } from '~/composables/useAuth'

export interface QuotationRequest {
  dev_user_ids: number[]
  risk_margin_pct: number
  profit_margin_pct: number
  task_ids?: string[]
  epic_ids?: string[]
}

export interface TaskCostLine {
  task_id: string
  title: string
  epic_title: string
  est_days: number
  mandays: number
  cost_per_manday: number
  cost: number
}

export interface QuotationResponse {
  project_id: string
  tasks: TaskCostLine[]
  subtotal: number
  risk_amount: number
  profit_amount: number
  vat: number
  grand_total: number
  cost_per_manday: number
  billable_days: number
  total_mandays: number
  currency: string
}

export interface CompanyCostConfig {
  id: number
  working_days_per_month: number
  working_hours_per_day: number
  overhead_multiplier: number
  default_profit_margin: number
  default_risk_buffer: number
  currency: string
  executive_expense: number
  company_expense: number
  created_at: string
  updated_at: string
}

export interface CompanyMandayRateResponse {
  total_monthly_salaries: number
  total_monthly_ss: number
  company_expense: number
  executive_expense: number
  overhead_role_salary_total: number
  company_expense_total: number
  total_monthly_burn_rate: number
  active_headcount: number
  working_days_per_month: number
  overhead_multiplier: number
  billable_days: number
  cost_per_manday: number
  cost_per_hour: number
  currency: string
}

export interface EmployeeSalaryWithUser {
  id: number
  user_id: number
  monthly_salary: number
  currency: string
  effective_from: string
  effective_to: string | null
  employment_type: 'FULLTIME' | 'PARTTIME' | 'CONTRACTOR'
  cost_per_minute: number
  created_at: string
  updated_at: string
  user_email: string
  user_display_name: string
  user_role: string
  /** Min(monthly_salary * 5%, 875 THB) — ค่าประกันสังคม */
  ss_cost: number
}

export interface UpsertSalaryPayload {
  user_id: number
  monthly_salary: number
  currency?: string
  effective_from: string
  effective_to?: string
  employment_type?: 'FULLTIME' | 'PARTTIME' | 'CONTRACTOR'
}

export interface UpdateCostConfigPayload {
  working_days_per_month: number
  working_hours_per_day: number
  overhead_multiplier: number
  default_profit_margin: number
  default_risk_buffer: number
  currency?: string
  executive_expense?: number
  company_expense?: number
}

export type CostReportPeriod = 'current_month' | 'current_quarter' | 'ytd' | 'all' | 'custom'

export interface CostReportRequest {
  period: CostReportPeriod
  date_from?: string   // YYYY-MM-DD, only used when period = 'custom'
  date_to?: string     // YYYY-MM-DD, only used when period = 'custom'
  project_ids?: string[]
}

function usePricingApi() {
  const { fetchWithAuth, token, apiBase } = useAuth()

  // ── Quotation ──────────────────────────────────────────────────────────────

  async function calculateQuotation(projectId: string, payload: QuotationRequest): Promise<QuotationResponse> {
    return await fetchWithAuth<QuotationResponse>(`/sentinel/projects/${projectId}/quotation/calculate`, {
      method: 'POST',
      body: payload,
      timeoutMs: 60000,
    })
  }

  async function exportQuotationPDF(projectId: string, payload: QuotationRequest): Promise<string> {
    const fullUrl = `${apiBase.value}/sentinel/projects/${projectId}/quotation/export`
    const res = await fetch(fullUrl, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token.value}`,
      },
      body: JSON.stringify(payload),
      signal: AbortSignal.timeout(120_000),
    })
    if (!res.ok) {
      const errText = await res.text()
      throw new Error(`PDF export failed (${res.status}): ${errText}`)
    }
    const blob = await res.blob()
    return URL.createObjectURL(blob)
  }

  // ── Admin: Company Cost Config ─────────────────────────────────────────────

  async function getCostConfig(): Promise<CompanyCostConfig> {
    return await fetchWithAuth<CompanyCostConfig>('/pricing/config')
  }

  async function getCompanyMandayRate(): Promise<CompanyMandayRateResponse> {
    return await fetchWithAuth<CompanyMandayRateResponse>('/pricing/manday-rate')
  }

  async function updateCostConfig(payload: UpdateCostConfigPayload): Promise<CompanyCostConfig> {
    return await fetchWithAuth<CompanyCostConfig>('/pricing/config', {
      method: 'PUT',
      body: payload,
    })
  }

  // ── Admin: Employee Salaries ───────────────────────────────────────────────

  async function listSalaries(): Promise<EmployeeSalaryWithUser[]> {
    const res = await fetchWithAuth<{ data: EmployeeSalaryWithUser[] }>('/pricing/salaries')
    return res?.data ?? []
  }

  async function upsertSalary(payload: UpsertSalaryPayload): Promise<EmployeeSalaryWithUser> {
    return await fetchWithAuth<EmployeeSalaryWithUser>('/pricing/salaries', {
      method: 'POST',
      body: payload,
    })
  }

  async function deleteSalary(id: number): Promise<void> {
    await fetchWithAuth(`/pricing/salaries/${id}`, { method: 'DELETE' })
  }

  async function exportCostReport(payload: CostReportRequest): Promise<string> {
    const fullUrl = `${apiBase.value}/pricing/report/export`
    const res = await fetch(fullUrl, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token.value}`,
      },
      body: JSON.stringify(payload),
      signal: AbortSignal.timeout(180_000),
    })
    if (!res.ok) {
      const errText = await res.text()
      throw new Error(`Cost report export failed (${res.status}): ${errText}`)
    }
    const blob = await res.blob()
    return URL.createObjectURL(blob)
  }

  return {
    calculateQuotation,
    exportQuotationPDF,
    getCostConfig,
    getCompanyMandayRate,
    updateCostConfig,
    listSalaries,
    upsertSalary,
    deleteSalary,
    exportCostReport,
  }
}

export { usePricingApi }

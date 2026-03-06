import { useAuth } from '~/composables/useAuth'

export interface MonthlyEntry {
  id: number
  year: number
  month: number
  revenue: number
  expenses: number
  cash_balance: number
  note: string
  created_at: string
  updated_at: string
}

export interface FinanceSummary {
  cash_balance: number
  runway_months: number
  burn_rate: number
  last_month_mrr: number
  net_new_arr: number
  currency: string
  last_entry_year: number
  last_entry_month: number
}

export interface CreateOrUpdateEntryPayload {
  year: number
  month: number
  revenue: number
  expenses: number
  cash_balance: number
  note?: string
}

function useFinanceApi() {
  const { fetchWithAuth } = useAuth()

  async function getEntries(limit = 24): Promise<MonthlyEntry[]> {
    const res = await fetchWithAuth<{ data: MonthlyEntry[] }>(`/finance/entries?limit=${limit}`)
    return res?.data ?? []
  }

  async function createOrUpdateEntry(payload: CreateOrUpdateEntryPayload): Promise<MonthlyEntry> {
    return await fetchWithAuth<MonthlyEntry>('/finance/entries', {
      method: 'POST',
      body: payload,
    })
  }

  async function getSummary(): Promise<FinanceSummary | null> {
    const res = await fetchWithAuth<FinanceSummary>('/finance/summary').catch(() => null)
    return res ?? null
  }

  return {
    getEntries,
    createOrUpdateEntry,
    getSummary,
  }
}

export { useFinanceApi }

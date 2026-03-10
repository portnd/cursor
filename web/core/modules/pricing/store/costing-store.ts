import { defineStore } from 'pinia'
import { usePricingApi, type QuotationRequest, type QuotationResponse, type CostReportRequest } from '../infrastructure/pricing-api'

interface CostingState {
  result: QuotationResponse | null
  loading: boolean
  exporting: boolean
  error: string | null
}

export const useCostingStore = defineStore('costing', {
  state: (): CostingState => ({
    result: null,
    loading: false,
    exporting: false,
    error: null,
  }),

  getters: {
    hasResult: (state): boolean => state.result !== null,
    grandTotal: (state): number => state.result?.grand_total ?? 0,
    currency: (state): string => state.result?.currency ?? 'THB',
  },

  actions: {
    async calculateQuotation(projectId: string, req: QuotationRequest) {
      const api = usePricingApi()
      this.loading = true
      this.error = null
      try {
        this.result = await api.calculateQuotation(projectId, req)
      } catch (e: any) {
        this.error = e?.data?.error ?? e?.message ?? 'Calculation failed'
        this.result = null
      } finally {
        this.loading = false
      }
    },

    async exportPDF(projectId: string, req: QuotationRequest) {
      const { token, apiBase } = useAuth()
      this.exporting = true
      this.error = null
      try {
        if (!token.value) throw new Error('Not authenticated')

        const url = `${apiBase.value}/sentinel/projects/${projectId}/quotation/export`
        const response = await fetch(url, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            Authorization: `Bearer ${token.value}`,
          },
          body: JSON.stringify(req),
          signal: AbortSignal.timeout(120_000),
        })

        if (!response.ok) {
          let msg = `PDF generation failed (${response.status})`
          try { const j = await response.json(); msg = j.error || j.message || msg } catch {}
          throw new Error(msg)
        }

        const blob = await response.blob()
        const objectUrl = URL.createObjectURL(blob)

        // Open in new tab — browser displays PDF natively (same as timeline export)
        const link = document.createElement('a')
        link.href = objectUrl
        link.target = '_blank'
        link.rel = 'noopener noreferrer'
        document.body.appendChild(link)
        link.click()
        document.body.removeChild(link)
        setTimeout(() => URL.revokeObjectURL(objectUrl), 60_000)
      } catch (e: any) {
        this.error = e?.message ?? 'PDF export failed'
      } finally {
        this.exporting = false
      }
    },

    reset() {
      this.result = null
      this.error = null
    },

    async exportCostReport(req: CostReportRequest) {
      const { token, apiBase } = useAuth()
      this.exporting = true
      this.error = null
      try {
        if (!token.value) throw new Error('Not authenticated')

        const url = `${apiBase.value}/pricing/report/export`
        const response = await fetch(url, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            Authorization: `Bearer ${token.value}`,
          },
          body: JSON.stringify(req),
          signal: AbortSignal.timeout(180_000),
        })

        if (!response.ok) {
          let msg = `Cost report export failed (${response.status})`
          try { const j = await response.json(); msg = j.error || j.message || msg } catch {}
          throw new Error(msg)
        }

        const blob = await response.blob()
        const objectUrl = URL.createObjectURL(blob)
        const link = document.createElement('a')
        link.href = objectUrl
        link.download = `cost-analysis-report-${new Date().toISOString().slice(0, 10)}.pdf`
        link.target = '_blank'
        link.rel = 'noopener noreferrer'
        document.body.appendChild(link)
        link.click()
        document.body.removeChild(link)
        setTimeout(() => URL.revokeObjectURL(objectUrl), 60_000)
      } catch (e: any) {
        this.error = e?.message ?? 'Cost report export failed'
      } finally {
        this.exporting = false
      }
    },
  },
})

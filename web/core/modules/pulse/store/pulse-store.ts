import { defineStore } from 'pinia'
import {
  usePulseApi,
  type CompanyPulseResponse,
  type DailyStandup,
  type SubmitStandupRequest,
} from '../infrastructure/pulse-api'

interface PulseState {
  pulse: CompanyPulseResponse | null
  todayStandup: DailyStandup | null
  loading: boolean
  submitting: boolean
  error: string | null
  /** ISO date string (YYYY-MM-DD) last successfully fetched */
  lastFetchedDate: string | null
}

export const usePulseStore = defineStore('pulse', {
  state: (): PulseState => ({
    pulse: null,
    todayStandup: null,
    loading: false,
    submitting: false,
    error: null,
    lastFetchedDate: null,
  }),

  getters: {
    /** True if the current user has already checked in today */
    hasCheckedInToday: (state): boolean => {
      if (!state.pulse) return false
      const auth = useAuth()
      const uid = auth.currentUser.value?.user_id
      if (!uid) return false
      return state.pulse.members.some(
        (m) => m.user_id === uid && m.standup !== null,
      )
    },

    membersWithBlockers: (state) =>
      state.pulse?.members.filter((m) => m.has_blocker) ?? [],

    checkinRate: (state): number => {
      if (!state.pulse || state.pulse.total_members === 0) return 0
      return Math.round((state.pulse.checked_in / state.pulse.total_members) * 100)
    },
  },

  actions: {
    async fetchDailyPulse(date?: string) {
      const api = usePulseApi()
      this.loading = true
      this.error = null
      try {
        this.pulse = await api.getDailyPulse(date)
        this.lastFetchedDate = date ?? new Date().toISOString().slice(0, 10)
      } catch (e: any) {
        this.error = e?.data?.error ?? e?.message ?? 'Failed to load pulse'
      } finally {
        this.loading = false
      }
    },

    async submitStandup(payload: SubmitStandupRequest) {
      const api = usePulseApi()
      this.submitting = true
      this.error = null
      try {
        this.todayStandup = await api.submitStandup(payload)
        // Refresh the pulse board after successful checkin
        await this.fetchDailyPulse(payload.date)
        return true
      } catch (e: any) {
        this.error = e?.data?.error ?? e?.message ?? 'Failed to submit standup'
        return false
      } finally {
        this.submitting = false
      }
    },

    reset() {
      this.pulse = null
      this.todayStandup = null
      this.error = null
      this.lastFetchedDate = null
    },
  },
})

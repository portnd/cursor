import { defineStore } from 'pinia'
import {
  useAttendanceApi,
  type AttendanceRecord,
  type OffsiteCheckInRequest,
  type OffsiteCheckOutRequest,
  type OfficeConfig,
  type UpsertOfficeConfigPayload,
} from '../infrastructure/attendance-api'

interface AttendanceState {
  todayRecord: AttendanceRecord | null
  officeConfig: OfficeConfig | null
  todayOffsiteRequest: OffsiteCheckInRequest | null
  todayOffsiteCheckOutRequest: OffsiteCheckOutRequest | null
  isRemote: boolean
  history: AttendanceRecord[]
  historyCursor: string | null
  historyHasMore: boolean
  loading: boolean
  actionLoading: boolean
  error: string | null
}

export const useAttendanceStore = defineStore('attendance', {
  state: (): AttendanceState => ({
    todayRecord: null,
    officeConfig: null,
    todayOffsiteRequest: null,
    todayOffsiteCheckOutRequest: null,
    isRemote: false,
    history: [],
    historyCursor: null,
    historyHasMore: false,
    loading: false,
    actionLoading: false,
    error: null,
  }),

  getters: {
    canCheckIn(state): boolean {
      return !state.todayRecord?.check_in_at
    },
    canCheckOut(state): boolean {
      return !!state.todayRecord?.check_in_at && !state.todayRecord?.check_out_at
    },
  },

  actions: {
    async fetchToday() {
      const api = useAttendanceApi()
      this.loading = true
      this.error = null
      try {
        const res = await api.getToday()
        this.todayRecord = res.record
        this.officeConfig = res.office_config
        this.todayOffsiteRequest = res.offsite_checkin_request
        this.todayOffsiteCheckOutRequest = res.offsite_checkout_request
        this.isRemote = res.is_remote ?? false
      } catch (e: any) {
        this.error = e?.data?.error ?? e?.message ?? 'Failed to load attendance'
      } finally {
        this.loading = false
      }
    },

    async fetchHistory(append = false) {
      const api = useAttendanceApi()
      this.loading = !append
      this.error = null
      try {
        const cursor = append ? this.historyCursor ?? undefined : undefined
        const res = await api.getHistory(cursor, 20)
        if (append) {
          this.history = [...this.history, ...res.items]
        } else {
          this.history = res.items
        }
        this.historyCursor = res.next_cursor ?? null
        this.historyHasMore = !!res.next_cursor
      } catch (e: any) {
        this.error = e?.data?.error ?? e?.message ?? 'Failed to load history'
      } finally {
        this.loading = false
      }
    },

    async loadMoreHistory() {
      if (!this.historyHasMore || !this.historyCursor) return
      await this.fetchHistory(true)
    },

    async checkIn(lat: number, lng: number) {
      const api = useAttendanceApi()
      this.actionLoading = true
      this.error = null
      try {
        this.todayRecord = await api.checkIn(lat, lng)
        await this.fetchHistory(false)
        return true
      } catch (e: any) {
        this.error = e?.data?.error ?? e?.message ?? 'Check-in failed'
        return false
      } finally {
        this.actionLoading = false
      }
    },

    async checkOut() {
      const api = useAttendanceApi()
      this.actionLoading = true
      this.error = null
      try {
        this.todayRecord = await api.checkOut()
        await this.fetchHistory(false)
        return true
      } catch (e: any) {
        this.error = e?.data?.error ?? e?.message ?? 'Check-out failed'
        return false
      } finally {
        this.actionLoading = false
      }
    },

    async requestOffsiteCheckIn(lat: number, lng: number, reason: string) {
      const api = useAttendanceApi()
      this.actionLoading = true
      this.error = null
      try {
        this.todayOffsiteRequest = await api.requestOffsiteCheckIn({ lat, lng, reason })
        await this.fetchToday()
        return true
      } catch (e: any) {
        this.error = e?.data?.error ?? e?.message ?? 'Offsite check-in request failed'
        return false
      } finally {
        this.actionLoading = false
      }
    },

    async requestOffsiteCheckOut(lat: number, lng: number, reason: string) {
      const api = useAttendanceApi()
      this.actionLoading = true
      this.error = null
      try {
        this.todayOffsiteCheckOutRequest = await api.requestOffsiteCheckOut({ lat, lng, reason })
        await this.fetchToday()
        return true
      } catch (e: any) {
        this.error = e?.data?.error ?? e?.message ?? 'Offsite check-out request failed'
        return false
      } finally {
        this.actionLoading = false
      }
    },

    async adminSaveConfig(payload: UpsertOfficeConfigPayload) {
      const api = useAttendanceApi()
      this.actionLoading = true
      this.error = null
      try {
        this.officeConfig = await api.adminPutConfig(payload)
        return true
      } catch (e: any) {
        this.error = e?.data?.error ?? e?.message ?? 'Failed to save config'
        return false
      } finally {
        this.actionLoading = false
      }
    },

    async adminLoadConfig() {
      const api = useAttendanceApi()
      this.loading = true
      this.error = null
      try {
        this.officeConfig = await api.adminGetConfig()
      } catch (e: any) {
        this.error = e?.data?.error ?? e?.message ?? 'Failed to load config'
      } finally {
        this.loading = false
      }
    },
  },
})

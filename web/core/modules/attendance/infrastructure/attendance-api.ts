import { useAuth } from '~/composables/useAuth'

export interface OfficeConfig {
  id: number
  name: string
  latitude: number
  longitude: number
  radius_meters: number
  allowed_ips: string[]
  work_start_time: string
  work_end_time: string
  work_days: number[]
  wfh_days?: number[]
  is_active: boolean
  created_at?: string
  updated_at?: string
}

export interface AttendanceRecord {
  id: number
  user_id: number
  office_config_id: number
  attendance_date: string
  check_in_at?: string | null
  check_out_at?: string | null
  check_in_lat?: number | null
  check_in_lng?: number | null
  check_in_method: string
  check_in_ip: string
  is_late: boolean
  early_checkout: boolean
  status: string
  user_email?: string
  user_display_name?: string
}

export interface TodayResponse {
  record: AttendanceRecord | null
  office_config: OfficeConfig | null
}

export interface AttendanceHistoryResponse {
  items: AttendanceRecord[]
  next_cursor?: string
}

export interface UpsertOfficeConfigPayload {
  name: string
  latitude: number
  longitude: number
  radius_meters: number
  allowed_ips: string[]
  work_start_time: string
  work_end_time: string
  work_days: number[]
  wfh_days: number[]
  is_active: boolean
}

export interface AdminRecordsResponse {
  date: string
  records: AttendanceRecord[]
}

function useAttendanceApi() {
  const { fetchWithAuth } = useAuth()

  async function checkIn(lat: number, lng: number): Promise<AttendanceRecord> {
    return await fetchWithAuth<AttendanceRecord>('/attendance/check-in', {
      method: 'POST',
      body: { lat, lng },
    })
  }

  async function checkOut(): Promise<AttendanceRecord> {
    return await fetchWithAuth<AttendanceRecord>('/attendance/check-out', {
      method: 'POST',
    })
  }

  async function getToday(): Promise<TodayResponse> {
    return await fetchWithAuth<TodayResponse>('/attendance/today')
  }

  async function getHistory(cursor?: string, limit?: number): Promise<AttendanceHistoryResponse> {
    const q = new URLSearchParams()
    if (cursor) q.set('cursor', cursor)
    if (limit != null) q.set('limit', String(limit))
    const suffix = q.toString() ? `?${q.toString()}` : ''
    return await fetchWithAuth<AttendanceHistoryResponse>(`/attendance/history${suffix}`)
  }

  async function adminGetConfig(): Promise<OfficeConfig | null> {
    return await fetchWithAuth<OfficeConfig | null>('/attendance/admin/config')
  }

  async function adminPutConfig(payload: UpsertOfficeConfigPayload): Promise<OfficeConfig> {
    return await fetchWithAuth<OfficeConfig>('/attendance/admin/config', {
      method: 'PUT',
      body: payload,
    })
  }

  async function adminRecords(date?: string): Promise<AdminRecordsResponse> {
    const q = date ? `?date=${encodeURIComponent(date)}` : ''
    return await fetchWithAuth<AdminRecordsResponse>(`/attendance/admin/records${q}`)
  }

  return {
    checkIn,
    checkOut,
    getToday,
    getHistory,
    adminGetConfig,
    adminPutConfig,
    adminRecords,
  }
}

export { useAttendanceApi }

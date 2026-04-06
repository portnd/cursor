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

export interface LeaveRequest {
  id: number
  user_id: number
  start_date: string
  end_date: string
  days_requested: number
  leave_type: 'ANNUAL' | 'SICK' | 'PERSONAL' | 'UNPAID'
  reason: string
  status: 'PENDING' | 'APPROVED' | 'REJECTED'
  approver_id?: number | null
  manager_comment?: string
  approved_at?: string | null
  created_at: string
  updated_at: string
  user_email?: string
  user_display_name?: string
  approver_email?: string
  approver_name?: string
}

export interface LeaveListResponse {
  items: LeaveRequest[]
}

export interface CreateLeaveRequestPayload {
  start_date: string
  end_date: string
  leave_type: 'ANNUAL' | 'SICK' | 'PERSONAL' | 'UNPAID'
  reason: string
}

export interface ReviewLeaveRequestPayload {
  status: 'APPROVED' | 'REJECTED'
  comment?: string
}

export interface LeaveBalanceSummary {
  leave_type: string
  annual_quota_days: number
  carry_forward_days: number
  approved_days_taken: number
  remaining_days: number
}

export interface LeaveBalanceResponse {
  year: number
  items: LeaveBalanceSummary[]
}

export interface LeavePolicy {
  id: number
  leave_type: 'ANNUAL' | 'SICK' | 'PERSONAL' | 'UNPAID'
  annual_quota_days: number
  max_carry_forward_days: number
  is_active: boolean
}

export interface LeavePoliciesResponse {
  items: LeavePolicy[]
}

export interface HolidayCalendar {
  id: number
  date: string
  name: string
}

export interface HolidayListResponse {
  items: HolidayCalendar[]
}

export interface LeaveAuditLog {
  id: number
  leave_id: number
  action: string
  actor_id?: number
  actor_role: string
  old_status: string
  new_status: string
  comment: string
  metadata: string
  created_at: string
  actor_email?: string
  actor_name?: string
}

export interface LeaveAuditResponse {
  items: LeaveAuditLog[]
}

export interface LeaveNotification {
  id: number
  user_id: number
  leave_id: number
  channel: 'IN_APP' | 'EMAIL' | 'LINE'
  event: string
  title: string
  message: string
  is_read: boolean
  delivered_at?: string
  created_at: string
}

export interface LeaveNotificationResponse {
  items: LeaveNotification[]
}

export interface LeaveTrendPoint {
  month: string
  team_id?: number
  team_name?: string
  requested: number
  approved: number
  rejected: number
  total_days: number
}

export interface LeaveTrendResponse {
  items: LeaveTrendPoint[]
}

export interface LeavePolicyUpsertPayload {
  leave_type: 'ANNUAL' | 'SICK' | 'PERSONAL' | 'UNPAID'
  annual_quota_days: number
  max_carry_forward_days: number
  is_active: boolean
}

export interface HolidayUpsertPayload {
  date: string
  name: string
}

export interface LeaveBackfillItem {
  employee_email: string
  start_date: string
  end_date: string
  leave_type: 'ANNUAL' | 'SICK' | 'PERSONAL' | 'UNPAID'
  status: 'PENDING' | 'APPROVED' | 'REJECTED'
  reason: string
  comment?: string
}

export interface LeaveBackfillRequest {
  item: LeaveBackfillItem
}

export interface LeaveBackfillBulkRequest {
  items: LeaveBackfillItem[]
}

export interface LeaveBackfillBulkResultItem {
  index: number
  email: string
  status: string
  leave_id?: number
  error?: string
}

export interface LeaveBackfillBulkResponse {
  total: number
  succeeded: number
  failed: number
  results: LeaveBackfillBulkResultItem[]
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

  async function createLeaveRequest(payload: CreateLeaveRequestPayload): Promise<LeaveRequest> {
    return await fetchWithAuth<LeaveRequest>('/attendance/leaves', {
      method: 'POST',
      body: payload,
    })
  }

  async function getMyLeaveRequests(): Promise<LeaveListResponse> {
    return await fetchWithAuth<LeaveListResponse>('/attendance/leaves/my')
  }

  async function getPendingLeaveRequests(): Promise<LeaveListResponse> {
    return await fetchWithAuth<LeaveListResponse>('/attendance/admin/leaves/pending')
  }

  async function reviewLeaveRequest(id: number, payload: ReviewLeaveRequestPayload): Promise<LeaveRequest> {
    return await fetchWithAuth<LeaveRequest>(`/attendance/admin/leaves/${id}/review`, {
      method: 'PATCH',
      body: payload,
    })
  }

  async function getMyLeaveBalance(year?: number): Promise<LeaveBalanceResponse> {
    const q = year ? `?year=${year}` : ''
    return await fetchWithAuth<LeaveBalanceResponse>(`/attendance/leaves/balance${q}`)
  }

  async function getLeavePolicies(): Promise<LeavePoliciesResponse> {
    return await fetchWithAuth<LeavePoliciesResponse>('/attendance/admin/leaves/policies')
  }

  async function upsertLeavePolicy(payload: LeavePolicyUpsertPayload): Promise<LeavePolicy> {
    return await fetchWithAuth<LeavePolicy>('/attendance/admin/leaves/policies', {
      method: 'PUT',
      body: payload,
    })
  }

  async function getHolidays(from?: string, to?: string): Promise<HolidayListResponse> {
    const q = new URLSearchParams()
    if (from) q.set('from', from)
    if (to) q.set('to', to)
    const suffix = q.toString() ? `?${q.toString()}` : ''
    return await fetchWithAuth<HolidayListResponse>(`/attendance/holidays${suffix}`)
  }

  async function upsertHoliday(payload: HolidayUpsertPayload): Promise<HolidayCalendar> {
    return await fetchWithAuth<HolidayCalendar>('/attendance/admin/holidays', {
      method: 'PUT',
      body: payload,
    })
  }

  async function getLeaveAudit(leaveId: number): Promise<LeaveAuditResponse> {
    return await fetchWithAuth<LeaveAuditResponse>(`/attendance/admin/leaves/${leaveId}/audit`)
  }

  async function getLeaveNotifications(unreadOnly = false): Promise<LeaveNotificationResponse> {
    return await fetchWithAuth<LeaveNotificationResponse>(`/attendance/leaves/notifications?unread_only=${unreadOnly ? 'true' : 'false'}`)
  }

  async function markLeaveNotificationRead(id: number): Promise<{ ok: boolean }> {
    return await fetchWithAuth<{ ok: boolean }>(`/attendance/leaves/notifications/${id}/read`, {
      method: 'PATCH',
    })
  }

  async function getLeaveTrend(from?: string, to?: string): Promise<LeaveTrendResponse> {
    const q = new URLSearchParams()
    if (from) q.set('from', from)
    if (to) q.set('to', to)
    const suffix = q.toString() ? `?${q.toString()}` : ''
    return await fetchWithAuth<LeaveTrendResponse>(`/attendance/admin/leaves/trend${suffix}`)
  }

  async function backfillLeave(payload: LeaveBackfillRequest): Promise<LeaveRequest> {
    return await fetchWithAuth<LeaveRequest>('/attendance/admin/leaves/backfill', {
      method: 'POST',
      body: payload,
    })
  }

  async function backfillLeaveBulk(payload: LeaveBackfillBulkRequest): Promise<LeaveBackfillBulkResponse> {
    return await fetchWithAuth<LeaveBackfillBulkResponse>('/attendance/admin/leaves/backfill/bulk', {
      method: 'POST',
      body: payload,
    })
  }

  return {
    checkIn,
    checkOut,
    getToday,
    getHistory,
    adminGetConfig,
    adminPutConfig,
    adminRecords,
    createLeaveRequest,
    getMyLeaveRequests,
    getPendingLeaveRequests,
    reviewLeaveRequest,
    getMyLeaveBalance,
    getLeavePolicies,
    upsertLeavePolicy,
    getHolidays,
    upsertHoliday,
    getLeaveAudit,
    getLeaveNotifications,
    markLeaveNotificationRead,
    getLeaveTrend,
    backfillLeave,
    backfillLeaveBulk,
  }
}

export { useAttendanceApi }

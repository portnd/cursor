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
  offsite_checkin_request: OffsiteCheckInRequest | null
  offsite_checkout_request: OffsiteCheckOutRequest | null
  is_remote?: boolean
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

export interface OffsiteCheckInRequest {
  id: number
  user_id: number
  office_config_id: number
  attendance_date: string
  request_lat: number
  request_lng: number
  reason: string
  status: 'PENDING' | 'APPROVED' | 'REJECTED'
  approver_id?: number | null
  approver_note?: string
  requested_at: string
  approved_at?: string | null
  user_email?: string
  user_display_name?: string
  approver_email?: string
  approver_name?: string
}

export interface OffsiteCheckInListResponse {
  items: OffsiteCheckInRequest[]
}

export interface OffsiteCheckOutRequest {
  id: number
  user_id: number
  office_config_id: number
  attendance_date: string
  request_lat: number
  request_lng: number
  reason: string
  status: 'PENDING' | 'APPROVED' | 'REJECTED'
  approver_id?: number | null
  approver_note?: string
  requested_at: string
  approved_at?: string | null
  user_email?: string
  user_display_name?: string
  approver_email?: string
  approver_name?: string
}

export interface OffsiteCheckOutListResponse {
  items: OffsiteCheckOutRequest[]
}

export interface RequestOffsiteCheckInPayload {
  lat: number
  lng: number
  reason: string
}

export interface ReviewOffsiteCheckInPayload {
  status: 'APPROVED' | 'REJECTED'
  note?: string
}

export interface RequestOffsiteCheckOutPayload {
  lat: number
  lng: number
  reason: string
}

export interface ReviewOffsiteCheckOutPayload {
  status: 'APPROVED' | 'REJECTED'
  note?: string
}

export interface LeaveRequest {
  id: number
  user_id: number
  start_date: string
  end_date: string
  days_requested: number
  leave_type: 'ANNUAL' | 'SICK' | 'PERSONAL' | 'UNPAID'
  is_half_day?: boolean
  half_day_session?: 'AM' | 'PM' | ''
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
  is_half_day?: boolean
  half_day_session?: 'AM' | 'PM'
  reason: string
}

export interface ReviewLeaveRequestPayload {
  status: 'APPROVED' | 'REJECTED'
  comment?: string
}

export interface UpdateLeaveRequestPayload {
  start_date: string
  end_date: string
  leave_type: 'ANNUAL' | 'SICK' | 'PERSONAL' | 'UNPAID'
  is_half_day?: boolean
  half_day_session?: 'AM' | 'PM'
  reason: string
}

export interface CancelLeaveRequestPayload {
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
  user_id?: number
  user_name?: string
  user_email?: string
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
  is_half_day?: boolean
  half_day_session?: 'AM' | 'PM'
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

  async function requestOffsiteCheckIn(payload: RequestOffsiteCheckInPayload): Promise<OffsiteCheckInRequest> {
    return await fetchWithAuth<OffsiteCheckInRequest>('/attendance/offsite-check-in/request', {
      method: 'POST',
      body: payload,
    })
  }

  async function requestOffsiteCheckOut(payload: RequestOffsiteCheckOutPayload): Promise<OffsiteCheckOutRequest> {
    return await fetchWithAuth<OffsiteCheckOutRequest>('/attendance/offsite-check-out/request', {
      method: 'POST',
      body: payload,
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

  async function adminDeleteRecord(id: number): Promise<{ ok: boolean }> {
    return await fetchWithAuth<{ ok: boolean }>(`/attendance/admin/records/${id}`, {
      method: 'DELETE',
    })
  }

  async function adminListPendingOffsiteCheckIn(): Promise<OffsiteCheckInListResponse> {
    return await fetchWithAuth<OffsiteCheckInListResponse>('/attendance/admin/offsite-check-in/pending')
  }

  async function adminReviewOffsiteCheckIn(id: number, payload: ReviewOffsiteCheckInPayload): Promise<OffsiteCheckInRequest> {
    return await fetchWithAuth<OffsiteCheckInRequest>(`/attendance/admin/offsite-check-in/${id}/review`, {
      method: 'PATCH',
      body: payload,
    })
  }

  async function adminListPendingOffsiteCheckOut(): Promise<OffsiteCheckOutListResponse> {
    return await fetchWithAuth<OffsiteCheckOutListResponse>('/attendance/admin/offsite-check-out/pending')
  }

  async function adminReviewOffsiteCheckOut(id: number, payload: ReviewOffsiteCheckOutPayload): Promise<OffsiteCheckOutRequest> {
    return await fetchWithAuth<OffsiteCheckOutRequest>(`/attendance/admin/offsite-check-out/${id}/review`, {
      method: 'PATCH',
      body: payload,
    })
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

  async function getAdminLeaveRequests(): Promise<LeaveListResponse> {
    return await fetchWithAuth<LeaveListResponse>('/attendance/admin/leaves')
  }

  async function reviewLeaveRequest(id: number, payload: ReviewLeaveRequestPayload): Promise<LeaveRequest> {
    return await fetchWithAuth<LeaveRequest>(`/attendance/admin/leaves/${id}/review`, {
      method: 'PATCH',
      body: payload,
    })
  }

  async function updateAdminLeaveRequest(id: number, payload: UpdateLeaveRequestPayload): Promise<LeaveRequest> {
    return await fetchWithAuth<LeaveRequest>(`/attendance/admin/leaves/${id}`, {
      method: 'PATCH',
      body: payload,
    })
  }

  async function cancelAdminLeaveRequest(id: number, payload: CancelLeaveRequestPayload = {}): Promise<LeaveRequest> {
    return await fetchWithAuth<LeaveRequest>(`/attendance/admin/leaves/${id}/cancel`, {
      method: 'PATCH',
      body: payload,
    })
  }

  async function deleteAdminLeaveRequest(id: number): Promise<{ ok: boolean }> {
    return await fetchWithAuth<{ ok: boolean }>(`/attendance/admin/leaves/${id}`, {
      method: 'DELETE',
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
    requestOffsiteCheckIn,
    requestOffsiteCheckOut,
    getToday,
    getHistory,
    adminGetConfig,
    adminPutConfig,
    adminRecords,
    adminDeleteRecord,
    adminListPendingOffsiteCheckIn,
    adminReviewOffsiteCheckIn,
    adminListPendingOffsiteCheckOut,
    adminReviewOffsiteCheckOut,
    createLeaveRequest,
    getMyLeaveRequests,
    getPendingLeaveRequests,
    getAdminLeaveRequests,
    reviewLeaveRequest,
    updateAdminLeaveRequest,
    cancelAdminLeaveRequest,
    deleteAdminLeaveRequest,
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

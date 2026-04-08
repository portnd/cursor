import { useAuth } from '~/composables/useAuth'

export interface PersonalKPIs {
  user_id: number
  email: string
  role: string
  health_score: number
  delivery_rate_pct: number
  code_quality_index: number
  rework_rate_pct: number
  time_accuracy_pct: number
  sprint_velocity_sp: number
  velocity_trend: string
}

export interface TeamMemberKPI {
  user_id: number
  email: string
  role: string
  health_score: number
  delivery_rate_pct: number
  code_quality_index: number
  rework_rate_pct: number
  time_accuracy_pct: number
  sprint_velocity_sp: number
  composite_score: number
}

export interface TeamKPIsResponse {
  members: TeamMemberKPI[]
}

// ─── Discipline Dashboard types ────────────────────────────────────────────────

export interface DisciplineUserDayStat {
  date: string
  tasks_closed: number
  reworks: number
  logged_minutes: number
  has_daily_pulse: boolean
  deployments_completed: number // deployment_requests marked DEPLOYED by this reviewer
  // Attendance
  is_late: boolean
  early_checkout: boolean
  check_in_at?: string  // HH:MM ICT — empty if no record
  check_out_at?: string // HH:MM ICT — empty if no record
  attendance_status?: string // present | late | absent | wfh | ""
}

export interface DisciplineUser {
  user_id: number
  user_email: string
  user_display_name: string
  user_avatar_url?: string
  role: string
  missed_pulse_count: number
  total_tasks_closed: number
  total_reworks: number
  total_logged_hours: number
  total_deployments: number          // total deployments completed as reviewer (Chief Engineer)
  total_late_days: number            // times checked in late
  total_early_checkout_days: number  // times left before work_end
  days: DisciplineUserDayStat[]
}

export interface DisciplineAttendanceRecord {
  check_in_at?: string  // HH:MM ICT
  check_out_at?: string // HH:MM ICT
  is_late: boolean
  early_checkout: boolean
  status: string // present | late | absent | wfh
  check_in_method?: string
}

export interface DisciplineResponse {
  from_date: string
  to_date: string
  dates: string[]
  users: DisciplineUser[]
}

export interface DisciplineTimeLogEntry {
  task_id: string
  task_code: string
  task_title: string
  minutes: number
  hours: number
  description: string
}

export interface DisciplineCompletedTask {
  task_id: string
  task_code: string
  task_title: string
  story_points: number
  task_type: string
}

export interface DisciplineReworkEntry {
  task_id: string
  task_code: string
  task_title: string
  rejected_comment: string
}

export interface DisciplineDeployedRequest {
  id: number
  title: string
  branch: string
  environment: string
}

export interface DisciplineDayDetail {
  user_id: number
  user_email: string
  user_display_name: string
  date: string
  has_daily_pulse: boolean
  total_logged_minutes: number
  attendance?: DisciplineAttendanceRecord
  time_logs: DisciplineTimeLogEntry[]
  completed_tasks: DisciplineCompletedTask[]
  reworks: DisciplineReworkEntry[]
  deployed_requests: DisciplineDeployedRequest[]
}

/** Which metric the user opened the score breakdown from (team leaderboard). */
export type PerformanceBreakdownFocus = 'composite' | 'delivery' | 'quality' | 'rework'

export interface OverviewKPIs {
  engineering_health_index: number
  sprint_success_rate_pct: number
  project_on_track_rate_pct: number
  milestone_hit_rate_pct: number
  cursor_adoption_score: number
  team_velocity_trend_pct: number
}

function usePerformanceApi() {
  const { fetchWithAuth } = useAuth()

  async function getPersonalKPIs(): Promise<PersonalKPIs> {
    return await fetchWithAuth<PersonalKPIs>('/performance/me')
  }

  async function getTeamKPIs(): Promise<TeamKPIsResponse> {
    return await fetchWithAuth<TeamKPIsResponse>('/performance/team')
  }

  async function getOverviewKPIs(): Promise<OverviewKPIs | null> {
    const res = await fetchWithAuth<OverviewKPIs>('/performance/overview').catch(() => null)
    return res
  }

  async function resetReworkRate(userId: number): Promise<void> {
    await fetchWithAuth(`/performance/users/${userId}/reset-rework`, { method: 'POST' })
  }

  async function getDiscipline(from: string, to: string): Promise<DisciplineResponse> {
    return await fetchWithAuth<DisciplineResponse>(`/performance/discipline?from=${from}&to=${to}`)
  }

  async function getDisciplineDayDetail(userId: number, date: string): Promise<DisciplineDayDetail> {
    return await fetchWithAuth<DisciplineDayDetail>(`/performance/discipline/detail?user_id=${userId}&date=${date}`)
  }

  return {
    getPersonalKPIs,
    getTeamKPIs,
    getOverviewKPIs,
    resetReworkRate,
    getDiscipline,
    getDisciplineDayDetail,
  }
}

export { usePerformanceApi }

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

  return {
    getPersonalKPIs,
    getTeamKPIs,
    getOverviewKPIs,
    resetReworkRate,
  }
}

export { usePerformanceApi }

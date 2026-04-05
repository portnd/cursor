import { useAuth } from '~/composables/useAuth'
import type { Project } from '~/core/modules/projects/infrastructure/projects-api'

export interface Team {
  id: number
  name: string
  capital_balance: number
  bonus_percentage: number
  created_at: string
  updated_at: string
  users?: TeamUser[]
}

export interface TeamUser {
  id: number
  email: string
  role: 'CEO' | 'MANAGER' | 'PRODUCT_OWNER' | 'PM' | 'ENGINEER' | 'CHIEF_ENGINEER' | 'SUPPORT'
  display_name: string
  team_id: number | null
  health_score: number
}

export interface TeamMonthlyCost {
  team_id: number
  member_cost: number
  shared_overhead: number
  total_monthly_cost: number
  capital_balance: number
  bonus_percentage: number
  runway_months: number
}

export interface TeamTransaction {
  id: number
  team_id: number
  type: 'INJECTION' | 'BURN' | 'BONUS_PAYOUT'
  amount: number
  reference: string
  created_at: string
}

export interface InjectCapitalPayload {
  amount: number
  bonus_percentage: number
  note: string
}

export interface EditCapitalPayload {
  new_balance: number
  bonus_percentage?: number
  note: string
}

export interface CloseCycleResult {
  team_id: number
  balance_before: number
  bonus_percentage: number
  bonus_amount: number
  balance_after: number
}

export function useTeamsApi() {
  const { fetchWithAuth } = useAuth()

  const getTeams = (): Promise<Team[]> =>
    fetchWithAuth<{ data: Team[] }>('/auth/teams').then(r => r.data ?? [])

  const createTeam = (name: string): Promise<Team> =>
    fetchWithAuth<{ data: Team }>('/auth/teams', {
      method: 'POST',
      body: JSON.stringify({ name }),
    }).then(r => r.data)

  const updateTeam = (id: number, name: string): Promise<Team> =>
    fetchWithAuth<{ data: Team }>(`/auth/teams/${id}`, {
      method: 'PATCH',
      body: JSON.stringify({ name }),
    }).then(r => r.data)

  const deleteTeam = (id: number): Promise<void> =>
    fetchWithAuth<void>(`/auth/teams/${id}`, { method: 'DELETE' })

  const assignUserToTeam = (userId: number, teamId: number | null): Promise<void> =>
    fetchWithAuth<void>(`/auth/users/${userId}/assign-team`, {
      method: 'PATCH',
      body: JSON.stringify({ team_id: teamId }),
    })

  const assignProjectToTeam = (projectId: string, teamId: number | null): Promise<void> =>
    fetchWithAuth<void>(`/sentinel/projects/${projectId}/assign-team`, {
      method: 'PATCH',
      body: JSON.stringify({ team_id: teamId }),
    })

  /** When teams feature is disabled: set Product Owner user IDs as project owners (CEO/MANAGER). Body uses pm_user_ids for API compatibility. */
  const assignProjectPmOwners = (projectId: string, pmUserIds: number[]): Promise<Project> =>
    fetchWithAuth<{ data: Project }>(`/sentinel/projects/${projectId}/pm-owners`, {
      method: 'PATCH',
      body: JSON.stringify({ pm_user_ids: pmUserIds }),
    }).then((r) => r.data)

  // --- Team Finance / Internal VC Model ---

  const getTeamMonthlyCost = (teamId: number): Promise<TeamMonthlyCost> =>
    fetchWithAuth<{ data: TeamMonthlyCost }>(`/auth/teams/${teamId}/finance/cost`).then(r => r.data)

  const injectCapital = (teamId: number, payload: InjectCapitalPayload): Promise<Team> =>
    fetchWithAuth<{ data: Team }>(`/auth/teams/${teamId}/finance/inject`, {
      method: 'POST',
      body: JSON.stringify(payload),
    }).then(r => r.data)

  const editCapital = (teamId: number, payload: EditCapitalPayload): Promise<Team> =>
    fetchWithAuth<{ data: Team }>(`/auth/teams/${teamId}/finance/capital`, {
      method: 'PUT',
      body: JSON.stringify(payload),
    }).then(r => r.data)

  const closeCycle = (teamId: number): Promise<CloseCycleResult> =>
    fetchWithAuth<{ data: CloseCycleResult }>(`/auth/teams/${teamId}/finance/close-cycle`, {
      method: 'POST',
    }).then(r => r.data)

  const getTeamsFeatureEnabled = (): Promise<boolean> =>
    fetchWithAuth<{ data: { enabled: boolean } }>('/auth/settings/teams-feature').then(r => r.data?.enabled ?? true)

  const setTeamsFeatureEnabled = (enabled: boolean): Promise<void> =>
    fetchWithAuth<void>('/auth/settings/teams-feature', {
      method: 'PUT',
      body: JSON.stringify({ enabled }),
    })

  return {
    getTeams,
    createTeam,
    updateTeam,
    deleteTeam,
    assignUserToTeam,
    assignProjectToTeam,
    assignProjectPmOwners,
    getTeamMonthlyCost,
    injectCapital,
    editCapital,
    closeCycle,
    getTeamsFeatureEnabled,
    setTeamsFeatureEnabled,
  }
}

import { defineStore } from 'pinia'
import {
  useTeamsApi,
  type Team,
  type TeamMonthlyCost,
  type InjectCapitalPayload,
  type EditCapitalPayload,
  type CloseCycleResult,
} from '../infrastructure/teams-api'

export const useTeamsStore = defineStore('teams', {
  state: () => ({
    teams: [] as Team[],
    teamCosts: {} as Record<number, TeamMonthlyCost>,
    loading: false,
    costLoading: {} as Record<number, boolean>,
    error: null as string | null,
  }),

  getters: {
    teamById: (state) => (id: number) => state.teams.find(t => t.id === id),
    teamNameById: (state) => (id: number | null | undefined): string => {
      if (!id) return '—'
      return state.teams.find(t => t.id === id)?.name ?? '—'
    },
    teamCostById: (state) => (id: number): TeamMonthlyCost | null =>
      state.teamCosts[id] ?? null,
  },

  actions: {
    async fetchTeams() {
      this.loading = true
      this.error = null
      try {
        const api = useTeamsApi()
        this.teams = await api.getTeams()
      } catch (e: unknown) {
        this.error = e instanceof Error ? e.message : 'Failed to fetch teams'
      } finally {
        this.loading = false
      }
    },

    async createTeam(name: string): Promise<Team | null> {
      try {
        const api = useTeamsApi()
        const team = await api.createTeam(name)
        this.teams.unshift(team)
        return team
      } catch (e: unknown) {
        this.error = e instanceof Error ? e.message : 'Failed to create team'
        return null
      }
    },

    async updateTeam(id: number, name: string): Promise<Team | null> {
      try {
        const api = useTeamsApi()
        const updated = await api.updateTeam(id, name)
        const idx = this.teams.findIndex(t => t.id === id)
        if (idx !== -1) this.teams[idx] = { ...this.teams[idx], ...updated }
        return updated
      } catch (e: unknown) {
        this.error = e instanceof Error ? e.message : 'Failed to update team'
        return null
      }
    },

    async deleteTeam(id: number) {
      try {
        const api = useTeamsApi()
        await api.deleteTeam(id)
        this.teams = this.teams.filter(t => t.id !== id)
        delete this.teamCosts[id]
      } catch (e: unknown) {
        this.error = e instanceof Error ? e.message : 'Failed to delete team'
      }
    },

    async assignUserToTeam(userId: number, teamId: number | null) {
      const api = useTeamsApi()
      await api.assignUserToTeam(userId, teamId)
      for (const team of this.teams) {
        if (team.users) {
          const user = team.users.find(u => u.id === userId)
          if (user) user.team_id = teamId
        }
      }
    },

    async assignProjectToTeam(projectId: string, teamId: number | null) {
      const api = useTeamsApi()
      await api.assignProjectToTeam(projectId, teamId)
    },

    // --- Team Finance ---

    async fetchTeamCost(teamId: number): Promise<TeamMonthlyCost | null> {
      this.costLoading[teamId] = true
      try {
        const api = useTeamsApi()
        const cost = await api.getTeamMonthlyCost(teamId)
        this.teamCosts[teamId] = cost
        // Sync capital/bonus back onto the team object
        const team = this.teams.find(t => t.id === teamId)
        if (team) {
          team.capital_balance = cost.capital_balance
          team.bonus_percentage = cost.bonus_percentage
        }
        return cost
      } catch (e: unknown) {
        this.error = e instanceof Error ? e.message : 'Failed to fetch team cost'
        return null
      } finally {
        this.costLoading[teamId] = false
      }
    },

    async injectCapital(teamId: number, payload: InjectCapitalPayload): Promise<Team | null> {
      try {
        const api = useTeamsApi()
        const updated = await api.injectCapital(teamId, payload)
        const idx = this.teams.findIndex(t => t.id === teamId)
        if (idx !== -1) this.teams[idx] = { ...this.teams[idx], ...updated }
        // Refresh cost after injection
        await this.fetchTeamCost(teamId)
        return updated
      } catch (e: unknown) {
        this.error = e instanceof Error ? e.message : 'Failed to inject capital'
        return null
      }
    },

    async editCapital(teamId: number, payload: EditCapitalPayload): Promise<Team | null> {
      try {
        const api = useTeamsApi()
        const updated = await api.editCapital(teamId, payload)
        const idx = this.teams.findIndex(t => t.id === teamId)
        if (idx !== -1) this.teams[idx] = { ...this.teams[idx], ...updated }
        await this.fetchTeamCost(teamId)
        return updated
      } catch (e: unknown) {
        this.error = e instanceof Error ? e.message : 'Failed to edit capital'
        return null
      }
    },

    async closeCycle(teamId: number): Promise<CloseCycleResult | null> {
      try {
        const api = useTeamsApi()
        const result = await api.closeCycle(teamId)
        // Update team balance in store
        const team = this.teams.find(t => t.id === teamId)
        if (team) team.capital_balance = result.balance_after
        // Refresh cost after payout
        await this.fetchTeamCost(teamId)
        return result
      } catch (e: unknown) {
        this.error = e instanceof Error ? e.message : 'Failed to close cycle'
        return null
      }
    },
  },
})

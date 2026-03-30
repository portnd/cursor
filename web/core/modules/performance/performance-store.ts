import { defineStore } from 'pinia'
import { usePerformanceApi } from './performance-api'
import type { PersonalKPIs, TeamMemberKPI, OverviewKPIs } from './performance-api'

export const usePerformanceStore = defineStore('performance', {
  state: () => ({
    personal: null as PersonalKPIs | null,
    team: [] as TeamMemberKPI[],
    overview: null as OverviewKPIs | null,
    loading: false,
    error: null as string | null,
  }),

  getters: {
    hasPersonal: (state) => !!state.personal,
    hasTeam: (state) => state.team.length > 0,
    hasOverview: (state) => !!state.overview,
  },

  actions: {
    async fetchPersonal() {
      const api = usePerformanceApi()
      this.loading = true
      this.error = null
      try {
        this.personal = await api.getPersonalKPIs()
      } catch (e: any) {
        this.error = e?.message || 'Failed to load personal KPIs'
      } finally {
        this.loading = false
      }
    },

    async fetchTeam() {
      const api = usePerformanceApi()
      this.loading = true
      this.error = null
      try {
        const res = await api.getTeamKPIs()
        this.team = res.members || []
      } catch (e: any) {
        this.error = e?.message || 'Failed to load team KPIs'
      } finally {
        this.loading = false
      }
    },

    async fetchOverview() {
      const api = usePerformanceApi()
      this.loading = true
      this.error = null
      try {
        this.overview = await api.getOverviewKPIs() ?? null
      } catch (e: any) {
        this.error = e?.message || 'Failed to load overview'
      } finally {
        this.loading = false
      }
    },

    /** Loads role-appropriate KPIs in one loading cycle (parallel where possible). */
    async fetchAll(role: string) {
      const api = usePerformanceApi()
      const r = (role || 'DEV').toUpperCase()
      this.loading = true
      this.error = null
      try {
        if (r === 'CEO') {
          const [personal, teamRes, overview] = await Promise.all([
            api.getPersonalKPIs(),
            api.getTeamKPIs(),
            api.getOverviewKPIs().catch(() => null),
          ])
          this.personal = personal
          this.team = teamRes.members || []
          this.overview = overview
          return
        }
        if (r === 'PM') {
          const [personal, teamRes] = await Promise.all([
            api.getPersonalKPIs(),
            api.getTeamKPIs(),
          ])
          this.personal = personal
          this.team = teamRes.members || []
          this.overview = null
          return
        }
        this.personal = await api.getPersonalKPIs()
        this.team = []
        this.overview = null
      } catch (e: any) {
        this.error = e?.message || 'Failed to load performance'
        if (r === 'DEV') {
          this.personal = null
        }
        if (r === 'PM') {
          this.personal = null
          this.team = []
        }
        if (r === 'CEO') {
          this.team = []
          this.overview = null
        }
      } finally {
        this.loading = false
      }
    },
  },
})

import { defineStore } from 'pinia'
import { useProjectsApi } from '../infrastructure/projects-api'
import type { Project, Sprint, Milestone, ProjectAnalytics } from '../infrastructure/projects-api'

export const useProjectsStore = defineStore('projects', {
  state: () => ({
    projects: [] as Project[],
    currentProject: null as Project | null,
    sprints: [] as Sprint[],
    milestones: [] as Milestone[],
    analytics: null as ProjectAnalytics | null,
    loading: false,
    error: null as string | null,
  }),

  getters: {
    activeSprint: (state) => state.sprints.find((s) => s.status === 'ACTIVE') ?? null,
    upcomingMilestones: (state) =>
      state.milestones
        .filter((m) => m.status === 'PENDING' && m.due_date)
        .sort((a, b) => new Date(a.due_date!).getTime() - new Date(b.due_date!).getTime())
        .slice(0, 5),
  },

  actions: {
    async fetchProjects() {
      const api = useProjectsApi()
      this.loading = true
      this.error = null
      try {
        this.projects = await api.getProjects()
      } catch (e: any) {
        this.error = e.message
      } finally {
        this.loading = false
      }
    },

    async fetchProject(idOrCode: string) {
      const api = useProjectsApi()
      this.loading = true
      this.error = null
      try {
        this.currentProject = await api.getProject(idOrCode)
      } catch (e: any) {
        this.error = e.message
      } finally {
        this.loading = false
      }
    },

    async createProject(payload: { name: string; description?: string; status?: string }) {
      const api = useProjectsApi()
      const project = await api.createProject(payload)
      this.projects.unshift(project)
      return project
    },

    async deleteProject(id: string) {
      const api = useProjectsApi()
      await api.deleteProject(id)
      this.projects = this.projects.filter((p) => p.id !== id)
    },

    async fetchSprints(projectId: string) {
      const api = useProjectsApi()
      this.sprints = await api.getSprints(projectId)
    },

    async createSprint(payload: Parameters<ReturnType<typeof useProjectsApi>['createSprint']>[0]) {
      const api = useProjectsApi()
      const sprint = await api.createSprint(payload)
      this.sprints.unshift(sprint)
      return sprint
    },

    async updateSprint(id: string, payload: Partial<Sprint>) {
      const api = useProjectsApi()
      const sprint = await api.updateSprint(id, payload)
      const idx = this.sprints.findIndex((s) => s.id === id)
      if (idx !== -1) this.sprints[idx] = sprint
      return sprint
    },

    async startSprint(id: string) {
      const api = useProjectsApi()
      const sprint = await api.startSprint(id)
      const idx = this.sprints.findIndex((s) => s.id === id)
      if (idx !== -1) this.sprints[idx] = sprint
      return sprint
    },

    async completeSprint(id: string) {
      const api = useProjectsApi()
      const sprint = await api.completeSprint(id)
      const idx = this.sprints.findIndex((s) => s.id === id)
      if (idx !== -1) this.sprints[idx] = sprint
      return sprint
    },

    async deleteSprint(id: string) {
      const api = useProjectsApi()
      await api.deleteSprint(id)
      this.sprints = this.sprints.filter((s) => s.id !== id)
    },

    async addTasksToSprint(sprintId: string, taskIds: string[]) {
      const api = useProjectsApi()
      await api.addTasksToSprint(sprintId, taskIds)
    },

    async fetchMilestones(projectId: string) {
      const api = useProjectsApi()
      this.milestones = await api.getMilestones(projectId)
    },

    async createMilestone(payload: Parameters<ReturnType<typeof useProjectsApi>['createMilestone']>[0]) {
      const api = useProjectsApi()
      const milestone = await api.createMilestone(payload)
      this.milestones.push(milestone)
      return milestone
    },

    async updateMilestone(id: string, payload: Partial<Milestone>) {
      const api = useProjectsApi()
      const milestone = await api.updateMilestone(id, payload)
      const idx = this.milestones.findIndex((m) => m.id === id)
      if (idx !== -1) this.milestones[idx] = milestone
      return milestone
    },

    async deleteMilestone(id: string) {
      const api = useProjectsApi()
      await api.deleteMilestone(id)
      this.milestones = this.milestones.filter((m) => m.id !== id)
    },

    async fetchAnalytics(projectId: string) {
      const api = useProjectsApi()
      this.analytics = await api.getProjectAnalytics(projectId)
    },
  },
})

import { defineStore } from 'pinia'
import { useTasksApi } from '../infrastructure/tasks-api'
import type { Task } from '../../projects/infrastructure/projects-api'
import type { TaskComment, TimeLog } from '../infrastructure/tasks-api'

const TASK_STATUSES = ['PENDING', 'IN_PROGRESS', 'READY_FOR_TEST', 'REVIEW_PENDING', 'COMPLETED', 'BLOCKED'] as const
type TaskStatus = typeof TASK_STATUSES[number]

export const useTasksStore = defineStore('tasks', {
  state: () => ({
    allTasks: [] as Task[],
    currentTask: null as Task | null,
    comments: [] as TaskComment[],
    timeLogs: [] as TimeLog[],
    loading: false,
    error: null as string | null,
  }),

  getters: {
    tasksByStatus: (state): Record<TaskStatus, Task[]> => {
      const map = {} as Record<TaskStatus, Task[]>
      for (const s of TASK_STATUSES) map[s] = []
      for (const t of state.allTasks) {
        const s = t.status as TaskStatus
        if (map[s]) map[s].push(t)
      }
      return map
    },

    backlogTasks: (state) =>
      state.allTasks.filter((t) => !t.sprint_id && !t.parent_id),

    epicTasks: (state) =>
      state.allTasks.filter((t) => !t.parent_id),

    totalLoggedMinutes: (state) =>
      state.timeLogs.reduce((sum, l) => sum + l.minutes, 0),
  },

  actions: {
    async fetchProjectTasks(projectId: string) {
      const api = useTasksApi()
      this.loading = true
      this.error = null
      try {
        this.allTasks = await api.getAllTasks().then((tasks) =>
          tasks.filter((t) => t.project_id === projectId)
        )
      } catch (e: any) {
        this.error = e.message
      } finally {
        this.loading = false
      }
    },

    async fetchAllTasks() {
      const api = useTasksApi()
      this.loading = true
      this.error = null
      try {
        this.allTasks = await api.getAllTasks()
      } catch (e: any) {
        this.error = e.message
      } finally {
        this.loading = false
      }
    },

    async createTask(payload: Parameters<ReturnType<typeof useTasksApi>['createTask']>[0]) {
      const api = useTasksApi()
      const task = await api.createTask(payload)
      this.allTasks.unshift(task)
      return task
    },

    async updateTask(id: string, payload: Parameters<ReturnType<typeof useTasksApi>['updateTask']>[1]) {
      const api = useTasksApi()
      const task = await api.updateTask(id, payload)
      const idx = this.allTasks.findIndex((t) => t.id === id)
      if (idx !== -1) this.allTasks[idx] = task
      if (this.currentTask?.id === id) this.currentTask = task
      return task
    },

    updateTaskStatusLocally(taskId: string, status: TaskStatus) {
      const idx = this.allTasks.findIndex((t) => t.id === taskId)
      if (idx !== -1) this.allTasks[idx] = { ...this.allTasks[idx], status }
    },

    async deleteTask(id: string) {
      const api = useTasksApi()
      await api.deleteTask(id)
      this.allTasks = this.allTasks.filter((t) => t.id !== id)
    },

    async assignTask(id: string, devId: number) {
      const api = useTasksApi()
      await api.assignTask(id, devId)
      const idx = this.allTasks.findIndex((t) => t.id === id)
      if (idx !== -1) this.allTasks[idx].assigned_to = devId
    },

    async bulkUpdateStatus(taskIds: string[], status: TaskStatus) {
      const api = useTasksApi()
      await api.bulkUpdateStatus(taskIds, status)
      for (const id of taskIds) {
        const idx = this.allTasks.findIndex((t) => t.id === id)
        if (idx !== -1) this.allTasks[idx].status = status
      }
    },

    async fetchComments(taskId: string) {
      const api = useTasksApi()
      this.comments = await api.getComments(taskId)
    },

    async addComment(taskId: string, content: string) {
      const api = useTasksApi()
      const comment = await api.addComment(taskId, content)
      this.comments.push(comment)
      return comment
    },

    async fetchTimeLogs(taskId: string) {
      const api = useTasksApi()
      this.timeLogs = await api.getTimeLogs(taskId)
    },

    async logTime(taskId: string, minutes: number, description: string) {
      const api = useTasksApi()
      const log = await api.logTime(taskId, minutes, description)
      this.timeLogs.unshift(log)
      return log
    },
  },
})

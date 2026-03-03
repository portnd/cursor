import { useAuth } from '~/composables/useAuth'

export interface Project {
  id: string
  code: string
  name: string
  description: string
  status: 'ACTIVE' | 'COMPLETED' | 'ON_HOLD'
  created_at: string
  updated_at: string
  tasks?: Task[]
}

export interface Sprint {
  id: string
  project_id: string
  name: string
  goal: string
  start_date: string | null
  end_date: string | null
  status: 'PLANNING' | 'ACTIVE' | 'COMPLETED'
  created_at: string
  updated_at: string
  tasks?: Task[]
}

export interface Milestone {
  id: string
  project_id: string
  title: string
  description: string
  due_date: string | null
  status: 'PENDING' | 'REACHED' | 'MISSED'
  created_at: string
}

// Epic represents a large feature/goal (Hierarchy Dimension 1, sits above Tasks)
export interface Epic {
  id: string
  project_id: string
  title: string
  description: string
  status: 'PLANNING' | 'IN_PROGRESS' | 'DONE'
  color: string
  sort_order: number
  start_date: string | null
  end_date: string | null
  created_at: string
  updated_at: string
  tasks?: Task[]
}

export interface EpicTimelineData {
  epics: Epic[]
}

export interface SprintTimelineData {
  sprints: Sprint[]
}

export interface Task {
  id: string
  code: string
  title: string
  description: string
  status: 'PENDING' | 'IN_PROGRESS' | 'REVIEW_PENDING' | 'COMPLETED' | 'BLOCKED'
  priority: 'CRITICAL' | 'HIGH' | 'MEDIUM' | 'LOW'
  story_points: number
  progress: number
  project_id: string | null
  parent_id: string | null
  epic_id: string | null
  sprint_id: string | null
  milestone_id: string | null
  assigned_to: number | null
  created_by: number | null
  due_at: string | null
  start_date: string | null
  end_date: string | null
  started_at: string | null
  completed_at: string | null
  ai_estimated_minutes: number
  sub_tasks?: Task[]
  created_at: string
  updated_at: string
}

export interface ProjectAnalytics {
  project_id: string
  total_tasks: number
  completed_tasks: number
  total_story_points: number
  completed_story_points: number
  total_logged_minutes: number
  avg_cycle_time_days: number
  burndown: { day: string; ideal: number; remaining: number }[]
  velocity: { sprint_name: string; completed_sp: number; planned_sp: number }[]
  team_capacity: {
    user_id: number
    user_email: string
    assigned_tasks: number
    estimated_hours: number
    logged_hours: number
    utilization_pct: number
  }[]
}

function useProjectsApi() {
  const { fetchWithAuth } = useAuth()

  async function getProjects(): Promise<Project[]> {
    const data = await fetchWithAuth<{ data: Project[] }>('/sentinel/projects')
    return data.data || []
  }

  async function getProject(idOrCode: string): Promise<Project> {
    const data = await fetchWithAuth<{ data: Project }>(`/sentinel/projects/${idOrCode}`)
    return data.data
  }

  async function createProject(payload: { name: string; description?: string; status?: string }): Promise<Project> {
    const data = await fetchWithAuth<{ data: Project }>('/sentinel/projects', {
      method: 'POST',
      body: payload,
    })
    return data.data
  }

  async function deleteProject(id: string): Promise<void> {
    await fetchWithAuth(`/sentinel/projects/${id}`, { method: 'DELETE' })
  }

  async function getSprints(projectId: string): Promise<Sprint[]> {
    const data = await fetchWithAuth<{ data: Sprint[] }>(`/sentinel/sprints?project_id=${projectId}`)
    return data.data || []
  }

  async function createSprint(payload: {
    project_id: string
    name: string
    goal?: string
    start_date?: string
    end_date?: string
  }): Promise<Sprint> {
    const data = await fetchWithAuth<{ data: Sprint }>('/sentinel/sprints', {
      method: 'POST',
      body: payload,
    })
    return data.data
  }

  async function updateSprint(id: string, payload: Partial<Sprint>): Promise<Sprint> {
    const data = await fetchWithAuth<{ data: Sprint }>(`/sentinel/sprints/${id}`, {
      method: 'PATCH',
      body: payload,
    })
    return data.data
  }

  async function startSprint(id: string): Promise<Sprint> {
    const data = await fetchWithAuth<{ data: Sprint }>(`/sentinel/sprints/${id}/start`, { method: 'POST' })
    return data.data
  }

  async function completeSprint(id: string): Promise<Sprint> {
    const data = await fetchWithAuth<{ data: Sprint }>(`/sentinel/sprints/${id}/complete`, { method: 'POST' })
    return data.data
  }

  async function reopenSprint(id: string): Promise<Sprint> {
    const data = await fetchWithAuth<{ data: Sprint }>(`/sentinel/sprints/${id}/reopen`, { method: 'POST' })
    return data.data
  }

  async function deleteSprint(id: string): Promise<void> {
    await fetchWithAuth(`/sentinel/sprints/${id}`, { method: 'DELETE' })
  }

  async function addTasksToSprint(sprintId: string, taskIds: string[]): Promise<void> {
    await fetchWithAuth(`/sentinel/sprints/${sprintId}/tasks`, {
      method: 'POST',
      body: { task_ids: taskIds },
    })
  }

  async function getMilestones(projectId: string): Promise<Milestone[]> {
    const data = await fetchWithAuth<{ data: Milestone[] }>(`/sentinel/milestones?project_id=${projectId}`)
    return data.data || []
  }

  async function createMilestone(payload: {
    project_id: string
    title: string
    description?: string
    due_date?: string
  }): Promise<Milestone> {
    const data = await fetchWithAuth<{ data: Milestone }>('/sentinel/milestones', {
      method: 'POST',
      body: payload,
    })
    return data.data
  }

  async function updateMilestone(id: string, payload: Partial<Milestone>): Promise<Milestone> {
    const data = await fetchWithAuth<{ data: Milestone }>(`/sentinel/milestones/${id}`, {
      method: 'PATCH',
      body: payload,
    })
    return data.data
  }

  async function deleteMilestone(id: string): Promise<void> {
    await fetchWithAuth(`/sentinel/milestones/${id}`, { method: 'DELETE' })
  }

  async function getProjectAnalytics(projectId: string): Promise<ProjectAnalytics> {
    const data = await fetchWithAuth<{ data: ProjectAnalytics }>(`/sentinel/projects/${projectId}/analytics`)
    return data.data
  }

  // --- Epic API (Hierarchy Dimension 1) ---

  async function getEpics(projectId: string): Promise<Epic[]> {
    const data = await fetchWithAuth<{ data: Epic[] }>(`/sentinel/epics?project_id=${projectId}`)
    return data.data || []
  }

  async function createEpic(payload: {
    project_id: string
    title: string
    description?: string
    color?: string
    start_date?: string
    end_date?: string
  }): Promise<Epic> {
    const data = await fetchWithAuth<{ data: Epic }>('/sentinel/epics', {
      method: 'POST',
      body: payload,
    })
    return data.data
  }

  async function updateEpic(id: string, payload: Partial<Epic>): Promise<Epic> {
    const data = await fetchWithAuth<{ data: Epic }>(`/sentinel/epics/${id}`, {
      method: 'PATCH',
      body: payload,
    })
    return data.data
  }

  async function deleteEpic(id: string): Promise<void> {
    await fetchWithAuth(`/sentinel/epics/${id}`, { method: 'DELETE' })
  }

  // --- Timeline Views (Matrix Dimension) ---

  async function getEpicTimelineData(projectId: string): Promise<EpicTimelineData> {
    const data = await fetchWithAuth<{ data: EpicTimelineData }>(`/sentinel/projects/${projectId}/timeline/epic-view`)
    return data.data || { epics: [] }
  }

  async function getSprintTimelineData(projectId: string): Promise<SprintTimelineData> {
    const data = await fetchWithAuth<{ data: SprintTimelineData }>(`/sentinel/projects/${projectId}/timeline/sprint-view`)
    return data.data || { sprints: [] }
  }

  return {
    getProjects,
    getProject,
    createProject,
    deleteProject,
    getSprints,
    createSprint,
    updateSprint,
    startSprint,
    completeSprint,
    reopenSprint,
    deleteSprint,
    addTasksToSprint,
    getMilestones,
    createMilestone,
    updateMilestone,
    deleteMilestone,
    getProjectAnalytics,
    getEpics,
    createEpic,
    updateEpic,
    deleteEpic,
    getEpicTimelineData,
    getSprintTimelineData,
  }
}

export { useProjectsApi }

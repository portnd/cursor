import { useAuth } from '~/composables/useAuth'

export interface Project {
  id: string
  code: string
  name: string
  description: string
  status: 'ACTIVE' | 'COMPLETED' | 'ON_HOLD'
  color?: string
  // Squad Model
  team_id?: number | null
  team_name?: string
  // Internal VC — Project Capital
  capital_balance: number
  bonus_percentage: number
  created_at: string
  updated_at: string
  tasks?: Task[]
  /** Task counts from list API (so cards show correct numbers without loading full tasks) */
  task_total?: number
  task_completed?: number
  task_overdue?: number
  /** When teams/squads are off — CEO-assigned Product Owner users for this project (API field: pm_owners) */
  pm_owners?: ProjectPmOwner[]
}

export interface ProjectPmOwner {
  user_id: number
  email: string
  display_name?: string
}

export interface ProjectTransaction {
  id: number
  project_id: string
  type: 'INJECTION' | 'BURN' | 'BONUS_PAYOUT' | 'ADJUSTMENT'
  amount: number
  reference: string
  created_at: string
}

export interface ProjectCapitalResponse {
  project_id: string
  project_name: string
  team_id?: number | null
  team_monthly_cost: number
  capital_balance: number
  bonus_percentage: number
  runway_months: number
  transactions: ProjectTransaction[]
}

export interface InjectProjectCapitalPayload {
  amount: number
  bonus_percentage?: number
  note?: string
}

export interface EditProjectCapitalPayload {
  new_balance: number
  bonus_percentage?: number | null
  note?: string
}

export interface CloseProjectCycleResponse {
  project_id: string
  balance_before: number
  bonus_percentage: number
  bonus_amount: number
  balance_after: number
}

export interface ProjectBackup {
  id: string
  project_id: string
  label: string
  payload?: Record<string, unknown>
  created_by: number | null
  created_at: string
}

export interface Sprint {
  id: string
  project_id: string
  /** Populated on dev “my sprints” payloads */
  project_name?: string
  project_code?: string
  name: string
  goal: string
  start_date: string | null
  end_date: string | null
  status: 'PLANNING' | 'ACTIVE' | 'COMPLETED'
  sort_order?: number
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

export interface AIGeneratedPlan {
  epics: { title: string; description: string; color: string }[]
  milestones: { title: string; description: string; due_date: string }[]
  sprints: { name: string; goal: string; start_date: string; end_date: string }[]
  tasks: {
    title: string
    description: string
    priority: string
    story_points: number
    epic_index?: number | null
    sprint_index?: number | null
    milestone_index?: number | null
    start_date: string
    end_date: string
  }[]
}

export interface Task {
  id: string
  code: string
  title: string
  description: string
  task_type: 'FEATURE' | 'TASK' | 'BUG'
  status: 'PENDING' | 'IN_PROGRESS' | 'READY_FOR_TEST' | 'REVIEW_PENDING' | 'COMPLETED' | 'BLOCKED'
  priority: 'CRITICAL' | 'HIGH' | 'MEDIUM' | 'LOW'
  story_points: number
  progress: number
  project_id: string | null
  parent_id: string | null
  epic_id: string | null
  sprint_id: string | null
  milestone_id: string | null
  sort_order: number
  assigned_to: number | null
  assigned_to_display_name?: string
  assigned_to_email?: string
  /** GET /sentinel/tasks/my — display enrichment */
  project_name?: string
  project_color?: string
  sprint_name?: string
  effective_sprint_id?: string | null
  created_by: number | null
  due_at: string | null
  start_date: string | null
  end_date: string | null
  started_at: string | null
  completed_at: string | null
  estimated_minutes: number
  sub_tasks?: Task[]
  is_komgrip?: boolean
  created_at: string
  updated_at: string
}

export interface ProjectDetailsTasksMeta {
  limit: number
  returned: number
  has_more: boolean
}

export interface ProjectTaskPageCursor {
  created_at: string
  id: string
}

export interface ProjectTasksPageResponse {
  tasks: Task[]
  limit: number
  returned: number
  has_more: boolean
  next_cursor?: ProjectTaskPageCursor
  next_offset?: number
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
    user_display_name?: string
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
    const data = await fetchWithAuth<{ data: Project }>(`/sentinel/projects/${encodeURIComponent(idOrCode)}`)
    return data.data
  }

  /** Combined project + tasks + sprints + milestones + epics (1 round-trip, use for project page). */
  async function getProjectDetails(idOrCode: string, opts?: { tasksLimit?: number }): Promise<{ project: Project; tasks: Task[]; tasks_meta: ProjectDetailsTasksMeta; sprints: Sprint[]; milestones: Milestone[]; epics: Epic[] }> {
    const q = typeof opts?.tasksLimit === 'number' && opts.tasksLimit > 0
      ? `?tasks_limit=${Math.floor(opts.tasksLimit)}`
      : ''
    const data = await fetchWithAuth<{ data: { project: Project; tasks: Task[]; tasks_meta: ProjectDetailsTasksMeta; sprints: Sprint[]; milestones: Milestone[]; epics: Epic[] } }>(`/sentinel/projects/${encodeURIComponent(idOrCode)}/details${q}`)
    return data.data
  }

  /** Load additional project tasks page (page 2+) using cursor or offset. */
  async function getProjectTasksPage(idOrCode: string, opts?: { limit?: number; cursorCreatedAt?: string; cursorId?: string; offset?: number }): Promise<ProjectTasksPageResponse> {
    const params = new URLSearchParams()
    if (typeof opts?.limit === 'number' && opts.limit > 0) params.set('limit', String(Math.floor(opts.limit)))
    if (opts?.cursorCreatedAt && opts?.cursorId) {
      params.set('cursor_created_at', opts.cursorCreatedAt)
      params.set('cursor_id', opts.cursorId)
    } else if (typeof opts?.offset === 'number' && opts.offset >= 0) {
      params.set('offset', String(Math.floor(opts.offset)))
    }
    const qs = params.toString()
    const data = await fetchWithAuth<{ data: ProjectTasksPageResponse }>(`/sentinel/projects/${encodeURIComponent(idOrCode)}/tasks${qs ? `?${qs}` : ''}`)
    return data.data
  }

  async function createProject(payload: { name: string; description?: string; status?: string }): Promise<Project> {
    const data = await fetchWithAuth<{ data: Project }>('/sentinel/projects', {
      method: 'POST',
      body: payload,
    })
    return data.data
  }

  async function updateProject(idOrCode: string, payload: { name: string; description?: string; status?: string; update_code?: boolean }): Promise<Project> {
    const data = await fetchWithAuth<{ data: Project }>(`/sentinel/projects/${idOrCode}`, {
      method: 'PATCH',
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

  /** projectIdOrCode: UUID or project code (e.g. mims-hd-map) — both work for parallel load. */
  async function getEpicTimelineData(projectIdOrCode: string): Promise<EpicTimelineData> {
    const data = await fetchWithAuth<{ data: EpicTimelineData }>(`/sentinel/projects/${encodeURIComponent(projectIdOrCode)}/timeline/epic-view`)
    return data.data || { epics: [] }
  }

  /** projectIdOrCode: UUID or project code (e.g. mims-hd-map) — both work for parallel load. */
  async function getSprintTimelineData(projectIdOrCode: string): Promise<SprintTimelineData> {
    const data = await fetchWithAuth<{ data: SprintTimelineData }>(`/sentinel/projects/${encodeURIComponent(projectIdOrCode)}/timeline/sprint-view`)
    return data.data || { sprints: [] }
  }

  /** Run AI estimate on a task; updates task.estimated_minutes. Returns updated task. Creator / CEO / Product Owner only. */
  async function estimateTask(taskIdOrCode: string): Promise<Task> {
    const data = await fetchWithAuth<{ data: Task }>(`/sentinel/tasks/${encodeURIComponent(taskIdOrCode)}/estimate`, {
      method: 'POST',
    })
    return data.data
  }

  /** Clear project plan: remove all tasks, sprints, milestones, epics. CEO / Product Owner only. */
  async function clearProjectPlan(projectIdOrCode: string): Promise<void> {
    await fetchWithAuth(`/sentinel/projects/${encodeURIComponent(projectIdOrCode)}/clear-plan`, {
      method: 'POST',
    })
  }

  /** AI Agent: estimate time + arrange timeline for existing tasks (no new tasks created). CEO / Product Owner only. */
  async function scheduleProjectWithAI(projectIdOrCode: string): Promise<{ message: string; updated: number }> {
    return fetchWithAuth<{ message: string; updated: number }>(
      `/sentinel/projects/${encodeURIComponent(projectIdOrCode)}/ai-schedule`,
      { method: 'POST', timeoutMs: 120000 }
    )
  }

  /** (Optional) Generate new epics, milestones, sprints, tasks from project name/description. CEO / Product Owner only. */
  async function generateProjectPlan(projectIdOrCode: string): Promise<AIGeneratedPlan> {
    const data = await fetchWithAuth<{ data: AIGeneratedPlan }>(
      `/sentinel/projects/${encodeURIComponent(projectIdOrCode)}/ai-plan`,
      { method: 'POST', timeoutMs: 120000 }
    )
    return data.data
  }

  // --- Project Finance (Internal VC — per-project capital) ---

  async function getProjectCapital(projectId: string): Promise<ProjectCapitalResponse> {
    const data = await fetchWithAuth<{ data: ProjectCapitalResponse }>(`/sentinel/projects/${projectId}/finance/capital`)
    return data.data
  }

  async function getProjectCapitals(projectIds: string[]): Promise<ProjectCapitalResponse[]> {
    if (!projectIds.length) return []
    const uniqueIds = [...new Set(projectIds.filter(Boolean))]
    const query = encodeURIComponent(uniqueIds.join(','))
    const data = await fetchWithAuth<{ data: ProjectCapitalResponse[] }>(`/sentinel/projects/finance/capital?project_ids=${query}`)
    return data.data || []
  }

  async function injectProjectCapital(projectId: string, payload: InjectProjectCapitalPayload): Promise<Project> {
    const data = await fetchWithAuth<{ data: Project }>(`/sentinel/projects/${projectId}/finance/inject`, {
      method: 'POST',
      body: payload,
    })
    return data.data
  }

  async function editProjectCapital(projectId: string, payload: EditProjectCapitalPayload): Promise<Project> {
    const data = await fetchWithAuth<{ data: Project }>(`/sentinel/projects/${projectId}/finance/capital`, {
      method: 'PUT',
      body: payload,
    })
    return data.data
  }

  async function closeProjectCycle(projectId: string): Promise<CloseProjectCycleResponse> {
    const data = await fetchWithAuth<{ data: CloseProjectCycleResponse }>(`/sentinel/projects/${projectId}/finance/close-cycle`, {
      method: 'POST',
    })
    return data.data
  }

  async function deleteProjectTransaction(projectId: string, txId: number): Promise<void> {
    await fetchWithAuth(`/sentinel/projects/${projectId}/finance/transactions/${txId}`, {
      method: 'DELETE',
    })
  }

  // --- Project Backups ---

  async function getProjectBackups(projectId: string): Promise<ProjectBackup[]> {
    const data = await fetchWithAuth<{ data: ProjectBackup[] }>(`/sentinel/projects/${projectId}/backups`)
    return data.data || []
  }

  async function createProjectBackup(projectId: string, label?: string): Promise<ProjectBackup> {
    const data = await fetchWithAuth<{ data: ProjectBackup }>(`/sentinel/projects/${projectId}/backups`, {
      method: 'POST',
      body: { label: label || '' },
    })
    return data.data
  }

  async function restoreProjectBackup(projectId: string, backupId: string): Promise<void> {
    await fetchWithAuth(`/sentinel/projects/${projectId}/backups/${backupId}/restore`, {
      method: 'POST',
    })
  }

  async function deleteProjectBackup(projectId: string, backupId: string): Promise<void> {
    await fetchWithAuth(`/sentinel/projects/${projectId}/backups/${backupId}`, {
      method: 'DELETE',
    })
  }

  async function getProjectBackupPayload(projectId: string, backupId: string): Promise<Record<string, unknown>> {
    const data = await fetchWithAuth<{ data: Record<string, unknown> }>(`/sentinel/projects/${projectId}/backups/${backupId}/payload`)
    return data.data
  }

  async function importProjectFromBackup(name: string, payload: Record<string, unknown>): Promise<Project> {
    const data = await fetchWithAuth<{ data: Project }>('/sentinel/projects/import-backup', {
      method: 'POST',
      body: { name, payload },
    })
    return data.data
  }

  return {
    getProjects,
    getProject,
    getProjectDetails,
    getProjectTasksPage,
    createProject,
    updateProject,
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
    estimateTask,
    clearProjectPlan,
    scheduleProjectWithAI,
    generateProjectPlan,
    getProjectCapital,
    getProjectCapitals,
    injectProjectCapital,
    editProjectCapital,
    closeProjectCycle,
    deleteProjectTransaction,
    getProjectBackups,
    createProjectBackup,
    restoreProjectBackup,
    deleteProjectBackup,
    getProjectBackupPayload,
    importProjectFromBackup,
  }
}

export { useProjectsApi }

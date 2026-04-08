import { useAuth } from '~/composables/useAuth'
import type { Task } from '../../../modules/projects/infrastructure/projects-api'

export interface PPTXPreviewSlide { index: number; title: string; hidden?: boolean; suggested_task_title?: string }
export interface PPTXPreviewResult { title: string; slides: PPTXPreviewSlide[] }
export interface PPTXImportResult { created_count: number; page_count: number; title: string; tasks: Task[] }

export interface TaskComment {
  id: string
  task_id: string
  user_id: number
  user_email: string
  content: string
  attachments?: Array<{
    file_name: string
    mime_type: string
    size: number
    data_url: string
    is_image: boolean
  }>
  created_at: string
}

export interface TimeLog {
  id: string
  task_id: string
  user_id: number
  user_email: string
  minutes: number
  description: string
  work_type: string      // DEV | REVIEW | TESTING | MEETING | RESEARCH | OTHER
  logged_date: string    // YYYY-MM-DD
  is_timer_session: boolean
  logged_at: string
}

export interface DailyTimeLogSummary {
  date: string
  total_minutes: number
  entries: TimeLog[]
}

export const WORK_TYPES = [
  { value: 'DEV',      label: 'Dev',      emoji: '💻' },
  { value: 'REVIEW',   label: 'Review',   emoji: '👁' },
  { value: 'TESTING',  label: 'Test',     emoji: '🧪' },
  { value: 'MEETING',  label: 'Meeting',  emoji: '📅' },
  { value: 'RESEARCH', label: 'Research', emoji: '🔬' },
  { value: 'OTHER',    label: 'Other',    emoji: '📌' },
] as const

export interface BulkLogEntry {
  task_id: string
  minutes: number
  description?: string
  work_type?: string
  logged_date?: string
}

export interface BulkLogResult {
  task_id: string
  success: boolean
  log?: TimeLog
  error?: string
}

export interface BulkLogResponse {
  message: string
  success_count: number
  total: number
  results: BulkLogResult[]
}

export interface GlobalActiveTask extends Task {
  project_name: string
  project_color: string
  assigned_to_display_name?: string
  assigned_to_email?: string
}

export interface UATPayload {
  staging_url: string
  test_credentials: string
  release_notes: string
}

export interface FeatureRoadmapItem extends Task {
  project_name: string
  project_color: string
  project_code: string
  rollup_progress: number
  child_tasks: Task[]
  uat_payload?: UATPayload
}

function useTasksApi() {
  const { fetchWithAuth } = useAuth()

  async function getTasksByProject(projectId: string, opts?: { sprintId?: string }): Promise<Task[]> {
    const params = new URLSearchParams({ project_id: projectId })
    if (opts?.sprintId) params.set('sprint_id', opts.sprintId)
    const data = await fetchWithAuth<{ data: Task[] }>(`/sentinel/tasks?${params.toString()}`)
    return data.data || []
  }

  async function getAllTasks(): Promise<Task[]> {
    const data = await fetchWithAuth<{ data: Task[] }>('/sentinel/tasks')
    return data.data || []
  }

  async function getTask(idOrCode: string): Promise<Task> {
    const data = await fetchWithAuth<{ data: Task }>(`/sentinel/tasks/${idOrCode}`)
    return data.data
  }

  async function createTask(payload: {
    title: string
    description?: string
    task_type?: string
    project_id?: string
    parent_id?: string
    priority?: string
    story_points?: number
    sprint_id?: string
    milestone_id?: string
    epic_id?: string
    start_date?: string
    end_date?: string
    due_date?: string
    estimated_minutes?: number
  }): Promise<Task> {
    const data = await fetchWithAuth<{ data: Task }>('/sentinel/tasks', {
      method: 'POST',
      body: payload,
    })
    return data.data
  }

  async function updateTask(id: string, payload: Partial<{
    title: string
    description: string
    priority: string
    story_points: number
    sprint_id: string
    milestone_id: string
    parent_id: string
    epic_id: string
    sort_order: number
    start_date: string
    end_date: string
    progress: number
    status: string
    estimated_minutes: number
  }>): Promise<Task> {
    const data = await fetchWithAuth<{ data: Task }>(`/sentinel/tasks/${id}`, {
      method: 'PATCH',
      body: payload,
    })
    return data.data
  }

  async function updateTaskSlideResources(id: string, resourceUrls: Record<string, unknown>): Promise<Task> {
    const data = await fetchWithAuth<{ data: Task }>(`/sentinel/tasks/${id}/slide-resources`, {
      method: 'PATCH',
      body: { resource_urls: resourceUrls },
    })
    return data.data
  }

  async function deleteTask(id: string): Promise<void> {
    await fetchWithAuth(`/sentinel/tasks/${id}`, { method: 'DELETE' })
  }

  async function assignTask(id: string, devId: number): Promise<void> {
    await fetchWithAuth(`/sentinel/tasks/${id}/assign`, {
      method: 'POST',
      body: { dev_id: devId },
    })
  }

  async function bulkUpdateStatus(taskIds: string[], status: string): Promise<void> {
    await fetchWithAuth('/sentinel/tasks/bulk-status', {
      method: 'PATCH',
      body: { task_ids: taskIds, status },
    })
  }

  async function getComments(taskId: string): Promise<TaskComment[]> {
    const data = await fetchWithAuth<{ data: TaskComment[] }>(`/sentinel/tasks/${taskId}/comments`)
    return data.data || []
  }

  async function addComment(taskId: string, content: string, attachments: File[] = []): Promise<TaskComment> {
    const hasAttachments = attachments.length > 0
    const body = hasAttachments
      ? (() => {
        const formData = new FormData()
        formData.append('content', content)
        for (const file of attachments) formData.append('attachments', file)
        return formData
      })()
      : { content }
    const data = await fetchWithAuth<{ data: TaskComment }>(`/sentinel/tasks/${taskId}/comments`, {
      method: 'POST',
      body,
    })
    return data.data
  }

  async function getTimeLogs(taskId: string): Promise<TimeLog[]> {
    const data = await fetchWithAuth<{ data: TimeLog[] }>(`/sentinel/tasks/${taskId}/time-logs`)
    return data.data || []
  }

  async function logTime(
    taskId: string,
    minutes: number,
    description: string,
    workType = 'DEV',
    loggedDate?: string,
    isTimerSession = false,
  ): Promise<TimeLog> {
    const data = await fetchWithAuth<{ data: TimeLog }>(`/sentinel/tasks/${taskId}/time-logs`, {
      method: 'POST',
      body: { minutes, description, work_type: workType, logged_date: loggedDate, is_timer_session: isTimerSession },
    })
    return data.data
  }

  async function editTimeLog(logId: string, minutes: number, description: string, workType: string): Promise<TimeLog> {
    const data = await fetchWithAuth<{ data: TimeLog }>(`/sentinel/time-logs/${logId}`, {
      method: 'PATCH',
      body: { minutes, description, work_type: workType },
    })
    return data.data
  }

  async function deleteTimeLog(logId: string): Promise<void> {
    await fetchWithAuth(`/sentinel/time-logs/${logId}`, { method: 'DELETE' })
  }

  async function bulkLogTime(entries: BulkLogEntry[]): Promise<BulkLogResponse> {
    const data = await fetchWithAuth<BulkLogResponse>('/sentinel/time-logs/bulk', {
      method: 'POST',
      body: JSON.stringify(entries),
    })
    return data
  }

  async function getMyDailyTimeLogs(date?: string): Promise<DailyTimeLogSummary> {
    const q = date ? `?date=${date}` : ''
    const data = await fetchWithAuth<{ data: DailyTimeLogSummary }>(`/sentinel/users/me/time-logs${q}`)
    return data.data
  }

  async function getGanttData(projectId?: string): Promise<{ tasks: Task[]; dependencies: any[] }> {
    const url = projectId
      ? `/sentinel/tasks/gantt?project_id=${projectId}`
      : '/sentinel/tasks/gantt'
    const data = await fetchWithAuth<{ data: { tasks: Task[]; dependencies: any[] } }>(url)
    return data.data || { tasks: [], dependencies: [] }
  }

  async function getGlobalActiveTasks(): Promise<GlobalActiveTask[]> {
    const data = await fetchWithAuth<{ data: GlobalActiveTask[] }>('/sentinel/tasks/my-global-active')
    return data.data || []
  }

  async function getTeamActiveTasks(): Promise<GlobalActiveTask[]> {
    const data = await fetchWithAuth<{ data: GlobalActiveTask[] }>('/sentinel/tasks/team-active')
    return data.data || []
  }

  async function getActiveFeatures(projectId?: string): Promise<FeatureRoadmapItem[]> {
    const q = projectId ? `?project_id=${encodeURIComponent(projectId)}` : ''
    const data = await fetchWithAuth<{ data: FeatureRoadmapItem[] }>(`/sentinel/tasks/features${q}`)
    return data.data || []
  }

  async function previewGoogleSlides(payload: {
    presentation_url: string
    api_key?: string
  }): Promise<{
    presentation_title: string
    slides: { index: number; title: string; hidden?: boolean; suggested_task_title?: string }[]
    import_mode: string
    api_key_status: string
    api_key_error?: string
  }> {
    const data = await fetchWithAuth<{
      data: {
        presentation_title: string
        presentation_id?: string
        slides: { index: number; title: string; hidden?: boolean; suggested_task_title?: string }[]
        already_imported_slide_indices?: number[]
        import_mode: string
        api_key_status: string
        api_key_error?: string
      }
    }>('/sentinel/import/google-slides/preview', { method: 'POST', body: payload, timeoutMs: 60 * 1000 })
    return data.data
  }

  async function importGoogleSlides(payload: {
    presentation_url: string
    project_id: string
    sprint_id?: string
    epic_id?: string
    parent_id?: string
    api_key?: string
    priority?: string
    story_points?: number
    slide_indices?: number[]
    slides?: {
      slide_index: number
      title: string
      assignee_id?: number | null
      estimated_minutes: number
      priority: string
    }[]
  }): Promise<{ created_count: number; slide_count: number; presentation_title: string; tasks: Task[] }> {
    const data = await fetchWithAuth<{ data: { created_count: number; slide_count: number; presentation_title: string; tasks: Task[] } }>(
      '/sentinel/import/google-slides',
      { method: 'POST', body: payload, timeoutMs: 5 * 60 * 1000 }, // 5 min: download PPTX + slide images + create tasks
    )
    return data.data
  }

  async function previewGoogleSheets(payload: { sheet_url: string }): Promise<{
    sheet_title: string
    sheet_id: string
    rows: {
      row_index: number
      title: string
      due_date: string
      status: string
      raw_status: string
      notes: string
      header?: string
      header_link?: string
      request_method?: string
      payload?: string
      image_ref?: string
      detail_links?: string[]
    }[]
  }> {
    const data = await fetchWithAuth<{
      data: {
        sheet_title: string
        sheet_id: string
        rows: {
          row_index: number
          title: string
          due_date: string
          status: string
          raw_status: string
          notes: string
          header?: string
          header_link?: string
          request_method?: string
          payload?: string
          image_ref?: string
          detail_links?: string[]
        }[]
      }
    }>('/sentinel/import/google-sheets/preview', { method: 'POST', body: payload, timeoutMs: 60 * 1000 })
    return data.data
  }

  async function previewPPTXUpload(file: File): Promise<PPTXPreviewResult> {
    const formData = new FormData()
    formData.append('file', file)
    const data = await fetchWithAuth<{ data: PPTXPreviewResult }>('/sentinel/import/pptx/preview', {
      method: 'POST',
      body: formData,
      timeoutMs: 60 * 1000,
    })
    return data.data
  }

  async function importPPTXUpload(
    file: File,
    payload: {
      project_id: string
      sprint_id?: string
      epic_id?: string
      parent_id?: string
      priority?: string
      story_points?: number
      pages: { slide_index: number; title: string; assignee_id?: number | null; estimated_minutes: number; priority: string }[]
    },
  ): Promise<PPTXImportResult> {
    const formData = new FormData()
    formData.append('file', file)
    formData.append('payload', JSON.stringify(payload))
    const data = await fetchWithAuth<{ data: PPTXImportResult }>('/sentinel/import/pptx', {
      method: 'POST',
      body: formData,
      timeoutMs: 3 * 60 * 1000,
    })
    return data.data
  }

  async function importGoogleSheets(payload: {
    sheet_url: string
    sheet_title?: string
    project_id: string
    sprint_id?: string
    epic_id?: string
    parent_id?: string
    rows: {
      row_index: number
      title: string
      priority: string
      estimated_minutes: number
      due_date: string
      status: string
      notes: string
      header?: string
      header_link?: string
      request_method?: string
      payload?: string
      image_ref?: string
      detail_links?: string[]
    }[]
  }): Promise<{ created_count: number; sheet_title: string; tasks: Task[] }> {
    const data = await fetchWithAuth<{ data: { created_count: number; sheet_title: string; tasks: Task[] } }>(
      '/sentinel/import/google-sheets',
      { method: 'POST', body: payload, timeoutMs: 2 * 60 * 1000 },
    )
    return data.data
  }

  async function splitTask(taskId: string, splits: {
    title: string
    estimated_minutes: number
    assignee_id?: number | null
    priority?: string
  }[]): Promise<Task[]> {
    const data = await fetchWithAuth<{ data: Task[] }>(
      `/sentinel/tasks/${taskId}/split`,
      { method: 'POST', body: { splits } }
    )
    return data.data ?? []
  }

  async function submitUAT(taskId: string, payload: UATPayload): Promise<void> {
    await fetchWithAuth(`/sentinel/tasks/${taskId}/submit-uat`, {
      method: 'POST',
      body: payload,
    })
  }

  async function approveTask(taskId: string): Promise<void> {
    await fetchWithAuth(`/sentinel/tasks/${taskId}/approve`, { method: 'POST' })
  }

  async function rejectTask(taskId: string, reason: string): Promise<void> {
    await fetchWithAuth(`/sentinel/tasks/${taskId}/reject`, {
      method: 'POST',
      body: { reason },
    })
  }

  async function markReadyForTest(taskId: string): Promise<void> {
    await fetchWithAuth(`/sentinel/tasks/${taskId}/ready-for-test`, { method: 'POST' })
  }

  /** Product Owner first-stage approval: READY_FOR_TEST → READY_FOR_UAT, attaches test evidence for CEO. */
  async function pmApproveSubTask(taskId: string, testUrl: string, testSteps: string): Promise<void> {
    await fetchWithAuth(`/sentinel/tasks/${taskId}/pm-approve-sub`, {
      method: 'POST',
      body: { test_url: testUrl, test_steps: testSteps },
    })
  }

  /** CEO final approval: READY_FOR_UAT → COMPLETED. */
  async function approveSubTask(taskId: string): Promise<void> {
    await fetchWithAuth(`/sentinel/tasks/${taskId}/approve-sub`, { method: 'POST' })
  }

  async function rejectSubTask(taskId: string, reason: string): Promise<void> {
    await fetchWithAuth(`/sentinel/tasks/${taskId}/reject-sub`, {
      method: 'POST',
      body: { reason },
    })
  }

  async function getTasksReadyForTest(): Promise<GlobalActiveTask[]> {
    const res = await fetchWithAuth<{ data: GlobalActiveTask[] }>('/sentinel/tasks/ready-for-test')
    return res.data ?? []
  }

  /** CEO: fetch TASK/BUG tasks in READY_FOR_UAT awaiting final approval. */
  async function getTasksReadyForCEOApproval(): Promise<GlobalActiveTask[]> {
    const res = await fetchWithAuth<{ data: GlobalActiveTask[] }>('/sentinel/tasks/ceo-approval-queue')
    return res.data ?? []
  }

  return {
    getTasksByProject,
    getAllTasks,
    getTask,
    createTask,
    updateTask,
    updateTaskSlideResources,
    deleteTask,
    assignTask,
    bulkUpdateStatus,
    getComments,
    addComment,
    getTimeLogs,
    logTime,
    editTimeLog,
    deleteTimeLog,
    bulkLogTime,
    getMyDailyTimeLogs,
    getGanttData,
    getGlobalActiveTasks,
    getTeamActiveTasks,
    getActiveFeatures,
    previewGoogleSlides,
    importGoogleSlides,
    previewGoogleSheets,
    importGoogleSheets,
    previewPPTXUpload,
    importPPTXUpload,
    splitTask,
    submitUAT,
    approveTask,
    rejectTask,
    markReadyForTest,
    pmApproveSubTask,
    approveSubTask,
    rejectSubTask,
    getTasksReadyForTest,
    getTasksReadyForCEOApproval,
  }
}

export { useTasksApi }

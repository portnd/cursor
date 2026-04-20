import { computed, ref, type Ref } from 'vue'
import type { Project, Sprint, Task } from '~/core/modules/projects/infrastructure/projects-api'
import type { useProjectsApi } from '~/core/modules/projects/infrastructure/projects-api'
import type { useTasksApi } from '~/core/modules/tasks/infrastructure/tasks-api'

type ProjectsApi = ReturnType<typeof useProjectsApi>
type TasksApi = ReturnType<typeof useTasksApi>

interface UseProjectSprintsOptions {
  project: Ref<Project | null>
  sprints: Ref<Sprint[]>
  allTasks: Ref<Task[]>
  activeSprint: Ref<Sprint | null>
  projectsApi: ProjectsApi
  tasksApi: TasksApi
  loadAll: () => Promise<void>
  loadSprintTimeline: () => Promise<void>
  updateTaskInTimelineData: (taskId: string, startDate: string, endDate: string) => void
  taskDatesInSprintRange: (
    task: { start_date?: string | null; end_date?: string | null; due_at?: string | null },
    sprint: { start_date: string | null; end_date: string | null },
  ) => { start_date: string; end_date: string } | null
  showError: (message: string, title?: string) => void
}

export function useProjectSprints(options: UseProjectSprintsOptions) {
  const {
    project,
    sprints,
    allTasks,
    activeSprint,
    projectsApi,
    tasksApi,
    loadAll,
    loadSprintTimeline,
    updateTaskInTimelineData,
    taskDatesInSprintRange,
    showError,
  } = options

  function isoToDatetimeLocal(iso: string | null | undefined): string {
    if (!iso) return ''
    const d = new Date(iso)
    if (isNaN(d.getTime())) return ''
    const y = d.getFullYear()
    const m = String(d.getMonth() + 1).padStart(2, '0')
    const day = String(d.getDate()).padStart(2, '0')
    const h = String(d.getHours()).padStart(2, '0')
    const min = String(d.getMinutes()).padStart(2, '0')
    return `${y}-${m}-${day}T${h}:${min}`
  }

  function parseDatetimeLocal(s: string): Date | null {
    if (!s?.trim()) return null
    const d = new Date(s)
    return isNaN(d.getTime()) ? null : d
  }

  function dateToDatetimeLocal(d: Date): string {
    const y = d.getFullYear()
    const m = String(d.getMonth() + 1).padStart(2, '0')
    const day = String(d.getDate()).padStart(2, '0')
    const h = String(d.getHours()).padStart(2, '0')
    const min = String(d.getMinutes()).padStart(2, '0')
    return `${y}-${m}-${day}T${h}:${min}`
  }

  function nextSprintMondayStart(base: Date = new Date(), hour = 9, minute = 0): Date {
    const d = new Date(base)
    d.setSeconds(0, 0)
    d.setMilliseconds(0)
    d.setHours(hour, minute, 0, 0)
    const dow = d.getDay()
    let add = 0
    if (dow === 0) add = 1
    else if (dow !== 1) add = 8 - dow
    d.setDate(d.getDate() + add)
    return d
  }

  function startOfLocalDay(d: Date): Date {
    const x = new Date(d)
    x.setHours(0, 0, 0, 0)
    x.setMilliseconds(0)
    return x
  }

  function addLocalCalendarDays(d: Date, n: number): Date {
    const x = new Date(d)
    x.setDate(x.getDate() + n)
    return x
  }

  function mondayOnOrAfterCalendarDay(day: Date, hour = 9, minute = 0): Date {
    const x = startOfLocalDay(day)
    const dow = x.getDay()
    let add = 0
    if (dow === 0) add = 1
    else if (dow !== 1) add = 8 - dow
    x.setDate(x.getDate() + add)
    x.setHours(hour, minute, 0, 0)
    x.setSeconds(0, 0)
    return x
  }

  function sprintEffectiveEndTime(s: Sprint): number | null {
    if (s.end_date) {
      const t = new Date(s.end_date).getTime()
      return Number.isNaN(t) ? null : t
    }
    if (s.start_date) {
      const st = new Date(s.start_date)
      if (Number.isNaN(st.getTime())) return null
      const end = addLocalCalendarDays(startOfLocalDay(st), 13)
      end.setHours(17, 0, 0, 0)
      return end.getTime()
    }
    return null
  }

  function defaultNextSprintMondayStart(existing: Sprint[], now: Date = new Date()): Date {
    const fromToday = nextSprintMondayStart(now)
    let latestEndMs: number | null = null
    for (const sp of existing) {
      const t = sprintEffectiveEndTime(sp)
      if (t == null) continue
      if (latestEndMs == null || t > latestEndMs) latestEndMs = t
    }
    if (latestEndMs == null) return fromToday
    const latestEnd = new Date(latestEndMs)
    const dayAfterLast = addLocalCalendarDays(startOfLocalDay(latestEnd), 1)
    const afterPreviousSprints = mondayOnOrAfterCalendarDay(dayAfterLast)
    return new Date(Math.max(fromToday.getTime(), afterPreviousSprints.getTime()))
  }

  function sprintEndFromMondayStart(start: Date, weeks: number): Date {
    const w = Math.max(1, Math.floor(weeks))
    const end = new Date(start)
    end.setDate(end.getDate() + w * 7 - 1)
    end.setHours(17, 0, 0, 0)
    return end
  }

  function sprintInclusiveCalendarDays(start: Date, end: Date): number {
    const s = new Date(start.getFullYear(), start.getMonth(), start.getDate())
    const e = new Date(end.getFullYear(), end.getMonth(), end.getDate())
    return Math.round((e.getTime() - s.getTime()) / 86400000) + 1
  }

  function formatSprintDateRangeForName(start: Date, end: Date): string {
    const sy = start.getFullYear()
    const ey = end.getFullYear()
    const short: Intl.DateTimeFormatOptions = { month: 'short', day: 'numeric' }
    const startPart =
      sy === ey
        ? start.toLocaleDateString('en-US', short)
        : start.toLocaleDateString('en-US', { ...short, year: 'numeric' })
    const endPart = end.toLocaleDateString('en-US', { ...short, year: 'numeric' })
    return `${startPart} – ${endPart}`
  }

  function defaultSprintDisplayName(projectName: string, ordinal: number, start?: Date | null, end?: Date | null): string {
    const base = `Sprint ${ordinal}`
    if (start && end && !isNaN(start.getTime()) && !isNaN(end.getTime())) {
      return `${base} (${formatSprintDateRangeForName(start, end)})`
    }
    return base
  }

  function nextSprintOrdinal(existing: Sprint[]): number {
    return existing.length + 1
  }

  const sprintsOrdered = computed(() =>
    [...sprints.value].sort(
      (a, b) =>
        (a.sort_order ?? 0) - (b.sort_order ?? 0) ||
        new Date(a.created_at).getTime() - new Date(b.created_at).getTime(),
    ),
  )

  function getSprintStats(sprintId: string) {
    const tasks = allTasks.value.filter((t) => !t.parent_id && t.sprint_id === sprintId)
    return {
      total: tasks.length,
      done: tasks.filter((t) => t.status === 'COMPLETED').length,
      sp: tasks.reduce((s, t) => s + (t.story_points || 0), 0),
    }
  }

  const showSprintModal = ref(false)
  const editingSprint = ref<Sprint | null>(null)
  const sprintForm = ref({ name: '', goal: '', start_date: '', end_date: '' })
  const sprintDurationWeeks = ref(2)
  const sprintDurationOptions = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12]
  const isCreatingSprint = ref(false)
  const sprintError = ref('')

  function fillCreateSprintFormDefaults() {
    const w = 2
    const start = defaultNextSprintMondayStart(sprints.value, new Date())
    const end = sprintEndFromMondayStart(start, w)
    const pname = project.value?.name ?? ''
    sprintForm.value = {
      name: pname ? defaultSprintDisplayName(pname, nextSprintOrdinal(sprints.value), start, end) : '',
      goal: '',
      start_date: dateToDatetimeLocal(start),
      end_date: dateToDatetimeLocal(end),
    }
    sprintDurationWeeks.value = w
  }

  function refreshCreateSprintNameFromDates(start: Date, end: Date) {
    if (!project.value || editingSprint.value) return
    sprintForm.value.name = defaultSprintDisplayName(project.value.name, nextSprintOrdinal(sprints.value), start, end)
  }

  function syncSprintEndFromStartAndDuration() {
    if (editingSprint.value) return
    const start = parseDatetimeLocal(sprintForm.value.start_date)
    if (!start) return
    const w = Math.max(1, sprintDurationWeeks.value)
    const end = sprintEndFromMondayStart(start, w)
    sprintForm.value.end_date = dateToDatetimeLocal(end)
    refreshCreateSprintNameFromDates(start, end)
  }

  function onSprintDurationWeeksChange(ev: Event) {
    const raw = Number((ev.target as HTMLSelectElement).value)
    sprintDurationWeeks.value = Number.isFinite(raw) && raw >= 1 ? Math.floor(raw) : 1
    syncSprintEndFromStartAndDuration()
  }

  function applySuggestedSprintName() {
    if (!project.value || editingSprint.value) return
    const start = parseDatetimeLocal(sprintForm.value.start_date)
    const end = parseDatetimeLocal(sprintForm.value.end_date)
    sprintForm.value.name = defaultSprintDisplayName(project.value.name, nextSprintOrdinal(sprints.value), start, end)
  }

  function resetSprintModalToDefaults() {
    if (editingSprint.value) return
    fillCreateSprintFormDefaults()
    sprintError.value = ''
  }

  function openSprintModal() {
    editingSprint.value = null
    sprintError.value = ''
    fillCreateSprintFormDefaults()
    showSprintModal.value = true
  }

  function openEditSprintModal(sprint: Sprint) {
    editingSprint.value = sprint
    sprintForm.value = {
      name: sprint.name,
      goal: sprint.goal || '',
      start_date: isoToDatetimeLocal(sprint.start_date),
      end_date: isoToDatetimeLocal(sprint.end_date),
    }
    sprintError.value = ''
    showSprintModal.value = true
  }

  function closeSprintModal() {
    showSprintModal.value = false
    editingSprint.value = null
  }

  async function submitSprint() {
    if (!project.value) {
      sprintError.value = 'Project not loaded. Please refresh the page.'
      return
    }
    isCreatingSprint.value = true
    sprintError.value = ''
    try {
      const name = sprintForm.value.name.trim()
      const goal = sprintForm.value.goal?.trim() || undefined
      let start_date: string | undefined
      let end_date: string | undefined
      const startD = parseDatetimeLocal(sprintForm.value.start_date)
      const endD = parseDatetimeLocal(sprintForm.value.end_date)
      if (startD) start_date = startD.toISOString()
      if (endD) end_date = endD.toISOString()

      if (startD && endD) {
        if (endD.getTime() <= startD.getTime()) {
          sprintError.value = 'End must be after start.'
          return
        }
        if (!editingSprint.value && sprintInclusiveCalendarDays(startD, endD) < 7) {
          sprintError.value = 'Sprint must span at least 1 full week (7 calendar days, Mon–Sun style).'
          return
        }
      }

      if (editingSprint.value) {
        const updated = await projectsApi.updateSprint(editingSprint.value.id, { name, goal, start_date, end_date })
        const idx = sprints.value.findIndex((s) => s.id === editingSprint.value!.id)
        if (idx !== -1) sprints.value[idx] = updated
        closeSprintModal()
      } else {
        const sprint = await projectsApi.createSprint({
          project_id: project.value.id,
          name,
          goal,
          start_date,
          end_date,
        })
        sprints.value.unshift(sprint)
        closeSprintModal()
      }
    } catch (e: any) {
      const msg = e?.data?.message ?? e?.data?.error ?? e?.message ?? (editingSprint.value ? 'Failed to update sprint' : 'Failed to create sprint')
      sprintError.value = typeof msg === 'string' ? msg : 'Failed to save sprint'
    } finally {
      isCreatingSprint.value = false
    }
  }

  const sprintDragId = ref<string | null>(null)
  function onSprintDragStart(e: DragEvent, sprintId: string) {
    sprintDragId.value = sprintId
    e.dataTransfer?.setData?.('application/json', JSON.stringify({ type: 'sprint', id: sprintId }))
    e.dataTransfer!.effectAllowed = 'move'
  }
  function onSprintDragOver(e: DragEvent) {
    e.preventDefault()
    if (e.dataTransfer) e.dataTransfer.dropEffect = 'move'
  }
  function onSprintDrop(e: DragEvent, dropIndex: number) {
    e.preventDefault()
    sprintDragId.value = null
    let dragId: string | null = null
    try {
      const raw = e.dataTransfer?.getData('application/json')
      if (raw) {
        const p = JSON.parse(raw) as { type: string; id: string }
        if (p.type === 'sprint') dragId = p.id
      }
    } catch {}
    if (!dragId) return
    const ordered = [...sprintsOrdered.value]
    const fromIndex = ordered.findIndex((x) => x.id === dragId)
    if (fromIndex < 0 || fromIndex === dropIndex) return
    const [removed] = ordered.splice(fromIndex, 1)
    ordered.splice(dropIndex, 0, removed)
    reorderSprints(ordered)
  }

  async function reorderSprints(newOrder: Sprint[]) {
    try {
      for (let i = 0; i < newOrder.length; i++) {
        const s = newOrder[i]
        if ((s.sort_order ?? 0) === i) continue
        const updated = await projectsApi.updateSprint(s.id, { sort_order: i })
        const idx = sprints.value.findIndex((x) => x.id === s.id)
        if (idx !== -1) (sprints.value[idx] as Sprint).sort_order = updated.sort_order ?? i
      }
    } catch {
      await loadAll()
    }
  }

  const showAddTasksToSprintModal = ref(false)
  const sprintForAddTasks = ref<Sprint | null>(null)
  const selectedTaskIdsForSprint = ref<string[]>([])
  const addTasksToSprintError = ref('')
  const isAddingTasksToSprint = ref(false)

  const tasksNotInSprint = computed(() => {
    const sprintId = sprintForAddTasks.value?.id
    if (!sprintId) return []
    return allTasks.value
      .filter((t) => t.sprint_id !== sprintId)
      .sort((a, b) => (a.code ?? '').localeCompare(b.code ?? '', undefined, { numeric: true }))
  })

  function openAddTasksToSprintModal(sprint: Sprint) {
    sprintForAddTasks.value = sprint
    selectedTaskIdsForSprint.value = []
    addTasksToSprintError.value = ''
    showAddTasksToSprintModal.value = true
  }

  function closeAddTasksToSprintModal() {
    showAddTasksToSprintModal.value = false
    sprintForAddTasks.value = null
    selectedTaskIdsForSprint.value = []
    addTasksToSprintError.value = ''
  }

  async function confirmAddTasksToSprint() {
    if (!sprintForAddTasks.value || selectedTaskIdsForSprint.value.length === 0) return
    const sprint = sprintForAddTasks.value
    isAddingTasksToSprint.value = true
    addTasksToSprintError.value = ''
    try {
      await projectsApi.addTasksToSprint(sprint.id, selectedTaskIdsForSprint.value)
      for (const id of selectedTaskIdsForSprint.value) {
        const t = allTasks.value.find((x) => x.id === id)
        if (t) {
          t.sprint_id = sprint.id
          const dates = taskDatesInSprintRange(t, sprint)
          if (dates) {
            try {
              await tasksApi.updateTask(id, { start_date: dates.start_date, end_date: dates.end_date })
              t.start_date = dates.start_date
              t.end_date = dates.end_date
              updateTaskInTimelineData(id, dates.start_date, dates.end_date)
            } catch {
              // ignore per-task date update failure
            }
          }
        }
      }
      await loadSprintTimeline()
      closeAddTasksToSprintModal()
    } catch (e: any) {
      const err = e?.data?.message ?? e?.data?.error ?? e?.message ?? 'เพิ่มงานไม่สำเร็จ'
      addTasksToSprintError.value = typeof err === 'string' ? err : 'เพิ่มงานไม่สำเร็จ'
    } finally {
      isAddingTasksToSprint.value = false
    }
  }

  const showDeleteSprintModal = ref(false)
  const sprintToDelete = ref<Sprint | null>(null)
  const deleteSprintError = ref('')
  const isDeletingSprint = ref(false)

  function openDeleteSprintModal(sprint: Sprint) {
    sprintToDelete.value = sprint
    deleteSprintError.value = ''
    showDeleteSprintModal.value = true
  }

  function closeDeleteSprintModal() {
    showDeleteSprintModal.value = false
    sprintToDelete.value = null
    deleteSprintError.value = ''
  }

  async function confirmDeleteSprint() {
    if (!sprintToDelete.value) return
    isDeletingSprint.value = true
    deleteSprintError.value = ''
    try {
      await projectsApi.deleteSprint(sprintToDelete.value.id)
      sprints.value = sprints.value.filter((s) => s.id !== sprintToDelete.value!.id)
      closeDeleteSprintModal()
    } catch (e: any) {
      const err = e?.data?.message ?? e?.data?.error ?? e?.message ?? 'ลบไม่สำเร็จ'
      deleteSprintError.value = typeof err === 'string' ? err : 'ลบไม่สำเร็จ'
    } finally {
      isDeletingSprint.value = false
    }
  }

  const showCompleteSprintModal = ref(false)
  const sprintToComplete = ref<Sprint | null>(null)
  const completeSprintError = ref('')
  const isCompletingSprint = ref(false)
  const completeSprintCarryOver = ref<'next_sprint' | 'backlog'>('next_sprint')

  function openCompleteSprintModal(sprint: Sprint) {
    sprintToComplete.value = sprint
    completeSprintCarryOver.value = 'next_sprint'
    completeSprintError.value = ''
    showCompleteSprintModal.value = true
  }

  function closeCompleteSprintModal() {
    showCompleteSprintModal.value = false
    sprintToComplete.value = null
    completeSprintError.value = ''
  }

  async function confirmCompleteSprint() {
    if (!sprintToComplete.value) return
    isCompletingSprint.value = true
    completeSprintError.value = ''
    try {
      const sprint = sprintToComplete.value
      const unfinishedTasks = allTasks.value.filter((t) => t.sprint_id === sprint.id && t.status !== 'COMPLETED')
      const updated = await projectsApi.completeSprint(sprint.id)
      const idx = sprints.value.findIndex((s) => s.id === sprint.id)
      if (idx !== -1) sprints.value[idx] = updated

      if (unfinishedTasks.length > 0) {
        const targetSprintId = completeSprintCarryOver.value === 'next_sprint'
          ? (sprintsOrdered.value.find((s) => s.id !== sprint.id && (s.sort_order ?? 0) > (sprint.sort_order ?? 0))?.id ?? null)
          : null

        for (const task of unfinishedTasks) {
          try {
            await tasksApi.updateTask(task.id, { sprint_id: targetSprintId })
            task.sprint_id = targetSprintId
          } catch (err: any) {
            const msg = err?.data?.message ?? err?.data?.error ?? err?.message ?? 'Failed to move unfinished task'
            throw new Error(typeof msg === 'string' ? msg : 'Failed to move unfinished task')
          }
        }
      }

      closeCompleteSprintModal()
    } catch (e: any) {
      const err = e?.data?.message ?? e?.data?.error ?? e?.message ?? 'ปิด sprint ไม่สำเร็จ'
      completeSprintError.value = typeof err === 'string' ? err : 'ปิด sprint ไม่สำเร็จ'
    } finally {
      isCompletingSprint.value = false
    }
  }

  const showReopenSprintModal = ref(false)
  const sprintToReopen = ref<Sprint | null>(null)
  const reopenSprintError = ref('')
  const isReopeningSprint = ref(false)

  function openReopenSprintModal(sprint: Sprint) {
    sprintToReopen.value = sprint
    reopenSprintError.value = ''
    showReopenSprintModal.value = true
  }

  function closeReopenSprintModal() {
    showReopenSprintModal.value = false
    sprintToReopen.value = null
    reopenSprintError.value = ''
  }

  async function confirmReopenSprint() {
    if (!sprintToReopen.value) return
    isReopeningSprint.value = true
    reopenSprintError.value = ''
    try {
      const updated = await projectsApi.reopenSprint(sprintToReopen.value.id)
      const idx = sprints.value.findIndex((s) => s.id === sprintToReopen.value!.id)
      if (idx !== -1) sprints.value[idx] = updated
      closeReopenSprintModal()
    } catch (e: any) {
      const err = e?.data?.message ?? e?.data?.error ?? e?.message ?? 'เปิด sprint กลับไม่สำเร็จ'
      reopenSprintError.value = typeof err === 'string' ? err : 'เปิด sprint กลับไม่สำเร็จ'
    } finally {
      isReopeningSprint.value = false
    }
  }

  async function handleStartSprint(id: string) {
    if (activeSprint.value) return
    try {
      const updated = await projectsApi.startSprint(id)
      const idx = sprints.value.findIndex((s) => s.id === id)
      if (idx !== -1) sprints.value[idx] = updated
    } catch (e: any) {
      const msg = e?.data?.message ?? e?.data?.error ?? e?.message ?? 'Failed to start sprint'
      showError(typeof msg === 'string' ? msg : 'Failed to start sprint', 'Start sprint failed')
    }
  }

  return {
    sprintsOrdered,
    getSprintStats,
    showSprintModal,
    editingSprint,
    sprintForm,
    sprintDurationWeeks,
    sprintDurationOptions,
    isCreatingSprint,
    sprintError,
    syncSprintEndFromStartAndDuration,
    onSprintDurationWeeksChange,
    applySuggestedSprintName,
    resetSprintModalToDefaults,
    openSprintModal,
    openEditSprintModal,
    closeSprintModal,
    submitSprint,
    sprintDragId,
    onSprintDragStart,
    onSprintDragOver,
    onSprintDrop,
    showAddTasksToSprintModal,
    sprintForAddTasks,
    selectedTaskIdsForSprint,
    addTasksToSprintError,
    isAddingTasksToSprint,
    tasksNotInSprint,
    openAddTasksToSprintModal,
    closeAddTasksToSprintModal,
    confirmAddTasksToSprint,
    showDeleteSprintModal,
    sprintToDelete,
    deleteSprintError,
    isDeletingSprint,
    openDeleteSprintModal,
    closeDeleteSprintModal,
    confirmDeleteSprint,
    showCompleteSprintModal,
    sprintToComplete,
    completeSprintError,
    isCompletingSprint,
    openCompleteSprintModal,
    closeCompleteSprintModal,
    confirmCompleteSprint,
    showReopenSprintModal,
    sprintToReopen,
    reopenSprintError,
    isReopeningSprint,
    openReopenSprintModal,
    closeReopenSprintModal,
    confirmReopenSprint,
    handleStartSprint,
  }
}

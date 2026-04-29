import { computed, ref, type Ref } from 'vue'
import { useAuth } from '~/composables/useAuth'
import { useTeamsApi } from '~/core/modules/teams/infrastructure/teams-api'
import { isTaskAssigneeRole } from '~/utils/roles'

interface TaskLike {
  id: string
  title: string
  code?: string
  parent_id?: string | null
  epic_id?: string | null
}

interface ProjectLike {
  id: string
}

interface TeamsStoreLike {
  teamsFeatureEnabled: boolean
  fetchTeamsFeatureEnabled: () => Promise<void>
}

interface TasksApiLike {
  previewGoogleSlides: (payload: { presentation_url: string }) => Promise<any>
  importGoogleSlides: (payload: any) => Promise<any>
  previewPPTXUpload: (file: File) => Promise<any>
  importPPTXUpload: (file: File, payload: any) => Promise<any>
  previewGoogleSheets: (payload: { sheet_url: string }) => Promise<any>
  importGoogleSheets: (payload: any) => Promise<any>
}

interface BacklogImportAssignee { id: number; email: string; display_name: string; first_name?: string; last_name?: string; role: string }
interface BacklogTriagedSlide { title: string; assignee_id: number | null; estimated_minutes: number; priority: string }
interface SheetsTriagedRow {
  title: string
  priority: string
  estimated_minutes: number
  due_date: string
  status: string
  notes: string
}
interface IODTriagedRow {
  title: string
  priority: string
  estimated_minutes: number
  status: string
  notes: string
  header: string
  header_link: string
  request_method: string
  payload: string
  image_ref: string
  detail_links: string[]
}

export function useProjectImports(params: {
  allTasks: Ref<TaskLike[]>
  project: Ref<ProjectLike | null>
  currentUser: Ref<{ role?: string; user_id?: number } | null>
  tasksApi: TasksApiLike
  teamsStore: TeamsStoreLike
  loadAll: () => Promise<void>
}) {
  const { allTasks, project, currentUser, tasksApi, teamsStore, loadAll } = params

  // Backlog Import from Google Slides
  const showBacklogImportModal = ref(false)
  const backlogImportStep = ref<'form' | 'select' | 'result'>('form')
  const isBacklogImporting = ref(false)
  const isBacklogLoadingPreview = ref(false)
  const backlogImportError = ref('')
  const backlogImportResult = ref<{ created_count: number; slide_count: number; presentation_title: string; tasks: any[] } | null>(null)
  const backlogImportPreview = ref<{
    presentation_title: string
    slides: { index: number; title: string; hidden?: boolean; suggested_task_title?: string }[]
    already_imported_slide_indices?: number[]
  } | null>(null)
  const backlogImportSelectedIndices = ref<number[]>([])
  const backlogImportTriagedSlides = ref<Record<number, BacklogTriagedSlide>>({})
  const backlogImportAssignees = ref<BacklogImportAssignee[]>([])
  const backlogImportForm = ref({
    presentation_url: '',
    priority: 'MEDIUM' as const,
    story_points: 1,
    epic_id: '',
    parent_id: '',
  })

  const backlogParentTaskOptions = computed(() => {
    const epicId = backlogImportForm.value.epic_id
    return allTasks.value
      .filter((t) => {
        if (t.parent_id) return false
        if (epicId) return t.epic_id === epicId
        return true
      })
      .sort((a, b) => a.title.localeCompare(b.title))
  })

  function onBacklogImportEpicChange() {
    const currentParent = backlogParentTaskOptions.value.find((t) => t.id === backlogImportForm.value.parent_id)
    if (!currentParent) backlogImportForm.value.parent_id = ''
  }

  async function loadBacklogImportAssignees() {
    if (backlogImportAssignees.value.length > 0) return
    try {
      const { fetchWithAuth: fw } = useAuth()
      const role = (currentUser.value?.role || '').toUpperCase()
      if (role === 'PRODUCT_OWNER' || role === 'PM') {
        await teamsStore.fetchTeamsFeatureEnabled()
        if (teamsStore.teamsFeatureEnabled) {
          const { getTeams } = useTeamsApi()
          const teams = await getTeams()
          const userId = currentUser.value?.user_id
          const myTeam = teams.find((t: any) => t.users?.some((u: any) => u.id === userId))
          backlogImportAssignees.value = (myTeam?.users ?? [])
            .filter((u: any) => isTaskAssigneeRole(u.role))
            .map((u: any) => ({ id: u.id, email: u.email, display_name: u.display_name, first_name: u.first_name, last_name: u.last_name, role: u.role }))
        } else {
          const res = await fw<{ data: BacklogImportAssignee[] }>('/auth/users')
          backlogImportAssignees.value = (res.data ?? []).filter((u: BacklogImportAssignee) => isTaskAssigneeRole(u.role))
        }
      } else {
        const res = await fw<{ data: BacklogImportAssignee[] }>('/auth/users')
        backlogImportAssignees.value = (res.data ?? []).filter((u: BacklogImportAssignee) => isTaskAssigneeRole(u.role))
      }
    } catch {
      // non-critical
    }
  }

  function openBacklogImportModal() {
    backlogImportForm.value = { presentation_url: '', priority: 'MEDIUM', story_points: 1, epic_id: '', parent_id: '' }
    backlogImportStep.value = 'form'
    backlogImportError.value = ''
    backlogImportResult.value = null
    backlogImportPreview.value = null
    backlogImportSelectedIndices.value = []
    backlogImportTriagedSlides.value = {}
    showBacklogImportModal.value = true
    loadBacklogImportAssignees()
  }

  function closeBacklogImportModal() {
    showBacklogImportModal.value = false
    if (backlogImportResult.value) loadAll()
  }

  async function loadBacklogImportPreview() {
    if (!backlogImportForm.value.presentation_url.trim()) return
    isBacklogLoadingPreview.value = true
    backlogImportError.value = ''
    try {
      const data = await tasksApi.previewGoogleSlides({
        presentation_url: backlogImportForm.value.presentation_url.trim(),
      })
      backlogImportPreview.value = data
      const alreadySet = new Set(data.already_imported_slide_indices ?? [])
      backlogImportSelectedIndices.value = data.slides
        .filter((s: { index: number; hidden?: boolean }) => !s.hidden && !alreadySet.has(s.index))
        .map((s: { index: number }) => s.index)

      const triagedMap: Record<number, BacklogTriagedSlide> = {}
      for (const s of data.slides) {
        const st = s.suggested_task_title?.trim()
        triagedMap[s.index] = {
          title: st || `Slide ${s.index}`,
          assignee_id: null,
          estimated_minutes: 0,
          priority: backlogImportForm.value.priority || 'MEDIUM',
        }
      }
      backlogImportTriagedSlides.value = triagedMap
      backlogImportStep.value = 'select'
    } catch (e: any) {
      backlogImportError.value = e?.data?.message ?? e?.message ?? 'โหลดรายการไม่สำเร็จ'
    } finally {
      isBacklogLoadingPreview.value = false
    }
  }

  function backlogImportSelectAll() {
    if (backlogImportPreview.value) backlogImportSelectedIndices.value = backlogImportPreview.value.slides.map((s) => s.index)
  }

  function backlogImportDeselectAll() {
    backlogImportSelectedIndices.value = []
  }

  function backlogImportSelectOnlyNew() {
    if (!backlogImportPreview.value) return
    const alreadySet = new Set(backlogImportPreview.value.already_imported_slide_indices ?? [])
    backlogImportSelectedIndices.value = backlogImportPreview.value.slides
      .filter((s: { index: number; hidden?: boolean }) => !s.hidden && !alreadySet.has(s.index))
      .map((s: { index: number }) => s.index)
  }

  async function submitBacklogImport() {
    if (!project.value) return
    isBacklogImporting.value = true
    backlogImportError.value = ''
    try {
      const triageSlides = backlogImportSelectedIndices.value.map((idx) => {
        const t = backlogImportTriagedSlides.value[idx]
        return {
          slide_index: idx,
          title: t?.title || `Slide ${idx}`,
          assignee_id: t?.assignee_id ?? null,
          estimated_minutes: t?.estimated_minutes ?? 0,
          priority: t?.priority || 'MEDIUM',
        }
      })

      const payload: any = {
        presentation_url: backlogImportForm.value.presentation_url.trim(),
        project_id: project.value.id,
        slides: triageSlides,
      }
      if (backlogImportForm.value.epic_id) payload.epic_id = backlogImportForm.value.epic_id
      if (backlogImportForm.value.parent_id) payload.parent_id = backlogImportForm.value.parent_id
      backlogImportResult.value = await tasksApi.importGoogleSlides(payload)
      backlogImportStep.value = 'result'
    } catch (e: any) {
      backlogImportError.value = e?.data?.message ?? e?.message ?? 'Import failed'
    } finally {
      isBacklogImporting.value = false
    }
  }

  // PPTX File Upload Import
  const showPPTXImportModal = ref(false)
  const pptxImportStep = ref<'form' | 'select' | 'result'>('form')
  const isPPTXImporting = ref(false)
  const isPPTXLoadingPreview = ref(false)
  const pptxImportError = ref('')
  const pptxImportResult = ref<{ created_count: number; page_count: number; title: string; tasks: any[] } | null>(null)
  const pptxImportPreview = ref<{ title: string; slides: { index: number; title: string; hidden?: boolean; suggested_task_title?: string }[] } | null>(null)
  const pptxImportFile = ref<File | null>(null)
  const pptxDragOver = ref(false)
  const pptxImportSelectedIndices = ref<number[]>([])
  const pptxImportTriagedSlides = ref<Record<number, BacklogTriagedSlide>>({})
  const pptxImportForm = ref({ epic_id: '', parent_id: '' })

  const pptxParentTaskOptions = computed(() => {
    const epicId = pptxImportForm.value.epic_id
    return allTasks.value
      .filter((t) => {
        if (t.parent_id) return false
        if (epicId) return t.epic_id === epicId
        return true
      })
      .sort((a, b) => a.title.localeCompare(b.title))
  })

  function onPPTXImportEpicChange() {
    const currentParent = pptxParentTaskOptions.value.find((t) => t.id === pptxImportForm.value.parent_id)
    if (!currentParent) pptxImportForm.value.parent_id = ''
  }

  function onPPTXFileChange(e: Event) {
    const input = e.target as HTMLInputElement
    if (input.files && input.files[0]) {
      pptxImportFile.value = input.files[0]
      pptxImportError.value = ''
    }
  }

  function onPPTXFileDrop(e: DragEvent) {
    pptxDragOver.value = false
    const file = e.dataTransfer?.files[0]
    if (file && (file.name.endsWith('.pptx') || file.type.includes('presentationml'))) {
      pptxImportFile.value = file
      pptxImportError.value = ''
    } else {
      pptxImportError.value = 'กรุณาเลือกไฟล์ .pptx เท่านั้น'
    }
  }

  function openPPTXImportModal() {
    pptxImportForm.value = { epic_id: '', parent_id: '' }
    pptxImportStep.value = 'form'
    pptxImportError.value = ''
    pptxImportResult.value = null
    pptxImportPreview.value = null
    pptxImportFile.value = null
    pptxImportSelectedIndices.value = []
    pptxImportTriagedSlides.value = {}
    showPPTXImportModal.value = true
    loadBacklogImportAssignees()
  }

  function closePPTXImportModal() {
    showPPTXImportModal.value = false
    if (pptxImportResult.value) loadAll()
  }

  function pptxImportSelectAll() {
    if (pptxImportPreview.value) pptxImportSelectedIndices.value = pptxImportPreview.value.slides.map((s) => s.index)
  }

  function pptxImportDeselectAll() {
    pptxImportSelectedIndices.value = []
  }

  function pptxImportSelectOnlyVisible() {
    if (pptxImportPreview.value) {
      pptxImportSelectedIndices.value = pptxImportPreview.value.slides.filter((s) => !s.hidden).map((s) => s.index)
    }
  }

  async function loadPPTXImportPreview() {
    if (!pptxImportFile.value) return
    isPPTXLoadingPreview.value = true
    pptxImportError.value = ''
    try {
      const data = await tasksApi.previewPPTXUpload(pptxImportFile.value)
      pptxImportPreview.value = data
      pptxImportSelectedIndices.value = data.slides.filter((s: any) => !s.hidden).map((s: any) => s.index)
      const triagedMap: Record<number, BacklogTriagedSlide> = {}
      for (const s of data.slides) {
        triagedMap[s.index] = {
          title: s.suggested_task_title?.trim() || `Slide ${s.index}`,
          assignee_id: null,
          estimated_minutes: 0,
          priority: 'MEDIUM',
        }
      }
      pptxImportTriagedSlides.value = triagedMap
      pptxImportStep.value = 'select'
    } catch (e: any) {
      pptxImportError.value = e?.data?.message ?? e?.message ?? 'โหลดรายการไม่สำเร็จ'
    } finally {
      isPPTXLoadingPreview.value = false
    }
  }

  async function submitPPTXImport() {
    if (!project.value || !pptxImportFile.value) return
    isPPTXImporting.value = true
    pptxImportError.value = ''
    try {
      const pages = pptxImportSelectedIndices.value.map((idx) => {
        const t = pptxImportTriagedSlides.value[idx]
        return {
          slide_index: idx,
          title: t?.title || `Slide ${idx}`,
          assignee_id: t?.assignee_id ?? null,
          estimated_minutes: t?.estimated_minutes ?? 0,
          priority: t?.priority || 'MEDIUM',
        }
      })
      const payload: Record<string, unknown> = {
        project_id: project.value.id,
        priority: 'MEDIUM',
        story_points: 1,
        pages,
      }
      if (pptxImportForm.value.epic_id) payload.epic_id = pptxImportForm.value.epic_id
      if (pptxImportForm.value.parent_id) payload.parent_id = pptxImportForm.value.parent_id
      pptxImportResult.value = await tasksApi.importPPTXUpload(pptxImportFile.value, payload as any)
      pptxImportStep.value = 'result'
    } catch (e: any) {
      pptxImportError.value = e?.data?.message ?? e?.message ?? 'Import failed'
    } finally {
      isPPTXImporting.value = false
    }
  }

  // Sheets Import
  const showSheetsImportModal = ref(false)
  const sheetsImportStep = ref<'form' | 'select' | 'result'>('form')
  const isSheetsImporting = ref(false)
  const isSheetsLoadingPreview = ref(false)
  const sheetsImportError = ref('')
  const sheetsImportResult = ref<{ created_count: number; sheet_title: string; tasks: any[] } | null>(null)
  const sheetsImportPreview = ref<{
    sheet_title: string
    sheet_id: string
    rows: { row_index: number; title: string; due_date: string; status: string; raw_status: string; notes: string }[]
  } | null>(null)
  const sheetsImportSelectedRowIndices = ref<number[]>([])
  const sheetsImportTriagedRows = ref<Record<number, SheetsTriagedRow>>({})
  const sheetsImportForm = ref({ sheet_url: '', epic_id: '', parent_id: '' })

  const sheetsParentTaskOptions = computed(() => {
    const epicId = sheetsImportForm.value.epic_id
    return allTasks.value
      .filter((t) => {
        if (t.parent_id) return false
        if (epicId) return t.epic_id === epicId
        return true
      })
      .sort((a, b) => a.title.localeCompare(b.title))
  })

  function onSheetsImportEpicChange() {
    const currentParent = sheetsParentTaskOptions.value.find((t) => t.id === sheetsImportForm.value.parent_id)
    if (!currentParent) sheetsImportForm.value.parent_id = ''
  }

  function openSheetsImportModal() {
    sheetsImportForm.value = { sheet_url: '', epic_id: '', parent_id: '' }
    sheetsImportStep.value = 'form'
    sheetsImportError.value = ''
    sheetsImportResult.value = null
    sheetsImportPreview.value = null
    sheetsImportSelectedRowIndices.value = []
    sheetsImportTriagedRows.value = {}
    showSheetsImportModal.value = true
  }

  function closeSheetsImportModal() {
    showSheetsImportModal.value = false
    if (sheetsImportResult.value) loadAll()
  }

  function sheetsImportSelectAll() {
    if (sheetsImportPreview.value) {
      sheetsImportSelectedRowIndices.value = sheetsImportPreview.value.rows.map((r) => r.row_index)
    }
  }

  function sheetsImportDeselectAll() {
    sheetsImportSelectedRowIndices.value = []
  }

  async function loadSheetsImportPreview() {
    if (!sheetsImportForm.value.sheet_url.trim()) return
    isSheetsLoadingPreview.value = true
    sheetsImportError.value = ''
    try {
      const data = await tasksApi.previewGoogleSheets({ sheet_url: sheetsImportForm.value.sheet_url.trim() })
      sheetsImportPreview.value = data
      const triagedMap: Record<number, SheetsTriagedRow> = {}
      const selected: number[] = []
      for (const r of data.rows) {
        selected.push(r.row_index)
        triagedMap[r.row_index] = {
          title: r.title,
          priority: 'MEDIUM',
          estimated_minutes: 0,
          due_date: r.due_date || '',
          status: r.status,
          notes: r.notes,
        }
      }
      sheetsImportTriagedRows.value = triagedMap
      sheetsImportSelectedRowIndices.value = selected
      sheetsImportStep.value = 'select'
    } catch (e: any) {
      sheetsImportError.value = e?.data?.message ?? e?.message ?? 'โหลดรายการไม่สำเร็จ'
    } finally {
      isSheetsLoadingPreview.value = false
    }
  }

  async function submitSheetsImport() {
    if (!project.value || !sheetsImportPreview.value) return
    sheetsImportError.value = ''
    isSheetsImporting.value = true
    try {
      const rows = sheetsImportSelectedRowIndices.value.map((rowIndex) => {
        const t = sheetsImportTriagedRows.value[rowIndex]
        const rawEst = Number(t?.estimated_minutes)
        const estimatedMinutes = Number.isFinite(rawEst) && rawEst >= 0 ? Math.floor(rawEst) : 0
        return {
          row_index: rowIndex,
          title: t?.title?.trim() || '',
          priority: t?.priority || 'MEDIUM',
          estimated_minutes: estimatedMinutes,
          due_date: t?.due_date?.trim() || '',
          status: t?.status || 'PENDING',
          notes: t?.notes?.trim() || '',
        }
      })
      const payload: Record<string, unknown> = {
        sheet_url: sheetsImportForm.value.sheet_url.trim(),
        sheet_title: sheetsImportPreview.value.sheet_title,
        project_id: project.value.id,
        rows,
      }
      if (sheetsImportForm.value.epic_id) payload.epic_id = sheetsImportForm.value.epic_id
      if (sheetsImportForm.value.parent_id) payload.parent_id = sheetsImportForm.value.parent_id
      sheetsImportResult.value = await tasksApi.importGoogleSheets(payload as any)
      sheetsImportStep.value = 'result'
    } catch (e: any) {
      sheetsImportError.value = e?.data?.message ?? e?.message ?? 'Import failed'
    } finally {
      isSheetsImporting.value = false
    }
  }

  // IOD Import
  const showIODImportModal = ref(false)
  const iodImportStep = ref<'form' | 'select' | 'result'>('form')
  const isIODImporting = ref(false)
  const isIODLoadingPreview = ref(false)
  const iodImportError = ref('')
  const iodImportResult = ref<{ created_count: number; sheet_title: string; tasks: any[] } | null>(null)
  const iodImportPreview = ref<{
    sheet_title: string
    sheet_id: string
    rows: {
      row_index: number
      title: string
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
  } | null>(null)
  const iodImportSelectedRowIndices = ref<number[]>([])
  const iodImportTriagedRows = ref<Record<number, IODTriagedRow>>({})
  const iodImportForm = ref({ sheet_url: '', epic_id: '', parent_id: '' })

  const iodParentTaskOptions = computed(() => {
    const epicId = iodImportForm.value.epic_id
    return allTasks.value
      .filter((t) => {
        if (t.parent_id) return false
        if (epicId) return t.epic_id === epicId
        return true
      })
      .sort((a, b) => a.title.localeCompare(b.title))
  })

  function onIODImportEpicChange() {
    const cur = iodParentTaskOptions.value.find((t) => t.id === iodImportForm.value.parent_id)
    if (!cur) iodImportForm.value.parent_id = ''
  }

  function openIODImportModal() {
    iodImportForm.value = { sheet_url: '', epic_id: '', parent_id: '' }
    iodImportStep.value = 'form'
    iodImportError.value = ''
    iodImportResult.value = null
    iodImportPreview.value = null
    iodImportSelectedRowIndices.value = []
    iodImportTriagedRows.value = {}
    showIODImportModal.value = true
  }

  function closeIODImportModal() {
    showIODImportModal.value = false
    if (iodImportResult.value) loadAll()
  }

  async function loadIODImportPreview() {
    if (!iodImportForm.value.sheet_url.trim()) return
    isIODLoadingPreview.value = true
    iodImportError.value = ''
    try {
      const data = await tasksApi.previewGoogleSheets({ sheet_url: iodImportForm.value.sheet_url.trim() })
      iodImportPreview.value = data as any
      const triagedMap: Record<number, IODTriagedRow> = {}
      const selected: number[] = []
      for (const r of data.rows) {
        selected.push(r.row_index)
        triagedMap[r.row_index] = {
          title: r.title,
          priority: 'MEDIUM',
          estimated_minutes: 0,
          status: r.status,
          notes: r.notes,
          header: r.header || '',
          header_link: r.header_link || '',
          request_method: r.request_method || '',
          payload: r.payload || '',
          image_ref: r.image_ref || '',
          detail_links: r.detail_links || [],
        }
      }
      iodImportTriagedRows.value = triagedMap
      iodImportSelectedRowIndices.value = selected
      iodImportStep.value = 'select'
    } catch (e: any) {
      iodImportError.value = e?.data?.message ?? e?.message ?? 'โหลดรายการไม่สำเร็จ'
    } finally {
      isIODLoadingPreview.value = false
    }
  }

  async function submitIODImport() {
    if (!project.value || !iodImportPreview.value) return
    iodImportError.value = ''
    isIODImporting.value = true
    try {
      const rows = iodImportSelectedRowIndices.value.map((rowIndex) => {
        const t = iodImportTriagedRows.value[rowIndex]
        const rawEst = Number(t?.estimated_minutes)
        const estimatedMinutes = Number.isFinite(rawEst) && rawEst >= 0 ? Math.floor(rawEst) : 0
        return {
          row_index: rowIndex,
          title: t?.title?.trim() || '',
          priority: t?.priority || 'MEDIUM',
          estimated_minutes: estimatedMinutes,
          due_date: '',
          status: t?.status || 'PENDING',
          notes: t?.notes?.trim() || '',
          header: t?.header || '',
          header_link: t?.header_link || '',
          request_method: t?.request_method || '',
          payload: t?.payload || '',
          image_ref: t?.image_ref || '',
          detail_links: t?.detail_links || [],
        }
      })
      const payload: Record<string, unknown> = {
        sheet_url: iodImportForm.value.sheet_url.trim(),
        sheet_title: iodImportPreview.value.sheet_title,
        project_id: project.value.id,
        rows,
      }
      if (iodImportForm.value.epic_id) payload.epic_id = iodImportForm.value.epic_id
      if (iodImportForm.value.parent_id) payload.parent_id = iodImportForm.value.parent_id
      iodImportResult.value = await tasksApi.importGoogleSheets(payload as any)
      iodImportStep.value = 'result'
    } catch (e: any) {
      iodImportError.value = e?.data?.message ?? e?.message ?? 'Import failed'
    } finally {
      isIODImporting.value = false
    }
  }

  return {
    showBacklogImportModal,
    backlogImportStep,
    isBacklogImporting,
    isBacklogLoadingPreview,
    backlogImportError,
    backlogImportResult,
    backlogImportPreview,
    backlogImportSelectedIndices,
    backlogImportTriagedSlides,
    backlogImportAssignees,
    backlogImportForm,
    backlogParentTaskOptions,
    onBacklogImportEpicChange,
    openBacklogImportModal,
    closeBacklogImportModal,
    loadBacklogImportPreview,
    backlogImportSelectAll,
    backlogImportDeselectAll,
    backlogImportSelectOnlyNew,
    submitBacklogImport,

    showPPTXImportModal,
    pptxImportStep,
    isPPTXImporting,
    isPPTXLoadingPreview,
    pptxImportError,
    pptxImportResult,
    pptxImportPreview,
    pptxImportFile,
    pptxDragOver,
    pptxImportSelectedIndices,
    pptxImportTriagedSlides,
    pptxImportForm,
    pptxParentTaskOptions,
    onPPTXImportEpicChange,
    onPPTXFileChange,
    onPPTXFileDrop,
    openPPTXImportModal,
    closePPTXImportModal,
    pptxImportSelectAll,
    pptxImportDeselectAll,
    pptxImportSelectOnlyVisible,
    loadPPTXImportPreview,
    submitPPTXImport,

    showSheetsImportModal,
    sheetsImportStep,
    isSheetsImporting,
    isSheetsLoadingPreview,
    sheetsImportError,
    sheetsImportResult,
    sheetsImportPreview,
    sheetsImportSelectedRowIndices,
    sheetsImportTriagedRows,
    sheetsImportForm,
    sheetsParentTaskOptions,
    onSheetsImportEpicChange,
    openSheetsImportModal,
    closeSheetsImportModal,
    sheetsImportSelectAll,
    sheetsImportDeselectAll,
    loadSheetsImportPreview,
    submitSheetsImport,

    showIODImportModal,
    iodImportStep,
    isIODImporting,
    isIODLoadingPreview,
    iodImportError,
    iodImportResult,
    iodImportPreview,
    iodImportSelectedRowIndices,
    iodImportTriagedRows,
    iodImportForm,
    iodParentTaskOptions,
    onIODImportEpicChange,
    openIODImportModal,
    closeIODImportModal,
    loadIODImportPreview,
    submitIODImport,
  }
}

export interface BacklogTaskLike {
  id: string
  code?: string
  parent_id?: string | null
  sprint_id?: string | null
  sort_order?: number | null
  created_at: string
}

export function taskCodeSuffix(code: string | undefined): string {
  if (!code) return '–'
  const suffix = code.split('-').pop()
  return /^\d+$/.test(suffix || '') ? String(Number(suffix)).padStart(4, '0') : code
}

function codeSuffix3(code: string | undefined): string {
  if (!code) return '???'
  const suffix = code.split('-').pop()
  return suffix && /^\d+$/.test(suffix) ? String(Number(suffix)).padStart(3, '0') : (suffix || '???')
}

export function buildTaskDisplayCodeMap<T extends BacklogTaskLike>(tasksInDisplayOrder: T[], allTasks: T[]): Record<string, string> {
  const byId = Object.fromEntries(allTasks.map((t) => [t.id, t])) as Record<string, T | undefined>
  const out: Record<string, string> = {}

  for (const t of tasksInDisplayOrder) {
    const num = codeSuffix3(t.code)
    if (!t.parent_id) {
      out[t.id] = 'A' + num
      continue
    }

    const parent = byId[t.parent_id]
    out[t.id] = parent && parent.parent_id ? 'C' + num : 'B' + num
  }

  return out
}

export function backlogSprintOrderIndex(task: Pick<BacklogTaskLike, 'sprint_id'>, sprintOrderIds: string[]): number {
  if (!task.sprint_id) return 0
  const idx = sprintOrderIds.findIndex((id) => id === task.sprint_id)
  return idx === -1 ? 9999 : idx + 1
}

export function sortBacklogTasks<T extends BacklogTaskLike>(tasks: T[], sprintOrderIds: string[]): T[] {
  void sprintOrderIds
  return [...tasks].sort(
    (a, b) =>
      (a.sort_order ?? 0) - (b.sort_order ?? 0) ||
      new Date(a.created_at).getTime() - new Date(b.created_at).getTime(),
  )
}

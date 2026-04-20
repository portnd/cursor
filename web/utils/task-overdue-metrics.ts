/**
 * Deploy / UAT pipeline: past due_at should not count as dev "overdue"
 * (aligned with Kanban columns from Wait for Deploy onward).
 */
function workflowExcludesOverdueMetrics(task: { status: string; task_type?: string | null }): boolean {
  if (task.status === 'WAIT_FOR_DEPLOY' || task.status === 'READY_FOR_UAT') return true
  if (task.status === 'REVIEW_PENDING' && (task.task_type || 'TASK') === 'FEATURE') return true
  return false
}

export function isTaskOverdueForMetrics(task: { status: string; due_at?: string | null; task_type?: string | null }): boolean {
  if (!task.due_at || task.status === 'COMPLETED') return false
  if (workflowExcludesOverdueMetrics(task)) return false
  return new Date(task.due_at).getTime() < Date.now()
}

/** When past due but not counted as overdue — short label for Kanban cards. */
export function pastDueNeutralCaption(task: { status: string; task_type?: string | null }): string | null {
  if (!workflowExcludesOverdueMetrics(task)) return null
  if (task.status === 'WAIT_FOR_DEPLOY') return 'Awaiting deploy'
  if (task.status === 'READY_FOR_UAT') return 'In UAT'
  return 'In review (UAT path)'
}

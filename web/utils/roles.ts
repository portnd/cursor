/** Roles allowed as task assignees (matches API `users.role`). */
export const TASK_ASSIGNEE_ROLES: readonly string[] = [
  'ENGINEER',
  'CHIEF_ENGINEER',
  'DEV',
  'PRODUCT_OWNER',
  'PM',
  'MANAGER',
  'SUPPORT',
]

export function isEngineerLikeRole(role: string | null | undefined): boolean {
  const r = (role ?? '').toUpperCase()
  return r === 'ENGINEER' || r === 'CHIEF_ENGINEER' || r === 'DEV'
}

export function isTaskAssigneeRole(role: string | null | undefined): boolean {
  if (!role) return false
  return TASK_ASSIGNEE_ROLES.includes(role)
}

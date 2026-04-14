/** Roles allowed as task assignees (matches API `users.role`). */
export const TASK_ASSIGNEE_ROLES: readonly string[] = [
  'ENGINEER',
  'CHIEF_ENGINEER',
  'DEV',
  'PRODUCT_OWNER',
  'PM',
  'MANAGER',
  'SUPPORT',
  'CEO',
]

export function isEngineerLikeRole(role: string | null | undefined): boolean {
  const r = (role ?? '').toUpperCase()
  return r === 'ENGINEER' || r === 'CHIEF_ENGINEER' || r === 'DEV'
}

export function isTaskAssigneeRole(role: string | null | undefined): boolean {
  if (!role) return false
  return TASK_ASSIGNEE_ROLES.includes(role)
}

export function canSeeCeoAssigneeOption(role: string | null | undefined): boolean {
  return (role ?? '').toUpperCase() === 'CEO'
}

export function canAssignCeoAsTaskAssignee(viewerRole: string | null | undefined): boolean {
  return canSeeCeoAssigneeOption(viewerRole)
}

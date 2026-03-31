/** Snap to 0.1 h before converting to minutes (0.1 h = 6 min). */
function roundHoursToTenth(hours: number): number {
  return Math.round(hours * 10) / 10
}

/** Input binding: hours from API minutes, shown with 1 decimal (e.g. 90 min → 1.5). */
export function minutesToEffortHours(minutes: number): number {
  const m = Number(minutes) || 0
  return Number((m / 60).toFixed(1))
}

/** Persist to API: hours snapped to 0.1 h, then to integer minutes (non-negative). */
export function effortHoursToMinutes(hours: number): number {
  const h = Number(hours)
  if (!Number.isFinite(h) || h < 0) return 0
  const tenths = roundHoursToTenth(h)
  return Math.round(tenths * 60)
}

/** Read-only label from stored minutes (at most 1 decimal place). */
export function formatMinutesAsHours(minutes: number): string {
  const m = Number(minutes) || 0
  if (m === 0) return '0'
  const s = (m / 60).toFixed(1)
  if (s === '0.0') return '0'
  return s.endsWith('.0') ? s.slice(0, -2) : s
}

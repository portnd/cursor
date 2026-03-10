import { useAuth } from '~/composables/useAuth'

// ─── Request / Response types ─────────────────────────────────────────────────

export interface SubmitStandupRequest {
  date: string              // YYYY-MM-DD
  yesterday_summary: string
  today_task_ids: string[]
  blocker: string
}

export interface DailyStandup {
  id: string
  user_id: number
  date: string
  yesterday_summary: string
  today_task_ids: string[]
  blocker: string
  created_at: string
  user_email?: string
  user_display_name?: string
}

export interface ActivityItem {
  type: 'time_log' | 'submission'
  description: string
  minutes?: number
  ai_verdict?: string
  ai_score?: number
  occurred_at: string
}

export interface UserPulse {
  user_id: number
  user_email: string
  user_display_name: string
  standup: DailyStandup | null
  total_logged_minutes: number
  total_logged_hours: number
  latest_activities: ActivityItem[]
  has_blocker: boolean
}

export interface CompanyPulseResponse {
  date: string
  total_members: number
  checked_in: number
  total_minutes_logged: number
  members: UserPulse[]
}

// ─── API composable ───────────────────────────────────────────────────────────

function usePulseApi() {
  const { fetchWithAuth } = useAuth()

  async function submitStandup(payload: SubmitStandupRequest): Promise<DailyStandup> {
    return await fetchWithAuth<DailyStandup>('/pulse/standup', {
      method: 'POST',
      body: payload,
    })
  }

  async function getDailyPulse(date?: string): Promise<CompanyPulseResponse> {
    const query = date ? `?date=${date}` : ''
    return await fetchWithAuth<CompanyPulseResponse>(`/pulse/daily${query}`)
  }

  return { submitStandup, getDailyPulse }
}

export { usePulseApi }

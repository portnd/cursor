import { useAuth } from '~/composables/useAuth'

export interface B2BRequest {
  id: string
  title: string
  description: string
  estimated_minutes: number
  proposed_minutes: number
  negotiation_reason: string
  status: 'PENDING' | 'COUNTER_OFFERED' | 'ACCEPTED' | 'REJECTED'
  requester_team_id: number
  target_team_id: number
  requester_user_id: number
  created_task_id?: string | null
  requester_team_name?: string
  target_team_name?: string
  created_at: string
  updated_at: string
}

export interface CreateB2BRequestPayload {
  title: string
  description?: string
  estimated_minutes: number
  target_team_id: number
}

export interface CounterOfferPayload {
  action: 'COUNTER'
  proposed_minutes: number
  reason?: string
}

export interface RejectPayload {
  action: 'REJECT'
}

export function useB2BApi() {
  const { fetchWithAuth } = useAuth()

  const createRequest = (payload: CreateB2BRequestPayload) =>
    fetchWithAuth<{ data: B2BRequest; message: string }>('/sentinel/b2b/requests', {
      method: 'POST',
      body: payload,
    })

  const getRequests = (direction: 'inbound' | 'outbound') =>
    fetchWithAuth<{ data: B2BRequest[] }>(`/sentinel/b2b/requests?direction=${direction}`)

  const counterOffer = (id: string, proposedMinutes: number, reason: string) =>
    fetchWithAuth<{ data: B2BRequest }>(`/sentinel/b2b/requests/${id}`, {
      method: 'PATCH',
      body: { action: 'COUNTER', proposed_minutes: proposedMinutes, reason } satisfies CounterOfferPayload,
    })

  const rejectRequest = (id: string) =>
    fetchWithAuth<{ data: B2BRequest }>(`/sentinel/b2b/requests/${id}`, {
      method: 'PATCH',
      body: { action: 'REJECT' } satisfies RejectPayload,
    })

  const acceptRequest = (id: string) =>
    fetchWithAuth<{ data: { id: string; title: string }; message: string }>(`/sentinel/b2b/requests/${id}/accept`, {
      method: 'POST',
    })

  return { createRequest, getRequests, counterOffer, rejectRequest, acceptRequest }
}

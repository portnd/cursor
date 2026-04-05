import { useAuth } from '~/composables/useAuth'

export interface DeploymentRequest {
  id: number
  title: string
  description: string
  branch: string
  pr_url: string
  environment: 'STAGING' | 'PRE-PROD' | 'PRODUCTION'
  status: 'PENDING' | 'REVIEWING' | 'APPROVED' | 'REJECTED' | 'DEPLOYED'
  requester_id: number
  reviewer_id?: number
  task_id?: string
  task_ref?: string
  rejection_reason?: string
  review_notes?: string
  deployed_at?: string
  created_at: string
  updated_at: string
  requester_email?: string
  requester_display_name?: string
  reviewer_email?: string
  reviewer_display_name?: string
}

export interface DeploymentStats {
  total_pending: number
  total_reviewing: number
  total_approved: number
  total_deployed: number
  total_rejected: number
  deployed_today: number
}

export interface CreateDeploymentPayload {
  title: string
  description?: string
  branch: string
  pr_url?: string
  environment: 'STAGING' | 'PRE-PROD' | 'PRODUCTION'
  task_id?: string    // UUID of linked sentinel task
  task_ref?: string
  reviewer_id?: number // optional pre-assign to a specific Chief Engineer
}

export interface ReviewActionPayload {
  notes?: string
  reason?: string
}

export interface DeploymentUser {
  id: number
  email: string
  display_name: string
  role: string
}

export function useDeploymentApi() {
  const { fetchWithAuth } = useAuth()

  async function createRequest(payload: CreateDeploymentPayload): Promise<DeploymentRequest> {
    return fetchWithAuth<DeploymentRequest>('/deployment/requests', {
      method: 'POST',
      body: JSON.stringify(payload),
    })
  }

  async function listRequests(status?: string): Promise<DeploymentRequest[]> {
    const qs = status ? `?status=${status}` : ''
    const data = await fetchWithAuth<{ data: DeploymentRequest[] }>(`/deployment/requests${qs}`)
    return data.data ?? []
  }

  async function getRequest(id: number): Promise<DeploymentRequest> {
    return fetchWithAuth<DeploymentRequest>(`/deployment/requests/${id}`)
  }

  async function pickForReview(id: number): Promise<DeploymentRequest> {
    return fetchWithAuth<DeploymentRequest>(`/deployment/requests/${id}/pick`, { method: 'PATCH' })
  }

  async function approveRequest(id: number, payload?: ReviewActionPayload): Promise<DeploymentRequest> {
    return fetchWithAuth<DeploymentRequest>(`/deployment/requests/${id}/approve`, {
      method: 'PATCH',
      body: JSON.stringify(payload ?? {}),
    })
  }

  async function rejectRequest(id: number, payload: ReviewActionPayload): Promise<DeploymentRequest> {
    return fetchWithAuth<DeploymentRequest>(`/deployment/requests/${id}/reject`, {
      method: 'PATCH',
      body: JSON.stringify(payload),
    })
  }

  async function markDeployed(id: number, payload?: ReviewActionPayload): Promise<DeploymentRequest> {
    return fetchWithAuth<DeploymentRequest>(`/deployment/requests/${id}/deploy`, {
      method: 'PATCH',
      body: JSON.stringify(payload ?? {}),
    })
  }

  async function getByTaskId(taskId: string): Promise<DeploymentRequest | null> {
    try {
      return await fetchWithAuth<DeploymentRequest>(`/deployment/requests/by-task/${taskId}`)
    } catch {
      return null
    }
  }

  async function getStats(): Promise<DeploymentStats> {
    return fetchWithAuth<DeploymentStats>('/deployment/stats')
  }

  /** Fetch all users with CHIEF_ENGINEER role for assignee selector */
  async function fetchChiefEngineers(): Promise<DeploymentUser[]> {
    try {
      const res = await fetchWithAuth<{ data: DeploymentUser[] }>('/auth/users')
      return (res.data ?? []).filter(
        (u) => u.role?.toUpperCase() === 'CHIEF_ENGINEER'
      )
    } catch {
      return []
    }
  }

  return {
    createRequest,
    listRequests,
    getRequest,
    getByTaskId,
    pickForReview,
    approveRequest,
    rejectRequest,
    markDeployed,
    getStats,
    fetchChiefEngineers,
  }
}

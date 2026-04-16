/**
 * 🔐 Authentication Composable
 * God-Tier JWT Auth with Cookie Management
 */

interface LoginResponse {
  data: {
    token: string
    user: {
      id: number
      email: string
      role: string
    }
  }
  message: string
}

export const useAuth = () => {
  const config = useRuntimeConfig()
  const token = useCookie('token', {
    maxAge: 60 * 60 * 24 * 7, // 7 days
    sameSite: 'lax',
    secure: process.env.NODE_ENV === 'production'
  })

  /**
   * Base URL for API (from nuxt.config / NUXT_PUBLIC_API_BASE).
   * On SSR (e.g. Docker web container), use NUXT_PUBLIC_API_BASE_SERVER (e.g. http://api:8080/api/v1) so the server can reach the API container; otherwise browser would get ERR_EMPTY_RESPONSE when SSR fetches fail to localhost:8080.
   */
  const apiBase = computed(() => {
    const isServer = import.meta.server
    const serverBase = (config.public?.apiBaseServer as string) || ''
    if (isServer && serverBase) {
      return serverBase.replace(/\/$/, '')
    }
    const base = (config.public?.apiBase as string) || ''
    const trimmed = base.replace(/\/$/, '')
    return trimmed || '/api/v1'
  })

  const LOGIN_TIMEOUT_MS = 8000

  /**
   * Login user and save JWT token. Uses timeout so login does not spin forever if API is down.
   */
  const login = async (email: string, password: string) => {
    const url = `${apiBase.value}/auth/login`
    const controller = new AbortController()
    const timeoutId = setTimeout(() => controller.abort(), LOGIN_TIMEOUT_MS)
    try {
      const response = await $fetch<LoginResponse>(url, {
        method: 'POST',
        body: { email, password },
        signal: controller.signal
      })
      clearTimeout(timeoutId)
      if (response.data?.token) {
        token.value = response.data.token
        return { success: true, token: response.data.token, user: response.data.user }
      }
      return { success: false, error: 'No token received' }
    } catch (error: any) {
      clearTimeout(timeoutId)
      console.error('Login failed:', error)
      const msg = error?.message || ''
      if (error?.name === 'AbortError' || msg.includes('abort') || msg.includes('timeout')) {
        return { success: false, error: 'API did not respond. Start the API server (e.g. docker compose up -d api).' }
      }
      if (msg.includes('Failed to fetch') || msg.includes('no response')) {
        return { success: false, error: 'Cannot reach the API. Ensure the API server is running on port 8080.' }
      }
      return {
        success: false,
        error: error?.data?.message || error?.message || 'Login failed'
      }
    }
  }

  /**
   * Logout user and clear token
   */
  const logout = () => {
    token.value = null
    navigateTo('/login')
  }

  /**
   * Check if user is authenticated
   */
  const isAuthenticated = computed(() => !!token.value)

  /**
   * Decode JWT token to get user info
   */
  const decodeToken = () => {
    if (!token.value) return null

    try {
      const base64Url = token.value.split('.')[1]
      const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/')
      const jsonPayload = decodeURIComponent(
        atob(base64)
          .split('')
          .map((c) => '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2))
          .join('')
      )
      return JSON.parse(jsonPayload)
    } catch (error) {
      console.error('Failed to decode token:', error)
      return null
    }
  }

  /**
   * Get current user info from token
   */
  const currentUser = computed(() => decodeToken())

  const API_TIMEOUT_MS = 30000 // 30s so colder project details requests have more time to respond
  /** Timeout for long-running requests (e.g. Google Slides import: download + 62 slide images). */
  const API_LONG_TIMEOUT_MS = 5 * 60 * 1000 // 5 minutes

  /**
   * Fetch with automatic JWT authentication and timeout (so UI does not spin forever if API is down).
   * Options may include timeoutMs to override the default (e.g. timeoutMs: 300000 for import).
   */
  const fetchWithAuth = async <T = any>(url: string, options: any = {}) => {
    if (!token.value) {
      throw new Error('No authentication token found')
    }

    const timeoutMs = options.timeoutMs ?? API_TIMEOUT_MS
    const { timeoutMs: _drop, ...restOptions } = options

    const fullUrl = url.startsWith('http') ? url : `${apiBase.value}${url.startsWith('/') ? url : '/' + url}`
    const controller = new AbortController()
    const timeoutId = setTimeout(() => controller.abort(), timeoutMs)

    try {
      const result = await $fetch<T>(fullUrl, {
        ...restOptions,
        signal: controller.signal,
        headers: {
          ...restOptions.headers,
          Authorization: `Bearer ${token.value}`
        }
      })
      clearTimeout(timeoutId)
      return result
    } catch (e: any) {
      clearTimeout(timeoutId)
      const status = e?.statusCode ?? e?.response?.status
      if (status === 401 && typeof token.value !== 'undefined') {
        token.value = null
        if (import.meta.client) {
          navigateTo('/login?session=expired')
        }
        throw new Error('Session expired or invalid. Please log in again.')
      }
      const msg = e?.message || String(e)
      if (e?.name === 'AbortError' || (msg && (msg.includes('abort') || msg.includes('timeout') || msg.includes('timed out')))) {
        throw new Error('API did not respond in time. Start the API server (e.g. docker compose up -d api or go run ./cmd/server in the api folder).')
      }
      if (msg && (msg.includes('Failed to fetch') || msg.includes('no response') || msg.includes('ERR_CONNECTION_RESET') || msg.includes('ERR_EMPTY_RESPONSE'))) {
        throw new Error('Cannot reach the API. Ensure the API is running (port 8080). In Docker, ensure the api service is up and NUXT_PUBLIC_API_BASE_SERVER is set for SSR.')
      }
      throw e
    }
  }

  return {
    token,
    login,
    logout,
    isAuthenticated,
    fetchWithAuth,
    currentUser,
    apiBase
  }
}

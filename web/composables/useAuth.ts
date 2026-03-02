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
   * Login user and save JWT token
   */
  const login = async (email: string, password: string) => {
    try {
      // Use full URL for SSR compatibility
      const response = await $fetch<LoginResponse>('http://localhost:8080/api/v1/auth/login', {
        method: 'POST',
        body: { email, password }
      })

      if (response.data?.token) {
        token.value = response.data.token
        return { success: true, user: response.data.user }
      }

      return { success: false, error: 'No token received' }
    } catch (error: any) {
      console.error('Login failed:', error)
      return { 
        success: false, 
        error: error.data?.message || error.message || 'Login failed' 
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

  /**
   * Fetch with automatic JWT authentication
   * Usage: fetchWithAuth('/sentinel/tasks')
   */
  const fetchWithAuth = async <T = any>(url: string, options: any = {}) => {
    if (!token.value) {
      throw new Error('No authentication token found')
    }

    // Use full URL for SSR compatibility
    const fullUrl = url.startsWith('http') ? url : `http://localhost:8080/api/v1${url}`
    return await $fetch<T>(fullUrl, {
      ...options,
      headers: {
        ...options.headers,
        Authorization: `Bearer ${token.value}`
      }
    })
  }

  return {
    token,
    login,
    logout,
    isAuthenticated,
    fetchWithAuth,
    currentUser
  }
}

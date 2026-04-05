/**
 * Authentication API Integration
 * 
 * This layer handles all API calls related to authentication.
 * Uses the shared useHttp composable for consistent API communication.
 * 
 * Follows Feature-Sliced Design (FSD) - Infrastructure layer
 */

import { useAuth } from '~/composables/useAuth'
import { useHttp } from '~/core/shared/api/http'

interface RegisterRequest {
  email: string
  password: string
  confirm_password: string
}

interface LoginRequest {
  email: string
  password: string
}

export interface User {
  id: number
  email: string
  role: string // CEO, MANAGER, PRODUCT_OWNER, ENGINEER, CHIEF_ENGINEER, SUPPORT
  created_at: string
  updated_at: string
  display_name?: string
  health_score?: number
  tech_stack?: string[]
}

interface AuthResponse {
  token: string
  user: User
}

interface ApiResponse<T> {
  message: string
  data: T
}

/**
 * Auth API Module
 * All authentication-related API calls
 */
export const authApi = {
  /**
   * Register a new user
   * 
   * POST /auth/register
   * 
   * @param email - User email
   * @param password - User password
   * @param confirmPassword - Password confirmation
   * @returns Promise with auth response (token + user)
   */
  async register(email: string, password: string, confirmPassword: string) {
    const http = useHttp()
    const response = await http.post<ApiResponse<AuthResponse>>('/auth/register', {
      email,
      password,
      confirm_password: confirmPassword,
    })

    if (response.error) {
      throw new Error(response.error.message || 'Registration failed')
    }

    return response.data
  },

  /**
   * Login user
   * 
   * POST /auth/login
   * 
   * @param email - User email
   * @param password - User password
   * @returns Promise with auth response (token + user)
   */
  async login(email: string, password: string) {
    const http = useHttp()
    const response = await http.post<ApiResponse<AuthResponse>>('/auth/login', {
      email,
      password,
    })

    if (response.error) {
      throw new Error(response.error.message || 'Login failed')
    }

    return response.data
  },

  /**
   * Validate JWT token (future implementation)
   * 
   * GET /auth/validate
   * 
   * @param token - JWT token to validate
   * @returns Promise with user data if valid
   */
  async validateToken(token: string) {
    // TODO: Implement token validation endpoint in backend
    // This will be used to check if token is still valid on app initialization
    const http = useHttp()
    const response = await http.get<ApiResponse<User>>('/auth/validate', {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    })

    if (response.error) {
      throw new Error(response.error.message || 'Token validation failed')
    }

    return response.data
  },

  /**
   * Get current user profile (authenticated)
   * GET /auth/me
   */
  async getMe(): Promise<User> {
    const { fetchWithAuth } = useAuth()
    const res = await fetchWithAuth<{ message: string; data: User }>('/auth/me')
    return res.data
  },

  /**
   * Update own profile
   * PATCH /auth/me
   */
  async updateProfile(payload: { display_name?: string; tech_stack?: string[] }) {
    const { fetchWithAuth } = useAuth()
    const res = await fetchWithAuth<{ message: string; data: User }>('/auth/me', {
      method: 'PATCH',
      body: payload,
    })
    return res.data
  },

  /**
   * Change own password
   * PATCH /auth/me/password
   */
  async changePassword(currentPassword: string, newPassword: string) {
    const { fetchWithAuth } = useAuth()
    await fetchWithAuth<{ message: string }>('/auth/me/password', {
      method: 'PATCH',
      body: { current_password: currentPassword, new_password: newPassword },
    })
  },
}

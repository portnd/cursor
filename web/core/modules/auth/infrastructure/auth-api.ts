/**
 * Authentication API Integration
 * 
 * This layer handles all API calls related to authentication.
 * Uses the shared useHttp composable for consistent API communication.
 * 
 * Follows Feature-Sliced Design (FSD) - Infrastructure layer
 */

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

interface User {
  id: number
  email: string
  created_at: string
  updated_at: string
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
}

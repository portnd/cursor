/**
 * Authentication Store (Pinia)
 * 
 * This is the centralized state management for authentication.
 * Follows Feature-Sliced Design (FSD) architecture.
 * 
 * Features:
 * - User state management
 * - Token management with cookies
 * - Login/Register/Logout actions
 * - Auto-initialization from stored token
 */

import { defineStore } from 'pinia'
import { authApi } from '../infrastructure/auth-api'

interface User {
  id: number
  email: string
  role: string // CEO, PM, DEV
  created_at: string
  updated_at: string
}

interface AuthState {
  user: User | null
  token: string | null
  isAuthenticated: boolean
  isLoading: boolean
  error: string | null
}

export const useAuthStore = defineStore('auth', {
  state: (): AuthState => ({
    user: null,
    token: null,
    isAuthenticated: false,
    isLoading: false,
    error: null,
  }),

  getters: {
    /**
     * Check if user is logged in
     */
    isLoggedIn: (state) => state.isAuthenticated && state.user !== null,

    /**
     * Get current user email
     */
    userEmail: (state) => state.user?.email || '',

    /**
     * Get current user ID
     */
    userId: (state) => state.user?.id || null,
  },

  actions: {
    /**
     * Initialize auth state from stored token (on app load)
     * This checks if user has a valid token cookie
     */
    async initialize() {
      const tokenCookie = useCookie('token')
      const token = tokenCookie.value

      if (token) {
        this.token = token
        this.isAuthenticated = true
        
        // Try to restore user data from localStorage
        if (process.client) {
          const storedUser = localStorage.getItem('user')
          if (storedUser) {
            try {
              this.user = JSON.parse(storedUser)
            } catch (error) {
              console.error('Failed to parse stored user data:', error)
            }
          }
        }
      }
    },

    /**
     * Register a new user
     * 
     * @param email - User email
     * @param password - User password
     * @param confirmPassword - Password confirmation
     */
    async register(email: string, password: string, confirmPassword: string) {
      this.isLoading = true
      this.error = null

      try {
        const response = await authApi.register(email, password, confirmPassword)

        if (response && response.data) {
          // Store token in cookie (expires in 7 days)
          const tokenCookie = useCookie('token', {
            maxAge: 60 * 60 * 24 * 7, // 7 days
            path: '/',
            sameSite: 'lax',
          })
          tokenCookie.value = response.data.token

          // Update store state
          this.token = response.data.token
          this.user = response.data.user
          this.isAuthenticated = true

          // Store user data in localStorage for persistence
          if (process.client) {
            localStorage.setItem('user', JSON.stringify(response.data.user))
          }

          return { success: true }
        }
      } catch (error: any) {
        this.error = error.message || 'Registration failed'
        return { success: false, error: this.error }
      } finally {
        this.isLoading = false
      }
    },

    /**
     * Login user
     * 
     * @param email - User email
     * @param password - User password
     */
    async login(email: string, password: string) {
      this.isLoading = true
      this.error = null

      try {
        const response = await authApi.login(email, password)

        if (response && response.data) {
          // Store token in cookie (expires in 7 days)
          const tokenCookie = useCookie('token', {
            maxAge: 60 * 60 * 24 * 7, // 7 days
            path: '/',
            sameSite: 'lax',
          })
          tokenCookie.value = response.data.token

          // Update store state
          this.token = response.data.token
          this.user = response.data.user
          this.isAuthenticated = true

          // Store user data in localStorage for persistence
          if (process.client) {
            localStorage.setItem('user', JSON.stringify(response.data.user))
          }

          return { success: true }
        }
      } catch (error: any) {
        this.error = error.message || 'Login failed'
        return { success: false, error: this.error }
      } finally {
        this.isLoading = false
      }
    },

    /**
     * Logout user
     * Clears token, user data, and redirects to login page
     */
    async logout() {
      // Clear cookie
      const tokenCookie = useCookie('token')
      tokenCookie.value = null

      // Clear localStorage
      if (process.client) {
        localStorage.removeItem('user')
      }

      // Clear store state
      this.token = null
      this.user = null
      this.isAuthenticated = false
      this.error = null

      // Redirect to login page
      await navigateTo('/login')
    },

    /**
     * Clear error message
     */
    clearError() {
      this.error = null
    },
  },
})

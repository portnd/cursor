/**
 * Authentication Middleware
 * 
 * This middleware protects routes that require authentication.
 * If user is not logged in, redirect to login page.
 * 
 * Usage:
 * In your page, add:
 * definePageMeta({
 *   middleware: 'auth'
 * })
 */

export default defineNuxtRouteMiddleware((to, from) => {
  // Check if token exists in cookie
  const token = useCookie('token')

  // If no token, redirect to login
  if (!token.value) {
    return navigateTo('/login')
  }

  // TODO: Optionally validate token with API
  // This would require a backend endpoint to validate JWT
  // For now, we just check if token exists
})

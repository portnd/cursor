/**
 * 🔐 Authentication Middleware
 * Protects routes that require authentication
 */

export default defineNuxtRouteMiddleware((to, from) => {
  const { isAuthenticated } = useAuth()

  // If trying to access protected route without auth
  if (!isAuthenticated.value && to.path !== '/login') {
    return navigateTo('/login')
  }

  // If authenticated and trying to access login page, redirect to dashboard
  if (isAuthenticated.value && to.path === '/login') {
    return navigateTo('/dashboard')
  }
})

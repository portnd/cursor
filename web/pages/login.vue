<template>
  <div class="login-page min-h-screen flex items-center justify-center p-4 relative overflow-hidden" :class="isDark ? 'login-dark' : 'login-light'">

    <!-- Background orbs -->
    <div class="absolute inset-0 overflow-hidden pointer-events-none">
      <div class="login-orb-1 absolute rounded-full blur-3xl animate-pulse" />
      <div class="login-orb-2 absolute rounded-full blur-3xl animate-pulse" style="animation-delay:1.2s" />
      <div class="login-orb-3 absolute rounded-full blur-3xl animate-pulse" style="animation-delay:2.4s" />
    </div>

    <!-- Login Card -->
    <div class="relative w-full max-w-md">

      <!-- Logo & Title -->
      <div class="text-center mb-8">
        <div class="inline-flex items-center justify-center w-16 h-16 rounded-2xl login-icon-bg mb-4 shadow-lg">
          <span class="text-3xl">🛡️</span>
        </div>
        <h1 class="text-4xl font-bold mb-2 tracking-tight">
          <span class="bg-gradient-to-r from-purple-500 via-violet-500 to-pink-500 bg-clip-text text-transparent">
            The Sentinel
          </span>
        </h1>
        <p class="login-subtitle text-sm font-medium tracking-wide">AI-Powered Development OS</p>
      </div>

      <!-- Card -->
      <div class="login-card rounded-2xl p-8">

        <h2 class="login-heading text-xl font-bold mb-1">Welcome back</h2>
        <p class="login-subheading text-sm mb-6">Sign in to your workspace</p>

        <!-- Session expired notice -->
        <div v-if="sessionExpired" class="mb-4 p-3.5 rounded-xl border login-alert-warning flex items-start gap-3 text-sm">
          <svg class="w-4 h-4 mt-0.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01M10.29 3.86L1.82 18a2 2 0 001.71 3h16.94a2 2 0 001.71-3L13.71 3.86a2 2 0 00-3.42 0z"/>
          </svg>
          <span>Session expired. Please sign in again.</span>
        </div>

        <!-- Error Message -->
        <div v-if="errorMessage" class="mb-4 p-3.5 rounded-xl border login-alert-error flex items-start gap-3 text-sm">
          <svg class="w-4 h-4 mt-0.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"/>
          </svg>
          <span>{{ errorMessage }}</span>
        </div>

        <form @submit.prevent="handleLogin" class="space-y-4">
          <div>
            <label for="email" class="login-label block text-xs font-semibold uppercase tracking-wider mb-1.5">
              Email Address
            </label>
            <input
              id="email"
              v-model="email"
              type="email"
              required
              placeholder="mail@komgrip.com"
              class="login-input w-full px-4 py-3 rounded-xl text-sm transition-all duration-150 focus:outline-none"
            />
          </div>

          <div>
            <label for="password" class="login-label block text-xs font-semibold uppercase tracking-wider mb-1.5">
              Password
            </label>
            <input
              id="password"
              v-model="password"
              type="password"
              required
              placeholder="••••••••"
              class="login-input w-full px-4 py-3 rounded-xl text-sm transition-all duration-150 focus:outline-none"
            />
          </div>

          <div class="pt-1">
            <button
              type="submit"
              :disabled="isLoading"
              class="w-full py-3 px-4 bg-gradient-to-r from-violet-100 dark:from-violet-600 to-purple-100 dark:to-purple-600 hover:from-violet-200 dark:hover:from-violet-500 hover:to-purple-200 dark:hover:to-purple-500 text-gray-900 dark:text-white font-semibold rounded-xl shadow-lg hover:shadow-violet-500/40 transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2.5 text-sm"
            >
              <svg v-if="isLoading" class="animate-spin h-4 w-4" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"/>
              </svg>
              <svg v-else class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 16l-4-4m0 0l4-4m-4 4h14m-5 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h7a3 3 0 013 3v1"/>
              </svg>
              <span>{{ isLoading ? 'Signing in…' : 'Sign In' }}</span>
            </button>
          </div>
        </form>
      </div>

      <!-- Footer with theme toggle -->
      <div class="mt-5 flex items-center justify-between px-1">
        <p class="login-footer text-xs">Powered by AI · Sentinel OS</p>
        <button
          type="button"
          @click="toggleTheme"
          class="login-theme-btn flex items-center gap-1.5 text-xs px-2.5 py-1.5 rounded-lg transition-all duration-200"
          :title="isDark ? 'Switch to Light Mode' : 'Switch to Dark Mode'"
        >
          <svg v-if="isDark" class="w-3.5 h-3.5" fill="none" stroke="currentColor" stroke-width="1.75" viewBox="0 0 24 24">
            <circle cx="12" cy="12" r="4"/>
            <path stroke-linecap="round" d="M12 2v2M12 20v2M4.22 4.22l1.42 1.42M18.36 18.36l1.42 1.42M2 12h2M20 12h2M4.22 19.78l1.42-1.42M18.36 5.64l1.42-1.42"/>
          </svg>
          <svg v-else class="w-3.5 h-3.5" fill="none" stroke="currentColor" stroke-width="1.75" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" d="M21 12.79A9 9 0 1111.21 3 7 7 0 0021 12.79z"/>
          </svg>
          <span>{{ isDark ? 'Light' : 'Dark' }}</span>
        </button>
      </div>

    </div>
  </div>
</template>

<style scoped>
/* ── Dark ── */
.login-dark {
  background:
    radial-gradient(900px 700px at 75% -10%, rgba(139,92,246,.22), transparent 65%),
    radial-gradient(700px 500px at -5% 80%,  rgba(59,130,246,.16), transparent 65%),
    linear-gradient(160deg, #070b17 0%, #0d1428 50%, #090e1c 100%);
}
.login-dark .login-orb-1 {
  top:5%; left:15%; width:28rem; height:28rem;
  background: rgba(139,92,246,.18);
}
.login-dark .login-orb-2 {
  bottom:10%; right:10%; width:26rem; height:26rem;
  background: rgba(236,72,153,.14);
}
.login-dark .login-orb-3 {
  top:50%; left:50%; width:20rem; height:20rem;
  background: rgba(59,130,246,.10);
}
.login-dark .login-icon-bg {
  background: rgba(124,58,237,.2);
  border: 1px solid rgba(139,92,246,.3);
}
.login-dark .login-subtitle   { color: #94A3B8; }
.login-dark .login-card {
  background: rgba(255,255,255,.04);
  border: 1px solid rgba(255,255,255,.08);
  box-shadow: 0 32px 80px rgba(0,0,0,.5), 0 8px 24px rgba(0,0,0,.35), inset 0 1px 0 rgba(255,255,255,.05);
  backdrop-filter: blur(24px) saturate(140%);
}
.login-dark .login-heading    { color: #F1F5F9; }
.login-dark .login-subheading { color: #64748B; }
.login-dark .login-label      { color: #64748B; letter-spacing: .06em; }
.login-dark .login-input {
  background: rgba(255,255,255,.05);
  border: 1px solid rgba(255,255,255,.1);
  color: #F1F5F9;
}
.login-dark .login-input::placeholder { color: #475569; }
.login-dark .login-input:hover  { border-color: rgba(255,255,255,.16); }
.login-dark .login-input:focus  { border-color: #7C3AED; box-shadow: 0 0 0 3px rgba(124,58,237,.2); }
.login-dark .login-alert-warning {
  background: rgba(180,83,9,.15); border-color: rgba(251,191,36,.3); color: #FCD34D;
}
.login-dark .login-alert-error {
  background: rgba(185,28,28,.15); border-color: rgba(248,113,113,.3); color: #FCA5A5;
}
.login-dark .login-footer     { color: #374151; }
.login-dark .login-theme-btn {
  color: #6B7280;
  background: rgba(255,255,255,.04);
  border: 1px solid rgba(255,255,255,.06);
}
.login-dark .login-theme-btn:hover {
  color: #F1F5F9;
  background: rgba(255,255,255,.08);
  border-color: rgba(255,255,255,.12);
}

/* ── Light (Arctic Violet) ── */
.login-light {
  background:
    radial-gradient(1000px 700px at 75% -5%,  rgba(124,58,237,.09), transparent 70%),
    radial-gradient(700px  500px at -2% 85%,  rgba(219,39,119,.06), transparent 65%),
    linear-gradient(150deg, #EEF1F9 0%, #F2F0FB 40%, #EDF0F8 70%, #F0EEF9 100%);
}
.login-light .login-orb-1 {
  top:8%; left:20%; width:26rem; height:26rem;
  background: rgba(124,58,237,.07);
}
.login-light .login-orb-2 {
  bottom:12%; right:12%; width:24rem; height:24rem;
  background: rgba(219,39,119,.05);
}
.login-light .login-orb-3 {
  top:45%; left:45%; width:16rem; height:16rem;
  background: rgba(99,102,241,.04);
}
.login-light .login-icon-bg {
  background: #F5F3FF;
  border: 1px solid #DDD6FE;
  box-shadow: 0 4px 16px rgba(124,58,237,.12);
}
.login-light .login-subtitle   { color: #6B7A9A; }
.login-light .login-card {
  background: rgba(255,255,255,.92);
  border: 1px solid rgba(210,216,236,.8);
  box-shadow: 0 40px 100px rgba(14,17,36,.14), 0 16px 40px rgba(14,17,36,.08), inset 0 1px 0 rgba(255,255,255,.9);
  backdrop-filter: blur(24px) saturate(160%);
}
.login-light .login-heading    { color: #0E1117; }
.login-light .login-subheading { color: #6B7A9A; }
.login-light .login-label      { color: #9DA6BD; letter-spacing: .07em; }
.login-light .login-input {
  background: #FFFFFF;
  border: 1px solid #D2D8EC;
  color: #0E1117;
  box-shadow: 0 1px 2px rgba(14,17,36,.04) inset;
}
.login-light .login-input::placeholder { color: #9DA6BD; }
.login-light .login-input:hover  { border-color: #B8C0D8; }
.login-light .login-input:focus  {
  border-color: #7C3AED;
  box-shadow: 0 0 0 3px rgba(124,58,237,.12), 0 1px 2px rgba(14,17,36,.04) inset;
}
.login-light .login-alert-warning {
  background: #FFFBEB; border-color: #FCD34D; color: #B45309;
}
.login-light .login-alert-error {
  background: #FEF2F2; border-color: #FCA5A5; color: #DC2626;
}
.login-light .login-footer     { color: #9DA6BD; }
.login-light .login-theme-btn {
  color: #9DA6BD;
  background: rgba(14,17,36,.04);
  border: 1px solid rgba(14,17,36,.06);
}
.login-light .login-theme-btn:hover {
  color: #374261;
  background: rgba(124,58,237,.06);
  border-color: rgba(124,58,237,.12);
}
</style>

<script setup lang="ts">
import { useAuthStore } from '~/core/modules/auth/store/auth-store'

definePageMeta({ layout: false })

const { isDark, toggle: toggleTheme, initTheme, setTheme } = useTheme()
onMounted(() => initTheme())

const { login } = useAuth()
const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const email = ref('')
const password = ref('')
const isLoading = ref(false)
const errorMessage = ref('')

const sessionExpired = ref(false)
onMounted(() => {
  if (route.query.session === 'expired') {
    sessionExpired.value = true
    if (import.meta.client) router.replace({ path: '/login', query: {} })
  }
})

const handleLogin = async () => {
  isLoading.value = true
  errorMessage.value = ''

  try {
    const result = await login(email.value, password.value)
    
    if (result.success) {
      // Sync token + user into auth store so role-based UI (Edit button, etc.) works
      const tokenCookie = useCookie('token')
      if (result.token) authStore.token = result.token
      else if (tokenCookie.value) authStore.token = tokenCookie.value
      if (result.user) {
        authStore.user = result.user
        authStore.isAuthenticated = true
        if (import.meta.client) {
          localStorage.setItem('user', JSON.stringify(result.user))
        }
        // Apply user's saved theme preference immediately after login
        if (result.user.theme_preference) {
          initTheme(result.user.theme_preference as 'dark' | 'light')
        }
      }
      await router.push('/dashboard')
    } else {
      errorMessage.value = result.error || 'Login failed. Please check your credentials.'
    }
  } catch (error: any) {
    errorMessage.value = error.message || 'An unexpected error occurred'
  } finally {
    isLoading.value = false
  }
}
</script>

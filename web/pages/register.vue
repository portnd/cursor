<template>
  <div
    class="login-page min-h-screen flex items-center justify-center p-4 relative overflow-hidden"
    :class="isDark ? 'login-dark' : 'login-light'"
  >
    <div class="absolute inset-0 overflow-hidden pointer-events-none">
      <div class="login-orb-1 absolute rounded-full blur-3xl animate-pulse" />
      <div class="login-orb-2 absolute rounded-full blur-3xl animate-pulse" style="animation-delay:1.2s" />
    </div>

    <div class="relative z-10 w-full max-w-md">
      <div class="text-center mb-8">
        <div class="inline-flex items-center justify-center w-16 h-16 rounded-2xl login-icon-bg mb-4 shadow-lg">
          <span class="text-3xl">🛡️</span>
        </div>
        <h1 class="text-4xl font-bold mb-2 tracking-tight">
          <span class="bg-gradient-to-r from-purple-500 via-violet-500 to-pink-500 bg-clip-text text-transparent">
            The Sentinel
          </span>
        </h1>
        <p class="login-subtitle text-sm font-medium">Create your account</p>
      </div>

      <div class="login-card rounded-2xl p-8">
        <RegisterForm />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from '~/core/modules/auth/store/auth-store'
import RegisterForm from '~/core/modules/auth/ui/RegisterForm.vue'

const authStore = useAuthStore()
const router = useRouter()
const { isDark, initTheme } = useTheme()

onMounted(() => {
  initTheme()
  if (authStore.isLoggedIn) router.push('/')
})

definePageMeta({ layout: false })

useHead({ title: 'Register - Komgrip Task Management' })
</script>

<style scoped>
.login-dark {
  background:
    radial-gradient(900px 700px at 75% -10%, rgba(139,92,246,.22), transparent 65%),
    radial-gradient(700px 500px at -5% 80%,  rgba(59,130,246,.16), transparent 65%),
    linear-gradient(160deg, #070b17 0%, #0d1428 50%, #090e1c 100%);
}
.login-dark .login-orb-1 { top:5%; left:15%; width:28rem; height:28rem; background:rgba(139,92,246,.18); }
.login-dark .login-orb-2 { bottom:10%; right:10%; width:26rem; height:26rem; background:rgba(236,72,153,.14); }
.login-dark .login-icon-bg { background:rgba(124,58,237,.2); border:1px solid rgba(139,92,246,.3); }
.login-dark .login-subtitle { color:#94A3B8; }
.login-dark .login-card {
  background:rgba(255,255,255,.04); border:1px solid rgba(255,255,255,.08);
  box-shadow:0 32px 80px rgba(0,0,0,.5), 0 8px 24px rgba(0,0,0,.35), inset 0 1px 0 rgba(255,255,255,.05);
  backdrop-filter:blur(24px) saturate(140%);
}

.login-light {
  background:
    radial-gradient(1000px 700px at 75% -5%,  rgba(124,58,237,.09), transparent 70%),
    radial-gradient(700px  500px at -2% 85%,  rgba(219,39,119,.06), transparent 65%),
    linear-gradient(150deg, #EEF1F9 0%, #F2F0FB 40%, #EDF0F8 70%, #F0EEF9 100%);
}
.login-light .login-orb-1 { top:8%; left:20%; width:26rem; height:26rem; background:rgba(124,58,237,.07); }
.login-light .login-orb-2 { bottom:12%; right:12%; width:24rem; height:24rem; background:rgba(219,39,119,.05); }
.login-light .login-icon-bg { background:#F5F3FF; border:1px solid #DDD6FE; box-shadow:0 4px 16px rgba(124,58,237,.12); }
.login-light .login-subtitle { color:#6B7A9A; }
.login-light .login-card {
  background:rgba(255,255,255,.92); border:1px solid rgba(210,216,236,.8);
  box-shadow:0 40px 100px rgba(14,17,36,.14), 0 16px 40px rgba(14,17,36,.08), inset 0 1px 0 rgba(255,255,255,.9);
  backdrop-filter:blur(24px) saturate(160%);
}
</style>

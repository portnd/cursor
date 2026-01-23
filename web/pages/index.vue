<template>
  <div class="min-h-screen bg-gradient-to-br from-slate-900 via-purple-900 to-slate-900">
    <div class="container mx-auto px-4 py-16">
      <!-- Hero Section -->
      <div class="text-center mb-16">
        <div class="inline-block mb-6">
          <div class="flex items-center justify-center w-20 h-20 mx-auto bg-gradient-to-r from-purple-600 to-pink-600 rounded-2xl shadow-2xl">
            <svg class="w-12 h-12 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
            </svg>
          </div>
        </div>
        
        <h1 class="text-6xl md:text-7xl font-black text-white mb-6 tracking-tight">
          🛡️ KOMGRIP
        </h1>
        
        <p class="text-3xl md:text-4xl font-bold bg-gradient-to-r from-purple-400 to-pink-400 bg-clip-text text-transparent mb-4">
          God-Tier Starter Kit
        </p>
        
        <p class="text-xl text-gray-300 max-w-3xl mx-auto mb-8">
          Production-ready, scalable monorepo built for national-level applications.
          <br />
          <span class="text-purple-400 font-semibold">Hexagonal Architecture</span> + 
          <span class="text-pink-400 font-semibold"> Feature-Sliced Design</span>
        </p>

        <!-- Auth Section -->
        <div class="flex items-center justify-center gap-4 mb-12">
          <!-- If logged in -->
          <div v-if="authStore.isLoggedIn" class="flex items-center gap-4">
            <div class="bg-white/10 backdrop-blur-lg rounded-full px-6 py-3 border border-white/20">
              <span class="text-white">👋 Welcome, <span class="font-semibold text-purple-400">{{ authStore.userEmail }}</span></span>
            </div>
            <button
              @click="authStore.logout()"
              class="bg-gradient-to-r from-red-500 to-red-600 text-white font-semibold py-3 px-8 rounded-full hover:from-red-600 hover:to-red-700 transition-all duration-200 shadow-lg hover:shadow-xl"
            >
              Logout
            </button>
          </div>
          
          <!-- If not logged in -->
          <div v-else class="flex gap-4">
            <NuxtLink
              to="/login"
              class="bg-white/10 backdrop-blur-lg text-white font-semibold py-3 px-8 rounded-full hover:bg-white/20 transition-all duration-200 border border-white/20 shadow-lg hover:shadow-xl"
            >
              Sign In
            </NuxtLink>
            <NuxtLink
              to="/register"
              class="bg-gradient-to-r from-purple-600 to-pink-600 text-white font-semibold py-3 px-8 rounded-full hover:from-purple-700 hover:to-pink-700 transition-all duration-200 shadow-lg hover:shadow-xl"
            >
              Sign Up
            </NuxtLink>
          </div>
        </div>

        <!-- System Status -->
        <div v-if="healthData" class="max-w-4xl mx-auto mb-12">
          <div class="bg-white/10 backdrop-blur-lg rounded-2xl p-8 border border-white/20 shadow-2xl">
            <div class="flex items-center justify-center mb-6">
              <div 
                :class="[
                  'flex items-center gap-3 px-6 py-3 rounded-full font-bold text-lg',
                  healthData.status === 'UP' 
                    ? 'bg-green-500/20 text-green-300 border-2 border-green-500' 
                    : 'bg-yellow-500/20 text-yellow-300 border-2 border-yellow-500'
                ]"
              >
                <span class="relative flex h-4 w-4">
                  <span 
                    :class="[
                      'animate-ping absolute inline-flex h-full w-full rounded-full opacity-75',
                      healthData.status === 'UP' ? 'bg-green-400' : 'bg-yellow-400'
                    ]"
                  ></span>
                  <span 
                    :class="[
                      'relative inline-flex rounded-full h-4 w-4',
                      healthData.status === 'UP' ? 'bg-green-500' : 'bg-yellow-500'
                    ]"
                  ></span>
                </span>
                {{ healthData.status === 'UP' ? '🚀 System Online' : '⚠️ System Degraded' }}
              </div>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
              <!-- PostgreSQL Status -->
              <div class="bg-black/30 rounded-xl p-6 border border-white/10">
                <div class="flex items-center justify-between mb-2">
                  <span class="text-gray-300 font-semibold">PostgreSQL</span>
                  <span 
                    :class="[
                      'w-3 h-3 rounded-full',
                      healthData.services.postgres === 'UP' ? 'bg-green-500' : 'bg-red-500'
                    ]"
                  ></span>
                </div>
                <p class="text-sm text-gray-400">Primary Database</p>
                <p 
                  :class="[
                    'text-xs font-mono mt-2',
                    healthData.services.postgres === 'UP' ? 'text-green-400' : 'text-red-400'
                  ]"
                >
                  {{ healthData.services.postgres }}
                </p>
              </div>

              <!-- MongoDB Status -->
              <div class="bg-black/30 rounded-xl p-6 border border-white/10">
                <div class="flex items-center justify-between mb-2">
                  <span class="text-gray-300 font-semibold">MongoDB</span>
                  <span 
                    :class="[
                      'w-3 h-3 rounded-full',
                      healthData.services.mongodb === 'UP' ? 'bg-green-500' : 'bg-red-500'
                    ]"
                  ></span>
                </div>
                <p class="text-sm text-gray-400">Logs & Audits</p>
                <p 
                  :class="[
                    'text-xs font-mono mt-2',
                    healthData.services.mongodb === 'UP' ? 'text-green-400' : 'text-red-400'
                  ]"
                >
                  {{ healthData.services.mongodb }}
                </p>
              </div>

              <!-- Redis Status -->
              <div class="bg-black/30 rounded-xl p-6 border border-white/10">
                <div class="flex items-center justify-between mb-2">
                  <span class="text-gray-300 font-semibold">Redis</span>
                  <span 
                    :class="[
                      'w-3 h-3 rounded-full',
                      healthData.services.redis === 'UP' ? 'bg-green-500' : 'bg-red-500'
                    ]"
                  ></span>
                </div>
                <p class="text-sm text-gray-400">Cache (128MB)</p>
                <p 
                  :class="[
                    'text-xs font-mono mt-2',
                    healthData.services.redis === 'UP' ? 'text-green-400' : 'text-red-400'
                  ]"
                >
                  {{ healthData.services.redis }}
                </p>
              </div>
            </div>

            <div class="mt-6 text-center">
              <p class="text-xs text-gray-400 font-mono">
                Last checked: {{ new Date(healthData.timestamp).toLocaleString() }}
              </p>
            </div>
          </div>
        </div>

        <!-- Loading State -->
        <div v-else-if="pending" class="max-w-4xl mx-auto mb-12">
          <div class="bg-white/10 backdrop-blur-lg rounded-2xl p-8 border border-white/20">
            <div class="animate-pulse flex flex-col items-center">
              <div class="h-12 bg-gray-600 rounded-full w-48 mb-4"></div>
              <div class="h-4 bg-gray-600 rounded w-32"></div>
            </div>
          </div>
        </div>

        <!-- Error State -->
        <div v-else-if="error" class="max-w-4xl mx-auto mb-12">
          <div class="bg-red-500/10 backdrop-blur-lg rounded-2xl p-8 border border-red-500/50">
            <div class="flex items-center justify-center gap-3 text-red-400">
              <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <span class="text-lg font-semibold">Failed to connect to API</span>
            </div>
            <p class="text-gray-400 text-sm mt-2">
              Make sure the backend is running on {{ apiBase }}
            </p>
          </div>
        </div>
      </div>

      <!-- Features Grid -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 max-w-6xl mx-auto mb-16">
        <div class="bg-white/5 backdrop-blur-sm rounded-xl p-6 border border-white/10 hover:border-purple-500/50 transition-all duration-300 hover:scale-105">
          <div class="text-4xl mb-4">⚡</div>
          <h3 class="text-xl font-bold text-white mb-2">Blazing Fast</h3>
          <p class="text-gray-400">Hot reload with Air (Go) and Nuxt HMR. Changes reflect instantly.</p>
        </div>

        <div class="bg-white/5 backdrop-blur-sm rounded-xl p-6 border border-white/10 hover:border-pink-500/50 transition-all duration-300 hover:scale-105">
          <div class="text-4xl mb-4">🏗️</div>
          <h3 class="text-xl font-bold text-white mb-2">Clean Architecture</h3>
          <p class="text-gray-400">Hexagonal backend + Feature-Sliced frontend. Enterprise-grade structure.</p>
        </div>

        <div class="bg-white/5 backdrop-blur-sm rounded-xl p-6 border border-white/10 hover:border-blue-500/50 transition-all duration-300 hover:scale-105">
          <div class="text-4xl mb-4">🗄️</div>
          <h3 class="text-xl font-bold text-white mb-2">Multi-Database</h3>
          <p class="text-gray-400">PostgreSQL for ACID, MongoDB for logs, Redis for cache (128MB).</p>
        </div>

        <div class="bg-white/5 backdrop-blur-sm rounded-xl p-6 border border-white/10 hover:border-green-500/50 transition-all duration-300 hover:scale-105">
          <div class="text-4xl mb-4">🔒</div>
          <h3 class="text-xl font-bold text-white mb-2">Production Ready</h3>
          <p class="text-gray-400">Zero placeholders. Every line is battle-tested and deployment-ready.</p>
        </div>

        <div class="bg-white/5 backdrop-blur-sm rounded-xl p-6 border border-white/10 hover:border-yellow-500/50 transition-all duration-300 hover:scale-105">
          <div class="text-4xl mb-4">🐳</div>
          <h3 class="text-xl font-bold text-white mb-2">Docker Native</h3>
          <p class="text-gray-400">Full docker-compose setup. One command to rule them all.</p>
        </div>

        <div class="bg-white/5 backdrop-blur-sm rounded-xl p-6 border border-white/10 hover:border-red-500/50 transition-all duration-300 hover:scale-105">
          <div class="text-4xl mb-4">📊</div>
          <h3 class="text-xl font-bold text-white mb-2">Type Safe</h3>
          <p class="text-gray-400">TypeScript strict mode. Catch errors before runtime.</p>
        </div>
      </div>

      <!-- Tech Stack -->
      <div class="max-w-4xl mx-auto text-center">
        <h2 class="text-3xl font-bold text-white mb-8">Built With</h2>
        <div class="flex flex-wrap justify-center gap-4">
          <span class="px-6 py-3 bg-blue-500/20 text-blue-300 rounded-full font-semibold border border-blue-500/50">Go 1.23</span>
          <span class="px-6 py-3 bg-green-500/20 text-green-300 rounded-full font-semibold border border-green-500/50">Nuxt 3</span>
          <span class="px-6 py-3 bg-purple-500/20 text-purple-300 rounded-full font-semibold border border-purple-500/50">PostgreSQL 15</span>
          <span class="px-6 py-3 bg-green-600/20 text-green-400 rounded-full font-semibold border border-green-600/50">MongoDB 6</span>
          <span class="px-6 py-3 bg-red-500/20 text-red-300 rounded-full font-semibold border border-red-500/50">Redis 7</span>
          <span class="px-6 py-3 bg-cyan-500/20 text-cyan-300 rounded-full font-semibold border border-cyan-500/50">TailwindCSS</span>
          <span class="px-6 py-3 bg-yellow-500/20 text-yellow-300 rounded-full font-semibold border border-yellow-500/50">TypeScript</span>
        </div>
      </div>

      <!-- Footer -->
      <div class="text-center mt-16 text-gray-400">
        <p class="text-sm">
          Built with 💪 by <span class="text-purple-400 font-semibold">Thanandorn</span> for national-scale applications
        </p>
        <p class="text-xs mt-2">© 2026 Komgrip. All rights reserved.</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from '~/core/modules/auth/store/auth-store'

// Protect this page - require authentication
definePageMeta({
  middleware: 'auth'
})

interface HealthResponse {
  status: string
  timestamp: string
  services: {
    postgres: string
    mongodb: string
    redis: string
  }
}

// Initialize auth store
const authStore = useAuthStore()
onMounted(() => {
  authStore.initialize()
})

const config = useRuntimeConfig()
const apiBase = config.public.apiBase

const { data: healthData, pending, error } = await useFetch<HealthResponse>('/health', {
  baseURL: apiBase as string,
  server: false
})

useHead({
  title: 'Komgrip - God-Tier Starter Kit',
  meta: [
    { name: 'description', content: 'Production-ready, scalable monorepo starter kit built for national-level applications' }
  ]
})
</script>

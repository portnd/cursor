<template>
  <div class="flex h-screen bg-gray-900 text-gray-100">
    <!-- Sidebar -->
    <aside class="w-64 bg-gray-800 border-r border-gray-700 flex flex-col">
      <!-- Logo -->
      <div class="p-6 border-b border-gray-700">
        <h1 class="text-2xl font-bold bg-gradient-to-r from-purple-400 to-pink-600 bg-clip-text text-transparent">
          🛡️ The Sentinel
        </h1>
        <p class="text-xs text-gray-400 mt-1">AI-Powered Task Manager</p>
      </div>

      <!-- Navigation -->
      <nav class="flex-1 p-4 space-y-2">
        <NuxtLink
          to="/dashboard"
          class="flex items-center gap-3 px-4 py-3 rounded-lg transition-all hover:bg-gray-700 hover:translate-x-1"
          active-class="bg-gradient-to-r from-purple-600 to-pink-600 shadow-lg"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" />
          </svg>
          <span class="font-medium">Dashboard</span>
        </NuxtLink>

        <NuxtLink
          to="/tasks"
          class="flex items-center gap-3 px-4 py-3 rounded-lg transition-all hover:bg-gray-700 hover:translate-x-1"
          active-class="bg-gradient-to-r from-purple-600 to-pink-600 shadow-lg"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
          </svg>
          <span class="font-medium">📥 Tasks Inbox</span>
        </NuxtLink>

        <NuxtLink
          to="/create"
          class="flex items-center gap-3 px-4 py-3 rounded-lg transition-all hover:bg-gray-700 hover:translate-x-1"
          active-class="bg-gradient-to-r from-purple-600 to-pink-600 shadow-lg"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
          <span class="font-medium">Create Task</span>
        </NuxtLink>

        <!-- CEO Only: Team Management -->
        <NuxtLink
          v-if="currentUser?.role === 'CEO'"
          to="/admin/team"
          class="flex items-center gap-3 px-4 py-3 rounded-lg transition-all hover:bg-gray-700 hover:translate-x-1"
          active-class="bg-gradient-to-r from-purple-600 to-pink-600 shadow-lg"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
          </svg>
          <span class="font-medium">👑 Team Roster</span>
        </NuxtLink>

        <!-- CEO Only: AI Control Tower -->
        <NuxtLink
          v-if="currentUser?.role === 'CEO'"
          to="/admin/ai-settings"
          class="flex items-center gap-3 px-4 py-3 rounded-lg transition-all hover:bg-gradient-to-r hover:from-yellow-600/20 hover:to-orange-600/20 hover:translate-x-1 border border-transparent hover:border-yellow-500/50"
          active-class="bg-gradient-to-r from-yellow-600 to-orange-600 shadow-lg shadow-yellow-500/50 border-yellow-400"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z" />
          </svg>
          <span class="font-medium">⚙️ AI Control Tower</span>
        </NuxtLink>
      </nav>

      <!-- User Section & Logout -->
      <div class="p-4 border-t border-gray-700">
        <div class="flex items-center gap-3 px-4 py-3 bg-gray-700/50 rounded-lg mb-3">
          <div class="w-10 h-10 rounded-full bg-gradient-to-br from-purple-500 to-pink-500 flex items-center justify-center font-bold">
            {{ userInitial }}
          </div>
          <div class="flex-1 min-w-0">
            <p class="text-sm font-medium truncate">{{ userEmail }}</p>
            <p class="text-xs text-gray-400">{{ userRole }}</p>
          </div>
        </div>
        
        <button
          @click="handleLogout"
          class="w-full flex items-center gap-3 px-4 py-3 rounded-lg transition-all hover:bg-red-600/20 hover:text-red-400 text-gray-400"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
          </svg>
          <span class="font-medium">Logout</span>
        </button>
      </div>
    </aside>

    <!-- Main Content -->
    <main class="flex-1 overflow-auto bg-gradient-to-br from-gray-900 via-gray-800 to-gray-900">
      <slot />
    </main>
  </div>
</template>

<script setup lang="ts">
const { logout, currentUser } = useAuth()

const userEmail = computed(() => currentUser.value?.email || 'user@sentinel.com')
const userRole = computed(() => {
  const role = currentUser.value?.role || 'DEV'
  const roleMap: Record<string, string> = {
    'CEO': 'Chief Executive',
    'PM': 'Project Manager',
    'DEV': 'Developer'
  }
  return roleMap[role] || role
})
const userInitial = computed(() => userEmail.value.charAt(0).toUpperCase())

const handleLogout = () => {
  if (confirm('Are you sure you want to logout?')) {
    logout()
  }
}
</script>

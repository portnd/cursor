<template>
  <div class="flex h-screen bg-gray-900 text-gray-100" style="background-color: #111827;">
    <!-- Sidebar (collapsible) - inline fallback so first paint is never white -->
    <aside
      class="flex flex-col bg-gray-800 border-r border-gray-700 transition-[width] duration-200 ease-out shrink-0"
      :class="sidebarCollapsed ? 'w-[4.5rem]' : 'w-64'"
      style="background-color: #1f2937;"
    >
      <!-- Logo + Toggle -->
      <div class="p-4 border-b border-gray-700 flex items-center justify-between gap-2">
        <div class="flex items-center gap-3 min-w-0 overflow-hidden">
          <span class="text-xl shrink-0">🛡️</span>
          <div v-show="!sidebarCollapsed" class="min-w-0">
            <h1 class="text-lg font-bold bg-gradient-to-r from-purple-400 to-pink-600 bg-clip-text text-transparent truncate">The Sentinel</h1>
            <p class="text-[10px] text-gray-400 truncate">AI Task Manager</p>
          </div>
        </div>
        <button
          type="button"
          @click="sidebarCollapsed = !sidebarCollapsed"
          class="shrink-0 p-1.5 rounded-lg text-gray-400 hover:text-white hover:bg-gray-700 transition-colors"
          :title="sidebarCollapsed ? 'ขยายแถบด้านข้าง' : 'ย่อแถบด้านข้าง'"
        >
          <svg class="w-5 h-5 transition-transform" :class="sidebarCollapsed ? 'rotate-180' : ''" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 19l-7-7 7-7m8 14l-7-7 7-7" />
          </svg>
        </button>
      </div>

      <!-- Navigation -->
      <nav class="flex-1 p-3 space-y-1 overflow-hidden">
        <NuxtLink
          to="/dashboard"
          class="nav-link"
          active-class="bg-gradient-to-r from-purple-600 to-pink-600 shadow-lg"
          :title="sidebarCollapsed ? 'Dashboard' : undefined"
        >
          <svg class="w-5 h-5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" /></svg>
          <span v-show="!sidebarCollapsed" class="font-medium truncate">Dashboard</span>
        </NuxtLink>
        <NuxtLink
          to="/tasks"
          class="nav-link"
          active-class="bg-gradient-to-r from-purple-600 to-pink-600 shadow-lg"
          :title="sidebarCollapsed ? 'Tasks Inbox' : undefined"
        >
          <svg class="w-5 h-5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" /></svg>
          <span v-show="!sidebarCollapsed" class="font-medium truncate">Tasks Inbox</span>
        </NuxtLink>
        <NuxtLink
          to="/create"
          class="nav-link"
          active-class="bg-gradient-to-r from-purple-600 to-pink-600 shadow-lg"
          :title="sidebarCollapsed ? 'Create Task' : undefined"
        >
          <svg class="w-5 h-5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" /></svg>
          <span v-show="!sidebarCollapsed" class="font-medium truncate">Create Task</span>
        </NuxtLink>
        <NuxtLink
          to="/projects"
          class="nav-link"
          active-class="bg-gradient-to-r from-purple-600 to-pink-600 shadow-lg"
          :title="sidebarCollapsed ? 'Projects' : undefined"
        >
          <svg class="w-5 h-5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" /></svg>
          <span v-show="!sidebarCollapsed" class="font-medium truncate">Projects</span>
        </NuxtLink>
        <NuxtLink
          v-if="currentUser?.role === 'CEO'"
          to="/accounting"
          class="nav-link"
          active-class="bg-gradient-to-r from-purple-600 to-pink-600 shadow-lg"
          :title="sidebarCollapsed ? 'บัญชี' : undefined"
        >
          <svg class="w-5 h-5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 7h6m0 10v-3m-3 3h.01M9 17h.01M9 14h.01M12 14h.01M15 11h.01M12 11h.01M9 11h.01M7 21h10a2 2 0 002-2V5a2 2 0 00-2-2H7a2 2 0 00-2 2v14a2 2 0 002 2z" /></svg>
          <span v-show="!sidebarCollapsed" class="font-medium truncate">บัญชี</span>
        </NuxtLink>
        <NuxtLink
          v-if="currentUser?.role === 'CEO'"
          to="/performance"
          class="nav-link"
          active-class="bg-gradient-to-r from-purple-600 to-pink-600 shadow-lg"
          :title="sidebarCollapsed ? 'Team Performance' : undefined"
        >
          <svg class="w-5 h-5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7h8m0 0v8m0-8l-3 3m-6 2a2 2 0 11-4 0 2 2 0 014 0zM3 21V3m0 18v-4" /></svg>
          <span v-show="!sidebarCollapsed" class="font-medium truncate">Team Performance</span>
        </NuxtLink>
        <NuxtLink
          v-if="currentUser?.role === 'CEO'"
          to="/admin/team"
          class="nav-link"
          active-class="bg-gradient-to-r from-purple-600 to-pink-600 shadow-lg"
          :title="sidebarCollapsed ? 'User Management' : undefined"
        >
          <svg class="w-5 h-5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" /></svg>
          <span v-show="!sidebarCollapsed" class="font-medium truncate">User Management</span>
        </NuxtLink>
        <NuxtLink
          v-if="currentUser?.role === 'CEO'"
          to="/admin/ai-settings"
          class="nav-link nav-link-ai"
          active-class="bg-gradient-to-r from-yellow-600 to-orange-600 shadow-lg shadow-yellow-500/50 border-yellow-400"
          :title="sidebarCollapsed ? 'AI Control Tower' : undefined"
        >
          <svg class="w-5 h-5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z" /></svg>
          <span v-show="!sidebarCollapsed" class="font-medium truncate">AI Control Tower</span>
        </NuxtLink>
      </nav>

      <!-- User (link to Profile) + Logout -->
      <div class="p-3 border-t border-gray-700 space-y-1">
        <NuxtLink
          to="/profile"
          class="flex items-center gap-3 px-3 py-2.5 rounded-lg overflow-hidden transition-all hover:bg-gray-700"
          active-class="!bg-gradient-to-r !from-purple-600 !to-pink-600 shadow-lg"
          title="Profile &amp; Account"
        >
          <div class="w-9 h-9 rounded-full bg-gradient-to-br from-purple-500 to-pink-500 flex items-center justify-center font-bold shrink-0 text-sm">{{ userInitial }}</div>
          <div v-show="!sidebarCollapsed" class="min-w-0 flex-1">
            <p class="text-sm font-medium truncate">{{ userEmail }}</p>
            <p class="text-xs text-gray-400 truncate">{{ userRole }}</p>
          </div>
        </NuxtLink>
        <button
          @click="handleLogout"
          class="nav-link w-full hover:bg-red-600/20 hover:text-red-400 text-gray-400"
          title="Logout"
        >
          <svg class="w-5 h-5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" /></svg>
          <span v-show="!sidebarCollapsed" class="font-medium">Logout</span>
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
const { confirm } = useNotification()

const SIDEBAR_COLLAPSED_KEY = 'sentinel-sidebar-collapsed'
const sidebarCollapsed = ref(false)

onMounted(() => {
  if (import.meta.client) {
    const saved = localStorage.getItem(SIDEBAR_COLLAPSED_KEY)
    if (saved !== null) sidebarCollapsed.value = saved === '1'
  }
})
watch(sidebarCollapsed, (v) => {
  if (import.meta.client) localStorage.setItem(SIDEBAR_COLLAPSED_KEY, v ? '1' : '0')
})

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

const handleLogout = async () => {
  const ok = await confirm({
    title: 'Logout',
    message: 'Are you sure you want to logout?',
    confirmLabel: 'Logout',
    cancelLabel: 'Cancel',
    variant: 'primary'
  })
  if (ok) logout()
}
</script>

<style scoped>
.nav-link {
  @apply flex items-center gap-3 px-3 py-2.5 rounded-lg transition-all hover:bg-gray-700 hover:translate-x-0.5;
}
.nav-link-ai {
  @apply hover:bg-gradient-to-r hover:from-yellow-600/20 hover:to-orange-600/20 border border-transparent hover:border-yellow-500/50;
}
</style>

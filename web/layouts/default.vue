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
      <nav class="flex-1 p-3 space-y-1 overflow-y-auto scrollbar-thin scrollbar-thumb-gray-600 scrollbar-track-transparent">
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
          v-if="currentUser?.role !== 'SUPPORT'"
          to="/create"
          class="nav-link"
          active-class="bg-gradient-to-r from-purple-600 to-pink-600 shadow-lg"
          :title="sidebarCollapsed ? 'Create Task' : undefined"
        >
          <svg class="w-5 h-5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" /></svg>
          <span v-show="!sidebarCollapsed" class="font-medium truncate">Create Task</span>
        </NuxtLink>
        <NuxtLink
          v-if="currentUser?.role !== 'SUPPORT'"
          to="/projects"
          class="nav-link"
          active-class="bg-gradient-to-r from-purple-600 to-pink-600 shadow-lg"
          :title="sidebarCollapsed ? 'Projects' : undefined"
        >
          <svg class="w-5 h-5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" /></svg>
          <span v-show="!sidebarCollapsed" class="font-medium truncate">Projects</span>
        </NuxtLink>
        <NuxtLink
          to="/pulse"
          class="nav-link"
          active-class="bg-gradient-to-r from-violet-600 to-indigo-600 shadow-lg"
          :title="sidebarCollapsed ? 'Daily Standup' : undefined"
        >
          <svg class="w-5 h-5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" /></svg>
          <span v-show="!sidebarCollapsed" class="font-medium truncate">Daily Standup</span>
        </NuxtLink>

        <!-- Work Log (timer + quick log + EOD + history) -->
        <NuxtLink
          v-if="currentUser?.role !== 'SUPPORT'"
          to="/logtime"
          class="nav-link"
          active-class="bg-gradient-to-r from-purple-600 to-indigo-600 shadow-lg"
          :title="sidebarCollapsed ? 'Work Log' : undefined"
        >
          <svg class="w-5 h-5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <span v-show="!sidebarCollapsed" class="font-medium truncate">Work Log</span>
        </NuxtLink>
        <NuxtLink
          v-if="['PRODUCT_OWNER', 'PM', 'CEO', 'MANAGER'].includes(currentUser?.role ?? '')"
          to="/active-board"
          class="nav-link"
          active-class="bg-gradient-to-r from-indigo-600 to-purple-600 shadow-lg"
          :title="sidebarCollapsed ? 'Global Active Board' : undefined"
        >
          <svg class="w-5 h-5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 17V7m0 10a2 2 0 01-2 2H5a2 2 0 01-2-2V7a2 2 0 012-2h2a2 2 0 012 2m0 10a2 2 0 002 2h2a2 2 0 002-2M9 7a2 2 0 012-2h2a2 2 0 012 2m0 10V7m0 10a2 2 0 002 2h2a2 2 0 002-2V7a2 2 0 00-2-2h-2a2 2 0 00-2 2"/></svg>
          <span v-show="!sidebarCollapsed" class="font-medium truncate">Global Active Board</span>
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
        <!-- Deployment — visible to all engineers + management -->
        <NuxtLink
          v-if="['ENGINEER', 'CHIEF_ENGINEER', 'CEO', 'MANAGER', 'PRODUCT_OWNER', 'PM'].includes(currentUser?.role ?? '')"
          to="/deployment"
          class="nav-link relative"
          active-class="bg-gradient-to-r from-cyan-600 to-blue-600 shadow-lg"
          :title="sidebarCollapsed ? 'Deployment' : undefined"
        >
          <svg class="w-5 h-5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M9 19l3 3m0 0l3-3m-3 3V10" />
          </svg>
          <span v-show="!sidebarCollapsed" class="font-medium truncate">Deployment</span>
          <!-- badge: pending count shown only to CHIEF_ENGINEER — loaded reactively -->
          <span
            v-if="!sidebarCollapsed && currentUser?.role === 'CHIEF_ENGINEER' && deploymentPendingCount > 0"
            class="ml-auto text-[10px] font-bold px-1.5 py-0.5 rounded-full bg-yellow-500 text-gray-900 shrink-0"
          >{{ deploymentPendingCount }}</span>
        </NuxtLink>
        <NuxtLink
          v-if="['CEO', 'MANAGER', 'PRODUCT_OWNER', 'PM', 'ENGINEER', 'CHIEF_ENGINEER', 'SUPPORT'].includes(currentUser?.role ?? '')"
          to="/discipline"
          class="nav-link"
          active-class="bg-gradient-to-r from-orange-600 to-red-600 shadow-lg"
          :title="sidebarCollapsed ? 'Discipline Tracker' : undefined"
        >
          <svg class="w-5 h-5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01"/></svg>
          <span v-show="!sidebarCollapsed" class="font-medium truncate">Discipline Tracker</span>
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
        <NuxtLink
          v-if="currentUser?.role === 'CEO'"
          to="/admin/cost-config"
          class="nav-link"
          active-class="bg-gradient-to-r from-amber-600 to-orange-600 shadow-lg shadow-amber-500/25"
          :title="sidebarCollapsed ? 'Cost Configuration' : undefined"
        >
          <svg class="w-5 h-5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
          <span v-show="!sidebarCollapsed" class="font-medium truncate">Cost Configuration</span>
        </NuxtLink>
        <NuxtLink
          v-if="['CEO', 'MANAGER', 'PRODUCT_OWNER', 'PM', 'ENGINEER', 'CHIEF_ENGINEER', 'DEV', 'SUPPORT'].includes(currentUser?.role ?? '')"
          to="/attendance"
          class="nav-link"
          active-class="bg-gradient-to-r from-teal-600 to-emerald-600 shadow-lg"
          :title="sidebarCollapsed ? 'Attendance' : undefined"
        >
          <svg class="w-5 h-5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10m-12 9h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v11a2 2 0 002 2z" />
          </svg>
          <span v-show="!sidebarCollapsed" class="font-medium truncate">Attendance</span>
        </NuxtLink>
        <NuxtLink
          v-if="['CEO', 'MANAGER', 'PRODUCT_OWNER', 'PM', 'ENGINEER', 'CHIEF_ENGINEER', 'DEV', 'SUPPORT'].includes(currentUser?.role ?? '')"
          to="/leave"
          class="nav-link"
          active-class="bg-gradient-to-r from-cyan-600 to-blue-600 shadow-lg"
          :title="sidebarCollapsed ? 'Leave' : undefined"
        >
          <svg class="w-5 h-5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6M7 4h10a2 2 0 012 2v12a2 2 0 01-2 2H7a2 2 0 01-2-2V6a2 2 0 012-2z" />
          </svg>
          <span v-show="!sidebarCollapsed" class="font-medium truncate">Leave</span>
        </NuxtLink>
        <NuxtLink
          v-if="['CEO', 'SUPPORT'].includes(currentUser?.role ?? '')"
          to="/admin/leave"
          class="nav-link"
          active-class="bg-gradient-to-r from-cyan-600 to-blue-600 shadow-lg"
          :title="sidebarCollapsed ? 'Leave Admin' : undefined"
        >
          <svg class="w-5 h-5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6M7 4h10a2 2 0 012 2v12a2 2 0 01-2 2H7a2 2 0 01-2-2V6a2 2 0 012-2z" />
          </svg>
          <span v-show="!sidebarCollapsed" class="font-medium truncate">Leave Admin</span>
        </NuxtLink>
        <NuxtLink
          v-if="['CEO'].includes(currentUser?.role ?? '')"
          to="/admin/attendance-config"
          class="nav-link"
          active-class="bg-gradient-to-r from-teal-600 to-emerald-600 shadow-lg"
          :title="sidebarCollapsed ? 'Attendance config' : undefined"
        >
          <svg class="w-5 h-5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 11c0 3.517-1.009 6.799-2.753 9.571m-3.44-2.04l.054-.09A13.916 13.916 0 008 11a4 4 0 118 0c0 1.017-.07 2.019-.203 3m-2.118 6.844A21.88 21.88 0 0015.171 17m3.839 1.132l.857-6a4 4 0 00-3.838-4.659h-1.5a4 4 0 00-3.838 4.659l.857 6" />
          </svg>
          <span v-show="!sidebarCollapsed" class="font-medium truncate">Attendance config</span>
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

  <!-- Quick Log Time Modal (globally accessible) -->
  <TasksQuickLogTimeModal v-model="showQuickLog" />

  <!-- EOD Batch Log Modal -->
  <TasksBulkEodLoggerModal :show="showBulkLog" @close="showBulkLog = false" @done="showBulkLog = false" />
</template>

<script setup lang="ts">
import { useDeploymentApi } from '~/core/modules/deployment/infrastructure/deployment-api'

const { logout, currentUser } = useAuth()
const { confirm } = useNotification()

// Live pending deployment count badge for CHIEF_ENGINEER
const deploymentPendingCount = ref(0)
const deploymentApi = useDeploymentApi()

async function refreshDeploymentBadge() {
  if (currentUser.value?.role !== 'CHIEF_ENGINEER') return
  try {
    const s = await deploymentApi.getStats()
    deploymentPendingCount.value = (s.total_pending ?? 0) + (s.total_reviewing ?? 0)
  } catch { /* silent */ }
}

onMounted(() => {
  refreshDeploymentBadge()
  // Refresh badge every 60 s so Chief Engineer always sees live count
  const interval = setInterval(refreshDeploymentBadge, 60_000)
  onUnmounted(() => clearInterval(interval))
})

const SIDEBAR_COLLAPSED_KEY = 'sentinel-sidebar-collapsed'
const sidebarCollapsed = ref(false)
const showQuickLog = ref(false)
const showBulkLog = ref(false)

onMounted(() => {
  if (import.meta.client) {
    const saved = localStorage.getItem(SIDEBAR_COLLAPSED_KEY)
    if (saved !== null) sidebarCollapsed.value = saved === '1'
  }
})
watch(sidebarCollapsed, (v) => {
  if (import.meta.client) localStorage.setItem(SIDEBAR_COLLAPSED_KEY, v ? '1' : '0')
})

// ⌘+L → Quick Log Time | ⌘+Shift+L → EOD Batch Log
function onKeydown(e: KeyboardEvent) {
  if ((e.metaKey || e.ctrlKey) && e.key === 'l') {
    e.preventDefault()
    if (e.shiftKey) {
      showBulkLog.value = true
    } else {
      showQuickLog.value = true
    }
  }
}
onMounted(() => window.addEventListener('keydown', onKeydown))
onUnmounted(() => window.removeEventListener('keydown', onKeydown))

const userEmail = computed(() => currentUser.value?.email || 'user@sentinel.com')
const userRole = computed(() => {
  const role = currentUser.value?.role || 'ENGINEER'
  const roleMap: Record<string, string> = {
    'CEO': 'Chief Executive',
    'MANAGER': 'Manager',
    'PRODUCT_OWNER': 'Product Owner',
    'PM': 'Product Owner',
    'ENGINEER': 'Engineer',
    'CHIEF_ENGINEER': 'Chief Engineer',
    'DEV': 'Engineer',
    'SUPPORT': 'Support'
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

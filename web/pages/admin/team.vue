<template>
  <div class="min-h-screen p-8">
    <!-- Access Denied for Non-CEO -->
    <div v-if="!isCEO" class="flex items-center justify-center h-[80vh]">
      <div class="text-center">
        <div class="text-6xl mb-4">🚫</div>
        <h2 class="text-3xl font-bold text-red-400 mb-2">ACCESS DENIED</h2>
        <p class="text-gray-400 mb-6">This area is restricted to CEO access only.</p>
        <button
          @click="navigateTo('/dashboard')"
          class="px-6 py-3 bg-gray-700 hover:bg-gray-600 text-white rounded font-bold transition-colors"
        >
          Return to Dashboard
        </button>
      </div>
    </div>

    <!-- CEO Team Management -->
    <div v-else>
      <!-- Header -->
      <div class="mb-6 border-b border-gray-700 pb-4">
        <h1 class="text-2xl font-bold text-white">
          👑 TEAM ROSTER & ACCESS CONTROL
        </h1>
        <p class="text-sm text-gray-400 mt-1">Manage team members and permissions</p>
      </div>

      <!-- Metrics Bar -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-8">
        <div class="bg-gray-800 border border-gray-700 rounded p-4">
          <div class="text-xs text-gray-400 uppercase tracking-wide mb-1">TOTAL HEADCOUNT</div>
          <div class="text-3xl font-bold text-white">{{ teamMembers.length }}</div>
        </div>
        <div class="bg-gray-800 border border-gray-700 rounded p-4">
          <div class="text-xs text-gray-400 uppercase tracking-wide mb-1">PROJECT MANAGERS</div>
          <div class="text-3xl font-bold text-white">{{ pmCount }}</div>
        </div>
        <div class="bg-gray-800 border border-gray-700 rounded p-4">
          <div class="text-xs text-gray-400 uppercase tracking-wide mb-1">DEVELOPERS</div>
          <div class="text-3xl font-bold text-white">{{ devCount }}</div>
        </div>
      </div>

      <!-- Loading State -->
      <div v-if="isLoading" class="flex flex-col items-center justify-center min-h-[60vh]">
        <div class="animate-spin text-6xl mb-4">⚙️</div>
        <p class="text-sm text-gray-500">กำลังโหลดทีม...</p>
      </div>

      <!-- Error State -->
      <div v-else-if="error" class="bg-red-900/20 border border-red-500 rounded-lg p-6 mb-8">
        <div class="flex items-center gap-3">
          <span class="text-3xl">⚠️</span>
          <div>
            <h3 class="text-red-400 font-bold mb-1">Failed to load team data</h3>
            <p class="text-gray-300 text-sm">{{ error }}</p>
          </div>
        </div>
      </div>

      <!-- Team Table -->
      <div v-else class="bg-gray-800 border border-gray-700 rounded overflow-hidden">
        <!-- Table Header -->
        <div class="bg-gray-900 px-6 py-3">
          <div class="grid grid-cols-12 gap-4 text-xs text-gray-400 uppercase tracking-wide">
            <div class="col-span-4">Identity</div>
            <div class="col-span-2">Joined</div>
            <div class="col-span-2">Current Rank</div>
            <div class="col-span-2">Health Score</div>
            <div class="col-span-2">Command</div>
          </div>
        </div>

        <!-- Table Body -->
        <div class="divide-y divide-gray-700">
          <div
            v-for="member in teamMembers"
            :key="member.id"
            class="px-6 py-4 hover:bg-gray-700/50 transition-colors"
          >
            <div class="grid grid-cols-12 gap-4 items-center">
              <!-- Identity (Avatar + Email) -->
              <div class="col-span-4 flex items-center gap-3">
                <div
                  :class="[
                    'w-12 h-12 rounded-full flex items-center justify-center font-bold text-lg',
                    getRoleColor(member.role).avatar
                  ]"
                >
                  {{ member.email.charAt(0).toUpperCase() }}
                </div>
                <div class="min-w-0">
                  <p class="font-medium text-white truncate">{{ member.email }}</p>
                  <p class="text-xs text-gray-400">ID: {{ member.id }}</p>
                </div>
              </div>

              <!-- Joined Date -->
              <div class="col-span-2 text-sm text-gray-300">
                {{ formatDate(member.created_at) }}
              </div>

              <!-- Current Rank (Role Badge) -->
              <div class="col-span-2">
                <span
                  :class="[
                    'inline-flex items-center gap-1 px-3 py-1 rounded-full text-xs font-bold',
                    getRoleColor(member.role).badge
                  ]"
                >
                  <span>{{ getRoleIcon(member.role) }}</span>
                  <span>{{ member.role }}</span>
                </span>
              </div>

              <!-- Health Score -->
              <div class="col-span-2">
                <div class="flex items-center gap-2">
                  <div class="flex-1 h-2 bg-gray-700 rounded-full overflow-hidden">
                    <div
                      :class="[
                        'h-full transition-all',
                        member.health_score >= 80 ? 'bg-green-500' :
                        member.health_score >= 50 ? 'bg-yellow-500' : 'bg-red-500'
                      ]"
                      :style="{ width: `${member.health_score}%` }"
                    ></div>
                  </div>
                  <span class="text-sm font-medium text-white w-12 text-right">
                    {{ member.health_score }}%
                  </span>
                </div>
              </div>

              <!-- Command (Role Change Dropdown) -->
              <div class="col-span-2">
                <div class="flex items-center gap-2">
                  <select
                    v-model="member.role"
                    @change="handleRoleChange(member)"
                    :disabled="member.id === currentUser?.user_id || changingRoleFor === member.id"
                    :class="[
                      'flex-1 px-3 py-2 bg-gray-900 border border-gray-600 rounded text-sm font-medium transition-all',
                      'focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-transparent',
                      'disabled:opacity-50 disabled:cursor-not-allowed',
                      'text-gray-300'
                    ]"
                  >
                    <option value="CEO" class="bg-gray-900">👑 CEO</option>
                    <option value="PM" class="bg-gray-900">📋 PM</option>
                    <option value="DEV" class="bg-gray-900">💻 DEV</option>
                  </select>
                  
                  <div v-if="changingRoleFor === member.id" class="text-purple-400 animate-spin">
                    ⚙️
                  </div>
                  
                  <div v-if="member.id === currentUser?.user_id" class="text-xs text-gray-500">
                    (You)
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Empty State -->
          <div v-if="teamMembers.length === 0" class="px-6 py-12 text-center">
            <div class="text-4xl mb-3">👥</div>
            <p class="text-gray-400">No team members found</p>
          </div>
        </div>
      </div>

      <!-- Success Toast -->
      <div
        v-if="successMessage"
        class="fixed bottom-8 right-8 bg-gray-800 border-2 border-purple-500 text-white px-6 py-4 rounded-lg shadow-2xl flex items-center gap-3 z-50"
      >
        <span class="text-2xl">✅</span>
        <div>
          <p class="font-bold text-purple-400">Success!</p>
          <p class="text-sm text-gray-300">{{ successMessage }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  layout: 'default',
  middleware: 'auth'
})

interface TeamMember {
  id: number
  email: string
  role: string
  health_score: number
  created_at: string
  updated_at: string
}

const { fetchWithAuth, currentUser } = useAuth()

// State
const teamMembers = ref<TeamMember[]>([])
const isLoading = ref(true)
const error = ref('')
const changingRoleFor = ref<number | null>(null)
const successMessage = ref('')
const originalRoles = ref<Map<number, string>>(new Map())

// Computed
const isCEO = computed(() => currentUser.value?.role === 'CEO')

const pmCount = computed(() => 
  teamMembers.value.filter(m => m.role === 'PM').length
)

const devCount = computed(() => 
  teamMembers.value.filter(m => m.role === 'DEV').length
)

// Access Control
onMounted(async () => {
  if (!isCEO.value) {
    await navigateTo('/dashboard')
    return
  }
  
  await fetchTeamMembers()
})

// Fetch team members
const fetchTeamMembers = async () => {
  isLoading.value = true
  error.value = ''
  
  try {
    const response = await fetchWithAuth<{ data: TeamMember[] }>('/auth/users')
    teamMembers.value = response.data || []
    
    // Store original roles for comparison
    originalRoles.value.clear()
    teamMembers.value.forEach(member => {
      originalRoles.value.set(member.id, member.role)
    })
  } catch (err: any) {
    console.error('Failed to fetch team members:', err)
    error.value = err.data?.message || err.message || 'Failed to load team data'
  } finally {
    isLoading.value = false
  }
}

// Handle role change
const handleRoleChange = async (member: TeamMember) => {
  const originalRole = originalRoles.value.get(member.id)
  
  // Check if role actually changed
  if (member.role === originalRole) {
    return
  }

  // Confirm critical change
  if (member.role === 'CEO' || originalRole === 'CEO') {
    if (!confirm(`⚠️ CRITICAL: Change ${member.email} from ${originalRole} to ${member.role}?`)) {
      // Revert to original role
      member.role = originalRole || member.role
      return
    }
  }

  changingRoleFor.value = member.id
  
  try {
    await fetchWithAuth(`/auth/users/${member.id}/role`, {
      method: 'PATCH',
      body: { role: member.role }
    })
    
    // Update original role
    originalRoles.value.set(member.id, member.role)
    
    // Show success message
    successMessage.value = `${member.email} promoted to ${member.role}`
    setTimeout(() => {
      successMessage.value = ''
    }, 3000)
  } catch (err: any) {
    console.error('Failed to change role:', err)
    
    // Revert to original role on error
    member.role = originalRole || member.role
    
    // Show error
    alert(`Failed to change role: ${err.data?.message || err.message || 'Unknown error'}`)
  } finally {
    changingRoleFor.value = null
  }
}

// Helper functions
const getRoleIcon = (role: string) => {
  const icons: Record<string, string> = {
    'CEO': '👑',
    'PM': '📋',
    'DEV': '💻'
  }
  return icons[role] || '👤'
}

const getRoleColor = (role: string) => {
  const colors: Record<string, any> = {
    'CEO': {
      badge: 'bg-purple-700 text-purple-100',
      avatar: 'bg-purple-600 text-white',
      text: 'text-purple-400'
    },
    'PM': {
      badge: 'bg-blue-700 text-blue-100',
      avatar: 'bg-blue-600 text-white',
      text: 'text-blue-400'
    },
    'DEV': {
      badge: 'bg-green-700 text-green-100',
      avatar: 'bg-green-600 text-white',
      text: 'text-green-400'
    }
  }
  return colors[role] || {
    badge: 'bg-gray-700 text-gray-100',
    avatar: 'bg-gray-600 text-white',
    text: 'text-gray-400'
  }
}

const formatDate = (dateString: string) => {
  const date = new Date(dateString)
  return date.toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  })
}
</script>

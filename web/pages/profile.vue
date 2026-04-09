<template>
  <div class="min-h-screen p-6 md:p-8">
    <!-- Page Header -->
    <header class="mb-8 border-b border-gray-700/80 pb-6">
      <div class="flex flex-wrap items-center justify-between gap-4">
        <div>
          <h1 class="text-2xl md:text-3xl font-bold text-white tracking-tight">
            Profile &amp; Account
          </h1>
          <p class="text-sm text-gray-400 mt-1">
            Manage your identity, preferences, and security settings
          </p>
        </div>
        <div class="flex items-center gap-2 text-xs text-gray-500">
          <span class="inline-flex items-center gap-1.5 rounded-full bg-gray-800 px-2.5 py-1 border border-gray-700">
            <span class="w-1.5 h-1.5 rounded-full bg-emerald-500" />
            Account active
          </span>
        </div>
      </div>
    </header>

    <!-- Loading -->
    <div v-if="isLoading" class="flex flex-col items-center justify-center min-h-[50vh]">
      <div class="animate-spin w-10 h-10 border-2 border-purple-500 border-t-transparent rounded-full" />
      <p class="text-sm text-gray-500 mt-4">Loading profile...</p>
    </div>

    <!-- Error -->
    <div v-else-if="loadError" class="rounded-xl border border-red-500/50 bg-red-950/20 p-6 mb-8">
      <div class="flex items-start gap-3">
        <span class="text-2xl">⚠️</span>
        <div>
          <h3 class="text-red-400 font-semibold">Failed to load profile</h3>
          <p class="text-gray-300 text-sm mt-1">{{ loadError }}</p>
          <button
            type="button"
            @click="loadProfile"
            class="mt-3 px-4 py-2 bg-gray-700 hover:bg-gray-600 rounded-lg text-sm font-medium transition-colors"
          >
            Retry
          </button>
        </div>
      </div>
    </div>

    <template v-else>
      <div class="grid grid-cols-1 xl:grid-cols-3 gap-8">
        <!-- Left: Identity & Tech Stack -->
        <div class="xl:col-span-2 space-y-6">
          <!-- Identity Card -->
          <section class="rounded-xl border border-gray-700/80 bg-gray-800/60 overflow-hidden">
            <div class="px-6 py-4 border-b border-gray-700/80 bg-gray-900/50">
              <h2 class="text-lg font-semibold text-white flex items-center gap-2">
                <svg class="w-5 h-5 text-purple-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                </svg>
                Identity
              </h2>
              <p class="text-xs text-gray-400 mt-0.5">Display name and account identifier</p>
            </div>
            <div class="p-6">
              <div class="flex flex-col sm:flex-row gap-6">
                <!-- Avatar Upload -->
                <div class="flex flex-col items-center gap-3 shrink-0">
                  <button
                    type="button"
                    class="relative group w-20 h-20 rounded-full focus:outline-none focus:ring-2 focus:ring-purple-500 focus:ring-offset-2 focus:ring-offset-gray-800"
                    :title="uploadingAvatar ? 'Uploading...' : 'Click to change avatar'"
                    :disabled="uploadingAvatar"
                    @click="triggerAvatarInput"
                  >
                    <img
                      v-if="profile?.avatar_url"
                      :src="profile.avatar_url"
                      alt="Avatar"
                      class="w-20 h-20 rounded-full object-cover"
                    />
                    <div
                      v-else
                      class="w-20 h-20 rounded-full flex items-center justify-center text-2xl font-bold bg-gradient-to-br from-purple-500 to-pink-600 text-white"
                    >
                      {{ avatarLetter }}
                    </div>
                    <!-- Overlay on hover -->
                    <div
                      class="absolute inset-0 rounded-full bg-black/50 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity"
                    >
                      <svg v-if="!uploadingAvatar" class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z" />
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 13a3 3 0 11-6 0 3 3 0 016 0z" />
                      </svg>
                      <div v-else class="w-5 h-5 border-2 border-white border-t-transparent rounded-full animate-spin" />
                    </div>
                  </button>
                  <input
                    ref="avatarInputRef"
                    type="file"
                    accept="image/jpeg,image/png,image/webp,image/gif"
                    class="hidden"
                    @change="handleAvatarFileChange"
                  />
                  <span class="text-xs text-gray-500">Click to change</span>
                  <p v-if="avatarError" class="text-xs text-red-400 text-center max-w-[120px]">{{ avatarError }}</p>
                </div>
                <div class="flex-1 space-y-4 min-w-0">
                  <div>
                    <label class="block text-xs font-medium text-gray-400 uppercase tracking-wider mb-1.5">Display name</label>
                    <input
                      v-model="form.display_name"
                      type="text"
                      placeholder="e.g. John Doe"
                      class="w-full px-4 py-2.5 rounded-lg bg-gray-900 border border-gray-600 text-white placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-transparent"
                      maxlength="100"
                    />
                  </div>
                  <div>
                    <label class="block text-xs font-medium text-gray-400 uppercase tracking-wider mb-1.5">Email</label>
                    <input
                      :value="profile?.email"
                      type="email"
                      disabled
                      class="w-full px-4 py-2.5 rounded-lg bg-gray-900/80 border border-gray-700 text-gray-400 cursor-not-allowed"
                    />
                    <p class="text-xs text-gray-500 mt-1">Email is managed by your administrator.</p>
                  </div>
                  <div>
                    <label class="block text-xs font-medium text-gray-400 uppercase tracking-wider mb-1.5">Role</label>
                    <span
                      :class="[
                        'inline-flex items-center px-3 py-1 rounded-full text-xs font-medium',
                        roleBadgeClass
                      ]"
                    >
                      {{ roleLabel }}
                    </span>
                  </div>
                </div>
              </div>
              <div class="mt-6 pt-4 border-t border-gray-700/80 flex justify-end">
                <button
                  type="button"
                  :disabled="savingProfile"
                  @click="saveProfile"
                  class="px-4 py-2.5 bg-gradient-to-r from-purple-100 dark:from-purple-600 to-pink-100 dark:to-pink-600 hover:from-purple-200 dark:hover:from-purple-500 hover:to-pink-200 dark:hover:to-pink-500 disabled:opacity-50 text-gray-900 dark:text-white font-medium rounded-lg transition-all"
                >
                  {{ savingProfile ? 'Saving...' : 'Save identity' }}
                </button>
              </div>
            </div>
          </section>
        </div>

        <!-- Right: Security & Account info -->
        <div class="space-y-6">
          <!-- Security Card -->
          <section class="rounded-xl border border-gray-700/80 bg-gray-800/60 overflow-hidden">
            <div class="px-6 py-4 border-b border-gray-700/80 bg-gray-900/50">
              <h2 class="text-lg font-semibold text-white flex items-center gap-2">
                <svg class="w-5 h-5 text-amber-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
                </svg>
                Security
              </h2>
              <p class="text-xs text-gray-400 mt-0.5">Change your password</p>
            </div>
            <div class="p-6 space-y-4">
              <div>
                <label class="block text-xs font-medium text-gray-400 uppercase tracking-wider mb-1.5">Current password</label>
                <input
                  v-model="passwordForm.current"
                  type="password"
                  placeholder="••••••••"
                  autocomplete="current-password"
                  class="w-full px-4 py-2.5 rounded-lg bg-gray-900 border border-gray-600 text-white placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-purple-500"
                />
              </div>
              <div>
                <label class="block text-xs font-medium text-gray-400 uppercase tracking-wider mb-1.5">New password</label>
                <input
                  v-model="passwordForm.new"
                  type="password"
                  placeholder="Min 8 characters"
                  autocomplete="new-password"
                  class="w-full px-4 py-2.5 rounded-lg bg-gray-900 border border-gray-600 text-white placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-purple-500"
                />
              </div>
              <div>
                <label class="block text-xs font-medium text-gray-400 uppercase tracking-wider mb-1.5">Confirm new password</label>
                <input
                  v-model="passwordForm.confirm"
                  type="password"
                  placeholder="••••••••"
                  autocomplete="new-password"
                  class="w-full px-4 py-2.5 rounded-lg bg-gray-900 border border-gray-600 text-white placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-purple-500"
                />
              </div>
              <p v-if="passwordError" class="text-sm text-red-400">{{ passwordError }}</p>
              <p v-if="passwordSuccess" class="text-sm text-emerald-400">Password updated successfully.</p>
              <button
                type="button"
                :disabled="savingPassword"
                @click="changePassword"
                class="w-full px-4 py-2.5 bg-amber-100 dark:bg-amber-600/80 hover:bg-amber-100 dark:bg-amber-600 disabled:opacity-50 text-gray-900 dark:text-white font-medium rounded-lg transition-colors"
              >
                {{ savingPassword ? 'Updating...' : 'Update password' }}
              </button>
            </div>
          </section>

          <!-- Account info Card -->
          <section class="rounded-xl border border-gray-700/80 bg-gray-800/60 overflow-hidden">
            <div class="px-6 py-4 border-b border-gray-700/80 bg-gray-900/50">
              <h2 class="text-lg font-semibold text-white flex items-center gap-2">
                <svg class="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                Account info
              </h2>
            </div>
            <div class="p-6 space-y-3 text-sm">
              <div class="flex justify-between">
                <span class="text-gray-400">Member since</span>
                <span class="text-gray-200">{{ formatDate(profile?.created_at) }}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-gray-400">Last updated</span>
                <span class="text-gray-200">{{ formatDate(profile?.updated_at) }}</span>
              </div>
              <div v-if="typeof profile?.health_score === 'number'" class="flex justify-between">
                <span class="text-gray-400">Health score</span>
                <span class="text-gray-200">{{ profile.health_score }}%</span>
              </div>
              <div class="flex justify-between">
                <span class="text-gray-400">User ID</span>
                <span class="text-gray-500 font-mono text-xs">{{ profile?.id }}</span>
              </div>
            </div>
          </section>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { authApi } from '~/core/modules/auth/infrastructure/auth-api'
import type { User } from '~/core/modules/auth/infrastructure/auth-api'
import { useNotification } from '~/composables/useNotification'

definePageMeta({
  layout: 'default',
  middleware: 'auth',
})

const profile = ref<User | null>(null)
const isLoading = ref(true)
const loadError = ref<string | null>(null)
const savingProfile = ref(false)
const savingPassword = ref(false)
const passwordError = ref<string | null>(null)
const passwordSuccess = ref(false)
const techInputRef = ref<HTMLInputElement | null>(null)
const techInput = ref('')
const avatarInputRef = ref<HTMLInputElement | null>(null)
const uploadingAvatar = ref(false)
const avatarError = ref<string | null>(null)

// Shared with sidebar layout
const sidebarAvatarURL = useState<string>('sidebarAvatarURL', () => '')
const { showSuccess, showError } = useNotification()

const form = reactive({
  display_name: '',
  tech_stack: [] as string[],
})

const passwordForm = reactive({
  current: '',
  new: '',
  confirm: '',
})

const roleLabels: Record<string, string> = {
  CEO: 'Chief Executive',
  MANAGER: 'Manager',
  PRODUCT_OWNER: 'Product Owner',
  PM: 'Product Owner',
  ENGINEER: 'Engineer',
  CHIEF_ENGINEER: 'Chief Engineer',
  DEV: 'Engineer',
  SUPPORT: 'Support',
}

const roleLabel = computed(() => (profile.value ? roleLabels[profile.value.role] || profile.value.role : ''))

const roleBadgeClass = computed(() => {
  const role = profile.value?.role
  if (role === 'CEO') return 'bg-amber-500/20 text-amber-400 border border-amber-500/40'
  if (role === 'PRODUCT_OWNER' || role === 'PM') return 'bg-blue-500/20 text-blue-400 border border-blue-500/40'
  if (role === 'ENGINEER' || role === 'CHIEF_ENGINEER' || role === 'DEV') return 'bg-green-500/20 text-green-400 border border-green-500/40'
  return 'bg-gray-500/20 text-gray-300 border border-gray-500/40'
})

const avatarLetter = computed(() => {
  const name = form.display_name?.trim() || profile.value?.email
  if (!name) return '?'
  return name.charAt(0).toUpperCase()
})

function formatDate(value: string | undefined) {
  if (!value) return '—'
  try {
    const d = new Date(value)
    return d.toLocaleDateString(undefined, { year: 'numeric', month: 'short', day: 'numeric' })
  } catch {
    return value
  }
}

async function loadProfile() {
  isLoading.value = true
  loadError.value = null
  try {
    const data = await authApi.getMe()
    profile.value = data
    form.display_name = data.display_name ?? ''
    form.tech_stack = Array.isArray(data.tech_stack) ? [...data.tech_stack] : []
  } catch (e: any) {
    loadError.value = e?.message ?? 'Failed to load profile'
  } finally {
    isLoading.value = false
  }
}

async function saveProfile() {
  if (!profile.value) return
  savingProfile.value = true
  try {
    const updated = await authApi.updateProfile({
      display_name: form.display_name.trim() || undefined,
      tech_stack: form.tech_stack.length ? form.tech_stack : undefined,
    })
    profile.value = updated
    if (updated.display_name !== undefined) form.display_name = updated.display_name
    if (Array.isArray(updated.tech_stack)) form.tech_stack = [...updated.tech_stack]
    sidebarAvatarURL.value = updated.avatar_url || ''
    showSuccess('Your display name and preferences have been saved.', 'Identity updated')
  } catch (e: any) {
    showError(e?.message ?? 'Failed to update profile', 'Save failed')
  } finally {
    savingProfile.value = false
  }
}

function addTech() {
  const v = techInput.value.trim()
  if (!v || form.tech_stack.length >= 20) return
  if (form.tech_stack.includes(v)) return
  form.tech_stack.push(v)
  techInput.value = ''
  techInputRef.value?.focus()
}

function removeTech(tag: string) {
  form.tech_stack = form.tech_stack.filter((t) => t !== tag)
}

async function changePassword() {
  passwordError.value = null
  passwordSuccess.value = false
  if (!passwordForm.current || !passwordForm.new || !passwordForm.confirm) {
    passwordError.value = 'Fill all password fields.'
    return
  }
  if (passwordForm.new.length < 8) {
    passwordError.value = 'New password must be at least 8 characters.'
    return
  }
  if (passwordForm.new !== passwordForm.confirm) {
    passwordError.value = 'New password and confirmation do not match.'
    return
  }
  savingPassword.value = true
  try {
    await authApi.changePassword(passwordForm.current, passwordForm.new)
    passwordForm.current = ''
    passwordForm.new = ''
    passwordForm.confirm = ''
    passwordSuccess.value = true
  } catch (e: any) {
    passwordError.value = e?.message ?? 'Failed to change password'
  } finally {
    savingPassword.value = false
  }
}

function triggerAvatarInput() {
  avatarError.value = null
  avatarInputRef.value?.click()
}

function handleAvatarFileChange(event: Event) {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  if (!input) return
  input.value = ''

  if (!file) return

  const maxBytes = 2 * 1024 * 1024
  if (file.size > maxBytes) {
    avatarError.value = 'Image too large (max 2 MB)'
    return
  }

  const reader = new FileReader()
  reader.onload = async (e) => {
    const dataURL = e.target?.result as string
    if (!dataURL) return
    uploadingAvatar.value = true
    avatarError.value = null
    try {
      const updated = await authApi.uploadAvatar(dataURL)
      profile.value = updated
      sidebarAvatarURL.value = updated.avatar_url || ''
    } catch (err: any) {
      avatarError.value = err?.message ?? 'Failed to upload avatar'
    } finally {
      uploadingAvatar.value = false
    }
  }
  reader.onerror = () => {
    avatarError.value = 'Failed to read file'
  }
  reader.readAsDataURL(file)
}

onMounted(() => {
  loadProfile()
})
</script>

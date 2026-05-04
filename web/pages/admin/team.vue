<template>
  <div class="min-h-screen bg-gray-900 p-8">
    <!-- Access Denied for Non-CEO -->
    <div v-if="!isCEO" class="flex items-center justify-center h-[80vh]">
      <div class="text-center">
        <div class="text-6xl mb-4">🚫</div>
        <h2 class="text-3xl font-bold text-red-400 mb-2">ACCESS DENIED</h2>
        <p class="text-gray-400 mb-6">This area is restricted to CEO access only.</p>
        <button
          @click="navigateTo('/dashboard')"
          class="px-6 py-3 bg-gray-700 hover:bg-gray-600 text-gray-900 dark:text-white rounded font-bold transition-colors"
        >
          Return to Dashboard
        </button>
      </div>
    </div>

    <!-- CEO Team Management -->
    <div v-else>
      <!-- Header + Actions -->
      <div class="mb-6 border-b border-gray-700 pb-4 flex flex-wrap items-end justify-between gap-4">
        <div>
          <h1 class="text-2xl font-bold text-white">
            หน้าจัดการ User
          </h1>
          <p class="text-sm text-gray-400 mt-1">จัดการสมาชิกและสิทธิ์การใช้งาน</p>
        </div>
        <div class="flex items-center gap-3">
          <button
            type="button"
            @click="openImportPanel"
            class="inline-flex items-center gap-2 px-4 py-2.5 bg-gray-700 hover:bg-gray-600 border border-gray-600 text-gray-200 font-medium rounded-lg transition-colors focus:outline-none focus:ring-2 focus:ring-purple-500 focus:ring-offset-2 focus:ring-offset-gray-900"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12" /></svg>
            Import users
          </button>
          <button
            type="button"
            @click="showCreateModal = true"
            class="inline-flex items-center gap-2 px-4 py-2.5 bg-gradient-to-r from-purple-100 dark:from-purple-600 to-pink-100 dark:to-pink-600 hover:from-purple-200 dark:hover:from-purple-500 hover:to-pink-200 dark:hover:to-pink-500 text-gray-900 dark:text-white font-medium rounded-lg shadow-lg transition-all focus:outline-none focus:ring-2 focus:ring-purple-400 focus:ring-offset-2 focus:ring-offset-gray-900"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" /></svg>
            Add user
          </button>
        </div>
      </div>

      <!-- Metrics Bar -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-8">
        <div class="bg-gray-800 border border-gray-700 rounded p-4">
          <div class="text-xs text-gray-400 uppercase tracking-wide mb-1">TOTAL HEADCOUNT</div>
          <div class="text-3xl font-bold text-white">{{ teamMembers.length }}</div>
        </div>
        <div class="bg-gray-800 border border-gray-700 rounded p-4">
          <div class="text-xs text-gray-400 uppercase tracking-wide mb-1">PRODUCT OWNERS</div>
          <div class="text-3xl font-bold text-white">{{ pmCount }}</div>
        </div>
        <div class="bg-gray-800 border border-gray-700 rounded p-4">
          <div class="text-xs text-gray-400 uppercase tracking-wide mb-1">ENGINEERS</div>
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
            <div class="col-span-4">Command</div>
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
                    getRoleColor(normalizeUserRole(member.role)).avatar
                  ]"
                >
                  {{ getMemberAvatarInitial(member) }}
                </div>
                <div class="min-w-0">
                  <p class="font-medium text-white truncate">{{ getMemberDisplayName(member) }}</p>
                  <p class="text-xs text-gray-400 truncate">{{ member.email }}</p>
                  <p class="text-xs text-gray-500">ID: {{ member.id }}</p>
                </div>
              </div>

              <!-- Joined Date -->
              <div class="col-span-2 text-sm text-gray-300">
                {{ formatDate(member.created_at) }}
              </div>

              <!-- Current Rank (Role Badge) -->
              <div class="col-span-2">
                <div class="flex flex-col gap-1">
                  <span
                    :class="[
                      'inline-flex items-center gap-1 px-3 py-1 rounded-full text-xs font-bold',
                      getRoleColor(normalizeUserRole(member.role)).badge
                    ]"
                  >
                    <span>{{ getRoleIcon(normalizeUserRole(member.role)) }}</span>
                    <span>{{ formatRoleLabel(member.role) }}</span>
                  </span>
                  <span
                    v-if="member.is_remote"
                    class="inline-flex items-center gap-1 px-2 py-0.5 rounded text-xs font-medium bg-emerald-900/40 text-emerald-300 border border-emerald-700/50"
                  >
                    🌐 Remote
                  </span>
                </div>
              </div>

              <!-- Command (Role Change + Delete) -->
              <div class="col-span-4">
                <div class="flex items-center gap-2">
                  <select
                    v-model="member.role"
                    @change="handleRoleChange(member)"
                    :disabled="member.id === currentUser?.user_id || changingRoleFor === member.id"
                    :class="[
                      'flex-1 min-w-0 px-3 py-2 bg-gray-900 border border-gray-600 rounded text-sm font-medium transition-all',
                      'focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-transparent',
                      'disabled:opacity-50 disabled:cursor-not-allowed',
                      'text-gray-300'
                    ]"
                  >
                    <option value="CEO" class="bg-gray-900">👑 CEO</option>
                    <option value="MANAGER" class="bg-gray-900">🏢 MANAGER</option>
                    <option value="PRODUCT_OWNER" class="bg-gray-900">📋 Product Owner</option>
                    <option value="ENGINEER" class="bg-gray-900">💻 ENGINEER</option>
                    <option value="CHIEF_ENGINEER" class="bg-gray-900">⚙️ CHIEF ENGINEER</option>
                    <option value="SUPPORT" class="bg-gray-900">🎧 SUPPORT</option>
                  </select>
                  <button
                    v-if="member.id !== currentUser?.user_id && !deletingUserId && !resettingPasswordUserId"
                    type="button"
                    @click="handleRemoteToggle(member)"
                    :class="[
                      'p-2 rounded transition-colors',
                      member.is_remote
                        ? 'text-emerald-400 hover:bg-emerald-900/30'
                        : 'text-gray-400 hover:text-emerald-400 hover:bg-emerald-900/30'
                    ]"
                    :title="member.is_remote ? 'Unset remote worker' : 'Set as remote worker'"
                    :aria-label="member.is_remote ? 'Unset remote worker' : 'Set as remote worker'"
                    :disabled="changingRemoteFor === member.id"
                  >
                    <svg v-if="changingRemoteFor === member.id" class="w-5 h-5 animate-spin" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" /></svg>
                    <svg v-else class="w-5 h-5" fill="currentColor" viewBox="0 0 24 24"><path d="M21.7214 12.7517C21.7404 12.5036 21.75 12.2529 21.75 11.9999C21.75 10.4758 21.4003 9.03328 20.7767 7.74835C19.5396 8.92269 18.0671 9.85146 16.4374 10.4565C16.4789 10.9655 16.5 11.4803 16.5 11.9999C16.5 13.1011 16.4051 14.1802 16.2229 15.2293C18.2163 14.7277 20.0717 13.8792 21.7214 12.7517Z" /><path d="M14.6343 15.5501C14.874 14.4043 15 13.2168 15 11.9999C15 11.6315 14.9885 11.2659 14.9657 10.9032C14.0141 11.1299 13.021 11.2499 12 11.2499C10.979 11.2499 9.98594 11.1299 9.0343 10.9032C9.01155 11.2659 9 11.6315 9 11.9999C9 13.2168 9.12601 14.4043 9.3657 15.5501C10.2246 15.6817 11.1043 15.7499 12 15.7499C12.8957 15.7499 13.7754 15.6817 14.6343 15.5501Z" /><path d="M9.77224 17.119C10.5028 17.2054 11.2462 17.2499 12 17.2499C12.7538 17.2499 13.4972 17.2054 14.2278 17.119C13.714 18.7746 12.9575 20.3235 12 21.724C11.0425 20.3235 10.286 18.7746 9.77224 17.119Z" /><path d="M7.77705 15.2293C7.59493 14.1802 7.5 13.1011 7.5 11.9999C7.5 11.4803 7.52114 10.9655 7.56261 10.4565C5.93286 9.85146 4.46039 8.92269 3.22333 7.74835C2.59973 9.03328 2.25 10.4758 2.25 11.9999C2.25 12.2529 2.25964 12.5036 2.27856 12.7517C3.92826 13.8792 5.78374 14.7277 7.77705 15.2293Z" /><path d="M21.3561 14.7525C20.3404 18.2104 17.4597 20.8705 13.8776 21.5693C14.744 20.1123 15.4185 18.5278 15.8664 16.8508C17.8263 16.44 19.6736 15.7231 21.3561 14.7525Z" /><path d="M2.64395 14.7525C4.32642 15.7231 6.17372 16.44 8.13356 16.8508C8.58146 18.5278 9.25602 20.1123 10.1224 21.5693C6.54027 20.8705 3.65964 18.2104 2.64395 14.7525Z" /><path d="M13.8776 2.43055C16.3991 2.92245 18.5731 4.3862 19.9937 6.41599C18.9351 7.48484 17.6637 8.34251 16.2483 8.92017C15.862 6.58282 15.0435 4.39132 13.8776 2.43055Z" /><path d="M12 2.27588C13.4287 4.36548 14.4097 6.78537 14.805 9.39744C13.9083 9.62756 12.9684 9.74993 12 9.74993C11.0316 9.74993 10.0917 9.62756 9.19503 9.39744C9.5903 6.78537 10.5713 4.36548 12 2.27588Z" /><path d="M10.1224 2.43055C8.95648 4.39132 8.13795 6.58282 7.75171 8.92017C6.33629 8.34251 5.06489 7.48484 4.00635 6.41599C5.42689 4.3862 7.60085 2.92245 10.1224 2.43055Z" /></svg>
                  </button>
                  <button
                    v-if="member.id !== currentUser?.user_id && !deletingUserId && !resettingPasswordUserId"
                    type="button"
                    @click="openEditModal(member)"
                    class="p-2 text-gray-400 hover:text-blue-400 hover:bg-blue-900/30 rounded transition-colors"
                    :title="`Edit ${member.email}`"
                    aria-label="Edit user"
                  >
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" /></svg>
                  </button>
                  <button
                    v-if="member.id !== currentUser?.user_id && !deletingUserId && !resettingPasswordUserId"
                    type="button"
                    @click="openResetPasswordModal(member)"
                    class="p-2 text-gray-400 hover:text-amber-400 hover:bg-amber-900/30 rounded transition-colors"
                    :title="`Reset password for ${member.email}`"
                    aria-label="Reset password"
                  >
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z" /></svg>
                  </button>
                  <button
                    v-if="member.id !== currentUser?.user_id && !deletingUserId && !resettingPasswordUserId"
                    type="button"
                    @click="memberToDelete = member"
                    class="p-2 text-gray-400 hover:text-red-400 hover:bg-red-900/30 rounded transition-colors"
                    :title="`Remove ${member.email}`"
                    aria-label="Delete user"
                  >
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" /></svg>
                  </button>
                  <div v-if="resettingPasswordUserId === member.id" class="p-2 text-amber-400 animate-spin">
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" /></svg>
                  </div>
                  <div v-if="deletingUserId === member.id" class="p-2 text-red-400 animate-spin">
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" /></svg>
                  </div>
                  <div v-if="changingRoleFor === member.id && deletingUserId !== member.id && resettingPasswordUserId !== member.id" class="text-purple-400 animate-spin">
                    ⚙️
                  </div>
                  <div v-if="member.id === currentUser?.user_id" class="text-xs text-gray-500 shrink-0">
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

      <!-- Create User Modal -->
      <Teleport to="body">
        <div
          v-if="showCreateModal"
          class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-sm"
          role="dialog"
          aria-modal="true"
          aria-labelledby="create-user-title"
          @keydown.escape="showCreateModal = false"
          @click.self="showCreateModal = false"
        >
          <div
            class="w-full max-w-md bg-gray-800 border border-gray-700 rounded-xl shadow-2xl"
            @click.stop
          >
            <div class="px-6 py-4 border-b border-gray-700 flex items-center justify-between">
              <h2 id="create-user-title" class="text-lg font-bold text-white">Add user</h2>
              <button
                type="button"
                @click="showCreateModal = false"
                class="p-2 text-gray-400 hover:text-gray-900 dark:text-white rounded-lg hover:bg-gray-700 transition-colors"
                aria-label="Close"
              >
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" /></svg>
              </button>
            </div>
            <form @submit.prevent="submitCreateUser" class="p-6 space-y-4">
              <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
                <div>
                  <label for="create-first-name" class="block text-sm font-medium text-gray-300 mb-1">First name</label>
                  <input
                    id="create-first-name"
                    v-model="createForm.firstName"
                    type="text"
                    placeholder="John"
                    class="w-full px-4 py-2.5 bg-gray-900 border border-gray-600 rounded-lg text-white placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-transparent"
                  />
                </div>
                <div>
                  <label for="create-last-name" class="block text-sm font-medium text-gray-300 mb-1">Last name</label>
                  <input
                    id="create-last-name"
                    v-model="createForm.lastName"
                    type="text"
                    placeholder="Doe"
                    class="w-full px-4 py-2.5 bg-gray-900 border border-gray-600 rounded-lg text-white placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-transparent"
                  />
                </div>
              </div>
              <div>
                <label for="create-email" class="block text-sm font-medium text-gray-300 mb-1">Email</label>
                <input
                  id="create-email"
                  v-model="createForm.email"
                  type="email"
                  required
                  placeholder="user@company.com"
                  class="w-full px-4 py-2.5 bg-gray-900 border border-gray-600 rounded-lg text-white placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-transparent"
                />
              </div>
              <div>
                <label for="create-password" class="block text-sm font-medium text-gray-300 mb-1">Password</label>
                <input
                  id="create-password"
                  v-model="createForm.password"
                  type="password"
                  required
                  minlength="8"
                  placeholder="Min 8 characters"
                  class="w-full px-4 py-2.5 bg-gray-900 border border-gray-600 rounded-lg text-white placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-transparent"
                />
              </div>
              <div>
                <label for="create-confirm" class="block text-sm font-medium text-gray-300 mb-1">Confirm password</label>
                <input
                  id="create-confirm"
                  v-model="createForm.confirmPassword"
                  type="password"
                  required
                  minlength="8"
                  placeholder="Repeat password"
                  class="w-full px-4 py-2.5 bg-gray-900 border border-gray-600 rounded-lg text-white placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-transparent"
                />
                <p v-if="createForm.password && createForm.confirmPassword && createForm.password !== createForm.confirmPassword" class="mt-1 text-sm text-red-400">Passwords do not match</p>
              </div>
              <div>
                <label for="create-display-name" class="block text-sm font-medium text-gray-300 mb-1">Display name</label>
                <input
                  id="create-display-name"
                  v-model="createForm.displayName"
                  type="text"
                  placeholder="Optional"
                  class="w-full px-4 py-2.5 bg-gray-900 border border-gray-600 rounded-lg text-white placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-transparent"
                />
              </div>
              <div>
                <label for="create-role" class="block text-sm font-medium text-gray-300 mb-1">Role</label>
                <select
                  id="create-role"
                  v-model="createForm.role"
                  class="w-full px-4 py-2.5 bg-gray-900 border border-gray-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-transparent"
                >
                  <option value="ENGINEER">💻 ENGINEER</option>
                  <option value="CHIEF_ENGINEER">⚙️ CHIEF ENGINEER</option>
                  <option value="PRODUCT_OWNER">📋 Product Owner</option>
                  <option value="SUPPORT">🎧 SUPPORT</option>
                  <option value="MANAGER">🏢 MANAGER</option>
                  <option value="CEO">👑 CEO</option>
                </select>
              </div>
              <div v-if="createError" class="p-3 bg-red-900/30 border border-red-500/50 rounded-lg text-sm text-red-300">
                {{ createError }}
              </div>
              <div class="flex gap-3 pt-2">
                <button
                  type="button"
                  @click="showCreateModal = false"
                  class="flex-1 px-4 py-2.5 bg-gray-700 hover:bg-gray-600 text-gray-200 font-medium rounded-lg transition-colors"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  :disabled="createSubmitting || (createForm.password !== createForm.confirmPassword)"
                  class="flex-1 px-4 py-2.5 bg-gradient-to-r from-purple-100 dark:from-purple-600 to-pink-100 dark:to-pink-600 hover:from-purple-200 dark:hover:from-purple-500 hover:to-pink-200 dark:hover:to-pink-500 disabled:opacity-50 disabled:cursor-not-allowed text-gray-900 dark:text-white font-medium rounded-lg transition-colors"
                >
                  {{ createSubmitting ? 'Creating…' : 'Create user' }}
                </button>
              </div>
            </form>
          </div>
        </div>
      </Teleport>

      <!-- Edit User Modal -->
      <Teleport to="body">
        <div
          v-if="memberToEdit"
          class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-sm"
          role="dialog"
          aria-modal="true"
          aria-labelledby="edit-user-title"
          @keydown.escape="memberToEdit = null"
          @click.self="memberToEdit = null"
        >
          <div
            class="w-full max-w-md bg-gray-800 border border-gray-700 rounded-xl shadow-2xl"
            @click.stop
          >
            <div class="px-6 py-4 border-b border-gray-700 flex items-center justify-between">
              <h2 id="edit-user-title" class="text-lg font-bold text-white">Edit user</h2>
              <button
                type="button"
                @click="memberToEdit = null"
                class="p-2 text-gray-400 hover:text-gray-900 dark:text-white rounded-lg hover:bg-gray-700 transition-colors"
                aria-label="Close"
              >
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" /></svg>
              </button>
            </div>
            <form @submit.prevent="submitEditUser" class="p-6 space-y-4">
              <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
                <div>
                  <label for="edit-first-name" class="block text-sm font-medium text-gray-300 mb-1">First name</label>
                  <input
                    id="edit-first-name"
                    v-model="editForm.firstName"
                    type="text"
                    placeholder="John"
                    class="w-full px-4 py-2.5 bg-gray-900 border border-gray-600 rounded-lg text-white placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  />
                </div>
                <div>
                  <label for="edit-last-name" class="block text-sm font-medium text-gray-300 mb-1">Last name</label>
                  <input
                    id="edit-last-name"
                    v-model="editForm.lastName"
                    type="text"
                    placeholder="Doe"
                    class="w-full px-4 py-2.5 bg-gray-900 border border-gray-600 rounded-lg text-white placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  />
                </div>
              </div>
              <div>
                <label for="edit-display-name" class="block text-sm font-medium text-gray-300 mb-1">Display name</label>
                <input
                  id="edit-display-name"
                  v-model="editForm.displayName"
                  type="text"
                  placeholder="Optional"
                  class="w-full px-4 py-2.5 bg-gray-900 border border-gray-600 rounded-lg text-white placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                />
              </div>
              <div>
                <label for="edit-email" class="block text-sm font-medium text-gray-300 mb-1">Email</label>
                <input
                  id="edit-email"
                  v-model="editForm.email"
                  type="email"
                  required
                  placeholder="user@company.com"
                  class="w-full px-4 py-2.5 bg-gray-900 border border-gray-600 rounded-lg text-white placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                />
              </div>
              <div v-if="editError" class="p-3 bg-red-900/30 border border-red-500/50 rounded-lg text-sm text-red-300">
                {{ editError }}
              </div>
              <div class="flex gap-3 pt-2">
                <button
                  type="button"
                  @click="memberToEdit = null"
                  class="flex-1 px-4 py-2.5 bg-gray-700 hover:bg-gray-600 text-gray-200 font-medium rounded-lg transition-colors"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  :disabled="editSubmitting"
                  class="flex-1 px-4 py-2.5 bg-blue-600 hover:bg-blue-500 disabled:opacity-50 disabled:cursor-not-allowed text-white font-medium rounded-lg transition-colors"
                >
                  {{ editSubmitting ? 'Saving…' : 'Save changes' }}
                </button>
              </div>
            </form>
          </div>
        </div>
      </Teleport>

      <!-- Import Users Panel (slide-over) -->
      <Teleport to="body">
        <div
          v-if="showImportPanel"
          class="fixed inset-0 z-50 flex justify-end bg-black/60 backdrop-blur-sm"
          role="dialog"
          aria-modal="true"
          aria-labelledby="import-users-title"
          @keydown.escape="closeImportPanel"
          @click.self="closeImportPanel"
        >
          <div
            class="w-full max-w-2xl h-full bg-gray-800 border-l border-gray-700 shadow-2xl flex flex-col overflow-hidden"
            @click.stop
          >
            <div class="px-6 py-4 border-b border-gray-700 flex items-center justify-between shrink-0">
              <h2 id="import-users-title" class="text-lg font-bold text-white">Import users</h2>
              <button
                type="button"
                @click="closeImportPanel"
                class="p-2 text-gray-400 hover:text-gray-900 dark:text-white rounded-lg hover:bg-gray-700 transition-colors"
                aria-label="Close"
              >
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" /></svg>
              </button>
            </div>

            <div v-if="!importResult" class="flex-1 flex flex-col overflow-hidden p-6">
              <p class="text-sm text-gray-400 mb-3">
                One user per line. Format: <code class="px-1.5 py-0.5 bg-gray-700 rounded text-gray-300">email</code> or
                <code class="px-1.5 py-0.5 bg-gray-700 rounded text-gray-300">email,role</code> or
                <code class="px-1.5 py-0.5 bg-gray-700 rounded text-gray-300">email,role,password</code>. Role defaults to ENGINEER (legacy CSV may use DEV or PM; mapped to ENGINEER / PRODUCT_OWNER). Use CHIEF_ENGINEER for chief engineers. Max 500.
              </p>
              <textarea
                v-model="importRaw"
                rows="12"
                placeholder="dev1@company.com&#10;po@company.com,PRODUCT_OWNER&#10;ceo@company.com,CEO,SecurePass123"
                class="w-full px-4 py-3 bg-gray-900 border border-gray-600 rounded-lg text-white placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-purple-500 font-mono text-sm resize-y min-h-[200px]"
              />
              <div class="mt-3 flex items-center justify-between gap-4">
                <p class="text-sm text-gray-500">
                  <span v-if="parsedImport.length > 0">{{ parsedImport.length }} user(s) ready</span>
                  <span v-else>Paste or type lines above</span>
                </p>
                <div class="flex gap-3">
                  <button
                    type="button"
                    @click="closeImportPanel"
                    class="px-4 py-2 bg-gray-700 hover:bg-gray-600 text-gray-200 rounded-lg font-medium transition-colors"
                  >
                    Cancel
                  </button>
                  <button
                    type="button"
                    @click="submitImport"
                    :disabled="importSubmitting || parsedImport.length === 0 || parsedImport.length > 500"
                    class="px-4 py-2 bg-gradient-to-r from-purple-600 to-pink-600 hover:from-purple-500 hover:to-pink-500 disabled:opacity-50 disabled:cursor-not-allowed text-white font-medium rounded-lg transition-colors"
                  >
                    {{ importSubmitting ? 'Importing…' : 'Import' }}
                  </button>
                </div>
              </div>
              <p v-if="importError" class="mt-3 text-sm text-red-400">{{ importError }}</p>
            </div>

            <!-- Import results -->
            <div v-else class="flex-1 flex flex-col overflow-hidden p-6">
              <div class="grid grid-cols-2 sm:grid-cols-4 gap-3 mb-4">
                <div class="bg-gray-900 rounded-lg p-3 text-center">
                  <div class="text-2xl font-bold text-white">{{ importResult.created }}</div>
                  <div class="text-xs text-gray-400 uppercase">Created</div>
                </div>
                <div class="bg-gray-900 rounded-lg p-3 text-center">
                  <div class="text-2xl font-bold text-yellow-400">{{ importResult.skipped }}</div>
                  <div class="text-xs text-gray-400 uppercase">Skipped</div>
                </div>
                <div class="bg-gray-900 rounded-lg p-3 text-center">
                  <div class="text-2xl font-bold text-red-400">{{ importResult.errors }}</div>
                  <div class="text-xs text-gray-400 uppercase">Errors</div>
                </div>
                <div class="bg-gray-900 rounded-lg p-3 text-center">
                  <div class="text-2xl font-bold text-gray-300">{{ importResult.total }}</div>
                  <div class="text-xs text-gray-400 uppercase">Total</div>
                </div>
              </div>
              <div class="flex-1 overflow-auto border border-gray-700 rounded-lg">
                <table class="w-full text-sm">
                  <thead class="bg-gray-900 sticky top-0">
                    <tr class="text-left text-gray-400 uppercase tracking-wide">
                      <th class="px-4 py-2">Email</th>
                      <th class="px-4 py-2">Status</th>
                      <th class="px-4 py-2">Message</th>
                      <th class="px-4 py-2">Temp password</th>
                    </tr>
                  </thead>
                  <tbody class="divide-y divide-gray-700">
                    <tr
                      v-for="(row, i) in importResult.results"
                      :key="i"
                      :class="{
                        'text-green-400': row.status === 'created',
                        'text-yellow-400': row.status === 'skipped',
                        'text-red-400': row.status === 'error'
                      }"
                    >
                      <td class="px-4 py-2 text-white">{{ row.email }}</td>
                      <td class="px-4 py-2 font-medium">{{ row.status }}</td>
                      <td class="px-4 py-2 text-gray-400">{{ row.message || '—' }}</td>
                      <td class="px-4 py-2">
                        <span v-if="row.temp_password" class="font-mono text-gray-300">{{ row.temp_password }}</span>
                        <button
                          v-if="row.temp_password"
                          type="button"
                          @click="copyTempPassword(row.temp_password)"
                          class="ml-2 px-2 py-1 text-xs bg-gray-700 hover:bg-gray-600 rounded text-gray-300"
                        >
                          Copy
                        </button>
                        <span v-else class="text-gray-500">—</span>
                      </td>
                    </tr>
                  </tbody>
                </table>
              </div>
              <div class="mt-4 flex flex-wrap items-center justify-between gap-3 shrink-0">
                <button
                  type="button"
                  @click="exportImportResultsToExcel"
                  class="inline-flex items-center gap-2 px-4 py-2 bg-emerald-200 dark:bg-emerald-700 hover:bg-emerald-100 dark:bg-emerald-600 text-gray-900 dark:text-white rounded-lg font-medium transition-colors"
                >
                  <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" /></svg>
                  Export to Excel
                </button>
                <div class="flex gap-3">
                  <button
                    type="button"
                    @click="closeImportPanel"
                    class="px-4 py-2 bg-gray-700 hover:bg-gray-600 text-gray-200 rounded-lg font-medium transition-colors"
                  >
                    Done
                  </button>
                  <button
                    type="button"
                    @click="importResult = null; importRaw = ''"
                    class="px-4 py-2 bg-gradient-to-r from-purple-100 dark:from-purple-600 to-pink-100 dark:to-pink-600 text-gray-900 dark:text-white rounded-lg font-medium transition-colors"
                  >
                    Import more
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </Teleport>

      <!-- Delete user confirmation modal -->
      <Teleport to="body">
        <div
          v-if="memberToDelete"
          class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-sm"
          role="dialog"
          aria-modal="true"
          aria-labelledby="delete-confirm-title"
          @keydown.escape="memberToDelete = null"
          @click.self="memberToDelete = null"
        >
          <div
            class="w-full max-w-md bg-gray-800 border border-gray-700 rounded-xl shadow-2xl overflow-hidden"
            @click.stop
          >
            <div class="px-6 py-4 border-b border-gray-700 flex items-center gap-3">
              <div class="flex-shrink-0 w-12 h-12 rounded-full bg-red-900/50 flex items-center justify-center">
                <svg class="w-6 h-6 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" /></svg>
              </div>
              <div>
                <h2 id="delete-confirm-title" class="text-lg font-bold text-white">ยืนยันการลบผู้ใช้</h2>
                <p class="text-sm text-gray-400 mt-0.5">การดำเนินการนี้ไม่สามารถย้อนกลับได้</p>
              </div>
            </div>
            <div class="px-6 py-4">
              <p class="text-gray-300">
                คุณต้องการลบผู้ใช้
                <strong class="text-white">{{ memberToDelete?.email }}</strong>
                ออกจากทีมใช่หรือไม่?
              </p>
            </div>
            <div class="px-6 py-4 bg-gray-900/50 flex gap-3 justify-end">
              <button
                type="button"
                @click="memberToDelete = null"
                class="px-4 py-2.5 bg-gray-700 hover:bg-gray-600 text-gray-200 font-medium rounded-lg transition-colors"
              >
                ยกเลิก
              </button>
              <button
                type="button"
                @click="executeDeleteMember"
                :disabled="deletingUserId !== null"
                class="px-4 py-2.5 bg-red-100 dark:bg-red-600 hover:bg-red-100 dark:bg-red-500 disabled:opacity-50 disabled:cursor-not-allowed text-gray-900 dark:text-white font-medium rounded-lg transition-colors"
              >
                {{ deletingUserId !== null ? 'กำลังลบ…' : 'ลบผู้ใช้' }}
              </button>
            </div>
          </div>
        </div>
      </Teleport>

      <!-- Reset password modal -->
      <Teleport to="body">
        <div
          v-if="memberToReset"
          class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-sm"
          role="dialog"
          aria-modal="true"
          aria-labelledby="reset-password-title"
          @keydown.escape="closeResetPasswordModal"
          @click.self="closeResetPasswordModal"
        >
          <div
            class="w-full max-w-md bg-gray-800 border border-gray-700 rounded-xl shadow-2xl overflow-hidden"
            @click.stop
          >
            <div class="px-6 py-4 border-b border-gray-700 flex items-center gap-3">
              <div class="flex-shrink-0 w-12 h-12 rounded-full bg-amber-900/50 flex items-center justify-center">
                <svg class="w-6 h-6 text-amber-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z" /></svg>
              </div>
              <div>
                <h2 id="reset-password-title" class="text-lg font-bold text-white">รีเซ็ตรหัสผ่าน</h2>
                <p class="text-sm text-gray-400 mt-0.5">{{ memberToReset?.email }}</p>
              </div>
            </div>
            <div v-if="!resetPasswordResult" class="p-6 space-y-4">
              <p class="text-gray-300">ระบบจะสร้างรหัสผ่านชั่วคราว (12 ตัวอักษร) ให้อัตโนมัติ แจ้งให้ผู้ใช้เปลี่ยนรหัสหลังล็อกอินได้</p>
              <div v-if="resetPasswordError" class="p-3 bg-red-900/30 border border-red-500/50 rounded-lg text-sm text-red-300">
                {{ resetPasswordError }}
              </div>
              <div class="flex gap-3 pt-2">
                <button
                  type="button"
                  @click="closeResetPasswordModal"
                  class="flex-1 px-4 py-2.5 bg-gray-700 hover:bg-gray-600 text-gray-200 font-medium rounded-lg transition-colors"
                >
                  ยกเลิก
                </button>
                <button
                  type="button"
                  @click="submitResetPassword"
                  :disabled="resetPasswordSubmitting"
                  class="flex-1 px-4 py-2.5 bg-amber-100 dark:bg-amber-600 hover:bg-amber-100 dark:bg-amber-500 disabled:opacity-50 disabled:cursor-not-allowed text-gray-900 dark:text-white font-medium rounded-lg transition-colors"
                >
                  {{ resetPasswordSubmitting ? 'กำลังสร้างรหัส…' : 'รีเซ็ตรหัสผ่าน' }}
                </button>
              </div>
            </div>
            <div v-else class="p-6 space-y-4">
              <p class="text-sm text-gray-400">รหัสผ่านชั่วคราว (แจ้งผู้ใช้แล้ว Copy เก็บหรือส่งต่อ)</p>
              <div class="flex items-center gap-2">
                <code class="flex-1 px-4 py-2.5 bg-gray-900 border border-gray-600 rounded-lg text-white font-mono text-lg">
                  {{ resetPasswordResult.temp_password }}
                </code>
                <button
                  type="button"
                  @click="copyTempPassword(resetPasswordResult.temp_password)"
                  class="px-4 py-2.5 bg-amber-100 dark:bg-amber-600 hover:bg-amber-100 dark:bg-amber-500 text-gray-900 dark:text-white font-medium rounded-lg transition-colors"
                >
                  Copy
                </button>
              </div>
              <div class="flex justify-end pt-2">
                <button
                  type="button"
                  @click="closeResetPasswordModal"
                  class="px-4 py-2.5 bg-gray-700 hover:bg-gray-600 text-gray-200 font-medium rounded-lg transition-colors"
                >
                  เสร็จสิ้น
                </button>
              </div>
            </div>
          </div>
        </div>
      </Teleport>

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
import { isEngineerLikeRole } from '~/utils/roles'

definePageMeta({
  layout: 'default',
  middleware: 'auth'
})

interface TeamMember {
  id: number
  email: string
  first_name?: string
  last_name?: string
  display_name?: string
  role: string
  health_score: number
  created_at: string
  updated_at: string
  is_remote?: boolean
}

/** Legacy API/DB may use DEV or PM; canonical roles ENGINEER / PRODUCT_OWNER. */
function normalizeUserRole(role: string): string {
  if (role === 'DEV') return 'ENGINEER'
  if (role === 'PM') return 'PRODUCT_OWNER'
  return role
}

function formatRoleLabel(role: string): string {
  const n = normalizeUserRole(role)
  const labels: Record<string, string> = {
    CEO: 'CEO',
    MANAGER: 'MANAGER',
    PRODUCT_OWNER: 'Product Owner',
    PM: 'Product Owner',
    ENGINEER: 'ENGINEER',
    CHIEF_ENGINEER: 'Chief Engineer',
    SUPPORT: 'SUPPORT',
  }
  return labels[n] ?? n
}

const { fetchWithAuth, currentUser } = useAuth()
const { showError, showSuccess, confirm } = useNotification()

// State
const teamMembers = ref<TeamMember[]>([])
const isLoading = ref(true)
const error = ref('')
const changingRoleFor = ref<number | null>(null)
const successMessage = ref('')
const originalRoles = ref<Map<number, string>>(new Map())
const deletingUserId = ref<number | null>(null)
const memberToDelete = ref<TeamMember | null>(null)
const memberToReset = ref<TeamMember | null>(null)
const resetPasswordResult = ref<{ temp_password: string } | null>(null)
const resetPasswordError = ref('')
const resetPasswordSubmitting = ref(false)
const resettingPasswordUserId = ref<number | null>(null)
const changingRemoteFor = ref<number | null>(null)

// Create user modal
const showCreateModal = ref(false)
const createForm = ref({
  firstName: '',
  lastName: '',
  displayName: '',
  email: '',
  password: '',
  confirmPassword: '',
  role: 'ENGINEER'
})
const createError = ref('')
const createSubmitting = ref(false)

// Edit user modal
const memberToEdit = ref<TeamMember | null>(null)
const editForm = ref({
  firstName: '',
  lastName: '',
  displayName: '',
  email: ''
})
const editError = ref('')
const editSubmitting = ref(false)

// Import panel
const showImportPanel = ref(false)
const importRaw = ref('')
const importSubmitting = ref(false)
const importError = ref('')
const importResult = ref<{
  total: number
  created: number
  skipped: number
  errors: number
  results: Array<{ email: string; status: string; message?: string; temp_password?: string }>
} | null>(null)

// Computed
const isCEO = computed(() => currentUser.value?.role === 'CEO' || currentUser.value?.role === 'MANAGER')

const pmCount = computed(() =>
  teamMembers.value.filter(m => m.role === 'PRODUCT_OWNER' || m.role === 'PM').length
)

const devCount = computed(() =>
  teamMembers.value.filter(m => isEngineerLikeRole(m.role)).length
)

// Parse import text into users array (email, role?, password?)
const parsedImport = computed(() => {
  const lines = importRaw.value
    .split('\n')
    .map(l => l.trim())
    .filter(Boolean)
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  const roles = ['CEO', 'MANAGER', 'PRODUCT_OWNER', 'ENGINEER', 'CHIEF_ENGINEER', 'SUPPORT']
  return lines.map(line => {
    const parts = line.split(',').map(p => p.trim())
    const email = parts[0] || ''
    const rawRole = (parts[1] || '').toUpperCase()
    let normalized = rawRole === 'DEV' ? 'ENGINEER' : rawRole
    if (normalized === 'PM') normalized = 'PRODUCT_OWNER'
    const role = normalized && roles.includes(normalized) ? normalized : 'ENGINEER'
    const password = parts[2] || ''
    return {
      email: email.toLowerCase(),
      role,
      password,
      valid: emailRegex.test(email)
    }
  }).filter(u => u.valid)
})

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
    const rows = response.data || []
    teamMembers.value = rows.map((m) => ({
      ...m,
      role: normalizeUserRole(m.role)
    }))

    // Store original roles for comparison (normalized so DEV !== ENGINEER false positives are avoided)
    originalRoles.value.clear()
    teamMembers.value.forEach((member) => {
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
    const ok = await confirm({
      title: 'Confirm role change',
      message: `Change ${member.email} from ${originalRole} to ${member.role}? This is a critical change.`,
      confirmLabel: 'Yes, change role',
      cancelLabel: 'Cancel',
      variant: 'danger'
    })
    if (!ok) {
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
    
    showError(err.data?.message || err.message || 'Unknown error', 'Failed to change role')
  } finally {
    changingRoleFor.value = null
  }
}

// Toggle remote worker status (CEO only)
const handleRemoteToggle = async (member: TeamMember) => {
  if (member.id === currentUser.value?.user_id) return

  const newRemoteStatus = !member.is_remote
  changingRemoteFor.value = member.id

  try {
    await fetchWithAuth(`/auth/users/${member.id}/remote`, {
      method: 'PATCH',
      body: { is_remote: newRemoteStatus }
    })

    // Update local state
    member.is_remote = newRemoteStatus

    // Show success message
    showSuccess(
      newRemoteStatus
        ? `${member.email} set as remote worker`
        : `${member.email} unset as remote worker`
    )
  } catch (err: any) {
    console.error('Failed to toggle remote status:', err)
    showError(err.data?.message || err.message || 'Unknown error', 'Failed to change remote status')
  } finally {
    changingRemoteFor.value = null
  }
}

// Delete user (CEO only; cannot delete self) — opened from modal
const executeDeleteMember = async () => {
  const member = memberToDelete.value
  if (!member || member.id === currentUser.value?.user_id) {
    memberToDelete.value = null
    return
  }
  deletingUserId.value = member.id
  try {
    await fetchWithAuth(`/auth/users/${member.id}`, { method: 'DELETE' })
    successMessage.value = `${member.email} ถูกลบออกจากทีมแล้ว`
    setTimeout(() => { successMessage.value = '' }, 4000)
    memberToDelete.value = null
    await fetchTeamMembers()
  } catch (err: any) {
    showError(err?.data?.message || err?.message || 'ลบผู้ใช้ไม่สำเร็จ', 'ลบผู้ใช้ไม่สำเร็จ')
  } finally {
    deletingUserId.value = null
  }
}

// Reset password (CEO only) — system generates temp password
const openResetPasswordModal = (member: TeamMember) => {
  if (member.id === currentUser.value?.user_id) return
  memberToReset.value = member
  resetPasswordResult.value = null
  resetPasswordError.value = ''
}

const closeResetPasswordModal = () => {
  memberToReset.value = null
  resetPasswordResult.value = null
  resetPasswordError.value = ''
}

const submitResetPassword = async () => {
  const member = memberToReset.value
  if (!member) return
  resetPasswordError.value = ''
  resetPasswordSubmitting.value = true
  resettingPasswordUserId.value = member.id
  try {
    const res = await fetchWithAuth<{ temp_password: string }>(`/auth/users/${member.id}/password`, {
      method: 'PATCH',
      body: {}
    })
    resetPasswordResult.value = { temp_password: res.temp_password }
    successMessage.value = `รีเซ็ตรหัสผ่านของ ${member.email} แล้ว`
    setTimeout(() => { successMessage.value = '' }, 4000)
  } catch (err: any) {
    resetPasswordError.value = err?.data?.message || err?.message || 'รีเซ็ตรหัสผ่านไม่สำเร็จ'
  } finally {
    resetPasswordSubmitting.value = false
    resettingPasswordUserId.value = null
  }
}

// Create user (CEO)
const submitCreateUser = async () => {
  if (createForm.value.password !== createForm.value.confirmPassword) return
  if (createForm.value.password.length < 8) {
    createError.value = 'Password must be at least 8 characters'
    return
  }
  createError.value = ''
  createSubmitting.value = true
  try {
    await fetchWithAuth<{ data: TeamMember }>('/auth/users', {
      method: 'POST',
      body: {
        first_name: createForm.value.firstName.trim(),
        last_name: createForm.value.lastName.trim(),
        display_name: createForm.value.displayName.trim(),
        email: createForm.value.email.trim(),
        password: createForm.value.password,
        role: createForm.value.role
      }
    })
    successMessage.value = `User ${getMemberDisplayName({
      id: 0,
      email: createForm.value.email,
      first_name: createForm.value.firstName,
      last_name: createForm.value.lastName,
      display_name: createForm.value.displayName,
      role: createForm.value.role,
      health_score: 100,
      created_at: '',
      updated_at: ''
    })} created`
    setTimeout(() => { successMessage.value = '' }, 4000)
    showCreateModal.value = false
    createForm.value = { firstName: '', lastName: '', displayName: '', email: '', password: '', confirmPassword: '', role: 'ENGINEER' }
    await fetchTeamMembers()
  } catch (err: any) {
    createError.value = err?.data?.message || err?.message || 'Failed to create user'
  } finally {
    createSubmitting.value = false
  }
}

// Edit user (CEO)
const openEditModal = (member: TeamMember) => {
  memberToEdit.value = member
  editForm.value = {
    firstName: member.first_name || '',
    lastName: member.last_name || '',
    displayName: member.display_name || '',
    email: member.email
  }
  editError.value = ''
}

const submitEditUser = async () => {
  const member = memberToEdit.value
  if (!member) return
  editError.value = ''
  editSubmitting.value = true
  try {
    const payload: Record<string, string> = {}
    const fn = editForm.value.firstName.trim()
    const ln = editForm.value.lastName.trim()
    const dn = editForm.value.displayName.trim()
    const em = editForm.value.email.trim()
    if (fn !== (member.first_name || '')) payload.first_name = fn
    if (ln !== (member.last_name || '')) payload.last_name = ln
    if (dn !== (member.display_name || '')) payload.display_name = dn
    if (em !== member.email) payload.email = em

    if (Object.keys(payload).length > 0) {
      await fetchWithAuth(`/auth/users/${member.id}`, {
        method: 'PATCH',
        body: payload
      })
    }
    successMessage.value = `Updated ${em}`
    setTimeout(() => { successMessage.value = '' }, 4000)
    memberToEdit.value = null
    await fetchTeamMembers()
  } catch (err: any) {
    editError.value = err?.data?.message || err?.message || 'Failed to update user'
  } finally {
    editSubmitting.value = false
  }
}

// Import panel
const openImportPanel = () => {
  showImportPanel.value = true
  importRaw.value = ''
  importError.value = ''
  importResult.value = null
}

const closeImportPanel = () => {
  showImportPanel.value = false
  importResult.value = null
  importRaw.value = ''
  importError.value = ''
}

const submitImport = async () => {
  const users = parsedImport.value
  if (users.length === 0 || users.length > 500) return
  importError.value = ''
  importSubmitting.value = true
  try {
    const payload = {
      users: users.map(u => ({
        email: u.email,
        role: u.role,
        ...(u.password ? { password: u.password } : {})
      }))
    }
    const res = await fetchWithAuth<{ data: typeof importResult.value }>('/auth/users/import', {
      method: 'POST',
      body: payload,
      timeoutMs: 60000
    })
    importResult.value = res.data
    if (res.data.created > 0) {
      await fetchTeamMembers()
      successMessage.value = `Imported ${res.data.created} user(s)`
      setTimeout(() => { successMessage.value = '' }, 4000)
    }
  } catch (err: any) {
    importError.value = err?.data?.message || err?.message || 'Import failed'
  } finally {
    importSubmitting.value = false
  }
}

const copyTempPassword = (text: string) => {
  if (typeof navigator !== 'undefined' && navigator.clipboard?.writeText) {
    navigator.clipboard.writeText(text)
    successMessage.value = 'Copied to clipboard'
    setTimeout(() => { successMessage.value = '' }, 2000)
  }
}

// Export import results to CSV (Excel-compatible) for sharing with employees
const exportImportResultsToExcel = () => {
  if (!importResult.value) return
  const escapeCsv = (v: string): string => {
    const s = String(v ?? '')
    if (s.includes(',') || s.includes('"') || s.includes('\n')) {
      return '"' + s.replace(/"/g, '""') + '"'
    }
    return s
  }
  const loginUrl = typeof window !== 'undefined' ? `${window.location.origin}/login` : 'https://your-app-url/login'
  const rows: string[] = [
    'Email,Status,Message,Temp Password,Login URL',
    ...importResult.value.results.map((r) =>
      [r.email, r.status, r.message ?? '', r.temp_password ?? '', loginUrl].map(escapeCsv).join(',')
    )
  ]
  const csv = '\uFEFF' + rows.join('\n') // UTF-8 BOM for Excel
  const blob = new Blob([csv], { type: 'text/csv;charset=utf-8' })
  const name = `sentinel-import-users-${new Date().toISOString().slice(0, 19).replace(/[-:T]/g, '-')}.csv`
  const a = document.createElement('a')
  a.href = URL.createObjectURL(blob)
  a.download = name
  a.click()
  URL.revokeObjectURL(a.href)
  successMessage.value = 'Exported to CSV'
  setTimeout(() => { successMessage.value = '' }, 2500)
}

// Helper functions
const getRoleIcon = (role: string) => {
  const icons: Record<string, string> = {
    'CEO': '👑',
    'MANAGER': '🏢',
    'PRODUCT_OWNER': '📋',
    'PM': '📋',
    'ENGINEER': '💻',
    'CHIEF_ENGINEER': '⚙️',
    'DEV': '💻',
    'SUPPORT': '🎧'
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
    'MANAGER': {
      badge: 'bg-orange-700 text-orange-100',
      avatar: 'bg-orange-600 text-white',
      text: 'text-orange-400'
    },
    'PRODUCT_OWNER': {
      badge: 'bg-blue-700 text-blue-100',
      avatar: 'bg-blue-600 text-white',
      text: 'text-blue-400'
    },
    'PM': {
      badge: 'bg-blue-700 text-blue-100',
      avatar: 'bg-blue-600 text-white',
      text: 'text-blue-400'
    },
    'ENGINEER': {
      badge: 'bg-green-700 text-green-100',
      avatar: 'bg-green-600 text-white',
      text: 'text-green-400'
    },
    'CHIEF_ENGINEER': {
      badge: 'bg-green-700 text-green-100',
      avatar: 'bg-green-600 text-white',
      text: 'text-green-400'
    },
    'DEV': {
      badge: 'bg-green-700 text-green-100',
      avatar: 'bg-green-600 text-white',
      text: 'text-green-400'
    },
    'SUPPORT': {
      badge: 'bg-cyan-700 text-cyan-100',
      avatar: 'bg-cyan-600 text-white',
      text: 'text-cyan-400'
    }
  }
  return colors[role] || {
    badge: 'bg-gray-700 text-gray-100',
    avatar: 'bg-gray-600 text-white',
    text: 'text-gray-400'
  }
}

const getMemberDisplayName = (member: TeamMember) => {
  const parts = [member.first_name, member.last_name].map((part) => (part || '').trim()).filter(Boolean)
  if (parts.length > 0) return parts.join(' ')
  if (member.display_name?.trim()) return member.display_name.trim()
  return member.email
}

const getMemberAvatarInitial = (member: TeamMember) => {
  const name = getMemberDisplayName(member)
  return name.charAt(0).toUpperCase()
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

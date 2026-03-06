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
          class="px-6 py-3 bg-gray-700 hover:bg-gray-600 text-white rounded font-bold transition-colors"
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
            class="inline-flex items-center gap-2 px-4 py-2.5 bg-gradient-to-r from-purple-600 to-pink-600 hover:from-purple-500 hover:to-pink-500 text-white font-medium rounded-lg shadow-lg transition-all focus:outline-none focus:ring-2 focus:ring-purple-400 focus:ring-offset-2 focus:ring-offset-gray-900"
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
                    <option value="PM" class="bg-gray-900">📋 PM</option>
                    <option value="DEV" class="bg-gray-900">💻 DEV</option>
                  </select>
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
                class="p-2 text-gray-400 hover:text-white rounded-lg hover:bg-gray-700 transition-colors"
                aria-label="Close"
              >
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" /></svg>
              </button>
            </div>
            <form @submit.prevent="submitCreateUser" class="p-6 space-y-4">
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
                <label for="create-role" class="block text-sm font-medium text-gray-300 mb-1">Role</label>
                <select
                  id="create-role"
                  v-model="createForm.role"
                  class="w-full px-4 py-2.5 bg-gray-900 border border-gray-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-transparent"
                >
                  <option value="DEV">💻 DEV</option>
                  <option value="PM">📋 PM</option>
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
                  class="flex-1 px-4 py-2.5 bg-gradient-to-r from-purple-600 to-pink-600 hover:from-purple-500 hover:to-pink-500 disabled:opacity-50 disabled:cursor-not-allowed text-white font-medium rounded-lg transition-colors"
                >
                  {{ createSubmitting ? 'Creating…' : 'Create user' }}
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
                class="p-2 text-gray-400 hover:text-white rounded-lg hover:bg-gray-700 transition-colors"
                aria-label="Close"
              >
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" /></svg>
              </button>
            </div>

            <div v-if="!importResult" class="flex-1 flex flex-col overflow-hidden p-6">
              <p class="text-sm text-gray-400 mb-3">
                One user per line. Format: <code class="px-1.5 py-0.5 bg-gray-700 rounded text-gray-300">email</code> or
                <code class="px-1.5 py-0.5 bg-gray-700 rounded text-gray-300">email,role</code> or
                <code class="px-1.5 py-0.5 bg-gray-700 rounded text-gray-300">email,role,password</code>. Role defaults to DEV. Max 500.
              </p>
              <textarea
                v-model="importRaw"
                rows="12"
                placeholder="dev1@company.com&#10;pm1@company.com,PM&#10;ceo@company.com,CEO,SecurePass123"
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
                  class="inline-flex items-center gap-2 px-4 py-2 bg-emerald-700 hover:bg-emerald-600 text-white rounded-lg font-medium transition-colors"
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
                    class="px-4 py-2 bg-gradient-to-r from-purple-600 to-pink-600 text-white rounded-lg font-medium transition-colors"
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
                class="px-4 py-2.5 bg-red-600 hover:bg-red-500 disabled:opacity-50 disabled:cursor-not-allowed text-white font-medium rounded-lg transition-colors"
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
                  class="flex-1 px-4 py-2.5 bg-amber-600 hover:bg-amber-500 disabled:opacity-50 disabled:cursor-not-allowed text-white font-medium rounded-lg transition-colors"
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
                  class="px-4 py-2.5 bg-amber-600 hover:bg-amber-500 text-white font-medium rounded-lg transition-colors"
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

// Create user modal
const showCreateModal = ref(false)
const createForm = ref({
  email: '',
  password: '',
  confirmPassword: '',
  role: 'DEV'
})
const createError = ref('')
const createSubmitting = ref(false)

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
const isCEO = computed(() => currentUser.value?.role === 'CEO')

const pmCount = computed(() => 
  teamMembers.value.filter(m => m.role === 'PM').length
)

const devCount = computed(() => 
  teamMembers.value.filter(m => m.role === 'DEV').length
)

// Parse import text into users array (email, role?, password?)
const parsedImport = computed(() => {
  const lines = importRaw.value
    .split('\n')
    .map(l => l.trim())
    .filter(Boolean)
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  const roles = ['CEO', 'PM', 'DEV']
  return lines.map(line => {
    const parts = line.split(',').map(p => p.trim())
    const email = parts[0] || ''
    const role = parts[1] && roles.includes(parts[1]) ? parts[1] : 'DEV'
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
        email: createForm.value.email.trim(),
        password: createForm.value.password,
        role: createForm.value.role
      }
    })
    successMessage.value = `User ${createForm.value.email} created`
    setTimeout(() => { successMessage.value = '' }, 4000)
    showCreateModal.value = false
    createForm.value = { email: '', password: '', confirmPassword: '', role: 'DEV' }
    await fetchTeamMembers()
  } catch (err: any) {
    createError.value = err?.data?.message || err?.message || 'Failed to create user'
  } finally {
    createSubmitting.value = false
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

<template>
  <div class="min-h-screen bg-gray-900 p-8">
    <!-- Access Guard -->
    <div v-if="!isCEO" class="flex items-center justify-center h-[80vh]">
      <div class="text-center">
        <div class="text-6xl mb-4">🚫</div>
        <h2 class="text-3xl font-bold text-red-400 mb-2">ACCESS DENIED</h2>
        <p class="text-gray-400 mb-6">This area is restricted to CEO access only.</p>
        <button @click="navigateTo('/dashboard')" class="px-6 py-3 bg-gray-700 hover:bg-gray-600 text-white rounded font-bold transition-colors">
          Return to Dashboard
        </button>
      </div>
    </div>

    <!-- CEO Teams Management -->
    <div v-else>
      <!-- Header -->
      <div class="mb-6 border-b border-gray-700 pb-4 flex flex-wrap items-end justify-between gap-4">
        <div>
          <h1 class="text-2xl font-bold text-white">Internal Venture Capital</h1>
          <p class="text-sm text-gray-400 mt-1">Inject capital into squads, monitor burn rate & distribute milestone bonuses</p>
        </div>
        <div class="flex items-center gap-3">
          <button
            v-if="teamsStore.teamsFeatureEnabled"
            type="button"
            @click="confirmDisableFeatureTeam"
            :disabled="isDisablingFeature"
            class="inline-flex items-center gap-2 px-4 py-2.5 bg-gray-700 hover:bg-red-900/40 border border-gray-600 hover:border-red-600/50 text-gray-300 hover:text-red-400 font-medium rounded-lg transition-all"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636"/></svg>
            {{ isDisablingFeature ? 'Disabling...' : 'Disable feature team' }}
          </button>
          <button
            v-else
            type="button"
            @click="enableFeatureTeam"
            :disabled="isEnablingFeature"
            class="inline-flex items-center gap-2 px-4 py-2.5 bg-emerald-600/20 hover:bg-emerald-600/30 border border-emerald-500/40 text-emerald-400 font-medium rounded-lg transition-all"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>
            {{ isEnablingFeature ? 'Enabling...' : 'Enable feature team' }}
          </button>
          <button
            v-if="teamsStore.teamsFeatureEnabled"
            @click="showCreateTeamModal = true"
            class="inline-flex items-center gap-2 px-4 py-2.5 bg-gradient-to-r from-purple-600 to-pink-600 hover:from-purple-500 hover:to-pink-500 text-white font-medium rounded-lg shadow-lg transition-all"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/></svg>
            New Team
          </button>
        </div>
      </div>

      <!-- Feature disabled banner -->
      <div
        v-if="!teamsStore.teamsFeatureEnabled"
        class="mb-6 p-4 bg-amber-900/20 border border-amber-600/40 rounded-xl text-amber-200"
      >
        <p class="font-medium">Feature team is disabled</p>
        <p class="text-sm text-amber-200/80 mt-0.5">การมองเห็นทีม (Squads) ถูกปิดทั้งหมด — ไม่แสดงใน Dashboard, Projects และหน้าอื่นๆ กด "Enable feature team" เพื่อเปิดใช้กลับ</p>
        <p class="text-sm text-amber-200/80 mt-2">ในขณะที่ปิดทีม: ไปที่ <span class="text-amber-100 font-medium">Projects</span> เพื่อให้ CEO กำหนด Product Owner เจ้าของโปรเจกต์ได้หลายคนต่อโปรเจกต์</p>
      </div>

      <!-- Metrics (only when feature enabled) -->
      <div v-if="teamsStore.teamsFeatureEnabled" class="grid grid-cols-2 md:grid-cols-4 gap-4 mb-8">
        <div class="bg-gray-800 border border-gray-700 rounded-lg p-4">
          <div class="text-2xl font-bold text-purple-400">{{ teamsStore.teams.length }}</div>
          <div class="text-xs text-gray-500 uppercase mt-1">Total Teams</div>
        </div>
        <div class="bg-gray-800 border border-gray-700 rounded-lg p-4">
          <div class="text-2xl font-bold text-blue-400">{{ totalMembers }}</div>
          <div class="text-xs text-gray-500 uppercase mt-1">Assigned Members</div>
        </div>
        <div class="bg-gray-800 border border-gray-700 rounded-lg p-4">
          <div class="text-2xl font-bold text-emerald-400">{{ formatCurrency(totalCapital) }}</div>
          <div class="text-xs text-gray-500 uppercase mt-1">Total Capital</div>
        </div>
        <div class="bg-gray-800 border border-gray-700 rounded-lg p-4">
          <div class="text-2xl font-bold text-orange-400">{{ formatCurrency(totalMonthlyBurn) }}</div>
          <div class="text-xs text-gray-500 uppercase mt-1">Total Burn / Mo</div>
        </div>
      </div>

      <!-- Loading (only when feature enabled) -->
      <div v-if="teamsStore.teamsFeatureEnabled && teamsStore.loading" class="text-center py-16 text-gray-500">
        <div class="animate-spin text-4xl mb-4">⚙️</div>
        Loading teams...
      </div>

      <!-- Error -->
      <div v-else-if="teamsStore.teamsFeatureEnabled && teamsStore.error" class="p-4 bg-red-900/20 border border-red-600 rounded-lg text-red-400 mb-6">
        {{ teamsStore.error }}
      </div>

      <!-- Teams List (only when feature enabled) -->
      <div v-else-if="teamsStore.teamsFeatureEnabled" class="space-y-6">
        <div
          v-for="team in teamsStore.teams"
          :key="team.id"
          class="bg-gray-800 border border-gray-700 rounded-xl overflow-hidden"
        >
          <!-- Team Header -->
          <div class="flex items-center justify-between px-5 py-4">
            <div class="flex items-center gap-3">
              <div class="flex h-9 w-9 items-center justify-center rounded-full bg-purple-600/20 text-purple-400">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z"/></svg>
              </div>
              <div>
                <h3 class="text-white font-semibold">{{ team.name }}</h3>
                <p class="text-xs text-gray-500">{{ team.users?.length ?? 0 }} member{{ (team.users?.length ?? 0) !== 1 ? 's' : '' }}</p>
              </div>
              <button
                @click="openEditNameModal(team)"
                class="p-1.5 text-gray-500 hover:text-purple-400 rounded-lg transition-colors"
                title="Edit team name"
              >
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/></svg>
              </button>
            </div>
            <div class="flex items-center gap-2">
              <button
                @click="toggleExpand(team.id)"
                class="px-3 py-1.5 text-xs bg-gray-700 hover:bg-gray-600 text-gray-300 rounded-lg transition-colors"
              >
                {{ expandedTeam === team.id ? 'Collapse' : 'Manage Members' }}
              </button>
              <button
                @click="deleteTeam(team)"
                class="px-3 py-1.5 text-xs bg-red-900/40 hover:bg-red-700/40 text-red-400 rounded-lg transition-colors"
              >
                Delete
              </button>
            </div>
          </div>

          <!-- Financial Control Section -->
          <div class="border-t border-gray-700/60 bg-gray-900/30 px-5 py-4">
            <div class="flex items-center justify-between mb-4">
              <h4 class="text-xs font-semibold text-yellow-400/80 uppercase tracking-wider flex items-center gap-1.5">
                <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>
                Financial Control
              </h4>
              <button
                v-if="!teamsStore.teamCosts[team.id]"
                @click="teamsStore.fetchTeamCost(team.id)"
                :disabled="teamsStore.costLoading[team.id]"
                class="text-xs text-gray-500 hover:text-gray-300 transition-colors"
              >
                {{ teamsStore.costLoading[team.id] ? 'Loading...' : 'Load Financials' }}
              </button>
              <button
                v-else
                @click="teamsStore.fetchTeamCost(team.id)"
                :disabled="teamsStore.costLoading[team.id]"
                class="text-xs text-gray-600 hover:text-gray-400 transition-colors"
              >
                ↻ Refresh
              </button>
            </div>

            <div v-if="teamsStore.costLoading[team.id]" class="text-center py-4 text-gray-600 text-sm">
              Computing burn rate...
            </div>

            <div v-else-if="teamsStore.teamCosts[team.id]" class="space-y-4">
              <!-- Cost / Balance Row -->
              <div class="grid grid-cols-3 gap-3">
                <div class="bg-gray-800/80 rounded-lg p-3 text-center">
                  <div class="text-lg font-bold text-orange-400">{{ formatCurrency(teamsStore.teamCosts[team.id]!.total_monthly_cost) }}</div>
                  <div class="text-xs text-gray-500 mt-0.5">Burn / Month</div>
                </div>
                <div class="bg-gray-800/80 rounded-lg p-3 text-center">
                  <div class="text-lg font-bold text-emerald-400">{{ formatCurrency(teamsStore.teamCosts[team.id]!.capital_balance) }}</div>
                  <div class="text-xs text-gray-500 mt-0.5 flex items-center justify-center gap-1">
                    Capital Balance
                    <button
                      @click="openEditCapitalModal(team.id)"
                      class="text-gray-600 hover:text-gray-300 transition-colors ml-1"
                      title="Edit balance"
                    >
                      <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/></svg>
                    </button>
                  </div>
                </div>
                <div class="bg-gray-800/80 rounded-lg p-3 text-center">
                  <div :class="['text-lg font-bold', runwayColor(teamsStore.teamCosts[team.id]!.runway_months)]">
                    {{ teamsStore.teamCosts[team.id]!.runway_months > 0 ? teamsStore.teamCosts[team.id]!.runway_months.toFixed(1) + ' mo' : '—' }}
                  </div>
                  <div class="text-xs text-gray-500 mt-0.5">Runway</div>
                </div>
              </div>

              <!-- Runway Progress Bar -->
              <div>
                <div class="flex items-center justify-between text-xs text-gray-500 mb-1.5">
                  <span>Runway</span>
                  <span>{{ teamsStore.teamCosts[team.id]!.bonus_percentage }}% bonus target</span>
                </div>
                <div class="h-2 bg-gray-700 rounded-full overflow-hidden">
                  <div
                    :class="['h-full rounded-full transition-all duration-500', runwayBarColor(teamsStore.teamCosts[team.id]!.runway_months)]"
                    :style="{ width: runwayBarWidth(teamsStore.teamCosts[team.id]!.runway_months) }"
                  />
                </div>
                <div class="flex justify-between text-xs text-gray-600 mt-1">
                  <span>0 mo</span>
                  <span>6 mo</span>
                  <span>12+ mo</span>
                </div>
              </div>

              <!-- Cost Breakdown (collapsible) -->
              <details class="group">
                <summary class="text-xs text-gray-500 cursor-pointer hover:text-gray-300 transition-colors list-none flex items-center gap-1">
                  <svg class="w-3 h-3 group-open:rotate-90 transition-transform" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/></svg>
                  Cost breakdown
                </summary>
                <div class="mt-2 grid grid-cols-2 gap-2 text-xs">
                  <div class="flex justify-between bg-gray-800/50 rounded px-2 py-1.5">
                    <span class="text-gray-500">Member Salaries + SS</span>
                    <span class="text-gray-300">{{ formatCurrency(teamsStore.teamCosts[team.id]!.member_cost) }}</span>
                  </div>
                  <div class="flex justify-between bg-gray-800/50 rounded px-2 py-1.5">
                    <span class="text-gray-500">Shared Overhead</span>
                    <span class="text-gray-300">{{ formatCurrency(teamsStore.teamCosts[team.id]!.shared_overhead) }}</span>
                  </div>
                </div>
              </details>

              <!-- Action Buttons -->
              <div class="flex items-center gap-2 pt-1">
                <button
                  @click="openInjectModal(team.id)"
                  class="flex-1 flex items-center justify-center gap-2 py-2 text-sm font-medium bg-emerald-600/20 hover:bg-emerald-600/30 border border-emerald-600/40 text-emerald-400 rounded-lg transition-all"
                >
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/></svg>
                  Inject Capital
                </button>
              </div>
            </div>

            <!-- Not loaded yet -->
            <div v-else class="text-center py-3 text-gray-600 text-xs">
              Click "Load Financials" to view burn rate and capital data
            </div>
          </div>

          <!-- Expanded: Members List + Add Member -->
          <div v-if="expandedTeam === team.id" class="border-t border-gray-700 px-5 py-4">
            <h4 class="text-xs font-semibold text-gray-400 uppercase tracking-wider mb-3">Members</h4>
            <div v-if="team.users && team.users.length > 0" class="space-y-2 mb-4">
              <div
                v-for="member in team.users"
                :key="member.id"
                class="flex items-center justify-between py-2 px-3 bg-gray-700/40 rounded-lg"
              >
                <div class="flex items-center gap-3">
                  <div class="h-7 w-7 rounded-full bg-gray-600 flex items-center justify-center text-xs font-bold text-gray-300">
                    {{ (member.display_name || member.email).charAt(0).toUpperCase() }}
                  </div>
                  <div>
                    <p class="text-sm text-gray-200">{{ member.display_name || member.email }}</p>
                    <p class="text-xs text-gray-500">{{ member.email }} · {{ member.role }}</p>
                  </div>
                </div>
                <button
                  @click="removeMemberFromTeam(member.id)"
                  class="text-xs text-gray-500 hover:text-red-400 transition-colors"
                >
                  Remove
                </button>
              </div>
            </div>
            <p v-else class="text-sm text-gray-600 mb-4">No members in this team yet.</p>

            <!-- Add Member Dropdown -->
            <div class="flex items-center gap-2">
              <select
                v-model="addMemberUserID[team.id]"
                class="flex-1 bg-gray-700 border border-gray-600 rounded-lg px-3 py-1.5 text-sm text-gray-300 focus:outline-none focus:border-purple-500 transition-colors"
              >
                <option value="">— Select user to add —</option>
                <option
                  v-for="user in unassignedUsers"
                  :key="user.id"
                  :value="user.id"
                >
                  {{ user.display_name || user.email }} ({{ user.role }})
                </option>
              </select>
              <button
                @click="addMemberToTeam(team.id)"
                :disabled="!addMemberUserID[team.id]"
                class="px-4 py-1.5 text-sm bg-purple-600 hover:bg-purple-500 disabled:bg-gray-700 disabled:text-gray-600 text-white rounded-lg transition-colors"
              >
                Add
              </button>
            </div>
          </div>
        </div>

        <!-- Empty State -->
        <div v-if="teamsStore.teams.length === 0" class="text-center py-16 text-gray-600">
          <svg class="w-12 h-12 mx-auto mb-4 opacity-30" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z"/></svg>
          <p>No teams yet. Create your first squad!</p>
        </div>
      </div>
    </div>

    <!-- ===================== MODALS ===================== -->

    <!-- Create Team Modal -->
    <div
      v-if="showCreateTeamModal"
      class="fixed inset-0 bg-black/80 backdrop-blur-sm flex items-center justify-center z-50 p-4"
      @click.self="showCreateTeamModal = false"
    >
      <div class="bg-gray-800 border border-gray-700 rounded-2xl p-6 max-w-md w-full shadow-2xl">
        <div class="flex items-center justify-between mb-5">
          <h2 class="text-lg font-bold text-white">Create New Team</h2>
          <button @click="showCreateTeamModal = false" class="text-gray-500 hover:text-white">✕</button>
        </div>
        <div class="space-y-4">
          <div>
            <label class="block text-sm text-gray-400 mb-1.5">Team Name <span class="text-red-400">*</span></label>
            <input
              v-model="newTeamName"
              type="text"
              placeholder="e.g. Squad Alpha"
              class="w-full bg-gray-700 border border-gray-600 rounded-xl px-4 py-2.5 text-sm text-gray-200 placeholder-gray-500 focus:outline-none focus:border-purple-500 transition-colors"
              @keyup.enter="createTeam"
            />
          </div>
          <div v-if="createTeamError" class="p-3 bg-red-900/30 border border-red-600 rounded-lg text-red-400 text-sm">
            {{ createTeamError }}
          </div>
          <div class="flex gap-3 pt-1">
            <button
              @click="createTeam"
              :disabled="!newTeamName.trim() || isCreatingTeam"
              class="flex-1 py-2.5 bg-gradient-to-r from-purple-600 to-pink-600 hover:from-purple-500 hover:to-pink-500 disabled:from-gray-700 disabled:to-gray-700 text-white font-semibold rounded-xl transition-all"
            >
              {{ isCreatingTeam ? 'Creating...' : 'Create Team' }}
            </button>
            <button @click="showCreateTeamModal = false" class="px-5 py-2.5 bg-gray-700 hover:bg-gray-600 text-gray-300 rounded-xl transition-colors">
              Cancel
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Disable feature team confirm -->
    <div
      v-if="showDisableFeatureConfirm"
      class="fixed inset-0 bg-black/80 backdrop-blur-sm flex items-center justify-center z-50 p-4"
      @click.self="showDisableFeatureConfirm = false"
    >
      <div class="bg-gray-800 border border-gray-700 rounded-2xl p-6 max-w-md w-full shadow-2xl">
        <h2 class="text-lg font-bold text-white mb-3">Disable feature team</h2>
        <p class="text-gray-400 text-sm mb-5">
          ยกเลิกการมองเห็นเป็นทีมออกทั้งหมด — ลิงก์ Squads ใน Dashboard, คอลัมน์ทีมใน Projects และส่วนที่เกี่ยวกับทีมจะไม่แสดงจนกว่าจะกด Enable อีกครั้ง
        </p>
        <div class="flex gap-3">
          <button
            @click="doDisableFeatureTeam"
            :disabled="isDisablingFeature"
            class="flex-1 py-2.5 bg-red-600 hover:bg-red-700 disabled:opacity-50 text-white font-semibold rounded-xl transition-colors"
          >
            {{ isDisablingFeature ? 'Disabling...' : 'Disable' }}
          </button>
          <button @click="showDisableFeatureConfirm = false" class="px-5 py-2.5 bg-gray-700 hover:bg-gray-600 text-gray-300 rounded-xl transition-colors">
            Cancel
          </button>
        </div>
      </div>
    </div>

    <!-- Delete Confirm Modal -->
    <div
      v-if="teamToDelete"
      class="fixed inset-0 bg-black/80 backdrop-blur-sm flex items-center justify-center z-50 p-4"
      @click.self="teamToDelete = null"
    >
      <div class="bg-gray-800 border border-gray-700 rounded-2xl p-6 max-w-md w-full shadow-2xl">
        <h2 class="text-lg font-bold text-white mb-3">Delete Team</h2>
        <p class="text-gray-400 text-sm mb-5">
          Are you sure you want to delete <strong class="text-white">{{ teamToDelete.name }}</strong>?
          All members will be unassigned from this team.
        </p>
        <div class="flex gap-3">
          <button
            @click="confirmDeleteTeam"
            class="flex-1 py-2.5 bg-red-600 hover:bg-red-700 text-white font-semibold rounded-xl transition-colors"
          >
            Delete Team
          </button>
          <button @click="teamToDelete = null" class="px-5 py-2.5 bg-gray-700 hover:bg-gray-600 text-gray-300 rounded-xl transition-colors">
            Cancel
          </button>
        </div>
      </div>
    </div>

    <!-- Inject Capital Modal -->
    <div
      v-if="injectModal.open"
      class="fixed inset-0 bg-black/80 backdrop-blur-sm flex items-center justify-center z-50 p-4"
      @click.self="injectModal.open = false"
    >
      <div class="bg-gray-800 border border-emerald-700/40 rounded-2xl p-6 max-w-md w-full shadow-2xl">
        <div class="flex items-center justify-between mb-5">
          <div>
            <h2 class="text-lg font-bold text-white">Inject Capital</h2>
            <p class="text-xs text-gray-500 mt-0.5">{{ teamName(injectModal.teamId) }}</p>
          </div>
          <button @click="injectModal.open = false" class="text-gray-500 hover:text-white">✕</button>
        </div>

        <div class="space-y-4">
          <!-- วันที่ส่งงาน -->
          <div>
            <label class="block text-sm text-gray-400 mb-1.5">วันที่ส่งงาน <span class="text-red-400">*</span></label>
            <UiDatePicker
              v-model="injectModal.deliveryDate"
              placeholder="เลือกวันที่ส่งงาน…"
              :min="injectModalDeliveryDateMin"
            />
            <p class="text-xs text-gray-500 mt-1">ระบบจะคำนวณยอดเงินให้หมดพอดีวันส่งงาน (จาก Burn Rate ปัจจุบัน)</p>
          </div>

          <!-- Amount (auto-calculated, editable) -->
          <div>
            <label class="block text-sm text-gray-400 mb-1.5">Amount (THB) <span class="text-red-400">*</span></label>
            <input
              v-model.number="injectModal.amount"
              type="number"
              min="1"
              placeholder="คำนวณอัตโนมัติเมื่อเลือกวันที่ส่งงาน"
              class="w-full bg-gray-700 border border-gray-600 rounded-xl px-4 py-2.5 text-sm text-gray-200 placeholder-gray-500 focus:outline-none focus:border-emerald-500 transition-colors"
            />
            <p v-if="injectModal.deliveryDate && injectModalAmountFromDate !== null" class="text-xs text-emerald-400/90 mt-1">
              คำนวณจาก Burn {{ formatCurrency(teamsStore.teamCosts[injectModal.teamId]?.total_monthly_cost ?? 0) }}/mo × {{ injectModalMonthsToDelivery.toFixed(1) }} mo
            </p>
          </div>

          <!-- Note -->
          <div>
            <label class="block text-sm text-gray-400 mb-1.5">Note / Reference</label>
            <input
              v-model="injectModal.note"
              type="text"
              placeholder="e.g. งวดที่ 1 MIMS HD-MAP"
              class="w-full bg-gray-700 border border-gray-600 rounded-xl px-4 py-2.5 text-sm text-gray-200 placeholder-gray-500 focus:outline-none focus:border-emerald-500 transition-colors"
            />
          </div>

          <div v-if="injectModal.error" class="p-3 bg-red-900/30 border border-red-600 rounded-lg text-red-400 text-sm">
            {{ injectModal.error }}
          </div>

          <div class="flex gap-3 pt-1">
            <button
              @click="confirmInjectCapital"
              :disabled="!injectModal.amount || injectModal.amount <= 0 || injectModal.loading"
              class="flex-1 py-2.5 bg-emerald-600 hover:bg-emerald-500 disabled:bg-gray-700 disabled:text-gray-600 text-white font-semibold rounded-xl transition-all"
            >
              {{ injectModal.loading ? 'Injecting...' : 'Inject Capital' }}
            </button>
            <button @click="injectModal.open = false" class="px-5 py-2.5 bg-gray-700 hover:bg-gray-600 text-gray-300 rounded-xl transition-colors">
              Cancel
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Edit Team Name Modal -->
    <div
      v-if="editNameModal.open"
      class="fixed inset-0 bg-black/80 backdrop-blur-sm flex items-center justify-center z-50 p-4"
      @click.self="editNameModal.open = false"
    >
      <div class="bg-gray-800 border border-purple-700/40 rounded-2xl p-6 max-w-md w-full shadow-2xl">
        <div class="flex items-center justify-between mb-5">
          <h2 class="text-lg font-bold text-white">Edit Team Name</h2>
          <button @click="editNameModal.open = false" class="text-gray-500 hover:text-white">✕</button>
        </div>
        <div class="space-y-4">
          <div>
            <label class="block text-sm text-gray-400 mb-1.5">Team Name <span class="text-red-400">*</span></label>
            <input
              v-model="editNameModal.name"
              type="text"
              placeholder="e.g. Squad Alpha"
              class="w-full bg-gray-700 border border-gray-600 rounded-xl px-4 py-2.5 text-sm text-gray-200 placeholder-gray-500 focus:outline-none focus:border-purple-500 transition-colors"
              @keyup.enter="confirmEditName"
            />
          </div>
          <div v-if="editNameModal.error" class="p-3 bg-red-900/30 border border-red-600 rounded-lg text-red-400 text-sm">
            {{ editNameModal.error }}
          </div>
          <div class="flex gap-3 pt-1">
            <button
              @click="confirmEditName"
              :disabled="!editNameModal.name.trim() || editNameModal.loading"
              class="flex-1 py-2.5 bg-purple-600 hover:bg-purple-500 disabled:bg-gray-700 disabled:text-gray-600 text-white font-semibold rounded-xl transition-all"
            >
              {{ editNameModal.loading ? 'Saving...' : 'Save' }}
            </button>
            <button @click="editNameModal.open = false" class="px-5 py-2.5 bg-gray-700 hover:bg-gray-600 text-gray-300 rounded-xl transition-colors">
              Cancel
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Edit Capital Modal -->
    <div
      v-if="editCapitalModal.open"
      class="fixed inset-0 bg-black/80 backdrop-blur-sm flex items-center justify-center z-50 p-4"
      @click.self="editCapitalModal.open = false"
    >
      <div class="bg-gray-800 border border-blue-700/40 rounded-2xl p-6 max-w-md w-full shadow-2xl">
        <div class="flex items-center justify-between mb-5">
          <div>
            <h2 class="text-lg font-bold text-white">Edit Capital Balance</h2>
            <p class="text-xs text-gray-500 mt-0.5">{{ teamName(editCapitalModal.teamId) }} — แก้ไขยอด capital โดยตรง</p>
          </div>
          <button @click="editCapitalModal.open = false" class="text-gray-500 hover:text-white">✕</button>
        </div>

        <div class="space-y-4">
          <!-- Current vs New -->
          <div class="bg-gray-900/50 border border-gray-700 rounded-xl px-4 py-3 flex items-center justify-between text-sm mb-1">
            <span class="text-gray-400">Current Balance</span>
            <span class="text-emerald-400 font-semibold tabular-nums">
              {{ formatCurrency(teamsStore.teamCosts[editCapitalModal.teamId]?.capital_balance ?? 0) }}
            </span>
          </div>

          <!-- New Balance -->
          <div>
            <label class="block text-sm text-gray-400 mb-1.5">New Balance (THB) <span class="text-red-400">*</span></label>
            <input
              v-model.number="editCapitalModal.newBalance"
              type="number"
              min="0"
              step="1000"
              placeholder="e.g. 250000"
              class="w-full bg-gray-700 border border-gray-600 rounded-xl px-4 py-2.5 text-sm text-gray-200 placeholder-gray-500 focus:outline-none focus:border-blue-500 transition-colors"
            />
            <p v-if="editCapitalModal.newBalance !== null" class="text-xs mt-1.5"
              :class="(editCapitalModal.newBalance ?? 0) >= (teamsStore.teamCosts[editCapitalModal.teamId]?.capital_balance ?? 0)
                ? 'text-emerald-400'
                : 'text-red-400'">
              {{ (editCapitalModal.newBalance ?? 0) >= (teamsStore.teamCosts[editCapitalModal.teamId]?.capital_balance ?? 0)
                ? '↑ เพิ่มขึ้น ' + formatCurrency((editCapitalModal.newBalance ?? 0) - (teamsStore.teamCosts[editCapitalModal.teamId]?.capital_balance ?? 0))
                : '↓ ลดลง ' + formatCurrency((teamsStore.teamCosts[editCapitalModal.teamId]?.capital_balance ?? 0) - (editCapitalModal.newBalance ?? 0)) }}
            </p>
          </div>

          <!-- Bonus % (optional) -->
          <div>
            <label class="block text-sm text-gray-400 mb-1.5">
              Bonus % (optional)
              <span v-if="editCapitalModal.bonusPct !== null" class="text-yellow-400 font-bold ml-1">{{ editCapitalModal.bonusPct }}%</span>
            </label>
            <input
              v-model.number="editCapitalModal.bonusPct"
              type="number"
              min="0"
              max="100"
              step="1"
              placeholder="ไม่เปลี่ยนแปลง (leave blank)"
              class="w-full bg-gray-700 border border-gray-600 rounded-xl px-4 py-2.5 text-sm text-gray-200 placeholder-gray-500 focus:outline-none focus:border-blue-500 transition-colors"
            />
          </div>

          <!-- Note -->
          <div>
            <label class="block text-sm text-gray-400 mb-1.5">Reason / Note</label>
            <input
              v-model="editCapitalModal.note"
              type="text"
              placeholder="เช่น แก้ไขหลัง reconcile กับ Finance"
              class="w-full bg-gray-700 border border-gray-600 rounded-xl px-4 py-2.5 text-sm text-gray-200 placeholder-gray-500 focus:outline-none focus:border-blue-500 transition-colors"
            />
          </div>

          <div v-if="editCapitalModal.error" class="p-3 bg-red-900/30 border border-red-600 rounded-lg text-red-400 text-sm">
            {{ editCapitalModal.error }}
          </div>

          <div class="flex gap-3 pt-1">
            <button
              @click="confirmEditCapital"
              :disabled="editCapitalModal.newBalance === null || editCapitalModal.newBalance < 0 || editCapitalModal.loading"
              class="flex-1 py-2.5 bg-blue-600 hover:bg-blue-500 disabled:bg-gray-700 disabled:text-gray-600 text-white font-semibold rounded-xl transition-all"
            >
              {{ editCapitalModal.loading ? 'Saving...' : 'Save Changes' }}
            </button>
            <button @click="editCapitalModal.open = false" class="px-5 py-2.5 bg-gray-700 hover:bg-gray-600 text-gray-300 rounded-xl transition-colors">
              Cancel
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useAuth } from '~/composables/useAuth'
import { useTeamsStore } from '~/core/modules/teams/store/teams-store'
import type { Team, TeamUser } from '~/core/modules/teams/infrastructure/teams-api'
definePageMeta({ layout: 'default', middleware: 'auth' })

const { currentUser, fetchWithAuth } = useAuth()
const teamsStore = useTeamsStore()

const isCEO = computed(() => currentUser.value?.role === 'CEO' || currentUser.value?.role === 'MANAGER')

// --- State ---
const expandedTeam = ref<number | null>(null)
const addMemberUserID = ref<Record<number, number | ''>>({})
const allUsers = ref<TeamUser[]>([])

const showCreateTeamModal = ref(false)
const newTeamName = ref('')
const isCreatingTeam = ref(false)
const createTeamError = ref('')

const teamToDelete = ref<Team | null>(null)

const isDisablingFeature = ref(false)
const isEnablingFeature = ref(false)
const showDisableFeatureConfirm = ref(false)

// Inject Capital modal state
const injectModal = reactive({
  open: false,
  teamId: 0,
  deliveryDate: '' as string,
  amount: null as number | null,
  note: '',
  loading: false,
  error: '',
})

// Edit Capital modal state
const editCapitalModal = reactive({
  open: false,
  teamId: 0,
  newBalance: null as number | null,
  bonusPct: null as number | null,
  note: '',
  loading: false,
  error: '',
})

// Edit Team Name modal state
const editNameModal = reactive({
  open: false,
  teamId: 0,
  name: '',
  loading: false,
  error: '',
})

// --- Computed ---
const totalMembers = computed(() =>
  teamsStore.teams.reduce((sum, t) => sum + (t.users?.length ?? 0), 0)
)

const totalCapital = computed(() =>
  teamsStore.teams.reduce((sum, t) => sum + (t.capital_balance ?? 0), 0)
)

const totalMonthlyBurn = computed(() =>
  Object.values(teamsStore.teamCosts).reduce((sum, c) => sum + (c?.total_monthly_cost ?? 0), 0)
)

const assignedUserIDs = computed(() => {
  const ids = new Set<number>()
  for (const team of teamsStore.teams) {
    for (const u of (team.users ?? [])) ids.add(u.id)
  }
  return ids
})

const unassignedUsers = computed(() =>
  allUsers.value.filter(u => !assignedUserIDs.value.has(u.id) && u.role !== 'CEO' && u.role !== 'MANAGER')
)

const unassignedCount = computed(() => unassignedUsers.value.length)

// Inject modal: min date for delivery = today
const injectModalDeliveryDateMin = computed(() => new Date().toISOString().slice(0, 10))

// Months from today to delivery date (for display and calculation)
const injectModalMonthsToDelivery = computed(() => {
  if (!injectModal.deliveryDate || !injectModal.teamId) return 0
  const cost = teamsStore.teamCosts[injectModal.teamId]
  if (!cost?.total_monthly_cost) return 0
  const start = new Date()
  start.setHours(0, 0, 0, 0)
  const end = new Date(injectModal.deliveryDate)
  end.setHours(0, 0, 0, 0)
  const days = Math.max(0, (end.getTime() - start.getTime()) / (1000 * 60 * 60 * 24))
  return days / (365 / 12) // ~30.44 days per month
})

// Suggested amount so capital runs out exactly on delivery date
const injectModalAmountFromDate = computed(() => {
  if (!injectModal.teamId || injectModalMonthsToDelivery.value <= 0) return null
  const cost = teamsStore.teamCosts[injectModal.teamId]
  const monthly = cost?.total_monthly_cost ?? 0
  if (monthly <= 0) return null
  return Math.round(monthly * injectModalMonthsToDelivery.value)
})

// --- Helpers ---
function formatCurrency(val: number): string {
  return '฿' + (val ?? 0).toLocaleString('th-TH', { minimumFractionDigits: 0, maximumFractionDigits: 0 })
}

function teamName(id: number): string {
  return teamsStore.teams.find(t => t.id === id)?.name ?? ''
}

function runwayColor(months: number): string {
  if (months <= 0) return 'text-gray-600'
  if (months > 2) return 'text-emerald-400'
  if (months > 1) return 'text-yellow-400'
  return 'text-red-400'
}

function runwayBarColor(months: number): string {
  if (months <= 0) return 'bg-gray-600'
  if (months > 2) return 'bg-emerald-500'
  if (months > 1) return 'bg-yellow-500'
  return 'bg-red-500'
}

function runwayBarWidth(months: number): string {
  const pct = Math.min((months / 3) * 100, 100)
  return pct + '%'
}

// --- Actions ---
function toggleExpand(teamId: number) {
  expandedTeam.value = expandedTeam.value === teamId ? null : teamId
}

async function createTeam() {
  if (!newTeamName.value.trim()) return
  isCreatingTeam.value = true
  createTeamError.value = ''
  try {
    await teamsStore.createTeam(newTeamName.value.trim())
    newTeamName.value = ''
    showCreateTeamModal.value = false
  } catch (e: unknown) {
    createTeamError.value = e instanceof Error ? e.message : 'Failed to create team'
  } finally {
    isCreatingTeam.value = false
  }
}

function deleteTeam(team: Team) {
  teamToDelete.value = team
}

async function confirmDeleteTeam() {
  if (!teamToDelete.value) return
  await teamsStore.deleteTeam(teamToDelete.value.id)
  teamToDelete.value = null
}

function confirmDisableFeatureTeam() {
  showDisableFeatureConfirm.value = true
}

async function doDisableFeatureTeam() {
  showDisableFeatureConfirm.value = false
  isDisablingFeature.value = true
  try {
    await teamsStore.setTeamsFeatureEnabled(false)
  } finally {
    isDisablingFeature.value = false
  }
}

async function enableFeatureTeam() {
  isEnablingFeature.value = true
  try {
    await teamsStore.setTeamsFeatureEnabled(true)
    await teamsStore.fetchTeams()
    for (const team of teamsStore.teams) {
      teamsStore.fetchTeamCost(team.id)
    }
  } finally {
    isEnablingFeature.value = false
  }
}

async function addMemberToTeam(teamId: number) {
  const userId = addMemberUserID.value[teamId]
  if (!userId) return
  await teamsStore.assignUserToTeam(Number(userId), teamId)
  addMemberUserID.value[teamId] = ''
  await teamsStore.fetchTeams()
  await loadAllUsers()
  // Refresh cost since team composition changed
  if (teamsStore.teamCosts[teamId]) await teamsStore.fetchTeamCost(teamId)
}

async function removeMemberFromTeam(userId: number) {
  await teamsStore.assignUserToTeam(userId, null)
  await teamsStore.fetchTeams()
  await loadAllUsers()
}

// --- Finance Actions ---
function openInjectModal(teamId: number) {
  injectModal.teamId = teamId
  injectModal.deliveryDate = ''
  injectModal.amount = null
  injectModal.note = ''
  injectModal.error = ''
  injectModal.loading = false
  injectModal.open = true
}

// When delivery date is selected, auto-fill amount so capital lasts until that date
watch(
  () => injectModal.deliveryDate,
  () => {
    const suggested = injectModalAmountFromDate.value
    if (suggested != null && suggested > 0) injectModal.amount = suggested
  }
)

async function confirmInjectCapital() {
  if (!injectModal.amount || injectModal.amount <= 0) return
  injectModal.loading = true
  injectModal.error = ''
  const currentBonus = teamsStore.teamCosts[injectModal.teamId]?.bonus_percentage ?? 0
  try {
    await teamsStore.injectCapital(injectModal.teamId, {
      amount: injectModal.amount,
      bonus_percentage: currentBonus,
      note: injectModal.note || (injectModal.deliveryDate ? `ส่งงาน ${injectModal.deliveryDate}` : 'Capital injection'),
    })
    injectModal.open = false
  } catch (e: unknown) {
    injectModal.error = e instanceof Error ? e.message : 'Failed to inject capital'
  } finally {
    injectModal.loading = false
  }
}

async function loadAllUsers() {
  try {
    const data = await fetchWithAuth<{ data: TeamUser[] }>('/auth/users')
    allUsers.value = data.data ?? []
  } catch (_e) {
    // non-fatal
  }
}

function openEditNameModal(team: Team) {
  editNameModal.teamId = team.id
  editNameModal.name = team.name
  editNameModal.error = ''
  editNameModal.loading = false
  editNameModal.open = true
}

async function confirmEditName() {
  if (!editNameModal.name.trim()) return
  editNameModal.loading = true
  editNameModal.error = ''
  const updated = await teamsStore.updateTeam(editNameModal.teamId, editNameModal.name.trim())
  editNameModal.loading = false
  if (updated) {
    editNameModal.open = false
  } else {
    editNameModal.error = teamsStore.error || 'Failed to update team name'
  }
}

function openEditCapitalModal(teamId: number) {
  const cost = teamsStore.teamCosts[teamId]
  editCapitalModal.teamId = teamId
  editCapitalModal.newBalance = cost?.capital_balance ?? null
  editCapitalModal.bonusPct = cost?.bonus_percentage ?? null
  editCapitalModal.note = ''
  editCapitalModal.error = ''
  editCapitalModal.loading = false
  editCapitalModal.open = true
}

async function confirmEditCapital() {
  if (editCapitalModal.newBalance === null || editCapitalModal.newBalance < 0) return
  editCapitalModal.loading = true
  editCapitalModal.error = ''
  try {
    const payload: { new_balance: number; bonus_percentage?: number; note: string } = {
      new_balance: Number(editCapitalModal.newBalance),
      note: editCapitalModal.note || 'Manual capital adjustment',
    }
    // Only send bonus_percentage if it's a valid number (not null/empty)
    const bp = editCapitalModal.bonusPct
    if (bp !== null && bp !== undefined && !isNaN(Number(bp))) {
      payload.bonus_percentage = Number(bp)
    }
    await teamsStore.editCapital(editCapitalModal.teamId, payload)
    editCapitalModal.open = false
  } catch (e: unknown) {
    editCapitalModal.error = e instanceof Error ? e.message : 'Failed to edit capital'
  } finally {
    editCapitalModal.loading = false
  }
}

onMounted(async () => {
  await teamsStore.fetchTeamsFeatureEnabled()
  if (teamsStore.teamsFeatureEnabled) {
    await Promise.all([teamsStore.fetchTeams(), loadAllUsers()])
    for (const team of teamsStore.teams) {
      teamsStore.fetchTeamCost(team.id)
    }
  } else {
    await loadAllUsers()
  }
})
</script>

<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900 p-6">
  <DisciplineDayDetailModal
    :show="modalOpen"
    :user-id="modalUserId"
    :user-email="modalUserEmail"
    :display-name="modalUserDisplayName"
    :date="modalDate"
    @close="closeDayDetail"
  />

    <!-- Header -->
    <div class="mb-6 border-b border-gray-700 pb-4">
      <div class="flex items-center justify-between flex-wrap gap-4">
        <div>
          <h1 class="text-2xl font-bold text-white flex items-center gap-2">
            <span class="text-orange-400">⚡</span>
            Discipline Tracker
          </h1>
          <p class="text-sm text-gray-400 mt-1">ติดตามวินัยการทำงานของพนักงานรายวัน — tasks, rework, logtime, Daily Pulse</p>
        </div>
        <!-- Date range picker -->
        <div class="flex items-center gap-3 flex-wrap">
          <UiDatePicker v-model="fromDate" placeholder="จาก…" />
          <span class="text-gray-500">→</span>
          <UiDatePicker v-model="toDate" placeholder="ถึง…" />
          <!-- Quick presets -->
          <div class="flex gap-1">
            <button
              v-for="preset in datePresets"
              :key="preset.label"
              @click="applyPreset(preset)"
              class="px-2.5 py-1.5 text-xs rounded-md border transition-colors"
              :class="activePreset === preset.label
                ? 'bg-orange-600 border-orange-500 text-white'
                : 'bg-gray-800 border-gray-700 text-gray-400 hover:text-white hover:border-gray-600'"
            >
              {{ preset.label }}
            </button>
          </div>
          <button
            @click="loadData"
            :disabled="store.disciplineLoading"
            class="px-4 py-2 bg-orange-600 hover:bg-orange-500 disabled:opacity-50 text-white text-sm rounded-lg font-medium transition-colors flex items-center gap-2"
          >
            <svg v-if="store.disciplineLoading" class="animate-spin w-4 h-4" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8H4z"/>
            </svg>
            <span>{{ store.disciplineLoading ? 'Loading...' : 'โหลดข้อมูล' }}</span>
          </button>
        </div>
      </div>
    </div>

    <!-- Error -->
    <div v-if="store.disciplineError" class="rounded-lg border border-red-500/50 bg-red-900/20 p-4 text-red-400 mb-6">
      {{ store.disciplineError }}
    </div>

    <!-- Empty state -->
    <div v-else-if="!store.discipline && !store.disciplineLoading" class="text-center py-24 text-gray-500">
      <div class="text-5xl mb-4">📊</div>
      <p class="text-lg">เลือกช่วงวันที่แล้วกด "โหลดข้อมูล"</p>
    </div>

    <!-- Data loaded -->
    <template v-else-if="store.discipline">
      <!-- Summary cards -->
      <div class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-6 gap-4 mb-6">
        <div class="bg-gray-800 border border-gray-700 rounded-xl p-4">
          <div class="text-xs text-gray-400 mb-1">พนักงานทั้งหมด</div>
          <div class="text-3xl font-bold text-white">{{ store.discipline.users.length }}</div>
        </div>
        <div class="bg-gray-800 border border-emerald-700/50 rounded-xl p-4">
          <div class="text-xs text-gray-400 mb-1">Job Done รวม</div>
          <div class="text-3xl font-bold text-emerald-400">{{ totalTasksClosed }}</div>
        </div>
        <div class="bg-gray-800 border border-red-700/50 rounded-xl p-4">
          <div class="text-xs text-gray-400 mb-1">Rework รวม</div>
          <div class="text-3xl font-bold text-red-400">{{ totalReworks }}</div>
        </div>
        <div class="bg-gray-800 border border-yellow-700/50 rounded-xl p-4">
          <div class="text-xs text-gray-400 mb-1">Missed Pulse</div>
          <div class="text-3xl font-bold text-yellow-400">{{ totalMissedPulse }}</div>
          <div class="text-xs text-gray-500 mt-0.5">ครั้ง (ทีมรวม)</div>
        </div>
        <div class="bg-gray-800 border border-rose-700/50 rounded-xl p-4">
          <div class="text-xs text-gray-400 mb-1">🕐 สายรวม</div>
          <div class="text-3xl font-bold text-rose-400">{{ totalLateDays }}</div>
          <div class="text-xs text-gray-500 mt-0.5">ครั้ง (ทีมรวม)</div>
        </div>
        <div class="bg-gray-800 border border-amber-700/50 rounded-xl p-4">
          <div class="text-xs text-gray-400 mb-1">🚪 กลับก่อนรวม</div>
          <div class="text-3xl font-bold text-amber-400">{{ totalEarlyCheckoutDays }}</div>
          <div class="text-xs text-gray-500 mt-0.5">ครั้ง (ทีมรวม)</div>
        </div>
      </div>

      <!-- Legend -->
      <div class="flex flex-wrap gap-4 mb-4 text-xs text-gray-400">
        <div class="flex items-center gap-1.5"><span class="w-3 h-3 rounded-sm bg-emerald-600 inline-block"></span>Job Done</div>
        <div class="flex items-center gap-1.5"><span class="w-3 h-3 rounded-sm bg-red-600 inline-block"></span>Rework</div>
        <div class="flex items-center gap-1.5"><span class="w-3 h-3 rounded-sm bg-blue-600 inline-block"></span>Hours logged</div>
        <div class="flex items-center gap-1.5"><span class="w-3 h-3 rounded-sm bg-violet-600 inline-block"></span>Daily Pulse ✓</div>
        <div class="flex items-center gap-1.5"><span class="w-3 h-3 rounded-sm bg-gray-700 border border-red-500/50 inline-block"></span>ไม่มี Pulse</div>
        <div class="flex items-center gap-1.5"><span class="text-rose-400">🕐</span>สาย (เข้างานช้า)</div>
        <div class="flex items-center gap-1.5"><span class="text-amber-400">🚪</span>กลับก่อนเวลา</div>
      </div>

      <!-- Main discipline grid -->
      <div class="overflow-auto rounded-xl border border-gray-700">
        <table class="w-full text-sm border-collapse min-w-max">
          <thead>
            <tr class="bg-gray-800 border-b border-gray-700">
              <!-- Sticky user column -->
              <th class="sticky left-0 z-10 bg-gray-800 text-left px-4 py-3 text-gray-400 font-medium border-r border-gray-700 min-w-[180px]">
                พนักงาน
              </th>
              <th class="px-3 py-3 text-center text-gray-400 font-medium border-r border-gray-700 whitespace-nowrap">
                สรุป
              </th>
              <!-- Date columns -->
              <th
                v-for="d in store.discipline.dates"
                :key="d"
                class="px-2 py-3 text-center text-gray-400 font-medium min-w-[90px] border-r border-gray-700/50 last:border-r-0"
              >
                <div class="text-[11px]">{{ formatDateHeader(d) }}</div>
                <div class="text-[10px] text-gray-600">{{ dayOfWeek(d) }}</div>
              </th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="(user, idx) in sortedUsers"
              :key="user.user_id"
              class="border-b border-gray-700/50 hover:bg-gray-800/50 transition-colors"
              :class="idx % 2 === 0 ? 'bg-gray-900' : 'bg-gray-900/60'"
            >
              <!-- User info -->
              <td class="sticky left-0 z-10 pl-2 pr-3 py-0 border-r border-gray-700 bg-inherit w-[245px] min-w-[245px] max-w-[245px]">
                <div class="flex items-center gap-3">
                  <div class="w-16 h-16 rounded-full shrink-0 overflow-hidden">
                    <img
                      v-if="user.user_avatar_url"
                      :src="user.user_avatar_url"
                      :alt="user.user_display_name || user.user_email"
                      class="w-full h-full object-cover"
                    />
                    <div
                      v-else
                      class="w-full h-full flex items-center justify-center text-lg font-bold text-white"
                      :style="{ background: avatarColor(user.user_email) }"
                    >{{ userInitial(user) }}</div>
                  </div>
                  <div class="min-w-0 overflow-hidden">
                    <div class="text-white font-medium truncate text-xs">
                      {{ user.user_display_name || user.user_email.split('@')[0] }}
                    </div>
                    <div class="text-gray-500 text-[10px] truncate">{{ user.user_email }}</div>
                    <span
                      class="inline-block text-[9px] px-1.5 py-0.5 rounded mt-0.5"
                      :class="roleBadgeClass(user.role)"
                    >{{ user.role }}</span>
                  </div>
                </div>
              </td>

              <!-- Summary column -->
              <td class="px-3 py-3 border-r border-gray-700 text-center">
                <div class="space-y-1">
                  <div class="flex items-center justify-between gap-2 text-xs">
                    <span class="text-gray-500">Job Done</span>
                    <span class="font-bold text-emerald-400">{{ user.total_tasks_closed }}</span>
                  </div>
                  <div v-if="(user.total_deployments ?? 0) > 0" class="flex items-center justify-between gap-2 text-xs">
                    <span class="text-gray-500">🚀 Deploy</span>
                    <span class="font-bold text-orange-400">{{ user.total_deployments }}</span>
                  </div>
                  <div class="flex items-center justify-between gap-2 text-xs">
                    <span class="text-gray-500">Rework</span>
                    <span class="font-bold" :class="user.total_reworks > 0 ? 'text-red-400' : 'text-gray-600'">{{ user.total_reworks }}</span>
                  </div>
                  <div class="flex items-center justify-between gap-2 text-xs">
                    <span class="text-gray-500">Hrs</span>
                    <span class="font-bold text-blue-400">{{ user.total_logged_hours.toFixed(1) }}</span>
                  </div>
                  <div class="flex items-center justify-between gap-2 text-xs">
                    <span class="text-gray-500">No Pulse</span>
                    <span class="font-bold" :class="user.missed_pulse_count > 0 ? 'text-yellow-400' : 'text-gray-600'">{{ user.missed_pulse_count }}</span>
                  </div>
                  <div v-if="(user.total_late_days ?? 0) > 0" class="flex items-center justify-between gap-2 text-xs">
                    <span class="text-gray-500">🕐 สาย</span>
                    <span class="font-bold text-rose-400">{{ user.total_late_days }}</span>
                  </div>
                  <div v-if="(user.total_early_checkout_days ?? 0) > 0" class="flex items-center justify-between gap-2 text-xs">
                    <span class="text-gray-500">🚪 กลับก่อน</span>
                    <span class="font-bold text-amber-400">{{ user.total_early_checkout_days }}</span>
                  </div>
                </div>
              </td>

              <!-- Daily cells -->
              <td
                v-for="day in user.days"
                :key="day.date"
                class="px-2 py-2 border-r border-gray-700/30 last:border-r-0 align-top"
              >
                <div
                  class="rounded-lg p-1.5 min-h-[70px] flex flex-col gap-1 cursor-pointer hover:ring-1 hover:ring-orange-500/50 hover:brightness-110 transition-all"
                  :class="dayCellBg(day)"
                  @click="openDayDetail(user, day.date)"
                >
                  <!-- Daily Pulse badge -->
                  <div class="flex justify-end">
                    <span
                      class="text-[9px] px-1 py-0.5 rounded font-medium"
                      :class="day.has_daily_pulse
                        ? 'bg-violet-900/60 text-violet-300'
                        : 'bg-red-900/40 text-red-400'"
                    >
                      {{ day.has_daily_pulse ? '✓ Pulse' : '✗ Pulse' }}
                    </span>
                  </div>
                  <!-- Attendance row: check-in/out + late/early badges -->
                  <div v-if="day.attendance_status" class="flex items-center gap-1 flex-wrap mb-0.5">
                    <span
                      class="text-[9px] px-1 py-0.5 rounded font-semibold"
                      :class="day.is_late
                        ? 'bg-rose-900/60 text-rose-300'
                        : 'bg-teal-900/50 text-teal-300'"
                    >
                      {{ day.is_late ? '🕐 สาย' : '✓ ตรงเวลา' }}
                    </span>
                    <span
                      v-if="day.early_checkout"
                      class="text-[9px] px-1 py-0.5 rounded font-semibold bg-amber-900/60 text-amber-300"
                    >🚪 ก่อนเวลา</span>
                  </div>
                  <!-- Check-in / check-out times -->
                  <div v-if="day.check_in_at || day.check_out_at" class="flex items-center gap-1 text-[9px] text-gray-500 mb-0.5">
                    <span v-if="day.check_in_at" class="text-gray-600 dark:text-gray-400">↑{{ day.check_in_at }}</span>
                    <span v-if="day.check_in_at && day.check_out_at" class="text-gray-700">·</span>
                    <span v-if="day.check_out_at" class="text-gray-600 dark:text-gray-400">↓{{ day.check_out_at }}</span>
                  </div>
                  <!-- Metrics -->
                  <div class="space-y-0.5">
                    <div v-if="day.tasks_closed > 0" class="flex items-center gap-1 text-[10px]">
                      <span class="w-1.5 h-1.5 rounded-full bg-emerald-500 shrink-0"></span>
                      <span class="text-emerald-300">{{ day.tasks_closed }} task{{ day.tasks_closed > 1 ? 's' : '' }}</span>
                    </div>
                    <div v-if="(day.deployments_completed ?? 0) > 0" class="flex items-center gap-1 text-[10px]">
                      <span class="text-orange-400 shrink-0">🚀</span>
                      <span class="text-orange-300">{{ day.deployments_completed }} deploy{{ day.deployments_completed > 1 ? 's' : '' }}</span>
                    </div>
                    <div v-if="day.reworks > 0" class="flex items-center gap-1 text-[10px]">
                      <span class="w-1.5 h-1.5 rounded-full bg-red-500 shrink-0"></span>
                      <span class="text-red-300">{{ day.reworks }} rework{{ day.reworks > 1 ? 's' : '' }}</span>
                    </div>
                    <div class="flex items-center gap-1 text-[10px]">
                      <span class="w-1.5 h-1.5 rounded-full shrink-0" :class="day.logged_minutes > 0 ? 'bg-blue-500' : 'bg-gray-700'"></span>
                      <span :class="day.logged_minutes > 0 ? 'text-blue-300' : 'text-gray-600'">
                        {{ (day.logged_minutes / 60).toFixed(1) }}h
                      </span>
                    </div>
                  </div>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Per-user detail cards (mobile-friendly alternative) -->
      <div class="mt-8">
        <h2 class="text-base font-semibold text-white mb-4">สรุปรายบุคคล</h2>
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
          <div
            v-for="user in sortedUsers"
            :key="'card-' + user.user_id"
            class="bg-gray-800 border rounded-xl p-4 transition-colors"
            :class="userCardBorderClass(user)"
          >
            <div class="flex items-center gap-3 mb-3">
              <div class="w-16 h-16 rounded-full shrink-0 overflow-hidden">
                <img
                  v-if="user.user_avatar_url"
                  :src="user.user_avatar_url"
                  :alt="user.user_display_name || user.user_email"
                  class="w-full h-full object-cover"
                />
                <div
                  v-else
                  class="w-full h-full flex items-center justify-center text-lg font-bold text-white"
                  :style="{ background: avatarColor(user.user_email) }"
                >{{ userInitial(user) }}</div>
              </div>
              <div class="min-w-0">
                <div class="text-white font-semibold text-sm truncate">
                  {{ user.user_display_name || user.user_email.split('@')[0] }}
                </div>
                <div class="text-gray-500 text-xs truncate">{{ user.user_email }}</div>
              </div>
            </div>
            <div class="gap-2 text-xs" :class="(user.total_deployments ?? 0) > 0 ? 'grid grid-cols-2' : 'grid grid-cols-2'">
              <div class="bg-gray-900/60 rounded-lg p-2 text-center">
                <div class="text-gray-400 mb-0.5">Job Done</div>
                <div class="text-xl font-bold text-emerald-400">{{ user.total_tasks_closed }}</div>
              </div>
              <div class="bg-gray-900/60 rounded-lg p-2 text-center">
                <div class="text-gray-400 mb-0.5">Rework</div>
                <div class="text-xl font-bold" :class="user.total_reworks > 0 ? 'text-red-400' : 'text-gray-600'">
                  {{ user.total_reworks }}
                </div>
              </div>
              <div class="bg-gray-900/60 rounded-lg p-2 text-center">
                <div class="text-gray-400 mb-0.5">Logged</div>
                <div class="text-xl font-bold text-blue-400">{{ user.total_logged_hours.toFixed(1) }}h</div>
              </div>
              <div class="bg-gray-900/60 rounded-lg p-2 text-center">
                <div class="text-gray-400 mb-0.5">No Pulse</div>
                <div class="text-xl font-bold" :class="user.missed_pulse_count > 0 ? 'text-yellow-400' : 'text-emerald-400'">
                  {{ user.missed_pulse_count }}
                </div>
              </div>
              <div v-if="(user.total_deployments ?? 0) > 0" class="col-span-2 bg-orange-950/30 border border-orange-800/30 rounded-lg p-2 text-center">
                <div class="text-gray-400 mb-0.5">🚀 Deployed</div>
                <div class="text-xl font-bold text-orange-400">{{ user.total_deployments }}</div>
              </div>
            </div>
            <!-- Attendance mini-summary (show only when there's data) -->
            <div
              v-if="(user.total_late_days ?? 0) > 0 || (user.total_early_checkout_days ?? 0) > 0"
              class="mt-2 flex gap-2"
            >
              <div v-if="(user.total_late_days ?? 0) > 0" class="flex-1 bg-rose-950/30 border border-rose-800/30 rounded-lg px-2 py-1.5 text-center">
                <div class="text-[10px] text-gray-500">🕐 สาย</div>
                <div class="text-base font-bold text-rose-400">{{ user.total_late_days }} ครั้ง</div>
              </div>
              <div v-if="(user.total_early_checkout_days ?? 0) > 0" class="flex-1 bg-amber-950/30 border border-amber-800/30 rounded-lg px-2 py-1.5 text-center">
                <div class="text-[10px] text-gray-500">🚪 กลับก่อน</div>
                <div class="text-base font-bold text-amber-400">{{ user.total_early_checkout_days }} ครั้ง</div>
              </div>
            </div>
            <!-- Discipline score bar -->
            <div class="mt-3">
              <div class="flex justify-between text-[10px] text-gray-500 mb-1">
                <span>Discipline Score</span>
                <span :class="disciplineScoreColor(disciplineScore(user))">{{ disciplineScore(user) }}%</span>
              </div>
              <div class="h-1.5 bg-gray-700 rounded-full overflow-hidden">
                <div
                  class="h-full rounded-full transition-all duration-500"
                  :class="disciplineBarClass(disciplineScore(user))"
                  :style="{ width: disciplineScore(user) + '%' }"
                ></div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { usePerformanceStore } from '~/core/modules/performance/performance-store'
import type { DisciplineUser, DisciplineUserDayStat } from '~/core/modules/performance/performance-api'

// ─── Day Detail Modal state ───────────────────────────────────────────────────
const modalOpen = ref(false)
const modalUserId = ref(0)
const modalUserEmail = ref('')
const modalUserDisplayName = ref('')
const modalDate = ref('')

function openDayDetail(user: DisciplineUser, date: string) {
  modalUserId.value = user.user_id
  modalUserEmail.value = user.user_email
  modalUserDisplayName.value = user.user_display_name || ''
  modalDate.value = date
  modalOpen.value = true
}

function closeDayDetail() {
  modalOpen.value = false
}

definePageMeta({ middleware: 'auth', layout: 'default' })

const store = usePerformanceStore()

// ─── Date helpers ─────────────────────────────────────────────────────────────
import { localDateStr } from '~/composables/useLocalDate'

function todayStr() {
  return localDateStr()
}

function daysAgoStr(n: number) {
  const d = new Date()
  d.setDate(d.getDate() - n)
  return localDateStr(d)
}

const fromDate = ref(daysAgoStr(6))
const toDate = ref(todayStr())
const activePreset = ref('7 วัน')

const datePresets = [
  { label: '7 วัน', from: () => daysAgoStr(6), to: () => todayStr() },
  { label: '14 วัน', from: () => daysAgoStr(13), to: () => todayStr() },
  { label: '30 วัน', from: () => daysAgoStr(29), to: () => todayStr() },
  { label: 'สัปดาห์นี้', from: () => startOfWeekStr(), to: () => todayStr() },
]

function startOfWeekStr() {
  const d = new Date()
  const day = d.getDay()
  const diff = d.getDate() - day + (day === 0 ? -6 : 1)
  d.setDate(diff)
  return localDateStr(d)
}

function applyPreset(preset: typeof datePresets[0]) {
  fromDate.value = preset.from()
  toDate.value = preset.to()
  activePreset.value = preset.label
  loadData()
}

async function loadData() {
  activePreset.value = ''
  await store.fetchDiscipline(fromDate.value, toDate.value)
}

function formatDateHeader(d: string) {
  const date = new Date(d + 'T00:00:00')
  return date.toLocaleDateString('th-TH', { day: '2-digit', month: 'short' })
}

function dayOfWeek(d: string) {
  const days = ['อา', 'จ', 'อ', 'พ', 'พฤ', 'ศ', 'ส']
  return days[new Date(d + 'T00:00:00').getDay()]
}

// ─── Computed ─────────────────────────────────────────────────────────────────

const totalTasksClosed = computed(() =>
  store.discipline?.users.reduce((s, u) => s + u.total_tasks_closed, 0) ?? 0
)
const totalDeployments = computed(() =>
  store.discipline?.users.reduce((s, u) => s + (u.total_deployments ?? 0), 0) ?? 0
)
const totalReworks = computed(() =>
  store.discipline?.users.reduce((s, u) => s + u.total_reworks, 0) ?? 0
)
const totalMissedPulse = computed(() =>
  store.discipline?.users.reduce((s, u) => s + u.missed_pulse_count, 0) ?? 0
)
const totalLateDays = computed(() =>
  store.discipline?.users.reduce((s, u) => s + (u.total_late_days ?? 0), 0) ?? 0
)
const totalEarlyCheckoutDays = computed(() =>
  store.discipline?.users.reduce((s, u) => s + (u.total_early_checkout_days ?? 0), 0) ?? 0
)

const sortedUsers = computed(() => {
  if (!store.discipline) return []
  return [...store.discipline.users].sort((a, b) => disciplineScore(b) - disciplineScore(a))
})

// ─── Discipline score (0–100) ─────────────────────────────────────────────────

function disciplineScore(user: DisciplineUser): number {
  const days = store.discipline?.dates.length || 1
  // Pulse adherence 40%, logtime activity 40%, no rework 20%
  const pulsePct = Math.max(0, (days - user.missed_pulse_count) / days) * 40
  const logDays = user.days.filter(d => d.logged_minutes > 0).length
  const logPct = (logDays / days) * 40
  const totalSubs = user.total_tasks_closed + user.total_reworks
  const reworkPct = totalSubs > 0 ? (1 - user.total_reworks / totalSubs) * 20 : 20
  return Math.round(pulsePct + logPct + reworkPct)
}

// ─── Styling helpers ──────────────────────────────────────────────────────────

function dayCellBg(day: DisciplineUserDayStat): string {
  if (!day.has_daily_pulse && day.logged_minutes === 0 && day.tasks_closed === 0 && !day.attendance_status) {
    return 'bg-red-50 dark:bg-gray-800/40 border border-red-200 dark:border-red-900/30'
  }
  if (day.reworks > 0) return 'bg-red-50 dark:bg-red-950/30 border border-red-200 dark:border-red-700/20'
  if (day.is_late) return 'bg-rose-50 dark:bg-rose-950/25 border border-rose-200 dark:border-rose-800/30'
  if (day.early_checkout) return 'bg-amber-50 dark:bg-amber-950/20 border border-amber-200 dark:border-amber-800/25'
  if (day.logged_minutes > 0 || day.tasks_closed > 0) return 'bg-white dark:bg-gray-800/60 border border-gray-200 dark:border-gray-700/30'
  return 'bg-gray-100 dark:bg-gray-800/20 border border-gray-200 dark:border-gray-800'
}

function userCardBorderClass(user: DisciplineUser): string {
  const score = disciplineScore(user)
  if (score >= 80) return 'border-emerald-700/40'
  if (score >= 50) return 'border-yellow-700/40'
  return 'border-red-700/40'
}

function disciplineScoreColor(score: number): string {
  if (score >= 80) return 'text-emerald-400'
  if (score >= 50) return 'text-yellow-400'
  return 'text-red-400'
}

function disciplineBarClass(score: number): string {
  if (score >= 80) return 'bg-emerald-500'
  if (score >= 50) return 'bg-yellow-500'
  return 'bg-red-500'
}

function roleBadgeClass(role: string): string {
  const map: Record<string, string> = {
    CEO: 'bg-purple-900/60 text-purple-300',
    PRODUCT_OWNER: 'bg-blue-900/60 text-blue-300',
    PM: 'bg-blue-900/60 text-blue-300',
    ENGINEER: 'bg-gray-700/80 text-gray-300',
    CHIEF_ENGINEER: 'bg-gray-700/80 text-gray-300',
    DEV: 'bg-gray-700/80 text-gray-300',
    MANAGER: 'bg-indigo-900/60 text-indigo-300',
    SUPPORT: 'bg-green-900/60 text-green-300',
  }
  return map[role] || 'bg-gray-700 text-gray-400'
}

function userInitial(user: DisciplineUser): string {
  const name = user.user_display_name || user.user_email
  return name.charAt(0).toUpperCase()
}

const avatarColors = [
  'linear-gradient(135deg,#6366f1,#8b5cf6)',
  'linear-gradient(135deg,#ec4899,#f43f5e)',
  'linear-gradient(135deg,#14b8a6,#06b6d4)',
  'linear-gradient(135deg,#f59e0b,#f97316)',
  'linear-gradient(135deg,#10b981,#059669)',
  'linear-gradient(135deg,#3b82f6,#6366f1)',
]
function avatarColor(email: string): string {
  let hash = 0
  for (const c of email) hash = (hash * 31 + c.charCodeAt(0)) & 0xffffffff
  return avatarColors[Math.abs(hash) % avatarColors.length]
}

// ─── Auto-load on mount ───────────────────────────────────────────────────────

onMounted(() => {
  loadData()
})
</script>

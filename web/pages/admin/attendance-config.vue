<template>
  <div class="min-h-screen bg-gray-900 text-gray-100 p-6">
    <div v-if="!allowed" class="rounded-xl border border-red-800 bg-red-950/30 p-6 text-red-200">
      Access denied. CEO, Manager, or Support only. / สิทธิ์เฉพาะ CEO, Manager หรือ Support
    </div>

    <template v-else>
      <header class="mb-6 border-b border-gray-800 pb-6">
        <h1 class="text-2xl font-bold text-white">Attendance Admin</h1>
        <p class="text-sm text-gray-400 mt-1">
          Dashboard สำหรับดู check-in/check-out รายวัน และ Config สำหรับตั้งค่า Office attendance
        </p>
      </header>

      <div class="mb-6 flex gap-2 border-b border-gray-800 pb-3">
        <button
          type="button"
          class="px-4 py-2 rounded-lg text-sm font-medium transition"
          :class="activeTab === 'dashboard' ? 'bg-emerald-600 text-white' : 'bg-gray-800 text-gray-300 hover:bg-gray-700'"
          @click="activeTab = 'dashboard'"
        >
          Dashboard
        </button>
        <button
          type="button"
          class="px-4 py-2 rounded-lg text-sm font-medium transition"
          :class="activeTab === 'config' ? 'bg-violet-600 text-white' : 'bg-gray-800 text-gray-300 hover:bg-gray-700'"
          @click="activeTab = 'config'"
        >
          Config
        </button>
      </div>

      <section v-if="activeTab === 'dashboard'" class="space-y-5">
        <div class="rounded-xl border border-gray-800 bg-gray-850/70 p-4">
          <div class="flex flex-wrap items-end gap-3">
            <label class="block">
              <span class="text-xs text-gray-400">วันที่</span>
              <input v-model="recordsDate" type="date" class="mt-1 rounded-lg bg-gray-800 border border-gray-600 px-3 py-2 text-white text-sm" />
            </label>
            <button
              type="button"
              class="h-10 px-4 rounded-lg bg-emerald-600 hover:bg-emerald-500 text-white text-sm font-medium disabled:opacity-50"
              :disabled="recordsLoading"
              @click="loadRecords"
            >
              {{ recordsLoading ? 'Loading…' : 'Load records' }}
            </button>
          </div>
        </div>

        <div class="grid gap-3 sm:grid-cols-3">
          <div class="rounded-xl border border-gray-800 bg-gray-850/70 p-4">
            <p class="text-xs text-gray-400">Total</p>
            <p class="text-2xl font-semibold text-white mt-1">{{ adminRecords.length }}</p>
          </div>
          <div class="rounded-xl border border-gray-800 bg-gray-850/70 p-4">
            <p class="text-xs text-gray-400">Checked in</p>
            <p class="text-2xl font-semibold text-emerald-300 mt-1">{{ checkedInCount }}</p>
          </div>
          <div class="rounded-xl border border-gray-800 bg-gray-850/70 p-4">
            <p class="text-xs text-gray-400">Checked out</p>
            <p class="text-2xl font-semibold text-violet-300 mt-1">{{ checkedOutCount }}</p>
          </div>
        </div>

        <div class="rounded-xl border border-gray-800 bg-gray-850/70 overflow-hidden">
          <div class="overflow-x-auto">
            <table class="min-w-full text-sm">
              <thead class="bg-gray-800/80 text-gray-300">
                <tr>
                  <th class="text-left px-4 py-3 font-medium">Employee</th>
                  <th class="text-left px-4 py-3 font-medium">Check in</th>
                  <th class="text-left px-4 py-3 font-medium">Check out</th>
                  <th class="text-left px-4 py-3 font-medium">Status</th>
                  <th class="text-right px-4 py-3 font-medium">Action</th>
                </tr>
              </thead>
              <tbody>
                <tr v-if="recordsLoading">
                  <td colspan="5" class="px-4 py-8 text-center text-gray-400">Loading records…</td>
                </tr>
                <tr v-else-if="!adminRecords.length">
                  <td colspan="5" class="px-4 py-8 text-center text-gray-500">No records for selected date.</td>
                </tr>
                <tr v-for="r in adminRecords" :key="r.id" class="border-t border-gray-800">
                  <td class="px-4 py-3 text-gray-100">
                    {{ r.user_display_name || r.user_email || `User ${r.user_id}` }}
                  </td>
                  <td class="px-4 py-3 text-gray-300">{{ formatTime(r.check_in_at) }}</td>
                  <td class="px-4 py-3 text-gray-300">{{ formatTime(r.check_out_at) }}</td>
                  <td class="px-4 py-3">
                    <span class="inline-flex items-center rounded-full px-2.5 py-1 text-xs font-medium"
                      :class="statusPillClass(r.status)">
                      {{ r.status }}<span v-if="r.is_late"> · late</span>
                    </span>
                  </td>
                  <td class="px-4 py-3 text-right">
                    <button
                      type="button"
                      class="px-3 py-1.5 rounded-md bg-rose-700 hover:bg-rose-600 text-white text-xs font-medium disabled:opacity-50"
                      :disabled="recordsLoading || deletingId === r.id"
                      @click="onDeleteRecord(r.id)"
                    >
                      {{ deletingId === r.id ? 'Deleting…' : 'Delete' }}
                    </button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </section>

      <section v-else class="space-y-6">
        <div class="rounded-xl border border-gray-800 bg-gray-850/70 p-5">
          <h2 class="text-lg font-semibold text-white mb-1">Attendance configuration</h2>
          <p class="text-sm text-gray-400">
            Onsite days require GPS within office radius. WFH days still require check-in &amp; check-out but from any location.
            Days: 1=Mon … 7=Sun.
          </p>
        </div>

        <form class="max-w-2xl space-y-6" @submit.prevent="onSave">
          <div class="grid gap-4 sm:grid-cols-2">
            <label class="block sm:col-span-2">
              <span class="text-xs text-gray-400">Office name</span>
              <input v-model="form.name" required class="mt-1 w-full rounded-lg bg-gray-800 border border-gray-600 px-3 py-2 text-white" />
            </label>
            <label class="block">
              <span class="text-xs text-gray-400">Latitude</span>
              <input v-model.number="form.latitude" type="number" step="any" required class="mt-1 w-full rounded-lg bg-gray-800 border border-gray-600 px-3 py-2 text-white" />
            </label>
            <label class="block">
              <span class="text-xs text-gray-400">Longitude</span>
              <input v-model.number="form.longitude" type="number" step="any" required class="mt-1 w-full rounded-lg bg-gray-800 border border-gray-600 px-3 py-2 text-white" />
            </label>
            <label class="block">
              <span class="text-xs text-gray-400">Radius (meters)</span>
              <input v-model.number="form.radius_meters" type="number" min="1" max="5000" required class="mt-1 w-full rounded-lg bg-gray-800 border border-gray-600 px-3 py-2 text-white" />
            </label>
            <label class="flex items-center gap-2 mt-6">
              <input v-model="form.is_active" type="checkbox" class="rounded border-gray-600" />
              <span class="text-sm text-gray-300">Active (only one should be active)</span>
            </label>
          </div>


          <div class="grid gap-4 sm:grid-cols-2">
            <label class="block">
              <span class="text-xs text-gray-400">Work start (HH:MM)</span>
              <input v-model="form.work_start_time" required pattern="^([01]?[0-9]|2[0-3]):[0-5][0-9](:[0-5][0-9])?$" class="mt-1 w-full rounded-lg bg-gray-800 border border-gray-600 px-3 py-2 text-white" />
            </label>
            <label class="block">
              <span class="text-xs text-gray-400">Work end (HH:MM)</span>
              <input v-model="form.work_end_time" required pattern="^([01]?[0-9]|2[0-3]):[0-5][0-9](:[0-5][0-9])?$" class="mt-1 w-full rounded-lg bg-gray-800 border border-gray-600 px-3 py-2 text-white" />
            </label>
          </div>

          <fieldset class="rounded-lg border border-gray-700 p-4">
            <legend class="text-sm text-gray-400 px-1">Office days (onsite — geofence required)</legend>
            <div class="flex flex-wrap gap-3 mt-2">
              <label v-for="d in dayChoices" :key="`o-${d.v}`" class="flex items-center gap-2 text-sm text-gray-300">
                <input type="checkbox" :checked="form.work_days.includes(d.v)" @change="toggleWorkDay(d.v, ($event.target as HTMLInputElement).checked)" />
                {{ d.label }}
              </label>
            </div>
          </fieldset>

          <fieldset class="rounded-lg border border-sky-900/60 bg-sky-950/20 p-4">
            <legend class="text-sm text-sky-300/90 px-1">WFH days (check-in / check-out from anywhere)</legend>
            <p class="text-xs text-gray-500 mb-2">Same work hours apply (late / early rules). If a day is both office and WFH, WFH wins (no geofence).</p>
            <div class="flex flex-wrap gap-3 mt-1">
              <label v-for="d in dayChoices" :key="`w-${d.v}`" class="flex items-center gap-2 text-sm text-gray-300">
                <input type="checkbox" :checked="form.wfh_days.includes(d.v)" @change="toggleWfhDay(d.v, ($event.target as HTMLInputElement).checked)" />
                {{ d.label }}
              </label>
            </div>
          </fieldset>

          <p v-if="store.error" class="text-sm text-red-400">{{ store.error }}</p>

          <button
            type="submit"
            :disabled="store.actionLoading || (form.work_days.length === 0 && form.wfh_days.length === 0)"
            class="px-6 py-2.5 rounded-lg bg-gradient-to-r from-emerald-600 to-teal-600 text-white font-medium disabled:opacity-40"
          >
            {{ store.actionLoading ? 'Saving…' : 'Save configuration' }}
          </button>
        </form>
      </section>
    </template>
  </div>
</template>

<script setup lang="ts">
import { useAttendanceApi, type AttendanceRecord } from '~/core/modules/attendance/infrastructure/attendance-api'
import { useAttendanceStore } from '~/core/modules/attendance/store/attendance-store'

definePageMeta({ layout: 'default', middleware: 'auth' })

const { currentUser } = useAuth()
const store = useAttendanceStore()
const api = useAttendanceApi()

const activeTab = ref<'dashboard' | 'config'>('dashboard')

const allowed = computed(() => {
  const r = currentUser.value?.role
  return r === 'CEO' || r === 'MANAGER' || r === 'SUPPORT'
})

const dayChoices = [
  { v: 1, label: 'Mon' },
  { v: 2, label: 'Tue' },
  { v: 3, label: 'Wed' },
  { v: 4, label: 'Thu' },
  { v: 5, label: 'Fri' },
  { v: 6, label: 'Sat' },
  { v: 7, label: 'Sun' },
]

const form = reactive({
  name: 'Main Office',
  latitude: 13.7563,
  longitude: 100.5018,
  radius_meters: 150,
  work_start_time: '09:00',
  work_end_time: '18:00',
  work_days: [1, 2, 3, 4, 5] as number[],
  wfh_days: [] as number[],
  is_active: true,
})

const recordsDate = ref(new Date().toISOString().slice(0, 10))
const adminRecords = ref<AttendanceRecord[]>([])
const recordsLoading = ref(false)
const deletingId = ref<number | null>(null)

const checkedInCount = computed(() => adminRecords.value.filter((r) => !!r.check_in_at).length)
const checkedOutCount = computed(() => adminRecords.value.filter((r) => !!r.check_out_at).length)

function toggleWorkDay(v: number, on: boolean) {
  if (on && !form.work_days.includes(v)) {
    form.work_days = [...form.work_days, v].sort((a, b) => a - b)
  }
  if (!on) {
    form.work_days = form.work_days.filter((x) => x !== v)
  }
}

function toggleWfhDay(v: number, on: boolean) {
  if (on && !form.wfh_days.includes(v)) {
    form.wfh_days = [...form.wfh_days, v].sort((a, b) => a - b)
  }
  if (!on) {
    form.wfh_days = form.wfh_days.filter((x) => x !== v)
  }
}

function parseWorkDays(raw: unknown): number[] {
  if (Array.isArray(raw)) {
    return raw.map((x) => Number(x)).filter((n) => n >= 1 && n <= 7)
  }
  return [1, 2, 3, 4, 5]
}

function parseWfhDays(raw: unknown): number[] {
  if (Array.isArray(raw)) {
    return raw.map((x) => Number(x)).filter((n) => n >= 1 && n <= 7)
  }
  return []
}

function formatTime(v?: string | null): string {
  if (!v) return '—'
  return new Date(v).toLocaleTimeString()
}

function statusPillClass(status: string): string {
  const s = (status || '').toLowerCase()
  if (s === 'present') return 'bg-emerald-900/60 text-emerald-200 border border-emerald-700/60'
  if (s === 'late') return 'bg-amber-900/60 text-amber-200 border border-amber-700/60'
  if (s === 'absent') return 'bg-rose-900/60 text-rose-200 border border-rose-700/60'
  return 'bg-gray-800 text-gray-200 border border-gray-700'
}

onMounted(async () => {
  if (!allowed.value) return
  await store.adminLoadConfig()
  const c = store.officeConfig
  if (c) {
    form.name = c.name
    form.latitude = c.latitude
    form.longitude = c.longitude
    form.radius_meters = c.radius_meters
    form.work_start_time = c.work_start_time?.slice(0, 5) ?? '09:00'
    form.work_end_time = c.work_end_time?.slice(0, 5) ?? '18:00'
    form.work_days = parseWorkDays(c.work_days)
    form.wfh_days = parseWfhDays(c.wfh_days)
    form.is_active = c.is_active
  }
  await loadRecords()
})

async function onSave() {
  const ok = await store.adminSaveConfig({
    name: form.name,
    latitude: form.latitude,
    longitude: form.longitude,
    radius_meters: form.radius_meters,
    work_start_time: form.work_start_time,
    work_end_time: form.work_end_time,
    work_days: [...form.work_days],
    wfh_days: [...form.wfh_days],
    is_active: form.is_active,
  })
  if (ok) {
    activeTab.value = 'dashboard'
    await loadRecords()
  }
}

async function onDeleteRecord(id: number) {
  const ok = window.confirm('Delete this attendance record? This action cannot be undone.')
  if (!ok) return

  deletingId.value = id
  try {
    await api.adminDeleteRecord(id)
    await loadRecords()
  } catch {
    // keep existing list on failure
  } finally {
    deletingId.value = null
  }
}

async function loadRecords() {
  recordsLoading.value = true
  try {
    const res = await api.adminRecords(recordsDate.value)
    adminRecords.value = res.records || []
  } catch {
    adminRecords.value = []
  } finally {
    recordsLoading.value = false
  }
}
</script>

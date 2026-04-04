<template>
  <div class="min-h-screen bg-gray-900 text-gray-100 p-6">
    <div v-if="!allowed" class="rounded-xl border border-red-800 bg-red-950/30 p-6 text-red-200">
      Access denied. CEO or Manager only. / สิทธิ์เฉพาะ CEO หรือ Manager
    </div>

    <template v-else>
      <header class="mb-8 border-b border-gray-800 pb-6">
        <h1 class="text-2xl font-bold text-white">Attendance configuration</h1>
        <p class="text-sm text-gray-400 mt-1">
          Onsite days require GPS / office network. WFH days still require check-in &amp; check-out but from any location. Days: 1=Mon … 7=Sun.
        </p>
      </header>

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

        <label class="block">
          <span class="text-xs text-gray-400">Allowed IPs / CIDRs (one per line, office public IP)</span>
          <textarea
            v-model="allowedIPsText"
            rows="4"
            class="mt-1 w-full rounded-lg bg-gray-800 border border-gray-600 px-3 py-2 text-white font-mono text-sm"
            placeholder="203.0.113.0/24&#10;198.51.100.55"
          />
        </label>

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
            <label v-for="d in dayChoices" :key="'o-' + d.v" class="flex items-center gap-2 text-sm text-gray-300">
              <input type="checkbox" :checked="form.work_days.includes(d.v)" @change="toggleWorkDay(d.v, ($event.target as HTMLInputElement).checked)" />
              {{ d.label }}
            </label>
          </div>
        </fieldset>

        <fieldset class="rounded-lg border border-sky-900/60 bg-sky-950/20 p-4">
          <legend class="text-sm text-sky-300/90 px-1">WFH days (check-in / check-out from anywhere)</legend>
          <p class="text-xs text-gray-500 mb-2">Same work hours apply (late / early rules). If a day is both office and WFH, WFH wins (no geofence).</p>
          <div class="flex flex-wrap gap-3 mt-1">
            <label v-for="d in dayChoices" :key="'w-' + d.v" class="flex items-center gap-2 text-sm text-gray-300">
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

      <section class="mt-12 max-w-3xl">
        <h2 class="text-lg font-semibold text-white mb-3">Records by date</h2>
        <div class="flex gap-2 mb-4">
          <input v-model="recordsDate" type="date" class="rounded-lg bg-gray-800 border border-gray-600 px-3 py-2 text-white text-sm" />
          <button
            type="button"
            class="px-4 py-2 rounded-lg bg-gray-700 text-white text-sm"
            :disabled="recordsLoading"
            @click="loadRecords"
          >
            Load
          </button>
        </div>
        <div v-if="recordsLoading" class="text-gray-400 text-sm">Loading…</div>
        <ul v-else-if="adminRecords.length" class="space-y-2 text-sm">
          <li
            v-for="r in adminRecords"
            :key="r.id"
            class="border border-gray-700 rounded-lg p-3 flex flex-wrap justify-between gap-2"
          >
            <span class="text-gray-200">{{ r.user_display_name || r.user_email || 'User ' + r.user_id }}</span>
            <span class="text-gray-400">
              {{ r.check_in_at ? new Date(r.check_in_at).toLocaleTimeString() : '—' }}
              →
              {{ r.check_out_at ? new Date(r.check_out_at).toLocaleTimeString() : '—' }}
            </span>
            <span class="text-xs text-gray-500">{{ r.status }} {{ r.is_late ? '· late' : '' }}</span>
          </li>
        </ul>
        <p v-else class="text-gray-500 text-sm">No records for this date.</p>
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

const allowed = computed(() => {
  const r = currentUser.value?.role
  return r === 'CEO' || r === 'MANAGER'
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

const allowedIPsText = ref('')
const recordsDate = ref(new Date().toISOString().slice(0, 10))
const adminRecords = ref<AttendanceRecord[]>([])
const recordsLoading = ref(false)

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
    allowedIPsText.value = (c.allowed_ips || []).join('\n')
  }
  await loadRecords()
})

async function onSave() {
  const ips = allowedIPsText.value
    .split(/[\n,]+/)
    .map((s) => s.trim())
    .filter(Boolean)
  const ok = await store.adminSaveConfig({
    name: form.name,
    latitude: form.latitude,
    longitude: form.longitude,
    radius_meters: form.radius_meters,
    allowed_ips: ips,
    work_start_time: form.work_start_time,
    work_end_time: form.work_end_time,
    work_days: [...form.work_days],
    wfh_days: [...form.wfh_days],
    is_active: form.is_active,
  })
  if (ok) {
    await loadRecords()
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

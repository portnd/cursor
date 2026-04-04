<template>
  <div class="min-h-screen bg-gray-900 text-gray-100 p-6">
    <header class="mb-8 border-b border-gray-800 pb-6">
      <div class="flex items-center gap-3">
        <div class="flex h-12 w-12 items-center justify-center rounded-xl bg-gradient-to-br from-emerald-600 to-teal-700 text-2xl shadow-lg">
          📍
        </div>
        <div>
          <h1 class="text-2xl font-bold text-white">Office attendance</h1>
          <p class="text-sm text-gray-400 mt-1">
            Onsite days: office GPS or network. WFH days: check-in/out from anywhere. / วันเข้าออฟฟิศต้องอยู่ในรัศมีหรือเน็ตออฟฟิศ · วัน WFH เช็คอินได้ทุกที่
          </p>
        </div>
      </div>
    </header>

    <div v-if="store.loading && !store.todayRecord && !store.officeConfig" class="text-gray-400 text-sm">
      Loading…
    </div>

    <div v-else class="max-w-3xl space-y-6">
      <div
        v-if="!store.officeConfig"
        class="rounded-xl border border-amber-800/60 bg-amber-950/30 p-4 text-amber-200 text-sm"
      >
        Office attendance is not configured yet. Ask a CEO or Manager to set location and hours in Admin → Attendance config.
      </div>

      <div
        v-else
        class="rounded-xl border border-gray-700 bg-gray-800/50 p-4 text-sm text-gray-300 space-y-1"
      >
        <p><span class="text-gray-500">Office</span> {{ store.officeConfig.name }}</p>
        <p><span class="text-gray-500">Hours</span> {{ store.officeConfig.work_start_time?.slice(0, 5) }} – {{ store.officeConfig.work_end_time?.slice(0, 5) }}</p>
        <p><span class="text-gray-500">Radius</span> {{ store.officeConfig.radius_meters }} m</p>
        <p v-if="todayKind === 'wfh'" class="text-sky-300 text-xs mt-2">Today: WFH — no geofence.</p>
        <p v-else-if="todayKind === 'office'" class="text-amber-200/90 text-xs mt-2">Today: onsite — geofence required.</p>
        <p v-else class="text-gray-500 text-xs mt-2">Today: no check-in required (off schedule).</p>
      </div>

      <div
        v-if="store.todayRecord?.check_in_at"
        class="rounded-xl border border-emerald-800/50 bg-emerald-950/20 p-4 text-sm space-y-1"
      >
        <p class="text-emerald-300 font-medium">Today’s record</p>
        <p class="text-gray-300">Check-in: {{ fmtTime(store.todayRecord.check_in_at) }}</p>
        <p v-if="store.todayRecord.check_out_at" class="text-gray-300">
          Check-out: {{ fmtTime(store.todayRecord.check_out_at) }}
        </p>
        <p v-else class="text-amber-300">Not checked out yet.</p>
        <p class="text-gray-400">
          Status: {{ store.todayRecord.status }}
          <span v-if="store.todayRecord.is_late"> (late)</span>
          <span v-if="store.todayRecord.early_checkout"> (early checkout)</span>
        </p>
        <p class="text-gray-500 text-xs">Method: {{ store.todayRecord.check_in_method || '—' }}</p>
      </div>

      <CheckInPanel v-if="store.officeConfig" />

      <section class="rounded-xl border border-gray-700 bg-gray-800/40 p-6">
        <div class="flex items-center justify-between mb-4">
          <h2 class="text-lg font-semibold text-white">History</h2>
          <button
            type="button"
            class="text-xs text-violet-400 hover:text-violet-300"
            @click="store.fetchHistory(false)"
          >
            Refresh
          </button>
        </div>
        <ul v-if="store.history.length" class="space-y-2">
          <li
            v-for="row in store.history"
            :key="row.id"
            class="flex flex-wrap items-center justify-between gap-2 border-b border-gray-700/80 py-2 text-sm"
          >
            <span class="text-gray-200">{{ fmtDate(row.attendance_date) }}</span>
            <span class="text-gray-400">
              in {{ row.check_in_at ? fmtTime(row.check_in_at) : '—' }}
              · out {{ row.check_out_at ? fmtTime(row.check_out_at) : '—' }}
            </span>
            <span class="text-xs text-gray-500">{{ row.status }}</span>
          </li>
        </ul>
        <p v-else class="text-sm text-gray-500">No history yet.</p>
        <button
          v-if="store.historyHasMore"
          type="button"
          class="mt-4 w-full py-2 rounded-lg border border-gray-600 text-gray-300 text-sm hover:bg-gray-700/50"
          :disabled="store.loading"
          @click="store.loadMoreHistory()"
        >
          Load more
        </button>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import CheckInPanel from '~/core/modules/attendance/ui/CheckInPanel.vue'
import { useAttendanceStore } from '~/core/modules/attendance/store/attendance-store'

definePageMeta({ layout: 'default', middleware: 'auth' })

const store = useAttendanceStore()

function isoWeekdayToday(): number {
  const w = new Date().getDay()
  return w === 0 ? 7 : w
}

const todayKind = computed<'wfh' | 'office' | 'off'>(() => {
  const c = store.officeConfig
  if (!c) return 'off'
  const wd = isoWeekdayToday()
  const wfh = Array.isArray(c.wfh_days) ? c.wfh_days : []
  const office = Array.isArray(c.work_days) ? c.work_days : []
  if (wfh.includes(wd)) return 'wfh'
  if (office.includes(wd)) return 'office'
  return 'off'
})

function fmtDate(iso: string) {
  if (!iso) return '—'
  return iso.slice(0, 10)
}

function fmtTime(iso: string) {
  if (!iso) return '—'
  try {
    return new Date(iso).toLocaleString(undefined, {
      dateStyle: 'short',
      timeStyle: 'short',
    })
  } catch {
    return iso
  }
}

onMounted(async () => {
  await store.fetchToday()
  await store.fetchHistory(false)
})
</script>

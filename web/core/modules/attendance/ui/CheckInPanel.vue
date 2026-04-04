<template>
  <div class="rounded-xl border border-gray-700 bg-gray-800/80 p-6 space-y-4">
    <h2 class="text-lg font-semibold text-white">Check-in / check-out</h2>
    <p v-if="isWfhToday" class="text-sm text-sky-300/90">
      Today is a <strong>WFH</strong> day: check in from any location (GPS is still recorded for the log only).
    </p>
    <p v-else class="text-sm text-gray-400">
      GPS is used to verify you are within the office radius. On office WiFi, your network IP can also satisfy the check (configured by admin).
    </p>

    <div v-if="geoStatus" class="text-sm" :class="geoError ? 'text-amber-400' : 'text-gray-300'">
      {{ geoStatus }}
    </div>

    <div class="flex flex-wrap gap-3">
      <button
        type="button"
        :disabled="store.actionLoading || !store.canCheckIn"
        class="px-4 py-2 rounded-lg bg-gradient-to-r from-emerald-600 to-teal-600 text-white font-medium disabled:opacity-40 disabled:cursor-not-allowed"
        @click="onCheckIn"
      >
        {{ store.actionLoading && store.canCheckIn ? 'Checking in…' : 'Check in' }}
      </button>
      <button
        type="button"
        :disabled="store.actionLoading || !store.canCheckOut"
        class="px-4 py-2 rounded-lg bg-gradient-to-r from-indigo-600 to-violet-600 text-white font-medium disabled:opacity-40 disabled:cursor-not-allowed"
        @click="onCheckOut"
      >
        {{ store.actionLoading && store.canCheckOut ? 'Checking out…' : 'Check out' }}
      </button>
    </div>

    <p v-if="store.error" class="text-sm text-red-400">{{ store.error }}</p>
  </div>
</template>

<script setup lang="ts">
import { useAttendanceStore } from '../store/attendance-store'

const store = useAttendanceStore()

function isoWeekdayToday(): number {
  const w = new Date().getDay()
  return w === 0 ? 7 : w
}

const isWfhToday = computed(() => {
  const c = store.officeConfig
  const list = c?.wfh_days
  if (!list?.length) return false
  return list.includes(isoWeekdayToday())
})
const geoStatus = ref<string | null>(null)
const geoError = ref(false)

function readPosition(): Promise<GeolocationPosition> {
  return new Promise((resolve, reject) => {
    if (!import.meta.client || !navigator.geolocation) {
      reject(new Error('Geolocation is not available in this browser.'))
      return
    }
    navigator.geolocation.getCurrentPosition(resolve, reject, {
      enableHighAccuracy: true,
      timeout: 20000,
      maximumAge: 0,
    })
  })
}

async function onCheckIn() {
  geoError.value = false
  geoStatus.value = 'Requesting GPS position…'
  try {
    const pos = await readPosition()
    const lat = pos.coords.latitude
    const lng = pos.coords.longitude
    const acc = pos.coords.accuracy
    geoStatus.value =
      typeof acc === 'number'
        ? `Position acquired (accuracy ~${Math.round(acc)} m). Submitting check-in…`
        : 'Position acquired. Submitting check-in…'
    const ok = await store.checkIn(lat, lng)
    if (ok) {
      geoStatus.value = 'Check-in successful.'
    }
  } catch (e: any) {
    geoError.value = true
    geoStatus.value = e?.message || 'Could not read GPS. Allow location access and try again.'
  }
}

async function onCheckOut() {
  geoError.value = false
  geoStatus.value = null
  await store.checkOut()
}
</script>

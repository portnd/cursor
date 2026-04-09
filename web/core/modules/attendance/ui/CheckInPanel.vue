<template>
  <div class="rounded-xl border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800/80 p-6 space-y-4">
    <h2 class="text-lg font-semibold text-gray-900 dark:text-white">Check-in / check-out</h2>
    <p v-if="isWfhToday" class="text-sm text-sky-300/90">
      Today is a <strong>WFH</strong> day: check in from any location (GPS is still recorded for the log only).
    </p>
    <p v-else class="text-sm text-gray-400">
      GPS is used to verify you are within the office radius for check-in.
    </p>

    <div v-if="geoStatus" class="text-sm" :class="geoError ? 'text-amber-400' : 'text-gray-300'">
      {{ geoStatus }}
    </div>

    <div v-if="officeMap" class="rounded-lg border border-gray-700 bg-gray-900/70 p-3 space-y-2">
      <div class="flex flex-wrap items-center justify-between gap-2 text-xs">
        <p class="text-gray-300">
          Office location:
          <span class="text-gray-100 font-medium">{{ officeMap.name }}</span>
          <span class="text-gray-400">• Radius {{ officeMap.radiusMeters }}m</span>
        </p>
        <a
          :href="officeMap.link"
          target="_blank"
          rel="noopener noreferrer"
          class="text-violet-300 hover:text-violet-200"
        >
          Open office map
        </a>
      </div>

      <div class="relative h-64 w-full overflow-hidden rounded-md border border-gray-700">
        <iframe
          :src="officeMap.embedSrc"
          title="Office location and radius"
          class="h-full w-full pointer-events-none"
          loading="lazy"
          referrerpolicy="no-referrer-when-downgrade"
          tabindex="-1"
        />
        <div class="pointer-events-none absolute inset-0 flex items-center justify-center">
          <div
            class="rounded-full border-2 border-rose-300/90 bg-rose-300/15 shadow-[0_0_0_9999px_rgba(0,0,0,0.05)]"
            :style="{ width: `${officeMap.circleSizePx}px`, height: `${officeMap.circleSizePx}px` }"
          />
        </div>
        <div class="pointer-events-none absolute inset-0 flex items-center justify-center">
          <div class="h-9 w-9 rounded-full border border-slate-600 bg-slate-900/90 text-sky-300 shadow-lg flex items-center justify-center" aria-hidden="true">
            <svg viewBox="0 0 24 24" class="h-5 w-5 fill-current">
              <path d="M3 21h18v-2h-1V8.5a1.5 1.5 0 0 0-.88-1.36l-6-2.73a1.5 1.5 0 0 0-1.24 0l-6 2.73A1.5 1.5 0 0 0 5 8.5V19H3v2Zm4-2V9.14l5-2.27 5 2.27V19h-2v-3.5a1.5 1.5 0 0 0-1.5-1.5h-3A1.5 1.5 0 0 0 9 15.5V19H7Zm4 0v-3h2v3h-2Zm-2-6h2v-2H9v2Zm0-3h2V8H9v2Zm4 3h2v-2h-2v2Zm0-3h2V8h-2v2Z" />
            </svg>
          </div>
        </div>
        <div v-if="userMarkerStyle" class="pointer-events-none absolute inset-0">
          <div
            class="absolute -translate-x-1/2 -translate-y-1/2 h-3 w-3 rounded-full bg-emerald-400 ring-2 ring-emerald-200/90 shadow-[0_0_0_6px_rgba(16,185,129,0.15)]"
            :style="userMarkerStyle"
            aria-hidden="true"
          />
        </div>
      </div>
      <p class="text-[11px] text-gray-400">
        Circle overlay represents office check-in radius centered at office marker.
      </p>
      <p v-if="distanceInfo" class="text-xs" :class="distanceInfo.isOutside ? 'text-amber-300' : 'text-emerald-300'">
        <template v-if="distanceInfo.isOutside">
          Outside radius by {{ distanceInfo.outsideMeters }} m
          (distance {{ distanceInfo.distanceMeters }} m, radius {{ officeMap.radiusMeters }} m)
        </template>
        <template v-else>
          Inside radius (distance {{ distanceInfo.distanceMeters }} m / radius {{ officeMap.radiusMeters }} m)
        </template>
      </p>
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

    <div class="rounded-lg border border-amber-300 dark:border-amber-800/60 bg-amber-50 dark:bg-amber-950/20 p-3 space-y-2">
      <p class="text-sm text-amber-800 dark:text-amber-200">อยู่นอกสถานที่?</p>
      <p class="text-xs text-gray-700 dark:text-gray-400">
        ส่งคำขอ check-in นอกสถานที่พร้อมเหตุผล แล้วรอ CEO อนุมัติ เมื่ออนุมัติแล้วจะบันทึก check-in ให้อัตโนมัติ
      </p>
      <textarea
        v-model.trim="offsiteReason"
        rows="2"
        maxlength="1000"
        placeholder="เหตุผลที่ต้อง check-in นอกสถานที่"
        class="w-full rounded-md border border-gray-600 bg-gray-900 px-3 py-2 text-sm text-gray-100"
      />
      <button
        type="button"
        :disabled="store.actionLoading || !store.canCheckIn || offsiteReason.length < 5 || pendingOffsiteRequest"
        class="px-4 py-2 rounded-lg bg-gradient-to-r from-amber-600 to-orange-600 text-white font-medium disabled:opacity-40 disabled:cursor-not-allowed"
        @click="onRequestOffsite"
      >
        {{ store.actionLoading ? 'Submitting…' : 'ขอ check-in นอกสถานที่' }}
      </button>
      <p v-if="pendingOffsiteRequest" class="text-xs text-sky-300">
        มีคำขอค้างอนุมัติอยู่แล้ว ({{ pendingOffsiteRequest.status }})
      </p>
    </div>

    <div v-if="store.todayRecord?.check_in_at" class="rounded-lg border border-amber-800/60 bg-amber-950/20 p-3 space-y-2">
      <p class="text-sm text-amber-200">ต้อง check-out นอกสถานที่?</p>
      <p class="text-xs text-gray-400">
        ส่งคำขอ check-out นอกสถานที่พร้อมเหตุผล แล้วรอ CEO อนุมัติ เมื่ออนุมัติแล้วจะบันทึก check-out ให้อัตโนมัติ
      </p>
      <textarea
        v-model.trim="offsiteCheckoutReason"
        rows="2"
        maxlength="1000"
        placeholder="เหตุผลที่ต้อง check-out นอกสถานที่"
        class="w-full rounded-md border border-gray-600 bg-gray-900 px-3 py-2 text-sm text-gray-100"
      />
      <button
        type="button"
        :disabled="store.actionLoading || !store.canCheckOut || offsiteCheckoutReason.length < 5 || pendingOffsiteCheckoutRequest"
        class="px-4 py-2 rounded-lg bg-gradient-to-r from-amber-600 to-orange-600 text-white font-medium disabled:opacity-40 disabled:cursor-not-allowed"
        @click="onRequestOffsiteCheckout"
      >
        {{ store.actionLoading ? 'Submitting…' : 'ขอ check-out นอกสถานที่' }}
      </button>
      <p v-if="pendingOffsiteCheckoutRequest" class="text-xs text-sky-300">
        มีคำขอ check-out ค้างอนุมัติอยู่แล้ว ({{ pendingOffsiteCheckoutRequest.status }})
      </p>
    </div>

    <p v-if="store.error" class="text-sm text-red-400">{{ store.error }}</p>
  </div>
</template>

<script setup lang="ts">
import { useAttendanceStore } from '../store/attendance-store'

const store = useAttendanceStore()

type LatLng = {
  lat: number
  lng: number
}

function isoWeekdayToday(): number {
  const w = new Date().getDay()
  return w === 0 ? 7 : w
}

function metersToLatDelta(meters: number): number {
  return meters / 111_320
}

function metersToLngDelta(meters: number, lat: number): number {
  const latRad = (lat * Math.PI) / 180
  const metersPerDegree = 111_320 * Math.cos(latRad)
  if (!Number.isFinite(metersPerDegree) || Math.abs(metersPerDegree) < 1e-9) {
    return meters / 111_320
  }
  return meters / metersPerDegree
}

function haversineDistanceMeters(aLat: number, aLng: number, bLat: number, bLng: number): number {
  const toRad = (deg: number) => (deg * Math.PI) / 180
  const earthRadius = 6_371_000
  const dLat = toRad(bLat - aLat)
  const dLng = toRad(bLng - aLng)
  const aa =
    Math.sin(dLat / 2) * Math.sin(dLat / 2) +
    Math.cos(toRad(aLat)) * Math.cos(toRad(bLat)) * Math.sin(dLng / 2) * Math.sin(dLng / 2)
  const c = 2 * Math.atan2(Math.sqrt(aa), Math.sqrt(1 - aa))
  return earthRadius * c
}

const isWfhToday = computed(() => {
  const c = store.officeConfig
  const list = c?.wfh_days
  if (!list?.length) return false
  return list.includes(isoWeekdayToday())
})
const geoStatus = ref<string | null>(null)
const geoError = ref(false)
const lastPosition = ref<LatLng | null>(null)
const offsiteReason = ref('')
const offsiteCheckoutReason = ref('')

const pendingOffsiteRequest = computed(() => {
  const req = store.todayOffsiteRequest
  if (!req) return null
  if (req.status !== 'PENDING') return null
  return req
})

const pendingOffsiteCheckoutRequest = computed(() => {
  const req = store.todayOffsiteCheckOutRequest
  if (!req) return null
  if (req.status !== 'PENDING') return null
  return req
})

const officeMap = computed(() => {
  const cfg = store.officeConfig
  if (!cfg) return null

  const lat = cfg.latitude
  const lng = cfg.longitude
  const radiusMeters = Math.max(1, cfg.radius_meters || 1)

  const mapPaddingMeters = Math.max(120, radiusMeters * 2.2)
  const latDelta = metersToLatDelta(mapPaddingMeters)
  const lngDelta = metersToLngDelta(mapPaddingMeters, lat)

  const left = lng - lngDelta
  const right = lng + lngDelta
  const top = lat + latDelta
  const bottom = lat - latDelta

  const circleSizePx = Math.max(40, Math.min(180, (radiusMeters / mapPaddingMeters) * 256))

  return {
    name: cfg.name,
    lat,
    lng,
    left,
    right,
    top,
    bottom,
    radiusMeters,
    embedSrc: `https://www.openstreetmap.org/export/embed.html?bbox=${left}%2C${bottom}%2C${right}%2C${top}&layer=mapnik`,
    link: `https://www.openstreetmap.org/?mlat=${lat}&mlon=${lng}#map=17/${lat}/${lng}`,
    circleSizePx,
  }
})

const userMarkerStyle = computed(() => {
  if (!officeMap.value || !lastPosition.value) return null

  const { left, right, top, bottom } = officeMap.value
  const { lat, lng } = lastPosition.value

  const width = right - left
  const height = top - bottom
  if (Math.abs(width) < 1e-12 || Math.abs(height) < 1e-12) return null

  const x = ((lng - left) / width) * 100
  const y = ((top - lat) / height) * 100

  const clampedX = Math.max(0, Math.min(100, x))
  const clampedY = Math.max(0, Math.min(100, y))

  return {
    left: `${clampedX}%`,
    top: `${clampedY}%`,
  }
})

const distanceInfo = computed(() => {
  if (!officeMap.value || !lastPosition.value) return null

  const distMeters = haversineDistanceMeters(
    lastPosition.value.lat,
    lastPosition.value.lng,
    officeMap.value.lat,
    officeMap.value.lng,
  )
  const radiusMeters = officeMap.value.radiusMeters
  const isOutside = distMeters > radiusMeters

  return {
    distanceMeters: Math.round(distMeters),
    outsideMeters: Math.max(0, Math.round(distMeters - radiusMeters)),
    isOutside,
  }
})

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

    lastPosition.value = { lat, lng }

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

async function onRequestOffsite() {
  geoError.value = false
  geoStatus.value = 'Requesting GPS position for offsite request…'
  try {
    const pos = await readPosition()
    const lat = pos.coords.latitude
    const lng = pos.coords.longitude
    lastPosition.value = { lat, lng }
    geoStatus.value = 'Submitting offsite check-in request…'
    const ok = await store.requestOffsiteCheckIn(lat, lng, offsiteReason.value)
    if (ok) {
      geoStatus.value = 'ส่งคำขอสำเร็จ รอ CEO อนุมัติ'
      offsiteReason.value = ''
    }
  } catch (e: any) {
    geoError.value = true
    geoStatus.value = e?.message || 'Could not read GPS. Allow location access and try again.'
  }
}

async function onRequestOffsiteCheckout() {
  geoError.value = false
  geoStatus.value = 'Requesting GPS position for offsite check-out request…'
  try {
    const pos = await readPosition()
    const lat = pos.coords.latitude
    const lng = pos.coords.longitude
    lastPosition.value = { lat, lng }
    geoStatus.value = 'Submitting offsite check-out request…'
    const ok = await store.requestOffsiteCheckOut(lat, lng, offsiteCheckoutReason.value)
    if (ok) {
      geoStatus.value = 'ส่งคำขอ check-out สำเร็จ รอ CEO อนุมัติ'
      offsiteCheckoutReason.value = ''
    }
  } catch (e: any) {
    geoError.value = true
    geoStatus.value = e?.message || 'Could not read GPS. Allow location access and try again.'
  }
}
</script>

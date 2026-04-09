<template>
  <div
    :class="[
      'rounded-xl border transition-all',
      cardBorderClass,
      isUrgent ? 'shadow-lg shadow-yellow-900/20' : ''
    ]"
    class="bg-gray-800/60 p-5"
  >
    <!-- Top row: title + status badge + env badge -->
    <div class="flex items-start justify-between gap-3 mb-3">
      <div class="min-w-0 flex-1">
        <h3 class="font-semibold text-white text-sm truncate leading-tight">{{ request.title }}</h3>
        <div class="flex items-center gap-2 mt-1 flex-wrap">
          <!-- Branch -->
          <code class="text-xs px-2 py-0.5 rounded bg-gray-700/80 text-cyan-300 font-mono border border-gray-600/50">
            <span class="text-gray-500">⎇</span> {{ request.branch }}
          </code>
          <!-- Task ref -->
          <span v-if="request.task_ref" class="text-xs text-gray-500">
            {{ request.task_ref }}
          </span>
        </div>
      </div>
      <div class="flex items-center gap-2 shrink-0">
        <span :class="envBadge" class="text-[10px] font-bold px-2 py-1 rounded-md border uppercase">
          {{ request.environment }}
        </span>
        <span :class="statusBadge" class="text-[10px] font-bold px-2 py-1 rounded-md border uppercase">
          {{ statusLabel }}
        </span>
      </div>
    </div>

    <!-- Description -->
    <p v-if="request.description" class="text-xs text-gray-400 mb-3 line-clamp-2 leading-relaxed">
      {{ request.description }}
    </p>

    <!-- PR URL -->
    <a
      v-if="request.pr_url"
      :href="request.pr_url"
      target="_blank"
      rel="noopener"
      class="inline-flex items-center gap-1.5 text-xs text-blue-400 hover:text-blue-300 hover:underline mb-3 transition-colors"
    >
      <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
      </svg>
      View PR / Merge Request
    </a>

    <!-- Rejection reason / review notes -->
    <div v-if="request.rejection_reason" class="mb-3 px-3 py-2 rounded-lg bg-red-500/10 border border-red-500/20">
      <p class="text-[10px] font-bold text-red-400 uppercase tracking-wide mb-0.5">Rejection Reason</p>
      <p class="text-xs text-red-300">{{ request.rejection_reason }}</p>
    </div>
    <div v-if="request.review_notes && !request.rejection_reason" class="mb-3 px-3 py-2 rounded-lg bg-gray-700/40 border border-gray-600/30">
      <p class="text-[10px] font-bold text-gray-400 uppercase tracking-wide mb-0.5">Review Notes</p>
      <p class="text-xs text-gray-300">{{ request.review_notes }}</p>
    </div>

    <!-- Footer: meta + actions -->
    <div class="flex items-center justify-between gap-3 pt-3 border-t border-gray-700/50">
      <!-- Meta -->
      <div class="flex items-center gap-3 text-[11px] text-gray-500 min-w-0">
        <span class="flex items-center gap-1 truncate">
          <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
          </svg>
          <span class="truncate">{{ requesterName }}</span>
        </span>
        <span class="shrink-0">{{ timeAgo }}</span>
        <span v-if="request.deployed_at" class="flex items-center gap-1 text-green-500 shrink-0">
          🚀 {{ deployedAt }}
        </span>
      </div>

      <!-- Chief Engineer action buttons -->
      <div v-if="isChiefEngineer" class="flex items-center gap-2 shrink-0">
        <!-- PENDING: pick up -->
        <button
          v-if="request.status === 'PENDING'"
          @click="$emit('pick', request)"
          class="flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-blue-100 dark:bg-blue-600/20 border border-blue-300 dark:border-blue-500/30 text-blue-400 hover:bg-blue-100 dark:bg-blue-600/30 text-xs font-semibold transition-colors"
        >
          <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
          </svg>
          Pick Up
        </button>

        <!-- PENDING or REVIEWING: approve / reject -->
        <template v-if="request.status === 'PENDING' || request.status === 'REVIEWING'">
          <button
            @click="$emit('reject', request)"
            class="flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-red-100 dark:bg-red-600/20 border border-red-300 dark:border-red-500/30 text-red-400 hover:bg-red-100 dark:bg-red-600/30 text-xs font-semibold transition-colors"
          >
            ✗ Reject
          </button>
          <button
            @click="$emit('approve', request)"
            class="flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-green-100 dark:bg-green-600/20 border border-green-300 dark:border-green-500/30 text-green-400 hover:bg-green-100 dark:bg-green-600/30 text-xs font-semibold transition-colors"
          >
            ✓ Approve
          </button>
        </template>

        <!-- APPROVED: mark deployed -->
        <button
          v-if="request.status === 'APPROVED'"
          @click="$emit('deploy', request)"
          class="flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-gradient-to-r from-cyan-100 dark:from-cyan-600/30 to-blue-100 dark:to-blue-600/30 border border-cyan-300 dark:border-cyan-500/40 text-cyan-300 hover:from-cyan-100 dark:from-cyan-600/50 hover:to-blue-100 dark:to-blue-600/50 text-xs font-semibold transition-all"
        >
          🚀 Mark Deployed
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { DeploymentRequest } from '~/core/modules/deployment/infrastructure/deployment-api'

const props = defineProps<{
  request: DeploymentRequest
  isChiefEngineer: boolean
}>()

defineEmits<{
  pick: [req: DeploymentRequest]
  approve: [req: DeploymentRequest]
  reject: [req: DeploymentRequest]
  deploy: [req: DeploymentRequest]
}>()

const isUrgent = computed(() =>
  props.request.status === 'PENDING' || props.request.status === 'REVIEWING'
)

const cardBorderClass = computed(() => {
  switch (props.request.status) {
    case 'PENDING':   return 'border-yellow-700/40'
    case 'REVIEWING': return 'border-blue-700/40'
    case 'APPROVED':  return 'border-green-700/40'
    case 'REJECTED':  return 'border-red-700/40'
    case 'DEPLOYED':  return 'border-cyan-700/40'
    default:          return 'border-gray-700/40'
  }
})

const statusLabel = computed(() => {
  const map: Record<string, string> = {
    PENDING: 'Pending',
    REVIEWING: 'In Review',
    APPROVED: 'Approved',
    REJECTED: 'Rejected',
    DEPLOYED: 'Deployed',
  }
  return map[props.request.status] ?? props.request.status
})

const statusBadge = computed(() => {
  switch (props.request.status) {
    case 'PENDING':   return 'bg-yellow-500/10 text-yellow-400 border-yellow-500/30'
    case 'REVIEWING': return 'bg-blue-500/10 text-blue-400 border-blue-500/30'
    case 'APPROVED':  return 'bg-green-500/10 text-green-400 border-green-500/30'
    case 'REJECTED':  return 'bg-red-500/10 text-red-400 border-red-500/30'
    case 'DEPLOYED':  return 'bg-cyan-500/10 text-cyan-400 border-cyan-500/30'
    default:          return 'bg-gray-500/10 text-gray-400 border-gray-500/30'
  }
})

const envBadge = computed(() => {
  switch (props.request.environment) {
    case 'PRODUCTION': return 'bg-orange-500/10 text-orange-400 border-orange-500/30'
    case 'PRE-PROD':   return 'bg-amber-500/10 text-amber-400 border-amber-500/30'
    default:           return 'bg-violet-500/10 text-violet-400 border-violet-500/30'
  }
})

const requesterName = computed(() =>
  props.request.requester_display_name || props.request.requester_email || `User #${props.request.requester_id}`
)

function fmtAgo(dateStr: string): string {
  const diff = Date.now() - new Date(dateStr).getTime()
  const mins = Math.floor(diff / 60000)
  if (mins < 1)   return 'just now'
  if (mins < 60)  return `${mins}m ago`
  const hrs = Math.floor(mins / 60)
  if (hrs < 24)   return `${hrs}h ago`
  return `${Math.floor(hrs / 24)}d ago`
}

const timeAgo = computed(() => fmtAgo(props.request.created_at))
const deployedAt = computed(() =>
  props.request.deployed_at ? 'Deployed ' + fmtAgo(props.request.deployed_at) : ''
)
</script>

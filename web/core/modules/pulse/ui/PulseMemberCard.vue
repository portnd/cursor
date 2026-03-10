<template>
  <div
    :class="[
      'relative flex flex-col gap-3 rounded-xl border bg-gray-800/70 p-4 transition hover:border-gray-600',
      member.has_blocker ? 'border-red-700' : 'border-gray-700',
    ]"
  >
    <!-- Blocker banner -->
    <div
      v-if="member.has_blocker"
      class="absolute inset-x-0 top-0 rounded-t-xl bg-red-900/60 px-3 py-1 text-xs font-semibold text-red-300 flex items-center gap-1"
    >
      <span>🚧</span> Blocker reported
    </div>

    <!-- Header row -->
    <div :class="['flex items-center gap-3', member.has_blocker ? 'mt-5' : '']">
      <div class="flex h-9 w-9 flex-shrink-0 items-center justify-center rounded-full bg-indigo-700 text-sm font-bold text-white">
        {{ initials }}
      </div>
      <div class="min-w-0 flex-1">
        <p class="truncate text-sm font-semibold text-white">{{ displayName }}</p>
        <p class="truncate text-xs text-gray-400">{{ member.user_email }}</p>
      </div>
      <div class="flex-shrink-0 rounded-md bg-gray-700 px-2 py-0.5 text-xs text-gray-300">
        {{ member.total_logged_hours }}h
      </div>
    </div>

    <!-- No standup -->
    <div
      v-if="!member.standup"
      class="rounded-lg border border-gray-700 bg-gray-900/50 px-3 py-2 text-xs text-gray-500 italic"
    >
      No standup submitted yet.
    </div>

    <!-- Standup content -->
    <template v-else>
      <!-- Yesterday -->
      <div class="space-y-0.5">
        <p class="text-xs font-semibold uppercase tracking-wide text-gray-500">Yesterday</p>
        <p class="text-sm text-gray-200 leading-relaxed line-clamp-3">{{ member.standup.yesterday_summary }}</p>
      </div>

      <!-- Today task IDs -->
      <div v-if="member.standup.today_task_ids && member.standup.today_task_ids.length" class="space-y-1">
        <p class="text-xs font-semibold uppercase tracking-wide text-gray-500">Today's Tasks</p>
        <div class="flex flex-wrap gap-1">
          <span
            v-for="tid in member.standup.today_task_ids"
            :key="tid"
            class="rounded bg-indigo-900/60 px-1.5 py-0.5 text-xs font-mono text-indigo-300"
          >{{ tid }}</span>
        </div>
      </div>

      <!-- Blocker -->
      <div v-if="member.standup.blocker" class="rounded-lg border border-red-800 bg-red-50/5 px-3 py-2">
        <p class="text-xs font-semibold uppercase tracking-wide text-red-500">Blocker</p>
        <p class="mt-0.5 text-sm text-red-300 leading-relaxed">{{ member.standup.blocker }}</p>
      </div>
    </template>

    <!-- Activity feed -->
    <div v-if="member.latest_activities && member.latest_activities.length" class="border-t border-gray-700/60 pt-2 space-y-1">
      <p class="text-xs font-semibold uppercase tracking-wide text-gray-500">Activity</p>
      <ul class="space-y-1">
        <li
          v-for="(act, idx) in member.latest_activities"
          :key="idx"
          class="flex items-start gap-2 text-xs text-gray-400"
        >
          <span class="mt-0.5 flex-shrink-0 text-base leading-none">{{ act.type === 'time_log' ? '⏱' : '📦' }}</span>
          <span class="flex-1 leading-relaxed">
            {{ act.description }}
            <span v-if="act.type === 'time_log'" class="text-gray-500"> — {{ act.minutes }}m</span>
            <span
              v-if="act.type === 'submission' && act.ai_verdict"
              :class="verdictColor(act.ai_verdict)"
              class="ml-1 font-semibold"
            >{{ act.ai_verdict }} ({{ act.ai_score }})</span>
          </span>
          <span class="flex-shrink-0 text-gray-600">{{ formatTime(act.occurred_at) }}</span>
        </li>
      </ul>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { UserPulse } from '../infrastructure/pulse-api'

const props = defineProps<{ member: UserPulse }>()

const displayName = computed(
  () => props.member.user_display_name || props.member.user_email,
)

const initials = computed(() =>
  displayName.value
    .split(/[\s@]/)
    .slice(0, 2)
    .map((w: string) => w[0]?.toUpperCase() ?? '')
    .join(''),
)

function formatTime(iso: string) {
  return new Date(iso).toLocaleTimeString('en-GB', { hour: '2-digit', minute: '2-digit' })
}

function verdictColor(verdict: string) {
  if (verdict === 'PASS') return 'text-green-400'
  if (verdict === 'FAIL') return 'text-red-400'
  return 'text-yellow-400'
}
</script>

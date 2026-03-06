<template>
  <div class="rounded-lg border border-gray-700 bg-gray-800 overflow-hidden">
    <div class="px-4 py-3 border-b border-gray-700">
      <h3 class="text-lg font-bold text-white">Team Leaderboard</h3>
      <p class="text-xs text-gray-400">Ranked by composite KPI score</p>
    </div>
    <div class="overflow-x-auto">
      <table class="w-full text-sm">
        <thead>
          <tr class="border-b border-gray-700 text-left text-gray-400">
            <th class="px-4 py-2 font-medium">#</th>
            <th class="px-4 py-2 font-medium">Member</th>
            <th class="px-4 py-2 font-medium">Role</th>
            <th class="px-4 py-2 font-medium text-right">Delivery</th>
            <th class="px-4 py-2 font-medium text-right">Quality</th>
            <th class="px-4 py-2 font-medium text-right">Rework</th>
            <th class="px-4 py-2 font-medium text-right">Score</th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="(m, i) in members"
            :key="m.user_id"
            class="border-b border-gray-700/50 hover:bg-gray-700/30 transition-colors"
          >
            <td class="px-4 py-2 text-gray-500">{{ i + 1 }}</td>
            <td class="px-4 py-2 text-white font-medium">{{ m.email }}</td>
            <td class="px-4 py-2">
              <span
                class="rounded px-2 py-0.5 text-xs font-medium"
                :class="roleBadgeClass(m.role)"
              >
                {{ m.role }}
              </span>
            </td>
            <td class="px-4 py-2 text-right">
              <span :class="pctColor(m.delivery_rate_pct)">{{ m.delivery_rate_pct.toFixed(1) }}%</span>
            </td>
            <td class="px-4 py-2 text-right text-gray-300">{{ m.code_quality_index.toFixed(0) }}</td>
            <td class="px-4 py-2 text-right">
              <span :class="reworkColor(m.rework_rate_pct)">{{ m.rework_rate_pct.toFixed(1) }}%</span>
            </td>
            <td class="px-4 py-2 text-right font-bold" :class="scoreColor(m.composite_score)">
              {{ m.composite_score.toFixed(1) }}
            </td>
          </tr>
        </tbody>
      </table>
    </div>
    <div v-if="members.length === 0" class="px-4 py-8 text-center text-gray-500">
      No team members to display.
    </div>
  </div>
</template>

<script setup lang="ts">
import type { TeamMemberKPI } from '~/core/modules/performance/performance-api'

defineProps<{
  members: TeamMemberKPI[]
}>()

function roleBadgeClass(role: string) {
  if (role === 'CEO') return 'bg-purple-600/30 text-purple-300'
  if (role === 'PM') return 'bg-blue-600/30 text-blue-300'
  return 'bg-gray-600/30 text-gray-300'
}

function pctColor(pct: number) {
  if (pct >= 85) return 'text-emerald-400'
  if (pct >= 70) return 'text-amber-400'
  return 'text-red-400'
}

function reworkColor(pct: number) {
  if (pct <= 15) return 'text-emerald-400'
  if (pct <= 25) return 'text-amber-400'
  return 'text-red-400'
}

function scoreColor(score: number) {
  if (score >= 80) return 'text-emerald-400'
  if (score >= 60) return 'text-amber-400'
  return 'text-gray-400'
}
</script>

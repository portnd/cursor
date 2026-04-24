<template>
  <Teleport to="body">
    <div
      v-if="modelValue && member"
      class="fixed inset-0 z-[100] flex items-center justify-center p-4 bg-black/70"
      role="dialog"
      aria-modal="true"
      aria-labelledby="score-breakdown-title"
      @click.self="close"
    >
      <div
        class="max-h-[90vh] w-full max-w-lg overflow-y-auto rounded-xl border border-gray-600 bg-gray-900 shadow-2xl"
        @keydown.escape.prevent="close"
      >
        <div class="sticky top-0 flex items-start justify-between gap-3 border-b border-gray-700 bg-gray-900/95 px-4 py-3 backdrop-blur-sm">
          <div>
            <h2 id="score-breakdown-title" class="text-lg font-bold text-white">
              How this score is calculated
            </h2>
            <p class="text-xs text-gray-500 mt-0.5">
              วิธีคำนวณคะแนน — {{ member.email }}
            </p>
          </div>
          <button
            type="button"
            class="rounded-lg p-1.5 text-gray-400 hover:bg-gray-800 hover:text-gray-900 dark:text-white transition-colors"
            aria-label="Close"
            @click="close"
          >
            <svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
            </svg>
          </button>
        </div>

        <div class="px-4 py-4 space-y-4 text-sm text-gray-300">
          <!-- Product Owner (CEO leaderboard): team average -->
          <template v-if="member.role === 'PRODUCT_OWNER' || member.role === 'PM'">
            <p class="text-gray-300 leading-relaxed">
              This <strong class="text-white">Score</strong> is the <strong class="text-white">average composite KPI</strong>
              of every engineer who has at least one task assigned by this Product Owner. Each engineer’s composite uses the same
              weighted formula as in the engineer row (delivery, quality, rework, velocity, time accuracy), but only for tasks
              tied to that Product Owner.
            </p>
            <p class="text-xs text-gray-500 leading-relaxed">
              คะแนนของ Product Owner = ค่าเฉลี่ยคะแนนรวม (composite) ของ engineer ที่ Product Owner คนนี้เป็นคน assign task ให้อย่างน้อย 1 งาน
            </p>
            <p class="rounded-lg border border-gray-700/60 bg-gray-800/30 px-3 py-2 text-xs text-gray-400">
              If no engineers qualify, the API uses this Product Owner’s <strong class="text-gray-300">profile health score</strong> instead.
              <span class="block text-gray-600 mt-1">ถ้าไม่มี dev ที่เข้าเงื่อนไข ระบบใช้ health score จากโปรไฟล์</span>
            </p>
            <div class="rounded-lg border border-gray-700 bg-gray-800/50 px-3 py-2">
              <div class="text-xs text-gray-500 uppercase tracking-wide">Displayed score</div>
              <div class="text-2xl font-bold text-white tabular-nums">{{ member.composite_score.toFixed(1) }}</div>
            </div>
          </template>

          <!-- CEO / other non-engineer: health only -->
          <template v-else-if="!isEngineerLikeRole(member.role)">
            <p class="text-gray-300 leading-relaxed">
              For this role, <strong class="text-white">Score</strong> is the user’s
              <strong class="text-white">health score</strong> from their profile — not the engineer KPI blend.
            </p>
            <p class="text-xs text-gray-500 leading-relaxed">
              บทบาทนี้ใช้ค่า health score จากโปรไฟล์ ไม่ใช้สูตร composite ของ engineer
            </p>
            <div class="rounded-lg border border-gray-700 bg-gray-800/50 px-3 py-2">
              <div class="text-xs text-gray-500 uppercase tracking-wide">Score</div>
              <div class="text-2xl font-bold text-white tabular-nums">{{ member.composite_score.toFixed(1) }}</div>
            </div>
          </template>

          <!-- Engineer: full weighted breakdown -->
          <template v-else>
            <p
              v-if="focus === 'delivery'"
              :class="sectionHighlightClass('delivery')"
              class="rounded-lg px-2 py-1 -mx-2"
            >
              <strong class="text-white">Delivery</strong> — % of tasks that had a due date and reached
              <code class="text-cyan-400">job done</code> on or before that due date.
              <span class="block text-xs text-gray-500 mt-1">คิดจากเวลาที่งานถึงสถานะ job done ไม่ใช่รอ completed สุดท้าย</span>
            </p>
            <p
              v-if="focus === 'quality'"
              :class="sectionHighlightClass('quality')"
              class="rounded-lg px-2 py-1 -mx-2"
            >
              <strong class="text-white">Quality</strong> — approval rate from tasks with submissions:
              (1 − tasks with a <code class="text-cyan-400">[REJECTED]</code> comment ÷ tasks with submissions) × 100.
              <span class="block text-xs text-gray-500 mt-1">อัตรางานที่ไม่ถูก reject จากงานที่มี submission</span>
            </p>
            <p
              v-if="focus === 'rework'"
              :class="sectionHighlightClass('rework')"
              class="rounded-lg px-2 py-1 -mx-2"
            >
              <strong class="text-white">Rework</strong> — uses the same rework event as the
              <code class="text-cyan-400">discipline</code> page: every <code class="text-cyan-400">[REJECTED]</code> comment counts as 1 rework.
              The rate is <code class="text-cyan-400">rework events ÷ (job done + rework events)</code> × 100.
              <span class="block text-xs text-gray-500 mt-1">นับ rework แบบเดียวกับหน้า discipline คือทุกคอมเมนต์ [REJECTED] = 1 ครั้ง</span>
            </p>
            <p
              v-if="focus === 'composite'"
              class="text-gray-400 text-xs leading-relaxed"
            >
              Composite blends five signals (same weights as the API). Tap Delivery, Quality, or Rework in the table for those definitions.
              <span class="block text-gray-600 mt-1">คะแนนรวม = ถ่วงน้ำหนัก 5 องค์ประกอบตามด้านล่าง</span>
            </p>

            <div class="rounded-lg border border-gray-700 overflow-hidden">
              <table class="w-full text-xs">
                <thead>
                  <tr class="border-b border-gray-700 bg-gray-800/80 text-left text-gray-400">
                    <th class="px-3 py-2 font-medium">Input</th>
                    <th class="px-3 py-2 font-medium text-right">Weight</th>
                    <th class="px-3 py-2 font-medium text-right">Contribution</th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-gray-700/60">
                  <tr :class="rowHighlight('delivery')">
                    <td class="px-3 py-2">
                      <div class="text-white font-medium">Delivery rate</div>
                      <div class="text-gray-500 tabular-nums">{{ member.delivery_rate_pct.toFixed(1) }}%</div>
                    </td>
                    <td class="px-3 py-2 text-right text-gray-400">30%</td>
                    <td class="px-3 py-2 text-right font-mono text-emerald-300/90 tabular-nums">
                      {{ b.deliveryTerm.toFixed(2) }}
                    </td>
                  </tr>
                  <tr :class="rowHighlight('quality')">
                    <td class="px-3 py-2">
                      <div class="text-white font-medium">Quality (capped 100)</div>
                      <div class="text-gray-500 tabular-nums">min(index, 100) = {{ b.qualityNorm.toFixed(1) }}</div>
                    </td>
                    <td class="px-3 py-2 text-right text-gray-400">25%</td>
                    <td class="px-3 py-2 text-right font-mono text-emerald-300/90 tabular-nums">
                      {{ b.qualityTerm.toFixed(2) }}
                    </td>
                  </tr>
                  <tr :class="rowHighlight('rework')">
                    <td class="px-3 py-2">
                      <div class="text-white font-medium">Rework cushion</div>
                      <div class="text-gray-500 tabular-nums">max(0, 100 − {{ member.rework_rate_pct.toFixed(1) }}%) = {{ b.reworkNorm.toFixed(1) }}</div>
                    </td>
                    <td class="px-3 py-2 text-right text-gray-400">20%</td>
                    <td class="px-3 py-2 text-right font-mono text-emerald-300/90 tabular-nums">
                      {{ b.reworkTerm.toFixed(2) }}
                    </td>
                  </tr>
                  <tr :class="rowHighlight('velocity')">
                    <td class="px-3 py-2">
                      <div class="text-white font-medium">Velocity norm</div>
                      <div class="text-gray-500 tabular-nums">min(SP × 5, 100), SP = {{ member.sprint_velocity_sp.toFixed(2) }} → {{ b.velocityNorm.toFixed(1) }}</div>
                    </td>
                    <td class="px-3 py-2 text-right text-gray-400">15%</td>
                    <td class="px-3 py-2 text-right font-mono text-emerald-300/90 tabular-nums">
                      {{ b.velocityTerm.toFixed(2) }}
                    </td>
                  </tr>
                  <tr :class="rowHighlight('time')">
                    <td class="px-3 py-2">
                      <div class="text-white font-medium">Time accuracy</div>
                      <div class="text-gray-500 tabular-nums">{{ member.time_accuracy_pct.toFixed(1) }}%</div>
                      <div class="text-[10px] text-gray-600 mt-0.5">Avg. per task: 1 − |logged − est| / est (0–100%)</div>
                    </td>
                    <td class="px-3 py-2 text-right text-gray-400">10%</td>
                    <td class="px-3 py-2 text-right font-mono text-emerald-300/90 tabular-nums">
                      {{ b.timeTerm.toFixed(2) }}
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>

            <div class="flex items-center justify-between rounded-lg border border-gray-600 bg-gray-800/40 px-3 py-2">
              <span class="text-gray-400">Sum of contributions</span>
              <span class="font-mono font-bold text-white tabular-nums">{{ b.sum.toFixed(2) }}</span>
            </div>
            <div class="flex items-center justify-between text-xs text-gray-500">
              <span>API composite_score</span>
              <span class="font-mono tabular-nums">{{ member.composite_score.toFixed(2) }}</span>
            </div>
            <p v-if="Math.abs(b.sum - member.composite_score) > 0.05" class="text-[11px] text-amber-400/90">
              Minor difference vs API is floating-point rounding.
            </p>

            <div class="border-t border-gray-700 pt-3 space-y-2 text-xs text-gray-500">
              <p><strong class="text-gray-400">Velocity:</strong> average story points completed per sprint in the last 3 completed sprints.</p>
              <p><strong class="text-gray-400">Time accuracy:</strong> for each task with estimate and time logs, accuracy = 1 − |logged − estimated| / estimated (clamped), then averaged × 100.</p>
            </div>
          </template>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import type { PerformanceBreakdownFocus, TeamMemberKPI } from '~/core/modules/performance/performance-api'
import { isEngineerLikeRole } from '~/utils/roles'

const props = defineProps<{
  modelValue: boolean
  member: TeamMemberKPI | null
  focus: PerformanceBreakdownFocus
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', v: boolean): void
}>()

function close() {
  emit('update:modelValue', false)
}

const b = computed(() => breakdownForDev(props.member))

function breakdownForDev(m: TeamMemberKPI | null) {
  const z = {
    deliveryTerm: 0,
    qualityTerm: 0,
    reworkTerm: 0,
    velocityTerm: 0,
    timeTerm: 0,
    qualityNorm: 0,
    reworkNorm: 0,
    velocityNorm: 0,
    sum: 0,
  }
  if (!m || !isEngineerLikeRole(m.role)) return z

  let qualityNorm = m.code_quality_index
  if (qualityNorm > 100) qualityNorm = 100

  let reworkNorm = 100 - m.rework_rate_pct
  if (reworkNorm < 0) reworkNorm = 0

  let velocityNorm = m.sprint_velocity_sp * 5
  if (velocityNorm > 100) velocityNorm = 100

  const deliveryTerm = 0.3 * m.delivery_rate_pct
  const qualityTerm = 0.25 * qualityNorm
  const reworkTerm = 0.2 * reworkNorm
  const velocityTerm = 0.15 * velocityNorm
  const timeTerm = 0.1 * m.time_accuracy_pct
  const sum = deliveryTerm + qualityTerm + reworkTerm + velocityTerm + timeTerm

  return {
    deliveryTerm,
    qualityTerm,
    reworkTerm,
    velocityTerm,
    timeTerm,
    qualityNorm,
    reworkNorm,
    velocityNorm,
    sum,
  }
}

function rowHighlight(part: 'delivery' | 'quality' | 'rework' | 'velocity' | 'time') {
  const map: Record<PerformanceBreakdownFocus, string | null> = {
    composite: null,
    delivery: 'delivery',
    quality: 'quality',
    rework: 'rework',
  }
  const hit = map[props.focus]
  if (!hit) return ''
  return hit === part ? 'bg-cyan-950/40 ring-1 ring-inset ring-cyan-500/30' : 'opacity-80'
}

function sectionHighlightClass(part: PerformanceBreakdownFocus) {
  return props.focus === part ? 'bg-cyan-950/30 ring-1 ring-cyan-500/25' : ''
}

function onKey(e: KeyboardEvent) {
  if (e.key === 'Escape' && props.modelValue) close()
}

watch(
  () => props.modelValue,
  (open) => {
    if (import.meta.client) {
      if (open) {
        document.addEventListener('keydown', onKey)
        document.body.style.overflow = 'hidden'
      } else {
        document.removeEventListener('keydown', onKey)
        document.body.style.overflow = ''
      }
    }
  },
)

onUnmounted(() => {
  if (import.meta.client) {
    document.removeEventListener('keydown', onKey)
    document.body.style.overflow = ''
  }
})
</script>

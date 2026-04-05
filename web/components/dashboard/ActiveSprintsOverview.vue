<template>
  <section>
    <h2 class="section-label">Active Sprints — Squad Projects</h2>

    <div v-if="isLoading" class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
      <div
        v-for="i in 3"
        :key="i"
        class="h-40 animate-pulse rounded-2xl border border-gray-700 bg-gray-800/60"
      />
    </div>

    <div v-else-if="activeSprintCards.length === 0" class="rounded-2xl border border-gray-700 bg-gray-800/60 p-12 text-center">
      <svg class="h-10 w-10 text-gray-600 mx-auto mb-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 17V7m0 10a2 2 0 01-2 2H5a2 2 0 01-2-2V7a2 2 0 012-2h2a2 2 0 012 2m0 10a2 2 0 002 2h2a2 2 0 002-2M9 7a2 2 0 012-2h2a2 2 0 012 2m0 10V7m0 10a2 2 0 002 2h2a2 2 0 002-2V7a2 2 0 00-2-2h-2a2 2 0 00-2 2"/>
      </svg>
      <p class="text-sm text-gray-400 font-medium">No active sprints</p>
      <p class="text-xs text-gray-600 mt-0.5">Start a sprint in your team's projects to see it here</p>
    </div>

    <div v-else class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
      <div
        v-for="card in activeSprintCards"
        :key="card.sprint.id"
        class="rounded-2xl border bg-gray-800/60 p-5 shadow-lg hover:bg-gray-700/40 transition-colors cursor-pointer"
        :style="{ borderColor: card.project.color ? `${card.project.color}40` : undefined }"
        :class="!card.project.color ? 'border-gray-700' : ''"
        @click="navigateTo(`/projects/${card.project.id}`)"
      >
        <!-- Project + Sprint Header -->
        <div class="flex items-start justify-between gap-2 mb-3">
          <div class="min-w-0 flex-1">
            <div class="flex items-center gap-2 mb-0.5">
              <span
                v-if="card.project.color"
                class="w-2.5 h-2.5 rounded-full flex-shrink-0"
                :style="{ backgroundColor: card.project.color }"
              />
              <p class="text-xs font-medium text-gray-400 truncate">{{ card.project.name }}</p>
            </div>
            <p class="text-sm font-bold text-white truncate">{{ card.sprint.name }}</p>
            <p v-if="card.sprint.goal" class="text-xs text-gray-500 mt-0.5 line-clamp-2">{{ card.sprint.goal }}</p>
          </div>
          <span class="flex-shrink-0 px-2 py-0.5 rounded-full text-xs font-semibold bg-emerald-500/15 border border-emerald-500/30 text-emerald-400">
            ACTIVE
          </span>
        </div>

        <!-- Progress -->
        <div class="mb-3">
          <div class="flex items-center justify-between mb-1.5">
            <span class="text-xs text-gray-500">Progress</span>
            <span class="text-xs font-semibold tabular-nums" :class="progressColor(card.completedPct)">
              {{ card.completedPct.toFixed(0) }}%
            </span>
          </div>
          <div class="h-1.5 w-full rounded-full bg-gray-700 overflow-hidden">
            <div
              class="h-full rounded-full transition-all duration-500"
              :class="progressBarColor(card.completedPct)"
              :style="{ width: `${card.completedPct}%` }"
            />
          </div>
        </div>

        <!-- Task Count Breakdown -->
        <div class="flex items-center gap-3 text-xs flex-wrap">
          <span class="text-gray-500">
            <span class="font-semibold text-white">{{ card.completedCount }}</span>/{{ card.totalCount }} tasks
          </span>
          <template v-if="card.blockedCount > 0">
            <span class="h-3 w-px bg-gray-700"/>
            <span class="font-semibold text-red-400">{{ card.blockedCount }} blocked</span>
          </template>
          <template v-if="card.inProgressCount > 0">
            <span class="h-3 w-px bg-gray-700"/>
            <span class="font-semibold text-purple-400">{{ card.inProgressCount }} in progress</span>
          </template>
        </div>

        <!-- Dates -->
        <div v-if="card.sprint.start_date || card.sprint.end_date" class="flex items-center gap-2 mt-3 pt-3 border-t border-gray-700/60">
          <svg class="w-3 h-3 text-gray-600 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/>
          </svg>
          <span class="text-xs text-gray-500">
            {{ formatDate(card.sprint.start_date) }} → {{ formatDate(card.sprint.end_date) }}
          </span>
          <span v-if="card.daysLeft !== null" class="ml-auto text-xs" :class="card.daysLeft <= 2 ? 'text-red-400 font-semibold' : 'text-gray-500'">
            {{ card.daysLeft > 0 ? `${card.daysLeft}d left` : 'Overdue' }}
          </span>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { useProjectsApi } from '~/core/modules/projects/infrastructure/projects-api'
import type { Project, Sprint } from '~/core/modules/projects/infrastructure/projects-api'

const props = defineProps<{
  /** When set (Product Owner dashboard), avoids a duplicate getProjects() round-trip */
  projects: Project[]
}>()

interface SprintCard {
  project: Project
  sprint: Sprint
  totalCount: number
  completedCount: number
  inProgressCount: number
  blockedCount: number
  completedPct: number
  daysLeft: number | null
}

const { getSprints } = useProjectsApi()

const isLoading = ref(true)
const activeSprintCards = ref<SprintCard[]>([])

const progressColor = (pct: number) => {
  if (pct >= 80) return 'text-emerald-400'
  if (pct >= 40) return 'text-blue-400'
  return 'text-amber-400'
}

const progressBarColor = (pct: number) => {
  if (pct >= 80) return 'bg-emerald-500'
  if (pct >= 40) return 'bg-blue-500'
  return 'bg-amber-500'
}

const formatDate = (d: string | null) => {
  if (!d) return '—'
  return new Date(d).toLocaleDateString('en-US', { month: 'short', day: 'numeric' })
}

const getDaysLeft = (endDate: string | null): number | null => {
  if (!endDate) return null
  const diff = new Date(endDate).getTime() - Date.now()
  return Math.ceil(diff / 86400000)
}

const buildSprintCard = (project: Project, sprint: Sprint): SprintCard => {
  const tasks = sprint.tasks ?? []
  const total = tasks.length
  const completed = tasks.filter(t => t.status === 'COMPLETED').length
  const inProgress = tasks.filter(t => t.status === 'IN_PROGRESS').length
  const blocked = tasks.filter(t => t.status === 'BLOCKED').length
  return {
    project,
    sprint,
    totalCount: total,
    completedCount: completed,
    inProgressCount: inProgress,
    blockedCount: blocked,
    completedPct: total > 0 ? (completed / total) * 100 : 0,
    daysLeft: getDaysLeft(sprint.end_date),
  }
}

const fetchData = async () => {
  isLoading.value = true
  try {
    const projects = props.projects
    const activeProjects = projects.filter(p => p.status === 'ACTIVE')

    const sprintResults = await Promise.all(
      activeProjects.map(p => getSprints(p.id).then(sprints => ({ project: p, sprints })).catch(() => ({ project: p, sprints: [] })))
    )

    const cards: SprintCard[] = []
    for (const { project, sprints } of sprintResults) {
      for (const sprint of sprints) {
        if (sprint.status === 'ACTIVE') {
          cards.push(buildSprintCard(project, sprint))
        }
      }
    }
    activeSprintCards.value = cards
  } catch {
    // silent
  } finally {
    isLoading.value = false
  }
}

watch(
  () => props.projects.map(p => p.id).join(','),
  () => fetchData(),
  { immediate: true, flush: 'post' },
)
</script>

<style scoped>
.section-label {
  @apply text-xs font-semibold uppercase tracking-widest text-gray-500 mb-4;
}
</style>

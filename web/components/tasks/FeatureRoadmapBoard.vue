<template>
  <div class="w-full">
    <!-- Header -->
    <div class="flex items-center justify-between mb-6">
      <div class="flex items-center gap-3">
        <div class="w-8 h-8 rounded-lg bg-violet-500/15 border border-violet-500/30 flex items-center justify-center">
          <svg class="w-4 h-4 text-violet-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 20l-5.447-2.724A1 1 0 013 16.382V5.618a1 1 0 011.447-.894L9 7m0 13l6-3m-6 3V7m6 10l4.553 2.276A1 1 0 0021 18.382V7.618a1 1 0 00-.553-.894L15 4m0 13V4m0 0L9 7"/>
          </svg>
        </div>
        <div>
          <h2 class="text-sm font-bold text-white">Feature Roadmap Board</h2>
          <p class="text-xs text-gray-500">
            <template v-if="isProjectScoped && scopeProjectName">{{ scopeProjectName }} — features and delivery progress</template>
            <template v-else>Strategic features with roll-up delivery progress (Product Owner / CEO / Manager)</template>
          </p>
        </div>
      </div>

      <div class="flex flex-wrap items-center gap-3">
        <!-- Project filter (hidden when embedded on a project page) -->
        <div v-if="!isProjectScoped" class="flex items-center gap-2">
          <label class="text-[10px] font-bold uppercase tracking-wider text-gray-500 shrink-0">Project</label>
          <select
            v-model="selectedProjectId"
            class="min-w-[140px] max-w-[220px] rounded-lg border border-gray-700 bg-gray-800/80 px-2.5 py-2 text-xs font-medium text-gray-200 focus:border-violet-500/50 focus:outline-none focus:ring-1 focus:ring-violet-500/30"
          >
            <option value="">All projects</option>
            <option v-for="p in projectOptions" :key="p.id" :value="p.id">{{ p.name }}</option>
          </select>
        </div>

        <!-- Status filter -->
        <div class="flex rounded-lg border border-gray-700 bg-gray-800/60 overflow-hidden text-xs font-semibold">
          <button
            v-for="f in statusFilters"
            :key="f.value"
            @click="activeFilter = f.value"
            class="px-3 py-2 transition-colors"
            :class="activeFilter === f.value ? 'bg-violet-600 text-white' : 'text-gray-400 hover:text-white hover:bg-gray-700'"
          >{{ f.label }}</button>
        </div>

        <button
          @click="load"
          :disabled="loading"
          class="inline-flex items-center gap-2 px-3 py-2 rounded-lg border border-gray-700 bg-gray-800/60 text-xs font-medium text-gray-300 hover:border-gray-600 hover:bg-gray-700 hover:text-white transition-colors disabled:opacity-50"
        >
          <svg class="h-3.5 w-3.5" :class="loading ? 'animate-spin' : ''" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"/>
          </svg>
          Refresh
        </button>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="flex items-center justify-center py-16">
      <svg class="h-7 w-7 animate-spin text-violet-400 mr-3" fill="none" viewBox="0 0 24 24">
        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/>
      </svg>
      <span class="text-sm text-gray-500">Loading feature roadmap…</span>
    </div>

    <!-- Error -->
    <div v-else-if="error" class="flex items-start gap-3 rounded-xl border border-red-500/30 bg-red-900/20 px-5 py-4 text-red-400">
      <svg class="h-5 w-5 flex-shrink-0 mt-0.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01M10.29 3.86L1.82 18a2 2 0 001.71 3h16.94a2 2 0 001.71-3L13.71 3.86a2 2 0 00-3.42 0z"/>
      </svg>
      <div>
        <p class="text-sm font-semibold">Failed to load feature roadmap</p>
        <p class="text-xs text-red-300 mt-0.5">{{ error }}</p>
      </div>
    </div>

    <!-- Empty state -->
    <div
      v-else-if="filteredFeatures.length === 0"
      class="flex flex-col items-center justify-center rounded-2xl border border-dashed border-gray-700 bg-gray-800/30 py-16 text-center"
    >
      <div class="w-12 h-12 rounded-2xl bg-gray-700/50 flex items-center justify-center mb-3">
        <svg class="h-6 w-6 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 20l-5.447-2.724A1 1 0 013 16.382V5.618a1 1 0 011.447-.894L9 7m0 13l6-3m-6 3V7m6 10l4.553 2.276A1 1 0 0021 18.382V7.618a1 1 0 00-.553-.894L15 4m0 13V4m0 0L9 7"/>
        </svg>
      </div>
      <p class="text-sm font-semibold text-gray-300">No features found</p>
      <p class="text-xs text-gray-500 mt-1">Create FEATURE-type tasks and assign them to projects</p>
    </div>

    <!-- Summary stats -->
    <div v-else class="grid grid-cols-2 sm:grid-cols-4 gap-3 mb-6">
      <div class="rounded-xl border border-gray-700/60 bg-gray-800/50 px-4 py-3">
        <p class="text-xs font-semibold uppercase tracking-wider text-gray-500 mb-1">Total Features</p>
        <p class="text-xl font-black text-white tabular-nums">{{ filteredFeatures.length }}</p>
      </div>
      <div class="rounded-xl border border-emerald-500/30 bg-emerald-950/20 px-4 py-3">
        <p class="text-xs font-semibold uppercase tracking-wider text-gray-500 mb-1">Shipped</p>
        <p class="text-xl font-black text-emerald-400 tabular-nums">{{ shippedCount }}</p>
      </div>
      <div class="rounded-xl border border-amber-500/30 bg-amber-950/20 px-4 py-3">
        <p class="text-xs font-semibold uppercase tracking-wider text-gray-500 mb-1">Awaiting UAT</p>
        <p class="text-xl font-black text-amber-400 tabular-nums">{{ awaitingUATCount }}</p>
      </div>
      <div class="rounded-xl border border-violet-500/30 bg-violet-950/20 px-4 py-3">
        <p class="text-xs font-semibold uppercase tracking-wider text-gray-500 mb-1">Avg. Progress</p>
        <p class="text-xl font-black text-violet-400 tabular-nums">{{ avgProgress }}%</p>
      </div>
    </div>

    <!-- Feature list (grouped by project when viewing all projects) -->
    <div v-if="!loading && !error && filteredFeatures.length > 0" class="space-y-8">
      <section v-for="group in featureGroups" :key="group.key" class="space-y-3">
        <div
          v-if="!hideProjectGroupHeaders"
          class="flex items-center gap-2 pb-2 border-b border-gray-700/50"
        >
          <span class="w-2 h-2 rounded-full shrink-0" :style="{ backgroundColor: group.color }"/>
          <h3 class="text-xs font-bold uppercase tracking-wider text-gray-300">{{ group.name }}</h3>
          <span class="text-[10px] font-medium text-gray-600 tabular-nums">{{ group.features.length }} feature{{ group.features.length !== 1 ? 's' : '' }}</span>
        </div>
        <div class="space-y-3">
      <div
        v-for="feature in group.features"
        :key="feature.id"
        class="rounded-2xl border bg-gray-800/40 overflow-hidden transition-all"
        :class="featureBorderClass(feature)"
      >
        <!-- Feature card row -->
        <div class="px-5 py-4">
          <!-- Top row: project + priority + status -->
          <div class="flex items-center justify-between gap-3 mb-3">
            <div class="flex items-center gap-2 min-w-0">
              <!-- Project pill → links to project Kanban board -->
              <NuxtLink
                :to="`/projects/${feature.project_id}?tab=board`"
                class="inline-flex items-center gap-1.5 rounded-full border px-2 py-0.5 text-[10px] font-semibold leading-none shrink-0 hover:opacity-80 transition-opacity"
                :style="{
                  borderColor: feature.project_color || '#6366f1',
                  color: feature.project_color || '#6366f1',
                  backgroundColor: (feature.project_color || '#6366f1') + '18'
                }"
                @click.stop
              >
                <span class="w-1.5 h-1.5 rounded-full" :style="{ backgroundColor: feature.project_color || '#6366f1' }"/>
                {{ feature.project_name || 'Unknown Project' }}
              </NuxtLink>

              <!-- FEATURE badge -->
              <span class="text-[9px] font-bold uppercase tracking-wider px-1.5 py-0.5 rounded border border-violet-500/40 text-violet-400 bg-violet-500/10 shrink-0">
                FEATURE
              </span>
            </div>

            <div class="flex items-center gap-2 shrink-0">
              <!-- Priority -->
              <span class="text-[10px] font-bold uppercase tracking-wide px-1.5 py-0.5 rounded" :class="priorityClass(feature.priority)">
                {{ feature.priority }}
              </span>
              <!-- Status pill -->
              <span class="text-[10px] font-bold uppercase tracking-wider px-2 py-0.5 rounded-full border" :class="statusBadgeClass(feature.status)">
                {{ statusLabel(feature.status) }}
              </span>
            </div>
          </div>

          <!-- Title + code -->
          <div class="flex items-start gap-3 mb-3">
            <div class="min-w-0 flex-1">
              <p
                class="text-sm font-bold text-white leading-snug cursor-pointer hover:text-violet-300 transition-colors"
                @click="goToFeature(feature)"
              >{{ feature.title }}</p>
              <p v-if="feature.code" class="text-[10px] font-mono text-gray-600 mt-0.5">{{ feature.code }}</p>
            </div>
          </div>

          <!-- Roll-up progress bar -->
          <div class="mb-3">
            <div class="flex items-center justify-between text-xs mb-1.5">
              <span class="text-gray-500 font-medium">Delivery Progress</span>
              <span class="font-bold tabular-nums" :class="progressTextClass(feature.rollup_progress)">
                {{ feature.rollup_progress }}%
                <span class="text-gray-600 font-normal ml-1">({{ completedChildCount(feature) }}/{{ feature.child_tasks?.length ?? 0 }} tasks done)</span>
              </span>
            </div>
            <div class="relative h-2 bg-gray-700/80 rounded-full overflow-hidden">
              <div
                v-if="feature.rollup_progress >= 75"
                class="absolute inset-0 rounded-full blur-sm opacity-40"
                :class="progressBarClass(feature.rollup_progress)"
                :style="{ width: feature.rollup_progress + '%' }"
              />
              <div
                class="h-full rounded-full transition-all duration-700"
                :class="progressBarClass(feature.rollup_progress)"
                :style="{ width: feature.rollup_progress + '%' }"
              />
            </div>
          </div>

          <!-- Meta row + action buttons + accordion toggle -->
          <div class="flex items-center justify-between gap-2 text-[10px] text-gray-500">
            <div class="flex items-center gap-3">
              <span v-if="feature.due_at" :class="deadlineClass(feature.due_at)">
                Due {{ formatDeadline(feature.due_at) }}
              </span>
              <span v-else class="text-gray-700">No deadline</span>

              <span v-if="feature.story_points" class="text-gray-600">
                SP: <span class="font-bold text-gray-400">{{ feature.story_points }}</span>
              </span>
            </div>

            <div class="flex items-center gap-2">
              <!-- DEV: Submit for UAT button (shown when feature is READY_FOR_UAT) -->
              <button
                v-if="isDev && feature.status === 'READY_FOR_UAT'"
                @click="openSubmitUAT(feature)"
                class="inline-flex items-center gap-1.5 rounded-lg border border-amber-500/40 bg-amber-500/10 px-2.5 py-1.5 text-[10px] font-bold text-amber-400 hover:border-amber-400/60 hover:bg-amber-500/20 transition-all"
              >
                <svg class="h-3 w-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7l5 5m0 0l-5 5m5-5H6"/>
                </svg>
                Submit for UAT
              </button>

              <!-- Product Owner / CEO: Review UAT button (shown when feature is REVIEW_PENDING) -->
              <button
                v-if="isPMOrCEO && feature.status === 'REVIEW_PENDING'"
                @click="openUATReview(feature)"
                class="inline-flex items-center gap-1.5 rounded-lg border border-violet-500/40 bg-violet-500/10 px-2.5 py-1.5 text-[10px] font-bold text-violet-300 hover:border-violet-400/60 hover:bg-violet-500/20 transition-all"
              >
                <svg class="h-3 w-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
                </svg>
                Review UAT
              </button>

              <!-- Accordion toggle -->
              <button
                v-if="(feature.child_tasks?.length ?? 0) > 0"
                @click="toggleAccordion(feature.id)"
                class="inline-flex items-center gap-1.5 rounded-lg border border-gray-700 bg-gray-700/40 px-2.5 py-1.5 text-[10px] font-semibold text-gray-400 hover:border-violet-500/40 hover:text-violet-300 hover:bg-violet-500/10 transition-all"
              >
                <svg
                  class="h-3 w-3 transition-transform duration-200"
                  :class="expandedFeatures.has(feature.id) ? 'rotate-180' : ''"
                  fill="none" stroke="currentColor" viewBox="0 0 24 24"
                >
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M19 9l-7 7-7-7"/>
                </svg>
                {{ expandedFeatures.has(feature.id) ? 'Hide' : 'Show' }} {{ feature.child_tasks.length }} child task{{ feature.child_tasks.length !== 1 ? 's' : '' }}
              </button>
              <span v-else class="text-gray-700 italic">No child tasks yet</span>
            </div>
          </div>
        </div>

        <!-- Accordion: child task list -->
        <Transition
          enter-active-class="transition-all duration-200 ease-out"
          enter-from-class="opacity-0 -translate-y-1"
          enter-to-class="opacity-100 translate-y-0"
          leave-active-class="transition-all duration-150 ease-in"
          leave-from-class="opacity-100 translate-y-0"
          leave-to-class="opacity-0 -translate-y-1"
        >
          <div
            v-if="expandedFeatures.has(feature.id) && (feature.child_tasks?.length ?? 0) > 0"
            class="border-t border-gray-700/60 bg-gray-900/40"
          >
            <div class="px-5 py-3">
              <p class="text-[9px] font-bold uppercase tracking-widest text-gray-600 mb-2">Child Tasks</p>
              <div class="space-y-1.5">
                <div
                  v-for="child in feature.child_tasks"
                  :key="child.id"
                  class="flex items-center gap-3 rounded-lg px-3 py-2 border cursor-pointer hover:bg-gray-800/60 transition-colors"
                  :class="childRowBorderClass(child.status)"
                  @click="goToTask(child)"
                >
                  <!-- Status icon -->
                  <div class="shrink-0 w-4 h-4 flex items-center justify-center">
                    <svg v-if="child.status === 'COMPLETED'" class="w-4 h-4 text-emerald-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M5 13l4 4L19 7"/>
                    </svg>
                    <div v-else-if="child.status === 'IN_PROGRESS'" class="w-2.5 h-2.5 rounded-full bg-violet-400 animate-pulse"/>
                    <div v-else class="w-2.5 h-2.5 rounded-full border-2 border-gray-600"/>
                  </div>

                  <!-- Type badge -->
                  <span
                    class="text-[9px] font-bold uppercase tracking-wider px-1.5 py-0.5 rounded border shrink-0"
                    :class="child.task_type === 'BUG'
                      ? 'border-red-500/40 text-red-400 bg-red-500/10'
                      : 'border-indigo-500/40 text-indigo-400 bg-indigo-500/10'"
                  >{{ child.task_type || 'TASK' }}</span>

                  <!-- Title -->
                  <span class="text-xs text-gray-300 flex-1 truncate" :class="child.status === 'COMPLETED' ? 'line-through text-gray-600' : ''">
                    {{ child.title }}
                  </span>

                  <!-- Assignee badge -->
                  <span v-if="child.assigned_to_email || child.assigned_to_display_name" class="text-[10px] text-gray-500 shrink-0 hidden sm:block">
                    {{ child.assigned_to_display_name || child.assigned_to_email }}
                  </span>

                  <!-- Est. hours -->
                  <span v-if="child.estimated_minutes" class="text-[10px] text-indigo-400/80 shrink-0 tabular-nums">
                    {{ (child.estimated_minutes / 60).toFixed(1) }}h
                  </span>

                  <!-- Status pill -->
                  <span class="text-[9px] font-bold uppercase tracking-wider px-1.5 py-0.5 rounded-full border shrink-0" :class="statusBadgeClass(child.status)">
                    {{ statusLabel(child.status) }}
                  </span>
                </div>
              </div>
            </div>
          </div>
        </Transition>
      </div>
        </div>
      </section>
    </div>

    <!-- Modals -->
    <SubmitUATModal
      v-if="submitUATModalOpen"
      v-model="submitUATModalOpen"
      :feature="selectedFeature"
      @success="load"
    />

    <UATReviewModal
      v-if="uatReviewModalOpen"
      v-model="uatReviewModalOpen"
      :feature="selectedFeature"
      @success="load"
    />
  </div>
</template>

<script setup lang="ts">
import SubmitUATModal from '~/components/tasks/SubmitUATModal.vue'
import UATReviewModal from '~/components/tasks/UATReviewModal.vue'
import { useTasksApi } from '~/core/modules/tasks/infrastructure/tasks-api'
import type { FeatureRoadmapItem } from '~/core/modules/tasks/infrastructure/tasks-api'
import { useAuth } from '~/composables/useAuth'
import { isEngineerLikeRole } from '~/utils/roles'

const props = withDefaults(
  defineProps<{
    /** When set (e.g. on project detail page), only show features for this project. */
    scopeProjectId?: string
    scopeProjectCode?: string
    scopeProjectName?: string
    /** Optional preloaded feature list from parent page to avoid duplicate fetch. */
    prefetchedFeatures?: FeatureRoadmapItem[] | null
  }>(),
  {
    scopeProjectId: '',
    scopeProjectCode: '',
    scopeProjectName: '',
    prefetchedFeatures: null,
  }
)

const statusFilters = [
  { value: 'ALL', label: 'All' },
  { value: 'PENDING', label: 'Pending' },
  { value: 'IN_PROGRESS', label: 'In Progress' },
  { value: 'READY_FOR_UAT', label: 'Ready for UAT' },
  { value: 'REVIEW_PENDING', label: 'UAT Review' },
  { value: 'COMPLETED', label: 'Shipped' },
]

const { getActiveFeatures } = useTasksApi()
const { currentUser } = useAuth()

const features = ref<FeatureRoadmapItem[]>([])
const loading = ref(false)
const error = ref<string | null>(null)
const activeFilter = ref('ALL')
const expandedFeatures = ref<Set<string>>(new Set())

const selectedFeature = ref<FeatureRoadmapItem | null>(null)
const submitUATModalOpen = ref(false)
const uatReviewModalOpen = ref(false)

const userRole = computed(() => (currentUser.value?.role || '').toUpperCase())
const isDev = computed(() => isEngineerLikeRole(userRole.value))
/** Roadmap UAT review / approve — same status powers as PO & CEO for features in review */
const isPMOrCEO = computed(() =>
  userRole.value === 'PRODUCT_OWNER' || userRole.value === 'PM' || userRole.value === 'CEO' || userRole.value === 'MANAGER')

const isProjectScoped = computed(
  () => Boolean(props.scopeProjectId?.trim() || props.scopeProjectCode?.trim())
)

/** Stable project key for grouping / filter (API may use project_id or code). */
function projectKey(f: FeatureRoadmapItem): string {
  return String(f.project_id || f.project_code || '').trim() || '__none__'
}

function matchesScope(f: FeatureRoadmapItem): boolean {
  const id = props.scopeProjectId?.trim()
  const code = (props.scopeProjectCode || '').trim()
  const pk = projectKey(f)
  if (id && (f.project_id === id || pk === id)) return true
  if (code && (f.project_code === code || pk.toLowerCase() === code.toLowerCase())) return true
  return false
}

const hideProjectGroupHeaders = computed(
  () => isProjectScoped.value || Boolean(selectedProjectId.value)
)

const projectOptions = computed(() => {
  const m = new Map<string, { id: string; name: string }>()
  for (const f of features.value) {
    const id = projectKey(f)
    if (id === '__none__') continue
    if (!m.has(id)) m.set(id, { id, name: f.project_name || 'Unknown project' })
  }
  return Array.from(m.values()).sort((a, b) => a.name.localeCompare(b.name))
})

const selectedProjectId = ref('')

const filteredFeatures = computed(() => {
  let list = features.value
  if (activeFilter.value !== 'ALL') {
    list = list.filter(f => f.status === activeFilter.value)
  }
  if (isProjectScoped.value) {
    list = list.filter(f => matchesScope(f))
  } else if (selectedProjectId.value) {
    list = list.filter(f => projectKey(f) === selectedProjectId.value)
  }
  return list
})

const shippedCount = computed(() => filteredFeatures.value.filter(f => f.status === 'COMPLETED').length)
const awaitingUATCount = computed(() =>
  filteredFeatures.value.filter(f => f.status === 'READY_FOR_UAT' || f.status === 'REVIEW_PENDING').length
)
const avgProgress = computed(() => {
  if (!filteredFeatures.value.length) return 0
  return Math.round(filteredFeatures.value.reduce((s, f) => s + f.rollup_progress, 0) / filteredFeatures.value.length)
})

/** When showing all projects, group cards under project headings. */
const featureGroups = computed(() => {
  const list = filteredFeatures.value
  if (isProjectScoped.value || selectedProjectId.value) {
    const sample = list[0]
    const name =
      (isProjectScoped.value && props.scopeProjectName?.trim()) ||
      sample?.project_name ||
      'Project'
    const color = sample?.project_color || '#6366f1'
    const key = sample ? projectKey(sample) : (props.scopeProjectId || props.scopeProjectCode || 'scoped')
    return [{ key, name, color, features: list }]
  }
  const map = new Map<string, { key: string; name: string; color: string; features: FeatureRoadmapItem[] }>()
  for (const f of list) {
    const key = projectKey(f)
    const name = f.project_name || (key === '__none__' ? 'No project' : 'Unknown project')
    const color = f.project_color || '#6366f1'
    if (!map.has(key)) {
      map.set(key, { key, name, color, features: [] })
    }
    map.get(key)!.features.push(f)
  }
  return Array.from(map.values()).sort((a, b) => a.name.localeCompare(b.name))
})

function completedChildCount(feature: FeatureRoadmapItem): number {
  return feature.child_tasks?.filter(c => c.status === 'COMPLETED').length ?? 0
}

function toggleAccordion(featureId: string) {
  const s = new Set(expandedFeatures.value)
  if (s.has(featureId)) s.delete(featureId)
  else s.add(featureId)
  expandedFeatures.value = s
}

function openSubmitUAT(feature: FeatureRoadmapItem) {
  selectedFeature.value = feature
  submitUATModalOpen.value = true
}

function openUATReview(feature: FeatureRoadmapItem) {
  selectedFeature.value = feature
  uatReviewModalOpen.value = true
}

async function load() {
  loading.value = true
  error.value = null
  try {
    if (props.prefetchedFeatures) {
      features.value = props.prefetchedFeatures
      return
    }
    const scopedProjectId = props.scopeProjectId?.trim()
    features.value = await getActiveFeatures(scopedProjectId || undefined)
  } catch (e: any) {
    error.value = e?.data?.message || e?.message || 'Failed to load feature roadmap'
  } finally {
    loading.value = false
  }
}

function goToFeature(feature: FeatureRoadmapItem) {
  navigateTo(`/task/${feature.code || feature.id}?from=roadmap`)
}

function goToTask(task: { id: string; code?: string }) {
  navigateTo(`/task/${task.code || task.id}?from=roadmap`)
}

function featureBorderClass(feature: FeatureRoadmapItem): string {
  if (feature.status === 'READY_FOR_UAT') return 'border-amber-500/40'
  if (feature.status === 'REVIEW_PENDING') return 'border-violet-500/50'
  if (feature.rollup_progress === 100) return 'border-emerald-500/30'
  if (feature.rollup_progress >= 50) return 'border-violet-500/30'
  if (feature.rollup_progress > 0) return 'border-indigo-500/20'
  return 'border-gray-700/60'
}

function progressBarClass(pct: number): string {
  if (pct === 100) return 'bg-emerald-500'
  if (pct >= 75) return 'bg-violet-500'
  if (pct >= 40) return 'bg-indigo-500'
  return 'bg-gray-500'
}

function progressTextClass(pct: number): string {
  if (pct === 100) return 'text-emerald-400'
  if (pct >= 75) return 'text-violet-400'
  if (pct >= 40) return 'text-indigo-400'
  return 'text-gray-400'
}

function statusLabel(status: string): string {
  switch (status) {
    case 'READY_FOR_UAT': return 'READY FOR UAT'
    case 'REVIEW_PENDING': return 'UAT REVIEW'
    default: return status?.replace(/_/g, ' ') || ''
  }
}

function statusBadgeClass(status: string): string {
  switch (status) {
    case 'COMPLETED': return 'border-emerald-500/40 text-emerald-400 bg-emerald-500/10'
    case 'IN_PROGRESS': return 'border-violet-500/40 text-violet-400 bg-violet-500/10'
    case 'READY_FOR_UAT': return 'border-amber-500/40 text-amber-400 bg-amber-500/10'
    case 'REVIEW_PENDING': return 'border-violet-400/50 text-violet-300 bg-violet-500/15'
    default: return 'border-gray-600/60 text-gray-400 bg-gray-700/30'
  }
}

function childRowBorderClass(status: string): string {
  switch (status) {
    case 'COMPLETED': return 'border-emerald-700/30 bg-emerald-950/10'
    case 'IN_PROGRESS': return 'border-violet-700/30'
    default: return 'border-gray-700/40'
  }
}

function priorityClass(priority: string): string {
  switch (priority) {
    case 'CRITICAL': return 'bg-red-500/20 text-red-400'
    case 'HIGH': return 'bg-orange-500/20 text-orange-400'
    case 'MEDIUM': return 'bg-amber-500/20 text-amber-400'
    default: return 'bg-gray-500/20 text-gray-400'
  }
}

function deadlineClass(dueAt: string): string {
  const diff = new Date(dueAt).getTime() - Date.now()
  if (diff < 0) return 'text-red-400 font-semibold'
  if (diff < 3 * 24 * 3600 * 1000) return 'text-amber-400 font-semibold'
  return 'text-gray-500'
}

function formatDeadline(dueAt: string): string {
  return new Date(dueAt).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
}

onMounted(() => {
  load()
})
</script>

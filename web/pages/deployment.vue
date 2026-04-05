<template>
  <div class="min-h-screen bg-gray-900 text-gray-100">
    <!-- Header Banner -->
    <div class="relative overflow-hidden border-b border-gray-700/60 bg-gradient-to-r from-gray-900 via-slate-900 to-gray-900">
      <div class="absolute inset-0 bg-gradient-to-r from-cyan-900/20 via-blue-900/20 to-indigo-900/20 pointer-events-none" />
      <div class="relative px-6 py-6 flex flex-col gap-1 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <div class="flex items-center gap-3 mb-1">
            <span class="text-2xl">🚀</span>
            <h1 class="text-2xl font-bold tracking-tight bg-gradient-to-r from-cyan-400 to-blue-400 bg-clip-text text-transparent">
              {{ isChiefEngineer ? 'Deployment Command Center' : 'Deployment Requests' }}
            </h1>
          </div>
          <p class="text-sm text-gray-400 ml-10">
            <span v-if="isChiefEngineer">Review, approve, and merge code ready for deployment</span>
            <span v-else>Submit a deployment request for Chief Engineer review</span>
          </p>
        </div>
        <div class="flex items-center gap-3 ml-10 sm:ml-0">
          <!-- Stats chips (Chief Engineer) -->
          <template v-if="isChiefEngineer && stats">
            <div class="flex items-center gap-2 flex-wrap">
              <div v-if="stats.total_pending > 0" class="flex items-center gap-1.5 px-3 py-1.5 rounded-full bg-yellow-500/10 border border-yellow-500/30 text-yellow-400 text-xs font-bold">
                <span class="w-2 h-2 rounded-full bg-yellow-400 animate-pulse" />
                {{ stats.total_pending }} PENDING
              </div>
              <div v-if="stats.total_reviewing > 0" class="flex items-center gap-1.5 px-3 py-1.5 rounded-full bg-blue-500/10 border border-blue-500/30 text-blue-400 text-xs font-bold">
                <span class="w-2 h-2 rounded-full bg-blue-400 animate-pulse" />
                {{ stats.total_reviewing }} IN REVIEW
              </div>
              <div class="flex items-center gap-1.5 px-3 py-1.5 rounded-full bg-green-500/10 border border-green-500/30 text-green-400 text-xs font-bold">
                🟢 {{ stats.deployed_today }} TODAY
              </div>
            </div>
          </template>
          <!-- Refresh -->
          <button @click="refresh" :disabled="loading" class="p-2 rounded-lg text-gray-400 hover:text-white hover:bg-gray-700 transition-colors disabled:opacity-50">
            <svg class="w-5 h-5" :class="loading ? 'animate-spin' : ''" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
            </svg>
          </button>
        </div>
      </div>
    </div>

    <!-- Chief Engineer: action checklist -->
    <div v-if="isChiefEngineer && (stats?.total_pending ?? 0) + (stats?.total_reviewing ?? 0) > 0"
      class="mx-6 mt-4 px-4 py-3 rounded-xl bg-yellow-500/5 border border-yellow-500/20 flex items-start gap-3">
      <svg class="w-5 h-5 text-yellow-400 mt-0.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
      </svg>
      <p class="text-sm text-yellow-300/80">
        You have <strong class="text-yellow-300">{{ (stats?.total_pending ?? 0) + (stats?.total_reviewing ?? 0) }}</strong> request(s) awaiting your attention.
        Pick up a request, review the branch &amp; PR, then approve or reject before marking it deployed.
      </p>
    </div>

    <!-- Main content -->
    <div class="p-6">

      <!-- Filter tabs -->
      <div class="flex gap-1 p-1 bg-gray-800/60 rounded-xl border border-gray-700/50 w-fit mb-6">
        <button
          v-for="tab in tabs"
          :key="tab.value"
          @click="activeTab = tab.value"
          :class="[
            'px-3 py-1.5 rounded-lg text-xs font-semibold transition-all',
            activeTab === tab.value
              ? 'bg-gray-700 text-white shadow'
              : 'text-gray-400 hover:text-gray-200'
          ]"
        >
          {{ tab.label }}
          <span v-if="tab.count !== undefined" class="ml-1 opacity-70">({{ tab.count }})</span>
        </button>
      </div>

      <!-- Chief Engineer: two-column layout -->
      <div v-if="isChiefEngineer" class="grid grid-cols-1 xl:grid-cols-3 gap-6">
        <!-- Active queue (left, wider) -->
        <div class="xl:col-span-2 space-y-4">
          <h2 class="text-sm font-bold text-gray-400 uppercase tracking-widest">
            {{ activeTab === 'all' || activeTab === 'PENDING' || activeTab === 'REVIEWING' ? 'Active Queue' : 'Requests' }}
          </h2>
          <div v-if="loading && filteredRequests.length === 0" class="space-y-3">
            <div v-for="i in 3" :key="i" class="h-32 rounded-xl bg-gray-800/60 animate-pulse" />
          </div>
          <div v-else-if="filteredRequests.length === 0" class="flex flex-col items-center justify-center py-16 text-center">
            <div class="text-5xl mb-4">🎉</div>
            <p class="text-gray-400 font-medium">No requests here</p>
            <p class="text-gray-600 text-sm mt-1">Everything is clear in this view</p>
          </div>
          <DeploymentRequestCard
            v-for="req in filteredRequests"
            :key="req.id"
            :request="req"
            :is-chief-engineer="true"
            @pick="pickForReview"
            @approve="openApprove"
            @reject="openReject"
            @deploy="confirmDeploy"
          />
        </div>

        <!-- Right sidebar: stats + actions summary -->
        <div class="space-y-4">
          <h2 class="text-sm font-bold text-gray-400 uppercase tracking-widest">Overview</h2>
          <div v-if="stats" class="grid grid-cols-2 gap-3">
            <div class="bg-gray-800/60 border border-gray-700/50 rounded-xl p-4 text-center">
              <div class="text-2xl font-bold text-yellow-400">{{ stats.total_pending }}</div>
              <div class="text-xs text-gray-500 mt-0.5">Pending</div>
            </div>
            <div class="bg-gray-800/60 border border-gray-700/50 rounded-xl p-4 text-center">
              <div class="text-2xl font-bold text-blue-400">{{ stats.total_reviewing }}</div>
              <div class="text-xs text-gray-500 mt-0.5">Reviewing</div>
            </div>
            <div class="bg-gray-800/60 border border-gray-700/50 rounded-xl p-4 text-center">
              <div class="text-2xl font-bold text-green-400">{{ stats.total_deployed }}</div>
              <div class="text-xs text-gray-500 mt-0.5">Deployed</div>
            </div>
            <div class="bg-gray-800/60 border border-gray-700/50 rounded-xl p-4 text-center">
              <div class="text-2xl font-bold text-red-400">{{ stats.total_rejected }}</div>
              <div class="text-xs text-gray-500 mt-0.5">Rejected</div>
            </div>
          </div>

          <!-- Chief Engineer workflow guide -->
          <div class="bg-gray-800/40 border border-gray-700/40 rounded-xl p-4 space-y-3">
            <p class="text-xs font-bold text-gray-400 uppercase tracking-widest">Your Workflow</p>
            <ol class="space-y-2 text-sm">
              <li class="flex items-start gap-2">
                <span class="w-5 h-5 rounded-full bg-yellow-500/20 text-yellow-400 text-[10px] font-bold flex items-center justify-center shrink-0 mt-0.5">1</span>
                <span class="text-gray-300">Pick up a <span class="text-yellow-400 font-medium">PENDING</span> request to start review</span>
              </li>
              <li class="flex items-start gap-2">
                <span class="w-5 h-5 rounded-full bg-blue-500/20 text-blue-400 text-[10px] font-bold flex items-center justify-center shrink-0 mt-0.5">2</span>
                <span class="text-gray-300">Review branch, PR link, and release notes</span>
              </li>
              <li class="flex items-start gap-2">
                <span class="w-5 h-5 rounded-full bg-purple-500/20 text-purple-400 text-[10px] font-bold flex items-center justify-center shrink-0 mt-0.5">3</span>
                <span class="text-gray-300"><span class="text-green-400 font-medium">Approve</span> or <span class="text-red-400 font-medium">Reject</span> with notes</span>
              </li>
              <li class="flex items-start gap-2">
                <span class="w-5 h-5 rounded-full bg-green-500/20 text-green-400 text-[10px] font-bold flex items-center justify-center shrink-0 mt-0.5">4</span>
                <span class="text-gray-300">After merging, click <span class="text-cyan-400 font-medium">Mark Deployed</span></span>
              </li>
            </ol>
          </div>
        </div>
      </div>

      <!-- Engineer / others: my requests + create -->
      <div v-else>
        <div class="flex items-center justify-between mb-4">
          <h2 class="text-sm font-bold text-gray-400 uppercase tracking-widest">My Requests</h2>
          <button
            @click="openCreateModal()"
            class="flex items-center gap-2 px-4 py-2 rounded-lg bg-gradient-to-r from-cyan-600 to-blue-600 hover:from-cyan-500 hover:to-blue-500 text-white font-semibold text-sm transition-all"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
            </svg>
            New Request
          </button>
        </div>

        <div v-if="loading && filteredRequests.length === 0" class="space-y-3">
          <div v-for="i in 3" :key="i" class="h-32 rounded-xl bg-gray-800/60 animate-pulse" />
        </div>
        <div v-else-if="filteredRequests.length === 0" class="flex flex-col items-center justify-center py-20 text-center">
          <div class="text-5xl mb-4">🚀</div>
          <p class="text-gray-400 font-medium">No deployment requests yet</p>
          <p class="text-gray-600 text-sm mt-1">Submit a request when your code is ready to deploy</p>
          <button @click="openCreateModal()" class="mt-4 px-4 py-2 rounded-lg bg-cyan-600/20 border border-cyan-600/40 text-cyan-400 text-sm font-semibold hover:bg-cyan-600/30 transition-colors">
            Create First Request
          </button>
        </div>
        <div v-else class="space-y-4 max-w-3xl">
          <DeploymentRequestCard
            v-for="req in filteredRequests"
            :key="req.id"
            :request="req"
            :is-chief-engineer="false"
          />
        </div>
      </div>
    </div>

    <!-- ═══════════════════════════════════════════════════════════════════════ -->
    <!-- Create Request Modal -->
    <!-- ═══════════════════════════════════════════════════════════════════════ -->
    <Teleport to="body">
      <div v-if="showCreateModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-sm">
        <div class="bg-gray-800 rounded-2xl border border-gray-700 w-full max-w-lg shadow-2xl">
          <div class="px-6 pt-6 pb-4 border-b border-gray-700">
            <h2 class="text-lg font-bold text-white">New Deployment Request</h2>
            <p class="text-sm text-gray-400 mt-0.5">Submit a branch for Chief Engineer review</p>
          </div>
          <form @submit.prevent="submitCreate" class="px-6 py-4 space-y-4">
            <div>
              <label class="block text-xs font-semibold text-gray-400 uppercase tracking-wide mb-1.5">Title *</label>
              <input v-model="form.title" placeholder="e.g. feat: Add user auth — v2.3.0" required
                class="w-full px-3 py-2.5 rounded-lg bg-gray-700 border border-gray-600 text-white placeholder-gray-500 text-sm focus:outline-none focus:border-cyan-500 focus:ring-1 focus:ring-cyan-500/50" />
            </div>
            <div class="grid grid-cols-2 gap-3">
              <div>
                <label class="block text-xs font-semibold text-gray-400 uppercase tracking-wide mb-1.5">Branch *</label>
                <input v-model="form.branch" placeholder="feature/auth-v2" required
                  class="w-full px-3 py-2.5 rounded-lg bg-gray-700 border border-gray-600 text-white placeholder-gray-500 text-sm font-mono focus:outline-none focus:border-cyan-500 focus:ring-1 focus:ring-cyan-500/50" />
              </div>
              <div>
                <label class="block text-xs font-semibold text-gray-400 uppercase tracking-wide mb-1.5">Environment</label>
                <select v-model="form.environment"
                  class="w-full px-3 py-2.5 rounded-lg bg-gray-700 border border-gray-600 text-white text-sm focus:outline-none focus:border-cyan-500 focus:ring-1 focus:ring-cyan-500/50">
                  <option value="STAGING">STAGING</option>
                  <option value="PRE-PROD">PRE-PROD</option>
                  <option value="PRODUCTION">PRODUCTION</option>
                </select>
              </div>
            </div>
            <div>
              <label class="block text-xs font-semibold text-gray-400 uppercase tracking-wide mb-1.5">PR / MR URL</label>
              <input v-model="form.pr_url" type="url" placeholder="https://github.com/org/repo/pull/42"
                class="w-full px-3 py-2.5 rounded-lg bg-gray-700 border border-gray-600 text-white placeholder-gray-500 text-sm focus:outline-none focus:border-cyan-500 focus:ring-1 focus:ring-cyan-500/50" />
            </div>
            <div>
              <label class="block text-xs font-semibold text-gray-400 uppercase tracking-wide mb-1.5">Task Reference</label>
              <input v-model="form.task_ref" placeholder="TASK-123 or sprint goal description"
                class="w-full px-3 py-2.5 rounded-lg bg-gray-700 border border-gray-600 text-white placeholder-gray-500 text-sm focus:outline-none focus:border-cyan-500 focus:ring-1 focus:ring-cyan-500/50" />
            </div>
            <!-- Assignee selector -->
            <div>
              <label class="block text-xs font-semibold text-gray-400 uppercase tracking-wide mb-1.5">Assign to Chief Engineer</label>
              <select v-model="form.reviewer_id"
                class="w-full px-3 py-2.5 rounded-lg bg-gray-700 border border-gray-600 text-white text-sm focus:outline-none focus:border-cyan-500 focus:ring-1 focus:ring-cyan-500/50">
                <option value="">— Unassigned (anyone can pick up) —</option>
                <option v-for="ce in chiefEngineers" :key="ce.id" :value="String(ce.id)">
                  {{ ce.display_name || ce.email }}
                </option>
              </select>
              <p v-if="chiefEngineers.length === 0" class="text-[10px] text-gray-500 mt-1">No Chief Engineers found</p>
            </div>
            <div>
              <label class="block text-xs font-semibold text-gray-400 uppercase tracking-wide mb-1.5">Release Notes</label>
              <textarea v-model="form.description" rows="3" placeholder="What changed? Any migration or env var updates required?"
                class="w-full px-3 py-2.5 rounded-lg bg-gray-700 border border-gray-600 text-white placeholder-gray-500 text-sm resize-none focus:outline-none focus:border-cyan-500 focus:ring-1 focus:ring-cyan-500/50" />
            </div>
            <!-- pre-prod warning -->
            <div v-if="form.environment === 'PRE-PROD'" class="flex items-start gap-2 px-3 py-2.5 rounded-lg bg-amber-500/10 border border-amber-500/30">
              <svg class="w-4 h-4 text-amber-400 mt-0.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
              </svg>
              <p class="text-xs text-amber-300">You are requesting a <strong>PRE-PROD</strong> deployment. This will be verified against the pre-production environment before production.</p>
            </div>
            <!-- production warning -->
            <div v-if="form.environment === 'PRODUCTION'" class="flex items-start gap-2 px-3 py-2.5 rounded-lg bg-red-500/10 border border-red-500/30">
              <svg class="w-4 h-4 text-red-400 mt-0.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
              </svg>
              <p class="text-xs text-red-300">You are requesting a <strong>PRODUCTION</strong> deployment. The Chief Engineer will review this carefully before approving.</p>
            </div>
            <div class="flex gap-3 pt-2">
              <button type="button" @click="showCreateModal = false" class="flex-1 px-4 py-2.5 rounded-lg border border-gray-600 text-gray-400 hover:text-white hover:border-gray-500 text-sm font-semibold transition-colors">
                Cancel
              </button>
              <button type="submit" :disabled="submitting" class="flex-1 px-4 py-2.5 rounded-lg bg-gradient-to-r from-cyan-600 to-blue-600 hover:from-cyan-500 hover:to-blue-500 disabled:opacity-50 text-white font-semibold text-sm transition-all">
                {{ submitting ? 'Submitting…' : 'Submit Request' }}
              </button>
            </div>
          </form>
        </div>
      </div>

      <!-- Approve Modal -->
      <div v-if="approveTarget" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-sm">
        <div class="bg-gray-800 rounded-2xl border border-green-700/40 w-full max-w-md shadow-2xl">
          <div class="px-6 pt-6 pb-4 border-b border-gray-700">
            <h2 class="text-lg font-bold text-green-400">Approve Deployment</h2>
            <p class="text-sm text-gray-400 mt-0.5 truncate">{{ approveTarget.title }}</p>
          </div>
          <div class="px-6 py-4 space-y-4">
            <div>
              <label class="block text-xs font-semibold text-gray-400 uppercase tracking-wide mb-1.5">Review Notes (optional)</label>
              <textarea v-model="reviewNotes" rows="3" placeholder="LGTM — branch is clean, tests pass, ready to merge."
                class="w-full px-3 py-2.5 rounded-lg bg-gray-700 border border-gray-600 text-white placeholder-gray-500 text-sm resize-none focus:outline-none focus:border-green-500 focus:ring-1 focus:ring-green-500/50" />
            </div>
            <div class="flex gap-3">
              <button @click="approveTarget = null; reviewNotes = ''" class="flex-1 px-4 py-2.5 rounded-lg border border-gray-600 text-gray-400 hover:text-white text-sm font-semibold transition-colors">Cancel</button>
              <button @click="submitApprove" :disabled="submitting" class="flex-1 px-4 py-2.5 rounded-lg bg-green-600 hover:bg-green-500 disabled:opacity-50 text-white font-semibold text-sm transition-colors">
                {{ submitting ? 'Approving…' : '✓ Approve' }}
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Reject Modal -->
      <div v-if="rejectTarget" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-sm">
        <div class="bg-gray-800 rounded-2xl border border-red-700/40 w-full max-w-md shadow-2xl">
          <div class="px-6 pt-6 pb-4 border-b border-gray-700">
            <h2 class="text-lg font-bold text-red-400">Reject Deployment</h2>
            <p class="text-sm text-gray-400 mt-0.5 truncate">{{ rejectTarget.title }}</p>
          </div>
          <div class="px-6 py-4 space-y-4">
            <div>
              <label class="block text-xs font-semibold text-gray-400 uppercase tracking-wide mb-1.5">Rejection Reason *</label>
              <textarea v-model="rejectReason" rows="3" placeholder="Describe what needs to be fixed before re-submitting…"
                class="w-full px-3 py-2.5 rounded-lg bg-gray-700 border border-red-600/50 text-white placeholder-gray-500 text-sm resize-none focus:outline-none focus:border-red-500 focus:ring-1 focus:ring-red-500/50" />
            </div>
            <div>
              <label class="block text-xs font-semibold text-gray-400 uppercase tracking-wide mb-1.5">Review Notes (optional)</label>
              <textarea v-model="reviewNotes" rows="2" placeholder="Additional context for the engineer…"
                class="w-full px-3 py-2.5 rounded-lg bg-gray-700 border border-gray-600 text-white placeholder-gray-500 text-sm resize-none focus:outline-none focus:border-red-500 focus:ring-1 focus:ring-red-500/50" />
            </div>
            <div class="flex gap-3">
              <button @click="rejectTarget = null; rejectReason = ''; reviewNotes = ''" class="flex-1 px-4 py-2.5 rounded-lg border border-gray-600 text-gray-400 hover:text-white text-sm font-semibold transition-colors">Cancel</button>
              <button @click="submitReject" :disabled="submitting || !rejectReason.trim()" class="flex-1 px-4 py-2.5 rounded-lg bg-red-600 hover:bg-red-500 disabled:opacity-50 text-white font-semibold text-sm transition-colors">
                {{ submitting ? 'Rejecting…' : '✗ Reject' }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { useDeploymentApi, type DeploymentRequest } from '~/core/modules/deployment/infrastructure/deployment-api'

definePageMeta({ layout: 'default' })

const { currentUser } = useAuth()
const { showSuccess, showError } = useNotification()
const deploymentApi = useDeploymentApi()

const isChiefEngineer = computed(() => currentUser.value?.role === 'CHIEF_ENGINEER')

// ─── State ────────────────────────────────────────────────────────────────────

const loading = ref(false)
const submitting = ref(false)
const requests = ref<DeploymentRequest[]>([])
const stats = ref<import('~/core/modules/deployment/infrastructure/deployment-api').DeploymentStats | null>(null)
const activeTab = ref<string>('all')

const showCreateModal = ref(false)
const approveTarget = ref<DeploymentRequest | null>(null)
const rejectTarget = ref<DeploymentRequest | null>(null)
const reviewNotes = ref('')
const rejectReason = ref('')

const form = reactive({
  title: '',
  branch: '',
  pr_url: '',
  environment: 'STAGING' as 'STAGING' | 'PRE-PROD' | 'PRODUCTION',
  task_ref: '',
  description: '',
  reviewer_id: '' as string, // '' = unassigned
})

// Chief Engineers list for assignee selector
const chiefEngineers = ref<import('~/core/modules/deployment/infrastructure/deployment-api').DeploymentUser[]>([])

async function openCreateModal() {
  showCreateModal.value = true
  if (chiefEngineers.value.length === 0) {
    chiefEngineers.value = await deploymentApi.fetchChiefEngineers()
  }
}

// ─── Tabs ─────────────────────────────────────────────────────────────────────

const tabs = computed(() => {
  const allCount = requests.value.length
  const countFor = (s: string) => requests.value.filter(r => r.status === s).length
  return [
    { value: 'all', label: 'All', count: allCount },
    { value: 'PENDING', label: 'Pending', count: countFor('PENDING') },
    { value: 'REVIEWING', label: 'In Review', count: countFor('REVIEWING') },
    { value: 'APPROVED', label: 'Approved', count: countFor('APPROVED') },
    { value: 'DEPLOYED', label: 'Deployed', count: countFor('DEPLOYED') },
    { value: 'REJECTED', label: 'Rejected', count: countFor('REJECTED') },
  ]
})

const filteredRequests = computed(() => {
  if (activeTab.value === 'all') return requests.value
  return requests.value.filter(r => r.status === activeTab.value)
})

// ─── Data fetching ────────────────────────────────────────────────────────────

async function refresh() {
  loading.value = true
  try {
    const [reqs, st] = await Promise.all([
      deploymentApi.listRequests(),
      deploymentApi.getStats(),
    ])
    requests.value = reqs
    stats.value = st
  } catch (e: any) {
    showError(e?.message ?? 'Failed to load deployment requests')
  } finally {
    loading.value = false
  }
}

onMounted(refresh)

// ─── Create ───────────────────────────────────────────────────────────────────

async function submitCreate() {
  submitting.value = true
  try {
    const result = await deploymentApi.createRequest({
      title: form.title,
      description: form.description,
      branch: form.branch,
      pr_url: form.pr_url || undefined,
      environment: form.environment,
      task_ref: form.task_ref || undefined,
      reviewer_id: form.reviewer_id ? Number(form.reviewer_id) : undefined,
    })
    requests.value.unshift(result)
    showSuccess('Deployment request submitted. Awaiting Chief Engineer review.')
    showCreateModal.value = false
    Object.assign(form, { title: '', branch: '', pr_url: '', environment: 'STAGING', task_ref: '', description: '', reviewer_id: '' })
    await refresh()
  } catch (e: any) {
    showError(e?.message ?? 'Failed to submit request')
  } finally {
    submitting.value = false
  }
}

// ─── Chief Engineer actions ───────────────────────────────────────────────────

async function pickForReview(req: DeploymentRequest) {
  try {
    const updated = await deploymentApi.pickForReview(req.id)
    updateInList(updated)
    showSuccess(`Started reviewing: ${req.title}`)
    stats.value = await deploymentApi.getStats()
  } catch (e: any) {
    showError(e?.message ?? 'Failed to pick up request')
  }
}

function openApprove(req: DeploymentRequest) {
  approveTarget.value = req
  reviewNotes.value = ''
}

async function submitApprove() {
  if (!approveTarget.value) return
  submitting.value = true
  try {
    const updated = await deploymentApi.approveRequest(approveTarget.value.id, { notes: reviewNotes.value })
    updateInList(updated)
    showSuccess(`Approved: ${approveTarget.value.title}`)
    approveTarget.value = null
    reviewNotes.value = ''
    stats.value = await deploymentApi.getStats()
  } catch (e: any) {
    showError(e?.message ?? 'Failed to approve')
  } finally {
    submitting.value = false
  }
}

function openReject(req: DeploymentRequest) {
  rejectTarget.value = req
  rejectReason.value = ''
  reviewNotes.value = ''
}

async function submitReject() {
  if (!rejectTarget.value || !rejectReason.value.trim()) return
  submitting.value = true
  try {
    const updated = await deploymentApi.rejectRequest(rejectTarget.value.id, {
      reason: rejectReason.value,
      notes: reviewNotes.value,
    })
    updateInList(updated)
    showSuccess(`Rejected: ${rejectTarget.value.title}`)
    rejectTarget.value = null
    rejectReason.value = ''
    reviewNotes.value = ''
    stats.value = await deploymentApi.getStats()
  } catch (e: any) {
    showError(e?.message ?? 'Failed to reject')
  } finally {
    submitting.value = false
  }
}

async function confirmDeploy(req: DeploymentRequest) {
  try {
    const updated = await deploymentApi.markDeployed(req.id, { notes: 'Deployed successfully' })
    updateInList(updated)
    showSuccess(`🚀 Deployed: ${req.title}`)
    stats.value = await deploymentApi.getStats()
  } catch (e: any) {
    showError(e?.message ?? 'Failed to mark as deployed')
  }
}

function updateInList(updated: DeploymentRequest) {
  const idx = requests.value.findIndex(r => r.id === updated.id)
  if (idx !== -1) requests.value[idx] = updated
}
</script>

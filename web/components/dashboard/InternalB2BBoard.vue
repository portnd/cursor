<template>
  <section>
    <div class="flex items-center justify-between mb-4">
      <h2 class="section-label">Internal B2B Trading Floor</h2>
      <button
        v-if="canCreate"
        @click="showNewRequestModal = true"
        class="inline-flex items-center gap-1.5 rounded-lg bg-blue-100 dark:bg-blue-600/20 border border-blue-300 dark:border-blue-500/30 hover:bg-blue-100 dark:bg-blue-600/30 px-3 py-1.5 text-xs font-semibold text-blue-300 transition-colors"
      >
        <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
        </svg>
        New Request
      </button>
    </div>

    <div class="rounded-2xl border border-gray-700 bg-gray-800/60 overflow-hidden shadow-lg">
      <!-- Tabs -->
      <div class="flex border-b border-gray-700">
        <button
          v-for="tab in tabs"
          :key="tab.id"
          @click="switchTab(tab.id)"
          class="flex items-center gap-2 px-5 py-3.5 text-xs font-semibold uppercase tracking-wider transition-colors relative"
          :class="activeTab === tab.id
            ? 'text-white after:absolute after:bottom-0 after:left-0 after:right-0 after:h-0.5 after:bg-blue-500'
            : 'text-gray-500 hover:text-gray-300'"
        >
          <span class="w-2 h-2 rounded-full" :class="tab.dotClass"/>
          {{ tab.label }}
          <span
            class="ml-1 px-1.5 py-0.5 rounded-full text-xs font-bold"
            :class="activeTab === tab.id ? tab.badgeActive : tab.badgeInactive"
          >{{ tab.id === 'inbound' ? inboundRequests.length : outboundRequests.length }}</span>
        </button>
      </div>

      <!-- Loading -->
      <div v-if="isLoading" class="flex items-center justify-center py-16">
        <svg class="h-6 w-6 animate-spin text-blue-400" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/>
        </svg>
      </div>

      <div v-else class="p-4">

        <!-- INBOUND TAB -->
        <div v-if="activeTab === 'inbound'">
          <div v-if="inboundRequests.length === 0" class="flex flex-col items-center justify-center py-12 text-center">
            <svg class="h-8 w-8 text-gray-600 mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4"/>
            </svg>
            <p class="text-sm text-gray-400 font-medium">No inbound requests</p>
            <p class="text-xs text-gray-600 mt-0.5">Cross-team assignments will appear here</p>
          </div>
          <ul v-else class="space-y-2">
            <li
              v-for="req in inboundRequests"
              :key="req.id"
              class="rounded-xl border border-blue-500/15 bg-gray-900/50 p-4 transition-all hover:border-blue-500/30"
            >
              <div class="flex items-start justify-between gap-3">
                <div class="min-w-0 flex-1">
                  <p class="font-semibold text-white text-sm truncate">{{ req.title }}</p>
                  <p v-if="req.description" class="text-xs text-gray-500 mt-0.5 line-clamp-2">{{ req.description }}</p>
                  <div class="flex items-center gap-3 mt-2">
                    <span class="inline-flex items-center gap-1 text-xs text-gray-400">
                      <svg class="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
                      </svg>
                      Requested:
                      <span class="font-semibold text-blue-300">{{ formatMinutes(req.estimated_minutes) }}</span>
                    </span>
                    <span v-if="req.requester_team_name" class="text-xs text-gray-500">
                      from <span class="text-gray-300">{{ req.requester_team_name }}</span>
                    </span>
                    <span
                      class="px-2 py-0.5 rounded-full text-xs font-medium"
                      :class="statusClass(req.status)"
                    >{{ req.status }}</span>
                  </div>
                  <div v-if="req.status === 'COUNTER_OFFERED'" class="mt-1.5 flex items-center gap-2">
                    <span class="text-xs text-amber-400">Counter proposed:
                      <span class="font-semibold">{{ formatMinutes(req.proposed_minutes) }}</span>
                    </span>
                    <span v-if="req.negotiation_reason" class="text-xs text-gray-500 italic">"{{ req.negotiation_reason }}"</span>
                  </div>
                </div>
                <!-- Actions for target team Product Owner -->
                <div v-if="req.status === 'PENDING' || req.status === 'COUNTER_OFFERED'" class="flex flex-col gap-2 flex-shrink-0">
                  <button
                    @click="acceptInbound(req)"
                    :disabled="actionLoading === req.id"
                    class="inline-flex items-center gap-1.5 rounded-lg bg-emerald-100 dark:bg-emerald-600 hover:bg-emerald-100 dark:bg-emerald-500 px-3 py-1.5 text-xs font-semibold text-gray-900 dark:text-white transition-colors disabled:opacity-50"
                  >
                    <svg v-if="actionLoading === req.id" class="h-3 w-3 animate-spin" fill="none" viewBox="0 0 24 24">
                      <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
                      <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/>
                    </svg>
                    Accept
                  </button>
                  <button
                    @click="openCounterModal(req)"
                    :disabled="actionLoading === req.id"
                    class="inline-flex items-center gap-1.5 rounded-lg border border-amber-300 dark:border-amber-500/40 bg-amber-100 dark:bg-amber-500/10 hover:bg-amber-100 dark:bg-amber-500/20 px-3 py-1.5 text-xs font-semibold text-amber-300 transition-colors disabled:opacity-50"
                  >
                    Counter
                  </button>
                  <button
                    @click="rejectRequest(req)"
                    :disabled="actionLoading === req.id"
                    class="inline-flex items-center gap-1.5 rounded-lg border border-red-300 dark:border-red-500/30 bg-red-100 dark:bg-red-500/10 hover:bg-red-100 dark:bg-red-500/20 px-3 py-1.5 text-xs font-semibold text-red-300 transition-colors disabled:opacity-50"
                  >
                    Reject
                  </button>
                </div>
                <!-- Accepted: show link to created task -->
                <div v-else-if="req.status === 'ACCEPTED'" class="flex-shrink-0">
                  <span class="px-2 py-1 rounded-lg bg-emerald-900/40 text-emerald-400 text-xs font-semibold">Task Created</span>
                </div>
              </div>
            </li>
          </ul>
        </div>

        <!-- OUTBOUND TAB -->
        <div v-if="activeTab === 'outbound'">
          <div v-if="outboundRequests.length === 0" class="flex flex-col items-center justify-center py-12 text-center">
            <svg class="h-8 w-8 text-gray-600 mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M8 4H6a2 2 0 00-2 2v12a2 2 0 002 2h12a2 2 0 002-2V6a2 2 0 00-2-2h-2m-4-1v8m0 0l3-3m-3 3L9 8m-5 5h2.586a1 1 0 01.707.293l2.414 2.414a1 1 0 00.707.293h3.172a1 1 0 00.707-.293l2.414-2.414A1 1 0 0014.414 13H17"/>
            </svg>
            <p class="text-sm text-gray-400 font-medium">No outbound tasks</p>
            <p class="text-xs text-gray-600 mt-0.5">Tasks outsourced to other teams appear here</p>
          </div>
          <ul v-else class="space-y-2">
            <li
              v-for="req in outboundRequests"
              :key="req.id"
              class="rounded-xl border border-purple-500/15 bg-gray-900/50 p-4 transition-all hover:border-purple-500/30"
            >
              <div class="flex items-start justify-between gap-3">
                <div class="min-w-0 flex-1">
                  <p class="font-semibold text-white text-sm truncate">{{ req.title }}</p>
                  <div class="flex items-center gap-3 mt-2">
                    <span class="text-xs text-gray-400">
                      Est: <span class="font-semibold text-purple-300">{{ formatMinutes(req.estimated_minutes) }}</span>
                    </span>
                    <span v-if="req.target_team_name" class="text-xs text-gray-500">
                      to <span class="text-gray-300">{{ req.target_team_name }}</span>
                    </span>
                    <span class="text-xs px-2 py-0.5 rounded-full" :class="statusClass(req.status)">{{ req.status }}</span>
                  </div>
                  <!-- Counter-offer from target team -->
                  <div v-if="req.status === 'COUNTER_OFFERED'" class="mt-2 p-2.5 rounded-lg bg-amber-900/20 border border-amber-500/30">
                    <p class="text-xs text-amber-300 font-semibold mb-0.5">Counter-offer received</p>
                    <p class="text-xs text-amber-400">
                      Proposed: <span class="font-bold">{{ formatMinutes(req.proposed_minutes) }}</span>
                    </p>
                    <p v-if="req.negotiation_reason" class="text-xs text-gray-400 mt-0.5 italic">"{{ req.negotiation_reason }}"</p>
                    <div class="flex gap-2 mt-2">
                      <button
                        @click="acceptOutbound(req)"
                        :disabled="actionLoading === req.id"
                        class="inline-flex items-center gap-1 rounded-lg bg-emerald-100 dark:bg-emerald-600 hover:bg-emerald-100 dark:bg-emerald-500 px-2.5 py-1 text-xs font-semibold text-gray-900 dark:text-white transition-colors disabled:opacity-50"
                      >Accept Counter</button>
                      <button
                        @click="rejectRequest(req)"
                        :disabled="actionLoading === req.id"
                        class="inline-flex items-center gap-1 rounded-lg border border-red-300 dark:border-red-500/30 bg-red-100 dark:bg-red-500/10 hover:bg-red-100 dark:bg-red-500/20 px-2.5 py-1 text-xs font-semibold text-red-300 transition-colors disabled:opacity-50"
                      >Reject</button>
                    </div>
                  </div>
                </div>
              </div>
            </li>
          </ul>
        </div>

      </div>
    </div>

    <!-- Counter-Offer Modal -->
    <Teleport to="body">
      <Transition name="modal">
        <div
          v-if="showCounterModal"
          class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/70 backdrop-blur-sm"
          @click.self="closeCounterModal"
        >
          <div class="relative w-full max-w-md rounded-2xl border border-amber-500/30 bg-gray-800 shadow-2xl" @click.stop>
            <div class="border-b border-gray-700 px-6 py-4">
              <div class="flex items-center gap-3">
                <div class="w-8 h-8 rounded-lg bg-amber-500/15 border border-amber-500/30 flex items-center justify-center">
                  <svg class="w-4 h-4 text-amber-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"/>
                  </svg>
                </div>
                <div>
                  <h2 class="text-sm font-bold text-white">Counter-Offer</h2>
                  <p class="text-xs text-gray-500 line-clamp-1 mt-0.5">{{ counterTarget?.title }}</p>
                </div>
              </div>
            </div>
            <form class="p-6 space-y-4" @submit.prevent="submitCounter">
              <div>
                <p class="text-xs text-gray-500 mb-1">
                  Requested time: <span class="font-semibold text-blue-300">{{ formatMinutes(counterTarget?.estimated_minutes ?? 0) }}</span>
                </p>
              </div>
              <div>
                <label class="block text-xs font-medium text-gray-400 mb-1.5">Your Counter-Offer (minutes)</label>
                <input
                  v-model.number="counterMinutes"
                  type="number"
                  min="1"
                  placeholder="e.g. 120"
                  class="w-full rounded-lg border border-gray-600 bg-gray-900/50 px-3 py-2.5 text-sm text-white placeholder-gray-500 focus:border-amber-500 focus:outline-none focus:ring-1 focus:ring-amber-500/50 transition-colors"
                />
              </div>
              <div>
                <label class="block text-xs font-medium text-gray-400 mb-1.5">Reason / Notes</label>
                <textarea
                  v-model="counterReason"
                  rows="3"
                  placeholder="Explain the complexity, dependencies…"
                  class="w-full rounded-lg border border-gray-600 bg-gray-900/50 px-3 py-2.5 text-sm text-white placeholder-gray-500 focus:border-amber-500 focus:outline-none focus:ring-1 focus:ring-amber-500/50 transition-colors resize-none"
                />
              </div>
              <div class="flex gap-3">
                <button type="button" @click="closeCounterModal" class="btn-ghost flex-1">Cancel</button>
                <button type="submit" :disabled="isCountering || !counterMinutes" class="btn-amber flex-1">
                  <svg v-if="isCountering" class="mr-2 h-3.5 w-3.5 animate-spin" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/>
                  </svg>
                  Send Counter-Offer
                </button>
              </div>
            </form>
          </div>
        </div>
      </Transition>
    </Teleport>

    <!-- New Outsource Request Modal -->
    <OutsourceRequestModal
      v-model="showNewRequestModal"
      @created="onRequestCreated"
    />
  </section>
</template>

<script setup lang="ts">
import { useB2BApi, type B2BRequest } from '~/core/modules/b2b/infrastructure/b2b-api'
import OutsourceRequestModal from '~/components/b2b/OutsourceRequestModal.vue'

const { getRequests, counterOffer, rejectRequest: apiReject, acceptRequest } = useB2BApi()
const { currentUser } = useAuth()
const { showSuccess, showError } = useNotification()

const isLoading = ref(true)
const inboundRequests = ref<B2BRequest[]>([])
const outboundRequests = ref<B2BRequest[]>([])
const activeTab = ref<'inbound' | 'outbound'>('inbound')
const actionLoading = ref<string | null>(null)
const showNewRequestModal = ref(false)

// Counter-offer modal state
const showCounterModal = ref(false)
const counterTarget = ref<B2BRequest | null>(null)
const counterMinutes = ref<number | null>(null)
const counterReason = ref('')
const isCountering = ref(false)

const tabs = [
  {
    id: 'inbound',
    label: 'Inbound (Revenue)',
    dotClass: 'bg-blue-400',
    badgeActive: 'bg-blue-500/20 text-blue-300',
    badgeInactive: 'bg-gray-700 text-gray-500',
  },
  {
    id: 'outbound',
    label: 'Outbound (Cost)',
    dotClass: 'bg-purple-400',
    badgeActive: 'bg-purple-500/20 text-purple-300',
    badgeInactive: 'bg-gray-700 text-gray-500',
  },
]

const userRole = computed(() => (currentUser.value as any)?.role?.toUpperCase() ?? '')
const canCreate = computed(() =>
  ['CEO', 'PRODUCT_OWNER', 'PM', 'MANAGER'].includes(userRole.value))

const formatMinutes = (mins: number) => {
  if (!mins || mins <= 0) return '—'
  const h = Math.floor(mins / 60)
  const m = mins % 60
  if (h === 0) return `${m}m`
  if (m === 0) return `${h}h`
  return `${h}h ${m}m`
}

const statusClass = (status: string) => {
  const map: Record<string, string> = {
    PENDING: 'bg-amber-500/15 text-amber-300 border border-amber-500/30',
    COUNTER_OFFERED: 'bg-blue-500/15 text-blue-300 border border-blue-500/30',
    ACCEPTED: 'bg-emerald-500/15 text-emerald-300 border border-emerald-500/30',
    REJECTED: 'bg-red-500/15 text-red-300 border border-red-500/30',
  }
  return map[status] || 'bg-gray-700 text-gray-400'
}

const fetchAll = async () => {
  isLoading.value = true
  try {
    const [inbRes, outbRes] = await Promise.all([
      getRequests('inbound'),
      getRequests('outbound'),
    ])
    inboundRequests.value = inbRes?.data ?? []
    outboundRequests.value = outbRes?.data ?? []
  } catch {
    // silent
  } finally {
    isLoading.value = false
  }
}

const switchTab = (tab: 'inbound' | 'outbound') => {
  activeTab.value = tab
}

// Inbound: target team accepts → task is created
const acceptInbound = async (req: B2BRequest) => {
  actionLoading.value = req.id
  try {
    await acceptRequest(req.id)
    showSuccess('Request accepted! Task created in your project.', 'Done')
    await fetchAll()
  } catch (err: any) {
    showError(err?.data?.message ?? err?.message ?? 'Failed to accept')
  } finally {
    actionLoading.value = null
  }
}

// Outbound with counter-offer: requester accepts the counter-offer → task is created
const acceptOutbound = async (req: B2BRequest) => {
  actionLoading.value = req.id
  try {
    await acceptRequest(req.id)
    showSuccess('Counter-offer accepted! Task created in their project.', 'Done')
    await fetchAll()
  } catch (err: any) {
    showError(err?.data?.message ?? err?.message ?? 'Failed to accept')
  } finally {
    actionLoading.value = null
  }
}

const rejectRequest = async (req: B2BRequest) => {
  actionLoading.value = req.id
  try {
    await apiReject(req.id)
    showSuccess('Request rejected.', 'Done')
    await fetchAll()
  } catch (err: any) {
    showError(err?.data?.message ?? err?.message ?? 'Failed to reject')
  } finally {
    actionLoading.value = null
  }
}

const openCounterModal = (req: B2BRequest) => {
  counterTarget.value = req
  counterMinutes.value = req.estimated_minutes || null
  counterReason.value = ''
  showCounterModal.value = true
}

const closeCounterModal = () => {
  showCounterModal.value = false
  counterTarget.value = null
  counterMinutes.value = null
  counterReason.value = ''
}

const submitCounter = async () => {
  if (!counterTarget.value || !counterMinutes.value) return
  isCountering.value = true
  try {
    await counterOffer(counterTarget.value.id, counterMinutes.value, counterReason.value)
    showSuccess('Counter-offer sent.', 'Done')
    closeCounterModal()
    await fetchAll()
  } catch (err: any) {
    showError(err?.data?.message ?? err?.message ?? 'Failed to send counter-offer')
  } finally {
    isCountering.value = false
  }
}

const onRequestCreated = async () => {
  showSuccess('Outsource request sent!', 'Sent')
  activeTab.value = 'outbound'
  await fetchAll()
}

onMounted(() => fetchAll())
</script>

<style scoped>
.section-label {
  @apply text-xs font-semibold uppercase tracking-widest text-gray-500;
}
.btn-ghost {
  @apply inline-flex items-center justify-center rounded-lg border border-gray-600 bg-transparent
         px-4 py-2.5 text-sm font-semibold text-gray-300 transition-colors
         hover:bg-gray-700 hover:text-white disabled:cursor-not-allowed disabled:opacity-50;
}
.btn-amber {
  @apply inline-flex items-center justify-center rounded-lg bg-gradient-to-r from-amber-600 to-orange-600
         px-4 py-2.5 text-sm font-semibold text-white shadow-lg transition-all
         hover:from-amber-500 hover:to-orange-500 disabled:cursor-not-allowed disabled:opacity-50;
}
.modal-enter-active, .modal-leave-active { transition: opacity 0.2s ease; }
.modal-enter-from, .modal-leave-to { opacity: 0; }
</style>

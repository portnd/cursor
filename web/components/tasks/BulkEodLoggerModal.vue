<template>
  <Teleport to="body">
    <Transition name="fade">
      <div
        v-if="show"
        class="fixed inset-0 z-[80] flex items-end sm:items-center justify-center p-0 sm:p-4"
        @keydown.escape="$emit('close')"
      >
        <!-- Backdrop -->
        <div class="fixed inset-0 bg-black/80 backdrop-blur-sm" @click="$emit('close')" />

        <!-- Panel -->
        <div
          class="relative w-full sm:max-w-2xl min-h-0 max-h-[95dvh] sm:max-h-[90vh] bg-gray-900 border border-gray-700 sm:rounded-2xl rounded-t-2xl shadow-2xl flex flex-col overflow-hidden"
          @click.stop
        >
          <!-- Header -->
          <div class="flex items-center justify-between gap-3 px-5 pt-5 pb-4 border-b border-gray-700 shrink-0">
            <div>
              <h2 class="text-white font-bold text-lg flex items-center gap-2">
                <span class="text-2xl">📋</span>
                EOD Batch Log
              </h2>
              <p class="text-xs text-gray-400 mt-0.5">Log หลาย tasks ในครั้งเดียว — สำหรับสรุปท้ายวัน</p>
            </div>
            <button type="button" @click="$emit('close')" class="p-2 rounded-lg text-gray-400 hover:text-white hover:bg-gray-800 transition-colors">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
              </svg>
            </button>
          </div>

          <!-- Daily quota indicator -->
          <div class="px-5 py-3 bg-gray-800/50 border-b border-gray-700/50 shrink-0">
            <div class="flex items-center justify-between mb-1.5">
              <span class="text-xs text-gray-400">วันนี้ log แล้ว</span>
              <span class="text-xs font-bold" :class="dailyTotalH >= 8 ? 'text-green-400' : 'text-blue-400'">
                {{ dailyTotalH.toFixed(1) }}h / 8h
                <span class="ml-1 font-normal text-gray-500">(+ {{ pendingMinutes }}m จาก entries นี้)</span>
              </span>
            </div>
            <div class="h-1.5 bg-gray-700 rounded-full overflow-hidden">
              <div
                class="h-full rounded-full transition-all duration-500"
                :class="combinedPercent >= 100 ? 'bg-green-500' : combinedPercent >= 60 ? 'bg-blue-500' : 'bg-blue-600'"
                :style="{ width: Math.min(combinedPercent, 100) + '%' }"
              />
            </div>
          </div>

          <!-- Entries list -->
          <div class="flex-1 min-h-0 overflow-y-auto px-5 py-4 space-y-3">
            <!-- Row for each entry -->
            <div
              v-for="(entry, idx) in entries"
              :key="idx"
              class="group relative bg-gray-800/60 border rounded-xl p-3 transition-all"
              :class="getRowClass(entry, idx)"
            >
              <!-- Result overlay on success/fail -->
              <div v-if="results[idx]" class="absolute inset-0 flex items-center justify-center rounded-xl z-10"
                :class="results[idx].success ? 'bg-green-900/70' : 'bg-red-900/70'">
                <div class="text-center">
                  <div class="text-2xl mb-1">{{ results[idx].success ? '✅' : '❌' }}</div>
                  <div class="text-xs font-medium" :class="results[idx].success ? 'text-green-300' : 'text-red-300'">
                    {{ results[idx].success ? 'Saved' : results[idx].error }}
                  </div>
                </div>
              </div>

              <div class="flex items-start gap-3">
                <!-- Row number -->
                <div class="w-6 h-6 rounded-full bg-gray-700 flex items-center justify-center text-[10px] text-gray-400 font-bold shrink-0 mt-0.5">
                  {{ idx + 1 }}
                </div>

                <div class="flex-1 min-w-0 space-y-2">
                  <!-- Task picker -->
                  <div>
                    <button
                      v-if="entry.task_id"
                      type="button"
                      class="w-full flex items-center gap-1.5 px-2 py-1.5 rounded-lg bg-indigo-900/40 border border-indigo-700/40 text-xs text-indigo-300 hover:bg-indigo-900/60 transition-colors"
                      @click="entry.task_id = ''; entry.search = ''"
                    >
                      <span class="w-1.5 h-1.5 rounded-full shrink-0" :style="{ background: entry.project_color || '#6366f1' }" />
                      <span class="font-mono text-[9px] text-indigo-400 shrink-0">{{ entry.task_code }}</span>
                      <span class="truncate flex-1 text-left">{{ entry.task_title }}</span>
                      <span class="text-gray-500 shrink-0">✕</span>
                    </button>
                    <template v-if="!entry.task_id">
                      <input
                        v-model="entry.search"
                        type="text"
                        placeholder="🔍 ค้นหา task..."
                        class="w-full bg-gray-700/60 border border-gray-600 rounded-lg px-3 py-2 text-sm text-white placeholder-gray-500 focus:outline-none focus:border-indigo-500 transition-colors"
                      />
                      <div class="mt-1 bg-gray-800 border border-gray-700 rounded-xl overflow-hidden max-h-36 overflow-y-auto">
                        <div v-if="tasksLoading" class="flex items-center gap-2 px-3 py-2.5 text-xs text-gray-500">
                          <svg class="w-3.5 h-3.5 animate-spin shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                          </svg>
                          กำลังโหลด tasks...
                        </div>
                        <template v-else-if="filteredTasks(entry.search).length > 0">
                          <button
                            v-for="task in filteredTasks(entry.search)"
                            :key="task.id"
                            type="button"
                            class="w-full flex items-center gap-2.5 px-3 py-2 hover:bg-gray-700/60 text-left transition-colors border-b border-gray-700/40 last:border-0"
                            @click="selectTask(entry, task)"
                          >
                            <span class="w-2 h-2 rounded-full shrink-0" :style="{ background: task.project_color || '#6366f1' }" />
                            <span class="text-[10px] text-gray-500 font-mono shrink-0">{{ task.code }}</span>
                            <span class="text-sm text-white truncate flex-1">{{ task.title }}</span>
                            <span v-if="task.assigned_to_display_name || task.assigned_to_email" class="text-[10px] text-indigo-400/80 shrink-0 hidden sm:block">
                              {{ task.assigned_to_display_name || task.assigned_to_email }}
                            </span>
                            <span class="text-[10px] text-gray-500 shrink-0">{{ task.project_name }}</span>
                          </button>
                        </template>
                        <div v-else class="px-3 py-2.5 text-xs text-gray-500 text-center">
                          {{ entry.search ? `ไม่พบ "${entry.search}"` : 'ไม่มี task ใน sprint' }}
                        </div>
                      </div>
                    </template>
                  </div>

                  <!-- Minutes + Work type row -->
                  <div class="flex items-center gap-2 flex-wrap">
                    <!-- Time presets -->
                    <div class="flex gap-1">
                      <button
                        v-for="preset in [15, 30, 60, 120]"
                        :key="preset"
                        type="button"
                        @click="entry.minutes = preset"
                        class="text-[10px] px-2 py-1 rounded-md border transition-colors"
                        :class="entry.minutes === preset
                          ? 'bg-indigo-600 border-indigo-500 text-white'
                          : 'bg-gray-700 border-gray-600 text-gray-400 hover:text-white'"
                      >+{{ preset >= 60 ? (preset / 60) + 'h' : preset + 'm' }}</button>
                    </div>
                    <!-- Minutes input -->
                    <input
                      v-model.number="entry.minutes"
                      type="number"
                      min="1"
                      max="960"
                      placeholder="นาที"
                      class="w-20 bg-gray-700/60 border border-gray-600 rounded-lg px-2 py-1.5 text-sm text-white placeholder-gray-500 focus:outline-none focus:border-indigo-500 text-center"
                    />
                    <span class="text-xs text-gray-500">
                      {{ entry.minutes >= 60 ? (entry.minutes / 60).toFixed(1) + 'h' : '' }}
                    </span>

                    <!-- Work type -->
                    <div class="flex gap-1 ml-auto">
                      <button
                        v-for="wt in WORK_TYPES"
                        :key="wt.value"
                        type="button"
                        @click="entry.work_type = wt.value"
                        class="text-[10px] px-2 py-1 rounded-md border transition-all"
                        :class="entry.work_type === wt.value
                          ? workTypeActiveClass(wt.value)
                          : 'bg-gray-700/50 border-gray-600 text-gray-500 hover:text-gray-300'"
                      >{{ wt.emoji }} {{ wt.label }}</button>
                    </div>
                  </div>

                  <!-- Description -->
                  <input
                    v-model="entry.description"
                    type="text"
                    placeholder="คำอธิบาย (optional)"
                    class="w-full bg-gray-700/40 border border-gray-700 rounded-lg px-3 py-1.5 text-xs text-gray-300 placeholder-gray-600 focus:outline-none focus:border-gray-500 transition-colors"
                  />
                </div>

                <!-- Remove row button -->
                <button
                  type="button"
                  @click="removeEntry(idx)"
                  class="shrink-0 p-1.5 rounded-lg text-gray-600 hover:text-red-400 hover:bg-red-900/20 transition-colors opacity-0 group-hover:opacity-100"
                  :disabled="entries.length <= 1"
                >
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
                  </svg>
                </button>
              </div>
            </div>

            <!-- Add row button -->
            <button
              type="button"
              @click="addEntry"
              :disabled="entries.length >= 20 || submitted"
              class="w-full py-2.5 border border-dashed border-gray-700 rounded-xl text-sm text-gray-500 hover:text-gray-300 hover:border-gray-600 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
            >
              + เพิ่ม task ({{ entries.length }}/20)
            </button>
          </div>

          <!-- Footer -->
          <div class="px-5 py-4 border-t border-gray-700 shrink-0 flex items-center gap-3">
            <div class="flex-1 text-xs text-gray-500">
              <span v-if="submitted">
                ✅ {{ successCount }}/{{ entries.length }} entries บันทึกสำเร็จ
              </span>
              <span v-else>
                {{ validEntries }} entries พร้อม submit
              </span>
            </div>
            <button
              type="button"
              @click="$emit('close')"
              class="px-4 py-2 text-sm text-gray-400 hover:text-white hover:bg-gray-800 rounded-lg transition-colors"
            >
              {{ submitted ? 'ปิด' : 'ยกเลิก' }}
            </button>
            <button
              v-if="!submitted"
              type="button"
              @click="submitAll"
              :disabled="loading || validEntries === 0"
              class="px-5 py-2 text-sm font-semibold rounded-lg transition-all disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
              :class="loading ? 'bg-indigo-700 text-white' : 'bg-indigo-600 hover:bg-indigo-500 text-white'"
            >
              <svg v-if="loading" class="animate-spin w-4 h-4" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8H4z"/>
              </svg>
              <span>{{ loading ? 'กำลัง Log...' : `Log ${validEntries} entries` }}</span>
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { useTasksApi, WORK_TYPES } from '~/core/modules/tasks/infrastructure/tasks-api'
import type { GlobalActiveTask, BulkLogEntry } from '~/core/modules/tasks/infrastructure/tasks-api'

const props = defineProps<{ show: boolean }>()
const emit = defineEmits<{ close: [], done: [] }>()

const { bulkLogTime, getTeamActiveTasks, getMyDailyTimeLogs } = useTasksApi()

interface EntryRow {
  task_id: string
  task_title: string
  task_code: string
  project_color: string
  minutes: number
  description: string
  work_type: string
  search: string
}

function freshEntry(): EntryRow {
  return {
    task_id: '', task_title: '', task_code: '', project_color: '',
    minutes: 30, description: '', work_type: 'DEV',
    search: '',
  }
}

const entries = ref<EntryRow[]>([freshEntry()])
const results = ref<Record<number, { success: boolean; error?: string }>>({})
const loading = ref(false)
const submitted = ref(false)
const tasks = ref<GlobalActiveTask[]>([])
const dailyMinutes = ref(0)
const tasksLoading = ref(false)

async function loadData() {
  tasksLoading.value = true
  try {
    const [taskRes, dailyRes] = await Promise.all([getTeamActiveTasks(), getMyDailyTimeLogs()])
    tasks.value = taskRes
    dailyMinutes.value = dailyRes.total_minutes
  } catch { /* non-critical */ }
  finally { tasksLoading.value = false }
}

watch(() => props.show, (v) => {
  if (v) {
    entries.value = [freshEntry()]
    results.value = {}
    submitted.value = false
    loadData()
  }
}, { immediate: true })

const dailyTotalH = computed(() => dailyMinutes.value / 60)
const pendingMinutes = computed(() => entries.value.reduce((s, e) => s + (e.minutes || 0), 0))
const combinedPercent = computed(() => ((dailyMinutes.value + pendingMinutes.value) / 480) * 100)

const validEntries = computed(() =>
  entries.value.filter(e => e.task_id && e.minutes > 0).length
)

const successCount = computed(() =>
  Object.values(results.value).filter(r => r.success).length
)

function filteredTasks(search: string): GlobalActiveTask[] {
  if (!search) return tasks.value.slice(0, 10)
  const q = search.toLowerCase()
  return tasks.value
    .filter(t =>
      t.title.toLowerCase().includes(q) ||
      t.code?.toLowerCase().includes(q) ||
      t.project_name?.toLowerCase().includes(q) ||
      t.assigned_to_display_name?.toLowerCase().includes(q) ||
      t.assigned_to_email?.toLowerCase().includes(q)
    )
    .slice(0, 10)
}

function selectTask(entry: EntryRow, task: GlobalActiveTask) {
  entry.task_id = task.id
  entry.task_title = task.title
  entry.task_code = task.code || ''
  entry.project_color = task.project_color || ''
  entry.search = ''
}

function addEntry() { entries.value.push(freshEntry()) }

function removeEntry(idx: number) {
  if (entries.value.length > 1) entries.value.splice(idx, 1)
}

async function submitAll() {
  const payload: BulkLogEntry[] = entries.value
    .filter(e => e.task_id && e.minutes > 0)
    .map(e => ({
      task_id: e.task_id,
      minutes: e.minutes,
      description: e.description || undefined,
      work_type: e.work_type || 'DEV',
    }))
  if (!payload.length) return

  loading.value = true
  try {
    const res = await bulkLogTime(payload)
    // Map results back to entry indices
    let payloadIdx = 0
    for (let i = 0; i < entries.value.length; i++) {
      const e = entries.value[i]
      if (e.task_id && e.minutes > 0) {
        const r = res.results[payloadIdx++]
        if (r) results.value[i] = { success: r.success, error: r.error }
      }
    }
    submitted.value = true
    emit('done')
  } catch (err: any) {
    alert(err?.message || 'เกิดข้อผิดพลาด')
  } finally {
    loading.value = false
  }
}

function getRowClass(entry: EntryRow, idx: number) {
  if (results.value[idx]) {
    return results.value[idx].success ? 'border-green-700/50' : 'border-red-700/50'
  }
  return entry.task_id ? 'border-gray-700' : 'border-gray-700/40'
}

function workTypeActiveClass(wt: string): string {
  const map: Record<string, string> = {
    DEV: 'bg-blue-600 border-blue-500 text-white',
    REVIEW: 'bg-cyan-600 border-cyan-500 text-white',
    TESTING: 'bg-green-600 border-green-500 text-white',
    MEETING: 'bg-orange-600 border-orange-500 text-white',
    RESEARCH: 'bg-purple-600 border-purple-500 text-white',
    OTHER: 'bg-gray-600 border-gray-500 text-white',
  }
  return map[wt] || 'bg-indigo-600 border-indigo-500 text-white'
}
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>

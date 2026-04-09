<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-lg font-bold text-white flex items-center gap-2">
          <span class="text-2xl">🗄</span> Backup & Restore
        </h2>
        <p class="text-sm text-gray-400 mt-1">
          สร้าง snapshot ของข้อมูลโครงการ (tasks, sprints, milestones, epics) เพื่อป้องกันข้อมูลสูญหาย
        </p>
      </div>
      <button
        @click="openCreateModal"
        class="flex items-center gap-2 px-4 py-2 bg-gradient-to-r from-purple-100 dark:from-purple-600 to-pink-100 dark:to-pink-600 text-gray-900 dark:text-white text-sm font-semibold rounded-lg hover:opacity-90 transition-opacity"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
        สร้าง Backup
      </button>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="flex items-center justify-center py-20">
      <div class="w-8 h-8 border-2 border-purple-500 border-t-transparent rounded-full animate-spin" />
    </div>

    <!-- Empty state -->
    <div v-else-if="backups.length === 0" class="text-center py-20">
      <div class="text-5xl mb-4">🗄</div>
      <h3 class="text-lg font-semibold text-gray-300 mb-2">ยังไม่มี Backup</h3>
      <p class="text-sm text-gray-500 mb-6">สร้าง backup แรกเพื่อปกป้องข้อมูลโครงการของคุณ</p>
      <button
        @click="openCreateModal"
        class="px-4 py-2 bg-purple-100 dark:bg-purple-600/20 border border-purple-300 dark:border-purple-500/30 text-purple-300 text-sm rounded-lg hover:bg-purple-100 dark:bg-purple-600/30 transition-colors"
      >
        สร้าง Backup แรก
      </button>
    </div>

    <!-- Backup list -->
    <div v-else class="space-y-3">
      <div
        v-for="backup in backups"
        :key="backup.id"
        class="bg-gray-800/60 border border-gray-700/60 rounded-xl p-4 flex items-center gap-4 hover:border-gray-600/60 transition-colors group"
      >
        <!-- Icon -->
        <div class="w-10 h-10 rounded-lg bg-purple-600/20 border border-purple-500/30 flex items-center justify-center shrink-0">
          <svg class="w-5 h-5 text-purple-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 8h14M5 8a2 2 0 110-4h14a2 2 0 110 4M5 8v10a2 2 0 002 2h10a2 2 0 002-2V8m-9 4h4" />
          </svg>
        </div>

        <!-- Info -->
        <div class="flex-1 min-w-0">
          <p class="text-sm font-semibold text-white truncate">{{ backup.label || 'Backup ไม่มีชื่อ' }}</p>
          <p class="text-xs text-gray-400 mt-0.5">{{ formatDate(backup.created_at) }}</p>
        </div>

        <!-- Actions -->
        <div class="flex items-center gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
          <button
            @click="downloadBackup(backup)"
            :disabled="downloading === backup.id"
            class="flex items-center gap-1.5 px-3 py-1.5 bg-blue-100 dark:bg-blue-600/20 border border-blue-300 dark:border-blue-500/30 text-blue-300 text-xs font-medium rounded-lg hover:bg-blue-100 dark:bg-blue-600/30 transition-colors disabled:opacity-50"
            title="ดาวน์โหลดเป็นไฟล์ JSON"
          >
            <div v-if="downloading === backup.id" class="w-3.5 h-3.5 border border-blue-300 border-t-transparent rounded-full animate-spin" />
            <svg v-else class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
            </svg>
            Export
          </button>
          <button
            @click="confirmRestore(backup)"
            class="flex items-center gap-1.5 px-3 py-1.5 bg-green-100 dark:bg-green-600/20 border border-green-300 dark:border-green-500/30 text-green-300 text-xs font-medium rounded-lg hover:bg-green-100 dark:bg-green-600/30 transition-colors"
          >
            <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
            </svg>
            Restore
          </button>
          <button
            @click="confirmDelete(backup)"
            class="flex items-center gap-1.5 px-3 py-1.5 bg-red-100 dark:bg-red-600/20 border border-red-300 dark:border-red-500/30 text-red-300 text-xs font-medium rounded-lg hover:bg-red-100 dark:bg-red-600/30 transition-colors"
          >
            <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
            </svg>
            ลบ
          </button>
        </div>
      </div>
    </div>

    <!-- Info box -->
    <div class="bg-yellow-500/10 border border-yellow-500/20 rounded-xl p-4 flex gap-3">
      <svg class="w-5 h-5 text-yellow-400 shrink-0 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
      </svg>
      <div class="text-sm text-yellow-300/80 space-y-1">
        <p class="font-semibold text-yellow-300">คำเตือน: การ Restore</p>
        <p>การ restore จะ<strong>ลบ</strong> tasks, sprints, milestones, และ epics ปัจจุบันทั้งหมด แล้วนำข้อมูลใน snapshot กลับมา ควร backup ข้อมูลปัจจุบันก่อนทำการ restore</p>
        <p class="text-yellow-400/60 text-xs pt-1">💡 ใช้ปุ่ม <strong>Export</strong> เพื่อดาวน์โหลด backup เป็นไฟล์ <code>.json</code> แล้วนำไป <strong>Import</strong> ตอนสร้างโครงการใหม่ได้</p>
      </div>
    </div>
  </div>

  <!-- Create Backup Modal -->
  <Teleport to="body">
    <div v-if="showCreateModal" class="fixed inset-0 bg-black/80 backdrop-blur-sm flex items-center justify-center z-50 p-4" @click.self="showCreateModal = false">
      <div class="bg-gray-800 border border-gray-700 rounded-2xl shadow-2xl w-full max-w-md">
        <div class="flex items-center justify-between px-6 pt-5 pb-4 border-b border-gray-700/60">
          <h3 class="text-base font-bold text-white">สร้าง Backup ใหม่</h3>
          <button @click="showCreateModal = false" class="text-gray-500 hover:text-white transition-colors">✕</button>
        </div>
        <div class="px-6 py-5 space-y-4">
          <div>
            <label class="block text-xs font-medium text-gray-400 mb-1.5">ชื่อ Backup (ไม่บังคับ)</label>
            <input
              v-model="newBackupLabel"
              type="text"
              placeholder="เช่น Before Sprint 3, Pre-AI Plan, ..."
              class="w-full bg-gray-700 border border-gray-600 text-white text-sm rounded-lg px-3 py-2.5 focus:outline-none focus:ring-2 focus:ring-purple-500 placeholder-gray-500"
              @keydown.enter="handleCreate"
            />
          </div>
          <p class="text-xs text-gray-500">Snapshot จะบันทึก: Project metadata, Tasks, Sprints, Milestones, Epics</p>
        </div>
        <div class="px-6 pb-5 flex gap-3 justify-end">
          <button @click="showCreateModal = false" class="px-4 py-2 text-sm text-gray-400 hover:text-white transition-colors">ยกเลิก</button>
          <button
            @click="handleCreate"
            :disabled="creating"
            class="px-4 py-2 bg-gradient-to-r from-purple-100 dark:from-purple-600 to-pink-100 dark:to-pink-600 text-gray-900 dark:text-white text-sm font-semibold rounded-lg hover:opacity-90 transition-opacity disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
          >
            <div v-if="creating" class="w-3.5 h-3.5 border-2 border-white border-t-transparent rounded-full animate-spin" />
            {{ creating ? 'กำลังสร้าง...' : 'สร้าง Backup' }}
          </button>
        </div>
      </div>
    </div>
  </Teleport>

  <!-- Restore Confirmation Modal -->
  <Teleport to="body">
    <div v-if="backupToRestore" class="fixed inset-0 bg-black/80 backdrop-blur-sm flex items-center justify-center z-50 p-4" @click.self="backupToRestore = null">
      <div class="bg-gray-800 border border-gray-700 rounded-2xl shadow-2xl w-full max-w-md">
        <div class="flex items-center justify-between px-6 pt-5 pb-4 border-b border-gray-700/60">
          <h3 class="text-base font-bold text-white">ยืนยันการ Restore</h3>
          <button @click="backupToRestore = null" class="text-gray-500 hover:text-white transition-colors">✕</button>
        </div>
        <div class="px-6 py-5 space-y-3">
          <div class="bg-red-500/10 border border-red-500/20 rounded-xl p-4 text-sm text-red-300">
            <p class="font-semibold mb-1">การดำเนินการนี้ไม่สามารถยกเลิกได้</p>
            <p>ข้อมูลปัจจุบัน (tasks, sprints, milestones, epics) จะถูกลบและแทนที่ด้วยข้อมูลจาก snapshot นี้</p>
          </div>
          <div class="bg-gray-700/50 rounded-lg p-3">
            <p class="text-xs text-gray-400">Restore จาก:</p>
            <p class="text-sm font-semibold text-white mt-0.5">{{ backupToRestore.label || 'Backup ไม่มีชื่อ' }}</p>
            <p class="text-xs text-gray-500 mt-0.5">{{ formatDate(backupToRestore.created_at) }}</p>
          </div>
        </div>
        <div class="px-6 pb-5 flex gap-3 justify-end">
          <button @click="backupToRestore = null" class="px-4 py-2 text-sm text-gray-400 hover:text-white transition-colors">ยกเลิก</button>
          <button
            @click="handleRestore"
            :disabled="restoring"
            class="px-4 py-2 bg-green-100 dark:bg-green-600 hover:bg-green-100 dark:bg-green-500 text-gray-900 dark:text-white text-sm font-semibold rounded-lg transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
          >
            <div v-if="restoring" class="w-3.5 h-3.5 border-2 border-white border-t-transparent rounded-full animate-spin" />
            {{ restoring ? 'กำลัง Restore...' : 'ยืนยัน Restore' }}
          </button>
        </div>
      </div>
    </div>
  </Teleport>

  <!-- Delete Confirmation Modal -->
  <Teleport to="body">
    <div v-if="backupToDelete" class="fixed inset-0 bg-black/80 backdrop-blur-sm flex items-center justify-center z-50 p-4" @click.self="backupToDelete = null">
      <div class="bg-gray-800 border border-gray-700 rounded-2xl shadow-2xl w-full max-w-md">
        <div class="flex items-center justify-between px-6 pt-5 pb-4 border-b border-gray-700/60">
          <h3 class="text-base font-bold text-white">ลบ Backup</h3>
          <button @click="backupToDelete = null" class="text-gray-500 hover:text-white transition-colors">✕</button>
        </div>
        <div class="px-6 py-5">
          <p class="text-sm text-gray-300">ต้องการลบ backup <strong class="text-white">{{ backupToDelete.label || 'ไม่มีชื่อ' }}</strong> ใช่หรือไม่? การดำเนินการนี้ไม่สามารถยกเลิกได้</p>
        </div>
        <div class="px-6 pb-5 flex gap-3 justify-end">
          <button @click="backupToDelete = null" class="px-4 py-2 text-sm text-gray-400 hover:text-white transition-colors">ยกเลิก</button>
          <button
            @click="handleDelete"
            :disabled="deleting"
            class="px-4 py-2 bg-red-100 dark:bg-red-600 hover:bg-red-100 dark:bg-red-500 text-gray-900 dark:text-white text-sm font-semibold rounded-lg transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
          >
            <div v-if="deleting" class="w-3.5 h-3.5 border-2 border-white border-t-transparent rounded-full animate-spin" />
            {{ deleting ? 'กำลังลบ...' : 'ลบ Backup' }}
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useProjectsApi } from '../infrastructure/projects-api'
import type { ProjectBackup } from '../infrastructure/projects-api'

const props = defineProps<{
  projectId: string
}>()

const emit = defineEmits<{
  (e: 'restored'): void
}>()

const projectsApi = useProjectsApi()

const backups = ref<ProjectBackup[]>([])
const loading = ref(false)
const creating = ref(false)
const restoring = ref(false)
const deleting = ref(false)
const downloading = ref<string | null>(null)

const showCreateModal = ref(false)
const newBackupLabel = ref('')
const backupToRestore = ref<ProjectBackup | null>(null)
const backupToDelete = ref<ProjectBackup | null>(null)

async function loadBackups() {
  loading.value = true
  try {
    backups.value = await projectsApi.getProjectBackups(props.projectId)
  } catch (err) {
    console.error('Failed to load backups:', err)
  } finally {
    loading.value = false
  }
}

function openCreateModal() {
  newBackupLabel.value = ''
  showCreateModal.value = true
}

async function handleCreate() {
  creating.value = true
  try {
    await projectsApi.createProjectBackup(props.projectId, newBackupLabel.value || undefined)
    showCreateModal.value = false
    await loadBackups()
  } catch (err) {
    console.error('Failed to create backup:', err)
    alert('ไม่สามารถสร้าง backup ได้: ' + (err instanceof Error ? err.message : String(err)))
  } finally {
    creating.value = false
  }
}

async function downloadBackup(backup: ProjectBackup) {
  downloading.value = backup.id
  try {
    const full = await projectsApi.getProjectBackupPayload(props.projectId, backup.id)
    const blob = new Blob([JSON.stringify(full, null, 2)], { type: 'application/json' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    const safeName = (backup.label || 'backup').replace(/[^a-z0-9]/gi, '_').toLowerCase()
    a.download = `${safeName}_${backup.id.slice(0, 8)}.sentinel.json`
    a.click()
    URL.revokeObjectURL(url)
  } catch (err) {
    console.error('Failed to download backup:', err)
    alert('ไม่สามารถดาวน์โหลด backup ได้: ' + (err instanceof Error ? err.message : String(err)))
  } finally {
    downloading.value = null
  }
}

function confirmRestore(backup: ProjectBackup) {
  backupToRestore.value = backup
}

async function handleRestore() {
  if (!backupToRestore.value) return
  restoring.value = true
  try {
    await projectsApi.restoreProjectBackup(props.projectId, backupToRestore.value.id)
    backupToRestore.value = null
    emit('restored')
    await loadBackups()
  } catch (err) {
    console.error('Failed to restore backup:', err)
    alert('ไม่สามารถ restore ได้: ' + (err instanceof Error ? err.message : String(err)))
  } finally {
    restoring.value = false
  }
}

function confirmDelete(backup: ProjectBackup) {
  backupToDelete.value = backup
}

async function handleDelete() {
  if (!backupToDelete.value) return
  deleting.value = true
  try {
    await projectsApi.deleteProjectBackup(props.projectId, backupToDelete.value.id)
    backupToDelete.value = null
    await loadBackups()
  } catch (err) {
    console.error('Failed to delete backup:', err)
    alert('ไม่สามารถลบ backup ได้: ' + (err instanceof Error ? err.message : String(err)))
  } finally {
    deleting.value = false
  }
}

function formatDate(iso: string): string {
  const d = new Date(iso)
  return d.toLocaleString('th-TH', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

onMounted(loadBackups)
</script>

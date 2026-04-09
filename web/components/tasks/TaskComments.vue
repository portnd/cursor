<template>
  <div class="task-comments">
    <!-- Comments List -->
    <div class="space-y-4 mb-6 max-h-96 overflow-y-auto">
      <div v-if="!comments.length" class="text-center py-8 text-gray-500 text-sm">
        No comments yet. Start the discussion.
      </div>
      <div
        v-for="comment in comments"
        :key="comment.id"
        class="flex gap-3"
      >
        <div class="flex-shrink-0 w-8 h-8 rounded-full overflow-hidden">
          <img
            v-if="comment.user_avatar_url"
            :src="comment.user_avatar_url"
            :alt="comment.user_email || String(comment.user_id)"
            class="w-full h-full object-cover"
          />
          <div
            v-else
            class="w-full h-full bg-purple-600 flex items-center justify-center text-white text-sm font-bold"
          >
            {{ (comment.user_display_name || comment.user_email || String(comment.user_id)).charAt(0).toUpperCase() }}
          </div>
        </div>
        <div class="flex-1 bg-gray-800 rounded-xl px-4 py-3 border border-gray-700/50">
          <div class="flex items-center justify-between mb-1.5 gap-2">
            <span class="text-sm font-medium text-gray-300">{{ comment.user_display_name || comment.user_email || `User #${comment.user_id}` }}</span>
            <div class="flex items-center gap-2">
              <span class="text-xs text-gray-500">{{ formatTime(comment.created_at) }}</span>
              <button
                v-if="canEditComment(comment) && editingCommentId !== comment.id"
                class="text-xs text-purple-300 hover:text-purple-200"
                @click="startEdit(comment)"
              >Edit</button>
            </div>
          </div>

          <template v-if="editingCommentId === comment.id">
            <textarea
              v-model="editingContent"
              rows="3"
              class="w-full bg-gray-900 border border-gray-700 rounded-lg px-3 py-2 text-sm text-gray-200 placeholder-gray-500 focus:outline-none focus:border-purple-500 resize-none"
            />
            <div class="mt-2 flex items-center justify-end gap-2">
              <button class="px-3 py-1.5 rounded-lg text-xs border border-gray-600 text-gray-300 hover:bg-gray-700/50" @click="cancelEdit">Cancel</button>
              <button
                class="px-3 py-1.5 rounded-lg text-xs bg-purple-100 dark:bg-purple-600 hover:bg-purple-100 dark:bg-purple-500 text-gray-900 dark:text-white disabled:opacity-50"
                :disabled="editSubmitting || !editingContent.trim()"
                @click="submitEdit(comment.id)"
              >{{ editSubmitting ? 'Saving...' : 'Save' }}</button>
            </div>
          </template>

          <template v-else>
            <p v-if="comment.content" class="text-sm text-gray-300 leading-relaxed whitespace-pre-wrap">{{ comment.content }}</p>
            <div class="mt-1 flex items-center gap-2" v-if="comment.edited_at || (comment.edit_history?.length || 0) > 0">
              <span class="text-[11px] text-amber-300">edited</span>
              <button
                class="text-[11px] text-purple-300 hover:text-purple-200"
                @click="toggleHistory(comment.id)"
              >{{ showHistoryForId === comment.id ? 'Hide history' : 'View history' }}</button>
            </div>
            <div v-if="showHistoryForId === comment.id" class="mt-2 rounded-lg border border-gray-700/70 bg-gray-900/60 p-2.5 space-y-2">
              <div
                v-for="(h, idx) in (comment.edit_history || [])"
                :key="`${comment.id}-hist-${idx}`"
                class="text-xs text-gray-300"
              >
                <div class="text-gray-400 mb-1">{{ formatTime(h.edited_at) }}</div>
                <div class="grid grid-cols-1 md:grid-cols-2 gap-2">
                  <div>
                    <div class="text-[11px] text-gray-500 mb-0.5">Before</div>
                    <div class="whitespace-pre-wrap text-gray-400">{{ h.old_content }}</div>
                  </div>
                  <div>
                    <div class="text-[11px] text-gray-500 mb-0.5">After</div>
                    <div class="whitespace-pre-wrap">{{ h.new_content }}</div>
                  </div>
                </div>
              </div>
            </div>
          </template>

          <div v-if="comment.attachments?.length" class="mt-3 space-y-2">
            <div v-for="(att, idx) in comment.attachments" :key="`${comment.id}-${idx}`">
              <img
                v-if="att.is_image"
                :src="att.data_url"
                :alt="att.file_name"
                class="max-h-64 rounded-lg border border-gray-700 object-contain bg-gray-900 cursor-zoom-in"
                @click="openImagePreview(att.data_url, att.file_name)"
              />
              <a
                v-else
                :href="att.data_url"
                :download="att.file_name"
                class="inline-flex items-center gap-2 text-xs text-purple-300 hover:text-purple-200"
              >
                <span>📎</span>
                <span>{{ att.file_name }}</span>
                <span class="text-gray-500">({{ formatFileSize(att.size) }})</span>
              </a>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Add Comment Input -->
    <div class="flex gap-3">
      <div class="flex-shrink-0 w-8 h-8 rounded-full overflow-hidden">
        <img
          v-if="currentUserAvatar"
          :src="currentUserAvatar"
          alt="You"
          class="w-full h-full object-cover"
        />
        <div
          v-else-if="currentUserInitial"
          class="w-full h-full bg-gradient-to-br from-purple-500 to-pink-500 flex items-center justify-center text-white text-sm font-bold"
        >{{ currentUserInitial }}</div>
        <div
          v-else
          class="w-full h-full bg-gray-700 flex items-center justify-center text-gray-400 text-sm"
        >You</div>
      </div>
      <div class="flex-1 flex flex-col gap-2">
        <textarea
          v-model="newComment"
          placeholder="Write a comment... (Shift+Enter for new line)"
          rows="2"
          class="flex-1 bg-gray-800 border border-gray-700 rounded-xl px-4 py-2.5 text-sm text-gray-200 placeholder-gray-500 focus:outline-none focus:border-purple-500 resize-none"
          @keydown.enter.prevent="handleEnter"
        ></textarea>
        <div class="flex items-center justify-between gap-2">
          <label class="text-xs text-gray-400 hover:text-gray-200 cursor-pointer">
            <input type="file" multiple class="hidden" @change="handleFileChange" />
            📎 Attach files / images
          </label>
          <button
            @click="submitComment"
            :disabled="(!newComment.trim() && !attachments.length) || loading"
            class="btn-primary self-end px-4 py-2 text-sm disabled:opacity-40"
          >
            {{ loading ? '...' : 'Post' }}
          </button>
        </div>
        <div v-if="attachments.length" class="space-y-1.5">
          <div
            v-for="(file, idx) in attachments"
            :key="`${file.name}-${idx}`"
            class="flex items-center justify-between rounded-md border border-gray-700 px-2 py-1 text-xs text-gray-300"
          >
            <span class="truncate">{{ file.name }} ({{ formatFileSize(file.size) }})</span>
            <button class="text-red-300 hover:text-red-200" @click="removeAttachment(idx)">Remove</button>
          </div>
        </div>
      </div>
    </div>
  </div>
  <div
    v-if="previewImageUrl"
    class="fixed inset-0 z-[9999] bg-black/80 flex items-center justify-center p-4"
    @click="closeImagePreview"
  >
    <div class="relative max-w-[95vw] max-h-[95vh]" @click.stop>
      <button
        class="absolute -top-3 -right-3 w-8 h-8 rounded-full bg-gray-800 border border-gray-600 text-gray-200 hover:bg-gray-700"
        @click="closeImagePreview"
      >
        ×
      </button>
      <img
        :src="previewImageUrl"
        :alt="previewImageName"
        class="max-w-[95vw] max-h-[95vh] rounded-lg border border-gray-700 object-contain"
      />
      <div class="mt-2 text-center text-xs text-gray-300">{{ previewImageName }}</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { TaskComment } from '~/core/modules/tasks/infrastructure/tasks-api'

const props = defineProps<{
  comments: TaskComment[]
  loading?: boolean
  currentUserAvatar?: string
  currentUserInitial?: string
  currentUserId?: number
}>()

const emit = defineEmits<{
  (e: 'add-comment', payload: { content: string; attachments: File[] }): void
  (e: 'edit-comment', payload: { commentId: string; content: string }): void
}>()

const newComment = ref('')
const attachments = ref<File[]>([])
const previewImageUrl = ref('')
const previewImageName = ref('')
const editingCommentId = ref<string | null>(null)
const editingContent = ref('')
const editSubmitting = ref(false)
const showHistoryForId = ref<string | null>(null)

function formatTime(dateStr: string) {
  const d = new Date(dateStr)
  const now = new Date()
  const diff = now.getTime() - d.getTime()
  const mins = Math.floor(diff / 60000)
  if (mins < 1) return 'just now'
  if (mins < 60) return `${mins}m ago`
  const hrs = Math.floor(mins / 60)
  if (hrs < 24) return `${hrs}h ago`
  return d.toLocaleDateString('en-US', { month: 'short', day: 'numeric' })
}

function handleEnter(e: KeyboardEvent) {
  if (!e.shiftKey) {
    submitComment()
  } else {
    newComment.value += '\n'
  }
}

function submitComment() {
  const content = newComment.value.trim()
  if (!content && attachments.value.length === 0) return
  emit('add-comment', { content, attachments: [...attachments.value] })
  newComment.value = ''
  attachments.value = []
}

function handleFileChange(e: Event) {
  const input = e.target as HTMLInputElement
  const files = Array.from(input.files || [])
  if (!files.length) return
  attachments.value.push(...files)
  input.value = ''
}

function removeAttachment(idx: number) {
  attachments.value.splice(idx, 1)
}

function formatFileSize(size: number) {
  if (size < 1024) return `${size} B`
  if (size < 1024 * 1024) return `${(size / 1024).toFixed(1)} KB`
  return `${(size / (1024 * 1024)).toFixed(1)} MB`
}

function openImagePreview(url: string, name: string) {
  previewImageUrl.value = url
  previewImageName.value = name
}

function canEditComment(comment: TaskComment) {
  return Number(props.currentUserId || 0) > 0 && Number(comment.user_id) === Number(props.currentUserId)
}

function startEdit(comment: TaskComment) {
  editingCommentId.value = comment.id
  editingContent.value = comment.content || ''
}

function cancelEdit() {
  editingCommentId.value = null
  editingContent.value = ''
  editSubmitting.value = false
}

function submitEdit(commentId: string) {
  const content = editingContent.value.trim()
  if (!content) return
  editSubmitting.value = true
  emit('edit-comment', { commentId, content })
  setTimeout(() => {
    editSubmitting.value = false
    cancelEdit()
  }, 250)
}

function toggleHistory(commentId: string) {
  showHistoryForId.value = showHistoryForId.value === commentId ? null : commentId
}

function closeImagePreview() {
  previewImageUrl.value = ''
  previewImageName.value = ''
}
</script>

<style scoped>
.btn-primary {
  @apply bg-gradient-to-r from-purple-600 to-pink-600 hover:from-purple-700 hover:to-pink-700 text-white rounded-lg font-medium transition-colors;
}
</style>

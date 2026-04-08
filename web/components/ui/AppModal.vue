<template>
  <Teleport to="body">
    <!-- Notification modal (success / error / info) -->
    <Transition name="modal">
      <div
        v-if="notification.visible"
        class="fixed inset-0 z-[100] flex items-center justify-center p-4"
        role="dialog"
        aria-modal="true"
        :aria-labelledby="notificationTitleId"
        @keydown.escape="closeNotification"
      >
        <div
          class="fixed inset-0 bg-black/70 backdrop-blur-sm"
          aria-hidden="true"
          @click="closeNotification"
        />
        <div
          :class="[
            'relative w-full max-w-md rounded-2xl border-2 shadow-2xl overflow-hidden',
            notificationClasses
          ]"
          role="document"
          @click.stop
        >
          <div class="p-6">
            <div class="flex items-start gap-4">
              <span
                :class="[
                  'flex h-12 w-12 shrink-0 items-center justify-center rounded-full text-2xl',
                  iconBgClass
                ]"
                aria-hidden="true"
              >
                {{ iconEmoji }}
              </span>
              <div class="min-w-0 flex-1">
                <h2
                  :id="notificationTitleId"
                  :class="['text-lg font-semibold', titleClass]"
                >
                  {{ notification.title }}
                </h2>
                <p :class="['mt-2 text-sm leading-relaxed', messageClass]">
                  {{ notification.message }}
                </p>
              </div>
            </div>
            <div class="mt-6 flex justify-end">
              <button
                type="button"
                :class="[
                  'px-5 py-2.5 rounded-xl font-medium transition-colors',
                  okButtonClass
                ]"
                @click="closeNotification"
              >
                OK
              </button>
            </div>
          </div>
        </div>
      </div>
    </Transition>

    <!-- Confirm modal -->
    <Transition name="modal">
      <div
        v-if="confirmState.visible"
        class="fixed inset-0 z-[100] flex items-center justify-center p-4"
        role="dialog"
        aria-modal="true"
        aria-labelledby="confirmTitleId"
        @keydown.escape="closeConfirm(false)"
      >
        <div
          class="fixed inset-0 bg-black/70 backdrop-blur-sm"
          aria-hidden="true"
          @click="closeConfirm(false)"
        />
        <div
          :class="['relative w-full max-w-md rounded-2xl border-2 shadow-2xl overflow-hidden', confirmPanelClass]"
          role="document"
          @click.stop
        >
          <div class="p-6">
            <h2
              id="confirmTitleId"
              :class="['text-lg font-semibold', confirmTitleTextClass]"
            >
              {{ confirmState.title }}
            </h2>
            <p :class="['mt-3 text-sm leading-relaxed', confirmMessageTextClass]">
              {{ confirmState.message }}
            </p>
            <div class="mt-6 flex flex-wrap gap-3 justify-end">
              <button
                type="button"
                :class="['px-5 py-2.5 rounded-xl font-medium transition-colors', cancelBtnClass]"
                @click="closeConfirm(false)"
              >
                {{ confirmState.cancelLabel }}
              </button>
              <button
                type="button"
                :class="[
                  'px-5 py-2.5 rounded-xl font-medium transition-colors',
                  confirmButtonClass
                ]"
                @click="closeConfirm(true)"
              >
                {{ confirmState.confirmLabel }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import type { ConfirmVariant } from '~/composables/useNotification'

const notificationTitleId = 'app-modal-notification-title'
const confirmTitleId = 'app-modal-confirm-title'

const { notification, confirmState, closeNotification, closeConfirm } = useNotification()
const { isDark } = useTheme()

const iconEmoji = computed(() => {
  switch (notification.value.type) {
    case 'success': return '✓'
    case 'error': return '✕'
    case 'info': return 'ℹ'
    default: return '•'
  }
})

const notificationClasses = computed(() => {
  const t = notification.value.type
  const bg = isDark.value ? 'bg-gray-800' : 'bg-white'
  if (t === 'success') return `border-emerald-500/60 ${bg}`
  if (t === 'error') return `border-red-500/60 ${bg}`
  return `border-gray-500/60 ${bg}`
})

const iconBgClass = computed(() => {
  const t = notification.value.type
  if (t === 'success') return 'bg-emerald-600/90 text-white'
  if (t === 'error') return 'bg-red-600/90 text-white'
  return isDark.value ? 'bg-gray-600/90 text-white' : 'bg-slate-200 text-slate-700'
})

const titleClass = computed(() => {
  const t = notification.value.type
  if (t === 'success') return isDark.value ? 'text-emerald-400' : 'text-emerald-600'
  if (t === 'error') return isDark.value ? 'text-red-400' : 'text-red-600'
  return isDark.value ? 'text-gray-300' : 'text-slate-700'
})

const messageClass = computed(() => isDark.value ? 'text-gray-300' : 'text-slate-600')

const okButtonClass = computed(() => {
  const t = notification.value.type
  if (t === 'success') return 'bg-emerald-600 text-white hover:bg-emerald-500'
  if (t === 'error') return 'bg-red-600 text-white hover:bg-red-500'
  return isDark.value ? 'bg-gray-600 text-white hover:bg-gray-500' : 'bg-slate-200 text-slate-800 hover:bg-slate-300'
})

const confirmPanelClass = computed(() =>
  isDark.value
    ? 'border-gray-600 bg-gray-800'
    : 'border-slate-200 bg-white'
)
const confirmTitleTextClass = computed(() => isDark.value ? 'text-white' : 'text-slate-900')
const confirmMessageTextClass = computed(() => isDark.value ? 'text-gray-300' : 'text-slate-600')
const cancelBtnClass = computed(() =>
  isDark.value
    ? 'bg-gray-700 text-gray-200 hover:bg-gray-600'
    : 'bg-slate-100 text-slate-700 hover:bg-slate-200'
)

const confirmButtonClass = computed(() => {
  const v: ConfirmVariant = confirmState.value.variant ?? 'primary'
  if (v === 'danger') return 'bg-red-600 text-white hover:bg-red-500'
  if (v === 'primary') return 'bg-gradient-to-r from-purple-600 to-pink-600 text-white hover:from-purple-500 hover:to-pink-500'
  return isDark.value ? 'bg-gray-600 text-white hover:bg-gray-500' : 'bg-slate-200 text-slate-800 hover:bg-slate-300'
})
</script>

<style scoped>
.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.2s ease, transform 0.2s ease;
}
.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}
.modal-enter-from > div:last-child,
.modal-leave-to > div:last-child {
  transform: scale(0.95);
}
.modal-enter-active > div:last-child,
.modal-leave-active > div:last-child {
  transform: scale(1);
  transition: transform 0.2s ease;
}
</style>

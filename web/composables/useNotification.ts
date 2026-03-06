/**
 * Unified notification and confirmation modal.
 * Use showSuccess/showError/showInfo for result feedback;
 * use confirm() for yes/no confirmations (returns Promise<boolean>).
 */

export type NotificationType = 'success' | 'error' | 'info'
export type ConfirmVariant = 'danger' | 'primary' | 'neutral'

export interface NotificationState {
  type: NotificationType
  title: string
  message: string
  visible: boolean
}

export interface ConfirmOptions {
  title: string
  message: string
  confirmLabel?: string
  cancelLabel?: string
  variant?: ConfirmVariant
}

export interface ConfirmState extends ConfirmOptions {
  visible: boolean
  resolve: ((value: boolean) => void) | null
}

// Singleton state so all callers share the same modal
const notification = ref<NotificationState>({
  type: 'info',
  title: '',
  message: '',
  visible: false
})

const confirmState = ref<ConfirmState>({
  title: '',
  message: '',
  confirmLabel: 'Confirm',
  cancelLabel: 'Cancel',
  variant: 'primary',
  visible: false,
  resolve: null
})

export function useNotification() {
  function closeNotification() {
    notification.value.visible = false
  }

  function showSuccess(message: string, title = 'Success') {
    notification.value = {
      type: 'success',
      title,
      message,
      visible: true
    }
  }

  function showError(message: string, title = 'Error') {
    notification.value = {
      type: 'error',
      title,
      message,
      visible: true
    }
  }

  function showInfo(message: string, title = 'Info') {
    notification.value = {
      type: 'info',
      title,
      message,
      visible: true
    }
  }

  function confirm(options: ConfirmOptions): Promise<boolean> {
    return new Promise((resolve) => {
      confirmState.value = {
        title: options.title,
        message: options.message,
        confirmLabel: options.confirmLabel ?? 'Confirm',
        cancelLabel: options.cancelLabel ?? 'Cancel',
        variant: options.variant ?? 'primary',
        visible: true,
        resolve
      }
    })
  }

  function closeConfirm(result: boolean) {
    const r = confirmState.value.resolve
    confirmState.value = {
      ...confirmState.value,
      visible: false,
      resolve: null
    }
    r?.(result)
  }

  return {
    notification: readonly(notification),
    confirmState: readonly(confirmState),
    closeNotification,
    showSuccess,
    showError,
    showInfo,
    confirm,
    closeConfirm
  }
}

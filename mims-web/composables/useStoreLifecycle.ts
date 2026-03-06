import type { Store } from 'pinia'

interface StoreLifecycleOptions {
	/** Manual cleanup to run before unmount (e.g. clearInterval, unsubscribe) */
	cleanup?: () => void
	/** When true, calls $reset() in onBeforeMount to clear stale data. Default: true */
	resetOnEnter?: boolean
}

/**
 * Replaces the broken `onUnmounted(() => store.$dispose())` pattern.
 *
 * $dispose() kills the Pinia effect scope which breaks computed getters
 * when Nuxt's transition overlap causes the old page to dispose a store
 * while the new page is already referencing it.
 *
 * This composable uses $reset() on enter instead, keeping the effect scope alive.
 */
export function useStoreLifecycle(
	stores: Store | Store[],
	options: StoreLifecycleOptions = {}
) {
	const { cleanup, resetOnEnter = true } = options
	const storeArray = Array.isArray(stores) ? stores : [stores]

	if (resetOnEnter) {
		onBeforeMount(() => {
			storeArray.forEach((s) => s.$reset())
		})
	}

	onUnmounted(() => {
		cleanup?.()
	})
}

<script lang="ts" setup>
import { useLogoutStore } from "~~/core/modules/auth/logout/store"
import { useAuthStore } from "~~/core/modules/auth/authentication/store"

definePageMeta({
	layout: "auth",
})

const store = useLogoutStore()
useStoreLifecycle(store)
const logout = await store.logout()

if (logout) {
	const authStore = useAuthStore()
	authStore.logout()
}

</script>

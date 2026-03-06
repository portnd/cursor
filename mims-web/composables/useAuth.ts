import { useAuthStore } from "~~/core/modules/auth/authentication/store"

const useAuth = () => {
	const isLoggedIn = () => {
		const authStore = useAuthStore()
		return authStore.isLoggedIn
	}

	return {
		isLoggedIn,
	}
}

export default useAuth

export default defineNuxtRouteMiddleware(async () => {
	const { isLoggedIn } = useAuth()
	if (!isLoggedIn()) {
		return await navigateTo("/auth/login")
	}
	return await navigateTo(useDefaultLandingUrl())
})

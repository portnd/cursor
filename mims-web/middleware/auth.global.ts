import useAuth from "~~/composables/useAuth"

export default defineNuxtRouteMiddleware(async (to, _) => {
	const path = to.path
	const excepts = ["auth"] // Path ที่ละเว้นการตรวจสอบ

	if (!exceptPath(path, excepts)) {
		const { isLoggedIn } = useAuth()
		if (!isLoggedIn()) {
			return await navigateTo("/auth/login")
		}
	}
})

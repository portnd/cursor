import usePermission from '../../../composables/usePermission';
<script setup lang="ts">
import httpStatusCode from "~~/core/shared/http/HttpStatusCode"

definePageMeta({
	validate: (route) => {
		if (!isNumber(route.params.roadId)) {
			return false
		}
		return true
	},
	middleware: [
		function (to) {
			const roadId = to.params.roadId
			console.log("❄️❄️🇹🇭 roadid=,", roadId)
			console.log("❄️❄️🇹🇭 to=,", to.path)
			// if (from.fullPath.startsWith(`/roads/${roadId}/`)) {
			// 	return navigateTo(`/roads`)
			// } else {
			const canAccessObj = [
				{
					canAccess: usePermission().checkCanAccessPermissionBYKey(IUserRolesAccess.view_road_in_assets),
					url: `/roads/${to.params.roadId}/in-asset`,
				},
				{
					canAccess: usePermission().checkCanAccessPermissionBYKey(IUserRolesAccess.view_road_condition),
					url: `/roads/${to.params.roadId}/condition`,
				},
				{
					canAccess: usePermission().checkCanAccessPermissionBYKey(IUserRolesAccess.view_road_retro),
					url: `/roads/${to.params.roadId}/reflective-strip`,
				},
				{
					canAccess: usePermission().checkCanAccessPermissionBYKey(IUserRolesAccess.view_road_damage),
					url: `/roads/${to.params.roadId}/damage`,
				},
			]

			console.log("⛺️ canAccess sm:", canAccessObj)
			const canAccess = usePermission().checkCanAccessPermissionBYKey(IUserRolesAccess.view_road_summary)
			if (canAccess) {
				console.log("💚 canAccess ", canAccess)
				return navigateTo(`/roads/${roadId}/summary`)
			} else {
				for (let i = 0; i < canAccessObj.length; i++) {
					const access = canAccessObj[i]
					if (access.canAccess) {
						console.log("access url:", access.url)
						return navigateTo(access.url)
					}
					if (i === canAccessObj.length - 1) {
						showError({ statusCode: httpStatusCode.ACCESS_DENIED })
					}
				}
			}
			// }
		},
	],
})
</script>

<template>
	<div></div>
</template>

<style scoped></style>

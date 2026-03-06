/**
 * Returns the default URL for a logged-in user based on permissions.
 * Use after initUser has run (e.g. post-login or on landing page).
 */
export function useDefaultLandingUrl(): string {
	const permission = usePermission()
	const canAccessObj: { canAccess: boolean; url: string }[] = [
		{
			canAccess:
				permission.checkCanAccessPermissionBYKey(IUserRolesAccess.view_total_dashboard) ||
				permission.checkCanAccessPermissionBYKey(IUserRolesAccess.view_all_asset_dashboard) ||
				permission.checkCanAccessPermissionBYKey(IUserRolesAccess.view_owner_asset_dashboard) ||
				permission.checkCanAccessPermissionBYKey(IUserRolesAccess.view_all_road_condition_dashboard) ||
				permission.checkCanAccessPermissionBYKey(IUserRolesAccess.view_owner_road_condition_dashboard) ||
				permission.checkCanAccessPermissionBYKey(IUserRolesAccess.view_all_surface_dashboard) ||
				permission.checkCanAccessPermissionBYKey(IUserRolesAccess.view_owner_surface_dashboard) ||
				permission.checkCanAccessPermissionBYKey(IUserRolesAccess.view_all_maint_history_dashboard) ||
				permission.checkCanAccessPermissionBYKey(IUserRolesAccess.view_owner_maint_history_dashboard),
			url: "/dashboards",
		},
		{
			canAccess:
				permission.checkCanAccessPermissionBYKey(IUserRolesAccess.view_all_road) ||
				permission.checkCanAccessPermissionBYKey(IUserRolesAccess.view_owner_road),
			url: "/roads",
		},
		{
			canAccess:
				permission.checkCanAccessPermissionBYKey(IUserRolesAccess.view_all_maint_history) ||
				permission.checkCanAccessPermissionBYKey(IUserRolesAccess.view_owner_maint_history) ||
				permission.checkCanAccessPermissionBYKey(IUserRolesAccess.manage_all_maint_history) ||
				permission.checkCanAccessPermissionBYKey(IUserRolesAccess.manage_owner_maint_history),
			url: "/maintenances",
		},
		{
			canAccess:
				permission.checkCanAccessPermissionBYKey(IUserRolesAccess.view_all_maintenance_analysis) ||
				permission.checkCanAccessPermissionBYKey(IUserRolesAccess.view_owner_maintenance_analysis) ||
				permission.checkCanAccessPermissionBYKey(IUserRolesAccess.manage_myself_maintenance_analysis),
			url: "/analyses",
		},
		{ canAccess: true, url: "/profile" },
	]
	for (const access of canAccessObj) {
		if (access.canAccess) return access.url
	}
	return "/profile"
}

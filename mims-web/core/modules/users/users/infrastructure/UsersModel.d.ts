export interface IUsersRolesData {
	current_page: number
	next_page: number
	previous_page: number
	size_per_page: number
	total_pages: number
	total_items: number
	items: IUserRolesItem[]
}

export interface IUserRolesItem {
	id: number
	name: string
}

export interface IUsersDepartmentsData {
	current_page: number
	next_page: number
	previous_page: number
	size_per_page: number
	total_pages: number
	total_items: number
	items: IUsersDepartmentItems[]
}

export interface IUsersDepartmentItems {
	id: number
	name: string
	can_delete: boolean
}

export interface IDefaultUsersData {
	id: number
	email: string
	username: string
	department_id: number
	department: IDefaultUsersDepartment
	firstname: string
	lastname: string
	profile_img_path: string
	status: boolean
	tel: string
	created_by: number
	updated_by: number
	roles: IDefaultUsersRole[]
	ref_user_owner_id: number | null
	ref_depot_id: number | null
}

export interface IDefaultUsersDepartment {
	id: number
	name: string
	can_delete: boolean
}

export interface IDefaultUsersRole {
	id: number
	name: string
	is_checked: boolean
}

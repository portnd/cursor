export interface IRequestUserSearch {
	fullname: string
	username: string
	ref_user_owner_id: number | null
	ref_depot_id: number | null
	permission: number | string
	status: boolean | string
}
export interface IRequestUsers {
	created_by: number
	ref_user_owner_id: number | null
	ref_depot_id: number | null
	email: string
	firstname: string
	lastname: string
	profile_img_path: string
	roles: number[]
	tel: string
	status: boolean
	updated_by: number
	username: string
	password: string
}

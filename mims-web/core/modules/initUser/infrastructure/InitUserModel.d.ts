export interface IUser {
	id: number
	email: string
	title_name: string
	firstname: string
	lastname: string
	profile_img_part: string
	access_control: IAccessControl[]
	username: string
	tel: string
	ref_user_owner_id: number
	ref_depot_id: number
	ref_user_owner: IRefUserOwner
	ref_depot: IDepot
	roles: IUserRole[]
}
export interface IAccessControl {
	id: number
	access_title: string
	access_desc: string
	access_grp_id: number
	access_key: string
}

export interface IUserRole {
	id: number
	name: string
}

export interface IDepot {
	id: number
	name: string
	depot_code: string
}

export interface IRefUserOwner {
	id: number
	email: string
}

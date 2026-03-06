export interface IAccount {
	id: number
	email: string
	department_id: number
	title_name: string
	firstname: string
	lastname: string
	profile_img_path: string
	department: IDepartment
	access_control: IAccessControl[]
	username: string
	roles: IAccountRole[]
	tel: string
	ref_depot: IDepot
	ref_user_owner: IRefUserOwner
}

export interface IAccountRole {
	id: number
	name: string
}

export interface IAccessControl {
	id: number
	access_title: string
	access_desc: string
	access_grp_id: number
	access_key: string
}

export interface IDepartment {
	id: number
	name: string
	can_delete: boolean
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

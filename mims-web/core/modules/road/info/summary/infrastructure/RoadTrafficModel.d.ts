export interface ITrafficModel {
	year: number
	items: ITrafficItem[]
}

export interface ITrafficItem {
	id: number
	status: string
	year: number
	revision: number
	id_parent: number
	surveyed_date: string
}

export interface ITrafficDetail {
	id: number
	road_id: number
	year: number
	created_by: number
	created_date: string
	updated_date: string
	revision: number
	id_parent: number
	reject_reason: string
	veh1: number
	veh2: number
	veh3: number
	total: number
	aadt: number
	esal: number
	yax: number
	surveyed_date: Date
	hash_data: string
	status: string
	status_code: string
	updated_by: IUpdatedBy
}

export interface IUpdatedBy {
	id: number
	email: string
	username: string
	department_id: number
	firstname: string
	lastname: string
	profile_img_path: string
	status: boolean
	tel: string
	created_by: number
	updated_by: number
	department: IDepartment
}

export interface IDepartment {
	id: number
	name: string
	can_delete: boolean
}

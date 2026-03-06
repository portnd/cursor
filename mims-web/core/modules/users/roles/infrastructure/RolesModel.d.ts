export interface IRoles {
	id: number
	name: string
	can_delete: boolean
}

export interface IRolesDetail {
	access_group: IRolesAccessGroup[]
	id: number
	role: string
}

export interface IRolesAccessGroup {
	menu: IRolesAccessControl[]
	name: string
}

export interface IRolesAccessControl {
	access_control: IRolesAccessDetail[]
	name: string
}

export interface IRolesAccessDetail {
	id: number
	is_check: boolean
	name: string
	access_key: string
}

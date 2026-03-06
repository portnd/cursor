export interface IAccessRoles {
	menu: IAccessRolesMenu[]
	name: string
}

export interface IAccessRolesMenu {
	access_control: IAccessRolesAccessControl[]
	name: string
}

export interface IAccessRolesAccessControl {
	id: number
	is_check: boolean
	name: string
	access_key: string
}

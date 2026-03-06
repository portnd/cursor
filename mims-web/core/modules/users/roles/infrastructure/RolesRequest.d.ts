export interface IRequestRoles {
	access_group: IRequestUpdateRoles[]
}

export interface IRequestRolesAccess {
	menu: IRequestUpdateRoles[]
	name: string
}

export interface IRequestUpdateRoles {
	access_control: IRequestRolesAccessDetail[]
	name: string
}

export interface IRequestRolesAccessDetail {
	access_control_id: number
}

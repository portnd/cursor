export interface ITree {
	id: string
	label: string
	children?: ITree[]
}

export interface Child {
	id: string
	label: string
}

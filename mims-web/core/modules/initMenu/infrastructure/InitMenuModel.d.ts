export interface IInitMenu {
	id: number
	parent_id: number
	title: string
	name: string
	route: string
	icon: string
	is_children: boolean
	children: IInitMenu[] | null
}

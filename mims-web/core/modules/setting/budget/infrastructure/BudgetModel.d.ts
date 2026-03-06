export interface IBudget {
	id: number
	name: string
	can_delete: boolean
	budget: IItemBudget[]
}
export interface IItemBudget {
	id: number
	cost_per_unit: number | null
	method_name: string
	is_show_method?: boolean
}

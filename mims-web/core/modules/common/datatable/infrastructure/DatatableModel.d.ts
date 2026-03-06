export interface IDatatable {
	items: Item[]
	current_page: number
	next_page: number
	previous_page: number
	size_per_page: number
	total_pages: number
	total_items: number
}

export interface Item {
	id: number
	[key: string]: any
	can_delete?: boolean
}

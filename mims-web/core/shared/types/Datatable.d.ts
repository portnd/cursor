import type { Item, ServerOptions, ClickRowArgument } from "vue3-easy-data-table"

export type THeader = {
	text: string
	value: string
	sortable?: boolean
	fixed?: boolean
	width?: number
	align?: string
	type?: string
	fixedEnd?: boolean
}

export type TItem = Item
export type TServerOptions = ServerOptions
export type TClickRowArgument = ClickRowArgument

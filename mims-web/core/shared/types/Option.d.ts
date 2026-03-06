export interface IOption {
	value?: number | string | boolean
	label: string
	disabled?: boolean
	color?: string
	image?: Image
	isSquare?: boolean
}

interface Image {
	src: string
	width?: number
}

interface IFromOptions {
	id: number
	name: string
}

interface IToOptions {
	value: number
	label: string
}

export const toOptions = (options: IFromOptions[] | undefined): IToOptions[] => {
	if (options) {
		return options.map((option: IFromOptions) => ({
			value: option.id,
			label: option.name,
		}))
	} else {
		return []
	}
}

import { v4 as uuidv4 } from "uuid"

export const getUuid = (): string => {
	return uuidv4()
}

export const toNumber = (number: number | undefined, digit: number | null = null): string => {
	if (number !== undefined) {
		if (digit !== null) {
			return new Intl.NumberFormat("en-US", {
				maximumFractionDigits: digit,
				minimumFractionDigits: digit,
				useGrouping: true,
			}).format(number)
		} else {
			return new Intl.NumberFormat("en-US", {
				useGrouping: true,
			}).format(number)
		}
	} else {
		return ""
	}
}

export const generateNumber = (value: number | null) => {
	if (value === null) {
		return "0"
	}

	if (Math.abs(value) < 1000) {
		return value.toString()
	}

	const hasDecimals = value % 1 !== 0
	let formattedValue: number | string = value

	if (hasDecimals) {
		const parts = value.toString().split(".")
		const wholePartFormatted = parseInt(parts[0], 10).toLocaleString()
		formattedValue = `${wholePartFormatted}.${parts[1]}`
	} else {
		formattedValue = value.toLocaleString()
	}

	return formattedValue
}

export const isNumber = (value: number | string | string[]): boolean => {
	const number = Number(value)
	if (!Number.isInteger(number)) {
		return false
	}
	return true
}

export const formatNumberWithSuffix = (value: number) => {
	if (value >= 1000000) {
		return (value / 1000000).toLocaleString(undefined, { minimumFractionDigits: 0, maximumFractionDigits: 1 }) + "M"
	} else if (value >= 1000) {
		return (value / 1000).toLocaleString(undefined, { minimumFractionDigits: 0, maximumFractionDigits: 1 }) + "K"
	} else {
		return value.toLocaleString()
	}
}

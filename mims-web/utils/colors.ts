/**
 * Returns a draw priority for a hex color — higher value = drawn last = appears on top.
 * Red colors (high R, low G) receive the highest priority.
 * Formula: R - G channel difference.
 */
export const getColorDrawPriority = (hex: string): number => {
	const clean = hex.replace(/^#/, "").padEnd(6, "0")
	const bigint = parseInt(clean, 16)
	const r = (bigint >> 16) & 255
	const g = (bigint >> 8) & 255
	return r - g
}

export const convertHexToRGBA = (hex: string, opacity = 1) => {
	console.log(hex)
	hex = hex.replace(/^#/, "")

	// const hexRegex = /^[0-9A-Fa-f]{6}$/
	// if (!hexRegex.test(hex) && hex === "") {
	// 	// throw new Error("Invalid hex color value")
	// 	console.log(hex)
	// 	hex = "000000"
	// }

	hex = hex === "" ? "000000" : hex

	const bigint = parseInt(hex, 16)

	const r = (bigint >> 16) & 255
	const g = (bigint >> 8) & 255
	const b = bigint & 255

	return `rgba(${r}, ${g}, ${b}, ${opacity})`
}

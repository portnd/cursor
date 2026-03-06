export const generateSearchParams = <T>(params: T) => {
	const urlParams = new URLSearchParams()

	for (const key in params) {
		const value = params[key as keyof T]

		if (typeof value === "boolean") {
			urlParams.append(key, value.toString())
		} else if (value && value === false && value !== "") {
			urlParams.append(key, value.toString())
		}
	}

	return urlParams.toString() ? "?" + urlParams.toString() : ""
}

export const textTruncate = (str: string, maxLength = 20) => {
	if (str.length > maxLength) {
		return str.substring(0, maxLength) + "..."
	}

	return str
}

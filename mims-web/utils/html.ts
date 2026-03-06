export const encodeHTML = (str: string): string => {
	return str?.replace(/[&<>"']/g, (match) => {
		switch (match) {
			case "&":
				return "&amp;"
			case "<":
				return "&lt;"
			case ">":
				return "&gt;"
			case '"':
				return "&quot;"
			case "'":
				return "&#39;"
			default:
				return match
		}
	})
}

export const decodeHTML = (str: string): string => {
	const decodedStr = str?.replace(/&(amp|lt|gt|quot|#39);/g, (match, entity) => {
		switch (entity) {
			case "amp":
				return "&"
			case "lt":
				return "<"
			case "gt":
				return ">"
			case "quot":
				return '"'
			case "#39":
				return "'"
			default:
				return match
		}
	})

	return decodedStr
}

export const getBulletedList = (error: string[]) => {
	let html = ""

	if (error.length > 0) {
		html += "<ul>"

		error.forEach((item) => {
			html += `<li>${item}</li>`
		})

		html += "</ul>"
	}

	return html
}

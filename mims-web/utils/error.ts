export const scrollIntoInvalidField = (): void => {
	const invalidFields = document.querySelectorAll(".is-invalid")
	if (invalidFields.length > 0) {
		invalidFields[0].scrollIntoView({
			behavior: "smooth",
			block: "center",
			inline: "center",
		})
	}
}

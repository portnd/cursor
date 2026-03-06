import { IValidate } from "~/core/shared/types/Validate"

export interface IReflectivityRuleList {
	grade_id: number
	left_value: number | null
	right_value: number | null
	left_symbol: string
	right_symbol: string
	left_name: string
	right_name: string
}
export interface IReflectivityRuleType {
	reflectivity_list: IReflectivityRuleList[]
}

export interface IReflectivityRule {
	type: string
	name: string
	leftUnit: string
	rightUnit: string
	white: IReflectivityRuleType
	yellow: IReflectivityRuleType
	id?: number
}

// ข้อมูลตั้งต้น

const useReflectivityRule = (): IReflectivityRule[] => {
	return [
		{
			type: "B",
			name: "Reflex",
			leftUnit: "(ม./กม.)",
			rightUnit: "(ม./กม.)",
			white: {
				reflectivity_list: [
					{
						grade_id: 15,
						left_value: null,
						right_value: null,
						left_symbol: "<=",
						right_symbol: "<",
						left_name: "whiteReflexmin1",
						right_name: "whiteReflexmax1",
					},
					{
						grade_id: 16,
						left_value: null,
						right_value: null,
						left_symbol: "<=",
						right_symbol: "<",
						left_name: "whiteReflexmin2",
						right_name: "whiteReflexmax2",
					},
				],
			},
			yellow: {
				reflectivity_list: [
					{
						grade_id: 15,
						left_value: null,
						right_value: null,
						left_symbol: "<=",
						right_symbol: "<",
						left_name: "yellowReflexmin1",
						right_name: "yellowReflexmax1",
					},
					{
						grade_id: 16,
						left_value: null,
						right_value: null,
						left_symbol: "<=",
						right_symbol: "<",
						left_name: "yellowReflexmin2",
						right_name: "yellowReflexmax2",
					},
				],
			},
		},
	]
}

// ให้ค่าที่กรอกตามเงื่อนไข
const handleFieldReflectivityCondition = (data: IReflectivityRule[], findName = ""): void => {
	// Check if findName is not an empty string, then filter the data by name. Otherwise, set result to data.
	const result = findName !== "" ? data.filter((item: IReflectivityRule) => item.name === findName) : data

	// Loop through each value in result and get its type and conditionList.
	result.forEach((value: IReflectivityRule, i: number) => {
		const type: string = value.type // Type of recipe

		// const surfaceTypes: Array<keyof ISurveyRule> = ["ac", "cc"]
		// surfaceTypes.forEach((st: string) => {
		// Loop through each element in the conditionList of the current value.

		value.white.reflectivity_list.forEach((_: any, j: number) => {
			const beforeElement = result[i].white.reflectivity_list[j] // The element before the current element.
			const afterElement = result[i].white.reflectivity_list[j + 1] // The element after the current element.

			createInputEvent(beforeElement, afterElement, type)
		})

		value.yellow.reflectivity_list.forEach((_: any, j: number) => {
			const beforeElement = result[i].yellow.reflectivity_list[j] // The element before the current element.
			const afterElement = result[i].yellow.reflectivity_list[j + 1] // The element after the current element.

			createInputEvent(beforeElement, afterElement, type)
		})
		// })
	})
}

const createInputEvent = (beforeElement: IReflectivityRuleList, afterElement: IReflectivityRuleList, type: string) => {
	// Check if there is an element after the current element.
	if (afterElement !== undefined) {
		// If the type is A, get the right and left input elements and add event listeners for their input events.
		if (type === "A") {
			const rightInput = document.getElementById(beforeElement.right_name) as HTMLInputElement
			const leftInput = document.getElementById(afterElement.left_name) as HTMLInputElement

			rightInput?.addEventListener("input", (e: any) => {
				leftInput.value = isNaN(Number(e.target.value)) ? "" : e.target.value
				afterElement.left_value = Number(e.target.value)
			})
			leftInput?.addEventListener("input", (e: any) => {
				beforeElement.right_value = Number(e.target.value)
			})
		} else if (type === "B") {
			// If the type is B, get the left and right input elements and add event listeners for their input events.
			const leftInput = document.getElementById(beforeElement.left_name) as HTMLInputElement
			const rightInput = document.getElementById(afterElement.right_name) as HTMLInputElement

			leftInput?.addEventListener("input", (e: any) => {
				afterElement.right_value = Number(e.target.value)
			})
			rightInput?.addEventListener("input", (e: any) => {
				beforeElement.left_value = Number(e.target.value)
			})
		}
	}
}

// การ Validation
const handleFieldReflectivityValidation = (data: IReflectivityRule[], findName = ""): IValidate => {
	const validations: IValidate = {}

	// Check if findName is not an empty string, then filter the data by name. Otherwise, set result to data.
	const result = findName !== "" ? data.filter((item: IReflectivityRule) => item.name === findName) : data

	// Loop through each survey rule in the data array.
	result.forEach((survey: IReflectivityRule) => {
		// Loop through each condition in the current survey rule's conditionList.
		survey.white.reflectivity_list.forEach((reflectRule: IReflectivityRuleList) => {
			// Add a "required" validation for the name field.

			// Add a "required" validation for the left_name field.
			validations[reflectRule.left_name] = "required|max_value:999"
			// Add a "required" validation and a "max_field_value" validation for the right_name field.
			validations[reflectRule.right_name] = `required|max_value:999|max_field_value:@${reflectRule.left_name}`
		})

		// Loop through each condition in the current survey rule's conditionList.
		survey.yellow.reflectivity_list.forEach((reflectRule: IReflectivityRuleList) => {
			// Add a "required" validation for the name field.

			// Add a "required" validation for the left_name field.
			validations[reflectRule.left_name] = "required|max_value:999"
			// Add a "required" validation and a "max_field_value" validation for the right_name field.
			validations[reflectRule.right_name] = `required|max_value:999|max_field_value:@${reflectRule.left_name}`
		})
	})
	validations.name = "required"

	// Return the object containing all the validation rules.
	return validations
}

export { useReflectivityRule, handleFieldReflectivityCondition, handleFieldReflectivityValidation }

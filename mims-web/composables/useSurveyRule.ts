// import {
// 	IConditionGrade,
// 	IConditionGroup,
// 	IConditionInit,
// 	IGrade,
// } from "~/core/modules/initData/infrastructure/data/ConditionGrade"
import { IValidate } from "~/core/shared/types/Validate"

export type LeftSymbol = "<="
export type RightSymbol = "<"

export interface ConditionList {
	grade_id: number
	left_value: number | null
	right_value: number | null
	left_symbol: LeftSymbol
	right_symbol: RightSymbol
	left_name: string
	right_name: string
}

export interface ISurface {
	conditionList: ConditionList[]
}

export interface ISurveyRule {
	type: string
	name: string
	leftUnit: string
	rightUnit: string
	ac: ISurface
	cc: ISurface
	id?: number
}

// หา Grade
// const useCalculateGrade = (conditionId: number, conditionType: string, value: number): IGrade => {
// 	const conditionGrade = useInitData().conditionGrade()

// 	const owner: IConditionGrade = conditionGrade.find((item) => {
// 		return item.id === conditionId
// 	})

// 	const conditions: IConditionInit[] =
// 		owner?.condition_list.find((item: IConditionGroup) => item.condition_type === conditionType)?.surface_type. ?? []

// 	let result = {
// 		color: "#ddd",
// 		id: 0,
// 		name: "ไม่มีข้อมูล",
// 	}
// 	conditions.forEach((item) => {
// 		const expression = `${item.left_value} ${item.left_condition} ${value} && ${value} ${item.right_condition} ${item.right_value}`
// 		if (eval(expression)) {
// 			result = item.grade
// 		}
// 	})

// 	return result
// }

// ข้อมูลตั้งต้น
const useSurveyRule = (rangId: number): ISurveyRule[] => {
	const range25 = 1
	const range100 = 2
	const range1000 = 3

	switch (rangId) {
		case range25:
			return [
				{
					type: "A",
					name: "IRI",
					leftUnit: "(ม./กม.)",
					rightUnit: "(ม./กม.)",
					ac: {
						conditionList: [
							{
								grade_id: 1,
								left_value: null,
								right_value: null,
								left_symbol: "<=",
								right_symbol: "<",
								left_name: "acIRImin1",
								right_name: "acIRImax1",
							},
							{
								grade_id: 2,
								left_value: null,
								right_value: null,
								left_symbol: "<=",
								right_symbol: "<",
								left_name: "acIRImin2",
								right_name: "acIRImax2",
							},
							{
								grade_id: 3,
								left_value: null,
								right_value: null,
								left_symbol: "<=",
								right_symbol: "<",
								left_name: "acIRImin3",
								right_name: "acIRImax3",
							},
							{
								grade_id: 4,
								left_value: null,
								right_value: null,
								left_symbol: "<=",
								right_symbol: "<",
								left_name: "acIRImin4",
								right_name: "acIRImax4",
							},
						],
					},
					cc: {
						conditionList: [
							{
								grade_id: 1,
								left_value: null,
								right_value: null,
								left_symbol: "<=",
								right_symbol: "<",
								left_name: "ccIRImin1",
								right_name: "ccIRImax1",
							},
							{
								grade_id: 2,
								left_value: null,
								right_value: null,
								left_symbol: "<=",
								right_symbol: "<",
								left_name: "ccIRImin2",
								right_name: "ccIRImax2",
							},
							{
								grade_id: 3,
								left_value: null,
								right_value: null,
								left_symbol: "<=",
								right_symbol: "<",
								left_name: "ccIRImin3",
								right_name: "ccIRImax3",
							},
							{
								grade_id: 4,
								left_value: null,
								right_value: null,
								left_symbol: "<=",
								right_symbol: "<",
								left_name: "ccIRImin4",
								right_name: "ccIRImax4",
							},
						],
					},
				},
				{
					type: "B",
					name: "MPD",
					leftUnit: "(มม.)",
					rightUnit: "(มม.)",
					ac: {
						conditionList: [
							{
								grade_id: 8,
								left_value: null,
								right_value: null,
								left_symbol: "<=",
								right_symbol: "<",
								left_name: "acMPDmin1",
								right_name: "acMPDmax1",
							},
							{
								grade_id: 9,
								left_value: null,
								right_value: null,
								left_symbol: "<=",
								right_symbol: "<",
								left_name: "acMPDmin2",
								right_name: "acMPDmax2",
							},
							{
								grade_id: 10,
								left_value: null,
								right_value: null,
								left_symbol: "<=",
								right_symbol: "<",
								left_name: "acMPDmin3",
								right_name: "acMPDmax3",
							},
						],
					},
					cc: {
						conditionList: [
							{
								grade_id: 8,
								left_value: null,
								right_value: null,
								left_symbol: "<=",
								right_symbol: "<",
								left_name: "ccMPDmin1",
								right_name: "ccMPDmax1",
							},
							{
								grade_id: 9,
								left_value: null,
								right_value: null,
								left_symbol: "<=",
								right_symbol: "<",
								left_name: "ccMPDmin2",
								right_name: "ccMPDmax2",
							},
							{
								grade_id: 10,
								left_value: null,
								right_value: null,
								left_symbol: "<=",
								right_symbol: "<",
								left_name: "ccMPDmin3",
								right_name: "ccMPDmax3",
							},
						],
					},
				},
				{
					type: "A",
					name: "RUT",
					leftUnit: "(มม.)",
					rightUnit: "(มม.)",
					ac: {
						conditionList: [
							{
								grade_id: 11,
								left_value: null,
								right_value: null,
								left_symbol: "<=",
								right_symbol: "<",
								left_name: "acRUTmin1",
								right_name: "acRUTmax1",
							},
							{
								grade_id: 12,
								left_value: null,
								right_value: null,
								left_symbol: "<=",
								right_symbol: "<",
								left_name: "acRUTmin2",
								right_name: "acRUTmax2",
							},
							{
								grade_id: 13,
								left_value: null,
								right_value: null,
								left_symbol: "<=",
								right_symbol: "<",
								left_name: "acRUTmin3",
								right_name: "acRUTmax3",
							},
							{
								grade_id: 14,
								left_value: null,
								right_value: null,
								left_symbol: "<=",
								right_symbol: "<",
								left_name: "acRUTmin4",
								right_name: "acRUTmax4",
							},
						],
					},
					cc: {
						conditionList: [
							{
								grade_id: 11,
								left_value: null,
								right_value: null,
								left_symbol: "<=",
								right_symbol: "<",
								left_name: "ccRUTmin1",
								right_name: "ccRUTmax1",
							},
							{
								grade_id: 12,
								left_value: null,
								right_value: null,
								left_symbol: "<=",
								right_symbol: "<",
								left_name: "ccRUTmin2",
								right_name: "ccRUTmax2",
							},
							{
								grade_id: 13,
								left_value: null,
								right_value: null,
								left_symbol: "<=",
								right_symbol: "<",
								left_name: "ccRUTmin3",
								right_name: "ccRUTmax3",
							},
							{
								grade_id: 14,
								left_value: null,
								right_value: null,
								left_symbol: "<=",
								right_symbol: "<",
								left_name: "ccRUTmin4",
								right_name: "ccRUTmax4",
							},
						],
					},
				},
				{
					type: "B",
					name: "IFI",
					leftUnit: "",
					rightUnit: "",
					ac: {
						conditionList: [
							{
								grade_id: 15,
								left_value: null,
								right_value: null,
								left_symbol: "<=",
								right_symbol: "<",
								left_name: "acIFImin1",
								right_name: "acIFImax1",
							},
							{
								grade_id: 16,
								left_value: null,
								right_value: null,
								left_symbol: "<=",
								right_symbol: "<",
								left_name: "acIFImin2",
								right_name: "acIFImax2",
							},
						],
					},
					cc: {
						conditionList: [
							{
								grade_id: 15,
								left_value: null,
								right_value: null,
								left_symbol: "<=",
								right_symbol: "<",
								left_name: "ccIFImin1",
								right_name: "ccIFImax1",
							},
							{
								grade_id: 16,
								left_value: null,
								right_value: null,
								left_symbol: "<=",
								right_symbol: "<",
								left_name: "ccIFImin2",
								right_name: "ccIFImax2",
							},
						],
					},
				},
			]

		case range100:
			return [
				{
					type: "A",
					name: "IRI",
					leftUnit: "(ม./กม.)",
					rightUnit: "(ม./กม.)",
					ac: {
						conditionList: [
							{
								grade_id: 15,
								left_value: null,
								right_value: null,
								left_symbol: "<=",
								right_symbol: "<",
								left_name: "acIRImin1",
								right_name: "acIRImax1",
							},
							{
								grade_id: 16,
								left_value: null,
								right_value: null,
								left_symbol: "<=",
								right_symbol: "<",
								left_name: "acIRImin2",
								right_name: "acIRImax2",
							},
						],
					},
					cc: {
						conditionList: [
							{
								grade_id: 15,
								left_value: null,
								right_value: null,
								left_symbol: "<=",
								right_symbol: "<",
								left_name: "ccIRImin1",
								right_name: "ccIRImax1",
							},
							{
								grade_id: 16,
								left_value: null,
								right_value: null,
								left_symbol: "<=",
								right_symbol: "<",
								left_name: "ccIRImin2",
								right_name: "ccIRImax2",
							},
						],
					},
				},

				{
					type: "A",
					name: "RUT",
					leftUnit: "(มม.)",
					rightUnit: "(มม.)",
					ac: {
						conditionList: [
							{
								grade_id: 15,
								left_value: null,
								right_value: null,
								left_symbol: "<=",
								right_symbol: "<",
								left_name: "acRUTmin1",
								right_name: "acRUTmax1",
							},
							{
								grade_id: 16,
								left_value: null,
								right_value: null,
								left_symbol: "<=",
								right_symbol: "<",
								left_name: "acRUTmin2",
								right_name: "acRUTmax2",
							},
						],
					},
					cc: {
						conditionList: [
							{
								grade_id: 15,
								left_value: null,
								right_value: null,
								left_symbol: "<=",
								right_symbol: "<",
								left_name: "ccRUTmin1",
								right_name: "ccRUTmax1",
							},
							{
								grade_id: 16,
								left_value: null,
								right_value: null,
								left_symbol: "<=",
								right_symbol: "<",
								left_name: "ccRUTmin2",
								right_name: "ccRUTmax2",
							},
						],
					},
				},
				{
					type: "B",
					name: "IFI",
					leftUnit: "",
					rightUnit: "",
					ac: {
						conditionList: [
							{
								grade_id: 15,
								left_value: null,
								right_value: null,
								left_symbol: "<=",
								right_symbol: "<",
								left_name: "acIFImin1",
								right_name: "acIFImax1",
							},
							{
								grade_id: 16,
								left_value: null,
								right_value: null,
								left_symbol: "<=",
								right_symbol: "<",
								left_name: "acIFImin2",
								right_name: "acIFImax2",
							},
						],
					},
					cc: {
						conditionList: [
							{
								grade_id: 15,
								left_value: null,
								right_value: null,
								left_symbol: "<=",
								right_symbol: "<",
								left_name: "ccIFImin1",
								right_name: "ccIFImax1",
							},
							{
								grade_id: 16,
								left_value: null,
								right_value: null,
								left_symbol: "<=",
								right_symbol: "<",
								left_name: "ccIFImin2",
								right_name: "ccIFImax2",
							},
						],
					},
				},
			]

		case range1000:
			return [
				{
					type: "A",
					name: "IRI",
					leftUnit: "(ม./กม.)",
					rightUnit: "(ม./กม.)",
					ac: {
						conditionList: [
							{
								grade_id: 15,
								left_value: null,
								right_value: null,
								left_symbol: "<=",
								right_symbol: "<",
								left_name: "acIRImin1",
								right_name: "acIRImax1",
							},
							{
								grade_id: 16,
								left_value: null,
								right_value: null,
								left_symbol: "<=",
								right_symbol: "<",
								left_name: "acIRImin2",
								right_name: "acIRImax2",
							},
						],
					},
					cc: {
						conditionList: [
							{
								grade_id: 15,
								left_value: null,
								right_value: null,
								left_symbol: "<=",
								right_symbol: "<",
								left_name: "ccIRImin1",
								right_name: "ccIRImax1",
							},
							{
								grade_id: 16,
								left_value: null,
								right_value: null,
								left_symbol: "<=",
								right_symbol: "<",
								left_name: "ccIRImin2",
								right_name: "ccIRImax2",
							},
						],
					},
				},
			]
		default:
			return []
	}
}

// ให้ค่าที่กรอกตามเงื่อนไข
const handleFieldCondition = (data: ISurveyRule[], findName = ""): void => {
	// Check if findName is not an empty string, then filter the data by name. Otherwise, set result to data.
	const result = findName !== "" ? data.filter((item: ISurveyRule) => item.name === findName) : data

	// Loop through each value in result and get its type and conditionList.
	result.forEach((value: any, i: number) => {
		const type: string = value.type // Type of recipe

		const surfaceTypes: Array<keyof ISurveyRule> = ["ac", "cc"]

		surfaceTypes.forEach((st: string) => {
			// Loop through each element in the conditionList of the current value.
			value[st].conditionList.forEach((_: any, j: number) => {
				const beforeElement = (result[i] as any)[st].conditionList[j] // The element before the current element.
				const afterElement = (result[i] as any)[st].conditionList[j + 1] // The element after the current element.

				// Check if there is an element after the current element.
				if (afterElement !== undefined) {
					// If the type is A, get the right and left input elements and add event listeners for their input events.
					if (type === "A") {
						const rightInput = document.getElementById(beforeElement.right_name) as HTMLInputElement
						const leftInput = document.getElementById(afterElement.left_name) as HTMLInputElement

						rightInput?.addEventListener("input", (e: any) => {
							console.log("type A rightInput:", e.target.value)
							leftInput.value = isNaN(Number(e.target.value)) ? "" : e.target?.value
							afterElement.left_value = Number(e.target.value)
						})
						leftInput?.addEventListener("input", (e: any) => {
							console.log("type A leftInput:", e.target.value)
							beforeElement.right_value = Number(e.target.value)
						})
					} else if (type === "B") {
						// If the type is B, get the left and right input elements and add event listeners for their input events.
						const leftInput = document.getElementById(beforeElement.left_name) as HTMLInputElement
						const rightInput = document.getElementById(afterElement.right_name) as HTMLInputElement

						leftInput?.addEventListener("input", (e: any) => {
							console.log("type B leftInput:", e.target.value)
							afterElement.right_value = Number(e.target.value)
						})
						rightInput?.addEventListener("input", (e: any) => {
							console.log("type B rightInput:", e.target.value)
							beforeElement.left_value = Number(e.target.value)
						})
					}
				}
			})
		})
	})
}

// การ Validation
const handleFieldValidation = (data: ISurveyRule[], findName = ""): IValidate => {
	const validations: IValidate = {}

	// Check if findName is not an empty string, then filter the data by name. Otherwise, set result to data.
	const result = findName !== "" ? data.filter((item: ISurveyRule) => item.name === findName) : data
	// Loop through each survey rule in the data array.
	result.forEach((survey: ISurveyRule) => {
		// Loop through each condition in the current survey rule's conditionList.
		survey.ac.conditionList.forEach((condition: ConditionList) => {
			// Add a "required" validation for the name field.

			// Add a "required" validation for the left_name field.
			validations[condition.left_name] = "required|max_value:100"
			// Add a "required" validation and a "max_field_value" validation for the right_name field.
			validations[condition.right_name] = `required|max_value:100|max_field_value:@${condition.left_name}`
		})

		// Loop through each condition in the current survey rule's conditionList.
		survey.cc.conditionList.forEach((condition: ConditionList) => {
			// Add a "required" validation for the name field.

			// Add a "required" validation for the left_name field.
			validations[condition.left_name] = "required|max_value:100"
			// Add a "required" validation and a "max_field_value" validation for the right_name field.
			validations[condition.right_name] = `required|max_value:100|max_field_value:@${condition.left_name}`
		})
	})
	validations.name = "required"
	// Return the object containing all the validation rules.
	return validations
}

export { useSurveyRule, handleFieldCondition, handleFieldValidation, useCalculateGrade }

import { Form, Field, ErrorMessage, FieldArray, defineRule, configure } from "vee-validate"
import AllRules from "@vee-validate/rules"
import { localize } from "@vee-validate/i18n"
import th from "@vee-validate/i18n/dist/locale/th.json"

export default defineNuxtPlugin((nuxtApp) => {
	nuxtApp.vueApp.component("VForm", Form)
	nuxtApp.vueApp.component("VField", Field)
	nuxtApp.vueApp.component("VErrorMessage", ErrorMessage)
	nuxtApp.vueApp.component("VFieldArray", FieldArray)

	Object.keys(AllRules).forEach((rule) => {
		defineRule(rule, AllRules[rule])
	})

	// Activate the locale
	localize({ th })

	configure({
		generateMessage: localize("th", {
			messages: {
				alpha: "ต้องเป็นตัวอักษรเท่านั้น",
				alpha_dash: "สามารถมีตัวอักษร ตัวเลข เครื่องหมายขีดกลาง (-) และเครื่องหมายขีดล่าง (_)",
				alpha_num: "ต้องเป็นตัวอักษร และตัวเลขเท่านั้น",
				alpha_spaces: "ต้องเป็นตัวอักษร และช่องว่างเท่านั้น",
				between: "ต้องเป็นค่าระหว่าง 0:{min} และ 1:{max}",
				confirmed: "การยืนยันข้อมูลของ ไม่ตรงกัน",
				digits: "ต้องเป็นตัวเลขจำนวน 0:{length} หลักเท่านั้น",
				dimensions: "ต้องมีขนาด 0:{width}x{height} px",
				email: "ต้องเป็นรูปแบบอีเมล",
				not_one_of: "ต้องเป็นค่าที่กำหนดเท่านั้น",
				ext: "สกุลไฟล์ไม่ถูกต้อง",
				image: "ต้องเป็นรูปภาพเท่านั้น",
				one_of: "ต้องเป็นค่าที่กำหนดเท่านั้น",
				integer: "ต้องเป็นเลขจำนวนเต็ม",
				length: "ต้องมีความยาว 0:{length}",
				max: "ต้องมีความยาวไม่เกิน 0:{length} ตัวอักษร",
				max_value: "ต้องมีค่าไม่เกิน 0:{max}",
				mimes: "ประเภทไฟล์ไม่ถูกต้อง",
				min: "ต้องมีความยาวอย่างน้อย 0:{length} ตัวอักษร",
				min_value: "ต้องมีค่าตั้งแต่ 0:{min} ขึ้นไป",
				numeric: "ต้องเป็นตัวเลขเท่านั้น",
				regex: "รูปแบบ ไม่ถูกต้อง",
				required: "โปรดระบุ",
				required_if: "โปรดระบุ",
				size: "ต้องมีขนาดไฟล์ไม่เกิน 0:{size}KB",
			},
		}),
	})

	// Rules
	defineRule("max_field_value", (value: number, data: number) => {
		if (Number(value) === Number(data) || value < Number(data)) {
			return `ต้องมีค่ามากกว่า ${data}`
		}
		return true
	})

	defineRule("km", (value: any) => {
		const splitText = value?.split("+")
		if (
			splitText?.length !== 2 ||
			splitText[0]?.length === 0 ||
			typeof Number(splitText[0]) !== "number" ||
			isNaN(Number(splitText[0])) ||
			typeof Number(splitText[1]) !== "number" ||
			isNaN(splitText[1]) ||
			splitText[1]?.length !== 3
		) {
			return "รูปแบบ กม. ไม่ถูกต้อง"
		} else {
			return true
		}
	})

	defineRule("username", (value: any) => {
		const input = value
		let message: string | string[] = []

		const onlyEnglishAndNumbers = /^[A-Za-z0-9_\-]+$/.test(input)

		if (!onlyEnglishAndNumbers && input !== "") {
			message.push("สามารถประกอบด้วยตัวอักษร A-Z, a-z และตัวเลข 0-9 ตามเงื่อนไขเท่านั้น")
		}

		const result = message

		return getBulletedList(result)
	})

	defineRule("password", (value: any) => {
		let message: string[] = []

		const input = value

		const onlyEnglishAndNumbers = /^[A-Za-z0-9_\-]+$/.test(input)
		const hasUpper = /[A-Z]/.test(input)
		const hasLower = /[a-z]/.test(input)
		const hasNumber = /[0-9]/.test(input)
		const inputLength = input?.length

		if (inputLength < 8) {
			message.push("ต้องมีจำนวน 8 ตัวอักษรขึ้นไป")
		}

		if (!onlyEnglishAndNumbers) {
			message.push("ต้องประกอบด้วยตัวอักษร A-Z, a-z และตัวเลข 0-9 ตามเงื่อนไขเท่านั้น")
		}

		if (!hasUpper) {
			message.push("ต้องประกอบด้วยตัวอักษร A-Z ตัวพิมพ์ใหญ่")
		}

		if (!hasLower) {
			message.push("ต้องประกอบด้วยตัวอักษร  a-z ตัวพิมพ์เล็ก")
		}

		if (!hasNumber) {
			message.push("ต้องประกอบด้วยตัวเลข 0-9 อย่างน้อย 1 ตัว")
		}

		const result = message

		return getBulletedList(result)
	})
})

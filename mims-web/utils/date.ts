const monthNames = [
	"มกราคม",
	"กุมภาพันธ์",
	"มีนาคม",
	"เมษายน",
	"พฤษภาคม",
	"มิถุนายน",
	"กรกฎาคม",
	"สิงหาคม",
	"กันยายน",
	"ตุลาคม",
	"พฤศจิกายน",
	"ธันวาคม",
]

const monthShortNames = [
	"ม.ค.",
	"ก.พ.",
	"มี.ค.",
	"เม.ย.",
	"พ.ค.",
	"มิ.ย.",
	"ก.ค.",
	"ส.ค.",
	"ก.ย.",
	"ต.ค.",
	"พ.ย.",
	"ธ.ค.",
]

const dayOfWeekNames = ["วันอาทิตย์", "วันจันทร์", "วันอังคาร", "วันพุธ", "วันพฤหัสบดี", "วันศุกร์", "วันเสาร์"]
const dayOfWeekShortNames = ["อา.", "จ.", "อ.", "พ.", "พฤ.", "ศ.", "ส."]

const twoDigitPad = (num: number) => (num < 10 ? "0" + num : num)

export const buddhistFormatDate = (dateString: string | Date | undefined, patternStr = "yyyy-mm-dd") => {
	const isBuddhist = true
	return formatDate(dateString, patternStr, isBuddhist)
}

export const formatDate = (dateString: string | Date | undefined, patternStr = "yyyy-mm-dd", isBuddhist = false) => {
	if (typeof dateString === "undefined" || dateString === null) {
		return ""
	}

	if (dateString === "") {
		return ""
	}

	const date: Date = dateString instanceof Date ? dateString : new Date(dateString)

	if (isNaN(date.getTime())) {
		return ""
	}

	const day = date.getDate()
	const month = date.getMonth()
	const year = isBuddhist ? date.getFullYear() + 543 : date.getFullYear()
	const hour = date.getHours()
	const minute = date.getMinutes()
	const second = date.getSeconds()
	const miliseconds = date.getMilliseconds()
	const h = hour % 12
	const hh = twoDigitPad(h)
	const HH = twoDigitPad(hour)
	const ii = twoDigitPad(minute)
	const ss = twoDigitPad(second)
	const aaa = hour < 12 ? "AM" : "PM"
	const EEEE = dayOfWeekNames[date.getDay()]
	const EEE = dayOfWeekShortNames[date.getDay()]
	const dd = twoDigitPad(day)
	const m = month + 1
	const mm = twoDigitPad(m)
	const mmmm = monthNames[month]
	const mmm = monthShortNames[month]
	const yyyy = year + ""
	const yy = yyyy.substring(yyyy.length - 2)

	patternStr = patternStr
		.replace(/hh/g, `${hh}`)
		.replace(/h/g, `${h}`)
		.replace(/HH/g, `${HH}`)
		.replace(/H/g, `${hour}`)
		.replace(/ii/g, `${ii}`)
		.replace(/i/g, `${minute}`)
		.replace(/ss/g, `${ss}`)
		.replace(/s/g, `${second}`)
		.replace(/S/g, `${miliseconds}`)
		.replace(/dd/g, `${dd}`)
		.replace(/d/g, `${day}`)
		.replace(/EEEE/g, `${EEEE}`)
		.replace(/EEE/g, `${EEE}`)
		.replace(/yyyy/g, `${yyyy}`)
		.replace(/yy/g, `${yy}`)
		.replace(/aaa/g, `${aaa}`)
	if (patternStr.includes("mmm")) {
		patternStr = patternStr.replace(/mmmm/g, `${mmmm}`).replace(/mmm/g, `${mmm}`)
	} else {
		patternStr = patternStr.replace(/mm/g, `${mm}`).replace(/m/g, `${m}`)
	}

	return patternStr
}

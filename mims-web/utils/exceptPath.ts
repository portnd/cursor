// ฟังก์ชั่นตรวจสอบ Path ที่ละเว้น
export const exceptPath = (path: string, excepts: string[]): boolean => {
	for (let i = 0; i < excepts.length; i++) {
		if (path.includes(excepts[i])) {
			return true
		}
	}
	return false
}

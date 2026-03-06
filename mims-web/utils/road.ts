export const convertMeterToKm = (number: number): string => {
	let prefix = 0
	let meter = 0

	const zeroPad = (num: number, places: number) => String(num).padStart(places, "0")
	if (number > 0) {
		number = Math.round(number)
		if (number / 1000.0 >= 1) {
			prefix = Math.floor(number / 1000.0)
			meter = number % 1000
		} else {
			prefix = 0
			meter = number % 1000
		}

		return prefix + "+" + zeroPad(meter, 3)
	} else {
		return "0+000"
	}
}

export const convertStringToKm = (string: string): number => {
	const a = Number(string?.split("+")[0])
	const b = Number(string?.split("+")[1])
	const c = a * 1000 + b
	return c
}

export const calculateDistance = (start: number, end: number): string => {
	const floatStart = start / 1000
	const floatEnd = end / 1000
	const diff = Math.abs(floatEnd - floatStart)
	const diffString = diff.toString()
	const splitString: Array<string> = diffString.split(".")

	const zeroPad = (num: String, places: number) => num.padEnd(places, "0")

	let result = ""
	if (splitString.length === 1) {
		result = diff + ".000"
	} else if (splitString[1].length < 3) {
		result = splitString[0] + "." + zeroPad(splitString[1], 3)
	} else {
		result = diff.toFixed(3)
	}

	return result
}

export const getRoadTypeIcon = (roadTypeIconId: number): string => {
	const refRoadTypeIcons = useInitData().refRoadTypeIcon()

	let icon = ""
	if (refRoadTypeIcons) {
		const refRoadTypeIcon = refRoadTypeIcons.find((obj) => obj.id === roadTypeIconId)
		icon = refRoadTypeIcon?.icon as string
	}

	return icon
}

export const getItemRowRevision = (type: string) => {
	let html = ""
	switch (type) {
		case "A": // add
			html = `<span class="badge badge-success rounded-1 px-3 fs-8">เพิ่มใหม่</span>`
			break
		case "M": // edit
			html = `<span class="badge badge-warning rounded-1 px-3 fs-8">แก้ไข</span>`
			break
		case "D": // delete
			html = `<span class="badge badge-danger rounded-1 px-3 fs-8">ลบ</span>`
			break
	}
	return html
}

interface LatLng {
	lat: number
	lon: number
}

export const getLatLong = (geom: string): LatLng => {
	const latLng: LatLng = {
		lat: 0,
		lon: 0,
	}

	const location = geom?.match(/\(([^)]+)\)/)?.[1]?.split(" ")
	if (location) {
		latLng.lat = parseFloat(location[1])
		latLng.lon = parseFloat(location[0])
	}

	return latLng
}

export const parseLineString = (lineString: string): LatLng[] => {
	// ลบคำว่า "LINESTRING(" ที่จุดเริ่มต้น และ ")" ที่จุดสิ้นสุด
	const coordsString = lineString.replace("LINESTRING(", "").replace(")", "")

	// แยกพิกัดด้วย ","
	const coordsArray = coordsString.split(",")

	// แปลงพิกัดให้เป็น array ของ objects ที่มี key `lat` และ `lon`
	const coordinates: LatLng[] = coordsArray.map((coord) => {
		const [lon, lat] = coord.trim().split(" ").map(Number)
		return { lat, lon }
	})

	return coordinates
}

export const getLatLongObject = (geom: Object): LatLng => {
	const latLng: LatLng = {
		lat: 0,
		lon: 0,
	}

	if ((geom as any).type.toLowerCase() === "point") {
		latLng.lat = (geom as any).coordinates[1]
		latLng.lon = (geom as any).coordinates[0]
		return latLng
	}

	const location = (geom as any).coordinates[0]

	if (location) {
		latLng.lat = parseFloat(location[1])
		latLng.lon = parseFloat(location[0])
	}

	return latLng
}

// ฟังก์ชันคำนวณระยะทางระหว่างสองพิกัด
export function haversineDistance(coord1: LatLng, coord2: LatLng): number {
	const R = 6371000 // รัศมีของโลกในหน่วยเมตร
	const dLat = (coord2.lat - coord1.lat) * (Math.PI / 180)
	const dLon = (coord2.lon - coord1.lon) * (Math.PI / 180)
	const lat1 = coord1.lat * (Math.PI / 180)
	const lat2 = coord2.lat * (Math.PI / 180)

	const a =
		Math.sin(dLat / 2) * Math.sin(dLat / 2) + Math.sin(dLon / 2) * Math.sin(dLon / 2) * Math.cos(lat1) * Math.cos(lat2)
	const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a))

	return R * c
}

// ฟังก์ชันหาไดนามิกพิกัดที่ระยะทางที่กำหนด
export function getLatLngByDistance(coords: LatLng[], distance: number): LatLng | null {
	let accumulatedDistance = 0

	for (let i = 0; i < coords.length - 1; i++) {
		const segmentDistance = haversineDistance(coords[i], coords[i + 1])
		if (accumulatedDistance + segmentDistance >= distance) {
			const remainingDistance = distance - accumulatedDistance
			const ratio = remainingDistance / segmentDistance
			const lat = coords[i].lat + (coords[i + 1].lat - coords[i].lat) * ratio
			const lon = coords[i].lon + (coords[i + 1].lon - coords[i].lon) * ratio
			return { lat, lon }
		}
		accumulatedDistance += segmentDistance
	}

	return null
}

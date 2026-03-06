package helpers

import (
	"fmt"
	"strings"
	"time"
)

func TimeSlash(time time.Time) string {
	// format for 31/12/66
	return fmt.Sprintf("%02d/%02d/%02d", time.Day(), int(time.Month()), (time.Year()+543)%100)
}

func TimeThai(time time.Time) string {
	// 31 มกราคม 2566
	currentTime := time.AddDate(543, 0, 0).Format("2 Jan 2006")
	month := time.Format("Jan")
	monthThai := map[string]string{
		"Jan": "ม.ค.",
		"Feb": "ก.พ.",
		"Mar": "มี.ค.",
		"Apr": "เม.ย.",
		"May": "พ.ค.",
		"Jun": "มิ.ย.",
		"Jul": "ก.ค.",
		"Aug": "ส.ค.",
		"Sep": "ก.ย.",
		"Oct": "ต.ค.",
		"Nov": "พ.ย.",
		"Dec": "ธ.ค.",
	}
	currentTime = strings.Replace(currentTime, month, monthThai[month], 1)

	return currentTime
}

func TimeThai2DigitYear(time time.Time) string {
	// 31 มกราคม 2566
	currentTime := time.AddDate(543, 0, 0).Format("2 Jan 06")
	month := time.Format("Jan")
	monthThai := map[string]string{
		"Jan": "ม.ค.",
		"Feb": "ก.พ.",
		"Mar": "มี.ค.",
		"Apr": "เม.ย.",
		"May": "พ.ค.",
		"Jun": "มิ.ย.",
		"Jul": "ก.ค.",
		"Aug": "ส.ค.",
		"Sep": "ก.ย.",
		"Oct": "ต.ค.",
		"Nov": "พ.ย.",
		"Dec": "ธ.ค.",
	}
	currentTime = strings.Replace(currentTime, month, monthThai[month], 1)

	return currentTime
}

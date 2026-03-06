package helpers

import (
	"fmt"
	"time"
)

func ConvertToThaiCalendar(t time.Time) string {
	thaiYear := t.Year() + 543
	thaiMonth := ""

	switch t.Month() {
	case time.January:
		thaiMonth = "ม.ค."
	case time.February:
		thaiMonth = "ก.พ."
	case time.March:
		thaiMonth = "มี.ค."
	case time.April:
		thaiMonth = "เม.ย."
	case time.May:
		thaiMonth = "พ.ค."
	case time.June:
		thaiMonth = "มิ.ย."
	case time.July:
		thaiMonth = "ก.ค."
	case time.August:
		thaiMonth = "ส.ค."
	case time.September:
		thaiMonth = "ก.ย."
	case time.October:
		thaiMonth = "ต.ค."
	case time.November:
		thaiMonth = "พ.ย."
	case time.December:
		thaiMonth = "ธ.ค."
	}

	return fmt.Sprintf("%s %02d", thaiMonth, thaiYear%100)
}

func ConvertToThaiFullCalendar(t time.Time) string {
	thaiYear := t.Year() + 543
	thaiMonth := ""

	switch t.Month() {
	case time.January:
		thaiMonth = "ม.ค."
	case time.February:
		thaiMonth = "ก.พ."
	case time.March:
		thaiMonth = "มี.ค."
	case time.April:
		thaiMonth = "เม.ย."
	case time.May:
		thaiMonth = "พ.ค."
	case time.June:
		thaiMonth = "มิ.ย."
	case time.July:
		thaiMonth = "ก.ค."
	case time.August:
		thaiMonth = "ส.ค."
	case time.September:
		thaiMonth = "ก.ย."
	case time.October:
		thaiMonth = "ต.ค."
	case time.November:
		thaiMonth = "พ.ย."
	case time.December:
		thaiMonth = "ธ.ค."
	}

	return fmt.Sprintf("%02d %s %02d %02d:%02d:%02d", t.Day(), thaiMonth, thaiYear, t.Hour(), t.Minute(), t.Second())
}

func ConvertToThaiFullCalendarNoTime(t time.Time) string {
	thaiYear := t.Year() + 543
	thaiMonth := ""

	switch t.Month() {
	case time.January:
		thaiMonth = "ม.ค."
	case time.February:
		thaiMonth = "ก.พ."
	case time.March:
		thaiMonth = "มี.ค."
	case time.April:
		thaiMonth = "เม.ย."
	case time.May:
		thaiMonth = "พ.ค."
	case time.June:
		thaiMonth = "มิ.ย."
	case time.July:
		thaiMonth = "ก.ค."
	case time.August:
		thaiMonth = "ส.ค."
	case time.September:
		thaiMonth = "ก.ย."
	case time.October:
		thaiMonth = "ต.ค."
	case time.November:
		thaiMonth = "พ.ย."
	case time.December:
		thaiMonth = "ธ.ค."
	}

	return fmt.Sprintf("%02d %s %02d", t.Day(), thaiMonth, thaiYear)
}

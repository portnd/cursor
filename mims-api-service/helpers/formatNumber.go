package helpers

import (
	"strconv"
	"strings"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func FormatNumber(num int) string {
	numStr := strconv.Itoa(num)
	formattedNum := ""

	for i, digit := range numStr {
		if (len(numStr)-i)%3 == 0 && i != 0 {
			formattedNum += ","
		}
		formattedNum += string(digit)
	}

	return formattedNum
}

func FormatFloatNumber(num float64) string {
	p := message.NewPrinter(language.English)
	text := p.Sprintf("%.0f", num)
	return text
}

func AddCommasToNumber(num string) string {
	if strings.Contains(num, ".") {
		integer := strings.Split(num, ".")[0]
		decimal := strings.Split(num, ".")[1]
		length := len(integer)
		numWithCommas := ""
		for i, char := range integer {
			numWithCommas += string(char)
			if (length-i-1)%3 == 0 && i != length-1 {
				numWithCommas += ","
			}
		}
		if decimal == "00" {
			return numWithCommas
		} else {
			return numWithCommas + "." + decimal
		}
	} else {
		length := len(num)
		numWithCommas := ""
		for i, char := range num {
			numWithCommas += string(char)
			if (length-i-1)%3 == 0 && i != length-1 {
				numWithCommas += ","
			}
		}
		return numWithCommas
	}

}

func FormatNumberFloat(num float64) string {
	numStr := strconv.FormatFloat(num, 'f', 3, 64)
	formattedNum := ""

	intPart, decPart := splitNumber(numStr)

	intPart = formatIntegerPart(intPart)

	if num < 0 {
		formattedNum += "-"
	}
	// if decPart == "000" {
	// 	formattedNum += intPart
	// } else {
	formattedNum += intPart + "." + decPart
	//}

	return formattedNum
}

func splitNumber(numStr string) (string, string) {
	parts := strings.Split(numStr, ".")
	intPart := parts[0]
	decPart := ""
	if len(parts) > 1 {
		decPart = parts[1]
	}
	return intPart, decPart
}

func formatIntegerPart(intPart string) string {
	formattedIntPart := ""
	for i, digit := range reverseString(intPart) {
		if i != 0 && i%3 == 0 {
			formattedIntPart += ","
		}
		formattedIntPart += string(digit)
	}
	return reverseString(formattedIntPart)
}

func reverseString(str string) string {
	runes := []rune(str)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func FormatNumberFloatWithDigiInput(num float64, digi int) string {
	numStr := strconv.FormatFloat(num, 'f', digi, 64)
	formattedNum := ""

	intPart, decPart := splitNumber(numStr)

	intPart = formatIntegerPart(intPart)

	if num < 0 {
		formattedNum += "-"
	}

	formattedNum += intPart + "." + decPart

	return formattedNum
}

package main

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	decimalString string = "Decimal: "
	hexString     string = "Hex: "
	binaryString  string = "Binary: \n"
)

func binaryDisplayString(num int64) []string {
	var rows [4][4][]string
	var j, k, w int
	for i := 63; i > -1; i-- {
		value := (int(((1 << i) & num) >> i))
		value = max(value, -value)
		valueString := strconv.Itoa(value)
		if rows[j][k] == nil { //nolint:gosec // we are doing some math to ensure we stay in bounds
			rows[j][k] = make([]string, 4)
		}
		rows[j][k][w] = valueString
		w = (w + 1) % 4
		if w == 0 {
			k = (k + 1) % 4
			if k == 0 {
				j = (j + 1) % 4
			}
		}
	}
	display := []string{binaryString}
	for i := range 4 {
		displayValue := strconv.Itoa((64 - (16 * i)))
		var inner []string
		for j := range 4 {
			inner = append(inner, strings.Join(rows[i][j], " "))
		}
		innerString := strings.Join(inner, "  ")

		display = append(display, displayValue+": "+innerString+"\n")
	}
	return display
}

func decimalDisplayString(num int64) string {
	return decimalString + strconv.Itoa(int(num)) + "\n"
}

func hexDisplayString(num int64) string {
	return hexString + fmt.Sprintf("%x\n", num)
}

func displayString(num int64) []string {
	return append([]string{ascii(num), decimalDisplayString(num), hexDisplayString(num)}, binaryDisplayString(num)...)
}

func ascii(num int64) string {
	switch num {
	case 12:
		return "ascii: \n"
	case 7:
		return "ascii: \n"
	case 10:
		return "ascii: \\n\n"
	case 11:
		return "ascii: \\r\n"
	default:
		return "ascii: " + string(rune(num)) + "\n"
	}
}

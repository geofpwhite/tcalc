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

func binaryDisplayString(num int64) string {
	var rows [4][4][]string
	var j, k, w int
	for i := 63; i > -1; i-- {
		value := strconv.Itoa(int(((1 << i) & num) >> i))
		if rows[j][k] == nil {
			rows[j][k] = make([]string, 4)
		}
		rows[j][k][w] = value
		w = (w + 1) % 4
		if w == 0 {
			k = (k + 1) % 4
			if k == 0 {
				j = (j + 1) % 4
			}
		}
	}
	display := binaryString
	for i := range 4 {
		displayValue := strconv.Itoa((64 - (16 * i)))
		var inner []string
		for j := range 4 {
			inner = append(inner, strings.Join(rows[i][j], " "))
		}
		innerString := strings.Join(inner, "  ")

		display += displayValue + ": " + innerString + "\n"
	}
	return display
}

func decimalDisplayString(num int64) string {
	return decimalString + strconv.Itoa(int(num)) + "\n"
}

func hexDisplayString(num int64) string {
	return hexString + fmt.Sprintf("%x\n", num)
}

func displayString(num int64) string {
	return ascii(num) + decimalDisplayString(num) + hexDisplayString(num) + binaryDisplayString(num)
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

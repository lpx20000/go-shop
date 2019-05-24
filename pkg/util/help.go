package util

import (
	"fmt"
	"strconv"
)

func Multiplication(f1 float64) (f float64) {
	f, _ = FormatNumber(FormatPreciseNumber(f1))
	return
}

func FormatNumber(s string) (f float64, err error) {
	return strconv.ParseFloat(s, 64)
}

func FormatPreciseNumber(f float64) string {
	return fmt.Sprintf("%.2f", f)
}

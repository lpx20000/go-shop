package util

import (
	"crypto/md5"
	"encoding/hex"
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

func Md5V(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

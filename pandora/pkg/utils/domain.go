package utils

import "strings"

// CheckStringAlphabet
//
//	@Description: 检查字符串是否为域名
//	@param str
//	@return bool
func CheckStringAlphabet(str string) bool {
	L := len(str)
	if L < 1 {
		return false
	}
	charVariable := str[L-1:]
	// ipv6
	if strings.Contains(str, ":") {
		return false
	}
	// ipv4
	if charVariable >= "0" && charVariable <= "9" {
		return false
	}

	return true
}

// Reverse
//
//	@Description: 反转域名字符串
//	@param s
//	@return string
func Reverse(s string) string {
	if !CheckStringAlphabet(s) {
		return s
	}
	r := []rune(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}

	return string(r)
}

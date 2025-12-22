package utils

import "strings"

/**
 * 去除字符串中的 NULL 字符
 */
func CleanString(s string) string {
	return strings.ReplaceAll(s, "\x00", "")
}

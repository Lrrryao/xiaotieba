package util

import (
	"strings"
)

// KeyMatch2 是一个用于模式匹配的函数，用于 Casbin 的策略匹配。
// 它支持两种特殊字符的匹配：* 表示匹配任意字符，? 表示匹配单个字符。
// 例如，"/book/:id" 可以匹配 "/book/1" 和 "/book/2"。
func KeyMatch2(key1, key2 string) bool {
	key2 = strings.ReplaceAll(key2, "*", "_casbin_match_any_")
	key2 = strings.ReplaceAll(key2, "?", "_casbin_match_one_")
	return key1 == key2
}

package tools

// 校验手机号是否正确

import (
	"regexp"
)

// 校验手机号是否正确
func IsPhone(phone string) bool {
	reg := regexp.MustCompile(`^1[3456789]\d{9}$`)
	return reg.MatchString(phone)
}

package tools

import (
	"crypto/md5"
	"fmt"
)

// 密码哈希

func PasswordHash(password string, salt string) string {
	sum := md5.Sum([]byte(password + salt))
	return fmt.Sprintf("%x", sum)
}

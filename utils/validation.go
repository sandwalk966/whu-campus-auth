package utils

import "regexp"

var (
	emailRegex   = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	phoneRegex   = regexp.MustCompile(`^1[3-9]\d{9}$`)
)

// IsValidEmail 验证邮箱格式
func IsValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}

// IsValidPhone 验证手机号格式
func IsValidPhone(phone string) bool {
	return phoneRegex.MatchString(phone)
}

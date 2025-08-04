package security

import (
	"errors"
	"regexp"
	"strings"
)

// PasswordValidator 密码验证器
type PasswordValidator struct {
	minLength    int
	requireUpper bool
	requireLower bool
	requireDigit bool
	requireSpecial bool
	weakPasswords []string
}

// NewPasswordValidator 创建密码验证器
func NewPasswordValidator() *PasswordValidator {
	return &PasswordValidator{
		minLength:    8,
		requireUpper: true,
		requireLower: true,
		requireDigit: true,
		requireSpecial: true,
		weakPasswords: []string{
			"12345678", "password", "qwerty123", "admin123",
			"123456789", "password123", "qwertyuiop", "admin1234",
			"abcd1234", "1234abcd", "password1", "123qwe123",
			"aaaaaa", "111111", "000000", "123123",
		},
	}
}

// ValidatePassword 验证密码强度
func (v *PasswordValidator) ValidatePassword(password string) error {
	// 检查长度
	if len(password) < v.minLength {
		return errors.New("密码长度不能少于8位")
	}
	
	// 检查最大长度
	if len(password) > 128 {
		return errors.New("密码长度不能超过128位")
	}
	
	var (
		hasUpper   = false
		hasLower   = false
		hasDigit   = false
		hasSpecial = false
	)
	
	// 检查字符类型
	for _, char := range password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= '0' && char <= '9':
			hasDigit = true
		case strings.ContainsRune("@$!%*?&#+=-_()[]{}|;:,.<>", char):
			hasSpecial = true
		}
	}
	
	// 检查必需的字符类型
	if v.requireUpper && !hasUpper {
		return errors.New("密码必须包含至少一个大写字母")
	}
	
	if v.requireLower && !hasLower {
		return errors.New("密码必须包含至少一个小写字母")
	}
	
	if v.requireDigit && !hasDigit {
		return errors.New("密码必须包含至少一个数字")
	}
	
	if v.requireSpecial && !hasSpecial {
		return errors.New("密码必须包含至少一个特殊字符(@$!%*?&#+=-_()[]{}|;:,.<>)")
	}
	
	// 检查是否为弱密码
	lowerPassword := strings.ToLower(password)
	for _, weak := range v.weakPasswords {
		if lowerPassword == weak {
			return errors.New("密码过于简单，请使用更复杂的密码")
		}
	}
	
	// 检查重复字符
	if v.hasRepeatingChars(password, 3) {
		return errors.New("密码不能包含3个或更多连续相同的字符")
	}
	
	// 检查连续字符
	if v.hasSequentialChars(password, 4) {
		return errors.New("密码不能包含4个或更多连续的字符序列")
	}
	
	return nil
}

// hasRepeatingChars 检查是否有重复字符
func (v *PasswordValidator) hasRepeatingChars(password string, maxRepeat int) bool {
	if len(password) < maxRepeat {
		return false
	}
	
	count := 1
	for i := 1; i < len(password); i++ {
		if password[i] == password[i-1] {
			count++
			if count >= maxRepeat {
				return true
			}
		} else {
			count = 1
		}
	}
	
	return false
}

// hasSequentialChars 检查是否有连续字符
func (v *PasswordValidator) hasSequentialChars(password string, maxSequential int) bool {
	if len(password) < maxSequential {
		return false
	}
	
	for i := 0; i <= len(password)-maxSequential; i++ {
		// 检查递增序列
		isIncreasing := true
		for j := 1; j < maxSequential; j++ {
			if password[i+j] != password[i+j-1]+1 {
				isIncreasing = false
				break
			}
		}
		
		// 检查递减序列
		isDecreasing := true
		for j := 1; j < maxSequential; j++ {
			if password[i+j] != password[i+j-1]-1 {
				isDecreasing = false
				break
			}
		}
		
		if isIncreasing || isDecreasing {
			return true
		}
	}
	
	return false
}

// GetPasswordStrength 获取密码强度评分 (0-100)
func (v *PasswordValidator) GetPasswordStrength(password string) int {
	score := 0
	
	// 长度评分 (最多30分)
	if len(password) >= 8 {
		score += 10
	}
	if len(password) >= 12 {
		score += 10
	}
	if len(password) >= 16 {
		score += 10
	}
	
	// 字符类型评分 (每种类型10分，最多40分)
	if regexp.MustCompile(`[a-z]`).MatchString(password) {
		score += 10
	}
	if regexp.MustCompile(`[A-Z]`).MatchString(password) {
		score += 10
	}
	if regexp.MustCompile(`\d`).MatchString(password) {
		score += 10
	}
	if regexp.MustCompile(`[@$!%*?&#+=-_()[\]{}|;:,.<>]`).MatchString(password) {
		score += 10
	}
	
	// 复杂度评分 (最多30分)
	if !v.hasRepeatingChars(password, 3) {
		score += 10
	}
	if !v.hasSequentialChars(password, 4) {
		score += 10
	}
	
	// 检查是否为弱密码
	lowerPassword := strings.ToLower(password)
	isWeak := false
	for _, weak := range v.weakPasswords {
		if lowerPassword == weak {
			isWeak = true
			break
		}
	}
	if !isWeak {
		score += 10
	}
	
	return score
}

// GetPasswordStrengthText 获取密码强度文本描述
func (v *PasswordValidator) GetPasswordStrengthText(score int) string {
	switch {
	case score >= 80:
		return "强"
	case score >= 60:
		return "中等"
	case score >= 40:
		return "弱"
	default:
		return "很弱"
	}
}
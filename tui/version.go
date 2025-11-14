package tui

import (
	"strings"
)

// Version information injected by main package
var (
	AppVersion string
)

// GetVersion returns the application version string
func GetVersion() string {
	if AppVersion == "" {
		return "dev"
	}
	// 简化版本格式：v1.0.0-1-g81049f3 -> v1.0.0-g81049f3
	// 移除中间的提交计数
	simplified := simplifyVersion(AppVersion)
	return simplified
}

// simplifyVersion 简化版本格式
// 输入: v1.0.0-1-g81049f3
// 输出: v1.0.0-g81049f3
func simplifyVersion(version string) string {
	// 如果包含标准的git describe格式，简化它
	parts := strings.Split(version, "-")
	if len(parts) >= 3 {
		// 格式: version-count-hash -> version-hash
		// 检查第三个部分是否以g开头（git提交哈希）
		if len(parts[2]) > 1 && parts[2][0] == 'g' {
			return parts[0] + "-" + parts[2]
		}
	}
	// 其他情况保持原样（如纯版本号v1.0.0）
	return version
}
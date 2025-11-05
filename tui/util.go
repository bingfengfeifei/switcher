package tui

import "strings"

func maskAPIKey(key string) string {
	if len(key) <= 8 {
		return strings.Repeat("*", len(key))
	}
	return key[:4] + strings.Repeat("*", len(key)-8) + key[len(key)-4:]
}

// sanitizeInput removes bracketed paste markers and stray brackets/newlines.
func sanitizeInput(s string) string {
	if s == "" {
		return s
	}
	// Remove ANSI bracketed paste start/end sequences if present.
	s = strings.ReplaceAll(s, "\x1b[200~", "")
	s = strings.ReplaceAll(s, "\x1b[201~", "")
	// Remove any literal square brackets introduced by paste (not by single keypress).
	if len([]rune(s)) > 1 || strings.Contains(s, "\x1b") {
		s = strings.ReplaceAll(s, "[", "")
		s = strings.ReplaceAll(s, "]", "")
	}
	// Strip CR/LF to avoid breaking fields when pasting.
	s = strings.ReplaceAll(s, "\r", "")
	s = strings.ReplaceAll(s, "\n", "")
	return s
}

func hostOf(u string) string {
	s := strings.TrimSpace(u)
	if s == "" {
		return ""
	}
	if idx := strings.Index(s, "://"); idx != -1 {
		s = s[idx+3:]
	}
	if idx := strings.Index(s, "/"); idx != -1 {
		s = s[:idx]
	}
	return s
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

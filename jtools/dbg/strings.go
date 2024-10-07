package dbg

import (
	"fmt"
	"strings"
)

func TrimToLower(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}
func TrimToLowerPtr(s *string) string {
	*s = TrimToLower(*s)
	return *s
}
func TrimToUpper(s string) string {
	return strings.ToUpper(strings.TrimSpace(s))
}
func TrimToUpperPtr(s *string) string {
	*s = TrimToUpper(*s)
	return *s
}

func IsTrue(p interface{}) bool {
	switch v := p.(type) {
	case bool:
		return v
	case []bool:
		if len(v) > 0 {
			return v[0]
		}
	case string:
		return TrimToLower(v) == "true"
	case []interface{}:
		if len(v) > 0 {
			return IsTrue(v[0])
		}
	}
	return false
}

func Void(a interface{}) string {
	return fmt.Sprintf("%v", a)
}
func D(a ...interface{}) {}

func Cat(a ...interface{}) string {
	sl := make([]string, len(a))
	for i, v := range a {
		sl[i] = fmt.Sprintf("%v", v)
	}
	return strings.Join(sl, "")
}

// Key : a,b,c -> a_b_c
func Key(a ...interface{}) string {
	sl := make([]string, len(a))
	for i, v := range a {
		sl[i] = fmt.Sprintf("%v", v)
	}
	return strings.Join(sl, "_")
}

func ShortAddress(v string, is_skips ...bool) string {
	if IsTrue(is_skips) {
		return v
	}

	if len(v) <= 10 {
		return v
	}
	s := v[:6] + "...."
	s += v[len(v)-4:]
	return s
}

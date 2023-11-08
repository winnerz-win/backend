package dbg

import "strings"

//EmailCheck :
func EmailCheck(email string) bool {
	if strings.Contains(email, " ") {
		Red(1)
		return false
	}
	ss := strings.Split(email, "@")
	if len(ss) != 2 {
		Red(2)
		return false
	}
	if len(ss[0]) == 0 || len(ss[1]) == 0 {
		Red(3)
		return false
	}
	filters := func(c byte) bool {
		if c <= 32 { // NULL ~ SP
			Red(string(c))
			return true
		}
		switch c {
		case 34, 39, 40, 41, 42, 44, 47, 58, 59, 60, 61, 62, 63, 91, 92, 93, 94:
			Red(string(c))
			return true
		}
		if c < 0 || c > 127 {
			Purple(string(c))
			return true
		}
		return false
	}

	buf := []byte(email)
	for _, c := range buf {
		if filters(c) {
			Red(4)
			return false
		}
	}

	//Yellow(ss, len(ss))
	return true
}

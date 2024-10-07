package ebcm

import "strings"

func Is429Error(err error) bool {
	return strings.Contains(err.Error(), "429 Too Many Requests")
}

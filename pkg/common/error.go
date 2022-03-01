package common

import (
	"fmt"
)

func ErrorJson(s string) string {
	return fmt.Sprintf(`{"error": "%s"}`, s)
}

func ErrorString(s string, isRaw bool) string {
	if isRaw {
		return s
	} else {
		return ErrorJson(s)
	}
}

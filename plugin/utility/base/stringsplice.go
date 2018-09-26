package base

import "strings"

func StringSplice(src []string) string {
	var s strings.Builder
	s.Grow(len(src))

	for _, v := range src {
		s.WriteString(v)
	}

	return s.String()
}

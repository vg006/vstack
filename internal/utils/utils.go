package utils

func EchoPath(path string) string {
	result := ""
	for _, ch := range path {
		if ch == '{' {
			result += ":"
		} else if ch != '}' {
			result += string(ch)
		}
	}
	return result
}

func FiberPath(path string) string {
	result := ""
	for _, ch := range path {
		if ch == '{' {
			result += ":"
		} else if ch != '}' {
			result += string(ch)
		}
	}
	return result
}

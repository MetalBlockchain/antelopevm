package utils

func Assert(expression bool, message string, args ...interface{}) {
	if !expression {
		panic(message)
	}
}

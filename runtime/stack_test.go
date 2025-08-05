package runtimeext

func nested(level int) Frame {
	return StackLevel(level)
}

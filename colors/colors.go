package colors

func Bold(str string) string {
	return "\x1b[1m" + str + "\x1b[0m"
}

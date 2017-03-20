package picon

type Ansi string

const (
	fgBlack   Ansi = "\033[0;30m"
	fgRed     Ansi = "\033[0;31m"
	fgGreen   Ansi = "\033[0;32m"
	fgYellow  Ansi = "\033[0;33m"
	fgBlue    Ansi = "\033[0;34m"
	fgMagenta Ansi = "\033[0;35m"
	fgCyan    Ansi = "\033[0;36m"
	fgWhite   Ansi = "\033[0;37m"
	attrReset Ansi = "\033[0m"
)

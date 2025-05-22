package color

type Color string

const (
	Reset   Color = "\033[0m"
	Red     Color = "\033[31m"
	Cyan    Color = "\033[36m"
	Blue    Color = "\033[34m"
	Green   Color = "\033[32m"
	Yellow  Color = "\033[33m"
	Magenta Color = "\033[35m"
)

func (c Color) String() string {
	return string(c)
}

func Colorize(c Color, s string) string {
	return c.String() + s + Reset.String()
}

package logfx

type Color string

const (
	ColorReset        Color = "\033[0m"
	ColorRed          Color = "\033[31m"
	ColorGreen        Color = "\033[32m"
	ColorYellow       Color = "\033[33m"
	ColorBlue         Color = "\033[34m"
	ColorMagenta      Color = "\033[35m"
	ColorCyan         Color = "\033[36m"
	ColorGray         Color = "\033[37m"
	ColorDimGray      Color = "\033[90m"
	ColorLightRed     Color = "\033[91m"
	ColorLightGreen   Color = "\033[92m"
	ColorLightYellow  Color = "\033[93m"
	ColorLightBlue    Color = "\033[94m"
	ColorLightMagenta Color = "\033[95m"
	ColorLightCyan    Color = "\033[96m"
	ColorLightGray    Color = "\033[97m"
)

func Colored(color Color, message string) string {
	return string(color) + message + string(ColorReset)
}

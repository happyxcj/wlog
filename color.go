package wlog

import "fmt"

// Foreground colors for log level prefix.
const (
	NoColor Color = 0
	Red     Color = 30 + iota
	Green
	Yellow
	Blue
)

// Color is the type defined for level prefix color.
type Color uint8

// With returns a string with the specified color.
func (c Color) With(str string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", uint8(c), str)
}

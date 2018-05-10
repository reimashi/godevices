package displays

import "github.com/reimashi/godevices"

type Display interface {
	godevices.Device
}

type TextDisplay interface {
	GetLineCount() int
	GetLineSize(line int) int
	Clear() error
	ClearLine(line int) error
	Write(str string, line int, offset int, clearLine bool) error
	Display
}

type PixelDisplay interface {
	GetWidth() uint64
	GetHeight() uint64
	Display
}
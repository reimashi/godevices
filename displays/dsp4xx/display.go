package dsp4xx

import (
	goserial "github.com/tarm/goserial"
	"io"
	"runtime"
	"sync"
	"time"
	"github.com/reimashi/godevices/displays"
	"github.com/reimashi/godevices/interfaces"
	"github.com/reimashi/godevices/interfaces/serial"
	"errors"
)

type displayCommand byte

const (
	COM_SET_BAUDRATE displayCommand = 0x42
	COM_CLEAR        displayCommand = 0x43
	COM_SHOW_DEMO    displayCommand = 0x44
	COM_SET_CURSOR   displayCommand = 0x50
	COM_SAVE_DEMO    displayCommand = 0x53
	COM_GET_SCREEN   displayCommand = 0x54
	COM_SET_MODE     displayCommand = 0x7e
)

type displayMode byte

const (
	MODE_DSP_T         displayMode = 0x00
	MODE_EPSON_ESC_POS displayMode = 0x01
	MODE_UTC_STANDARD  displayMode = 0x02
	MODE_UTC_ENHANCED  displayMode = 0x03
	MODE_AEDEX         displayMode = 0x04
	MODE_ICD_2002      displayMode = 0x05
	MODE_CD_5220       displayMode = 0x06
	MODE_DSP_800       displayMode = 0x07
)

const (
	model = "DSP-4xx"
	vendor = "unknown"
	lines = 2
	refresh = 80
)

var (
	lineLength = []int{20, 20}
)

type Dsp4xx struct {
	serial.Device
	displays.TextDisplay

	config *serial.Config
	serialPort io.ReadWriteCloser
	mutexWr *sync.Mutex

	RefreshTime int
}

func Open(config *serial.Config) (*Dsp4xx, error) {
	this := &Dsp4xx{ config: config }

	c := &goserial.Config{Name: this.config.GetPortAddress(), Baud: this.config.GetBaud()}

	s, err := goserial.OpenPort(c)

	if err != nil {
		return nil, err
	}

	this.serialPort = s
	this.mutexWr = &sync.Mutex{}
	this.deviceMode(MODE_DSP_T)

	return this, nil
}

// TextDisplay interface

func (t *Dsp4xx) Clear() error {
	return t.deviceClear(1, 40)
}

func (t *Dsp4xx) ClearLine(line int) error {
	if line >= 0 && line < lines {
		startpos := (line * lineLength[line]) + 1
		return t.deviceClear(startpos, startpos + lineLength[line]-1)
	}

	return errors.New("invalid line")
}

func (t *Dsp4xx) GetLineCount() int {
	return lines
}

func (t *Dsp4xx) GetLineSize(line int) int {
	if line >= 0 && line < lines {
		return lineLength[line]
	}

	return 0
}

func (t *Dsp4xx) Write(str string, line int, offset int, clearLine bool) error {
	if line >= 0 && line < lines {
		if clearLine {
			err := t.ClearLine(line)
			if err != nil { return err }
		}

		if offset >= 0 && offset < lineLength[line] {
			t.deviceSetCursor(line, offset)

			countw := len(str)
			if offset+countw > lineLength[line] {
				countw = lineLength[line] - offset
			}

			towrite := str
			if len(str) > countw {
				towrite = str[:countw]
			}

			return t.deviceWrite([]byte(towrite), refresh)
		} else {
			return errors.New("invalid offset")
		}
	} else {
		return errors.New("invalid line")
	}
}

func (t *Dsp4xx) dspCodificationNumber(i int) byte {
	if i >= 0 || i <= lineLength[1] * lines {
		return byte(i + 48)
	} else {
		return 0
	}
}

func (t *Dsp4xx) deviceMode(mode displayMode) {
	code := []byte{byte(COM_SET_MODE), byte(mode), byte(COM_SET_MODE)}
	t.deviceWrite(code, 1200)
}

func (t *Dsp4xx) deviceClear(start int, end int) error {
	startb := t.dspCodificationNumber(start)
	endb := t.dspCodificationNumber(end)
	code := []byte{0x04, 0x01, byte(COM_CLEAR), startb, endb, 0x17}
	return t.deviceWrite(code, refresh)
}

func (t *Dsp4xx) deviceSetCursor(line int, offset int) {
	cpos := t.dspCodificationNumber((line * lineLength[line]) + offset + 1)
	code := []byte{0x04, 0x01, byte(COM_SET_CURSOR), cpos, 0x17}
	t.deviceWrite(code, refresh)
}

func (t *Dsp4xx) deviceWrite(data []byte, refresh int) error {
	t.mutexWr.Lock()

	_, err := t.serialPort.Write(data)

	if err == nil {
		time.Sleep(time.Duration(refresh) * time.Millisecond)
	}

	t.mutexWr.Unlock()
	runtime.Gosched()

	return err
}

// Device interface

func (t *Dsp4xx) GetModel() string {
	return model
}

func (t *Dsp4xx) GetVendor() string {
	return vendor
}

func (t *Dsp4xx) GetInterfaceType() interfaces.Type {
	return interfaces.Serial
}
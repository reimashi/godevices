package displays

import (
	"github.com/tarm/goserial"
	"io"
	"runtime"
	"sync"
	"time"
)

type IDisplay interface {
	Clear()
	ClearLine(line int)
	GetLineLength(line int) int
	GetName() string
	GetNumLines() int
	Init() error
	Write(str string, line int, offset int, clear bool)
}

type DisplayInfo struct {
	Name        string
	Lines       int
	LineLength  []int
	RefreshTime int
}

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

type DisplayDsp4xx struct {
	Port       string
	Baud       int
	serialPort io.ReadWriteCloser
	mutexWr    *sync.Mutex
	info       DisplayInfo
}

func Open(port string, baud int) (*DisplayDsp4xx, error) {
	this := &DisplayDsp4xx{
		Port: port,
		Baud: baud,
	}

	c := &serial.Config{Name: this.Port, Baud: this.Baud}

	s, err := serial.OpenPort(c)

	if err != nil {
		return nil, err
	}

	this.serialPort = s
	this.mutexWr = &sync.Mutex{}
	this.info = DisplayInfo{
		Name:        "DSP-4XX",
		Lines:       2,
		LineLength:  []int{20, 20},
		RefreshTime: 80}

	this.deviceMode(MODE_DSP_T)

	return this, nil
}

func (t DisplayDsp4xx) Clear() {
	t.deviceClear(1, 40)
}

func (t DisplayDsp4xx) ClearLine(line int) {
	if line >= 0 && line < t.info.Lines {
		startpos := (line * t.info.LineLength[line]) + 1
		t.deviceClear(startpos, startpos+t.info.LineLength[line]-1)
	}
}

func (t DisplayDsp4xx) GetLineLength(line int) int {
	if line >= 0 && line < t.info.Lines {
		return t.info.LineLength[line]
	}

	return 0
}

func (t DisplayDsp4xx) GetName() string {
	return t.info.Name
}

func (t DisplayDsp4xx) GetNumLines() int {
	return t.info.Lines
}

func (t DisplayDsp4xx) Write(str string, line int, offset int, clear bool) {
	if line >= 0 && line < t.info.Lines {
		if clear {
			t.ClearLine(line)
		}

		if offset >= 0 && offset < t.info.LineLength[line] {
			t.deviceSetCursor(line, offset)

			countw := len(str)
			if offset+countw > t.info.LineLength[line] {
				countw = t.info.LineLength[line] - offset
			}

			towrite := str
			if len(str) > countw {
				towrite = str[:countw]
			}

			t.deviceWrite([]byte(towrite), t.info.RefreshTime)
		}
	}
}

func (t DisplayDsp4xx) dspCodificationNumber(i int) byte {
	if i >= 0 || i <= t.info.LineLength[1]*t.info.Lines {
		return byte(i + 48)
	} else {
		return 0
	}
}

func (t DisplayDsp4xx) deviceMode(mode displayMode) {
	code := []byte{byte(COM_SET_MODE), byte(mode), byte(COM_SET_MODE)}
	t.deviceWrite(code, 1200)
}

func (t DisplayDsp4xx) deviceClear(start int, end int) {
	startb := t.dspCodificationNumber(start)
	endb := t.dspCodificationNumber(end)
	code := []byte{0x04, 0x01, byte(COM_CLEAR), startb, endb, 0x17}
	t.deviceWrite(code, t.info.RefreshTime)
}

func (t DisplayDsp4xx) deviceSetCursor(line int, offset int) {
	cpos := t.dspCodificationNumber((line * t.info.LineLength[line]) + offset + 1)
	code := []byte{0x04, 0x01, byte(COM_SET_CURSOR), cpos, 0x17}
	t.deviceWrite(code, t.info.RefreshTime)
}

func (t DisplayDsp4xx) deviceWrite(data []byte, refresh int) error {
	t.mutexWr.Lock()

	_, err := t.serialPort.Write(data)

	if err == nil {
		time.Sleep(time.Duration(refresh) * time.Millisecond)
	}

	t.mutexWr.Unlock()
	runtime.Gosched()

	return err
}

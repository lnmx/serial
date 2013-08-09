package serial

import (
	"fmt"
	"strings"
	"syscall"
	"unsafe"
)

type Port struct {
	config Config
	handle syscall.Handle
}

func (p *Port) configure(cfg Config) (err error) {
	p.config = cfg

	return
}

func (p *Port) open() (err error) {

	// open serial port handle
	//
	err = p.openHandle()

	if err != nil {
		return fmt.Errorf("error opening serial device %s: %s", p.config.Device, err)
	}

	// configure baud, data, parity
	//
	err = p.setCommState()

	if err != nil {
		return fmt.Errorf("error applying serial settings: %s", err)
	}

	// set read/write timeouts
	//
	err = p.setTimeouts()

	if err != nil {
		return fmt.Errorf("error setting timeouts: %s", err)
	}

	return
}

func (p *Port) openHandle() (err error) {
	device := p.config.Device

	// add a "\\.\" prefix, e.g. "\\.\COM2"
	// optional for COM1-9, required for COM10 and up
	//
	if !strings.HasPrefix(device, `\\.\`) {
		device = `\\.\` + device
	}

	device_utf16, err := syscall.UTF16PtrFromString(device)

	if err != nil {
		return
	}

	p.handle, err = syscall.CreateFile(device_utf16,
		syscall.GENERIC_READ|syscall.GENERIC_WRITE,
		0,
		nil,
		syscall.OPEN_EXISTING,
		syscall.FILE_ATTRIBUTE_NORMAL,
		0)

	return
}

func (p *Port) setCommState() (err error) {

	dcb := &_DCB{
		DCBlength: 28,
	}

	err = _GetCommState(p.handle, dcb)

	if err != nil {
		return
	}

	if v, ok := settings[p.config.BaudRate]; ok {
		dcb.BaudRate = uint32(v)
	} else {
		return fmt.Errorf("unsupported baud rate: %d", p.config.BaudRate)
	}

	if v, ok := settings[p.config.DataBits]; ok {
		dcb.ByteSize = byte(v)
	} else {
		return fmt.Errorf("unsupported data bits: %d", p.config.DataBits)
	}

	if v, ok := settings[p.config.Parity]; ok {
		dcb.Parity = byte(v)
	} else {
		return fmt.Errorf("unsupported parity: %d", p.config.Parity)
	}

	if v, ok := settings[p.config.StopBits]; ok {
		dcb.StopBits = byte(v)
	} else {
		return fmt.Errorf("unsupported stop bits: %d", p.config.StopBits)
	}

	err = _SetCommState(p.handle, dcb)

	return
}

func (p *Port) setTimeouts() (err error) {
	timeouts := &_COMMTIMEOUTS{
		ReadIntervalTimeout:         1,
		ReadTotalTimeoutConstant:    1,
		ReadTotalTimeoutMultiplier:  1,
		WriteTotalTimeoutConstant:   1,
		WriteTotalTimeoutMultiplier: 1,
	}

	err = _SetCommTimeouts(p.handle, timeouts)

	return
}

func (p *Port) close() (err error) {
	return syscall.CloseHandle(p.handle)
}

func (p *Port) read(b []byte) (n int, err error) {
	var done uint32
	err = syscall.ReadFile(p.handle, b, &done, nil)
	n = int(done)

	return
}

func (p *Port) write(b []byte) (n int, err error) {
	var done uint32
	err = syscall.WriteFile(p.handle, b, &done, nil)
	n = int(done)

	return
}

func (p *Port) flush() (err error) {
	return syscall.FlushFileBuffers(p.handle)
}

var (
	// handles to some serial-related functions not supported by syscall
	//
	kernel32, _           = syscall.LoadLibrary("kernel32.dll")
	getCommState, _       = syscall.GetProcAddress(kernel32, "GetCommState")
	setCommState, _       = syscall.GetProcAddress(kernel32, "SetCommState")
	setCommTimeouts, _    = syscall.GetProcAddress(kernel32, "SetCommTimeouts")
	getCommTimeouts, _    = syscall.GetProcAddress(kernel32, "GetCommTimeouts")
	escapeCommFunction, _ = syscall.GetProcAddress(kernel32, "EscapeCommFunction")

	// map the generic constants in serial.go to values for DCB
	// (ref. WinBase.h)
	//
	settings = map[int]int{
		DataBits_5:      5,
		DataBits_6:      6,
		DataBits_7:      7,
		DataBits_8:      8,
		StopBits_1:      0,
		StopBits_1_5:    1,
		StopBits_2:      2,
		Parity_None:     0,
		Parity_Odd:      1,
		Parity_Even:     2,
		Parity_Mark:     3,
		Parity_Space:    4,
		BaudRate_9600:   9600,
		BaudRate_14400:  14400,
		BaudRate_19200:  19200,
		BaudRate_38400:  38400,
		BaudRate_57600:  57600,
		BaudRate_115200: 115200,
	}
)

// sizeof(DCB) = 28 bytes

type _DCB struct {
	DCBlength uint32
	BaudRate  uint32

	// Flags:
	//   1 binary
	//   1 enable parity
	//   1 out cts flow control
	//   1 out dtr flow control
	//   2 dtr flow control
	//   1 dsr sensitivity
	//   1 continue tx on xoff
	//   1 enable output x-on/x-off
	//   1 enable input x-on/x-off
	//   1 enable err replacement
	//   1 enable null stripping
	//   2 rts flow control
	//   1 abort reads/writes on error
	//  17 reserved
	Flags     uint32
	Reserved1 uint16
	XonLim    uint16
	XoffLim   uint16
	ByteSize  byte
	Parity    byte
	StopBits  byte
	XonChar   byte
	XoffChar  byte
	ErrorChar byte
	EofChar   byte
	EvtChar   byte
	Reserved2 uint16
}

// COMMTIMEOUTS
//
type _COMMTIMEOUTS struct {
	ReadIntervalTimeout         uint32
	ReadTotalTimeoutMultiplier  uint32
	ReadTotalTimeoutConstant    uint32
	WriteTotalTimeoutMultiplier uint32
	WriteTotalTimeoutConstant   uint32
}

func _GetCommTimeouts(handle syscall.Handle, timeouts *_COMMTIMEOUTS) (err error) {
	// BOOL GetCommTimeouts( HANDLE hFile, LPCOMMTIMEOUTS lpCommTimeouts);

	r0, _, e1 := syscall.Syscall(getCommTimeouts, 2, uintptr(handle), uintptr(unsafe.Pointer(timeouts)), 0)

	if r0 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}

	return
}

func _SetCommTimeouts(handle syscall.Handle, timeouts *_COMMTIMEOUTS) (err error) {
	// BOOL SetCommTimeouts( HANDLE hFile, LPCOMMTIMEOUTS lpCommTimeouts);

	r0, _, e1 := syscall.Syscall(setCommTimeouts, 2, uintptr(handle), uintptr(unsafe.Pointer(timeouts)), 0)

	if r0 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}

	return
}

func _GetCommState(handle syscall.Handle, dcb *_DCB) (err error) {
	// BOOL GetCommState( HANDLE hFile, LPDCB lpDCB );

	r0, _, e1 := syscall.Syscall(getCommState, 2, uintptr(handle), uintptr(unsafe.Pointer(dcb)), 0)

	if r0 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}

	return
}

func _SetCommState(handle syscall.Handle, dcb *_DCB) (err error) {
	// BOOL SetCommState( HANDLE hFile, LPDCB lpDCB );

	r0, _, e1 := syscall.Syscall(setCommState, 2, uintptr(handle), uintptr(unsafe.Pointer(dcb)), 0)

	if r0 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}

	return
}

// operations for EscapeCommFunction
// (not exposed in the API yet)
//
type escapeFn int

const (
	_SETXOFF  escapeFn = 1
	_SETXON            = 2
	_SETRTS            = 3
	_CLRRTS            = 4
	_SETDTR            = 5
	_CLRDTR            = 6
	_RESETDEV          = 7
	_SETBREAK          = 8
	_CLRBREAK          = 9
)

func _EscapeCommFunction(handle syscall.Handle, escape escapeFn) (err error) {
	// BOOL EscapeCommFunction( HANDLE hFile, DWORD dwFunc );

	r0, _, e1 := syscall.Syscall(escapeCommFunction, 2, uintptr(handle), uintptr(escape), 0)

	if r0 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}

	return
}

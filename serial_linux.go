package serial

/*

Serial syscalls adapted from github.com/dustin/go-rs232 via github.com/schleibinger/sio

1) Original: Copyright (c) 2005-2008 Dustin Sallings <dustin@spy.net>.

2) Mods: Copyright (c) 2012 Schleibinger Gerä Teubert u. Greim GmbH
<info@schleibinger.com>. Blame: Jan Mercl

All rights reserved.  Use of this source code is governed by a MIT-style
license that can be found in the LICENSE file.

*/

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

type Port struct {
	config Config
	handle *os.File
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
		p.Close()

		return fmt.Errorf("error opening serial device %s: %s", p.config.Device, err)
	}

	return
}

func (p *Port) openHandle() (err error) {

	rate := uint32(settings[p.config.BaudRate])

	p.handle, err = os.OpenFile(p.config.Device, syscall.O_RDWR|syscall.O_NOCTTY|syscall.O_NDELAY, 0666)

	if err != nil {
		return
	}

	fd := p.handle.Fd()

	term := syscall.Termios{
		Iflag:  syscall.IGNPAR,
		Cflag:  syscall.CS8 | syscall.CREAD | syscall.CLOCAL | rate,
		Cc:     [32]uint8{syscall.VMIN: p.config.VMIN, syscall.VTIME: p.config.VTIME},
		Ispeed: rate,
		Ospeed: rate,
	}

	_, _, errno := syscall.Syscall6(
		syscall.SYS_IOCTL,
		uintptr(fd),
		uintptr(syscall.TCSETS),
		uintptr(unsafe.Pointer(&term)),
		0,
		0,
		0,
	)

	if errno != 0 {
		err = errno
	}

	if err != nil {
		return
	}

	err = syscall.SetNonblock(int(fd), false)

	if err != nil {
		return
	}

	return
}

func (p *Port) close() (err error) {
	if p.handle != nil {
		err = p.handle.Close()
	}

	return
}

func (p *Port) read(b []byte) (n int, err error) {
	return p.handle.Read(b)
}

func (p *Port) write(b []byte) (n int, err error) {
	return p.handle.Write(b)
}

var (

	// map the generic constants in serial.go to values for termios
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
		BaudRate_9600:   syscall.B9600,
		BaudRate_19200:  syscall.B19200,
		BaudRate_38400:  syscall.B38400,
		BaudRate_57600:  syscall.B57600,
		BaudRate_115200: syscall.B115200,
	}
)

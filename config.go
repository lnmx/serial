package serial

type Config struct {

	// Device Name (ex. "COM1" on Windows, "/dev/ttyUSB0" on Linux)
	//
	Device string

	// Baud Rate, use BaudRate_*, default 9600
	//
	BaudRate int

	// Data Bits per byte, use DataBits_*, default 8
	//
	DataBits int

	// Stop Bits per byte, use StopBits_*, default 1
	//
	StopBits int

	// Parity, use Parity_*, default None
	//
	Parity int

	// termios VMIN, default 1
	//
	// see 'man termios(3)'
	//
	VMIN uint8

	// termios VTIME, default 0
	//
	VTIME uint8

	// win32 COMMTIMEOUTS
	//
	// see http://msdn.microsoft.com/en-us/library/windows/desktop/aa363190(v=vs.85).aspx/html
	//
	ReadIntervalTimeout         uint32
	ReadTotalTimeoutMultiplier  uint32
	ReadTotalTimeoutConstant    uint32
	WriteTotalTimeoutMultiplier uint32
	WriteTotalTimeoutConstant   uint32
}

const (
	BaudRate_9600   = 9600
	BaudRate_14400  = 14400
	BaudRate_19200  = 19200
	BaudRate_38400  = 38400
	BaudRate_57600  = 57600
	BaudRate_115200 = 115200
)

const (
	DataBits_5 int = 2000000
	DataBits_6     = iota
	DataBits_7
	DataBits_8

	StopBits_1
	StopBits_1_5
	StopBits_2

	Parity_None
	Parity_Odd
	Parity_Even
	Parity_Mark
	Parity_Space
)

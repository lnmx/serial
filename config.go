package serial

type Config struct {
	Device   string
	BaudRate int
	DataBits int
	StopBits int
	Parity   int
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

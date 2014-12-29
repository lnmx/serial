package main

import (
	"fmt"
	"github.com/lnmx/serial"
	"time"
)

func main() {
	err := run()

	if err != nil {
		fmt.Println(err)
	}
}

func run() (err error) {

	var device string = "COM3"
	var baud int = 9600

	var read_constant uint32 = 500
	var read_multiplier uint32 = 100
	var read_length int = 10

	port := serial.NewPort()

	config := port.Config()
	config.Device = device
	config.BaudRate = baud
	config.ReadTotalTimeoutConstant = read_constant
	config.ReadTotalTimeoutMultiplier = read_multiplier

	err = port.Configure(config)

	if err != nil {
		return err
	}

	err = port.Open()

	if err != nil {
		return err
	}

	defer port.Close()

	buf := make([]byte, read_length)

	begin := time.Now()
	n, err := port.Read(buf)
	elapsed := int64(time.Since(begin) / time.Millisecond)

	fmt.Printf("read = %v, time = %v ms, err = %v", n, elapsed, err)

	return nil
}

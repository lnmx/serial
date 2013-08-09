package main

import (
	"bufio"
	"fmt"
	"github.com/lnmx/serial"
	"os"
)

func main() {
	device := "COM9"
	baud := 9600

	fmt.Println("open", device, "at", baud)

	port, err := serial.Open(device, baud)

	if err != nil {
		fmt.Println("open failed:", err)
		return
	}

	defer port.Close()

	fmt.Println("ready")

	// display data from serial
	//
	go func() {
		buf := make([]byte, 32)

		for {
			n, err := port.Read(buf)

			if err != nil {
				fmt.Println("serial read error:", err)
				return
			}

			if n > 0 {
				fmt.Println(n, ">", string(buf[:n]))
			}
		}
	}()

	// send user input (by line) to serial
	//
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		_, err := port.Write([]byte(scanner.Text() + "\n"))

		if err != nil {
			fmt.Println("serial write error:", err)
		}
	}

}

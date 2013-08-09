package main

import (
	"bufio"
	"fmt"
	"github.com/lnmx/serial"
	"os"
)

func main() {
	device := "COM4"
	baud := 115200

	fmt.Println("open", device, "at", baud)

	port, err := serial.Open(device, baud)

	if err != nil {
		fmt.Println("open failed:", err)
		return
	}

	defer port.Close()

	// display data from serial
	//
	go func() {
		scanner := bufio.NewScanner(port)

		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			fmt.Println("serial read error:", err)
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

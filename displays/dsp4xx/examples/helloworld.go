package main

import (
	"github.com/reimashi/godevices/displays/dsp4xx"
	"github.com/reimashi/godevices/interfaces/serial"
	"fmt"
)

func main() {
	devConn, err := dsp4xx.Open(serial.NewConfig("1", serial.Baud19200))

	fmt.Println("Device name: " + devConn.GetModel())

	if err != nil {
		fmt.Println("Error while connecting to device: " + err.Error())
		return
	}

	devConn.Write("Hello", 0, 0, true)
	devConn.Write("World", 1, 6, true)
}

package main

import (
	"image"
	"os"

	_ "image/png"

	"golang.org/x/exp/io/spi"
	"github.com/reimashi/godevices/displays/ssd1306"
)

func main() {
	rc, err := os.Open("./test.png")
	if err != nil {
		panic(err)
	}
	defer rc.Close()

	m, _, err := image.Decode(rc)
	if err != nil {
		panic(err)
	}

	d, err := ssd1306.Open(&spi.Devfs{
		Dev:      "/dev/spidev0.0",
		MaxSpeed: int64(8000000),
	},128, 64)

	if err != nil {
		panic(err)
	}
	defer d.Close()

	// clear the display before putting on anything
	if err := d.Clear(); err != nil {
		panic(err)
	}
	if err := d.SetImage(0, 0, m); err != nil {
		panic(err)
	}
	if err := d.Draw(); err != nil {
		panic(err)
	}
}

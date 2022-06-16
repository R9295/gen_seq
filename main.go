package main

import (
	"fmt"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers/rtmididrv"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv" // autoregisters driver
)

func main() {
	driver, err := rtmididrv.New()
	if err != nil {
		panic(err)
	}
	inPort, err := driver.OpenVirtualIn("sequencer")
	if err != nil {
		panic(err)
	}
	_, err = midi.ListenTo(inPort, func(msg midi.Message, timestamps int32) {
		fmt.Println(msg)
	}, midi.UseTimeCode())
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}
	for true {

	}
}

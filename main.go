package main

import (
	"fmt"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers/rtmididrv"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv" // autoregisters driver
)

func pulseIsQuarter(pulse int) bool {
	switch pulse {
	case 0,
		24,
		48,
		72,
		96:
		return true
	}
	return false
}

func main() {
	driver, err := rtmididrv.New()
	if err != nil {
		panic(err)
	}
	inPort, err := driver.OpenVirtualIn("sequencer_in")
	outPort, err := driver.OpenVirtualOut("sequencer_out")
	if err != nil {
		panic(err)
	}
	pulses := 0
	pulsesPerQuarter := 24
	pulsesPerBar := pulsesPerQuarter * 4
	bars := 1
	_, err = midi.ListenTo(inPort, func(msg midi.Message, timestamps int32) {
		switch msg.Type() {
		case midi.TimingClockMsg:
			pulses++
			if pulses == pulsesPerBar {
				bars++
				pulses = 0
				fmt.Printf("Bar %d\n", bars)
			}
			if pulseIsQuarter(pulses) {
				fmt.Println("Quarter")
				outPort.Send(midi.NoteOn(1, midi.A(2), 127))
			}
		case midi.StopMsg:
			pulses = 0
			bars = 1
		}

	}, midi.UseTimeCode())
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}
	for true {
	}
}

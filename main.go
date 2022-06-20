package main

import (
	"fmt"
	"math/rand"
	"time"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers/rtmididrv"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv" // autoregisters driver
)

func pulseIsQuarter(pulse int) bool {
	if pulse%6 == 0 {
		return true
	}
	return false
}

func getRandomInt(a int, b int) uint8 {
	rand.Seed(time.Now().UnixNano())
	return uint8(a + rand.Intn(b-a+1))
}

func getRandomNote(octave uint8) uint8 {
	notes := [12]uint8{
		midi.C(octave),
		midi.Db(octave),
		midi.D(octave),
		midi.Eb(octave),
		midi.E(octave),
		midi.F(octave),
		midi.Gb(octave),
		midi.G(octave),
		midi.Ab(octave),
		midi.A(octave),
		midi.Bb(octave),
		midi.B(octave),
	}
	return notes[getRandomInt(0, 11)]
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
	var prev_note uint8 = getRandomNote(3)
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
				note := getRandomNote(3)
				fmt.Printf("playing note %d\n", note)
				outPort.Send(midi.NoteOff(1, prev_note))
				outPort.Send(midi.NoteOn(1, note, 127))
				prev_note = note
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

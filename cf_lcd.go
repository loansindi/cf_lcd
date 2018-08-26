// Package cf_lcd provides functionality for CrystalFontz CFA533 series displays
package cf_lcd

import (
	"encoding/binary"
	"errors"
	"github.com/tarm/serial"
	"log"
)

func makecrc(ptr []byte) uint16 {
	var crc uint16
	crc = 0xFFFF

	for j := 0; j < len(ptr); j++ {
		data := ptr[j]
		for i := 8; i != 0; i-- {
			if (crc^uint16(data))&0x01 != 0 {
				crc >>= 1
				crc ^= 0x8408
			} else {
				crc >>= 1
			}
			data >>= 1
		}
	}
	return ^crc
}

func emptyBuffer(p *serial.Port) {

	buf := make([]byte, 4)
	_, err := p.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
}

func Backlight(p *serial.Port, b int) {
	bright := make([]byte, 5)
	bright[0] = 0x0E
	bright[1] = 0x01
	bright[2] = byte(b)
	crc := makecrc(bright[:3])
	binary.LittleEndian.PutUint16(bright[3:], crc)
	send := bright[0:5]
	p.Write([]byte(send[0:]))
}

func Clear(p *serial.Port) {
	clear := make([]byte, 4)
	clear[0] = 0x06
	clear[1] = 0x00
	b := makecrc(clear[:2])
	binary.LittleEndian.PutUint16(clear[2:], b)
	send := clear[0:4]
	p.Write([]byte(send[0:]))
}

func Write(p *serial.Port, row int, col int, message string) (err error) {
	msg := make([]byte, len(message)+6)
	if len(message) > 16 {
		return errors.New("Message too long!")
	}
	if row != 0 && row != 1 {
		return errors.New("Invalid row selection!")
	}
	if col > 16 {
		return errors.New("Invalid column selection!")
	}
	msg[0] = 0x1F
	msg[1] = byte(len(message) + 2)
	msg[2] = byte(col)
	msg[3] = byte(row)
	for i, character := range message {
		msg[i+4] = byte(character)
	}
	c := makecrc(msg[0 : len(message)+4])
	binary.LittleEndian.PutUint16(msg[len(message)+4:], c)
	send := msg[0 : len(message)+6]
	p.Write([]byte(send[0:]))
	return nil
}

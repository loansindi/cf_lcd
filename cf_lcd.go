// Package cf_lcd provides functionality for CrystalFontz CFA533 series displays
package cf_lcd

import (
	"encoding/binary"
	"fmt"
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

func Clear(p *serial.Port) {
	clear := make([]byte, 4)
	clear[0] = 0x06
	clear[1] = 0x00
	b := makecrc(clear[:2])
	binary.LittleEndian.PutUint16(clear[2:], b)
	send := clear[0:4]
	p.Write([]byte(send))
}

func WriteLine1(p *serial.Port, message string) {
	msg := make([]byte, 22)
	var line2 string
	if len(message) > 16 {
		line2 = message[16:]
		WriteLine2(p, line2)
		message = message[:16]

	}
	msg[0] = 0x1F
	msg[1] = byte(len(message) + 2)
	msg[2] = 0x00
	msg[3] = 0x00
	for i, character := range message {
		msg[i+4] = byte(character)
	}
	c := makecrc(msg[0 : len(message)+4])
	binary.LittleEndian.PutUint16(msg[len(message)+4:], c)
	send := msg[0 : len(message)+6]
	p.Write([]byte(send[0:]))
	fmt.Println([]byte(send))

}

func WriteLine2(p *serial.Port, message string) {
	if len(message) > 16 {
		message = message[:16]
	}
	msg := make([]byte, len(message)+6)
	msg[0] = 0x1F
	msg[1] = byte(len(message) + 2)
	msg[2] = 0x00
	msg[3] = 0x01
	for i, character := range message {
		msg[i+4] = byte(character)
	}
	c := makecrc(msg[0 : len(message)+4])
	binary.LittleEndian.PutUint16(msg[len(message)+4:], c)
	send := msg[0 : len(message)+6]
	p.Write([]byte(send[0:]))
	fmt.Println([]byte(send))

}

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

func handleBuffer(p *serial.Port) ([]byte, error) {
	buf := make([]byte, 16)
	_, err := p.Read(buf)
	if err != nil {

		log.Fatal(err)
	}
	for index := range buf {
		if buf[index] == 87 || buf[index] == 88 {
			length := buf[1]
			packet := make([]byte, length+2)
			packet = buf[2 : len(buf)-2]
			if err != nil {
				log.Fatal("Error reading packet")
			}
			return packet, err

		} else {
			return nil, err
		}
	}
	return nil, err
}

func Backlight(p *serial.Port, b int) ([]byte, error) {
	bright := make([]byte, 5)
	bright[0] = 0x0E
	bright[1] = 0x01
	bright[2] = byte(b)
	crc := makecrc(bright[:3])
	binary.LittleEndian.PutUint16(bright[3:], crc)
	send := bright[0:5]
	p.Write([]byte(send[0:]))

	packet, err := handleBuffer(p)
	return packet, err
}

func Clear(p *serial.Port) ([]byte, error) {
	clear := make([]byte, 4)
	clear[0] = 0x06
	clear[1] = 0x00
	b := makecrc(clear[:2])
	binary.LittleEndian.PutUint16(clear[2:], b)
	send := clear[0:4]
	p.Write([]byte(send[0:]))
	packet, err := handleBuffer(p)
	return packet, err

}

func Write(p *serial.Port, row int, col int, message string) (packet []byte, err error) {
	msg := make([]byte, len(message)+6)
	if len(message) > 16 {
		return nil, errors.New("Message too long!")
	}
	if row != 0 && row != 1 {
		return nil, errors.New("Invalid row selection!")
	}
	if col > 16 {
		return nil, errors.New("Invalid column selection!")
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
	packet, err = handleBuffer(p)
	return packet, err

}

func KeyReporting(p *serial.Port, mask []byte) (packet []byte, err error) {
	msg := make([]byte, 6)
	msg[0] = 0x17
	msg[1] = 0x02
	msg[2] = mask[0]
	msg[3] = mask[1]
	c := makecrc(msg[0:4])
	binary.LittleEndian.PutUint16(msg[4:], c)
	p.Write([]byte(msg[0:]))
	packet, err = handleBuffer(p)
	return packet, err
}

func GetKeys(p *serial.Port) ([]byte, error) {
	msg := make([]byte, 4)
	msg[0] = 0x18
	msg[1] = 0x00
	c := makecrc(msg[:2])
	binary.LittleEndian.PutUint16(msg[2:], c)
	p.Write([]byte(msg[0:]))
	packet, err := handleBuffer(p)
	return packet, err
}

func Flush(p *serial.Port) {
	p.Flush()
}

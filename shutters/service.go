package shutters

import (
	"github.com/damek86/go-somfy-rts"
)

type Service interface {
	MoveDown(addr uint32, rollingCode uint16)
	MoveUp(addr uint32, rollingCode uint16)
	MoveMy(addr uint32, rollingCode uint16)
	Program(addr uint32, rollingCode uint16)
}

type service struct {
	encryptionKey byte
	writeCommand  func(pulses []somfy.Pulse)
}

func NewService(encryptionKey byte, writeCommand func(pulses []somfy.Pulse)) Service {
	return &service{
		encryptionKey: encryptionKey,
		writeCommand:  writeCommand,
	}
}

func (s *service) MoveDown(addr uint32, rollingCode uint16) {
	s.move(addr, rollingCode, somfy.DownValue)
}

func (s *service) MoveUp(addr uint32, rollingCode uint16) {
	s.move(addr, rollingCode, somfy.UpValue)
}

func (s *service) MoveMy(addr uint32, rollingCode uint16) {
	s.move(addr, rollingCode, somfy.MyValue)
}

func (s *service) Program(addr uint32, rollingCode uint16) {
	data := somfy.NewPayload(s.encryptionKey)
	data.Address = addr
	data.Control = somfy.ProgramValue
	data.RollingCode = rollingCode

	s.sendCommand(data, 1)
}

func (s *service) move(addr uint32, rollingCode uint16, control somfy.Control) {
	data := somfy.NewPayload(s.encryptionKey)
	data.Address = addr
	data.Control = control
	data.RollingCode = rollingCode

	s.sendCommand(data, 2)
}

func (s *service) sendCommand(data somfy.PayloadData, frameRepetition int) {
	pulseWave := somfy.NewEncoder().Marshal(data, frameRepetition)
	s.writeCommand(pulseWave)
}

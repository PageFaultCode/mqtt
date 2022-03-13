// Package mqtt is the main encompassing mqtt package
package mqtt

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

const (
	singleByteLength           = 127
	singleByteCount            = 1
	twoByteLength              = 128
	twoByteCount               = 2
	twoByteUpperValue   uint8  = 1
	threeByteLength            = 16384
	threeByteCount             = 3
	threeByteUpperValue uint8  = 1
	fourByteLength             = 2097152
	fourByteCount              = 4
	fourByteUpperValue  uint8  = 1
	index0                     = 0
	index1                     = 1
	index2                     = 2
	index3                     = 3
	carryByte           uint8  = 0x80
	maxByte             uint8  = 0xff
	twoByteValue        uint32 = 0x3fff
	threeByteValue      uint32 = 2097151
	fourByteValue       uint32 = 268435455
	errorValue          uint32 = 0
)

type PacketTest struct {
	suite.Suite
}

func (m *PacketTest) TestLength() {
	message1 := ControlPacket{length: singleByteLength}

	encodedLength := message1.Length()

	m.Assert().Equal(singleByteCount, len(encodedLength))

	message2 := ControlPacket{length: twoByteLength}
	encodedLength = message2.Length()

	m.Assert().Equal(twoByteCount, len(encodedLength))
	m.Assert().Equal(carryByte, encodedLength[index0])
	m.Assert().Equal(twoByteUpperValue, encodedLength[index1])

	message3 := ControlPacket{length: threeByteLength}
	encodedLength = message3.Length()

	m.Assert().Equal(threeByteCount, len(encodedLength))
	m.Assert().Equal(carryByte, encodedLength[index0])
	m.Assert().Equal(carryByte, encodedLength[index1])
	m.Assert().Equal(twoByteUpperValue, encodedLength[index2])

	message4 := ControlPacket{length: fourByteLength}
	encodedLength = message4.Length()

	m.Assert().Equal(fourByteCount, len(encodedLength))
	m.Assert().Equal(carryByte, encodedLength[index0])
	m.Assert().Equal(carryByte, encodedLength[index1])
	m.Assert().Equal(carryByte, encodedLength[index2])
	m.Assert().Equal(fourByteUpperValue, encodedLength[index3])
}

func (m *PacketTest) TestDecodeLength() {
	message1 := ControlPacket{data: []byte{singleByteLength}}

	length, bytesUsed := message1.decodeLength()
	m.Assert().Equal(uint32(singleByteLength), length)
	m.Assert().Equal(1, bytesUsed)

	message2 := ControlPacket{data: []byte{maxByte, singleByteLength}}

	length, bytesUsed = message2.decodeLength()
	m.Assert().Equal(twoByteValue, length)
	m.Assert().Equal(2, bytesUsed)

	message3 := ControlPacket{data: []byte{maxByte, maxByte, singleByteLength}}

	length, bytesUsed = message3.decodeLength()
	m.Assert().Equal(threeByteValue, length)
	m.Assert().Equal(3, bytesUsed)

	message4 := ControlPacket{data: []byte{maxByte, maxByte, maxByte, singleByteLength}}

	length, bytesUsed = message4.decodeLength()
	m.Assert().Equal(fourByteValue, length)
	m.Assert().Equal(4, bytesUsed)

	messageError := ControlPacket{data: []byte{maxByte, maxByte, maxByte, maxByte, singleByteLength}}

	length, bytesUsed = messageError.decodeLength()
	m.Assert().Equal(errorValue, length)
	m.Assert().Equal(0, bytesUsed)
}

func (m *PacketTest) TestDecodePacket() {
	data := []byte{}

	packet, err := DecodePacket(data)

	// Expected to error due to bad length/packet
	m.Assert().NotNil(err)
	m.Assert().Nil(packet)

	data = []byte{0}

	packet, err = DecodePacket(data)

	// Expected to error due to bad length/packet
	m.Assert().NotNil(err)
	m.Assert().Nil(packet)

	data = []byte{0x0, 0x0, 0x0}

	packet, err = DecodePacket(data)

	// Expected to error due to bad packet type
	m.Assert().NotNil(err)
	m.Assert().Nil(packet)
}

func (m *PacketTest) TestDecodeConnectPacket() {
	data := []byte{Connect << ctrlTypeShift}

	packet, err := DecodePacket(data)

	// Expected to error due to bad length/packet
	m.Assert().NotNil(err)
	m.Assert().Nil(packet)
}

func TestPacketSuite(t *testing.T) {
	suite.Run(t, new(PacketTest))
}

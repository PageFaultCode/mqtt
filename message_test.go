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

type MessageTest struct {
	suite.Suite
}

func (m *MessageTest) TestLength() {
	message1 := message{length: singleByteLength}

	encodedLength := message1.Length()

	m.Assert().Equal(singleByteCount, len(encodedLength))

	message2 := message{length: twoByteLength}
	encodedLength = message2.Length()

	m.Assert().Equal(twoByteCount, len(encodedLength))
	m.Assert().Equal(carryByte, encodedLength[index0])
	m.Assert().Equal(twoByteUpperValue, encodedLength[index1])

	message3 := message{length: threeByteLength}
	encodedLength = message3.Length()

	m.Assert().Equal(threeByteCount, len(encodedLength))
	m.Assert().Equal(carryByte, encodedLength[index0])
	m.Assert().Equal(carryByte, encodedLength[index1])
	m.Assert().Equal(twoByteUpperValue, encodedLength[index2])

	message4 := message{length: fourByteLength}
	encodedLength = message4.Length()

	m.Assert().Equal(fourByteCount, len(encodedLength))
	m.Assert().Equal(carryByte, encodedLength[index0])
	m.Assert().Equal(carryByte, encodedLength[index1])
	m.Assert().Equal(carryByte, encodedLength[index2])
	m.Assert().Equal(fourByteUpperValue, encodedLength[index3])
}

func (m *MessageTest) TestDecodeLength() {
	message1 := message{data: []byte{singleByteLength}}

	m.Assert().Equal(uint32(singleByteLength), message1.decodeLength())

	message2 := message{data: []byte{maxByte, singleByteLength}}

	m.Assert().Equal(twoByteValue, message2.decodeLength())

	message3 := message{data: []byte{maxByte, maxByte, singleByteLength}}

	m.Assert().Equal(threeByteValue, message3.decodeLength())

	message4 := message{data: []byte{maxByte, maxByte, maxByte, singleByteLength}}

	m.Assert().Equal(fourByteValue, message4.decodeLength())

	messageError := message{data: []byte{maxByte, maxByte, maxByte, maxByte, singleByteLength}}

	m.Assert().Equal(errorValue, messageError.decodeLength())
}

func TestMessageSuite(t *testing.T) {
	suite.Run(t, new(MessageTest))
}

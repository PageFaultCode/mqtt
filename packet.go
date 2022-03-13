// Package mqtt is the main encompassing mqtt package
package mqtt

import "fmt"

const (
	messageLengthMask       = 127
	messageLengthMax        = 128
	maxByteMultiplierCount  = 3
	ctrlTypePosition        = 0
	ctrlTypeShift           = 4
	ctrlFlagMask            = 0xF0
	ctrlPacketMinimumLength = 2
)

type ControlPacket struct {
	ctrlType byte
	flags    byte
	length   uint32
	data     []byte
}

// EncodePacket encodes the data for sending onto the wire
func (cp *ControlPacket) EncodePacket() ([]byte, error) {
	data := make([]byte, ctrlPacketMinimumLength) // we need enough room for ctrl/flags and length

	data[ctrlTypePosition] = cp.ctrlType<<ctrlTypeShift | cp.flags

	return data, nil
}

// DecodePacket decodes an incoming packet from the wire
func DecodePacket(data []byte) (*ControlPacket, error) {
	if len(data) <= ctrlPacketMinimumLength {
		return nil, fmt.Errorf("invalid mqtt packet: +%v", data)
	}

	// setup type and flags
	// length and data will be decoded as we know more
	ctrlPacket := ControlPacket{
		ctrlType: data[ctrlTypePosition] >> ctrlTypeShift,
		flags:    data[ctrlTypePosition] & ctrlFlagMask,
	}

	// Length is byte 2 -> 5 depending on length

	switch ctrlPacket.ctrlType {
	case reserved:
		return nil, fmt.Errorf("invalid control packet type: %v", data[ctrlTypePosition]>>ctrlTypeShift)
	case Connect:
		ctrlPacket.flags = 0
	case ConnectAck:
		ctrlPacket.flags = 0
	}

	return &ctrlPacket, nil
}

// Length returns the byte encoded length of the message
func (cp *ControlPacket) Length() []byte {
	encodedLength := make([]byte, 0)
	lengthToEncode := cp.length

	for {
		encodedByte := lengthToEncode % messageLengthMax
		lengthToEncode /= messageLengthMax
		if lengthToEncode > 0 {
			encodedByte |= messageLengthMax
		}
		encodedLength = append(encodedLength, byte(encodedByte))

		// if length to encode gets to 0
		// we are done
		if lengthToEncode == 0 {
			break
		}
	}
	return encodedLength
}

func (cp *ControlPacket) decodeLength() (uint32, int) {
	multiplier := 1
	index := 0
	var length uint32

	for {
		encodedByte := cp.data[index]
		length += uint32(encodedByte&byte(messageLengthMask)) * uint32(multiplier)
		if encodedByte&messageLengthMax == 0 {
			break
		}
		index++
		if index > maxByteMultiplierCount {
			return 0, 0
		}
		multiplier *= messageLengthMax
	}
	return length, index + 1
}

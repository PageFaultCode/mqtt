// Package mqtt is the main encompassing mqtt package
package mqtt

const (
	messageLengthMask      = 127
	messageLengthMax       = 128
	maxByteMultiplierCount = 3
)

type fixedHeader struct {
	// 8 bits for ctrl type and flag
	// upper 4 is type, lower 4 is flags
	ctrlTypeFlags   byte
	remainingLength []byte
}

type numberedPacket struct {
	packetNumber uint16
}

type message struct {
	length uint32
	data   []byte
}

// Length returns the byte encoded length of the message
func (m *message) Length() []byte {
	encodedLength := make([]byte, 0)
	lengthToEncode := m.length

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

func (m *message) decodeLength() uint32 {
	multiplier := 1
	index := 0
	var length uint32

	for {
		encodedByte := m.data[index]
		length += uint32(encodedByte&byte(messageLengthMask)) * uint32(multiplier)
		if encodedByte&messageLengthMax == 0 {
			break
		}
		index++
		if index > maxByteMultiplierCount {
			length = 0
			break
		}
		multiplier *= messageLengthMax
	}
	return length
}

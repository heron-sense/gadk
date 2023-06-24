package extension

const (
	Digits = `12345678AFHKSWXY`
)

var (
	Places = []byte{
		0b00000000, 0b00000000, 0b00000000, 0b00000000,
		0b00000000, 0b00000000, 0b11111110, 0b00000001,
		0b01000010, 0b00001001, 0b10001000, 0b00000011,
		0b01000010, 0b00001001, 0b10001000, 0b00000011,
	}
)

func EncodeHgiPadding(data [20]byte) []byte {
	var remaining, bits uint32

	padding := make([]byte, 0, len(data)*8/5)
	pos := 0

	for n := 0; n < len(data) || remaining > 0; {
		if remaining >= 5 {
			remaining -= 5
			val := bits >> remaining
			bits &= 0xFFFF >> (16 - remaining)
			padding = append(padding, Digits[val])
			pos++
		} else {
			bits = (bits << 8) | uint32(data[n])
			remaining += 8
			n++
		}
	}

	return padding
}

func He32ofRaw(raw []byte) []byte {
	he32 := make([]byte, 0, len(raw)*2)
	var remaining uint8
	var bits uint16
	for next := 0; next < len(raw) || remaining > 0; {
		switch {
		case remaining >= 5:
			remaining -= 5
			he32 = append(he32, Digits[bits>>(remaining)])
			bits = bits & (0xffff >> (16 - remaining))
		case next < len(raw):
			bits = (bits << 8) | uint16(raw[next])
			remaining += 8
			next++
		default:
			bits = bits << (5 - remaining)
			he32 = append(he32, Digits[bits>>(remaining)])
			remaining = 0
			bits = 0
		}
	}
	return he32
}

func IsHe32Uuid(uuid string) bool {
	raw := []byte(uuid)
	for _, chr := range raw {
		if chr < 0x80 {
			if Places[chr/8]&(1<<(chr%8)) != 0 {
				continue
			}
		}
		if chr == '-' || chr == '.' {
			continue
		}
		return false
	}
	return true
}

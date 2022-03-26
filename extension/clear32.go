package extension

import (
	"fmt"
)

const (
	Digits = "01236789CLSW45BDGJKMNPQTXYAEFHIR"
	Values = "0123<=4567-------J>8?KL@MNAB9CD-EFO:G--;HI"
)

func He32ofUint64(val uint64) string {
	if val == 0 {
		return "0"
	}

	const capacity = 19
	rep := make([]byte, capacity+1)
	occupied := 0

	for val > 0 {
		remainder := val & 0x1f
		rep[capacity-occupied] = Digits[remainder]
		occupied++
		val = val >> 5
	}
	return string(rep[capacity+1-occupied:])
}

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

func ParseUint64FromHe32(rep string) (uint64, error) {
	var val uint64
	for loc, chr := range rep {
		if chr < '0' || chr > 'Y' {
			return 0, fmt.Errorf("invalid byte[%c] at location[%d] encountered", chr, loc)
		}

		tmp := Values[chr-'0'] - '0'
		if tmp < 0 {
			return 0, fmt.Errorf("invalid byte[%c] at location[%d] encountered", chr, loc)
		}

		val = val<<5 + uint64(tmp)
	}

	return val, nil
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
	var flags uint32
	raw := []byte(uuid)
	for idx := 0; idx < len(raw); idx++ {
		chr := raw[idx]
		if chr == '-' {
			continue
		}
		if chr < '0' || chr > 'Y' {
			return false
		}

		if tmp := Values[chr-'0'] - '0'; tmp < 0 {
			return false
		} else {
			flags |= 1 << tmp
		}
	}
	occurs := 0
	for ; flags != 0; flags = flags >> 1 {
		if flags&1 != 0 {
			occurs++
		}
	}
	return len(raw) == 32 && occurs > 1
}

package facilities

import "unsafe"

func str2bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func bytes2str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

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

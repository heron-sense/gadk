package facilities

import (
	"fmt"
	"github.com/heron-sense/gadk/extension"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsHe32Uuid(t *testing.T) {
	for chr := 0; chr < 0x80; chr++ {
		rep := string([]byte{byte(chr)})
		result := extension.IsHe32Uuid(rep)
		switch chr {
		case '1', '2', '3', '4', '5', '6', '7', '8':
		case 'A', 'F', 'H', 'K', 'S', 'W', 'X', 'Y':
		case 'a', 'f', 'h', 'k', 's', 'w', 'x', 'y':
			assert.True(t, result, "expecting when rep=%c return true", byte(chr))
		default:
			assert.False(t, result, "expecting when rep=%c return false", byte(chr))
		}
	}
}

func BenchmarkIsHe32Uuid(b *testing.B) {
	res := false
	rep := string("12345678AFHKSWXY")
	for n := 0; n < b.N; n++ {
		res = extension.IsHe32Uuid(rep)
	}
	fmt.Printf("res:=%v\n", res)
}

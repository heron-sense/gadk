package maf

import (
	"bytes"
)

type FieldSymbol = uint8

const (
	DataTypeNil  FieldSymbol = 'n'
	DataTypeBool FieldSymbol = 'b' //sample content:= "publish=b1:t"
	DataTypeInt  FieldSymbol = 'i'
	DataTypeFlt  FieldSymbol = 'f'
	DataTypeRaw  FieldSymbol = 'r'
	DataTypeObj  FieldSymbol = 'o'
	DataTypeLst  FieldSymbol = 'L'
)

//64进制: . _ 0-9 a-z A-Z
type _octTreeNode struct {
	PathAbbr uint64
	child    [8]*_octTreeNode //mod by 8
	Symbol   FieldSymbol
	abbrLen  uint8
	Data     interface{}
}

/**
 *
 */
func track(str string) ([]byte, bool) {
	track := make([]byte, 0, 2*len(str))

	for _, chr := range []byte(str) {
		switch {
		case chr >= 'a':
			if chr <= 'z' {
				chr = chr - 'a' + 11
			} else if chr < 0x7f {
				//['{','|','}','~'] => [[63,7],[63,8],[63,9],[63,10]]
				track = append(track, 0x3f)
				chr = chr - 'z' + 6
			} else {
				return track, false
			}
			break
		case chr <= '9':
			if chr >= '/' {
				chr = chr - '/'
			} else if chr >= ' ' {
				//[32,...,46] => [[63,11],...,[63,25]]
				track = append(track, 0x3f)
				chr = chr - ' ' + 11
			} else {
				return track, false
			}
			break
		case chr >= 'A' && chr <= 'Z':
			chr = chr - 'A' + 37
			break
		default:
			//[':',';','<','=','>','?','@'] => [[63,0],[63,1],[63,2],[63,3],[63,4],[63,5],[63,6]]
			//['[','\\',']','^','_', '`'] => [[63,33],[63,34],[63,35],[63,36],[63,37],[63,38]]
			track = append(track, 0x3f)
			chr -= ':'
		}
		track = append(track, chr)
	}

	return track, true
}

func (tree *_octTreeNode) Set(key string, data interface{}) bool {
	kBytes, ok := track(key)
	if !ok {
		return false
	}
	pos := 0

	for ptr, match := tree, uint8(0); ; {
		switch remain := ptr.abbrLen - match; {
		case remain > 3:
			if kBytes[pos] == uint8(ptr.PathAbbr>>match) {
				match += 6
			} else {
				//patch
			}
		case remain == 3:
		default:

		}
	}
}

func (tree *_octTreeNode) Serialize(track []byte, buf *bytes.Buffer) {
	if tree.Symbol != 0 {
		buf.Write(track)
	}

	track = append(track, ' ')
	for idx := 0; idx < len(tree.child); idx++ {
		if sub := tree.child[idx]; sub != nil {
			track[len(track)-1] = 'a'
			sub.Serialize(track, buf)
		}
	}
}

func (tree *_octTreeNode) Del(key string) {

}

func (tree *_octTreeNode) Find(key string) (interface{}, bool) {
	return nil, false
}

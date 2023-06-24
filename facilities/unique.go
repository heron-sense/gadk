package facilities

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

func HasDuplicatedElem(arr []int64) bool {
	left := 0
	right := 0
	var tmp int64
	for pos := len(arr)>>1 - 1; pos >= 0; pos-- {
		for parent := pos; parent < len(arr); {
			left = parent<<1 + 1
			right = left + 1

			var pp = -1
			if right < len(arr) {
				delta := arr[left] - arr[right]
				if delta < 0 {
					if arr[parent] < arr[right] {
						pp = right
					}
				} else if delta == 0 {
					return true
				} else {
					if arr[parent] < arr[left] {
						pp = left
					}
				}
			} else if left < len(arr) {
				if arr[parent] < arr[left] {
					pp = left
				}
			}

			if pp >= 0 {
				tmp = arr[parent]
				arr[parent] = arr[pp]
				arr[pp] = tmp
				parent = pp
			} else {
				break
			}
		}
	}

	for pos := len(arr) - 1; pos >= 0; pos-- {
		if pos+1 < len(arr) && arr[pos+1] == arr[0] {
			return true
		}
		tmp = arr[pos]
		arr[pos] = arr[0]
		arr[0] = tmp

		for parent := 0; parent <= pos; {
			left = (parent << 1) + 1
			right = left + 1

			var pp = -1
			if right <= pos-1 {
				delta := arr[left] - arr[right]
				if delta < 0 {
					if arr[parent] < arr[right] {
						pp = right
					}
				} else if delta == 0 {
					return true
				} else {
					if arr[parent] < arr[left] {
						pp = left
					}
				}
			} else if left <= pos-1 {
				if arr[parent] < arr[left] {
					pp = left
				}
			}

			if pp >= 0 {
				tmp = arr[parent]
				arr[parent] = arr[pp]
				arr[pp] = tmp
				parent = pp
			} else {
				break
			}
		}
	}

	return false
}

func HasDuplicateInt64(slice []int64) bool {
	if len(slice) == 0 {
		return false
	}
	set := make(map[int64]bool)
	for _, n := range slice {
		if set[n] {
			return true
		}
		set[n] = true
	}
	return false
}

func HasDuplicatedElemCondensed(arr []int64) bool {
	var leftChildPos, right int64

	size := len(arr)
	nextInitPos := int64(size>>1) - 1
	adjustPos := int64(size)
	var pp int64 = 0
	var adjustLen int64 = 0

	for nextInitPos >= 0 || adjustPos > 0 {
		if nextInitPos < 0 {
			adjustPos--
			pp = arr[adjustPos]
			arr[adjustPos] = arr[0]
			arr[0] = pp
			pp = 0
			adjustLen = adjustPos - 1
		} else {
			pp = nextInitPos
			adjustLen = int64(size)
			nextInitPos--
		}

		parent := pp
		for leftChildPos = parent<<1 + 1; leftChildPos < adjustLen; leftChildPos = parent<<1 + 1 {
			pval := arr[parent]
			pp = parent
			right = leftChildPos + 1

			if pval < arr[leftChildPos] {
				pp = leftChildPos
			} else if pval == arr[leftChildPos] {
				return true
			}
			if right < adjustLen {
				if arr[pp] < arr[right] {
					pp = right
				} else if arr[right] == arr[leftChildPos] || pval == arr[right] {
					return true
				}
			}

			if pp != parent {
				arr[parent] = arr[pp]
				arr[pp] = pval
				parent = pp
				continue
			}
			break
		}
	}

	return false
}

func main() {
	//arr1 := []int64{6, 3, 55, 6, 7, 6, 8, 5, 3, 3, 2, 2, 5, 456, 6, 7, 876, 8, 9, 9, 90, 0, 0}
	//arr2 := []int64{6, 3, 55, 6, 7, 6, 8, 5, 3, 3, 2, 2, 5, 456, 6, 7, 876, 8, 9, 9, 90, 0, 0}
	//arr := []int64{6, 3, 55, 6, 7, 6, 8, 5, 0}

	var v11Elapse int64 = 0
	var v22Elapse int64 = 0
	var v33Elapse int64

	for n := 0; n < 200000; n++ {
		arr1 := make([]int64, 0)
		arr2 := make([]int64, 0)
		arr3 := make([]int64, 0)

		for x := 0; x < 50; x++ {
			r := int64(rand.Uint32() % 100)
			arr1 = append(arr1, r)
			arr2 = append(arr2, r)
			arr3 = append(arr3, r)
		}

		sort.Slice(arr1, func(i, j int) bool {
			return arr1[i] < arr1[j]
		})

		sort.Slice(arr2, func(i, j int) bool {
			return arr2[i] < arr2[j]
		})

		sort.Slice(arr3, func(i, j int) bool {
			return arr3[i] < arr3[j]
		})

		begin := time.Now().UnixNano()
		v11 := HasDuplicateInt64(arr1)
		v11Elapse += time.Now().UnixNano() - begin

		begin = time.Now().UnixNano()
		v22 := HasDuplicatedElem(arr2)
		v22Elapse += time.Now().UnixNano() - begin

		begin = time.Now().UnixNano()
		v33 := HasDuplicatedElemCondensed(arr3)
		v33Elapse += time.Now().UnixNano() - begin

		if v11 != v33 || v22 != v33 {
			panic("ff")
		}
	}

	fmt.Printf("v11Elapse   %+v\n", v11Elapse/1e3)
	fmt.Printf("v22Elapse   %+v\n", v22Elapse/1e3)
	fmt.Printf("v33Elapse   %+v\n", v33Elapse/1e3)
}

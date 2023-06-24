package facilities

func HeapSort(arr []int64) {
	left := 0
	right := 0
	var tmp int64
	for pos := len(arr)>>1 - 1; pos >= 0; pos-- {
		for parent := pos; parent < len(arr); {
			left = parent<<1 + 1
			right = left + 1

			var pp = -1
			if right < len(arr) {
				if arr[left] < arr[right] {
					if arr[parent] < arr[right] {
						pp = right
					}
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
		tmp = arr[pos]
		arr[pos] = arr[0]
		arr[0] = tmp

		for parent := 0; parent <= pos; {
			left = (parent << 1) + 1
			right = left + 1

			var pp = -1
			if right <= pos-1 {
				if arr[left] < arr[right] {
					if arr[parent] < arr[right] {
						pp = right
					}
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
}

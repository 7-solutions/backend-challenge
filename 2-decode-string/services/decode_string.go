package services

func DecodeV2(s string) []int {

	if s == "" {
		return []int{}
	}

	left := []int{1, 0}
	right := []int{0, 1}

	result := []int{}

	lenInput := len(s)
	prev := -1

	for _, char := range s {

		if lenInput == 1 {
			if char == 'L' {
				result = append(result, left...)
			} else if char == 'R' {
				result = append(result, right...)
			}
		} else {
			if char == 'L' {
				if prev == -1 {
					prev = 1
				} else {
					prev++
				}
				result = append(result, prev)
				lenInput--
			} else if char == 'R' {
				if prev == -1 {
					prev = 0
				} else {
					prev--
				}
				result = append(result, prev)
				lenInput--
			} else if char == '=' {
				result = append(result, result[len(result)-1])
			}

		}
	}
	return result
}

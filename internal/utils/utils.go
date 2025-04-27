package utils

func IsValidLuhn(number string) bool {
	var sum int
	alt := false
	for i := len(number) - 1; i >= 0; i-- {
		n := int(number[i] - '0')
		if alt {
			n *= 2
			if n > 9 {
				n -= 9
			}
		}
		sum += n
		alt = !alt
	}
	return sum%10 == 0
}

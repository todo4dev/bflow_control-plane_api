package util

func LuhnValid(number string) bool {
	var sum int
	double := false

	for index := len(number) - 1; index >= 0; index-- {
		digit := int(number[index] - '0')

		if double {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}

		sum += digit
		double = !double
	}

	return sum%10 == 0
}

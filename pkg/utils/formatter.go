package utils

import "strconv"

func ParseInt(num string) int64 {
	n, err := strconv.ParseInt(num, 10, 64)
	if err != nil {
		return 0
	}
	return n
}
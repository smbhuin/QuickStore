package main

import "strconv"

func Stoi(s string, defaultValue int) int {
	var n int = defaultValue
	if s != "" {
		if parsed, err := strconv.Atoi(s); err == nil {
			n = parsed
		}
	}
	return n
}

func StoiStrict(s string) (int, error) {
	return strconv.Atoi(s)
}

func rangeBound(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

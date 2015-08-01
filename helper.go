package main

import "strings"

func toSnake(str string) string {
	if len(str) <= 1 {
		return strings.ToLower(str)
	}

	runes := []rune(str)
	result := []rune{runes[0]}
	for i, l := range runes[1 : len(runes)-1] {
		if isUpper(l) {
			if !isUpper(runes[i]) || !isUpper(runes[i+2]) {
				result = append(result, rune('_'))
			}
		}

		result = append(result, l)
	}

	last := runes[len(runes)-1]
	if isUpper(last) && !isUpper(runes[len(runes)-2]) {
		result = append(result, rune('_'))
	}
	result = append(result, last)

	return strings.ToLower(string(result))
}

func isUpper(l rune) bool {
	return rune('A') <= l && l <= rune('Z')
}

func camelCase(str string) string {
	return strings.ToLower(str[:1]) + str[1:]
}

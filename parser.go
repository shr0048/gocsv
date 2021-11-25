package gocsv

import (
	"strings"
)

func ParseLine(row []string, separator string) []string {
	var findComma = 0
	var substr string
	var realSplit []string

	for _, element := range row {
		comma := strings.Contains(element, "\"")
		if comma == true || findComma != 0 {
			if comma == true && strings.Count(element, "\"")%2 == 1 {
				findComma++
			} else if comma == true && strings.Count(element, "\"")%2 == 0 && findComma != 1 {
				realSplit = append(realSplit, element)
				continue
			}

			if findComma == 1 {
				substr = substr + element + separator
			} else {
				substr = substr + element
			}

			if findComma == 2 {
				realSplit = append(realSplit, substr)

				substr = ""
				findComma = 0
			}
		} else if findComma == 0 {
			realSplit = append(realSplit, element)
		}
	}
	return realSplit
}

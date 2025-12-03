package day2

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func Solve(part int, isTest bool) (any, error) {
	f := "day_2/input.txt"
	if isTest {
		f = "day_2/input-test.txt"
	}

	body, err := os.ReadFile(f)
	if err != nil {
		return nil, err
	}

	switch part {
	case 1:
		return partOne(string(body))
	case 2:
		return partTwo(string(body))
	}

	return nil, fmt.Errorf("part should be only 1 or 2")
}

func partOne(data string) (any, error) {
	idRanges, err := parseProductIDRanges(data)
	if err != nil {
		return nil, err
	}

	ans := 0
	for _, idRange := range idRanges {
		for _, invalidID := range idRange.PartOneInvalidIDs() {
			ans += invalidID
		}
	}

	return ans, nil
}

func partTwo(data string) (any, error) {
	idRanges, err := parseProductIDRanges(data)
	if err != nil {
		return nil, err
	}

	ans := 0
	for _, idRange := range idRanges {
		for _, invalidID := range idRange.PartTwoInvalidIDs() {
			ans += invalidID
		}
	}

	return ans, nil
}

type ProductIDRange struct {
	First int
	Last  int
}

func (idRange *ProductIDRange) PartOneInvalidIDs() []int {
	invalidIDs := []int{}
	for id := idRange.First; id <= idRange.Last; id++ {
		if repeatSequenceOfDigitsTwice(id) {
			invalidIDs = append(invalidIDs, id)
		}
	}
	return invalidIDs
}

func (idRange *ProductIDRange) PartTwoInvalidIDs() []int {
	invalidIDs := []int{}
	for id := idRange.First; id <= idRange.Last; id++ {
		if repeatSequenceOfDigitsAtLeastTwice(id) {
			invalidIDs = append(invalidIDs, id)
		}
	}
	return invalidIDs
}

func repeatSequenceOfDigitsTwice(number int) bool {
	s := strconv.Itoa(number)
	lenght := len(s)

	if lenght%2 != 0 {
		return false
	}

	halfLen := lenght / 2
	return s[:halfLen] == s[halfLen:]
}

func repeatSequenceOfDigitsAtLeastTwice(number int) bool {
	s := strconv.Itoa(number)
	lenght := len(s)

	for size := 1; size <= lenght/2; size++ {
		if lenght%size != 0 {
			continue
		}
		pattern := s[:size]

		repeated := strings.Repeat(pattern, lenght/size)

		if repeated == s {
			return true
		}
	}

	return false
}

func parseProductIDRanges(data string) ([]ProductIDRange, error) {
	var ranges []ProductIDRange
	for _, rangeStr := range strings.Split(data, ",") {
		firstStr, lastStr, found := strings.Cut(rangeStr, "-")
		if !found {
			return nil, fmt.Errorf("invalid range: %v", rangeStr)
		}

		first, err := strconv.Atoi(firstStr)
		if err != nil {
			return nil, fmt.Errorf("invalid first number in range %v: %v", rangeStr, err)
		}

		last, err := strconv.Atoi(lastStr)
		if err != nil {
			return nil, fmt.Errorf("invalid last number in range %v: %v", rangeStr, err)
		}

		ranges = append(ranges, ProductIDRange{
			First: first,
			Last:  last,
		})
	}

	return ranges, nil
}

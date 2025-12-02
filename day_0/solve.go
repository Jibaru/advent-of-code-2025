package day0

import (
	"fmt"
	"os"
)

func Solve(part int, isTest bool) (any, error) {
	f := "day_0/input.txt"
	if isTest {
		f = "day_0/input-test.txt"
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

func partOne(_ string) (any, error) {
	return "part 1 ok", nil
}

func partTwo(_ string) (any, error) {
	return "part 2 ok", nil
}

package day0

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func Solve(part int, isTest bool) (any, error) {
	f := "day_1/input.txt"
	if isTest {
		f = "day_1/input-test.txt"
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
	rotations, err := parseRotations(data)
	if err != nil {
		return nil, err
	}

	position := 50
	zeroTimes := 0
	for _, rotation := range rotations {
		position = newPosition(position, rotation)
		if position == 0 {
			zeroTimes++
		}
	}

	return zeroTimes, nil
}

func partTwo(data string) (any, error) {
	rotations, err := parseRotations(data)
	if err != nil {
		return nil, err
	}

	position := 50
	zeroTimes := 0
	for _, rotation := range rotations {
		zeroTimes += countsWhenRotationPassesFromZero(position, rotation)
		position = newPosition(position, rotation)
	}

	return zeroTimes, nil
}

type Direction int

type Rotation struct {
	Dir   Direction
	Times int
}

const (
	Left Direction = iota
	Right
)

func parseRotations(data string) ([]Rotation, error) {
	var rotations []Rotation

	for _, line := range strings.Split(data, "\n") {
		numberStr := line[1:]
		dirStr := line[:1]

		r := Rotation{}
		if dirStr == "L" {
			r.Dir = Left
		} else {
			r.Dir = Right
		}

		v, err := strconv.Atoi(numberStr)
		if err != nil {
			return nil, err
		}
		r.Times = v

		rotations = append(rotations, r)
	}

	return rotations, nil
}

func newPosition(current int, rotation Rotation) int {
	if rotation.Dir == Right {
		return (current + rotation.Times) % 100
	}

	// Left
	return (current - rotation.Times + 100) % 100
}

func countsWhenRotationPassesFromZero(current int, rotation Rotation) int {
	counts := 0
	if rotation.Dir == Right {
		for i := 1; i <= rotation.Times; i++ {
			if (current+i)%100 == 0 {
				counts++
			}
		}
	} else {
		for i := 1; i <= rotation.Times; i++ {
			if (current-i+100)%100 == 0 {
				counts++
			}
		}
	}
	return counts
}

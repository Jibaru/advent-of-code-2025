package day3

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func Solve(part int, isTest bool) (any, error) {
	f := "day_3/input.txt"
	if isTest {
		f = "day_3/input-test.txt"
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
	banks, err := parseBank(data)
	if err != nil {
		return nil, err
	}

	ans := 0
	for _, bank := range banks {
		ans += largestVoltage(bank)
	}

	return ans, nil
}

func partTwo(data string) (any, error) {
	banks, err := parseBank(data)
	if err != nil {
		return nil, err
	}

	ans := 0
	for _, bank := range banks {
		ans += largestVoltageK(bank, 12)
	}

	return ans, nil
}

type Battery = int

type Bank []Battery

func parseBank(data string) ([]Bank, error) {
	var banks []Bank
	for _, line := range strings.Split(data, "\n") {
		var bank Bank
		for _, batteryStr := range line {
			b, err := strconv.Atoi(string(batteryStr))
			if err != nil {
				return nil, err
			}
			bank = append(bank, Battery(b))
		}
		banks = append(banks, bank)
	}
	return banks, nil
}

func largestVoltage(bank Bank) int {
	m1 := -1
	idx := -1
	for i := 0; i < len(bank)-1; i++ {
		if bank[i] > m1 {
			m1 = bank[i]
			idx = i
		}
	}

	m2 := -1
	for i := idx + 1; i < len(bank); i++ {
		m2 = max(m2, bank[i])
	}

	return m1*10 + m2
}

func largestVoltageK(bank Bank, size int) int {
	numberOfBatteries := len(bank)
	if size >= numberOfBatteries {
		return bankToInt(bank)
	}

	toRemove := numberOfBatteries - size
	stack := NewStack[Battery](numberOfBatteries)

	for _, curr := range bank {
		for toRemove > 0 && !stack.IsEmpty() && stack.Top() < curr {
			stack.Pop()
			toRemove--
		}

		stack.Push(curr)
	}

	stack.TrimLast(toRemove)

	slice := stack.Slice()
	return bankToInt(slice[:size])
}

// bankToInt converts a bank to number, like {9,8,7} to 987
func bankToInt(b Bank) int {
	var result int = 0
	for _, d := range b {
		result = result*10 + int(d)
	}
	return result
}

type Stack[T any] struct {
	data []T
}

func NewStack[T any](capacity int) *Stack[T] {
	return &Stack[T]{data: make([]T, 0, capacity)}
}

func (s *Stack[T]) Push(v T) {
	s.data = append(s.data, v)
}

func (s *Stack[T]) Pop() T {
	var zero T
	if len(s.data) == 0 {
		return zero
	}
	last := s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]
	return last
}

func (s *Stack[T]) Top() T {
	var zero T
	if len(s.data) == 0 {
		return zero
	}
	return s.data[len(s.data)-1]
}

func (s *Stack[T]) Len() int {
	return len(s.data)
}

func (s *Stack[T]) IsEmpty() bool {
	return len(s.data) == 0
}

// TrimLast removes the last n elements
func (s *Stack[T]) TrimLast(n int) {
	if n <= 0 || n > len(s.data) {
		return
	}
	s.data = s.data[:len(s.data)-n]
}

// Slice returns a copy of the internal slice
func (s *Stack[T]) Slice() []T {
	out := make([]T, len(s.data))
	copy(out, s.data)
	return out
}

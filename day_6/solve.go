package day6

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func Solve(part int, isTest bool) (any, error) {
	f := "day_6/input.txt"
	if isTest {
		f = "day_6/input-test.txt"
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
	ws, err := parseWorksheet(data)
	if err != nil {
		return nil, err
	}

	total := 0
	for _, p := range ws {
		res := p.Ans()
		total += res
	}
	return total, nil
}

func partTwo(data string) (any, error) {
	ws, err := parseWorksheetCephalopod(data)
	if err != nil {
		return nil, err
	}

	total := 0
	for _, p := range ws {
		res := p.Ans()
		total += res
	}
	return total, nil
}

type Problem struct {
	Numbers  []int
	Operator rune
}

func (p Problem) Ans() int {
	switch p.Operator {
	case '*':
		ans := 1
		for _, n := range p.Numbers {
			ans *= n
		}
		return ans
	case '+':
		ans := 0
		for _, n := range p.Numbers {
			ans += n
		}
		return ans
	}
	return 0
}

func parseWorksheet(data string) ([]Problem, error) {
	data = strings.TrimSpace(data)

	lines := strings.Split(data, "\n")
	opLine := strings.TrimSpace(lines[len(lines)-1]) // last line

	opRe := regexp.MustCompile(`[\+\*]`)
	operators := opRe.FindAllString(opLine, -1)

	matrix := parseNumbersStr(lines[:len(lines)-1])
	matrix = append(matrix, operators)
	matrix = transposeMatrix(matrix)

	var problems []Problem
	for _, row := range matrix {
		p := Problem{}
		for _, s := range row {
			if s == "*" || s == "+" {
				p.Operator = rune(s[0])
			} else {
				n, err := strconv.Atoi(s)
				if err != nil {
					return nil, err
				}
				p.Numbers = append(p.Numbers, n)
			}
		}
		problems = append(problems, p)
	}

	return problems, nil
}

func parseNumbersStr(lines []string) [][]string {
	numRe := regexp.MustCompile(`\d+`)
	allNums := [][]string{}
	for _, line := range lines {
		nums := numRe.FindAllString(line, -1)
		allNums = append(allNums, nums)
	}
	return allNums
}

func transposeMatrix[T any](matrix [][]T) [][]T {
	if len(matrix) == 0 {
		return nil
	}
	rows := len(matrix)
	cols := len(matrix[0])
	out := make([][]T, cols)
	for c := range cols {
		col := make([]T, rows)
		for r := range rows {
			col[r] = matrix[r][c]
		}
		out[c] = col
	}
	return out
}

func parseWorksheetCephalopod(input string) ([]Problem, error) {
	lines := strings.Split(strings.TrimRight(input, "\n"), "\n")
	if len(lines) < 2 {
		return nil, fmt.Errorf("invalid input")
	}

	opLine := lines[len(lines)-1]
	dataLines := lines[:len(lines)-1]

	cuts := computeSlicesFromOperators(opLine)

	var matrix [][]string
	for _, ln := range dataLines {
		row := sliceByCuts(ln, cuts)
		matrix = append(matrix, row)
	}

	inv := transposeMatrix(matrix)

	// extract operators
	opSlices := sliceByCuts(opLine, cuts)
	var ops []rune
	for _, s := range opSlices {
		s = strings.TrimSpace(s)
		if len(s) == 1 {
			ops = append(ops, rune(s[0]))
		} else {
			ops = append(ops, '+') // default
		}
	}

	// to problem struct
	var out []Problem
	for i, row := range inv {
		out = append(out, cephalopodToProblem(row, ops[i]))
	}

	return out, nil
}

type Cut struct {
	Start int
	End   int
}

func computeSlicesFromOperators(line string) []Cut {
	var cuts []Cut
	var idxs []int

	for i, r := range line {
		switch r {
		case '+', '*':
			idxs = append(idxs, i)
		}
	}

	// intervals [start, end)
	for i := range idxs {
		start := idxs[i]
		var end int
		if i+1 < len(idxs) {
			end = idxs[i+1]
		} else {
			end = len(line)
		}
		cuts = append(cuts, Cut{Start: start, End: end})
	}

	return cuts
}

func sliceByCuts(line string, cuts []Cut) []string {
	res := make([]string, len(cuts))
	for i, c := range cuts {
		s := c.Start
		e := c.End
		if s >= len(line) {
			res[i] = ""
		} else if e > len(line) {
			res[i] = line[s:]
		} else {
			res[i] = line[s:e]
		}
	}
	return res
}

func cephalopodToProblem(parts []string, op rune) Problem {
	// parts are the block substrings for each numeric row (operator excluded).
	// Find max width and pad-right so we can index columns.
	maxW := 0
	for _, s := range parts {
		if len(s) > maxW {
			maxW = len(s)
		}
	}
	padded := make([]string, len(parts))
	for i, s := range parts {
		if len(s) < maxW {
			s = s + strings.Repeat(" ", maxW-len(s))
		}
		padded[i] = s
	}

	// build numbers
	nums := []int{}
	for c := maxW - 1; c >= 0; c-- {
		var b strings.Builder
		for r := 0; r < len(padded); r++ { // top to bottom
			ch := padded[r][c]
			if ch != ' ' {
				b.WriteByte(ch)
			}
		}
		if b.Len() == 0 {
			// empty column must skip
			continue
		}
		n, err := strconv.Atoi(b.String())
		if err != nil {
			n = 0 // TODO: return error
		}
		nums = append(nums, n)
	}

	return Problem{Numbers: nums, Operator: op}
}

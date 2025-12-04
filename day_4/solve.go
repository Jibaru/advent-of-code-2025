package day4

import (
	"fmt"
	"os"
	"strings"
)

func Solve(part int, isTest bool) (any, error) {
	f := "day_4/input.txt"
	if isTest {
		f = "day_4/input-test.txt"
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
	grid, err := parseGrid(data)
	if err != nil {
		return nil, err
	}

	ans := 0
	for r := 0; r < len(grid); r++ {
		for c := 0; c < len(grid[r]); c++ {
			if grid[r][c] == '@' {
				if grid.IsAccesibleAt(r, c) {
					ans++
				}
			}
		}
	}

	return ans, nil
}

func partTwo(data string) (any, error) {
	grid, err := parseGrid(data)
	if err != nil {
		return nil, err
	}

	type Pos struct {
		r, c int
	}

	needsCheck := true
	ans := 0
	for needsCheck {
		positions := []Pos{}
		for r := range grid {
			for c := 0; c < len(grid[r]); c++ {
				if grid[r][c] == '@' {
					if grid.IsAccesibleAt(r, c) {
						positions = append(positions, Pos{
							r: r,
							c: c,
						})
					}
				}
			}
		}
		if len(positions) > 0 {
			ans += len(positions)
			needsCheck = true

			for _, pos := range positions {
				grid.ReplaceAt(pos.r, pos.c, '.')
			}
		} else {
			needsCheck = false
		}
	}

	return ans, nil
}

type Cell rune
type Row []Cell
type Grid []Row

func (g Grid) At(r, c int) (Cell, error) {
	if r < 0 || r >= len(g) {
		return 0, fmt.Errorf("row index out of bounds: %d", r)
	}
	if c < 0 || c >= len(g[r]) {
		return 0, fmt.Errorf("column index out of bounds: %d", c)
	}
	return g[r][c], nil
}

func (g Grid) IsAccesibleAt(r, c int) bool {
	return totalRollOfPapersAdjacent(g, r, c) < 4
}

func (g Grid) ReplaceAt(r, c int, val Cell) error {
	if r < 0 || r >= len(g) {
		return fmt.Errorf("row index out of bounds: %d", r)
	}
	if c < 0 || c >= len(g[r]) {
		return fmt.Errorf("column index out of bounds: %d", c)
	}
	g[r][c] = val
	return nil
}

func parseGrid(data string) (Grid, error) {
	var grid []Row

	for _, line := range strings.Split(strings.TrimSpace(data), "\n") {
		if len(strings.TrimSpace(line)) == 0 {
			continue
		}

		var row Row
		for _, char := range line {
			if char != '.' && char != '@' {
				return nil, fmt.Errorf("invalid cell character: %q", char)
			}
			row = append(row, Cell(char))
		}
		grid = append(grid, row)
	}

	return grid, nil
}

func totalRollOfPapersAdjacent(grid Grid, r, c int) int {
	directions := [][2]int{
		{-1, -1}, {-1, 0}, {-1, 1},
		{0, -1}, {0, 1},
		{1, -1}, {1, 0}, {1, 1},
	}

	count := 0
	for _, dir := range directions {
		nr, nc := r+dir[0], c+dir[1]
		cell, err := grid.At(nr, nc)
		if err != nil {
			continue
		}
		if cell == '@' {
			count++
		}
	}

	return count
}

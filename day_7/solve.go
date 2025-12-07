package day7

import (
	"fmt"
	"os"
	"strings"
)

func Solve(part int, isTest bool) (any, error) {
	f := "day_7/input.txt"
	if isTest {
		f = "day_7/input-test.txt"
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
	grid := parseGrid(data)
	start := grid.StartPosition()
	beams := NewUniqueQueue()
	beams.Put(start)
	splits := 0

	for !beams.IsEmpty() {
		pos, _ := beams.Pop()
		pos = pos.Down()
		if !grid.InBounds(pos) {
			continue
		}

		if grid.InSplitter(pos) {
			left, right := pos.Left(), pos.Right()
			beams.Put(left)
			beams.Put(right)
			splits++
			continue
		}

		beams.Put(pos)
	}

	return splits, nil
}

func partTwo(data string) (any, error) {
	grid := parseGrid(data)
	start := grid.StartPosition()

	// propagation tail: positions + multiplicity
	q := NewTimelineQueue()
	q.Put(start, 1)

	// how many timelines reach each cell before propagating
	counts := make(map[Pos]int)

	timelines := int(0)

	for !q.IsEmpty() {
		pos, val := q.Pop()

		// If this value had been accumulated before, add it.
		val += counts[pos]
		delete(counts, pos)

		down := pos.Down()
		if !grid.InBounds(down) {
			timelines += val
			continue
		}

		if grid.InSplitter(down) {
			left := down.Left()
			right := down.Right()

			// Only queue if it's the first time that cell has been touched
			// If there's already an accumulated value, add it later
			if q.Has(left) {
				counts[left] += val
			} else {
				q.Put(left, val)
			}
			if q.Has(right) {
				counts[right] += val
			} else {
				q.Put(right, val)
			}
			continue
		}

		// movimiento normal hacia abajo
		if q.Has(down) {
			counts[down] += val
		} else {
			q.Put(down, val)
		}
	}

	return timelines, nil
}

type Grid [][]rune

type Pos struct {
	Row int
	Col int
}

func parseGrid(data string) Grid {
	var g Grid
	for _, line := range strings.Split(data, "\n") {
		g = append(g, []rune(line))
	}
	return g
}

func (g Grid) StartPosition() Pos {
	var p Pos
	for row, line := range g {
		for col, char := range line {
			if char == 'S' {
				p.Row = row
				p.Col = col
				return p
			}
		}
	}
	return p
}

func (g Grid) InBounds(p Pos) bool {
	return p.Row >= 0 && p.Row < len(g) && p.Col >= 0
}

func (g Grid) InSplitter(p Pos) bool {
	return g.InBounds(p) && g[p.Row][p.Col] == '^'
}

func (p Pos) Down() Pos {
	return Pos{p.Row + 1, p.Col}
}

func (p Pos) Left() Pos {
	return Pos{p.Row, p.Col - 1}
}

func (p Pos) Right() Pos {
	return Pos{p.Row, p.Col + 1}
}

type UniqueQueue struct {
	elements map[Pos]bool
	queue    []Pos
}

func NewUniqueQueue() *UniqueQueue {
	return &UniqueQueue{
		elements: make(map[Pos]bool),
		queue:    make([]Pos, 0),
	}
}

func (uq *UniqueQueue) Put(pos Pos) {
	if !uq.elements[pos] {
		uq.elements[pos] = true
		uq.queue = append(uq.queue, pos)
	}
}

func (uq *UniqueQueue) Pop() (Pos, bool) {
	if len(uq.queue) == 0 {
		return Pos{}, false
	}
	pos := uq.queue[0]
	uq.queue = uq.queue[1:]
	return pos, true
}

func (uq *UniqueQueue) IsEmpty() bool {
	return len(uq.queue) == 0
}

// TimelineQueue maintains a queue of positions with associated multiplicity.
// If a position is already in the queue, it is not duplicated: it is added to the counts.
type TimelineQueue struct {
	queue    []Pos
	elements map[Pos]bool
	vals     map[Pos]int
}

func NewTimelineQueue() *TimelineQueue {
	return &TimelineQueue{
		queue:    make([]Pos, 0),
		elements: make(map[Pos]bool),
		vals:     make(map[Pos]int),
	}
}

func (t *TimelineQueue) Put(p Pos, v int) {
	if !t.elements[p] {
		t.queue = append(t.queue, p)
		t.elements[p] = true
	}
	t.vals[p] += v
}

func (t *TimelineQueue) Has(p Pos) bool {
	return t.elements[p]
}

func (t *TimelineQueue) Pop() (Pos, int) {
	if len(t.queue) == 0 {
		return Pos{}, 0
	}

	p := t.queue[0]
	t.queue = t.queue[1:]
	delete(t.elements, p)

	v := t.vals[p]
	delete(t.vals, p)

	return p, v
}

func (t *TimelineQueue) IsEmpty() bool {
	return len(t.queue) == 0
}

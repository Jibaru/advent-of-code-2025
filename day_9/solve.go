package day9

import (
	"fmt"
	"math"
	"os"
	"sort"
	"strings"
)

func Solve(part int, isTest bool) (any, error) {
	f := "day_9/input.txt"
	if isTest {
		f = "day_9/input-test.txt"
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
	points, err := parsePoints(data)
	if err != nil {
		return nil, err
	}

	lenght := len(points)
	maxArea := 0

	for i := 0; i < lenght; i++ {
		for j := i + 1; j < lenght; j++ {
			area := rectArea(points[i], points[j])
			if area > maxArea {
				maxArea = area
			}
		}
	}

	return maxArea, nil
}

func partTwo(data string) (any, error) {
	points, err := parsePoints(data)
	if err != nil {
		return nil, err
	}

	edges := buildEdges(points)
	rects := buildRects(points)

	sortEdgesBySize(edges)
	sortRectsBySize(rects)

	return findMaxVisibleRect(edges, rects), nil
}

type Point struct {
	X int
	Y int
}

func parsePoints(data string) ([]Point, error) {
	var points []Point
	for _, line := range strings.Split(data, "\n") {
		var point Point
		_, err := fmt.Sscanf(line, "%d,%d", &point.X, &point.Y)
		if err != nil {
			return nil, err
		}
		points = append(points, point)
	}
	return points, nil
}

func rectArea(p1, p2 Point) int {
	width := int(math.Abs(float64(p1.X-p2.X))) + 1
	height := int(math.Abs(float64(p1.Y-p2.Y))) + 1
	return width * height
}

type Edge struct {
	A Point
	B Point
}

type Rect struct {
	Size int
	A    Point
	B    Point
}

func buildEdges(points []Point) []Edge {
	n := len(points)
	edges := make([]Edge, 0, n)

	for i := 0; i < n; i++ {
		prev := (i - 1 + n) % n
		a, b := points[i], points[prev]

		if less(b, a) {
			a, b = b, a
		}
		edges = append(edges, Edge{A: a, B: b})
	}
	return edges
}

func buildRects(points []Point) []Rect {
	n := len(points)
	rects := make([]Rect, 0, n*n)

	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			a, b := points[i], points[j]
			if less(b, a) {
				a, b = b, a
			}
			rects = append(rects, Rect{
				Size: rectArea(a, b),
				A:    a,
				B:    b,
			})
		}
	}
	return rects
}

func sortEdgesBySize(edges []Edge) {
	sort.Slice(edges, func(i, j int) bool {
		return rectArea(edges[i].A, edges[i].B) > rectArea(edges[j].A, edges[j].B)
	})
}

func sortRectsBySize(rects []Rect) {
	sort.Slice(rects, func(i, j int) bool {
		return rects[i].Size > rects[j].Size
	})
}

func findMaxVisibleRect(edges []Edge, rects []Rect) int {
	for _, r := range rects {
		x1, y1 := r.A.X, r.A.Y
		x2, y2 := r.B.X, r.B.Y

		if y1 > y2 {
			y1, y2 = y2, y1
		}

		if !isCovered(edges, x1, x2, y1, y2) {
			return r.Size
		}
	}
	return 0
}

func isCovered(edges []Edge, x1, x2, y1, y2 int) bool {
	for _, e := range edges {
		x3, y3 := e.A.X, e.A.Y
		x4, y4 := e.B.X, e.B.Y

		if x4 > x1 && x3 < x2 && y4 > y1 && y3 < y2 {
			return true
		}
	}
	return false
}

// Simple lexicographical order: first X, then Y
func less(a, b Point) bool {
	if a.X != b.X {
		return a.X < b.X
	}
	return a.Y < b.Y
}

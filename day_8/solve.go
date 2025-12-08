package day8

import (
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

func Solve(part int, isTest bool) (any, error) {
	f := "day_8/input.txt"
	if isTest {
		f = "day_8/input-test.txt"
	}

	body, err := os.ReadFile(f)
	if err != nil {
		return nil, err
	}

	switch part {
	case 1:
		return partOne(string(body), isTest)
	case 2:
		return partTwo(string(body))
	}

	return nil, fmt.Errorf("part should be only 1 or 2")
}

func partOne(data string, isTest bool) (any, error) {
	positions, err := parsePositions(data)
	if err != nil {
		return nil, err
	}

	edges := buildEdges(positions)

	// sort by distance asc
	sort.Slice(edges, func(i, j int) bool {
		return edges[i].Distance < edges[j].Distance
	})

	limit := 1000
	if isTest {
		limit = 10
	}

	uf := NewUnionFind(len(positions))
	for i := 0; i < limit; i++ {
		uf.Union(edges[i].A, edges[i].B)
	}

	// count sizes
	counts := map[int]int{}
	for i := 0; i < len(positions); i++ {
		r := uf.Find(i)
		counts[r]++
	}

	// put all sizes in a slice
	var sizes []int
	for _, v := range counts {
		sizes = append(sizes, v)
	}

	// sort sizes desc
	sort.Slice(sizes, func(i, j int) bool {
		return sizes[i] > sizes[j]
	})

	if len(sizes) < 3 {
		return nil, fmt.Errorf("not enough circuits")
	}

	// largest three sizes
	return sizes[0] * sizes[1] * sizes[2], nil
}

func partTwo(data string) (any, error) {
	positions, err := parsePositions(data)
	if err != nil {
		return nil, err
	}

	edges := buildEdges(positions)

	// sort by distance asc
	sort.Slice(edges, func(i, j int) bool {
		return edges[i].Distance < edges[j].Distance
	})

	uf := NewUnionFind(len(positions))
	var x1, x2 int
	for _, e := range edges {
		if uf.Union(e.A, e.B) {
			if uf.UnionsCount == 1 {
				x1 = positions[e.A].X
				x2 = positions[e.B].X
				break
			}
		}
	}

	return x1 * x2, nil
}

type Pos struct {
	X, Y, Z int
}

type Edge struct {
	A, B     int
	Distance float64
}

func (p Pos) Distance(other Pos) float64 {
	dx := float64(p.X - other.X)
	dy := float64(p.Y - other.Y)
	dz := float64(p.Z - other.Z)
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

func parsePositions(data string) ([]Pos, error) {
	var positions []Pos
	for _, line := range strings.Split(data, "\n") {
		nums := strings.Split(line, ",")

		x, err := strconv.Atoi(nums[0])
		if err != nil {
			return nil, err
		}
		y, err := strconv.Atoi(nums[1])
		if err != nil {
			return nil, err
		}
		z, err := strconv.Atoi(nums[2])
		if err != nil {
			return nil, err
		}

		positions = append(positions, Pos{X: x, Y: y, Z: z})
	}

	return positions, nil
}

func buildEdges(positions []Pos) []Edge {
	n := len(positions)
	edges := make([]Edge, 0, n*n)

	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			d := positions[i].Distance(positions[j])
			edges = append(edges, Edge{A: i, B: j, Distance: d})
		}
	}

	return edges
}

type UnionFind struct {
	Parent      []int
	Size        []int
	UnionsCount int
}

func NewUnionFind(numberOfPositions int) *UnionFind {
	parent := make([]int, numberOfPositions)
	size := make([]int, numberOfPositions)
	for i := 0; i < numberOfPositions; i++ {
		parent[i] = i
		size[i] = 1
	}
	return &UnionFind{
		Parent:      parent,
		Size:        size,
		UnionsCount: numberOfPositions,
	}
}

func (u *UnionFind) Find(idx int) int {
	for idx != u.Parent[idx] {
		u.Parent[idx] = u.Parent[u.Parent[idx]]
		idx = u.Parent[idx]
	}
	return idx
}

func (u *UnionFind) Union(a, b int) bool {
	ra := u.Find(a)
	rb := u.Find(b)
	if ra == rb {
		return false
	}

	if u.Size[ra] < u.Size[rb] {
		ra, rb = rb, ra
	}

	u.Parent[rb] = ra
	u.Size[ra] += u.Size[rb]
	u.UnionsCount--
	return true
}

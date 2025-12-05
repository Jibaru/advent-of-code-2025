package day5

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

func Solve(part int, isTest bool) (any, error) {
	f := "day_5/input.txt"
	if isTest {
		f = "day_5/input-test.txt"
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
	ranges, ids, err := parseIngredientDB(data)
	if err != nil {
		return nil, err
	}

	ans := 0
	for _, id := range ids {
		for _, r := range ranges {
			if r.Inside(id) {
				ans++
				break
			}
		}
	}

	return ans, nil
}

func partTwo(data string) (any, error) {
	ranges, _, err := parseIngredientDB(data)
	if err != nil {
		return nil, err
	}

	// Sort by start
	sort.Slice(ranges, func(i, j int) bool {
		return ranges[i].From < ranges[j].From
	})

	ans := 0
	current := ranges[0]

	for _, r := range ranges[1:] {
		if current.Overlaps(r) {
			current = current.Merge(r)
			continue
		}

		ans += current.Size()
		current = r
	}

	// Add the last range
	ans += current.Size()

	return ans, nil
}

type Range struct {
	From int
	To   int
}

type IDList []int

func (r Range) Inside(id int) bool {
	return id >= r.From && id <= r.To
}

func (r Range) Overlaps(other Range) bool {
	return r.From <= other.To && other.From <= r.To
}

func (r Range) Merge(other Range) Range {
	return Range{
		From: min(r.From, other.From),
		To:   max(r.To, other.To),
	}
}

func (r Range) Size() int {
	return r.To - r.From + 1
}

func parseIngredientDB(data string) ([]Range, IDList, error) {
	sections := strings.SplitN(strings.TrimSpace(data), "\n\n", 2)
	if len(sections) != 2 {
		return nil, nil, fmt.Errorf("expected two sections: ranges and ids")
	}

	ranges, err := parseRangesSection(sections[0])
	if err != nil {
		return nil, nil, err
	}

	ids, err := parseIDsSection(sections[1])
	if err != nil {
		return nil, nil, err
	}

	return ranges, ids, nil
}

func parseRangesSection(section string) ([]Range, error) {
	var ranges []Range

	for _, line := range strings.Split(section, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, "-", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid range: %q", line)
		}

		from, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("invalid range start %q: %w", parts[0], err)
		}

		to, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("invalid range end %q: %w", parts[1], err)
		}

		ranges = append(ranges, Range{From: from, To: to})
	}

	return ranges, nil
}

func parseIDsSection(section string) (IDList, error) {
	var ids IDList

	for _, line := range strings.Split(section, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		id, err := strconv.Atoi(line)
		if err != nil {
			return nil, fmt.Errorf("invalid ingredient id %q: %w", line, err)
		}

		ids = append(ids, id)
	}

	return ids, nil
}

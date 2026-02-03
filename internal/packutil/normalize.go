package packutil

import (
	"errors"
	"sort"
)

func NormalizePackSizes(packSizes []int, invalidErr error) ([]int, error) {
	if invalidErr == nil {
		invalidErr = errors.New("pack sizes must be positive")
	}
	if len(packSizes) == 0 {
		return nil, invalidErr
	}

	seen := make(map[int]struct{}, len(packSizes))
	for _, p := range packSizes {
		if p <= 0 {
			return nil, invalidErr
		}
		seen[p] = struct{}{}
	}

	out := make([]int, 0, len(seen))
	for p := range seen {
		out = append(out, p)
	}

	sort.Ints(out)
	return out, nil
}

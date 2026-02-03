package pack

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSolve_SpecExamples(t *testing.T) {
	packSizes := []int{250, 500, 1000, 2000, 5000}

	cases := []struct {
		name        string
		amount      int
		items       int
		expected    map[int]int
		packCount   int
		overage     int
	}{
		{
			name:      "order 1",
			amount:    1,
			items:     250,
			expected:  map[int]int{250: 1},
			packCount: 1,
			overage:   249,
		},
		{
			name:      "order 250",
			amount:    250,
			items:     250,
			expected:  map[int]int{250: 1},
			packCount: 1,
			overage:   0,
		},
		{
			name:      "order 251",
			amount:    251,
			items:     500,
			expected:  map[int]int{500: 1},
			packCount: 1,
			overage:   249,
		},
		{
			name:      "order 501",
			amount:    501,
			items:     750,
			expected:  map[int]int{500: 1, 250: 1},
			packCount: 2,
			overage:   249,
		},
		{
			name:      "order 12001",
			amount:    12001,
			items:     12250,
			expected:  map[int]int{5000: 2, 2000: 1, 250: 1},
			packCount: 4,
			overage:   249,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := Solve(tc.amount, packSizes)
			require.NoError(t, err)
			assert.Equal(t, tc.amount, res.Amount)
			assert.Equal(t, tc.items, res.ItemsShipped)
			assert.Equal(t, tc.overage, res.Overage)
			assert.Equal(t, tc.packCount, res.PackCount)
			assert.Equal(t, tc.expected, res.Packs)
		})
	}
}

func TestSolve_EdgeCaseLargeAmount(t *testing.T) {
	packSizes := []int{23, 31, 53}
	amount := 500000

	res, err := Solve(amount, packSizes)
	require.NoError(t, err)
	assert.Equal(t, amount, res.ItemsShipped)
	assert.Equal(t, 0, res.Overage)
	assert.Equal(t, map[int]int{23: 2, 31: 7, 53: 9429}, res.Packs)
}

func TestSolve_Validation(t *testing.T) {
	_, err := Solve(-1, []int{250})
	require.ErrorIs(t, err, ErrInvalidAmount)

	_, err = Solve(1, []int{})
	require.ErrorIs(t, err, ErrInvalidPackSizes)

	_, err = Solve(1, []int{0, 250})
	require.ErrorIs(t, err, ErrInvalidPackSizes)
}

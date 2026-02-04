package pack

import (
	"errors"
	"fmt"
	"math"
	"pack-calculator/internal/packutil"
)

var (
	ErrInvalidAmount    = errors.New("amount must be >= 0")
	ErrInvalidPackSizes = errors.New("pack sizes must be positive")
	ErrNoSolution       = errors.New("no solution found")
	ErrReconstruct      = errors.New("failed to reconstruct solution")
)

func Solve(amount int, packSizes []int) (Solution, error) {
	if amount < 0 {
		return Solution{}, ErrInvalidAmount
	}

	sizes, err := packutil.NormalizePackSizes(packSizes, ErrInvalidPackSizes)
	if err != nil {
		return Solution{}, fmt.Errorf("can't normalize packsize %w", err)
	}

	return solveNormalized(amount, sizes)
}

func solveNormalized(amount int, sizes []int) (Solution, error) {
	if amount == 0 {
		return Solution{
			Packs: map[int]int{},
		}, nil
	}

	maxPack := sizes[len(sizes)-1]
	limit := amount + maxPack - 1

	const inf = math.MaxInt32
	dp := make([]int, limit+1)
	prev := make([]int, limit+1)

	for i := 1; i <= limit; i++ {
		dp[i] = inf
		prev[i] = -1
	}

	for t := 1; t <= limit; t++ {
		for _, p := range sizes {
			if t >= p && dp[t-p]+1 < dp[t] {
				dp[t] = dp[t-p] + 1
				prev[t] = p
			}
		}
	}

	bestTotal := -1
	for t := amount; t <= limit; t++ {
		if dp[t] != inf {
			bestTotal = t
			break
		}
	}
	if bestTotal == -1 {
		return Solution{}, ErrNoSolution
	}

	packs := make(map[int]int)
	for t := bestTotal; t > 0; {
		p := prev[t]
		if p <= 0 {
			return Solution{}, ErrReconstruct
		}
		packs[p]++
		t -= p
	}

	return Solution{
		Packs: packs,
	}, nil
}

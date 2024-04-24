package main

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"sync"
)

const SAMPLES = 100_000

// Takes a C-value as the initial probability, and simulates a "dice roll".
// If the roll succeeds, the function returns the number of dice rolls that happened.
// If the roll doesn't succeed, the probability increases by C, and another dice roll occurs.
// Eventually, the probability will exceed 1.0, and the dice roll is forced to succeed.
func trialCValue(C float64) int64 {
	if C <= 0 {
		panic("C must be greater than zero")
	} else if C > 1 {
		return 1
	}

	N := int64(1)

	for {
		if rand.Float64() <= C*float64(N) {
			return N
		}
		N += 1
	}
}

// Approximates the expected value (EV) for a given C-value.
func calcEV(C float64) float64 {
	sum := int64(0)

	for i := 0; i < SAMPLES; i++ {
		sum += trialCValue(C)
	}

	return float64(sum) / float64(SAMPLES)
}

func main() {
	type Row struct {
		C  float64
		P  float64
		EV float64
	}

	data := []Row{}
	mut := sync.Mutex{}
	wg := sync.WaitGroup{}

	precision := 3
	iNextStepIncrease := int(math.Pow(10, float64(precision)))
	i := 1
	iStep := 1
	cStep := 1e-10
	C := float64(0)

	// Computes the expected value for a C-value and then adds the result to an array.
	// Each iteration the C-value increases by a small step amount. The step amount will
	// also increase over time to maintain a specific precision.
	for C < 1 {
		C = cStep * float64(i)

		if i == iNextStepIncrease {
			iStep *= 10
			iNextStepIncrease *= 10
		}

		wg.Add(1)
		go func(C float64, i int) {
			defer wg.Done()
			ev := calcEV(C)
			mut.Lock()
			data = append(data, Row{C, 1 / ev, ev})
			mut.Unlock()
		}(C, i)

		i += iStep
	}

	wg.Wait()
	sort.Slice(data, func(i, j int) bool { return data[i].C < data[j].C })

	for _, row := range data {
		fmt.Printf("%.12f,%.10f,%f\n", row.P, row.C, row.EV)
	}
}

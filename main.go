package main

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"sync"
)

const (
	// The number of samples for each C-value. Smaller = fast, larger = better approximation
	EV_SAMPLES = 200

	// The precision of C-values. Smaller = fast, larger = more data
	PRECISION = 3

	// The C-value that the simulation will start at, as well as the amount the C-value
	// is incremented by each iteration. The number of iterations needed starts to increase
	// dramatically for C-values smaller than 1e-07
	INITIAL_C_VALUE = 1e-08
)

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

	for i := 0; i < EV_SAMPLES; i++ {
		sum += trialCValue(C)
	}

	return float64(sum) / float64(EV_SAMPLES)
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

	iNextStepIncrease := int(math.Pow(10, float64(PRECISION)))
	i := 1
	iStep := 1
	C := float64(0)

	// Computes the expected value for a range of C-values. The precision of the
	// C-value is controlled in order to obtain a wide range of datapoints in a
	// reasonable amount of time
	for C < 1 {
		C = INITIAL_C_VALUE * float64(i)

		// The precision of the next C-value will be too high
		if i == iNextStepIncrease {
			// Increasing the step by a factor of 10 shifts everything over one digit
			iStep *= 10
			// After 10 iterations a new digit will be introduced (1, 2, 3, ..., 9, 10)
			// and the precision will need to be adjusted again
			iNextStepIncrease *= 10
		}

		wg.Add(1)

		// Run the EV calc in parallel
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

	// Sort results by C-value. They may be out of order due to running in parallel
	sort.Slice(data, func(i, j int) bool { return data[i].C < data[j].C })

	for _, row := range data {
		fmt.Printf("%.12f,%.12f,%f\n", row.C, row.P, row.EV)
	}
}

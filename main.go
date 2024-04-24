package main

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"sync"
)

const P_SAMPLES = 100_000

func rollN(C float64) int64 {
	if C <= 0 {
		panic("C must be greater than zero")
	} else if C > 1 {
		return 1
	}

	N := int64(0)

	for {
		if rand.Float64() <= C*float64(N) {
			return N
		}
		N += 1
	}
}

func calcEV(C float64) float64 {
	sum := int64(0)

	for i := 0; i < P_SAMPLES; i++ {
		sum += rollN(C)
	}

	return float64(sum) / float64(P_SAMPLES)
}

// func calcC(P float64) float64 {
// 	a := -2.266e-00
// 	b := 3.959e+00
// 	c := -1.208e+00
// 	d := 5.511e-01
// 	e := -2.854e-02

// 	return a*math.Pow(P, 4) + b*math.Pow(P, 3) + c*math.Pow(P, 2) + d*P + e
// }

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

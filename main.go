package main

import (
	"fmt"
	"math/rand"
	"sort"
	"sync"
)

const P_SAMPLES = 100

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
	wg := sync.WaitGroup{}
	initialIncr := float64(0.000001)
	incr := initialIncr
	C := float64(0)
	i := 0
	mut := sync.Mutex{}

	for C < 0.01 {
		if i == 100 || i == 1000 || i == 10000 {
			incr = initialIncr * float64(i/10)
		}
		C += incr

		wg.Add(1)
		go func(C float64, i int) {
			defer wg.Done()
			ev := calcEV(C)
			mut.Lock()
			data = append(data, Row{C, 1 / ev, ev})
			mut.Unlock()
		}(C, i)

		i++
	}

	wg.Wait()
	sort.Slice(data, func(i, j int) bool { return data[i].C < data[j].C })

	for _, row := range data {
		fmt.Printf("%f,%f,%f\n", row.C, row.P, row.EV)
	}
}

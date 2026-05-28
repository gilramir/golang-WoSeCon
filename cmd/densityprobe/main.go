// densityprobe sweeps a fixed word list across a range of grid sizes
// (so the "letter density" — total letters / grid cells — varies from
// generous to extreme) and reports how long Construct takes for each
// (grid, seed). It's the tool that produced the table in the README's
// "Density and search time" section. Re-run it after algorithm changes
// to refresh those numbers.
//
// All knobs are hard-coded so the run is reproducible. To experiment,
// edit the word list, grids, seeds, or per-run cap below.
package main

import (
	"fmt"
	"time"

	wosecon "github.com/gilramir/golang-WoSeCon/v2"
)

func main() {
	// 12 short English words, no palindromes, no reverse-pairs. The
	// total letter count (56) is what makes the density numbers tidy.
	words := []string{
		"PUZZLE", "WORDS", "GAMES", "FUN", "PLAY", "SEARCH",
		"GRID", "CELL", "ROW", "COLUMN", "DOWN", "ACROSS",
	}
	totalLetters := 0
	for _, w := range words {
		totalLetters += len(w)
	}

	// Grids span letter density from ~47% (lots of slack) to ~156%
	// (overlap mandatory and very tight).
	grids := []struct{ cols, rows int }{
		{12, 10},
		{10, 8},
		{9, 7},
		{8, 7},
		{7, 7},
		{7, 6},
		{6, 6},
	}
	seeds := []int64{0, 1, 2}
	perRunCap := 15 * time.Second

	fmt.Printf("Word list: %d words, %d total letters (no palindromes, no reverse-pairs)\n",
		len(words), totalLetters)
	fmt.Printf("Directions: NaturalLTR (Down, LTRHorizontal, LTRDescending, LTRAscending)\n")
	fmt.Printf("Per-run cap: %s\n\n", perRunCap)

	fmt.Printf("%-6s %5s %8s   ", "grid", "cells", "density")
	for _, s := range seeds {
		fmt.Printf("%-9s ", fmt.Sprintf("seed=%d", s))
	}
	fmt.Println()
	fmt.Println("------------------------------------------------------------")

	for _, g := range grids {
		cells := g.cols * g.rows
		density := 100.0 * float64(totalLetters) / float64(cells)
		fmt.Printf("%-6s %5d %7.0f%%   ",
			fmt.Sprintf("%dx%d", g.cols, g.rows), cells, density)

		for _, seed := range seeds {
			start := time.Now()
			_, err := wosecon.Construct(g.cols, g.rows, words,
				wosecon.FillUniformlyFromString("ABCDEFGHIJKLMNOPQRSTUVWXYZ"),
				wosecon.AddNaturalLTRDirections(),
				wosecon.RandomSeed(seed),
				wosecon.WithTimeLimit(perRunCap),
			)
			elapsed := time.Since(start)
			label := elapsed.Round(time.Millisecond).String()
			switch err {
			case nil:
				// solved; show elapsed
			case wosecon.ErrReachedTimeLimit:
				label = fmt.Sprintf(">%s", perRunCap)
			case wosecon.ErrCannotFitWords:
				label = "no-fit"
			default:
				label = err.Error()
			}
			fmt.Printf("%-9s ", label)
		}
		fmt.Println()
	}
}

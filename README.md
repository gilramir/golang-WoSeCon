# golang-WoSeCon: A Word Search Constructor library

This is based off of the WeSoCon algorithm desribed in
[this article](http://ijses.com/wp-content/uploads/2022/01/68-IJSES-V6N1.pdf),
written by Lefteris Moussiades.

In this implementation, however, words can be placed in all 8 directions,
at the choice of the user.

v1 of this API used runes (Unicode code poitns) as the type
for cells in the puzzle grid.

v2 now uses strings, to support puzzles that need more than one rune
in a cell.

The API supports filling in the non-solution part of the
puzzle with random values, from a list of possible values that you provide.

See the
[GoDoc documentation](https://pkg.go.dev/github.com/gilramir/golang-WoSeCon/v2)
for this module.

The algorithm performs an exhaustive search. As it attempts to place words randomly,
if it cannot place a word on the first try, it will backtrack, removing
words and replacing them, until the entire solution space has been exhausted.
Under normal conditions, this does not take a long time.

Other WoSeCon implementations:
* **C++** - The reference implementation, at https://github.com/lmous/WoSeCon
* **.NET** - https://github.com/martinrotter/wordsearch-wosecon

## Code

An example of using this API is in cmd/mkwordsearch in this repo.
It is a very simple CLI tool which is useful during development testing.

## API

Construct a WordSearch with one simple function call:

* **Construct(numCols, numRows, words, ...options)** - this adds all Right-to-left directions,
    and Up, to the list of possible directions for the solution.

```
import (
    wosecon "github.com/gilramir/golang-WoSeCon/v2"
)

numColumns := 8
numRows := 5
words := []string{"TOYS", "APPLE", "CAR"}
alphabet := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

wordSearch, err := wosecon.Construct(numColumns, numRows, words,
    wosecond.FillUniformlyFromString(alphabet))
```

The **WordSearch** object gives you the placement of each word. A placement
is the starting colum and row, and direction of the word. It also gives
you a matrix of strings for the solution (the puzzle with no filler), and
the complete puzzle (solution + filler)

If you print the **Solution**, it would look something like this. However,
empty cells in the grid are given as empty string, so you must take that
into account when printing.
```
    0  1  2  3  4  5  6  7

 0     A

 1        P  S     C

 2        Y  P     A

 3     O        L  R

 4  T              E
```

If you print the **Puzzle**, which has the **Solution** and the empty cells are filled
randomly, you might see this:
```
    0  1  2  3  4  5  6  7

 0  R  A  A  L  H  B  T  T

 1  G  C  P  S  X  C  J  Y

 2  K  X  Y  P  C  A  C  D

 3  K  O  N  N  L  R  B  Z

 4  T  D  L  Q  L  E  N  E
```

You can use the cmd/mkwordsearch/mkwordsearch example tool to experiment:
```
./mkwordsearch 8 5 small.txt
```

## Options

When running Constructor(), you can pass it options to control its behavior.
The options are:

* **AddNaturalLTRDirections()** - this adds all Left-to-right directions,
    and Down, to the list of possible directions for the solution.
    This is the default, and doesn't need to be used. This adds, so it
    can be combined with other options which also add directions.

* **AddUnnaturalLTRDirections()** - this adds all Right-to-left directions,
    and Up, to the list of possible directions for the solution.

* **AddDirection(direction Direction)** - this adds a single direction
    to the list of possible directions for the solution. The 8 possible
    directions are:
    * Down
    * LTRHorizontal
    * LTRAscending
    * LTRDescending
    * Up
    * RTLHorizontal
    * RTLAscending
    * RTLDescending
    As the directions are bit values, you can bitwise-OR more than
    one direction, as in:

```
    wosecon.AddDirection(wosecon.Down|wosecon.Up)
```

* **UseAllDirections()** - the solution can use all 8 of the possible
    directions.

* **RandomSeed(seed int64)** - seed the random number generator. Use this
    if you want repeatable randomness. Otherwise, the seed itself is
    already random.

* **FillUniformlyFromString(filler string)** - fill in the puzzle
    with equally-random choices of runes from this string.

* **FillUniformlyFromStringSlice(filler []string)** - fill in the puzzle
    with equally-random choices of strings from this slice.

* **FillUniformlyFromRuneSlice(filler []rune)** - fill in the puzzle
    with equally-random choices of runes from this slice.

* **FillWeighted(filler []FillerWeight)** - fill in the puzzle
    based on weighted values. Give a slice of RuneWeight objects,
    which are just a rune and its relative weight. The weights do not
    need to sum up to some specific value. The code will sum the
    weights and use their percentages to decide the probability that
    a rune will be chosen to fill in an empty cell in the puzzle.
    Example:

```
    fillerWeights := make([]wosecon.FillerWeight, 3)
    fillerWeights[0] = wosecon.FillerWeight{"E", 20}
    fillerWeights[1] = wosecon.FillerWeight{"S", 18}
    fillerWeights[2] = wosecon.FillerWeight{"T", 15}

	wordSearch, err := wosecon.Construct(10, 10, wordSlice,
		wosecon.FillWeighted(fillerWeights))
```

* **WithTimeLimit(time.Duration)** - if your input may contain very large
    lists of words that cannot create a puzzle due to size, this sets
    a maximum time that the algorithm will run. If the time limit is reached,
    ErrReachedTimeLimit is returned.

* **WithMemoizationLimit(maxEntries int)** - the Constructor can cache
    the cell states it has already proven to be dead ends, so the same
    dead end is not re-derived through a different backtrack path. The
    cache is **off by default** (`maxEntries == 0`) because it only
    helps in a narrow set of inputs; for most word lists, the per-call
    fingerprinting cost outweighs the savings, so enabling it makes the
    construction *slower*, not faster.

    Turn the cache on when your word list can produce the same grid of
    letters via more than one combination of word placements. The two
    common cases are:

    1. **A palindromic word in a bidirectional grid.** "LEVEL" placed
        `LTRHorizontal` at `(col=0, row=r)` lands the letters
        `L, E, V, E, L` in cells `(0,r)..(4,r)`. The same word placed
        `RTLHorizontal` at `(col=4, row=r)` lands the same letters in
        the same cells. The algorithm sees two distinct placements
        that produce identical grids; the cache lets the search skip
        the second subtree once the first is explored.

        English examples: `RACECAR`, `RADAR`, `LEVEL`, `NOON`,
        `KAYAK`, `MADAM`, `REFER`, `CIVIC`. Only relevant when you
        enable a direction *and* its reverse (e.g. both
        `LTRHorizontal` and `RTLHorizontal`, or both `Down` and
        `Up`). A puzzle that uses only the natural LTR directions
        won't see this.

    2. **A pair of words that are reverses of each other in a
        bidirectional grid.** If the word list contains both `STAR`
        and `RATS` and the grid allows both LTR and RTL movement,
        placing `STAR` LTR at `(0,r)` and placing `RATS` RTL at
        `(3,r)` produce the same letters in the same cells via two
        different placement choices.

        Pairs to watch for: `STAR`/`RATS`, `STOP`/`POTS`,
        `LIVE`/`EVIL`, `GOD`/`DOG`, `DESSERTS`/`STRESSED`.

    For *typical* puzzles — no palindromes, no word that is another
    word spelled backwards, or a uni-directional grid — leave this at
    `0`. Concretely, the dense Korean puzzle in `cmd/mkwordsearch`
    runs roughly 2× faster with the cache off than with it on.

    When you do enable it, each cached entry holds a fingerprint of
    the puzzle grid (one cell string per cell, NUL-separated) plus
    Go map overhead. As a rough rule, expect about
    `(numCols * numRows * bytesPerCellString) + 150` bytes per entry,
    where `bytesPerCellString` is 1 for ASCII single-rune cells and
    3 for typical Korean/CJK cells. Some sample sizings:

    | Grid    | ASCII cells (~2 B) | Korean cells (~4 B) |
    |---------|--------------------|----------------------|
    |  7 ×  5 | ~220 B/entry       | ~290 B/entry         |
    | 15 × 15 | ~600 B/entry       | ~1.05 KB/entry       |
    | 25 × 25 | ~1.4 KB/entry      | ~2.65 KB/entry       |

    Multiply by `maxEntries` to bound memory. A reasonable starting
    point for an English puzzle with several palindromes is
    `WithMemoizationLimit(10000)` (≈2 MB). Once the cap is reached,
    further dead ends stop being recorded but existing entries still
    short-circuit lookups.

## Errors

The **Constructor** functions will return only a handful of specific errors, provided
a constants starting with the name "Err". See the documentation for the values.

## Cells with more than one Rune

If you need to have more than one run per cell in the grid, you need to make a
Sequence for each "word", and use. I needed this when making a puzzle game for
the Thai language, in which one "character" can be a combination of multiple
Unicode code points.

* **ConstructFromSequences(numCols, numRows, sequences, ...options)**

Where a **Sequence** is
```
type Sequence interface {
	String() string
	Len() int
	Cmp(Sequence) int
	Index(int) string
	Items() iter.Seq2[int, string]
}
```
The **NewRuneSequence** and **NewStringSequence** functions provide Sequences for normal
words (a string of runes), and multi-rune-per-cell sequences
(NewStringSequence).  You can implement your own Sequence, of course.

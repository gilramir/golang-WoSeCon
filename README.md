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

## The WordSearch result

On success, `Construct` and `ConstructFromSequences` return a pointer to
a `WordSearch` struct that contains everything needed to render the
puzzle or check student answers:

```go
type WordSearch struct {
    NumCols int
    NumRows int

    Sequences      map[string]Sequence
    WordPlacements map[string]WordPlacement

    AllPossibleWordPlacements     map[string][]WordPlacement
    NumWordsWithMultipleSolutions int

    SolutionRows [][]string
    PuzzleRows   [][]string

    TotalLetters  int
    LetterDensity float64
}
```

Fields:

* **NumCols, NumRows** - grid size, the same values that were passed to
    `Construct`.

* **Sequences** - the input words keyed by their string form, paired
    with the `Sequence` object that the constructor used to split each
    word into cell-sized pieces. For the common case of one rune per
    cell, this is a `RuneSequence`; for multi-rune cells (see the
    "Cells with more than one Rune" section), this is whichever
    `Sequence` implementation you supplied. Useful when rendering, to
    iterate the cells in placement order without re-decoding the
    string.

* **WordPlacements** - the placement that the algorithm picked for each
    word. The key is the word string and the value is a `WordPlacement`
    containing the starting column, starting row, and direction:

    ```go
    type WordPlacement struct {
        Col       int
        Row       int
        Direction Direction
    }
    ```

    `WordPlacement.DirectionString()` returns a human-readable name
    ("Down", "LTRHorizontal", etc.) for the direction.

* **AllPossibleWordPlacements** - the same map shape as
    `WordPlacements`, but includes *all* placements of the word that
    can be found in the final puzzle, including ones that the random
    filler happened to create. Each word's value is a slice because
    one word may have several valid placements once filler is added.
    If a short word is a substring of a longer word, the match is
    credited to the longer word only.

* **NumWordsWithMultipleSolutions** - the number of words whose entry
    in `AllPossibleWordPlacements` has more than one element. A puzzle
    with no random filler typically has `0` here; a puzzle whose
    filler accidentally spells extra copies of a word will have
    `> 0`. Useful as a quality signal when generating puzzles for
    students.

* **SolutionRows** - a `[][]string` view of the puzzle showing only
    the words the algorithm placed; cells that would be filled by the
    random filler are empty strings (`""`). The outer slice is rows
    (y), the inner slice is columns (x), so to draw the grid you
    iterate `SolutionRows[row][col]`.

* **PuzzleRows** - the same shape as `SolutionRows`, but every cell is
    filled — the words plus whatever filler the `FillUniformly*` or
    `FillWeighted` option produced. Empty strings only appear if you
    didn't pass a filler option.

* **TotalLetters** - the sum of cell counts across the input word
    list (`sum of seq.Len()`). Distinct from the number of *occupied*
    grid cells, which is smaller when words overlap on shared
    characters.

* **LetterDensity** - `TotalLetters / (NumCols * NumRows)`. A value
    of `1.0` means letters and cells balance exactly; values above
    `1.0` mean the algorithm found an arrangement where words shared
    cells. The "Density and search time" section below explains how
    this number relates to construction time — `1.0` to `~1.3` is
    the regime where you might want `WithParallelism(n)`, and above
    `~1.3` is where you should expect `ErrCannotFitWords` for most
    word lists.

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

    Note also that when combined with `WithParallelism(n)`, each
    worker carries its own cache, so the worst-case memory cost is
    `n * maxEntries` entries.

* **WithParallelism(n int)** - race `n` independent searches in
    parallel, each starting from a distinct RNG seed. The first
    goroutine to find a solution wins; the others are cancelled. See
    the "Density and search time" section below for the regime where
    this helps and a measured table of the effect.

    Semantics:

    | `n` value | Behavior                                                |
    |-----------|---------------------------------------------------------|
    | `0`       | Use `runtime.NumCPU()` workers.                         |
    | `1`       | Single-threaded (equivalent to leaving the option off). |
    | `> 1`     | Race that many workers.                                 |
    | `< 0`     | Treated as `1`.                                         |

    Caveats:

    - **Determinism is lost.** Which seed wins depends on goroutine
      scheduling, so the returned solution will vary run to run even
      with the same `Construct` arguments.
    - **`RandomSeed` cannot be combined with `WithParallelism(n > 1)`.**
      Each worker needs its own seed, and silently overriding a
      caller-supplied seed would be surprising. Combining them
      returns `ErrSeedWithParallelism` immediately, before any
      search runs.
    - **Memory scales with `n`.** Each worker has its own solution
      matrix, locator, and (if enabled) dead-end cache. For typical
      grids this is well under 1 MB per worker, but it composes
      with `WithMemoizationLimit` as noted above.
    - **No help below ~100% letter density or well above ~130%.**
      Those regimes already finish in milliseconds. Parallelism
      pays off in the borderline regime in between, where one seed
      may wander for seconds while another stumbles onto a solution
      immediately.

## Density and search time

The WoSeCon algorithm is exhaustive backtracking — it walks the tree of
possible word placements with no global view of which subtrees can ever
lead to a valid puzzle. For grids that have plenty of room this is
fine, but as the puzzle gets denser the time to find (or rule out) a
solution grows sharply.

The useful single number is the **letter density**: total letters in
the word list divided by grid cells. Densities above 100% mean a valid
puzzle is only possible if words share cells.

In Go, computing the score for a list of words is straightforward.
Count *cells*, not bytes — use `len([]rune(w))` so multi-byte characters
(Korean, accented Latin, etc.) count as one cell each:

```go
totalCells := 0
for _, w := range words {
    totalCells += len([]rune(w))
}
density := 100.0 * float64(totalCells) / float64(numCols*numRows)
fmt.Printf("letter density: %.0f%%\n", density)
```

If you are using `ConstructFromSequences` with multi-rune cells (the
Thai case described below), substitute `seq.Len()` for
`len([]rune(w))` so each `Sequence` contributes its cell count rather
than its rune count.

Three regimes show up in practice:

1. **Below ~100% density.** Solutions are abundant; the first valid
   placement usually leads straight to a complete puzzle, and
   `Construct` finishes in milliseconds.

2. **Slightly above 100% (roughly 110%–130%).** Solutions *may*
   exist but require overlap. This is the slow regime. Many partial
   placements look fine but turn out to be dead ends only after the
   algorithm has explored a long subtree. Different random seeds at
   the same grid size can swing the runtime from milliseconds to
   "longer than you want to wait."

3. **Well above ~130%.** The constraints are usually tight enough
   that early placements force a contradiction quickly, the
   algorithm exhausts the space, and `Construct` returns
   `ErrCannotFitWords` in negligible time. (Whether a solution
   *actually* exists at these densities depends on how much the
   words share letters; for many word lists, they do not share
   enough.)

The `cmd/densityprobe` program in this repository measures all three
regimes against a fixed list of 12 short English words (56 letters
total, no palindromes, no reverse-pairs), sweeping grid size to vary
density. Three random seeds per grid, 15-second per-run cap. The
numbers below were measured on an AMD Ryzen 7 8845HS (16 logical CPUs,
boost ~5.1 GHz) running Go 1.25.1 on Linux — single-threaded, so the
core count doesn't matter; the per-core speed does. On slower silicon
expect the milliseconds to scale up roughly with single-thread
performance, and on a thermally-limited laptop expect the `>15 s`
rows to stay at the cap regardless:

| Grid  | Cells | Letter density | seed=0 | seed=1 | seed=2 |
|-------|------:|---------------:|--------|--------|--------|
| 12×10 |   120 |  47%           | <1 ms  | <1 ms  | <1 ms  |
| 10× 8 |    80 |  70%           | <1 ms  | <1 ms  | <1 ms  |
|  9× 7 |    63 |  89%           | <1 ms  | <1 ms  | <1 ms  |
|  8× 7 |    56 | 100%           | 3 ms   | <1 ms  | <1 ms  |
|  7× 7 |    49 | 114%           | 569 ms | >15 s  | >15 s  |
|  7× 6 |    42 | 133%           | no-fit | no-fit | no-fit |
|  6× 6 |    36 | 156%           | no-fit | no-fit | no-fit |

Note the cliff at 114%: one seed solves in half a second, two others
fail to converge inside 15 seconds. That is regime 2 in action — the
borderline case where solutions are possible but the search has to
wander.

**If your call is returning `ErrReachedTimeLimit`**, the puzzle is
almost certainly in regime 2: a solution likely exists, but the
algorithm cannot find it in the budget you allowed. The fix that
actually works is to enlarge the grid by one or two rows or columns,
which drops the density into regime 1 and usually solves in
milliseconds. The longer-form workarounds (different random seeds,
longer time limits, enabling more directions) help occasionally; a
slightly bigger grid almost always does.

### Parallel search

If you cannot enlarge the grid, the second-best lever is
`WithParallelism(n)`. Regime 2 is high-variance in the random seed:
one seed may take ~10 seconds while another solves the same puzzle
in milliseconds. Racing several seeds in parallel takes the
*minimum* of those times rather than the average.

Five trials per parallelism level on the 7×7 / 114% case from the
table above (each trial uses a fresh per-process seed), 15-second
per-run cap, same Ryzen 7 8845HS as above:

| Parallelism | Min   | Median | Max    |
|-------------|-------|--------|--------|
| 1           | 83 ms | 12 s   | >15 s  |
| 2           |  4 ms | 720 ms | 2.9 s  |
| 4           |  8 ms | 27 ms  | 86 ms  |
| 8           | 10 ms | 40 ms  | 2.6 s  |

Going from `WithParallelism(1)` to `WithParallelism(4)` cut the
median runtime by roughly 400×. Going from 4 to 8 didn't help on
this case (and was slightly worse on the worst trial) — once the
variance is mostly squeezed out, additional workers just compete
for CPU.

This won't make a puzzle that legitimately has no solution finish
any faster, won't help below ~100% density (the search already
finishes instantly), and won't rescue the well-above-130% no-fit
cases (the algorithm already proves them quickly). It is precisely
the borderline regime that benefits.

To re-measure on your own machine — or after changes to the
algorithm — build and run the probe:

```
go build -o densityprobe ./cmd/densityprobe
./densityprobe
```

The word list, grid range, seeds, and per-run cap are hard-coded at
the top of `cmd/densityprobe/main.go`; edit them to explore other
shapes of the same question.

## Errors

The **Constructor** functions will return only a handful of specific errors,
defined by a constant starting with the name "Err". See the `api.go` for the values.

## Cells with more than one Rune

If you need to have more than one rune per cell in the grid, you need to make a
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

# golang-WoSeCon: A Word Search Constructor library

This is based off of the WeSoCon algorithm desribed in
[this article](http://ijses.com/wp-content/uploads/2022/01/68-IJSES-V6N1.pdf),
written by Lefteris Moussiades.

In this implementation, however, words can be placed in all 8 directions,
at the choice of the user.

The API has been developed to support runes, as opposed to bytes, to any
unicode codepoint can be used in any cell of the grid,
even if the codepoint encodes into multiple bytes as UTF-8.

The API supports filling in the puzzle with random runes.

See the
[GoDoc documentation](https://pkg.go.dev/github.com/gilramir/golang-WoSeCon)
for this module.

The algorithm is exhaustive. As it attempts to place words randomly,
if it cannot place a word on the first try, it will backtrack, removing
words and replacing them, until the entire solution space has been exhausted.
Under normal conditions, this does not take a long time.

Other WoSeCon implementations:
* **.NET** - https://github.com/martinrotter/wordsearch-wosecon

## Code

An example of using this API is in cmd/mkwordsearch in this repo.
It is a very simple CLI tool which is useful during development testing.

## API

Construct a WordSearch with one simple function call:

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

The WordSearch object gives you the placement of each word. A placement
is the starting colum and row, and direction of the word. It also gives
you a matrix of runes for the solution (the puzzle with no filler), and
the complete puzzle (solution + filler)

If you print the Solution, it would look something like this. However,
empty cells in the grid are given as empty string, so you must take those
into account.
```
    0  1  2  3  4  5  6  7

 0     A

 1        P  S     C

 2        Y  P     A

 3     O        L  R

 4  T              E
```

If you print the Puzzle, which has the Solution and the empty cells are filled
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

* **UseAllDirections()** - the solution can use all 8 of the possible
    directions.

* **RandomSeed(seed int64)** - seed the random number generator. Use this
    if you want repeatable randomness. Otherwise, the seed itself is
    already random.

* **FillUniformlyFromString(filler string)** - fill in the puzzle
    with equally-random choices of runes from this string.

* **FillUniformlyFromStringSlice(filler []string)** - fill in the puzzble
    with equally-random choices of strings from this slice.

* **FillUniformlyFromRuneSlice(filler []rune)** - fill in the puzzble
    with equally-random choices of runes from this slice.

* **FillWeighted(filler []FillerWeight)** - fill in the puzzble
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

## Errors

The Constructor function wil return only a handful of specific errors, provided
a constants starting with the name "Err". See the documentation.

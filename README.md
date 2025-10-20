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
    wosecon "github.com/gilramir/golang-WoSeCon"
)

numColumns := 8
numRows := 5
words := []string{"TOYS", "APPLE", "CAR"}
alphabet := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

wordSearch, err := wosecon.Construct(numColumns, numRows, words,
    wosecond.FillEquallyString(alphabet))
```

The WordSearch object gives you the placement of each word. A placement
is the starting colum and row, and direction of the word. It also gives
you a matrix of runes for the solution (the puzzle with no filler), and
the complete puzzle (solution + filler)

If you print the Solution, it would look something like this:
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

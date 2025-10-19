# Word Search Constructor

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

numColumns := 10
numRows := 10
words := []string{"APPLE", "BOY", "CAR"}
alphabet := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

wordSearch, err := wosecon.Construct(numColumns, numRows, words,
    wosecond.FillEquallyString(alphabet))
```

The WordSearch object gives you the placement of each word. A placement
is the starting colum and row, and direction of the word. It also gives
you a matrix of runes for the solution (the puzzle with no filler), and
the complete puzzle (solution + filler)

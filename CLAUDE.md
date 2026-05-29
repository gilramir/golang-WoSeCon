# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

A Go library that constructs word search puzzles using the WoSeCon algorithm by Lefteris Moussiades (reference C++ implementation at https://github.com/lmous/WoSeCon, paper PDF at `docs/68-IJSES-V6N1.pdf`). This implementation supports all 8 directions and multi-rune-per-cell sequences (e.g. for Thai). The module path is `github.com/gilramir/golang-WoSeCon/v2` — v2 changed cell type from `rune` to `string`.

## Commands

Build the example CLI:
```
go build ./cmd/mkwordsearch
```

Run the full test suite (the project uses go-check / gopkg.in/check.v1, bridged via `init_test.go`):
```
go test ./...
```

Run a single gocheck test method — use the `-check.f` flag to filter by method name (NOT the standard `-run`):
```
go test -check.f TestSolutions01
go test -check.f 'TestSolutions.*'
```

Manual puzzle generation for poking at behavior:
```
./mkwordsearch [-v] [-w weights.txt] NUM_COLS NUM_ROWS WORD_LIST_FILE
```
Sample input files live in `cmd/mkwordsearch/` (`small.txt`, `big.txt`, `ko.txt`, `long-word.txt`, `weights-abc.txt`).

Benchmark how letter density affects construction time (the tool that generated the README's "Density and search time" table):
```
go build -o densityprobe ./cmd/densityprobe
./densityprobe
```
The word list, grid sweep, seeds, and per-run cap are hard-coded in `cmd/densityprobe/main.go`; edit them in place when experimenting.

## Architecture

The public API is intentionally narrow — `Construct()` and `ConstructFromSequences()` in `api.go` are the only entry points. Everything else is unexported and orchestrated by `wordSearchConstructor` in `wosecon.go`.

**The placement algorithm (wosecon.go)** is an exhaustive backtracking search with two modes (`forwardMode` / `backwardMode`):
1. Words are sorted longest-first (stable, alphabetical tiebreak) — long words constrain placement the most.
2. For each word, `locateOne` walks a pre-shuffled list of candidate `directedLocation`s and calls `validPlacement` until one fits.
3. **Forward-check after every successful `validPlacement`:** `remainingHaveFit(currentWordIndex)` iterates the still-unplaced words and uses `canFit` (a pure read-only variant of `validPlacement`) to confirm each has at least one fitting location in the global locator. If any word has zero fits, the just-applied placement is undone, the location is added to `testedLocations` so a later backward-mode visit will skip it, and the search continues with the next candidate. Sound because future placements can only consume more cells.
4. When `locateOne` returns false, the algorithm backtracks: the previous word's location is moved to its `testedLocations` (so it won't be retried for *that* word at *that* level) and re-attempted from remaining candidates.
5. If word 0 runs out of candidates, the whole grid is impossible → `ErrCannotFitWords`.

**`tryPlace` wraps `locateOne` with optional dead-end memoization.** Default is off (`defaultDeadEndCacheCap = 0`, `deadEnds` map nil). When enabled via `WithMemoizationLimit(n>0)`, the key is `(solutionMatrix.fingerprint(), currentWordIndex)`; on a cache hit, return false without calling `locateOne`. On a `locateOne` failure, insert the key (bounded by the cap; once full, further failures stop being recorded). The cache is **only sound when no two backtrack paths can reach the same cell state at the same word index** — palindromes in a bidirectional grid and reverse-pairs (`STAR`/`RATS`) are the canonical inputs where it helps. For typical word lists the cache never hits and the per-call fingerprint cost slows the search down; this is why the default is off.

**Three cooperating data structures drive placement:**
- `randomLocator` (`randomlocator.go`) — flat, shuffled slice of every legal `(col, row, direction)` tuple. Built once via `directedLocationMatrix` (`directed_location.go`), which masks out directions that would run off the edge per-cell (e.g. top-row cells get no `Up`). When a word is placed, its location is removed from the locator; on backtrack it's re-added. The `minus(tested)` method returns a transient locator that skips a word's already-tested locations — and **re-shuffles** the result, so each backward-mode entry probes positions in a fresh order (matches the C++ reference's `RandomLocator::minus()` behavior and keeps the search from re-trying the same biased order).
- `solutionMatrix` (`solution_matrix.go`) — per-cell `{seqString, count}` with reference counting. Two overlapping words sharing the same rune at a cell are both legal and bump `count`; `clearPlacement` decrements. `isCellAvailableFor(col, row, text)` is the overlap-compatibility check. `fingerprint()` returns a row-major NUL-separated string used as the dead-end cache key.
- `wordInfo` (`wordinfo.go`) — per-word placement state plus `testedLocations` for backtracking memory.

**Sequences (`sequence.go`)** abstract "one cell's worth of content." `RuneSequence` is the default (one rune per cell); `StringSequence` allows multi-rune cells for scripts like Thai. The `Sequence` interface is public, so callers can implement their own. Cell content is always compared as `string`.

**Directions (`directions.go`)** are bit flags on a `uint8`. A single `Direction` value can mean one direction or a set; combine with `|`, test with `&`. `GoesDownward` / `GoesUpward` / `GoesLTR` / `GoesRTL` are aggregate masks used throughout to compute `colAdj` / `rowAdj` deltas.

**The `validPlacement` boundary check** is subtle: for `Down`/`LTR`, `endRow`/`endCol` is "one past the last cell" so out-of-grid is `> numRows` / `> numCols`. For `Up`/`RTL`, `endRow`/`endCol` is "one *before* the last cell" so the lowest legal value is `-1` (last cell at index 0) — the check is `< -1`, not `< 0`. Getting this wrong silently drops every `Up`/`LTRAscending`/`RTLHorizontal`/etc. placement that ends on row 0 or column 0. Locked in by `TestAscendingPlacementAtRowZero` in `directions_test.go`.

**Parallel construction (`construct_parallel.go`)** is dispatched from `ConstructFromSequences` whenever the option-applied `parallelism > 1`. It spawns N goroutines, each with its own `wordSearchConstructor` initialized from the original options plus a distinct `RandomSeed(baseSeed + i)` (where `baseSeed` is a fresh per-call random). A shared `*atomic.Bool` is wired into each worker's `cancelled` field; the construct loop checks it once per iteration. First success flips the flag and the losing workers exit with the internal `errCancelled` sentinel (filtered out before the final error is returned). Combining `WithParallelism(n > 1)` with `RandomSeed(...)` returns `ErrSeedWithParallelism` immediately — workers each need an independent seed and silently overriding the caller's would be surprising.

**Result construction (`wordsearch.go` + `solutions.go`)** runs after placement succeeds:
1. `translateToWordSearch` builds the public `WordSearch` (`NumCols`/`NumRows`, `SolutionRows`, `PuzzleRows`, `WordPlacements`, `TotalLetters`, `LetterDensity`).
2. Filler is applied per the chosen option (`applyUniformFiller` or `applyWeightedFiller` — weighted uses binary search over a cumulative-sum table).
3. `findAllPossibleSolutions` re-scans the now-filled puzzle to detect *additional* placements the filler may have created (e.g. random letters happening to spell a word). Words are scanned longest-first so a shorter word that's a substring of a longer one isn't double-counted. Results populate `AllPossibleWordPlacements` and `NumWordsWithMultipleSolutions`.

`TotalLetters` is the sum of `seq.Len()` across the *deduplicated* word list. `LetterDensity` is `TotalLetters / (NumCols * NumRows)`; values around 1.1–1.3 are the high-variance "borderline" regime where the README's density section recommends `WithParallelism(n)`.

**Options (`options.go`)** use the functional-option pattern (`WordSearchOption func(*wordSearchConstructor) error`). Direction options *add* to `possibleDirections` (bitwise OR) except `UseAllDirections` which sets it outright. Default direction set when none specified is `NaturalLTRDirections` (Down + LTR horizontal/ascending/descending). Performance-related options: `WithTimeLimit(d)`, `WithMemoizationLimit(maxEntries)` (default 0 = disabled), `WithParallelism(n)` (default 1 = single-threaded; n=0 → `runtime.NumCPU()`).

## Errors

`Construct()` returns only the typed `ErrorString` constants defined in `api.go` (`ErrCannotFitWords`, `ErrSmallerThanMinimumSize`, `ErrWordIsTooLong`, `ErrFillerWeightNotPositive`, `ErrReachedTimeLimit`, `ErrSeedWithParallelism`). When adding error conditions, add a new constant rather than wrapping. `errCancelled` (lowercase) is an internal sentinel used by parallel workers — it never leaks to the public API; the parallel orchestrator filters it out.

## Test-coverage notes

A few non-obvious correctness paths have dedicated regression tests because the broader solutions tests don't exercise them:

- `TestAscendingPlacementAtRowZero` (`directions_test.go`) — `validPlacement` boundary off-by-one.
- `TestTotalLetters*` (`wordsearch_test.go`) — result-struct fields, including the deduplication semantics.
- `TestMemoization*` (`memoization_test.go`) — the `tryPlace` cache path that the default (cap=0) doesn't reach.
- `TestParallelism*` (`parallelism_test.go`) — `WithParallelism` + `RandomSeed` rejection (both option orderings) and a smoke test that parallel construction actually solves.

When adding new options or result fields, add a similar targeted test rather than relying on the existing high-level tests to incidentally exercise the new branch.

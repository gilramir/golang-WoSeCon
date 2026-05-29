// Copyright (c) 2025 by Gilbert Ramirez <gram@alumni.rice.edu>.
// All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package wosecon

import (
	"math/rand"
	"sync/atomic"
	"time"
)

// constructParallel races n independent searches, each with its own RNG
// seed. The first goroutine to find a solution wins; the others observe
// the shared cancelled flag and exit with errCancelled. If every worker
// fails, the first non-cancellation error is returned to the caller —
// typically ErrReachedTimeLimit or ErrCannotFitWords.
func constructParallel(numCols, numRows int, sequences []Sequence, opts []WordSearchOption, n int) (*WordSearch, error) {
	type result struct {
		ws  *WordSearch
		err error
	}

	var cancelled atomic.Bool
	results := make(chan result, n)

	// Choose one base seed per call so each invocation explores a
	// different region of the search space. Workers offset from the
	// base by their index, giving each an independent stream.
	baseSeed := rand.New(rand.NewSource(time.Now().UnixNano())).Int63()

	for i := 0; i < n; i++ {
		seed := baseSeed + int64(i)
		go func() {
			workerOpts := make([]WordSearchOption, 0, len(opts)+1)
			workerOpts = append(workerOpts, opts...)
			workerOpts = append(workerOpts, RandomSeed(seed))

			var ctor wordSearchConstructor
			ctor.cancelled = &cancelled
			if err := ctor.init(numCols, numRows, sequences, workerOpts...); err != nil {
				results <- result{nil, err}
				return
			}
			if err := ctor.construct(); err != nil {
				results <- result{nil, err}
				return
			}
			results <- result{ctor.translateToWordSearch(), nil}
		}()
	}

	var firstRealErr error
	for i := 0; i < n; i++ {
		r := <-results
		if r.err == nil {
			cancelled.Store(true)
			return r.ws, nil
		}
		if firstRealErr == nil && r.err != errCancelled {
			firstRealErr = r.err
		}
	}
	if firstRealErr != nil {
		return nil, firstRealErr
	}
	// All workers reported errCancelled with no solution. This should be
	// unreachable in practice — at least one worker would have to win
	// for the rest to be cancelled — but return something honest if it
	// somehow happens.
	return nil, ErrCannotFitWords
}

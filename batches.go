// Package batches provides a single function for processing a channel
// in batches.
//
// 	ctx := ctx.Background()
// 	ch := make(chan int)
//
// 	// Cancellation, up to 10 items, or after 5 seconds.
// 	batches := Batches(ctx, ch, 10, time.Second*5)
//
// 	go func() {
// 		for i := 0; i < 100; i++ {
// 			ch <- i
// 		}
// 	}()
//
// 	for batch := range batches {
// 		...
// 	}
package batches

import (
	"context"
	"time"
)

// Batches will read items from the values chan and write items in
// batches to the chan it returns.
//
// A slice is written to the returned chan when its length is size items
// or timeout is reached, whichever happens first. If ctx is cancelled a
// final batch of up to size items may be delivered.
func Batches[T any](ctx context.Context, values <-chan T, size int, timeout time.Duration) <-chan []T {
	batches := make(chan []T, 1)

	go func() {
		defer close(batches)

		for stop := false; !stop; {
			var batch []T
			expire := time.After(timeout)

			for {
				select {
				case <-ctx.Done():
					stop = true
					goto done
				case value, ok := <-values:
					if !ok {
						stop = true
						goto done
					}
					batch = append(batch, value)
					if len(batches) == size {
						goto done
					}
				case <-expire:
					goto done
				}
			}
		done:
			if len(batches) > 0 {
				batches <- batch
			}
		}
	}()

	return batches
}

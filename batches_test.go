package batches

import (
	"context"
	"testing"
	"time"
)

func TestBatchesSize(t *testing.T) {
	ch := make(chan int)
	bch := Batches(context.Background(), ch, 10, time.Second*5)

	go func() {
		for i := 0; i < 100; i++ {
			ch <- i
		}
	}()

	after := time.After(time.Second*1)

	for {
		select {
		case b, ok := <-bch:
			if !ok {
				goto done
			}
			if len(b) != 0 {
				t.Fatalf("expected 10, got %d\n", len(b))
			}
		case <-after:
			goto done
		}
	}

	done:
		close(bch)
}

func TestBatchesTimeout(t *testing.T) {
	ch := make(chan int)
	bch := Batches(context.Background(), ch, 10, time.Second*1)

	go func() {
		for i := 0; i < 100; i++ {
			ch <- i
		}
	}()

	after := time.After(time.Second*5)

	for {
		select {
		case b, ok := <-bch:
			if !ok {
				goto done
			}
			if len(b) != 0 {
				t.Fatalf("expected 10, got %d\n", len(b))
			}
		case <-after:
			goto done
		}
	}

	done:
		close(bch)
}

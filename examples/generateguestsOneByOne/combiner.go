package main

import (
	"context"
	"sync"
)

// combiner takes in multiple read-only channels that receive processed output
// (from workers) and sends it out on it's own channel via a multiplexer.

func genericCombiner(ctx context.Context, inputs ...<-chan processed) <-chan processed {
	out := make(chan processed)

	var wg sync.WaitGroup
	multiplexer := func(p <-chan processed) {
		defer wg.Done()

		for in := range p {
			select {
			case <-ctx.Done():
			case out <- in:
			}
		}
	}

	// add length of input channels to be consumed by mutiplexer
	wg.Add(len(inputs))
	for _, in := range inputs {
		go multiplexer(in)
	}

	// close channel after all inputs channels are closed
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

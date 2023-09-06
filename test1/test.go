package main

import (
	"fmt"
	"sync"
)

type idxPair struct {
	idx, val int
}

func main() {
	// add a done channel, an ability to stop the world by closing this.
	done := make(chan struct{})
	defer close(done)

	// create srcChan, this will be where the values go into the pipeline
	srcCh := make(chan idxPair)

	// create a slice of result channels, one for each of the go workers
	const numWorkers = 8
	resChans := make([]<-chan idxPair, numWorkers)

	// waitgroup to wait for all the workers to stop
	var wg sync.WaitGroup
	wg.Add(numWorkers)

	// start the workers, passing them each the src channel,
	// collecting the result channels they return
	for i := 0; i < numWorkers; i++ {
		resChans[i] = worker(done, &wg, srcCh)
	}

	// start a single goroutine to send values into the pipeline
	// all values are sent with an index, to be pieces back into order at the end.
	go func() {
		defer close(srcCh)
		for i := 1; i < 100; i++ {
			srcCh <- idxPair{idx: i, val: i}
		}
	}()

	// merge all the results channels into a single results channel
	// this channel is unordered.
	mergedCh := merge(done, resChans...)

	// order the values coming from the mergedCh according the the idxPair.idx field.
	orderedResults := order(100, mergedCh)

	// iterate over each of the ordered results
	for _, v := range orderedResults {
		fmt.Println(v)
	}
}

func order(len int, res <-chan idxPair) []int {
	results := make([]int, len)

	// collect all the values to order them
	for r := range res {
		results[r.idx] = r.val
	}

	return results
}

func worker(done <-chan struct{}, wg *sync.WaitGroup, src <-chan idxPair) <-chan idxPair {
	res := make(chan idxPair)

	go func() {
		defer wg.Done()
		defer close(res)
		sendValue := func(pair idxPair) {
			v := pair.val
			v *= v
			ip := idxPair{idx: pair.idx, val: v}
			select {
			case res <- ip:
			case <-done:
			}
		}

		for v := range src {
			sendValue(v)
		}
	}()

	return res
}

// example and explanation here: https://blog.golang.org/pipelines
func merge(done <-chan struct{}, cs ...<-chan idxPair) <-chan idxPair {
	var wg sync.WaitGroup
	out := make(chan idxPair)

	output := func(c <-chan idxPair) {
		defer wg.Done()
		for n := range c {
			select {
			case out <- n:
			case <-done:
				return
			}
		}
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

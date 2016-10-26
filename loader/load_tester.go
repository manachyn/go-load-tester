package loader

import (
	"math/rand"
	"sync"
	"time"
	"github.com/valyala/fasthttp"
)

type LoadTester struct {
	concurrency uint64
	stopch chan struct{}
}

func NewLoadTester(concurrency uint64) *LoadTester {
	t := &LoadTester{
		stopch: make(chan struct{}),
		concurrency: concurrency,
	}

	return t
}

func (t *LoadTester) Load(urls []string, rps uint64, requests uint64) <-chan *Result {
	var wg sync.WaitGroup
	results := make(chan *Result)
	ticks := make(chan time.Time)
	for i := uint64(0); i < t.concurrency; i++ {
		wg.Add(1)
		go t.runWorker(urls, &wg, ticks, results)
	}

	go func() {
		defer close(results)
		defer wg.Wait()
		defer close(ticks)
		interval := 1e9 / rps
		began, done := time.Now(), uint64(0)
		for {
			now, next := time.Now(), began.Add(time.Duration(done*interval))
			time.Sleep(next.Sub(now))
			select {
			case ticks <- maxTime(next, now):
				if done++; done == requests {
					return
				}
			case <-t.stopch:
				return
			default:
				wg.Add(1)
				go t.runWorker(urls, &wg, ticks, results)
			}
		}
	}()

	return results
}

func (t *LoadTester) Stop() {
	select {
	case <-t.stopch:
		return
	default:
		close(t.stopch)
	}
}

func (t *LoadTester) makeRequest(url string, tm time.Time) *Result {
	var (
		res = Result{Timestamp: tm}
		err error
	)

	defer func() {
		res.Duration = time.Since(tm)
		if err != nil {
			res.Error = err.Error()
		}
	}()

	statusCode, _, err := fasthttp.Get(nil, url)
	if err != nil {
		return &res
	}
	res.Code = uint16(statusCode)
	if statusCode != fasthttp.StatusOK {
		res.Error = fasthttp.StatusMessage(statusCode)
	}

	return &res
}

func (t *LoadTester) runWorker(urls []string, wg *sync.WaitGroup, ticks <-chan time.Time, results chan<- *Result) {
	defer wg.Done()
	for tm := range ticks {
		results <- t.makeRequest(urls[rand.Intn(len(urls))], tm)
	}
}

func maxTime(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}
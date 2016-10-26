package main

import (
	"flag"
	"github.com/manachyn/go-load-tester/loader"
	"os"
	"os/signal"
	//"fmt"
)

var (
	requests = flag.Uint64("n", 1000, "Number of requests")
	concurrency = flag.Uint64("c", 100, "Number of concurrent requests")
	rps = flag.Uint64("q", 100, "RPS")
)

func main() {
	flag.Parse()

	tester := loader.NewLoadTester(*concurrency)
	urls := []string{"https://www.google.com.ua", "https://www.yandex.ua"}
	stats := loader.NewStats()
	results := tester.Load(urls, *rps, *requests)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	for {
		select {
		case res, more := <-results:
			if more {
				stats.Add(res)
				//fmt.Println(res)
			} else {
				stats.Print()
				return
			}
		case <-sig:
			tester.Stop()
			stats.Print()
			return

		}
	}
}
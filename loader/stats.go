package loader

import (
	"time"
	"fmt"
)

type Stats struct {
	Start time.Time
	Finish time.Time
	Duration time.Duration
	TotalRequestsDuration time.Duration
	AvgRequestDuration time.Duration
	Requests uint64
	Rps float64
	StatusCodes map[string]int
	Errors []string
	SuccessRequests uint64
}

func NewStats() *Stats {
	return &Stats{}
}

func (s *Stats) Add(r *Result) {
	s.Requests++
	s.TotalRequestsDuration += r.Duration
	if s.Start.IsZero() || s.Start.After(r.Timestamp) {
		s.Start = r.Timestamp
	}
	if end := r.End(); end.After(s.Finish) {
		s.Finish = end
	}
	if r.Code >= 200 && r.Code < 400 {
		s.SuccessRequests++
	}
}

func (s *Stats) Print() {
	s.aggregate()
	fmt.Println("Statistic:")
	fmt.Printf("Duration:\t%.4f sec\n", s.Duration.Seconds())
	fmt.Printf("Requests:\t%d\n", s.Requests)
	fmt.Printf("Success sequests:\t%d\n", s.SuccessRequests)
	fmt.Printf("Requests per second:\t%.2f r/sec\n", s.Rps)
	fmt.Printf("AVG request duration:\t%.4f sec\n", s.AvgRequestDuration.Seconds())
}

func (s *Stats) aggregate() {
	s.Duration = s.Finish.Sub(s.Start)
	s.Rps = float64(s.Requests) / s.Duration.Seconds()
	s.AvgRequestDuration = time.Duration(float64(s.TotalRequestsDuration) / float64(s.Requests))
}
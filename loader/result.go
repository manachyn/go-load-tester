package loader

import "time"

type Result struct {
	Timestamp time.Time
	Duration time.Duration
	Code uint16
	Error string
}

func (r *Result) End() time.Time {
	return r.Timestamp.Add(r.Duration)
}

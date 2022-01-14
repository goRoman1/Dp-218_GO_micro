package services
//go:generate mockgen -source=clock_to_mock.go -destination=../mock/clock_mock.go -package=mock

import "time"

//Clock creates the wrapper above the internal Time service
type Clock interface {
	Now() time.Time
}
type realClock struct{}
func (realClock) Now() time.Time { return time.Now() }

//NewClock returns the realClock of type Clock.
func NewClock() Clock {
	return &realClock{}
}
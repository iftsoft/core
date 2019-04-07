package core

import "time"

// Used to get duration of process
type DurationTimer struct {
	start time.Time
}

// Save time of process beginning
func (timer *DurationTimer) StartTimer() {
	timer.start = time.Now()
}

// Get process duration in nanoseconds
func (timer *DurationTimer) Nanoseconds() int64 {
	if !timer.start.IsZero() {
		delta := time.Now().Sub(timer.start)
		return int64(delta)
	}
	return 0
}

// Get process duration in microseconds
func (timer *DurationTimer) Microseconds() int {
	if !timer.start.IsZero() {
		delta := time.Now().Sub(timer.start)
		return int(delta / time.Microsecond)
	}
	return 0
}

// Get process duration in milliseconds
func (timer *DurationTimer) Milliseconds() int {
	if !timer.start.IsZero() {
		delta := time.Now().Sub(timer.start)
		return int(delta / time.Millisecond)
	}
	return 0
}

// Get process duration in seconds
func (timer *DurationTimer) Seconds() int {
	if !timer.start.IsZero() {
		delta := time.Now().Sub(timer.start)
		return int(delta / time.Second)
	}
	return 0
}

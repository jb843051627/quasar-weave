package clock

import "time"

type Clock interface {
	Now() time.Time
	Sleep(time.Duration)
}

type System struct{}

func (System) Now() time.Time { return time.Now().UTC() }

func (System) Sleep(duration time.Duration) { time.Sleep(duration) }

type Fixed struct {
	Mu chan struct{}
	At time.Time
}

func NewFixed(at time.Time) *Fixed {
	return &Fixed{Mu: make(chan struct{}, 1), At: at.UTC()}
}

func (f *Fixed) Now() time.Time {
	return f.At
}

func (f *Fixed) Sleep(duration time.Duration) {
	f.At = f.At.Add(duration)
}

func Since(c Clock, at time.Time) time.Duration {
	if at.IsZero() {
		return 0
	}
	return c.Now().Sub(at)
}

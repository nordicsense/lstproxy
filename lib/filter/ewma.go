package filter

import (
	"math"
	"time"
)

type ewma struct {
	t           time.Time
	v           float64
	burnInSince time.Time
	halfLifeSec int
}

func EWMA(halfLifeSec int, future time.Time) *ewma {
	// Time does not matter if value is NaN: the logic in Compute will iterate until a good value found.
	return &ewma{
		t:           time.Now(),
		v:           math.NaN(),
		burnInSince: future,
		halfLifeSec: halfLifeSec,
	}
}

func (s *ewma) Compute(t time.Time, v float64) float64 {
	// No new value: return NaN and keep burning.
	// A good value in the future will be correctly handled by the clauses below.
	if math.IsNaN(v) {
		return math.NaN()
	}
	dt := float64(t.Unix() - s.t.Unix())
	s.t = t
	// If no state value assume new value as is and restart burning.
	// If no good values for long period assume new value as is and restart burning.
	if math.IsNaN(s.v) || dt >= 0.75*float64(s.halfLifeSec) {
		s.v = v
		s.burnInSince = t
	} else {
		// Normal flow, compute new value and keep burning.
		alpha := 1.0 - math.Exp(-dt/float64(s.halfLifeSec))
		s.v += alpha * (v - s.v)
	}
	if !math.IsNaN(s.v) && s.burnInSince.Add(time.Duration(int64(s.halfLifeSec)*int64(time.Second))).Before(s.t) {
		return s.v
	}
	return math.NaN()
}

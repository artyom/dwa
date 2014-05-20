package dwa

import (
	"sync"
	"time"
)

// DWA (decaying window average) is a decaying simple moving average. It differs
// from simple moving average by decaying factor: if samples do not get updated
// for a while, old samples are removed as time goes by.
type DWA struct {
	m           sync.Mutex
	lastUpdate  time.Time
	decayPeriod time.Duration
	idx         int
	sampleSize  int
	values      []int64
	sum         int64
}

// Add adds given items to DWA sample
func (d *DWA) Add(values ...int64) {
	d.m.Lock()
	defer d.m.Unlock()
	for _, v := range values {
		if d.idx < d.sampleSize && len(d.values) <= d.idx {
			d.values = append(d.values, v)
			d.idx++
			d.sum += v
			continue
		}
		d.idx = d.idx % d.sampleSize
		d.sum -= d.values[d.idx]
		d.values[d.idx] = v
		d.sum += v
		d.idx++
	}
	d.lastUpdate = time.Now()
}

// Value returns average of current DWA sample
func (d *DWA) Value() float64 {
	now := time.Now()
	d.m.Lock()
	defer d.m.Unlock()
	if len(d.values) == 0 {
		return 0
	}
	// as we don't use any background goroutines to maintain decaying, we're
	// zeroing out items here
	if d.decayPeriod > 0 && now.After(d.lastUpdate.Add(d.decayPeriod)) {
		n := int(now.Sub(d.lastUpdate) / d.decayPeriod)
		if n > len(d.values) {
			d.values = d.values[:0]
			d.sum = 0
			d.idx = 0
			return 0
		}
		for i := 0; i < n; i++ {
			d.idx = d.idx % d.sampleSize
			d.sum -= d.values[d.idx]
			d.values[d.idx] = 0
			d.idx++
		}
		d.lastUpdate = d.lastUpdate.Add(time.Duration(n) * d.decayPeriod)
	}
	return float64(d.sum) / float64(len(d.values))
}

// NewDWA returns initialized DWA with given sample size and decay period. For
// each decay period DWA does not receive an update, the oldest item would be
// zeroed.
//
// Function panics if sample size is non-positive or decay period is negative.
// Decay period of 0 means no decay.
func NewDWA(sampleSize int, decayPeriod time.Duration) *DWA {
	if sampleSize < 1 {
		panic("sample size should be positive")
	}
	if decayPeriod < 0 {
		panic("decay period cannot be negative")
	}
	return &DWA{
		sampleSize:  sampleSize,
		decayPeriod: decayPeriod,
		values:      make([]int64, 0, sampleSize),
	}
}

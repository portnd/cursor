package repository

import (
	"sync"
	"time"

	"github.com/portnd/the-sentinel-core/internal/modules/sentinel/domain"
)

const (
	defaultLimitRPM = 15  // free tier typical
	defaultLimitRPD = 250
	maxStored       = 2000
)

type memoryUsageTracker struct {
	mu    sync.Mutex
	times []time.Time
}

// NewMemoryUsageTracker returns a thread-safe usage tracker for Gemini API calls.
func NewMemoryUsageTracker() domain.UsageTracker {
	return &memoryUsageTracker{times: make([]time.Time, 0, 128)}
}

func (t *memoryUsageTracker) RecordRequest() {
	t.mu.Lock()
	defer t.mu.Unlock()
	now := time.Now().UTC()
	t.times = append(t.times, now)
	// Keep only last 24h of data; also cap length
	cut := time.Now().UTC().Add(-24 * time.Hour)
	for len(t.times) > 0 && t.times[0].Before(cut) {
		t.times = t.times[1:]
	}
	if len(t.times) > maxStored {
		t.times = t.times[len(t.times)-maxStored:]
	}
}

func (t *memoryUsageTracker) GetUsage(limitRPM, limitRPD int) domain.AIUsage {
	t.mu.Lock()
	defer t.mu.Unlock()
	if limitRPM <= 0 {
		limitRPM = defaultLimitRPM
	}
	if limitRPD <= 0 {
		limitRPD = defaultLimitRPD
	}
	now := time.Now().UTC()
	windowMin := now.Add(-1 * time.Minute)
	midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	var lastMin, today int
	for _, ts := range t.times {
		if ts.After(windowMin) {
			lastMin++
		}
		if ts.After(midnight) {
			today++
		}
	}
	remainingRPM := limitRPM - lastMin
	if remainingRPM < 0 {
		remainingRPM = 0
	}
	remainingRPD := limitRPD - today
	if remainingRPD < 0 {
		remainingRPD = 0
	}
	return domain.AIUsage{
		RequestsLastMinute: lastMin,
		RequestsToday:      today,
		LimitRPM:           limitRPM,
		LimitRPD:           limitRPD,
		RemainingRPM:       remainingRPM,
		RemainingRPD:       remainingRPD,
	}
}

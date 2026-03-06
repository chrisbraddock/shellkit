package data

import (
	"bufio"
	"encoding/json"
	"math"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// MetricEntry represents a single shell startup timing record.
type MetricEntry struct {
	Timestamp  time.Time `json:"ts"`
	DurationMs float64   `json:"duration_ms"`
	Host       string    `json:"host"`
}

// MetricsSummary holds computed statistics over the metrics history.
type MetricsSummary struct {
	Entries []MetricEntry
	Count   int
	Current float64
	Average float64
	Median  float64
	Min     float64
	Max     float64
	P95     float64
}

// LoadMetrics reads the JSONL metrics file and returns parsed entries.
// Returns nil, nil if the file does not exist.
func LoadMetrics(homeDir string) ([]MetricEntry, error) {
	path := filepath.Join(homeDir, ".local", "share", "shellkit", "metrics.jsonl")

	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()

	var entries []MetricEntry
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var entry MetricEntry
		if err := json.Unmarshal(line, &entry); err != nil {
			continue // skip malformed lines
		}
		entries = append(entries, entry)
	}

	return entries, scanner.Err()
}

// ComputeSummary calculates statistics from a slice of MetricEntry.
func ComputeSummary(entries []MetricEntry) MetricsSummary {
	s := MetricsSummary{
		Entries: entries,
		Count:   len(entries),
	}

	if len(entries) == 0 {
		return s
	}

	// Extract durations
	durations := make([]float64, len(entries))
	var sum float64
	for i, e := range entries {
		durations[i] = e.DurationMs
		sum += e.DurationMs
	}

	s.Current = durations[len(durations)-1]
	s.Average = sum / float64(len(durations))

	// Sort for min/max/median/p95
	sorted := make([]float64, len(durations))
	copy(sorted, durations)
	sort.Float64s(sorted)

	s.Min = sorted[0]
	s.Max = sorted[len(sorted)-1]
	s.Median = percentile(sorted, 0.5)
	s.P95 = percentile(sorted, 0.95)

	return s
}

// percentile returns the value at the given percentile from a sorted slice.
func percentile(sorted []float64, p float64) float64 {
	if len(sorted) == 0 {
		return 0
	}
	if len(sorted) == 1 {
		return sorted[0]
	}
	idx := p * float64(len(sorted)-1)
	lower := int(math.Floor(idx))
	upper := int(math.Ceil(idx))
	if lower == upper {
		return sorted[lower]
	}
	frac := idx - float64(lower)
	return sorted[lower]*(1-frac) + sorted[upper]*frac
}

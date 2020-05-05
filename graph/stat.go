package graph

import (
	"runtime"
	"time"
)

// Stat contains some information about the graph.
type Stat struct {
	StartTime     time.Time
	NodeCount     int
	EdgeCount     int
	NumCPU        int
	NumGoroutings int
	MemStats      runtime.MemStats
}

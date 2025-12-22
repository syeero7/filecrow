package main

import (
	"fmt"
)

func humanReadSize(s int64) string {
	const unit = 1000
	if s < unit {
		return fmt.Sprintf("%d B", s)
	}

	div, exp := int64(unit), 0
	for n := s / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.2f %cB", float64(s)/float64(div), "kMGTPE"[exp])
}

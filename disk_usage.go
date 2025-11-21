package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"runtime"

	"github.com/shirou/gopsutil/v4/disk"
)

type DiskUsage struct {
	Total      string
	Free       string
	Used       string
	Percentage string
}

func getOSMountPoint() string {
	if runtime.GOOS == "windows" {
		dir, _ := os.UserHomeDir()
		return filepath.VolumeName(dir)

	}

	return "/"
}

func getDiskUsage() *DiskUsage {
	usage, err := disk.Usage(getOSMountPoint())
	if err != nil {
		log.Println(err)
	}

	return &DiskUsage{
		Free:       humanReadSize(safeUint64ToInt64(usage.Free)),
		Used:       humanReadSize(safeUint64ToInt64(usage.Used)),
		Total:      humanReadSize(safeUint64ToInt64(usage.Total)),
		Percentage: fmt.Sprintf("%.2f%%", usage.UsedPercent),
	}
}

func safeUint64ToInt64(v uint64) int64 {
	if v > math.MaxInt64 {
		log.Println("uint64 value exceeds math.MaxInt64")
		return 0
	}

	return int64(v)
}

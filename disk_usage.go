package main

import (
	"fmt"
	"log"
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
		Free:       humanReadSize(int64(usage.Free)),
		Used:       humanReadSize(int64(usage.Used)),
		Total:      humanReadSize(int64(usage.Total)),
		Percentage: fmt.Sprintf("%.2f%%", usage.UsedPercent),
	}
}

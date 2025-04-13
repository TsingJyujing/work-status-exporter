package workstatus

import (
	"fmt"
	"github.com/shirou/gopsutil/v4/process"
	"strings"
	"work-status-exporter/logging"
)

func GetZoomMeetingStatus() (bool, error) {
	processList, err := process.Processes()
	if err != nil {
		return false, fmt.Errorf("failed to get process list: %w", err)
	}
	for _, proc := range processList {
		cmdline, err := proc.Cmdline()
		if err != nil {
			logging.Logger.WithError(err).Tracef("failed to get process cmdline: %v", proc)
			continue
		}
		if strings.Contains(cmdline, "zoom.us") {
			if strings.Contains(cmdline, "CptHost") {
				return true, nil
			}
		}
	}
	return false, nil
}

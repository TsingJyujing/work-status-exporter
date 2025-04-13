package workstatus

import (
	"fmt"
	"os/exec"
	"strings"
	"work-status-exporter/logging"
)

func GetMacOSCameraStatus() (bool, bool, error) {
	cmdStdOut, cmdError := exec.Command("ioreg", "-n", "AppleH13CamIn").Output()
	if cmdError != nil {
		return false, false, cmdError
	}
	activated, streaming := false, false
	hasActivated, hasStreaming := false, false
	for _, row := range strings.Split(string(cmdStdOut), "\n") {
		if strings.Contains(row, "FrontCameraActive") {
			activated = strings.Contains(row, "Yes")
			hasActivated = true
			logging.Logger.Debugf(fmt.Sprintf("Found in ioreg result: %s", row))
		}
		if strings.Contains(row, "FrontCameraStreaming") {
			streaming = strings.Contains(row, "Yes")
			hasStreaming = true
			logging.Logger.Debugf(fmt.Sprintf("Found in ioreg result: %s", row))
		}
	}
	if !hasStreaming || !hasActivated {
		return activated, streaming, fmt.Errorf("no streaming or activated camera found, found activated = %v and found streaming = %v", hasActivated, hasStreaming)
	}
	return activated, streaming, nil
}

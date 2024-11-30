package logs

import (
	"fmt"
	"time"

	"github.com/ROHITHSAKTHIVEL/GoatRobotics/models"
)

var (
	Logs map[time.Time]models.Logs
)

func init() {
	Logs = make(map[time.Time]models.Logs)
}

// AddLog adds a log entry to the Logs map
func AddLog(logEntry models.Logs) {
	Logs[logEntry.StartTime] = logEntry
}

// GetLog retrieves the log entry for a given clientID or timestamp
func GetLog(clientID string, startTime time.Time) ([]models.Logs, error) {
	var result []models.Logs
	for _, log := range Logs {
		if log.ClientID == clientID || log.StartTime == startTime {
			result = append(result, log)
		}
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no logs found for the given clientID or timestamp")
	}

	return result, nil
}

func GetAllLogs() ([]models.Logs, error) {
	var result []models.Logs

	// Iterate over the map and append logs to the result slice
	for _, log := range Logs {
		result = append(result, log)
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no logs available")
	}

	return result, nil
}
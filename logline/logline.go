package logline

import (
	"time"
	"strconv"
)

type LogLine struct {
	Timestamp string `json:"__REALTIME_TIMESTAMP"`
	Unit string `json:"_SYSTEMD_UNIT"`
	Message string `json:"MESSAGE"`
}

func (logLine *LogLine) FormatTimestamp() (string) {
	timeInt64, err := strconv.ParseInt(logLine.Timestamp, 10, 64)
	if err != nil {
		return "ERROR:time-unparsable"
	}
	// NanoTime is... correct? How does Go even modulus?
	return time.Unix(timeInt64/1000000, timeInt64%1000000).Format(time.RFC3339Nano)
}

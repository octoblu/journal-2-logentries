package journal

import (
	"time"
	"fmt"
	"strings"
	"strconv"
)

type LogLine struct {
	Timestamp string `json:"__REALTIME_TIMESTAMP"`
	Unit string `json:"_SYSTEMD_UNIT"`
	Message string `json:"MESSAGE"`
	Hostname string `json:"_HOSTNAME"`
}

func (logLine *LogLine) FormatTimestamp() (string) {
	timeInt64, err := strconv.ParseInt(logLine.Timestamp, 10, 64)
	if err != nil {
		return "ERROR:time-unparsable"
	}
	// NanoTime is... correct? How does Go even modulus?j
	return time.Unix(timeInt64/1000000, timeInt64%1000000).Format(time.RFC3339Nano)
}

func (logLine *LogLine) FormatLine() (string) {
	prettyOutput := fmt.Sprintf("%s %s %s %s", logLine.FormatTimestamp(), logLine.Hostname, logLine.Unit, logLine.Message)
	return strings.Replace(prettyOutput, "\n", "\\n", -1)
}

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
	Message string
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

func (logLine *LogLine) Parse(obj interface{}) error {
	objMap := obj.(map[string]interface{})

	message, err := parseMessage(objMap["MESSAGE"])
	if err != nil {
		return err
	}

	logLine.Message = message
	logLine.Timestamp = parseStringOrNil(objMap["__REALTIME_TIMESTAMP"])
	logLine.Unit = parseStringOrNil(objMap["_SYSTEMD_UNIT"])
	logLine.Hostname = parseStringOrNil(objMap["_HOSTNAME"])

	return nil
}

func (logLine *LogLine) SetError(errorStr string) {
	logLine.Timestamp = fmt.Sprintf("%v", int32(time.Now().Unix()))
	logLine.Message = errorStr
	logLine.Unit = "logline.go"
}

func parseByteArray(byteArrayInterface []interface{}) (string, error) {
	byteArray := make([]byte, len(byteArrayInterface))

	for index, value := range byteArrayInterface {
		byteArray[index] = byte(value.(float64))
	}

	return string(byteArray), nil
}

func parseMessage(msgInterface interface{}) (string, error) {
	switch msgInterface.(type) {
	case string:
		return msgInterface.(string), nil
	case []interface{}:
		return parseByteArray(msgInterface.([]interface{}))
	}
	return "", fmt.Errorf("Could not parse MESSAGE: %v", msgInterface)
}

func parseStringOrNil(strInterface interface{}) (string) {
	switch strInterface.(type) {
	case string:
		return strInterface.(string)
	}
	return ""
}

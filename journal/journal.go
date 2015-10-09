package journal

import (
	"fmt"
	"net"
	"net/http"
	"encoding/json"
	"github.com/octoblu/journal-2-logentries/logline"
)

const DefaultSocket = "/run/journald.sock"

func Follow(socket string) (<-chan logline.LogLine, error) {
	if socket == "" {
		socket = DefaultSocket
	}
	c := &http.Client{
		Transport: &http.Transport{
			Dial: func(network, addr string) (net.Conn, error) {
				return net.Dial("unix", socket)
			},
		},
	}
	req, err := http.NewRequest("GET", "http://journal/entries?follow", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non 200 response: %d", resp.StatusCode)
	}
	logs := make(chan logline.LogLine)
	decoder := json.NewDecoder(resp.Body)
	go func() {
		var logLine logline.LogLine
		for decoder.More() {
			if err := decoder.Decode(&logLine); err != nil {
				// MESSAGE might be a byte array... WHAT DO I DO?!
				logLine.Message = "ERROR:journal-2-logentries " + err.Error()
			}
			logs <- logLine
		}
	}()
	return logs, nil
}

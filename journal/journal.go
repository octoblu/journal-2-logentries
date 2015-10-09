package journal

import (
	"fmt"
	"net"
	"net/http"
	"encoding/json"
)

const DefaultSocket = "/run/journald.sock"

func Follow(socket string) (<-chan []byte, error) {
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
	logs := make(chan []byte)
	decoder := json.NewDecoder(resp.Body)
	go func() {
		var logLine LogLine
		for decoder.More() {
			if err := decoder.Decode(&logLine); err != nil {
				// MESSAGE might be a byte array... WHAT DO I DO?!
				logLine.Message = "ERROR:journal-2-logentries " + err.Error()
			}
			logs <- []byte(logLine.FormatLine())
		}
	}()
	return logs, nil
}

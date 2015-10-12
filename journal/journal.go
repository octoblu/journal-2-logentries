package journal

import (
	"fmt"
	"log"
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
		var obj interface{}
		var logLine LogLine
		for decoder.More() {

			if err := decoder.Decode(&obj); err != nil {
				log.Println(err.Error())
				logLine.SetError(err.Error())
			}

			if err := logLine.Parse(obj); err != nil {
				log.Fatalf("ERROR: %v", err.Error())
				logLine.SetError(err.Error())
			}

			logs <- []byte(logLine.FormatLine())
		}
	}()
	return logs, nil
}

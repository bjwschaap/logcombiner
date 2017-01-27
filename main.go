package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

// LogLine defines a single json log event
type LogLine struct {
	Time   string `json:"time"`
	Log    string `json:"log"`
	Stream string `json:"stream"`
}

func main() {
	logFile, err := os.Open("./docker_log_json_stacktrace.txt")
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(logFile)
	for scanner.Scan() {
		line := scanner.Text()
		var logLine LogLine
		err := json.Unmarshal([]byte(line), &logLine)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(strings.TrimRight(logLine.Log, "\n"))
	}
}

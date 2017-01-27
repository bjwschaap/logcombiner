package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/oleiade/lane"
)

// LogLine defines a single json log event
type LogLine struct {
	Time   string `json:"time"`
	Log    string `json:"log"`
	Stream string `json:"stream"`
}

var timeStampRegex = regexp.MustCompile(`^\[\d{4}-\d{1,2}-\d{1,2}.*\]`)

func main() {
	startTime := time.Now()
	logFile, err := os.Open("./docker_log_json_stacktrace.txt")
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(logFile)
	buffer := lane.NewQueue()
	for scanner.Scan() {
		line := scanner.Text()
		var logLine LogLine
		err := json.Unmarshal([]byte(line), &logLine)
		if err != nil {
			log.Fatal(err)
		}
		if isPrimaryLine(logLine.Log) && buffer.Size() > 0 {
			// flush out buffer to single new event
			newLogEvent, err := json.Marshal(flush(buffer))
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(newLogEvent))
		}
		buffer.Enqueue(logLine)
	}
	if buffer.Size() > 0 {
		newLogEvent, err := json.Marshal(flush(buffer))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(newLogEvent))
	}
	executionTime := time.Since(startTime)
	fmt.Printf("Parsing took: %dms\n", int64(executionTime.Nanoseconds()/1000000))
}

func isPrimaryLine(line string) bool {
	if timeStampRegex.FindString(line) != "" {
		return true
	}
	return false
}

func flush(buffer *lane.Queue) *LogLine {
	var stringBuffer bytes.Buffer
	time := buffer.Head().(LogLine).Time
	stream := buffer.Head().(LogLine).Stream
	for buffer.Head() != nil {
		logLine := buffer.Dequeue().(LogLine)
		stringBuffer.WriteString(logLine.Log)
	}
	return &LogLine{Time: time, Stream: stream, Log: stringBuffer.String()}
}

package combiner

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/oleiade/lane"
	"github.com/urfave/cli"
)

// LogLine defines a single json log event
type LogLine struct {
	Time   string `json:"time"`
	Log    string `json:"log"`
	Stream string `json:"stream"`
}

// Start is the main entrypoint for the CLI program
func Start(c *cli.Context) error {
	startTime := time.Now()
	timeStampRegex, err := regexp.Compile(c.String("regex"))
	if err != nil {
		return err
	}
	if c.String("logfile") == "" {
		cli.ShowAppHelp(c)
		fmt.Println("")
		return errors.New("--logfile (-l) is a required parameter")
	}
	logFile, err := os.Open(c.String("logfile"))
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(logFile)
	buffer := lane.NewQueue()
	for scanner.Scan() {
		line := scanner.Text()
		var logLine LogLine
		err := json.Unmarshal([]byte(line), &logLine)
		if err != nil {
			return err
		}
		if isPrimaryLine(logLine.Log, timeStampRegex) && buffer.Size() > 0 {
			// flush out buffer to single new event
			newLogEvent, err := json.Marshal(flush(buffer))
			if err != nil {
				return err
			}
			fmt.Println(string(newLogEvent))
		}
		buffer.Enqueue(logLine)
	}
	if buffer.Size() > 0 {
		newLogEvent, err := json.Marshal(flush(buffer))
		if err != nil {
			return err
		}
		fmt.Println(string(newLogEvent))
	}
	executionTime := time.Since(startTime)
	fmt.Printf("Parsing took: %dms\n", int64(executionTime.Nanoseconds()/1000000))
	return nil
}

func isPrimaryLine(line string, regex *regexp.Regexp) bool {
	if regex.FindString(line) != "" {
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

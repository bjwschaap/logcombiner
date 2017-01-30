package main

import (
	"log"
	"os"

	"github.com/bjwschaap/logcombiner/combiner"
	"github.com/urfave/cli"
)

func main() {
	log.SetOutput(os.Stdout)

	app := cli.NewApp()
	app.Name = "logcombiner"
	app.Usage = "Simple utility program to combine multiline JSON log-events into single line JSON log-events"
	app.Version = "0.1.0"
	app.Copyright = "(C)2017 Capgemini - PLP"
	app.Author = "Bastiaan Schaap"
	app.Email = "plp.nl@capgemini.com"
	app.UsageText = `E.g.: ./logcombiner -l mylog.txt -r '^[\\d{4'
		Specify the file to be parsed by using --logfile or -l. Optionally pass in
		--regex or -r to specify a regular expression to be used for detecting the
		first line of a multiline log event.`

	// Define the configuration flags the program can/should use
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "logfile,l",
			Usage:  "Path (relative or absolute) to the `LOGFILE` to be parsed",
			EnvVar: "PARSE_LOG_FILE",
		},
		cli.StringFlag{
			Name:   "regex,r",
			Value:  `^\[\d{4}-\d{1,2}-\d{1,2}.*\]`,
			Usage:  "`REGEX` to use for detecting a new log event",
			EnvVar: "PARSE_REGEX",
		},
	}

	// Set the main program logic
	app.Action = func(c *cli.Context) error {
		return combiner.Start(c)
	}

	// Now start doing stuff
	app.Run(os.Args)
}

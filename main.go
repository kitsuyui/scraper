package main

import (
	"fmt"
	"io"
	"os"

	"github.com/docopt/docopt-go"
	"github.com/kitsuyui/scraper/scraper"
	"github.com/kitsuyui/scraper/server"
)

var exit = os.Exit
var createServer = server.CreateServer
var standardOutput io.Writer = os.Stdout
var errorOutput io.Writer = os.Stderr

const (
	exitValidateConfig = 2
	exitServer         = 3
	exitInputFile      = 4
	exitOutputFile     = 5
	exitConfigFile     = 6
	exitScrape         = 7
)

var usage = `scraper: Swiss army knife for web scraping

Usage:
 scraper [-c <conf>] [-i <input>] [-o <output>]
 scraper validate [-c <conf>]
 scraper server [-d <conf-dir>] [-H <host>] [-p <port>]

Options:
 -c=<conf> --config=<conf>            Configuration file [default: scraper-config.json].
 -i=<input> --input=<input>           Input file specified instead of STDIN.
 -o=<output> --output=<output>        Output file specified instead of STDOUT.
 -H=<host> --host=<host>              Server mode host [default: 127.0.0.1].
 -p=<port> --port=<port>              Server mode port [default: 8080].
 -d=<conf-dir> --conf-dir=<conf-dir>  Configuration directory for server mode [default: .].
`

func reportError(err error, exitCode int) {
	fmt.Fprintln(errorOutput, err.Error())
	exit(exitCode)
}

func main() {
	opts, _ := docopt.ParseDoc(usage)

	if validate, _ := opts.Bool("validate"); validate {
		configFilepath, _ := opts.String("--config")
		c, err := os.Open(configFilepath)
		if err != nil {
			reportError(err, exitConfigFile)
			return
		}
		defer c.Close()
		err = scraper.ValidateConfigFile(c)
		if err != nil {
			reportError(err, exitValidateConfig)
			return
		}
	} else if serverMode, _ := opts.Bool("server"); serverMode {
		host, _ := opts.String("--host")
		port, _ := opts.Int("--port")
		confDir, _ := opts.String("--conf-dir")
		s, err := createServer(host, port, confDir)
		if err != nil {
			reportError(err, exitServer)
			return
		}
		if err := s.ListenAndServe(); err != nil {
			reportError(err, exitServer)
		}
		return
	} else {
		var confFile io.Reader
		var input io.Reader
		var output io.Writer

		input = os.Stdin
		output = standardOutput
		if inputFilepath, err := opts.String("--input"); err == nil {
			inputFile, err := os.Open(inputFilepath)
			if err != nil {
				reportError(err, exitInputFile)
				return
			}
			defer inputFile.Close()
			input = inputFile
		}

		if outputFilepath, err := opts.String("--output"); err == nil {
			outputFile, err := os.Create(outputFilepath)
			if err != nil {
				reportError(err, exitOutputFile)
				return
			}
			defer outputFile.Close()
			output = outputFile
		}

		configFilepath, _ := opts.String("--config")
		cf, err := os.Open(configFilepath)
		if err != nil {
			reportError(err, exitConfigFile)
			return
		}
		defer cf.Close()
		confFile = cf

		err = scraper.ScrapeByConfFile(confFile, input, output)
		if err != nil {
			reportError(err, exitScrape)
			return
		}
	}
}

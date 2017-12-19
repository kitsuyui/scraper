package commands

import (
	"fmt"
	"io"
	"os"

	"../scraper"
	"github.com/spf13/cobra"
)

var confFile io.Reader
var input io.Reader
var output io.Writer

var inputFilepath string
var outputFilepath string
var configFilepath string

var exit = os.Exit

func init() {
	CLIRoot.AddCommand(ValidateConfig)
	cobra.OnInitialize(initConfig)
	CLIRoot.PersistentFlags().StringVarP(
		&configFilepath, "configFile", "c", "scraper-config.json", "config file")
	CLIRoot.PersistentFlags().StringVarP(
		&inputFilepath, "in", "i", "", "input file (default STDIN)")
	CLIRoot.PersistentFlags().StringVarP(
		&outputFilepath, "out", "o", "", "output file (default STDOUT)")
}

func initConfig() {
	var err error
	if inputFilepath == "" {
		input = os.Stdin
	} else {
		input, err = os.Open(inputFilepath)
		if err != nil {
			fmt.Println(err.Error())
			exit(1)
		}
	}
	if outputFilepath == "" {
		output = os.Stdout
	} else {
		output, err = os.Open(outputFilepath)
		if err != nil {
			fmt.Println(err.Error())
			exit(1)
		}
	}
	confFile, err = os.Open(configFilepath)
	if err != nil {
		fmt.Println(err.Error())
		exit(1)
	}
}

// CLIRoot is top Cobra Object of scraper command
var CLIRoot = &cobra.Command{
	Use:  `scraper`,
	Long: `Scraper: Swiss Army Knife for Web scraping`,
	Run: func(cmd *cobra.Command, args []string) {
		err := scraper.ScrapeByConfFile(confFile, input, output)
		if err != nil {
			fmt.Println(err.Error())
			exit(1)
		}
	},
}

// ValidateConfig is a command for checking syntax errors in scraping configuration.
var ValidateConfig = &cobra.Command{
	Use:   `validate`,
	Short: `Validate config file`,
	Run: func(cmd *cobra.Command, args []string) {
		err := scraper.ValidateConfigFile(confFile)
		if err != nil {
			fmt.Println(err.Error())
			exit(1)
		}
		fmt.Println("OK")
	},
}

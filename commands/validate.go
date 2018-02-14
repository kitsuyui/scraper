package commands

import (
	"fmt"
	"os"

	"github.com/kitsuyui/scraper/scraper"
	"github.com/spf13/cobra"
)

func init() {
	ValidateConfig.Flags().StringVarP(
		&configFilepath, "configFile", "c", "scraper-config.json", "config file")
}

// ValidateConfig is a command for checking syntax errors in scraping configuration.
var ValidateConfig = &cobra.Command{
	Use:   `validate`,
	Short: `Validate config file`,
	Run: func(cmd *cobra.Command, args []string) {
		c, err := os.Open(configFilepath)
		if err != nil {
			fmt.Println(err.Error())
			exit(1)
		}
		confFile = c
		err = scraper.ValidateConfigFile(confFile)
		if err != nil {
			fmt.Println(err.Error())
			exit(1)
		}
		fmt.Println("OK")
	},
}

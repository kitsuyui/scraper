package commands

import (
	"fmt"

	"github.com/kitsuyui/scraper/server"
	"github.com/spf13/cobra"
)

var serveConfigDirectory string
var bindHost string
var bindPort int

func init() {
	Server.Flags().StringVarP(
		&serveConfigDirectory, "serveConfigDirectory", "d", ".", "config directory")
	Server.Flags().IntVarP(&bindPort, "port", "p", 8080, "bind port")
	Server.Flags().StringVarP(&bindHost, "host", "H", "127.0.0.1", "bind host")
}

var Server = &cobra.Command{
	Use:   `server`,
	Short: `Server mode`,
	Run: func(cmd *cobra.Command, args []string) {
		s, err := server.CreateServer(bindHost, bindPort, serveConfigDirectory)
		if err != nil {
			fmt.Println(err.Error())
			exit(1)
			return
		}
		s.ListenAndServe()
	},
}

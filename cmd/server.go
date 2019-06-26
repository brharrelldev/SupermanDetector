package cmd

import (
	"errors"
	"github.com/urfave/cli"
)

var ServerCmd = cli.Command{
	Name: "server",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name: "port",
			Usage: "port server will run on",
		},
		cli.StringFlag{
			Name: "mmdb",
		},
	},
	Before: func(app *cli.Context) error {
		if app.String("port") == "" || app.String("mmdb") == ""{
			return errors.New("port and mmdb should not be empty")
		}

		return nil
	},
	Action: func(app *cli.Context) {


	},
}

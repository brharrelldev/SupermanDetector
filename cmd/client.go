package cmd

import (
	"fmt"
	"github.com/brharrelldev/SupermanDetector/service"
	"github.com/urfave/cli"
	"log"
)

var ClientCmd = cli.Command{
	Name: "client",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name: "login",
			Usage: "login name to search",
		},
		cli.StringFlag{
			Name: "addr",
			Usage: "address of service",
		},
	},
	Before: func(app*cli.Context) error {
		if app.String("login") == "" || app.String("addr") == ""{
			if err := cli.ShowSubcommandHelp(app); err != nil{
				return fmt.Errorf("error displaying help %v", err)
			}
		}

		return nil

	},
	Action: func(app *cli.Context) {

		srv := service.Server{
			Port: app.String("port"),
		}

		if err := srv.StartServer(); err != nil{
			log.Fatalf("could not start server %v", err)
		}


	},

}

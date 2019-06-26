package cmd

import (
	"errors"
	"github.com/brharrelldev/SupermanDetector/databases"
	"github.com/brharrelldev/SupermanDetector/service"
	"github.com/urfave/cli"
	"log"
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
			EnvVar: "MMDB",
		},
		cli.StringFlag{
			Name: "logindb",
			EnvVar: "LOGIN_DB",

		},
	},
	Before: func(app *cli.Context) error {
		if app.String("port") == "" || app.String("mmdb") == ""{
			return errors.New("port and mmdb should not be empty")
		}

		return nil
	},
	Action: func(app *cli.Context) {

		srv := service.Server{
			Port: app.String("port"),
			SupermanDBs: &databases.SupermanDatabases{
				LoginDBClient: &databases.LoginDBClient{
					DBFile: app.String("logindb"),
				},
				MMDB: databases.MMDB{
					MMFile: app.String("mmdb"),
				},
			},

		}

		if err := srv.StartServer(); err != nil{
			log.Fatalf("could not start server %v", err)
		}
	},
}

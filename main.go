package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/brharrelldev/SupermanDetector/cmd"
	"github.com/brharrelldev/SupermanDetector/databases"
	"github.com/urfave/cli"
	"log"
	"os"
)

func main() {

	app := cli.NewApp()
	app.Name = "supes-cli"
	app.Flags = []cli.Flag{
		cli.StringFlag{ /**/
			Name:   "db-file",
			EnvVar: "LOGIN_DB",
		},
		cli.StringFlag{
			Name:   "data",
			EnvVar: "INPUT_DATA",
		},
		cli.StringFlag{
			Name: "login",
		},
	}
	app.Before = beforeHandler
	app.Action = actionHandler

	app.Commands = []cli.Command{
		cmd.ClientCmd,
		cmd.ServerCmd,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalln(err)
	}

}

func actionHandler(app *cli.Context) {




}

func beforeHandler(app *cli.Context) error {

	dbFile := app.GlobalString("db-file")
	if dbFile == "" {
		return errors.New("please specify db file")
	}

	if err := databases.CheckIfExists(dbFile); err != nil {
		l, err := databases.NewDBClient(&databases.SupermanDatabases{
			LoginDBClient: databases.LoginDBClient{
				DBFile: app.GlobalString("db-file"),
			},
		})

		if err != nil {
			return fmt.Errorf("could not create new DB object due to %v", err)
		}

		defer l.Close()

		if err := l.CreateDB(); err != nil {
			return fmt.Errorf("could not create new table %v", err)
		}
		data, err := os.Open(app.GlobalString("data"))
		if err != nil {
			return fmt.Errorf("error reading file %v", err)
		}
		defer data.Close()

		r := csv.NewReader(data)

		results, err := r.ReadAll()
		if err != nil {
			return fmt.Errorf("could not get results %v", err)
		}


		for _, recs := range results {
			ipaddress := recs[0]
			username := recs[1]
			timestamps := recs[2]

			if err := l.LoadDataset(ipaddress, username, timestamps); err != nil {
				return fmt.Errorf("error occured %v", err)
			}
		}

		fmt.Println("database is loaded")
	}

	return nil
}

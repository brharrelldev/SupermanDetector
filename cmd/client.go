package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/brharrelldev/SupermanDetector/service"
	"github.com/urfave/cli"
	"io/ioutil"
	"log"
	"net/http"
)

var ClientCmd = cli.Command{
	Name: "client",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "login",
			Usage: "login name to search",
		},
		cli.StringFlag{
			Name:  "addr",
			Usage: "address of service",
		},
	},
	Before: func(app *cli.Context) error {
		if app.String("login") == "" || app.String("addr") == "" {
			if err := cli.ShowSubcommandHelp(app); err != nil {
				return fmt.Errorf("error displaying help %v", err)
			}
		}

		return nil

	},
	Action: func(app *cli.Context) {

		var superResponse service.SuperResponse

		loginReq := service.SuperRequest{
			Username: app.String("login"),
		}

		body, err := json.Marshal(loginReq)
		if err != nil {
			log.Fatalf("error could not unmarshall %v", err)
		}



		req, err := http.NewRequest(http.MethodPost, app.String("addr"), bytes.NewBuffer(body))
		if err != nil {
			log.Fatalf("error building a new request %v", err)
		}

		httpClient := http.Client{}

		resp, err := httpClient.Do(req)
		if err != nil {
			log.Fatalf("error getting http response %v", err)
		}

		defer resp.Body.Close()

		r, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("could not read response body %v", err)
		}

		if err := json.Unmarshal(r, &superResponse); err != nil {
			log.Fatalf("could not decode response %v", err)
		}

		fmt.Println(superResponse)

	},
}

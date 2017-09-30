package main

import (
	"os"
	"github.com/urfave/cli"
	"dev.sigpipe.me/dashie/git.txt/cmd"
	"dev.sigpipe.me/dashie/git.txt/setting"
	"github.com/getsentry/raven-go"
	"fmt"
)

const appVersion = "0.2"

func init() {
	setting.AppVer = appVersion
	if os.Getenv("USE_RAVEN") == "true" {
		raven.SetDSN(os.Getenv("RAVEN_DSN"))
		fmt.Printf("Using Raven with DSN: %s\r\n", os.Getenv("RAVEN_DSN"))
	} else {
		fmt.Println("Running without Raven/Sentry support.")
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "git.txt"
	app.Usage = "paste stuff to the interweb with git backend"
	app.Version = appVersion
	app.Commands = []cli.Command{
		cmd.Web,
	}
	app.Flags = append(app.Flags, []cli.Flag{}...)
	app.Run(os.Args)
}
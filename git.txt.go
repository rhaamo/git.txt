package main

import (
	"os"
	"github.com/urfave/cli"
	"dev.sigpipe.me/dashie/git.txt/cmd"
	"dev.sigpipe.me/dashie/git.txt/setting"
)

const APP_VER = "0.1"

func init() {
	setting.AppVer = APP_VER
}

func main() {
	app := cli.NewApp()
	app.Name = "git.txt"
	app.Usage = "paste stuff to the interweb with git backend"
	app.Version = APP_VER
	app.Commands = []cli.Command{
		cmd.Web,
	}
	app.Flags = append(app.Flags, []cli.Flag{}...)
	app.Run(os.Args)
}
package main

import (
	"os"

	"github.com/tmtk75/cli"
)

var flags = []cli.Flag{
	cli.StringFlag{Name: "filter", Value: "describe", Usage: "String to filter operations"},
	cli.BoolFlag{Name: "debug", Usage: "Output to stdout"},
}

func main() {
	app := cli.NewApp()
	app.Name = "genope"
	app.Usage = "Generate Aliyun API operation file"
	app.Version = "0.1.0"
	app.Commands = []cli.Command{
		cli.Command{
			Name: "subcmd",
			Args: "<api-file>",
			Flags: append(flags, []cli.Flag{
				cli.StringFlag{Name: "dir,d", Value: "cli", Usage: "Output root directory"},
			}...),
			Action: func(c *cli.Context) {
				fn, _ := c.ArgFor("api-file")
				genSubcmd(fn, newOptions(c))
			},
		},
		cli.Command{
			Name: "service",
			Action: func(c *cli.Context) {
				genService()
			},
		},
		cli.Command{
			Name:  "ope",
			Usage: "Generate a .go file for operations of a service",
			Args:  "<api-file>",
			Flags: append(flags, []cli.Flag{
				cli.StringFlag{Name: "dir,d", Value: "api", Usage: "Output root directory"},
			}...),
			Action: func(c *cli.Context) {
				fn, _ := c.ArgFor("api-file")
				genOpe(fn, newOptions(c))
			},
		},
	}
	app.Run(os.Args)
}

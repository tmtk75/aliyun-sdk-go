// aly command
package main

import (
	"os"
	"time"

	"github.com/tmtk75/aliyun-sdk-go/cli/ecs"
	//"github.com/tmtk75/aliyun-sdk-go/cli/oss"
	"github.com/tmtk75/aliyun-sdk-go/cli/ram"
	"github.com/tmtk75/aliyun-sdk-go/cli/rds"
	"github.com/tmtk75/aliyun-sdk-go/cli/slb"
	"github.com/tmtk75/cli"
)

func main() {
	//fmt.Printf("%v\n", time.Now().UTC().Format("2006-01-02T15:04:03Z"))
	//os.Exit(1)

	app := cli.NewApp()
	app.Name = "aly"
	app.Usage = "Aliyun comamnd line interface"
	app.Version = "0.0.1"
	app.Author = "tmtk75"
	app.Commands = topCmds
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "region", Value: "cn-hangzhou", Usage: "The region to use", EnvVar: "ALY_DEFAULT_REGION"},
		cli.StringFlag{Name: "start-time", Usage: "Start time in UTC  e.g) 2016-03-05T09:00:00Z (RFC3339)"},
		cli.StringFlag{Name: "end-time", Value: time.Now().UTC().Format("2006-01-02T15:04:03Z"), Usage: "End time in UTC  e.g) 2016-03-06T09:00:00Z (RFC3339)"},
		cli.StringFlag{Name: "duration", Value: "1h", Usage: "Duration until end-time  e.g) 2h45m", EnvVar: "ALY_DEFAULT_DURATION"},
	}
	app.Before = func(c *cli.Context) error {
		return nil
	}
	app.Run(os.Args)
}

var topCmds = []cli.Command{
	cli.Command{
		Name:        "ecs",
		Usage:       "ECS services",
		Subcommands: ecs.Commands,
	},
	cli.Command{
		Name:        "rds",
		Usage:       "RDS services",
		Subcommands: rds.Commands,
	},
	cli.Command{
		Name:        "slb",
		Usage:       "SLB services",
		Subcommands: slb.Commands,
	},
	//cli.Command{
	//	Name:        "oss",
	//	Usage:       "OSS services",
	//	Subcommands: oss.Commands,
	//},
	cli.Command{
		Name:        "ram",
		Usage:       "RAM services",
		Subcommands: ram.Commands,
	},
}

//go:generate gen/gen service
//go:generate gen/gen ope node_modules/aliyun-sdk/apis/ecs-2014-05-26.json
//go:generate gen/gen ope node_modules/aliyun-sdk/apis/rds-2014-08-15.json
//go:generate gen/gen ope node_modules/aliyun-sdk/apis/ots-2014-08-08.json
//go:generate gen/gen ope node_modules/aliyun-sdk/apis/slb-2014-05-15.json
//go:generate gen/gen ope --filter list node_modules/aliyun-sdk/apis/oss-2013-10-15.json
//go:generate gen/gen ope --filter list node_modules/aliyun-sdk/apis/ram-2015-05-01.json
//go:generate gen/gen subcmd node_modules/aliyun-sdk/apis/ecs-2014-05-26.json
//go:generate gen/gen subcmd node_modules/aliyun-sdk/apis/rds-2014-08-15.json
//go:generate gen/gen subcmd node_modules/aliyun-sdk/apis/ots-2014-08-08.json
//go:generate gen/gen subcmd node_modules/aliyun-sdk/apis/slb-2014-05-15.json
//go:generate gen/gen subcmd --filter list node_modules/aliyun-sdk/apis/oss-2013-10-15.json
//go:generate gen/gen subcmd --filter list node_modules/aliyun-sdk/apis/ram-2015-05-01.json

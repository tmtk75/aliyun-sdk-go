package slb

import (
	"github.com/tmtk75/aliyun-sdk-go/cli/helper"
	"github.com/tmtk75/cli"
)

var Commands = helper.Merge(defaultCommands, []cli.Command{
	{
		Name:  "describe-load-balancers",
		Usage: "Describes load balancers",
	},
})

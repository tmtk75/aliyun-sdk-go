package rds

import (
	"fmt"

	"github.com/tmtk75/aliyun-sdk-go/api/rds"
	"github.com/tmtk75/aliyun-sdk-go/cli/helper"
	"github.com/tmtk75/cli"
)

var Commands = helper.Merge(append(defaultCommands,
	[]cli.Command{
		{
			Name: "describe-db-instance-performance-keys",
			Action: func(c *cli.Context) {
				for _, v := range rds.DBInstancePerformanceKeys {
					fmt.Println(v)
				}
			},
		},
	}...),
	[]cli.Command{
		{
			Name:  "describe-db-instances",
			Usage: "Returns information about provisioned RDS instances",
		},
		{
			Name:  "describe-db-instance-performance",
			Usage: "Describes performance information",
		},
		{
			Name:  "describe-db-instance-attribute",
			Usage: "Describes dataase instance attribute",
		},
		{
			Name:  "describe-parameters",
			Usage: "Describe parameters",
		},
		{
			Name:  "describe-resource-usage",
			Usage: "Describe resource usage",
		},
		{
			Name:  "describe-binlog-files",
			Usage: "Describes binlog files",
		},
	})

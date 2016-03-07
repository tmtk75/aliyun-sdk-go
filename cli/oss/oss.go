package oss

import (
	"github.com/tmtk75/aliyun-sdk-go/api/helper"
	"github.com/tmtk75/aliyun-sdk-go/api/oss"
	"github.com/tmtk75/cli"
)

var Commands = helper.Merge(defaultCommands, []cli.Command{
	{
		Name: "list-buckets",
		Action: func(c *cli.Context) {
			oss.New().Request(oss.ListBuckets{})
		},
	},
})

package ecs

import (
	"github.com/tmtk75/aliyun-sdk-go/api"
	"github.com/tmtk75/aliyun-sdk-go/api/ecs"
	"github.com/tmtk75/aliyun-sdk-go/cli/helper"
	"github.com/tmtk75/cli"
)

var Commands = helper.Merge(defaultCommands, []cli.Command{
	{
		Name:  "describe-instances",
		Usage: "Describes one or more of your ECS instances",
	},
	{
		Name:  "describe-eip-addresses",
		Usage: "Describes EIP addresses",
	},
	{
		Name:  "describe-images",
		Usage: "Describe images",
	},
	{
		Name:  "describe-instance-types",
		Usage: "",
	},
	{
		Name:  "describe-instance-vnc-passwd",
		Usage: "",
		Args:  "<instance-id>",
		Action: func(c *cli.Context) {
			conf := &api.Config{RegionId: c.GlobalString("region")}
			id, _ := c.ArgFor("instance-id")
			ecs.New(conf).Request(ecs.DescribeInstanceVncPasswd{
				InstanceId: id,
			})
		},
	},
})

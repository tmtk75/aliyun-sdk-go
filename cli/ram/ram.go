package ram

import (
	"github.com/tmtk75/aliyun-sdk-go/cli/helper"
	"github.com/tmtk75/cli"
)

var Commands = helper.Merge(defaultCommands, []cli.Command{
	{
		Name:  "list-access-keys",
		Usage: "List access keys",
	},
	{
		Name:  "list-groups",
		Usage: "List groups",
	},
	{
		Name:  "list-policies",
		Usage: "List policies",
	},
	{
		Name:  "list-roles",
		Usage: "List roles",
	},
	{
		Name:  "list-users",
		Usage: "List users",
	},
})

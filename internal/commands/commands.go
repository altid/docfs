package commands

import "github.com/altid/libs/service/commander"

var Commands = []*commander.Command{
	{
		Name:			"fetch",
		Description: 	"Retrieve feed for remote URL",
		Args:			[]string{"name"},
		Heading:		commander.DefaultGroup,
	},
}
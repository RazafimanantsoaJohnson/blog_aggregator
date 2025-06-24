package main

import (
	"fmt"
	"os"

	"github.com/RazafimanantsoaJohnson/blog_aggregator/internal/commands"
)

func ReceiveCommandFromCLI() (commands.Command, error) {
	resultCmd := commands.Command{}
	cliArgs := os.Args
	if len(cliArgs) < 2 {
		return resultCmd, fmt.Errorf("the command is invalid")
	}
	resultCmd.Name = cliArgs[1]
	resultCmd.Args = cliArgs[2:]
	return resultCmd, nil
}

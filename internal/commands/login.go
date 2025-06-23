package commands

import (
	//"github.com/RazafimanantsoaJohnson/blog_aggregator/internal/config"
	"fmt"
)

func HandlerLogin(state *State, cmd Command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("this command requires an argument")
	}
	err := state.config.SetUser(cmd.args[0])
	if err != nil {
		return err
	}

	fmt.Printf("The user has been set to: %v\n", cmd.args[0])
	return nil
}

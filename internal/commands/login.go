package commands

import (
	//"github.com/RazafimanantsoaJohnson/blog_aggregator/internal/config"
	"fmt"
)

func HandlerLogin(state *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("this command requires an argument")
	}
	fmt.Printf("Received args for login: %v\n", cmd.Args)
	err := state.Config.SetUser(cmd.Args[0])
	if err != nil {
		return err
	}

	fmt.Printf("The user has been set to: %v\n", cmd.Args[0])
	return nil
}

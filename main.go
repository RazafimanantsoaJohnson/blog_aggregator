package main

import (
	"fmt"
	"os"

	"github.com/RazafimanantsoaJohnson/blog_aggregator/internal/commands"
	"github.com/RazafimanantsoaJohnson/blog_aggregator/internal/config"
)

func main() {
	conf, err := config.Read()
	if err != nil {
		fmt.Errorf(err.Error())
		os.Exit(1)
	}
	state := commands.State{
		Config: &conf,
	}
	cmds := commands.Commands{
		List: make(map[string]func(*commands.State, commands.Command) error),
	}
	cmds.Register("login", commands.HandlerLogin) // the state and the command will be passed when running the handler

	//will turn into the REPL
	receivedCmd, err := ReceiveCommandFromCLI()
	if err != nil {
		fmt.Printf("an error happened: %v\n", err)
		os.Exit(1)
	}
	err = cmds.Run(&state, receivedCmd)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	// cmds.Run(&state, commands.Command{Name: "login", Args: []string{"hello"}})
	// fmt.Println(conf)
	// fmt.Printf("Params: %v", os.Args[1:])
}

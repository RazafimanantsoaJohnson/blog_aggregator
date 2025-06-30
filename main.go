package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/RazafimanantsoaJohnson/blog_aggregator/internal/commands"
	"github.com/RazafimanantsoaJohnson/blog_aggregator/internal/config"
	"github.com/RazafimanantsoaJohnson/blog_aggregator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	conf, err := config.Read()
	if err != nil {
		fmt.Errorf(err.Error())
		os.Exit(1)
	}
	cmds := registerCmds()
	// db connection
	db, err := sql.Open("postgres", conf.DbUrl)
	if err != nil {
		fmt.Errorf(err.Error())
		os.Exit(1)
	}
	dbQueries := database.New(db)
	// initialize the application state
	state := commands.State{
		Config:    &conf,
		DbQueries: dbQueries,
	}
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

func registerCmds() commands.Commands {
	result := commands.Commands{
		List: make(map[string]func(*commands.State, commands.Command) error),
	}
	result.Register("login", commands.HandlerLogin) // the state and the command will be passed when running the handler
	result.Register("register", commands.HandlerRegister)
	result.Register("reset", commands.HandlerReset)
	result.Register("users", commands.HandlerListUsers)
	result.Register("agg", commands.HandlerAggregate)
	result.Register("feeds", commands.HandlerListFeeds)
	result.Register("addfeed", middlewareLoggedIn(commands.HandlerAddFeed))
	result.Register("follow", middlewareLoggedIn(commands.HandlerFollow))
	result.Register("following", middlewareLoggedIn(commands.HandlerFollowing))
	result.Register("unfollow", middlewareLoggedIn(commands.HandlerUnfollow))
	result.Register("browse", middlewareLoggedIn(commands.HandlerBrowse))
	return result
}

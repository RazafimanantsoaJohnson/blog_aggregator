package main

import (
	"context"

	"github.com/RazafimanantsoaJohnson/blog_aggregator/internal/commands"
	"github.com/RazafimanantsoaJohnson/blog_aggregator/internal/database"
)

func middlewareLoggedIn(handler func(s *commands.State, cmd commands.Command, user database.User) error) func(*commands.State, commands.Command) error {
	resultHandler := func(state *commands.State, c commands.Command) error {
		curUser, err := state.DbQueries.GetUser(context.Background(), state.Config.CurUser)
		if err != nil {
			return err
		}
		return handler(state, c, curUser)
	}
	return resultHandler
}

package commands

import (
	//"github.com/RazafimanantsoaJohnson/blog_aggregator/internal/config"
	"context"
	"fmt"
	"time"

	"github.com/RazafimanantsoaJohnson/blog_aggregator/internal/database"
	"github.com/google/uuid"
)

func HandlerLogin(state *State, cmd Command) error {
	if len(cmd.Args) == 0 || len(cmd.Args) > 1 {
		return fmt.Errorf("this command requires one argument")
	}
	receivedUsername := cmd.Args[0]
	_, err := state.DbQueries.GetUser(context.Background(), receivedUsername)
	if err != nil {
		return err
	}
	err = state.Config.SetUser(receivedUsername)
	if err != nil {
		return err
	}

	fmt.Printf("The user has been set to: %v\n", receivedUsername)
	return nil
}

func HandlerRegister(state *State, cmd Command) error {
	if len(cmd.Args) == 0 || len(cmd.Args) > 1 {
		return fmt.Errorf("this command requires one argument")
	}
	createdUser, err := state.DbQueries.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
	})

	if err != nil {
		return err
	}

	err = state.Config.SetUser(cmd.Args[0])
	if err != nil {
		return err
	}
	fmt.Println("The user has been created")
	fmt.Printf("CreatedUser: %v\n", createdUser)

	return nil
}

func HandlerReset(state *State, cmd Command) error {
	err := state.DbQueries.DeleteAllUsers(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("users database has been reset")
	return nil
}

func HandlerListUsers(state *State, cmd Command) error {
	allUsers, err := state.DbQueries.GetAllUsers(context.Background())
	if err != nil {
		return err
	}
	loggedInUser := state.Config.CurUser
	for _, user := range allUsers {
		if user.Name == loggedInUser {
			fmt.Printf("\t* %v (current)\n", user.Name)
			continue
		}
		fmt.Printf("\t* %v\n", user.Name)
	}
	return nil
}

package main

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/RazafimanantsoaJohnson/blog_aggregator/internal/commands"
	"github.com/RazafimanantsoaJohnson/blog_aggregator/internal/database"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

func TestAddCommand(t *testing.T) {
	commandsToRegister := []string{"command1", "command2", "command3"}
	cmds := commands.Commands{
		List: make(map[string]func(*commands.State, commands.Command) error),
	}
	for _, cmd := range commandsToRegister {
		cmds.Register(cmd, func(s *commands.State, c commands.Command) error { return nil })
		if _, ok := cmds.List[cmd]; !ok {
			t.Errorf("the command %v was not registered in the commandList", cmd)
			return
		}
	}
}

func TestRunningCommand(t *testing.T) {
	// testing with an addition command
	var result int
	cmds := commands.Commands{
		List: map[string]func(*commands.State, commands.Command) error{
			"add": func(s *commands.State, cm commands.Command) error {
				a, err := strconv.Atoi(cm.Args[0])
				b, err := strconv.Atoi(cm.Args[1])
				if err != nil {
					return err
				}
				result = a + b
				fmt.Println(a, b, result)
				return nil
			},
		},
	}
	cases := []struct {
		input    []string
		expected int
	}{
		{
			input:    []string{"3", "3"},
			expected: 6,
		},
		{
			input:    []string{"4", "5"},
			expected: 9,
		},
	}

	for _, c := range cases {
		cm := commands.Command{
			Name: "add",
			Args: c.input,
		}
		cmds.Run(nil, cm)
		if result != c.expected {
			t.Errorf("returned value: %v, expected value: %v", result, c.expected)
			return
		}
	}
}

func TestRegister(t *testing.T) {
	cases := []string{"johnson", "anotherguy", "shatyt"}
	// db connection
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/gator")
	if err != nil {
		fmt.Println("error here")
		t.Errorf(err.Error())
		return
	}
	dbQueries := database.New(db)
	for _, c := range cases {
		createdUser, err := dbQueries.CreateUser(context.Background(), database.CreateUserParams{
			ID:        uuid.New(),
			Name:      c,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
		if err != nil {
			t.Errorf(err.Error())
			return
		}
		if createdUser.Name != c {
			t.Errorf("Inserted username: %v; created username: %v", c, createdUser.Name)
			return
		}
	}
}

package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/RazafimanantsoaJohnson/blog_aggregator/internal/database"
	"github.com/google/uuid"
)

func HandlerFollow(s *State, cmd Command, curUser database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("this command require 1 argument:\t url")
	}
	chosenFeed, err := s.DbQueries.GetFeedByURL(context.Background(), cmd.Args[0])
	if err != nil {
		fmt.Println("Arg: " + cmd.Args[0])
		return err
	}

	created_row, err := s.DbQueries.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    curUser.ID,
		FeedID:    chosenFeed.ID,
	})
	if err != nil {
		fmt.Println(created_row)
		return err
	}
	fmt.Printf("%v officially follow the feed: %v\n", curUser.Name, chosenFeed.Name)
	fmt.Println(created_row)
	return nil
}

func HandlerFollowing(s *State, cmd Command, curUser database.User) error {
	allFollows, err := s.DbQueries.GetFeedFollowsForUser(context.Background(), curUser.ID)
	if err != nil {
		return err
	}
	fmt.Println(curUser.Name)
	for _, feed := range allFollows {
		fmt.Printf("\t- '%v'\n", feed.FeedName)
	}
	return nil
}

func HandlerUnfollow(s *State, cmd Command, curUser database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("this command requires 1 argument:\t'feedUrl'")
	}
	err := s.DbQueries.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		UserID: curUser.ID,
		Url:    cmd.Args[0],
	})
	if err != nil {
		return err
	}
	fmt.Printf("User: %v, officially unfollow the %v feed.\n", curUser.Name, cmd.Args[0])
	return nil
}

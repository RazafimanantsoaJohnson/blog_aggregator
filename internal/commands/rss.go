package commands

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"

	"github.com/RazafimanantsoaJohnson/blog_aggregator/internal/database"
	"github.com/google/uuid"
)

func HandlerAggregate(s *State, cmd Command) error {
	return scrapeFeed(s)
}

func HandlerAddFeed(s *State, cmd Command, curUser database.User) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("this command requires 2 arguments: name, url \n")
	}
	createdFeed, err := s.DbQueries.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		UserID:    curUser.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
	})
	if err != nil {
		return err
	}
	_, err = s.DbQueries.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    curUser.ID,
		FeedID:    createdFeed.ID,
	})
	if err != nil {
		return err
	}

	fmt.Println("The feed was created")
	fmt.Println(createdFeed)
	return nil
}

func HandlerListFeeds(s *State, cmd Command) error {
	allFeeds, err := s.DbQueries.GetAllFeeds(context.Background())
	if err != nil {
		return err
	}
	for _, feed := range allFeeds {
		fmt.Printf("%v:\n", feed.Name)
		fmt.Printf("\t%v\n", feed.Url)
		fmt.Printf("\t%v\n", feed.UserName)
	}
	return nil
}

func scrapeFeed(s *State) error {
	nextFeedToFetch, err := s.DbQueries.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}
	fetchedFeed, err := fetchFeed(context.Background(), nextFeedToFetch.Url)
	if err != nil {
		return err
	}
	_, err = s.DbQueries.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
		ID:          nextFeedToFetch.ID,
		UpdatedAt:   time.Now(),
		LastFetched: sql.NullTime{Valid: true, Time: time.Now()},
	})
	if err != nil {
		return err
	}

	fmt.Printf("'%v' posts:\n", fetchedFeed.Channel.Title)
	for _, item := range fetchedFeed.Channel.Item {
		fmt.Printf("\t- '%v'\n", item.Title)
	}
	return nil
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	request, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("User-Agent", "gator")
	client := http.Client{}
	res, err := client.Do(request)
	if err != nil {
		//network error
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var result RSSFeed
	err = xml.Unmarshal(body, &result)
	formatFeed(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func formatFeed(feed *RSSFeed) { // will mutate the feed to get the values in a well formated way
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for i, item := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(item.Title)
		feed.Channel.Item[i].Description = html.UnescapeString(item.Description)
	}
}

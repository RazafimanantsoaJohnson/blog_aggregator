package commands

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/RazafimanantsoaJohnson/blog_aggregator/internal/database"
	"github.com/google/uuid"
)

func HandlerAggregate(s *State, cmd Command) error { // a command/function which might run forever (and it's OK because it is runned on specific intervals)
	// will take a time duration string (ex: 3h4m20s63ms) and print our scraped feed in a loop.
	if len(cmd.Args) != 1 {
		return fmt.Errorf("this command expects an argument:\t'duration'(xhxxmxxs)\n")
	}
	durationBetweenRequests, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return err
	}
	fmt.Println("Collecting feeds every " + cmd.Args[0])
	ticker := time.NewTicker(durationBetweenRequests) //creating the 'recurring event'
	for ; ; <-ticker.C {
		err = scrapeFeed(s)
		if err != nil {
			return err
		}
	}
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

func HandlerBrowse(s *State, cmd Command, curUser database.User) error {
	limit := 2
	if len(cmd.Args) != 0 {
		param, err := strconv.Atoi(cmd.Args[0])
		if err != nil {
			return err
		}
		limit = param
	}
	curUserPosts, err := s.DbQueries.GetUserPosts(context.Background(), database.GetUserPostsParams{UserID: curUser.ID, Limit: int32(limit)})
	if err != nil {
		return err
	}
	fmt.Printf("%v posts: \n", curUser.Name)
	for _, post := range curUserPosts {
		fmt.Printf("\t- %v, %v\n", post.Title, post.PublishedAt)
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

	fmt.Printf("'%v' posts:\n", nextFeedToFetch.Name)
	for _, item := range fetchedFeed.Channel.Item {
		dateLayouts := []string{time.Layout, time.ANSIC, time.UnixDate, time.RubyDate, time.RFC822, time.RFC822Z, time.RFC850, time.RFC1123, time.RFC1123Z, time.RFC3339,
			time.RFC3339Nano, time.Kitchen, time.Stamp, time.StampMicro, time.StampMilli, time.DateTime, time.DateOnly}
		var pubDate time.Time
		for _, layout := range dateLayouts {
			pubDate, err = time.Parse(layout, item.PubDate)
			if err == nil {
				break
			}
		}
		if err != nil {
			fmt.Println("Unable to parse the publication date for post: " + item.PubDate)
		}
		s.DbQueries.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: sql.NullString{Valid: true, String: item.Description},
			PublishedAt: pubDate,
			FeedID:      nextFeedToFetch.ID,
		})
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

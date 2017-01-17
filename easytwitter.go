package bigdatabot

// Twitter related stuff

import (
	"fmt"
	"github.com/dghubble/go-twitter/twitter"
)

// Look in twitter/statuses.go
// are cool methods to do things like looking up statuses with status id's

// This is currently callled with:
// tweetslice := search_twitter_for_keyword(c, messagewithouttwitter)

// Would be cool to call it like this:
// tweetslice , metadata := search_twitter_for_keyword(c, messagewithouttwitter)

func search_twitter_for_keyword(c *twitter.Client, querystring string) ([]twitter.Tweet, error) {

	// Search Tweets
	search, http_response, err := c.Search.Tweets(&twitter.SearchTweetParams{
		Query: querystring,
	})

	if err != nil {
		return nil, err
	}

	//searchquery := search.Metadata.Query
	//refreshurl := search.Metadata.RefreshURL
	//		fmt.Println("Tweet has count of: ", search.Metadata.Count)

	// Read in the slice of statuses.

	tweetslice := search.Statuses
	return tweetslice, nil

}

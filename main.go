package main

import (
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/coreos/pkg/flagutil"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/tucnak/telebot"
	"log"
	"os"
	"strings"
	"time"
)

func messages(bot *telebot.Bot, c *twitter.Client) {
	for message := range bot.Messages {
		log.Printf("Received a message from %s with the text: %s\n",
			message.Sender.Username, message.Text)

		// This works:
		if strings.HasPrefix(message.Text, "!twitter") {

			messagewithouttwitter := strings.TrimPrefix(message.Text, "!twitter")

			//Need to url decode string
			// Need to remove space prefixed to string

			// Dit werk nou add error checking

			tweetslice := search_twitter_for_keyword(c, messagewithouttwitter)

			for i := 0; i < len(tweetslice); i++ {

				// Get the text of each tweet.
				// for other struct attributes look at Tweet struct in twitter/statuses.go ...
				tweet_text := tweetslice[i].Text // The text contained in the tweet.
				//	tweet_user := tweetslice[i].User // user who tweeted something

				//tweet_gps := tweetslice[i].Coordinates
				// If there is GPS co-ordinates: Where The Tweet was made from.
				// If this has no value then it returns nil. Create a check to check if it was nil.
				tweet_creation_time := tweetslice[i].CreatedAt
				// New Stuff to check for:
				tweet_lang := tweetslice[i].Lang
				tweet_id_string := tweetslice[i].IDStr

				/*			stuff := TweetStuff{
							id:             tweetslice[i].IDStr,
							tweetcreatedat: tweetslice[i].CreatedAt,
							tweetlang:      tweetslice[i].Lang,
							tweetsource:    tweetslice[i].Source,
							tweettext:      tweetslice[i].Text,
						}  */

				fmt.Println("Tweet Language is:", tweet_lang)
				fmt.Println("Tweet id string is:", tweet_id_string)

				fmt.Println("Tweet text is: ", tweet_text)
				//fmt.Println("User who tweeted: %s\n", tweet_user)

				//fmt.Println("The Tweet was made using GPS co-ordinates of: %s\n", tweet_gps)
				fmt.Println("Tweet was created at: ", tweet_creation_time)

				// Log to dynamodb

				// Send messages with telegram

				bot.SendMessage(message.Chat,
					"You wanted something from twitter so here goes...", nil)
				bot.SendMessage(message.Chat,
					tweet_text, nil)
				bot.SendMessage(message.Chat,
					tweet_creation_time, nil)

				sess, err := session.NewSessionWithOptions(session.Options{
					Config:  aws.Config{Region: aws.String("eu-west-1")},
					Profile: "dynamodb-eclipse",
				})

				if err != nil {
					fmt.Println("failed to create session,", err)
					return
				}

				svc := dynamodb.New(sess)

				r := Record{
					ID:                     tweetslice[i].IDStr,
					TweetCreatedat:         tweetslice[i].CreatedAt,
					TweetLang:              tweetslice[i].Lang,
					TweetText:              tweetslice[i].Text,
					TweetQuotedStatusIDStr: tweetslice[i].QuotedStatusIDStr,
					TweetRetweetCount:      tweetslice[i].RetweetCount,
					TweetFavoriteCount:     tweetslice[i].FavoriteCount,
					TweetPossiblySensitive: tweetslice[i].PossiblySensitive,
				}

				item, err := dynamodbattribute.MarshalMap(r)
				if err != nil {
					fmt.Println("Failed to convert", err)
					return
				}

				result, err := svc.PutItem(&dynamodb.PutItemInput{
					Item:      item,
					TableName: aws.String("twitter3"),
				})

				fmt.Println(result)

			}

		}
	}
}

func queries(bot *telebot.Bot) {
	for query := range bot.Queries {
		log.Println("--- new query ---")
		// Could write some code to log this to sqlite?
		log.Println("from:", query.From.Username)
		// And this
		log.Println("text:", query.Text)

		// Figure uit wat doen hierdie deel en wat is die use case.

		// Create an article (a link) object to show in results.
		article := &telebot.InlineQueryResultArticle{
			Title: "Telebot",
			URL:   "https://github.com/tucnak/telebot",
			InputMessageContent: &telebot.InputTextMessageContent{
				Text:           "Telebot is a Telegram bot framework.",
				DisablePreview: false,
			},
		}

		// Build the list of results (make sure to pass pointers!).
		results := []telebot.InlineQueryResult{article}

		// Build a response object to answer the query.

		// Dis hoe mens 'n struct van imported package call: Baie goe voorbeeld
		response := telebot.QueryResponse{
			Results:    results,
			IsPersonal: true,
		}

		// Send it.
		if err := bot.AnswerInlineQuery(&query, &response); err != nil {
			log.Println("Failed to respond to query:", err)
		}
	}
}

func search_twitter_for_keyword(c *twitter.Client, querystring string) []twitter.Tweet {

	// Search Tweets
	search, _, _ := c.Search.Tweets(&twitter.SearchTweetParams{
		Query: querystring,
	})
	// Figure out how to print out better as it is url encoded at the moment
	fmt.Println("Tweets matching our query: ", search.Metadata.Query)

	// Convert to string before printing or use print format character for int.
	fmt.Println("Tweet has count of: ", search.Metadata.Count)

	// We can access stuff from results like this:

	//searchquery := search.Metadata.Query
	//refreshurl := search.Metadata.RefreshURL
	//	fmt.Println("Tweet has the following refreshurl: %s\n", refreshurl)

	// Read in the slice of statuses.

	tweetslice := search.Statuses
	return tweetslice

}

type Record struct {
	ID                     string
	TweetCreatedat         string
	TweetLang              string
	TweetText              string
	TweetQuotedStatusIDStr string
	TweetRetweetCount      int
	TweetFavoriteCount     int
	TweetPossiblySensitive bool
}

func main() {

	// Twitter client related code
	flags := flag.NewFlagSet("user-auth", flag.ExitOnError)
	consumerKey := flags.String("consumer-key", "", "Twitter Consumer Key")
	consumerSecret := flags.String("consumer-secret", "", "Twitter Consumer Secret")
	accessToken := flags.String("access-token", "", "Twitter Access Token")
	accessSecret := flags.String("access-secret", "", "Twitter Access Secret")
	flags.Parse(os.Args[1:])
	flagutil.SetFlagsFromEnv(flags, "TWITTER")

	if *consumerKey == "" || *consumerSecret == "" || *accessToken == "" || *accessSecret == "" {
		log.Fatal("Consumer key/secret and Access token/secret required")
	}

	config := oauth1.NewConfig(*consumerKey, *consumerSecret)
	token := oauth1.NewToken(*accessToken, *accessSecret)
	// OAuth1 http.Client will automatically authorize Requests
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter client
	client := twitter.NewClient(httpClient)

	// Telegram API client related code

	bot, err := telebot.NewBot(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatalln(err)
	}

	bot.Messages = make(chan telebot.Message, 100)
	bot.Queries = make(chan telebot.Query, 1000)

	go messages(bot, client)
	go queries(bot)

	bot.Start(1 * time.Second)
}

package main

import (
	"github.com/tucnak/telebot"
	"log"
	"os"
	"time"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"fmt"
)



type MessageLog struct {

	Username              string
	Messagetext             string
	Title							string
	Type						string
	Id     						int
}

func logmessages(message telebot.Message,svc *dynamodb.DynamoDB) {

	chattype := fmt.Sprintf("%s", message.Chat.Type)

loggedmessage := MessageLog{

	Username:              message.Chat.Username,
	Messagetext:              message.Text,
	Title: 									message.Chat.Title,
	Type:                     chattype,
	Id:											message.ID,

}



// Break this into functions too if you can
item, err := dynamodbattribute.MarshalMap(loggedmessage)
if err != nil {
	fmt.Println("Failed to convert struct into Marshalled Map for DynamoDB to Ingest", err)
	log.Fatal(err)
}

result, err := svc.PutItem(&dynamodb.PutItemInput{
	Item:      item,
	TableName: aws.String("GroupMessageLog"),
})
// up to here
fmt.Println(result)


}



/*
func messages(bot *telebot.Bot, c *twitter.Client, svc *dynamodb.DynamoDB) {
	for message := range bot.Messages {
		log.Printf("Received a message from %s with the text: %s\n",
			message.Sender.Username, message.Text)

		// This works:
		if strings.HasPrefix(message.Text, "!twitter") {

			messagewithouttwitter := strings.TrimPrefix(message.Text, "!twitter")

			//Need to url decode string
			// Need to remove space prefixed to string

			// Dit werk nou add error checking

			tweetslice, err := bigdatabot.SearchTwitterKeyword(c, messagewithouttwitter)

			if err != nil {
				log.Fatal(err)
			}

			for i := 0; i < len(tweetslice); i++ {

				// Get the text of each tweet.
				// for other struct attributes look at Tweet struct in twitter/statuses.go ...
				tweet_text := tweetslice[i].Text

				tweet_creation_time := tweetslice[i].CreatedAt
				// New Stuff to check for:
				tweet_lang := tweetslice[i].Lang
				tweet_id_string := tweetslice[i].IDStr

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

// Break this into functions too if you can
				item, err := dynamodbattribute.MarshalMap(r)
				if err != nil {
					fmt.Println("Failed to convert", err)
					return
				}

				result, err := svc.PutItem(&dynamodb.PutItemInput{
					Item:      item,
					TableName: aws.String("twitter3"),
				})
// up to here
				fmt.Println(result)

			}

		}

	}
}
*/

/*

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

*/

func main() {

	/*

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

	*/

	// DynamoDB Stuff
	// Commment out from here

		sess, err := session.NewSessionWithOptions(session.Options{
			Config:  aws.Config{Region: aws.String("eu-west-1")},
			Profile: "dynamodb-eclipse",
		})

		if err != nil {
			fmt.Println("failed to create session,", err)
			return
		}

		svc := dynamodb.New(sess)

		// DynamoDB Stuff ends here



	// Up to here ..............................................

	// Telegram API client related code

	bot, err := telebot.NewBot(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatalln(err)
	}

	//bot.Messages = make(chan telebot.Message, 100)
	//bot.Queries = make(chan telebot.Query, 1000)

	// Also passing twitter client to messages in go routine

	//go messages(bot, client, svc)
	//go queries(bot)

	messages := make(chan telebot.Message, 100)
	bot.Listen(messages, 1*time.Second)
	//bot.Start(1 * time.Second)
	for message := range messages {

		// See if we can call it here.
		// If this works we can make the function return an error and add propper error checking
		logmessages(message,svc)

		log.Printf("Received a message : %s\n", message.Text)
		log.Printf("Message type : %s\n", message.Chat.Type)
		log.Printf("Message type : %s\n", message.Chat.Title)
		log.Printf("Message ID: %s\n", message.ID)

	}

}

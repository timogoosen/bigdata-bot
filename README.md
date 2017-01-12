# bigdata-bot
Telegram bot that does various things...

Still crashes under lots of circumstances. Needs lots of work still.

README:

Install dependencies with:

This is to interact with telegram REST API with a Golang wrapper:
<pre>
$ go get github.com/tucnak/telebot
</pre>
Install AWS SDK With:
<pre>
$ go get github.com/aws/aws-sdk-go/aws
</pre>
Install FlagUtil from CoreOS:
<pre>
$ go get github.com/coreos/pkg/flagutil
</pre>
Install Go OAth Module:
<pre>
$ go get github.com/dghubble/oauth1
</pre>

Also for the stuff that interacts with DynamoDB make sure
to have your credentials in :
<pre>

$ vi ~/.aws/config

$ vi ~/.aws/credentials

</pre>


Install Twitter API Wrapper:

<pre>
$ go get github.com/dghubble/go-twitter/twitter
</pre>


Set env variables:
<pre>

export TWITTER_CONSUMER_KEY=aaaaaaa

export TWITTER_CONSUMER_SECRET=aaaaa

export TWITTER_ACCESS_TOKEN=aaaaa

export TWITTER_ACCESS_SECRET=sssss

</pre>




Export your TELEGRAM API KEY with:


<pre>
$ export BOT_TOKEN=aaaaaaa

</pre>

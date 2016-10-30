# Twittersounds

Tweet links to random books retrieved from Project Gutenburg. Tweets are scheduled
by specifying a time of day and the number of days in between tweets. Tweets are
sent out on the sheduled time with a variable delay.

To run, after installing with

```
go get github.com/aaronnorby/twittersounds
```
do the following:

```
import (
  "github.com/aaronnorby/twittersounds"
  "os"
)


os.Setenv("CONSUMER_KEY", "consumer_key")
os.Setenv("CONSUMER_SECRET", "consumer_secret")
os.Setenv("ACCESS_TOKEN", "access_token")
os.Setenv("ACCESS_TOKEN_SECRET", "access_token_secret")

twittersounds.Schedule(13, 2)
```
This will schedule a tweet for 1pm every two days. You must replace the
placeholders above with your own consumer key, etc. that you get from the Twitter
developer console.

`Schedule` is the entirety of the API.

It can also be run by putting the above code into a subdirectory called `tweetout`
and putting it in the `main` package and running the call to `Schedule` from inside
`func main()`. Then you can build a Docker image from the docker file, which
installs that as a binary executable and runs it (it also installs the dependencies
for you).

// Package twittersounds tweets some stuff from a chosen text to a designated
// twitter account once in a while.
package twittersounds

import (
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"os"
	"time"
)

func Tweet(text string) error {

	var (
		CONSUMER_KEY        = os.Getenv("CONSUMER_KEY")
		CONSUMER_SECRET     = os.Getenv("CONSUMER_SECRET")
		ACCESS_TOKEN        = os.Getenv("ACCESS_TOKEN")
		ACCESS_TOKEN_SECRET = os.Getenv("ACCESS_TOKEN_SECRET")
	)
	anaconda.SetConsumerKey(CONSUMER_KEY)
	anaconda.SetConsumerSecret(CONSUMER_SECRET)
	api := anaconda.NewTwitterApi(ACCESS_TOKEN, ACCESS_TOKEN_SECRET)

	tweet, err := api.PostTweet(text, nil)

	fmt.Printf("%+v", tweet)

	if err != nil {
		return err
	}

	return nil
}

func generateText() string {
	return "string"
}

func Initiate(timeToTweet time.Time, vary bool) {

}

// Package twittersounds tweets some stuff from a chosen text to a designated
// twitter account once in a while.
package twittersounds

import (
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"log"
	"math"
	"math/rand"
	"os"
	"time"
)

const gutenburgBaseUrl string = "https://www.gutenberg.org"

func Tweet(text string, hashtags []string) error {

	var (
		CONSUMER_KEY        = os.Getenv("CONSUMER_KEY")
		CONSUMER_SECRET     = os.Getenv("CONSUMER_SECRET")
		ACCESS_TOKEN        = os.Getenv("ACCESS_TOKEN")
		ACCESS_TOKEN_SECRET = os.Getenv("ACCESS_TOKEN_SECRET")
	)
	anaconda.SetConsumerKey(CONSUMER_KEY)
	anaconda.SetConsumerSecret(CONSUMER_SECRET)
	api := anaconda.NewTwitterApi(ACCESS_TOKEN, ACCESS_TOKEN_SECRET)

	for _, tag := range hashtags {
		text = text + " " + tag
	}
	tweet, err := api.PostTweet(text, nil)

	fmt.Printf("tweet: %+v", tweet)

	if err != nil {
		return err
	}

	return nil
}

func generateText() string {
	book, err := FindBook()
	if err != nil {
		log.Printf("Error getting book data: %v\n", err)
		return ""
	}

	title := book.Title
	author := book.Subtitle
	link := gutenburgBaseUrl + book.Href

	return "\"" + title + "\"," + author + " " + link
}

func Initiate(vary bool) interface{} {
	var delayMins int = 40
	var offsetMins int = 0

	if vary == true {
		seed := time.Now().UnixNano()
		rng := rand.New(rand.NewSource(seed))
		offsetMins = int(math.Floor(rng.NormFloat64()*float64(15) + 0.5))

		// keep it within +/-40mins of scheduled time
		if offsetMins < -40 {
			offsetMins = -40
		} else if offsetMins > 40 {
			offsetMins = 40
		}
	}

	delayMins = delayMins + offsetMins

	tweetText := generateText()
	// if tweetText is empty, there was an error so we'll wait until next time
	if tweetText == "" {
		log.Println("No tweet tweeted")
		return nil
	}

	Tweet(tweetText, []string{"#projectgutenberg", "#randombook"})

	return nil
}

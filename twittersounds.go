// Package twittersounds tweets some stuff from a chosen text to a designated
// twitter account once in a while.
package twittersounds

import (
	"github.com/ChimeraCoder/anaconda"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
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

	log.Printf("tweeted: %+v", tweet)

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

func Initiate(immediate, vary bool) interface{} {
	if immediate && vary {
		return "Error: cannot have immediate execution and variable delay"
	}

	var delayMins int = 40
	var offsetMins int = 0

	if immediate == true {
		delayMins = 0
	}

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

	delay, _ := time.ParseDuration(strconv.Itoa(delayMins) + "m")
	log.Printf("Tweet delayed for %v mins from now.", delayMins)
	timer := time.NewTimer(delay)
	<-timer.C
	timer.Stop()

	err := Tweet(tweetText, []string{"#ProjectGutenberg", "#RandomBook"})
	if err != nil {
		return err
	}

	return nil
}

// timeToTweet is hour of day to send tweet, interval [0,23]. Assumed to be "America/Los_Angeles"
// timezone
func Schedule(timeToTweet, daysBetweenTweets int) {
	tz, _ := time.LoadLocation("America/Los_Angeles")
	hoursBetweenTweets, _ := time.ParseDuration(strconv.Itoa(daysBetweenTweets*24) + "h")

	timeTilNextTweet := getTimeTilNextTweet(timeToTweet, hoursBetweenTweets, tz)
	timer := time.NewTimer(timeTilNextTweet)

	log.Printf("Scheduling tweet for %v from now", timeTilNextTweet)

	maxRetries := 3
	for {
		<-timer.C
		err := Initiate(false, true)
		if err != nil {
			log.Printf("Error in Initiate: %v", err)

			for retries := 0; retries < maxRetries; retries++ {
				log.Println("Retying Initiate")
				err := Initiate(true, false)
				if err == nil {
					break
				}

				if err != nil && retries == maxRetries-1 {
					log.Printf("Max retries exceeded. Stopping with error: %v", err)
					return
				}
			}
		}
		timeTilNextTweet = getTimeTilNextTweet(timeToTweet, hoursBetweenTweets, tz)
		log.Printf("Tweet scheduled for %v from now", timeTilNextTweet)

		timer.Reset(timeTilNextTweet)
	}
}

type Scheduler struct {
	TweetTime int
	Hours     time.Duration
	Tz        *time.Location
	Timer     time.Timer
	C         chan time.Time
}

func getTimeTilNextTweet(timeToTweet int, hoursBetweenTweets time.Duration, tz *time.Location) time.Duration {
	var timeTilNextTweet time.Duration
	nextTweetTime := timeTodayFromHour(timeToTweet, tz)

	if time.Now().After(nextTweetTime) {
		nextTweetTime = nextTweetTime.Add(hoursBetweenTweets)
	}

	timeTilNextTweet = nextTweetTime.Sub(time.Now())

	return timeTilNextTweet
}

func timeTodayFromHour(hour int, zone *time.Location) time.Time {
	return time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), hour, 0, 0, 0, zone)
}

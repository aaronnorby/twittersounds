package twittersounds

import (
	"golang.org/x/net/html"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestGutenberg(t *testing.T) {
	file, err := os.Open("./testdata/gutenberg_random.html")
	defer file.Close()
	if err != nil {
		t.Errorf("Problem opening test data: %v", err)
	}

	doc, err := html.Parse(file)

	if err != nil {
		t.Errorf("Problem parsing test html: %v", err)
	}

	books := parseBookHtml(doc)
	// t.Log(books)

	var expectedResultsLength int = 25

	if expectedResultsLength != len(books) {
		t.Errorf("Error in parseBookHtml. Number of books expected to be %v, was %v", expectedResultsLength, len(books))
	}
}

func TestGetTimeTilNextTweet(t *testing.T) {
	tz, _ := time.LoadLocation("America/Los_Angeles")
	hoursBetweenTweets, _ := time.ParseDuration(strconv.Itoa(2*24) + "h")

	timeTil := getTimeTilNextTweet(2, hoursBetweenTweets, tz)

	if timeTil.Hours() > hoursBetweenTweets.Hours() {
		t.Errorf("getTimeTilNextTweet should have returned value <= %v, returned %v", "48h", timeTil)
	}

	// t.Log(timeTil)

}

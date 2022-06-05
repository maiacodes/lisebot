package main

import (
	"bufio"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"strings"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.Info("Starting LiseBot downloader!")

	// Setup colly collector
	c := colly.NewCollector()
	c.UserAgent = "LiseBot"

	proxyURL := os.Getenv("PROXY")
	if proxyURL != "" {
		logrus.Debug("Using proxy")
		err := c.SetProxy(proxyURL)
		if err != nil {
			logrus.Fatal(err)
		}
	}

	// Open Tweets file
	file, err := os.Open("tweets.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	// Scan each file line
	for scanner.Scan() {
		t := scanner.Text()

		// Extract Tweet data
		text, name, username, id, mediaURLs := extractTweet(t, c)

		// Write file with Tweet info
		err := os.WriteFile(fmt.Sprintf("tweets/%v-%v-info.txt", username, id), []byte(fmt.Sprintf("%v: %v", name, text)), 0644)
		if err != nil {
			logrus.Error("Couldn't write info file for ", t, " - ", err)
		}

		// Download images
		for i, image := range mediaURLs {
			downloadImage(image, fmt.Sprintf("tweets/%v-%v-%v.jpg", username, id, i), c)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

// extractTweet pulls down the metadata of the Tweet
func extractTweet(url string, c *colly.Collector) (text string, name string, username string, id string, mediaURLs []string) {
	// Clone collector
	c = c.Clone()

	// Setup request logger
	c.OnRequest(func(r *colly.Request) {
		logrus.Info("Visiting ", r.URL)
	})

	// Meta tag scanner
	c.OnHTML("meta", func(e *colly.HTMLElement) {
		property := e.Attr("property")
		content := e.Attr("content")

		switch property {
		case "og:description":
			text = content
		case "og:image":
			mediaURLs = append(mediaURLs, content)
		case "og:title":
			name = content
		}
	})

	// Find @ username in Tweet
	c.OnHTML("b", func(e *colly.HTMLElement) {
		if e.Attr("class") == "u-linkComplex-target" {
			username = e.Text
		}
	})

	// Visit Tweet
	err := c.Visit(url)
	if err != nil {
		logrus.Error("Cannot pull tweet ", err)
	}

	// Calculate Tweet ID from URL
	id = strings.TrimPrefix(strings.ToLower(url), "https://twitter.com/"+strings.ToLower(username)+"/status/")

	return
}

// downloadImage pulls a media file from Twitter to the local disk
func downloadImage(url string, fileName string, c *colly.Collector) {
	// Clone collector
	c = c.Clone()

	// Setup request logger
	c.OnRequest(func(r *colly.Request) {
		logrus.Info("Downloading ", r.URL)
	})

	// On response, save image to disk
	c.OnResponse(func(r *colly.Response) {
		err := r.Save(fileName)
		if err != nil {
			logrus.Error("Cannot save image ", err)
		}
	})

	// Visit image URL
	err := c.Visit(url)
	if err != nil {
		logrus.Error("Cannot pull image ", err)
	}
}

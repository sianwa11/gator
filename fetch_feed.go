package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
)

type RSSFeed struct {
	Channel struct {
		Title 			string        `xml:"title"`
		Link 				string        `xml:"link"`
		Description string        `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {

	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("error processing request")
	}

	req.Header.Set("User-Agent", "gator")
	
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("error sending request: %w", err)
	}
	defer res.Body.Close()

	
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("error reading body: %w", err)
	}
	
	var feed RSSFeed
	err = xml.Unmarshal(data, &feed)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("failed to parse data: %w", err)
	}
	
	html.UnescapeString(feed.Channel.Title)
	html.UnescapeString(feed.Channel.Description)

	for _, item := range feed.Channel.Item {
		html.UnescapeString(item.Title)
		html.UnescapeString(item.Description)
	}
	
	return &feed, nil
}
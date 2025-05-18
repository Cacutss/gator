package rss

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       *string `xml:"title"`
	Link        *string `xml:"link"`
	Description *string `xml:"description"`
	PubDate     *string `xml:"pubDate"`
}

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("Error getting req:%w", err)
	}
	req.Header.Set("User-Agent", "gator")
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	defer res.Body.Close()
	feed := RSSFeed{}
	data, err := io.ReadAll(res.Body)
	if err = xml.Unmarshal(data, &feed); err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	for _, v := range feed.Channel.Item {
		*v.Description = html.UnescapeString(*v.Description)
		*v.Title = html.UnescapeString(*v.Title)
	}
	return &feed, nil
}

func ConvertDate(date *string) (time.Time, error) {
	if date == nil {
		return time.Time{}, fmt.Errorf("Date is missing")
	}
	formats := []string{
		time.RFC822,
		time.RFC822Z,
		time.RFC1123,
		time.RFC1123Z,
		time.UnixDate,
		time.ANSIC,
		time.RFC850,
		time.RFC3339,
		time.RFC3339Nano,
	}
	var err error
	result := time.Time{}
	for _, format := range formats {
		result, err := time.Parse(format, *date)
		if err == nil {
			return result, nil
		}
	}
	return result, fmt.Errorf("Could not parse date %q: %w", *date, err)
}

package rss

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoad(t *testing.T) {
	testcases := []struct {
		name        string
		title       string
		description string
		image       string
		link        string
	}{
		{
			name:        "checkly",
			title:       "The Checkly Blog",
			description: "The Checkly Blog is your go-to place for technical stories on building a SaaS, building a company and growing it from scratch.",
			image:       "https://blog.checklyhq.com/favicon.png",
			link:        "https://blog.checklyhq.com/",
		}, {
			name:        "hackernews",
			title:       "Hacker News",
			description: "Links for the intellectually curious, ranked by readers.",
			image:       "",
			link:        "https://news.ycombinator.com/",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			file := filepath.Join("testdata", tc.name+".xml")
			f, err := os.Open(file)
			if err != nil {
				t.Fatalf("could not open test file: %v", err)
			}
			defer f.Close()

			feed, err := Load(f)
			if err != nil {
				t.Fatalf("error loading feed: %v", err)
			}
			{
				got := strings.TrimSpace(feed.Channel.Title)
				if got != tc.title {
					t.Errorf("expected title %q, got %q", tc.title, got)
				}
			}
			{
				got := strings.TrimSpace(feed.Channel.Description)
				if got != tc.description {
					t.Errorf("expected description %q, got %q", tc.description, got)
				}
			}
			{
				got := feed.Channel.Image
				if got != tc.image {
					t.Errorf("expected image %q, got %q", tc.image, got)
				}
			}
			{
				got := feed.Channel.Link
				if got != tc.link {
					t.Errorf("expected link %q, got %q", tc.link, got)
				}
			}
		})
	}
}

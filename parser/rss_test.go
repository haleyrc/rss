package parser

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestProcessElement(t *testing.T) {
	mod := func(_ string) string {
		return "the replacement text"
	}
	testcases := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "missing end tag",
			input: `preamble <content:encoded>some stuff`,
			want:  "",
		},
		{
			name:  "no tag pairs",
			input: `some bare text`,
			want:  `some bare text`,
		},
		{
			name:  "single tag pair",
			input: `some preamble text <content:encoded>some stuff that needs quoting</content:encoded> some postamble text`,
			want:  `some preamble text <content:encoded>the replacement text</content:encoded> some postamble text`,
		},
		{
			name:  "multiple tag pairs",
			input: `some preamble text <content:encoded>some stuff that needs quoting</content:encoded> some <content:encoded>some stuff that needs quoting</content:encoded> postamble text`,
			want:  `some preamble text <content:encoded>the replacement text</content:encoded> some <content:encoded>the replacement text</content:encoded> postamble text`,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := ProcessElementText(tc.input, TagStartContentEncoded, TagEndContentEncoded, mod)
			if got != tc.want {
				t.Errorf("wanted %q, got %q", tc.want, got)
			}
		})
	}

}

func TestQuote(t *testing.T) {
	input, err := ioutil.ReadFile(filepath.Join("testdata", "checklymin.xml"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := Quote(string(input), TagStartContentEncoded, TagEndContentEncoded)
	feed, err := Load(strings.NewReader(output))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "https://blog.checklyhq.com/"
	if feed.Channel.Link != want {
		t.Errorf("wanted %q, got %q", want, feed.Channel.Link)
	}
}

func TestLoad(t *testing.T) {
	testcases := []struct {
		name        string
		title       string
		description string
		image       string
		link        string
		err         bool
	}{
		{
			name: "invalid",
			err:  true,
		},
		{
			name:        "checklymin",
			title:       "The Checkly Blog",
			description: "The Checkly Blog is your go-to place for technical stories on building a SaaS, building a company and growing it from scratch.",
			image:       "https://blog.checklyhq.com/favicon.png",
			link:        "https://blog.checklyhq.com/",
		},
		{
			name:        "checkly",
			title:       "The Checkly Blog",
			description: "The Checkly Blog is your go-to place for technical stories on building a SaaS, building a company and growing it from scratch.",
			image:       "https://blog.checklyhq.com/favicon.png",
			link:        "https://blog.checklyhq.com/",
		},
		{
			name:        "hackernews",
			title:       "Hacker News",
			description: "Links for the intellectually curious, ranked by readers.",
			image:       "",
			link:        "https://news.ycombinator.com/",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			file := filepath.Join("..", "testdata", tc.name+".xml")
			f, err := os.Open(file)
			if err != nil {
				t.Fatalf("could not open test file: %v", err)
			}
			defer f.Close()

			feed, err := Load(f)
			if tc.err {
				if err == nil {
					t.Fatalf("expected error, but got none")
				}
				return
			}
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

func TestRepositoryURL(t *testing.T) {
	testcases := []struct {
		name        string
		title       string
		description string
		image       string
		link        string
		url         string
		err         bool
	}{
		{
			name:        "checkly",
			title:       "The Checkly Blog",
			description: "The Checkly Blog is your go-to place for technical stories on building a SaaS, building a company and growing it from scratch.",
			image:       "https://blog.checklyhq.com/favicon.png",
			link:        "https://blog.checklyhq.com/",
			url:         "https://blog.checklyhq.com/rss/",
		},
		{
			name:        "hackernews",
			title:       "Hacker News",
			description: "Links for the intellectually curious, ranked by readers.",
			image:       "",
			link:        "https://news.ycombinator.com/",
			url:         "https://news.ycombinator.com/rss",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			feed, err := LoadURL(tc.url)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
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

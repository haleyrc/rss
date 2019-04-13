package transport_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"os"
	"strings"
	"testing"

	_ "github.com/lib/pq"

	"github.com/haleyrc/rss"
	"github.com/haleyrc/rss/mock"
	"github.com/haleyrc/rss/transport"
)

var (
	repo       rss.Repository
	controller transport.Controller
)

func TestMain(m *testing.M) {
	repo = mock.NewRepository()
	code := m.Run()
	os.Exit(code)
}

func TestCreateFeed(t *testing.T) {
	srv := transport.NewServer(repo)
	server := httptest.NewServer(srv)
	defer server.Close()

	createResponse, err := http.Post(server.URL+"/feeds", "application/json", strings.NewReader(`{"url":"https://blog.checklyhq.com/rss"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var resp transport.CreateFeedResponse
	if err := json.NewDecoder(createResponse.Body).Decode(&resp); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	url := fmt.Sprintf("%s/feeds/%d", server.URL, resp.Feed.ID)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	deleteResponse, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	b, _ := httputil.DumpResponse(deleteResponse, true)
	fmt.Println(string(b))
}

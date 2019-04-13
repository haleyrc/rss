package transport_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"strings"
	"testing"

	_ "github.com/lib/pq"

	"github.com/haleyrc/rss/repository"
	"github.com/haleyrc/rss/transport"
	"github.com/jmoiron/sqlx"
)

func TestCreateFeed(t *testing.T) {
	db := sqlx.MustConnect("postgres", "host=localhost user=postgres password=test port=5433 dbname=rss sslmode=disable")
	repo := repository.New(db)
	controller := transport.NewController(repo)
	endpoint := transport.NewEndpoint(
		controller.CreateFeed,
		transport.DecodeCreateFeedRequest,
		transport.EncodeResponse,
	)
	server := httptest.NewServer(endpoint)
	defer server.Close()

	resp, err := http.Post(server.URL, "application/json", strings.NewReader(`{"url":"https://blog.checklyhq.com/rss"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	//
	b, _ := httputil.DumpResponse(resp, true)
	fmt.Println(string(b))
}

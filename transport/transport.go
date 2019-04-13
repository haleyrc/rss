package transport

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/haleyrc/rss"
	"github.com/haleyrc/rss/parser"
)

func NewServer(repo rss.Repository) http.Handler {
	controller := NewController(repo)

	createFeedEndpoint := NewEndpoint(
		controller.CreateFeed,
		decodeCreateFeedRequest,
		encodeResponse,
	)

	removeFeedEndpoint := NewEndpoint(
		controller.RemoveFeed,
		decodeRemoveFeedRequest,
		encodeResponse,
	)

	r := mux.NewRouter()
	r.Handle("/feeds", createFeedEndpoint).Methods(http.MethodPost)
	r.Handle("/feeds/{id}", removeFeedEndpoint).Methods(http.MethodDelete)

	return r
}

type Error struct {
	Message string `json:"message"`
}

func encodeResponse(w http.ResponseWriter, data interface{}, err error) {
	if err != nil {
		json.NewEncoder(w).Encode(map[string]Error{
			"error": Error{
				Message: err.Error(),
			},
		})
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": data,
	})
}

func decodeCreateFeedRequest(r *http.Request) (interface{}, error) {
	var request createFeedRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

type Handler func(req interface{}) (interface{}, error)
type DecoderFunc func(r *http.Request) (interface{}, error)
type EncoderFunc func(w http.ResponseWriter, data interface{}, err error)

func NewEndpoint(h Handler, dec DecoderFunc, enc EncoderFunc) Endpoint {
	return Endpoint{
		h:   h,
		dec: dec,
		enc: enc,
	}
}

type Endpoint struct {
	h   Handler
	dec DecoderFunc
	enc EncoderFunc
}

func (e Endpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req, err := e.dec(r)
	if err != nil {
		e.enc(w, nil, err)
		return
	}
	data, err := e.h(req)
	e.enc(w, data, err)
}

func NewController(repo rss.Repository) Controller {
	return Controller{repo}
}

type Controller struct {
	repository rss.Repository
}

type createFeedRequest struct {
	URL string `json:"url"`
}

type CreateFeedResponse struct {
	Feed *rss.Feed `json:"feed"`
}

func (h *Controller) CreateFeed(request interface{}) (interface{}, error) {
	req := request.(createFeedRequest)
	xmlFeed, err := parser.LoadURL(req.URL)
	if err != nil {
		return CreateFeedResponse{}, err
	}

	feed, err := rss.NewFromChannel(xmlFeed.Channel)
	if err != nil {
		return CreateFeedResponse{}, err
	}

	if err := h.repository.CreateFeed(feed, feed.Items...); err != nil {
		return CreateFeedResponse{}, err
	}

	return CreateFeedResponse{Feed: feed}, nil
}

type removeFeedRequest struct {
	ID int64 `json:"id"`
}

type removeFeedResponse struct {
	Status string `json:"status"`
}

func decodeRemoveFeedRequest(r *http.Request) (interface{}, error) {
	var request removeFeedRequest
	id := mux.Vars(r)["id"]
	iid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, err
	}
	request.ID = iid
	return request, nil
}

func (c *Controller) RemoveFeed(request interface{}) (interface{}, error) {
	req := request.(removeFeedRequest)

	if err := c.repository.RemoveFeed(req.ID); err != nil {
		return removeFeedResponse{}, err
	}

	return removeFeedResponse{Status: "success"}, nil
}

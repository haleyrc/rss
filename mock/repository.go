package mock

import (
	"errors"
	"log"

	"github.com/haleyrc/rss"
)

func NewRepository() rss.Repository {
	return &repository{
		feeds: make(map[int64]*rss.Feed),
		items: make(map[int64]*rss.Item),
	}
}

type repository struct {
	lastID int64
	feeds  map[int64]*rss.Feed
	items  map[int64]*rss.Item
}

func (r *repository) CreateFeed(feed *rss.Feed, items ...*rss.Item) error {
	r.lastID++
	feed.ID = r.lastID
	r.feeds[feed.ID] = feed
	for _, item := range items {
		item.FeedID = feed.ID
		if err := r.CreateItem(item); err != nil {
			log.Printf("error creating item: %v: skipping\n", err)
		}
	}
	return nil
}

func (r *repository) RemoveFeed(id int64) error {
	for iid, item := range r.items {
		if item.FeedID == id {
			delete(r.items, iid)
		}
	}
	return nil
}

func (r *repository) CreateItem(item *rss.Item) error {
	r.lastID++
	item.ID = r.lastID
	r.items[item.ID] = item
	return nil
}

func (r *repository) GetItem(id int64) (*rss.Item, error) {
	item, ok := r.items[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return item, nil
}

func (r *repository) ReadItem(id int64) error {
	if _, ok := r.items[id]; !ok {
		return errors.New("not found")
	}
	r.items[id].Read = true
	return nil
}

func (r *repository) UnreadItem(id int64) error {
	if _, ok := r.items[id]; !ok {
		return errors.New("not found")
	}
	r.items[id].Read = false
	return nil
}

func (r *repository) IgnoreItem(id int64) error {
	if _, ok := r.items[id]; !ok {
		return errors.New("not found")
	}
	r.items[id].Ignored = true
	return nil
}

func (r *repository) UnignoreItem(id int64) error {
	if _, ok := r.items[id]; !ok {
		return errors.New("not found")
	}
	r.items[id].Ignored = true
	return nil
}

func (r *repository) StarItem(id int64) error {
	if _, ok := r.items[id]; !ok {
		return errors.New("not found")
	}
	r.items[id].Starred = true
	return nil
}

func (r *repository) UnstarItem(id int64) error {
	if _, ok := r.items[id]; !ok {
		return errors.New("not found")
	}
	r.items[id].Starred = true
	return nil
}

func (r *repository) ListItems(limit int) ([]*rss.Item, error) {
	var items []*rss.Item
	for _, item := range r.items {
		items = append(items, item)
	}
	if limit < len(items) {
		return items[:limit], nil
	}
	return items, nil
}

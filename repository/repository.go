package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/haleyrc/rss"
)

const (
	AllItems int = 0
)

type Repository interface {
	CreateFeed(feed *rss.Feed, items ...*rss.Item) error
	CreateItem(item *rss.Item) error
	GetItem(id int64) (*rss.Item, error)
	ReadItem(id int64) error
	UnreadItem(id int64) error
	IgnoreItem(id int64) error
	UnignoreItem(id int64) error
	StarItem(id int64) error
	UnstarItem(id int64) error
	ListItems(limit int) ([]*rss.Item, error)
}

func New(db *sqlx.DB) Repository {
	return &repository{db}
}

type repository struct {
	db *sqlx.DB
}

func (r *repository) ListItems(limit int) ([]*rss.Item, error) {
	q := `SELECT id, feed_id, title, link, publication_date, read, ignored, starred FROM items ORDER BY publication_date DESC`
	if limit != AllItems {
		q += fmt.Sprintf("LIMIT %d", limit)
	}
	var items []*rss.Item
	if err := r.db.Select(&items, q); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *repository) setItemRead(id int64, status bool) error {
	q := `UPDATE items SET read = $2 WHERE id = $1`
	_, err := r.db.Exec(q, id, status)
	return err
}

func (r *repository) ReadItem(id int64) error {
	return r.setItemRead(id, true)
}

func (r *repository) UnreadItem(id int64) error {
	return r.setItemRead(id, false)
}

func (r *repository) setItemIgnored(id int64, status bool) error {
	q := `UPDATE items SET ignored = $2 WHERE id = $1`
	_, err := r.db.Exec(q, id, status)
	return err
}

func (r *repository) IgnoreItem(id int64) error {
	return r.setItemIgnored(id, true)
}

func (r *repository) UnignoreItem(id int64) error {
	return r.setItemIgnored(id, false)
}
func (r *repository) setItemStarred(id int64, status bool) error {
	q := `UPDATE items SET starred = $2 WHERE id = $1`
	_, err := r.db.Exec(q, id, status)
	return err
}

func (r *repository) StarItem(id int64) error {
	return r.setItemStarred(id, true)
}

func (r *repository) UnstarItem(id int64) error {
	return r.setItemStarred(id, false)
}

type Getter interface {
	Get(dest interface{}, q string, args ...interface{}) error
}

func (r *repository) GetItem(id int64) (*rss.Item, error) {
	q := `SELECT feed_id, title, link, publication_date, read, ignored, starred FROM items WHERE id = $1`
	var item rss.Item
	if err := r.db.Get(&item, q, id); err != nil {
		return nil, err
	}
	return &item, nil
}

func createItem(g Getter, item *rss.Item) error {
	q := `INSERT INTO items (feed_id, title, link, publication_date) VALUES ($1, $2, $3, $4) ON CONFLICT (feed_id, link) DO UPDATE SET title=EXCLUDED.title, publication_date=EXCLUDED.publication_date RETURNING id`
	return g.Get(item, q, item.FeedID, item.Title, item.Link, item.PublicationDate)
}

func (r *repository) CreateFeed(feed *rss.Feed, items ...*rss.Item) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	q := `INSERT INTO feeds (title, description, link, icon) VALUES ($1, $2, $3, $4) ON CONFLICT (link) DO UPDATE SET description=EXCLUDED.description, title=EXCLUDED.title, icon=EXCLUDED.icon RETURNING id`
	if err := tx.Get(feed, q, feed.Title, feed.Description, feed.Link, feed.Image); err != nil {
		tx.Rollback()
		return err
	}

	for _, item := range items {
		item.FeedID = feed.ID
		if err := createItem(tx, item); err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (r *repository) CreateItem(item *rss.Item) error {
	return createItem(r.db, item)
}

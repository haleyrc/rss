CREATE TABLE IF NOT EXISTS feeds (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL DEFAULT '',
    description TEXT NOT NULL DEFAULT '',
    link TEXT NOT NULL UNIQUE,
    icon TEXT NOT NULL DEFAULT ''
);

CREATE TABLE items (
    id                  SERIAL      PRIMARY KEY,
    feed_id             INTEGER     NOT NULL REFERENCES feeds(id) ON DELETE CASCADE,
    title               TEXT        NOT NULL DEFAULT '',
    link                TEXT        NOT NULL UNIQUE,
    publication_date    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    read                BOOLEAN     NOT NULL DEFAULT false,
    starred             BOOLEAN     NOT NULL DEFAULT false,
    ignored             BOOLEAN     NOT NULL DEFAULT false,
    CONSTRAINT uniq_feed_link_pub UNIQUE (feed_id, link)
);
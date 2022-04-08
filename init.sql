CREATE TABLE IF NOT EXISTS articles
(
    id SERIAL,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
	author TEXT NOT NULL,
    CONSTRAINT articles_pkey PRIMARY KEY (id)
)
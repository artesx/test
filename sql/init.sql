CREATE TABLE stats (
       id SERIAL PRIMARY KEY,
       count INTEGER NOT NULL,
       banner_id INTEGER NOT NULL,
       timestamp TIMESTAMP NOT NULL,
       UNIQUE (banner_id, timestamp)
);

CREATE INDEX idx_stats_banner_timestamp ON stats (banner_id, timestamp);


DROP TABLE IF EXISTS ip;
CREATE EXTENSION IF NOT EXISTS CITEXT;

CREATE TABLE ip
(
    id_ip    BIGSERIAL PRIMARY KEY,
    ip_value CITEXT NOT NULL,
    counter  INT                      DEFAULT 1,
    time     TIMESTAMP WITH TIME ZONE DEFAULT NULL
);

CREATE UNIQUE INDEX idx_ip_uniqueip ON ip(ip_value);
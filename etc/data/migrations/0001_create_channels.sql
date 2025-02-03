-- +goose Up
CREATE TABLE channel (
    id TEXT PRIMARY KEY,
    name TEXT
);

INSERT INTO channel VALUES
('telegram', 'Telegram'),
('whatsapp', 'WhatsApp');

-- +goose Down
DROP TABLE channel;

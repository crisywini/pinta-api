CREATE TABLE outfits (
    id         TEXT PRIMARY KEY DEFAULT gen_random_uuid(),
    name       TEXT        NOT NULL,
    occasions  TEXT[]      NOT NULL DEFAULT '{}',
    seasons    TEXT[]      NOT NULL DEFAULT '{}',
    mood       TEXT        NOT NULL DEFAULT '',
    fragrance  TEXT        NOT NULL DEFAULT '',
    notes      TEXT        NOT NULL DEFAULT '',
    favorite   BOOLEAN     NOT NULL DEFAULT false,
    wear_count INT         NOT NULL DEFAULT 0,
    photo      TEXT        NOT NULL DEFAULT '',
    last_worn  TIMESTAMPTZ
);

CREATE TABLE outfit_items (
    outfit_id TEXT NOT NULL REFERENCES outfits(id) ON DELETE CASCADE,
    item_id   TEXT NOT NULL REFERENCES items(id),
    PRIMARY KEY (outfit_id, item_id)
);

CREATE INDEX idx_outfit_items_item_id ON outfit_items (item_id);
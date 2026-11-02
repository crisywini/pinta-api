CREATE TABLE items (
    id              TEXT PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         TEXT        NOT NULL,
    name            TEXT        NOT NULL,
    category        TEXT        NOT NULL CHECK (category IN (
        'top','bottom','shoes','hosiery','outerwear',
        'accessories','jewelry','innerwear','fragrance')),
    sub_category    TEXT        NOT NULL DEFAULT '',
    primary_color   TEXT        NOT NULL DEFAULT '',
    secondary_color TEXT        NOT NULL DEFAULT '',
    pattern         TEXT        NOT NULL DEFAULT '',
    material        TEXT        NOT NULL DEFAULT '',
    fit             TEXT        NOT NULL DEFAULT '',
    silhouette      TEXT        NOT NULL DEFAULT '',
    style_tags      TEXT[]      NOT NULL DEFAULT '{}',
    occasions       TEXT[]      NOT NULL DEFAULT '{}',
    season_tags     TEXT[]      NOT NULL DEFAULT '{}',
    image_asset_id  TEXT        NOT NULL DEFAULT '',
    available       BOOLEAN     NOT NULL DEFAULT true
);

CREATE INDEX idx_items_user_id ON items (id);
DROP TABLE words;

CREATE TABLE words (
    word text NOT NULL,
    translation text NOT NULL,
    createdate text DEFAULT (STRFTIME('%Y-%m-%d %H:%M:%f', 'NOW')),
    lastreviewdate text DEFAULT (STRFTIME('%Y-%m-%d %H:%M:%f', 'NOW')),
    reviewstatus INT DEFAULT 0,
    PRIMARY KEY(word, lastreviewdate)
);

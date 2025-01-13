CREATE TABLE songs (
                       id SERIAL PRIMARY KEY,
                       "group" VARCHAR(255) NOT NULL,
                       song VARCHAR(255) NOT NULL
);

CREATE TABLE song_lyrics (
                             id SERIAL PRIMARY KEY,
                             song_id INT NOT NULL,
                             lyrics TEXT NOT NULL,
                             FOREIGN KEY (song_id) REFERENCES songs(id) ON DELETE CASCADE
);
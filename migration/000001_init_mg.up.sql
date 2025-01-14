CREATE TABLE groups (
                        id SERIAL PRIMARY KEY,
                        name VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE songs (
                       id SERIAL PRIMARY KEY,
                       group_id INT NOT NULL,
                       name VARCHAR(255) NOT NULL,
                       release_date DATE,
                       link VARCHAR(255),
                       FOREIGN KEY (group_id) REFERENCES groups(id)
);

CREATE TABLE song_lyrics (
                             id SERIAL PRIMARY KEY,
                             song_id INT NOT NULL,
                             lyrics_line TEXT NOT NULL,
                             FOREIGN KEY (song_id) REFERENCES songs(id)
);


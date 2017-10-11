-- SONG
CREATE TABLE song (
  id serial PRIMARY KEY,
  name text NOT NULL,
  link text NOT NULL,
  provider text NOT NULL,
  thumbnail text NOT NULL
);


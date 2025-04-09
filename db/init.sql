CREATE TABLE IF NOT EXISTS series (
  id SERIAL PRIMARY KEY,
  title TEXT NOT NULL,
  description TEXT,
  status TEXT NOT NULL,
  current_episode INTEGER NOT NULL DEFAULT 0,
  score INTEGER NOT NULL DEFAULT 0
);
-- +goose UP
CREATE TABLE chirps (
  id UUID PRIMARY KEY,
  body TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  user_id UUID NOT NULL,
  FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

-- +goose DOWN
DROP TABLE chirps;
-- +goose Up
CREATE TABLE oauth_states (
  token TEXT PRIMARY KEY,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  expires_at TIMESTAMPTZ NOT NULL,
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  service TEXT NOT NULL
    CHECK (service IN ('youtube', 'tiktok', 'instagram', 'twitch', 'twitter')),
  channel_id TEXT NOT NULL,
  channel_handle TEXT NOT NULL
);

CREATE TABLE connections (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  service TEXT NOT NULL
    CHECK (service IN ('youtube', 'tiktok', 'instagram', 'twitch', 'twitter')),
  channel_id TEXT NOT NULL,
  channel_handle TEXT NOT NULL,
  access_token TEXT NOT NULL,
  refresh_token TEXT,
  token_expires_at TIMESTAMPTZ,
  scopes TEXT NOT NULL, UNIQUE (user_id, service, channel_id)
);

-- +goose Down
DROP TABLE connections;
DROP TABLE oauth_states;
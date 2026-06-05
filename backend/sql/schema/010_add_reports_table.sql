-- +goose Up
CREATE TABLE reports (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL,
  tracking_start TIMESTAMPTZ NOT NULL DEFAULT now(),
  tracking_end TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE user_reports (
  user_id UUID REFERENCES users(id),
  report_id UUID REFERENCES reports(id),
  PRIMARY KEY (user_id, report_id),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE reports;
DROP TABLE user_reports;
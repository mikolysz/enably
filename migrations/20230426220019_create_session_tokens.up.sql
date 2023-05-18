CREATE TABLE session_tokens (
  id BIGSERIAL PRIMARY KEY,
  email_address text NOT NULL,
  token TEXT NOT NULL,
  created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);
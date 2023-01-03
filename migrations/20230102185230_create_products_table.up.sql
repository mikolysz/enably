CREATE TABLE products (
  id BIGSERIAL PRIMARY KEY,
  category_slug text NOT NULL,
  data jsonb NOT NULL,
  created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
  updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);
ALTER TABLE products
ADD COLUMN approved BOOLEAN NOT NULL DEFAULT FALSE;

UPDATE products
SET approved=TRUE;
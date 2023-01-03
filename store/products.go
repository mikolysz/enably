package store

import (
	"context"
	"fmt"

	"github.com/mikolysz/enably/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresProductsStore is a ProductsStore that uses a Postgres database.
type PostgresProductsStore struct {
	db *pgxpool.Pool
}

// NewPostgresProductsStore returns a PostgresProductsStore using the given connection pool.
func NewPostgresProductsStore(pool *pgxpool.Pool) *PostgresProductsStore {
	return &PostgresProductsStore{pool}
}

// AddProduct inserts a product into the database.
// The returned product will have the "id" field filled in with the ID of the new product.
func (s PostgresProductsStore) AddProduct(c context.Context, p model.Product) (model.Product, error) {
	query := "INSERT INTO products(category_slug, data) VALUES($1, $2) RETURNING id"
	row := s.db.QueryRow(c, query, p.CategorySlug, p.Data)
	if err := row.Scan(&p.ID); err != nil {
		return model.Product{}, fmt.Errorf("error when inserting product: %s", err)
	}
	return p, nil
}

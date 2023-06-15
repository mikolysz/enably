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

// GetProductsByCategory returns all products in the category with the given slug.
// ONLY approved products are returned.
func (s PostgresProductsStore) GetProductsByCategory(c context.Context, slug string) ([]model.Product, error) {
	query := "SELECT id, data, approved FROM products WHERE category_slug = $1 AND approved = true"

	rows, err := s.db.Query(c, query, slug)
	if err != nil {
		return nil, fmt.Errorf("error when querying products: %s", err)
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var p model.Product
		if err := rows.Scan(&p.ID, &p.Data, &p.Approved); err != nil {
			return nil, fmt.Errorf("error when scanning product: %s", err)
		}
		p.CategorySlug = slug
		products = append(products, p)
	}
	return products, nil
}

// GetProductByID returns the product with the given ID.
func (s PostgresProductsStore) GetProductByID(c context.Context, id int) (model.Product, error) {
	query := "SELECT id, category_slug, data, approved FROM products WHERE id = $1"

	var p model.Product
	row := s.db.QueryRow(c, query, id)
	if err := row.Scan(&p.ID, &p.CategorySlug, &p.Data, &p.Approved); err != nil {
		return model.Product{}, fmt.Errorf("error when querying product: %s", err)
	}

	return p, nil
}

// GetProductsNeedingApproval 		returns all products that need approval by the mod team.
func (s PostgresProductsStore) GetProductsNeedingApproval(c context.Context) ([]model.Product, error) {
	query := "SELECT id, category_slug, data, approved FROM products WHERE approved = false"

	rows, err := s.db.Query(c, query)
	if err != nil {
		return nil, fmt.Errorf("error when querying products: %s", err)
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var p model.Product
		if err := rows.Scan(&p.ID, &p.CategorySlug, &p.Data, &p.Approved); err != nil {
			return nil, fmt.Errorf("error when scanning product: %s", err)
		}
		products = append(products, p)
	}
	return products, nil
}

// ApproveProduct 		approves the product with the given ID.
func (s PostgresProductsStore) ApproveProduct(c context.Context, id int) error {
	query := "UPDATE products SET approved = true WHERE id = $1"
	_, err := s.db.Exec(c, query, id)
	if err != nil {
		return fmt.Errorf("error when approving product: %s", err)
	}
	return nil
}

// RejectProduct 		rejects the product with the given ID.
func (s PostgresProductsStore) RejectProduct(c context.Context, id int) error {
	query := "DELETE FROM products WHERE id = $1"
	_, err := s.db.Exec(c, query, id)
	if err != nil {
		return fmt.Errorf("error when deleting product: %s", err)
	}
	return nil
}

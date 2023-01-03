package model

type Product struct {
	ID           int
	CategorySlug string

	// The JSON representation of the product data.
	Data []byte
}

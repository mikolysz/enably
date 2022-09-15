package model

type Category struct {
	// Slug is how we refer to this category in URLs, schemas and so on.
	Slug string `json:"slug"`

	// Name is human-readable.
	Name string `json:"name"`

	// Parent is the slug of the parent category, if any.
	// If this is a top-level category, it is empty.
	Parent string `json:"parent,omitempty"`

	// Subcategories is a list of slugs of subcategories.
	Subcategories []string `json:"subcategories"`
}

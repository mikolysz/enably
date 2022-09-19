package model

type Category struct {
	// Slug is how we refer to this category in URLs, schemas and so on.
	Slug string `json:"slug"`

	// Name is human-readable.
	Name string `json:"name"`

	// Parent is the slug of the parent category, if any.
	// If this is a top-level category, it is empty.
	Parent string `json:"parent,omitempty"`

	// Subcategories is a list of categories which this category is a parent of.
	Subcategories []*SubcategoryInfo `json:"subcategories"`
}

// SubcategoryInfo contains only the information needed to display a category in a list of subcategories.
type SubcategoryInfo struct {
	Slug string `json:"slug"`
	Name string `json:"name"`

	// IsLeafcategory is true if this subcategory has no further subcategories, but contains products instead.
	IsLeafCategory bool `json:"is_leaf_category"`
}

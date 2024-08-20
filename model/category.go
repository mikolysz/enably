package model

type Category struct {
	// Slug is how we refer to this category in URLs, schemas and so on.
	Slug string `json:"slug"`

	// Name is human-readable.
	Name string `json:"name"`

	ShortDescription string `json:"short_description"`

	// Parent is the slug of the parent category, if any.
	// If this is a top-level category, it is empty.
	Parent string `json:"parent,omitempty"`

	// Subcategories is a list of categories which this category is a parent of.
	Subcategories []*SubcategoryInfo `json:"subcategories"`

	// Fieldsets is a slice of fieldsets, each describing a single aspect of the products in this category.
	// This includes fieldsets from parent categories.
	Fieldsets []*Fieldset `json:"fieldsets"`

	// NameField and DescriptionField indicate which fields in the products' JSON representations contain their names and descriptions.
	// This is used when displaying a list of products.
	// The format is fieldset_slug.field_slug.
	NameField, DescriptionField string

	// FeaturedFields should be presented when displaying a list of products.
	FeaturedFields []string
}

// SubcategoryInfo contains only the information needed to display a category in a list of subcategories.
type SubcategoryInfo struct {
	Slug string `json:"slug"`
	Name string `json:"name"`

	// IsLeafCategory is true if this subcategory has no further subcategories, but contains products instead.
	// Non-leaf categories can contain other subcategories, but not products directly.
	IsLeafCategory bool `json:"is_leaf_category"`
}

// IsLeafCategory returns true if this category has no subcategories.
// Leaf categories can only contain products, while non-leaf categories can contain other subcategories, but not products directly.
func (c *Category) IsLeafCategory() bool {
	return len(c.Subcategories) > 0
}

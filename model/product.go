package model

type Product struct {
	ID           int    `json:"id"`
	CategorySlug string `json:"category_slug"`
	Approved     bool   `json:"approved"`

	// maps fieldset slugs to maps of field names to their values
	Data map[string]map[string]any `json:"data"`

	// The fields below aren't stored in the database,
	// as they can be derived from the JSON data and the schema.
	Name           string         `json:"name"`
	Description    string         `json:"description"`
	FeaturedFields map[string]any `json:"featured_fields"`
}

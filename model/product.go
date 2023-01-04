package model

type Product struct {
	ID           int
	CategorySlug string

	// maps fieldset slugs to maps of field names to their values
	Data map[string]map[string]any

	// The fields below aren't stored in the database,
	// as they can be derived from the JSON data and the schema.
	Name           string         `json:"name"`
	Description    string         `json:"description"`
	FeaturedFields map[string]any `json:"featured_fields"`
}

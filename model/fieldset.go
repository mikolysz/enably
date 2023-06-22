package model

// Fieldset is a set of fields that describes a certain aspect of a product.
type Fieldset struct {
	Slug string `json:"slug"` // used by computers
	Name string `json:"name"` // used by humans

	Fields []*Field `json:"fields"`
}

// Field describes a form field.
type Field struct {
	Name     string   `json:"name"`  // the name used in JSON schemas and the database.
	Label    string   `json:"label"` // What is actually shown on the page.
	Type     string   `json:"type"`  // TODO:provide constants for the types we accept.
	Optional bool     `json:"optional"`
	Options  []string `json:"options"`
}

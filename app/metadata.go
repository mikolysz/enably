package app

import (
	"github.com/mikolysz/enably/model"
)

// MetadataService provides information about categories and fieldsets.
type MetadataService struct {
	store MetadataStore
}

// MetadataStore lets you retrieve information about categories and fieldsets.
type MetadataStore interface {
	TopLevelCategories() []*model.SubcategoryInfo
	CategoryBySlug(slug string) (*model.Category, error)
	FieldsetBySlug(slug string) (*model.Fieldset, error)
	AllCategories() ([]*model.Category, error)
	AllFieldsets() ([]*model.Fieldset, error)
}

// NewMetadataService returns a MetadataService that 	uses the given MetadataStore.
func NewMetadataService(store MetadataStore) *MetadataService {
	return &MetadataService{store}
}

// GetRootCategory returns a dummy category that contains all top-level categories.
func (s *MetadataService) GetRootCategory() *model.Category {
	return &model.Category{
		Slug:          "root",
		Name:          "Root",
		Parent:        "",
		Subcategories: s.store.TopLevelCategories(),
	}
}

// GetCategory returns the category with the given slug.
func (s *MetadataService) GetCategory(slug string) (*model.Category, error) {
	return s.store.CategoryBySlug(slug)
}

// GetAllCategories returns all categories.
func (s *MetadataService) GetAllCategories() ([]*model.Category, error) {
	return s.store.AllCategories()
}

// GetAllFIeldsets returns all fieldsets.
func (s *MetadataService) GetAllFieldsets() ([]*model.Fieldset, error) {
	return s.store.AllFieldsets()
}

// GetSchemasForCategory returns the JSON schemas for all fieldsets in the given category.
func (s *MetadataService) GetSchemasForCategory(category *model.Category) (map[string]any, error) {
	fsets := make(map[string]any)

	for _, fieldset := range category.Fieldsets {
		fieldsetSchema, err := s.getSchemaForFieldset(fieldset)
		if err != nil {
			return nil, err
		}
		fsets[fieldset.Slug] = fieldsetSchema
	}
	return fsets, nil
}

// GetSchemaForFieldset returns the JSON schema for the given fieldset.
func (s *MetadataService) getSchemaForFieldset(fieldset *model.Fieldset) (map[string]any, error) {
	props := make(map[string]any)

	for _, field := range fieldset.Fields {
		props[field.Name] = getSchemaForField(field)
	}

	required := make([]string, 0, len(fieldset.Fields))
	for _, field := range fieldset.Fields {
		if !field.Optional {
			required = append(required, field.Name)
		}
	}

	return map[string]any{
		"type":       "object",
		"properties": props,
		"required":   required,
	}, nil
}

// getSchemaForField returns the JSON schema for the given field.
func getSchemaForField(field *model.Field) map[string]any {
	schema := map[string]any{
		"title": field.Label,
	}

	switch field.Type {
	case "short-text", "textarea", "url", "radio-buttons", "dropdown":
		schema["type"] = "string"
	case "checkbox":
		schema["type"] = "boolean"
	}

	if field.Type == "url" {
		schema["format"] = "uri"
	}

	if field.Type == "radio-buttons" || field.Type == "dropdown" {
		schema["enum"] = field.Options
	}

	return schema
}

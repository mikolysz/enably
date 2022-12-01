package app

import (
	"fmt"

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
	FieldsetNamesForCategory(slug string) ([]string, error)
	AllFieldsets() ([]*model.Fieldset, error)
}

// NewMetadataService returns a MetadataService which 	uses the given MetadataStore.
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

// GetFieldsetsForCategory 		returns the fieldsets that products in this category are covered by.
// These include the fieldsets associate with this category, as well as the fieldsets associated with all its parent categories.
func (s *MetadataService) GetFieldsetsForCategory(slug string) ([]*model.Fieldset, error) {
	cat, err := s.store.CategoryBySlug(slug)
	if err != nil {
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	// 	We want the fieldsets from the parent categories to come first.
	var fsets []*model.Fieldset
	if cat.Parent != "" {
		fsets, err = s.GetFieldsetsForCategory(cat.Parent)
		if err != nil {
			// We don't have to wrap the error here,as we're calling ourselves recursively.
			return nil, err
		}
	}

	// Get the fieldsets for this category.
	names, err := s.store.FieldsetNamesForCategory(slug)
	if err != nil {
		// If we can't get stuff from the schema, something must have gone seriously wrong.
		panic(fmt.Sprintf("failed to get fieldsets for category: %v", err))
	}

	for _, name := range names {
		fset, err := s.store.FieldsetBySlug(name)
		if err != nil {
			panic(fmt.Sprintf("failed to get fieldset: %v", err))
		}

		fsets = append(fsets, fset)
	}
	return fsets, nil
}

// GetAllFieldsets returns all defined fieldsets.
func (s *MetadataService) GetAllFieldsets() ([]*model.Fieldset, error) {
	return s.store.AllFieldsets()
}

// GetSchemaForFieldset returns the JSON schema for the given fieldset.
func (s *MetadataService) GetSchemaForFieldset(fieldset *model.Fieldset) (map[string]any, error) {
	props := make(map[string]any)

	for _, field := range fieldset.Fields {
		props[field.Name] = map[string]any{
			"type":  getSchemaType(field.Type),
			"title": field.Label,
		}
	}
	return map[string]any{
		"type":       "object",
		"properties": props,
	}, nil
}

// get	SchemaType returns the JSON schema type for the given field type.
func getSchemaType(fieldType string) string {
	switch fieldType {
	case "short-text", "textarea":
		return "string"
	default:
		panic(fmt.Sprintf("unknown field type: %v", fieldType))
	}
}

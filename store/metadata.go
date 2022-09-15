package store

import (
	"fmt"
	"net/http"

	"github.com/pelletier/go-toml/v2"

	"github.com/mikolysz/enably/model"
)

// TOMLMetadataStore lets you retrieve category and fieldset metadata stored in the schema. toml file.
type TOMLMetadataStore struct {
	s schema
}

type schema struct {
	Categories map[string]*model.Category
}

// NewTOMLMetadataStore returns a TOMLMetadataStore which uses the given TOML schema.
func NewTOMLMetadataStore(schemaData []byte) *TOMLMetadataStore {
	var s schema

	if err := toml.Unmarshal(schemaData, &s); err != nil {
		panic(fmt.Sprintf("failed to unmarshal schema: %v", err))
	}

	// Fill in the Slug field of each category.
	for slug, cat := range s.Categories {
		cat.Slug = slug
	}

	// Populate the list of subcategories for each category.
	for _, cat := range s.Categories {
		if cat.Parent == "" {
			continue
		}

		parent, ok := s.Categories[cat.Parent]
		if !ok {
			panic(fmt.Sprintf("category %q has parent %q, but no such category exists", cat.Slug, cat.Parent))
		}

		parent.Subcategories = append(parent.Subcategories, cat.Slug)
	}

	return &TOMLMetadataStore{s}
}

// GetTopLevelCategories returns all categories that have no parent.
func (s *TOMLMetadataStore) GetTopLevelCategories() []*model.Category {
	var cats []*model.Category

	for _, cat := range s.s.Categories {
		if cat.Parent == "" {
			cats = append(cats, cat)
		}
	}

	return cats
}

// GetCategory returns the category with the given slug.
func (s *TOMLMetadataStore) GetCategory(slug string) (*model.Category, error) {
	cat, ok := s.s.Categories[slug]
	if !ok {
		return nil, &model.UserFacingError{
			HTTPStatusCode:    http.StatusNotFound,
			UserFacingMessage: fmt.Sprintf("no such category: %q", slug),
		}
	}

	return cat, nil
}

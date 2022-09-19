package store

import (
	"fmt"
	"net/http"

	"github.com/pelletier/go-toml/v2"

	"github.com/mikolysz/enably/model"
)

// TOMLMetadataStore lets you retrieve category and fieldset metadata stored in the schema. toml file.
type TOMLMetadataStore struct {
	categories         map[string]*model.Category
	topLevelCategories []*model.SubcategoryInfo
}

type schema struct {
	Categories map[string]category
}

// category is the TOML representation of a model.Category.
type category struct {
	// We take the slug from the categories map, so we don't need to store it here.

	Name   string
	Parent string //empty if this is a top-level category
}

// NewTOMLMetadataStore returns a TOMLMetadataStore which uses the given TOML schema.
func NewTOMLMetadataStore(schemaData []byte) *TOMLMetadataStore {
	var s schema

	if err := toml.Unmarshal(schemaData, &s); err != nil {
		panic(fmt.Sprintf("failed to unmarshal schema: %v", err))
	}

	st := &TOMLMetadataStore{

		categories: make(map[string]*model.Category),
	}

	// Populate the categories map.
	for slug, cat := range s.Categories {
		st.categories[slug] = &model.Category{
			Slug:   slug,
			Name:   cat.Name,
			Parent: cat.Parent,
		}
	}

	// Populate the subcategory list of each category.
	for _, cat := range st.categories {
		if cat.Parent == "" {
			continue
		}

		parent, ok := st.categories[cat.Parent]
		if !ok {
			panic(fmt.Sprintf("category %q has parent %q, but no such category exists", cat.Slug, cat.Parent))
		}

		parent.Subcategories = append(parent.Subcategories, &model.SubcategoryInfo{
			Slug: cat.Slug,
			Name: cat.Name,
		})
	}

	// Determine which subcategories are leaf categories.
	for _, cat := range st.categories {
		for _, subcat := range cat.Subcategories {
			subcat.IsLeafCategory = len(st.categories[subcat.Slug].Subcategories) == 0
		}
	}

	// Populate the topLevelCategories slice.
	for _, cat := range st.categories {
		if cat.Parent == "" {
			st.topLevelCategories = append(st.topLevelCategories, &model.SubcategoryInfo{
				Slug:           cat.Slug,
				Name:           cat.Name,
				IsLeafCategory: len(cat.Subcategories) == 0,
			})
		}
	}

	return st
}

// GetTopLevelCategories returns all categories that have no parent.
func (s *TOMLMetadataStore) GetTopLevelCategories() []*model.SubcategoryInfo {
	return s.topLevelCategories
}

// GetCategory returns the category with the given slug.
func (s *TOMLMetadataStore) GetCategory(slug string) (*model.Category, error) {
	cat, ok := s.categories[slug]
	if !ok {
		return nil, &model.UserFacingError{
			HTTPStatusCode:    http.StatusNotFound,
			UserFacingMessage: fmt.Sprintf("no such category: %q", slug),
		}
	}

	return cat, nil
}

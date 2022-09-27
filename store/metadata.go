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

	fieldsets map[string]*model.Fieldset

	// categoryFieldsets maps category slugs to slices of fieldset slugs.
	// Fieldsets from parent categories aren't included.
	categoryFieldsets map[string][]string
}

type schema struct {
	Categories map[string]category

	Fieldsets map[string]*model.Fieldset

	// Fields contains the fields from each fieldset.
	// This is done so that, instead of writing [[fieldsets.<name>.fields]] in the TOML file,
	// we can write [[fields.<name>]].
	// These fields will be placed in the Fieldset itself at store creation.
	Fields map[string][]*model.Field
}

// category is the TOML representation of a model.Category.
type category struct {
	// We take the slug from the categories map, so we don't need to store it here.

	Name      string
	Parent    string //empty if this is a top-level category
	Fieldsets []string
}

// NewTOMLMetadataStore returns a TOMLMetadataStore which uses the given TOML schema.
func NewTOMLMetadataStore(schemaData []byte) *TOMLMetadataStore {
	// Deserialize the TOML schema.
	var s schema

	if err := toml.Unmarshal(schemaData, &s); err != nil {
		panic(fmt.Sprintf("failed to unmarshal schema: %v", err))
	}

	st := &TOMLMetadataStore{
		categories:        make(map[string]*model.Category),
		fieldsets:         make(map[string]*model.Fieldset),
		categoryFieldsets: make(map[string][]string),
	}

	st.populateFieldsets(s)
	st.populateCategories(s)
	st.associateFieldsetsWithCategories(s)

	return st
}

// TODO: Refactor so that these functions return categories, fieldsets etc.

func (st *TOMLMetadataStore) populateFieldsets(s schema) {
	// Populate the fieldsets map in the store.
	for slug, fs := range s.Fieldsets {
		fs.Slug = slug
		st.fieldsets[slug] = fs
	}

	// Put the fields from s.Fields in their respective fieldsets.
	for slug, fields := range s.Fields {
		fs, ok := st.fieldsets[slug]

		if !ok {
			panic(fmt.Sprintf("no fieldset with slug %q", slug))
		}

		fs.Fields = fields
	}
}

func (st *TOMLMetadataStore) populateCategories(s schema) {
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
}

func (st *TOMLMetadataStore) associateFieldsetsWithCategories(s schema) {
	for slug, cat := range s.Categories {
		// FIXME: verify that the fieldsets mentioned here actually exist.
		st.categoryFieldsets[slug] = cat.Fieldsets
	}
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

// GetFieldset returns the fieldset with the given slug.
func (s *TOMLMetadataStore) GetFieldset(slug string) (*model.Fieldset, error) {
	fs, ok := s.fieldsets[slug]
	if !ok {
		return nil, &model.UserFacingError{
			HTTPStatusCode:    http.StatusNotFound,
			UserFacingMessage: fmt.Sprintf("no such fieldset: %q", slug),
		}
	}

	return fs, nil
}

// GetFieldsetsForCategory returnsa slice of fieldset slugs
// for the given category.
// Fieldsets from parent categories are not included.
func (s *TOMLMetadataStore) GetFieldsetsForCategory(slug string) ([]string, error) {
	fieldsets, ok := s.categoryFieldsets[slug]
	if !ok {
		return nil, &model.UserFacingError{
			HTTPStatusCode:    http.StatusNotFound,
			UserFacingMessage: fmt.Sprintf("no such category: %q", slug),
		}
	}

	return fieldsets, nil
}

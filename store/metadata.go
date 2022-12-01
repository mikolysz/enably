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

	Name          string
	Parent        string   //empty if this is a top-level category
	FieldsetSlugs []string `toml:"fieldsets"`
}

// NewTOMLMetadataStore returns a TOMLMetadataStore which uses the given TOML schema.
func NewTOMLMetadataStore(schemaData []byte) (*TOMLMetadataStore, error) {
	// Deserialize the TOML schema.
	var s schema

	if err := toml.Unmarshal(schemaData, &s); err != nil {
		return nil, fmt.Errorf("failed to parse TOML schema: %w", err)
	}

	st := &TOMLMetadataStore{
		categories:        make(map[string]*model.Category),
		fieldsets:         make(map[string]*model.Fieldset),
		categoryFieldsets: make(map[string][]string),
	}

	fieldsets, err := st.populateFieldsets(s)
	if err != nil {
		return nil, fmt.Errorf("failed to populate fieldsets: %w", err)
	}
	st.fieldsets = fieldsets

	st.categories, err = st.populateCategories(s)
	if err != nil {
		return nil, fmt.Errorf("failed to populate categories: %w", err)
	}

	st.topLevelCategories, err = st.populateTopLevelCategories()
	if err != nil {
		return nil, fmt.Errorf("failed to populate top-level categories: %w", err)
	}

	st.categoryFieldsets, err = st.associateFieldsetsWithCategories(s)
	if err != nil {
		return nil, fmt.Errorf("failed to set up field-category associations: %w", err)
	}

	return st, nil
}

func (st *TOMLMetadataStore) populateFieldsets(s schema) (map[string]*model.Fieldset, error) {
	fsets := make(map[string]*model.Fieldset)
	for slug, fs := range s.Fieldsets {
		fs.Slug = slug
		fsets[slug] = fs
	}

	// Put the fields from s.Fields in their respective fieldsets.
	for slug, fields := range s.Fields {
		fs, ok := fsets[slug]

		if !ok {
			return nil, fmt.Errorf("found fields block for nonexistent fieldset with slug %q", slug)
		}

		fs.Fields = fields
	}
	return fsets, nil
}

func (st *TOMLMetadataStore) populateCategories(s schema) (map[string]*model.Category, error) {
	cats := make(map[string]*model.Category)
	for slug, cat := range s.Categories {
		cats[slug] = &model.Category{
			Slug:   slug,
			Name:   cat.Name,
			Parent: cat.Parent,
		}
	}

	// Populate the subcategory list of each category.
	for _, cat := range cats {
		if cat.Parent == "" {
			continue
		}

		parent, ok := cats[cat.Parent]
		if !ok {
			panic(fmt.Sprintf("category %q has parent %q, but no such category exists", cat.Slug, cat.Parent))
		}

		parent.Subcategories = append(parent.Subcategories, &model.SubcategoryInfo{
			Slug: cat.Slug,
			Name: cat.Name,
		})
	}

	// Determine which subcategories are leaf categories.
	for _, cat := range cats {
		for _, subcat := range cat.Subcategories {
			subcat.IsLeafCategory = len(cats[subcat.Slug].Subcategories) == 0
		}
	}
	return cats, nil
}

func (st *TOMLMetadataStore) associateFieldsetsWithCategories(s schema) (map[string][]string, error) {
	catFieldsets := make(map[string][]string)
	for catSlug, cat := range s.Categories {
		// verify that the fieldsets exist
		for _, fsSlug := range cat.FieldsetSlugs {
			if _, ok := st.fieldsets[fsSlug]; !ok {
				return nil, fmt.Errorf("category %q has nonexistent fieldset %q", catSlug, fsSlug)
			}
		}

		catFieldsets[catSlug] = cat.FieldsetSlugs
	}
	return catFieldsets, nil
}

func (st *TOMLMetadataStore) populateTopLevelCategories() ([]*model.SubcategoryInfo, error) {
	var topLevel []*model.SubcategoryInfo

	for _, cat := range st.categories {
		if cat.Parent != "" {
			continue
		}

		topLevel = append(topLevel, &model.SubcategoryInfo{
			Slug:           cat.Slug,
			Name:           cat.Name,
			IsLeafCategory: len(cat.Subcategories) == 0,
		})
	}

	return topLevel, nil
}

// TopLevelCategories returns all categories that have no parent.
func (s *TOMLMetadataStore) TopLevelCategories() []*model.SubcategoryInfo {
	return s.topLevelCategories
}

// CategoryBySlug returns the category with the given slug.
func (s *TOMLMetadataStore) CategoryBySlug(slug string) (*model.Category, error) {
	cat, ok := s.categories[slug]
	if !ok {
		return nil, &model.UserFacingError{
			HTTPStatusCode:    http.StatusNotFound,
			UserFacingMessage: fmt.Sprintf("no such category: %q", slug),
		}
	}

	return cat, nil
}

// FieldsetBySlug returns the fieldset with the given slug.
func (s *TOMLMetadataStore) FieldsetBySlug(slug string) (*model.Fieldset, error) {
	fs, ok := s.fieldsets[slug]
	if !ok {
		return nil, &model.UserFacingError{
			HTTPStatusCode:    http.StatusNotFound,
			UserFacingMessage: fmt.Sprintf("no such fieldset: %q", slug),
		}
	}

	return fs, nil
}

// FieldsetNamesForCategory returns a slice of fieldset slugs
// for the given category.
// Fieldsets from parent categories are not included.
func (s *TOMLMetadataStore) FieldsetNamesForCategory(slug string) ([]string, error) {
	fieldsets, ok := s.categoryFieldsets[slug]
	if !ok {
		return nil, &model.UserFacingError{
			HTTPStatusCode:    http.StatusNotFound,
			UserFacingMessage: fmt.Sprintf("no such category: %q", slug),
		}
	}

	return fieldsets, nil
}

// AllFieldsets returns all defined fieldsets.
func (s *TOMLMetadataStore) AllFieldsets() ([]*model.Fieldset, error) {
	fsets := make([]*model.Fieldset, 0, len(s.fieldsets))
	for _, fset := range s.fieldsets {
		fsets = append(fsets, fset)
	}

	return fsets, nil
}

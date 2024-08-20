package store

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/pelletier/go-toml/v2"

	"github.com/mikolysz/enably/model"
)

// TOMLMetadataStore lets you retrieve category and fieldset metadata stored in the schema. toml file.
type TOMLMetadataStore struct {
	categories         map[string]*model.Category
	topLevelCategories []*model.SubcategoryInfo
	fieldsets          map[string]*model.Fieldset
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
	// We take the slug from the section header, so we don't need to store it here.

	Name             string
	ShortDescription string   `toml:"short_description"`
	Parent           string   //empty if this is a top-level category
	FieldsetSlugs    []string `toml:"fieldsets"`

	// field names of the form fieldset.field_name.
	// If empty, take from parent.
	NameField        string   `toml:"name_field"` // indicate which field to use for the product name
	DescriptionField string   `toml:"description_field"`
	FeaturedFields   []string `toml:"featured_fields"` // the fields to show in the products list, overides parent if present.
}

// NewTOMLMetadataStore returns a TOMLMetadataStore which uses the given TOML schema.
func NewTOMLMetadataStore(schemaData []byte) (*TOMLMetadataStore, error) {
	// Deserialize the TOML schema.
	var s schema

	if err := toml.Unmarshal(schemaData, &s); err != nil {
		return nil, fmt.Errorf("failed to parse TOML schema: %w", err)
	}

	st := &TOMLMetadataStore{
		categories: make(map[string]*model.Category),
		fieldsets:  make(map[string]*model.Fieldset),
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
			Slug:             slug,
			Name:             cat.Name,
			ShortDescription: cat.ShortDescription,
			Parent:           cat.Parent,
			NameField:        cat.NameField,
			DescriptionField: cat.DescriptionField,
			FeaturedFields:   cat.FeaturedFields,
		}

		if cat.ShortDescription == "" {
			cats[slug].ShortDescription = cat.Name
		}
	}

	// Set up parent-child relationships.
	for _, cat := range cats {
		if cat.Parent == "" {
			st.topLevelCategories = append(st.topLevelCategories, &model.SubcategoryInfo{
				Slug: cat.Slug,
				Name: cat.Name,
			})
		} else {
			parent, ok := cats[cat.Parent]
			if !ok {
				panic(fmt.Sprintf("category %q has parent %q, but no such category exists", cat.Slug, cat.Parent))
			}

			parent.Subcategories = append(parent.Subcategories, &model.SubcategoryInfo{
				Slug: cat.Slug,
				Name: cat.Name,
			})
		}
	}

	// Determine which categories are leaf categories.

	for _, cat := range cats {
		for _, subcat := range cat.Subcategories {
			subcat.IsLeafCategory = len(cats[subcat.Slug].Subcategories) == 0
		}
	}

	// And now do the same for the top-level categories.
	for _, cat := range st.topLevelCategories {
		cat.IsLeafCategory = len(cats[cat.Slug].Subcategories) == 0
	}

	// Make categories inherit fields from their parents.
	// We start with the top-level categories and descend recursively.
	for _, cat := range st.topLevelCategories {
		if err := st.inheritFields(cat.Slug, cats, s.Categories); err != nil {
			return nil, err
		}
	}
	return cats, nil
}

func (st *TOMLMetadataStore) inheritFields(slug string, cats map[string]*model.Category, schemaCats map[string]category) error {
	cat := cats[slug]

	if cat.Parent != "" {
		parent := cats[cat.Parent]

		cat.Fieldsets = append(parent.Fieldsets[:], cat.Fieldsets...)

		if cat.NameField == "" {
			cat.NameField = parent.NameField
		}

		if cat.DescriptionField == "" {
			cat.DescriptionField = parent.DescriptionField
		}

		if len(cat.FeaturedFields) == 0 {
			cat.FeaturedFields = parent.FeaturedFields
		}
	}

	// Validate that all the fields are present.
	// Non-leaf categories are exempt, as they can't directly contain products anyway.
	if len(cat.Subcategories) == 0 {
		if cat.NameField == "" {
			return fmt.Errorf("category %q has no name field", cat.Slug)
		}

		if !st.verifyFieldsetFieldExists(cat.NameField) {
			return fmt.Errorf("category %q has nonexistent name field %q", cat.Slug, cat.NameField)
		}

		if cat.DescriptionField == "" {
			return fmt.Errorf("category %q has no description field", cat.Slug)
		}

		if !st.verifyFieldsetFieldExists(cat.DescriptionField) {
			return fmt.Errorf("category %q has nonexistent description field %q", cat.Slug, cat.DescriptionField)
		}

		if len(cat.FeaturedFields) == 0 {
			return fmt.Errorf("category %q has no featured fields", cat.Slug)
		}

		for _, field := range cat.FeaturedFields {
			if !st.verifyFieldsetFieldExists(field) {
				return fmt.Errorf("category %q has nonexistent featured field %q", cat.Slug, field)
			}
		}
	}

	// Add fieldsets referenced by this category.
	ownFieldsets := schemaCats[cat.Slug].FieldsetSlugs
	for _, slug := range ownFieldsets {
		fset, ok := st.fieldsets[slug]
		if !ok {
			return fmt.Errorf("category %q references nonexistent fieldset %q", cat.Slug, slug)
		}
		cat.Fieldsets = append(cat.Fieldsets, fset)
	}

	// Recurse into subcategories.
	// FIXME: handle cycles, in case we accidentally put one in our schema.
	for _, subcat := range cat.Subcategories {
		if err := st.inheritFields(subcat.Slug, cats, schemaCats); err != nil {
			return err
		}
	}

	return nil
}

func (st *TOMLMetadataStore) verifyFieldsetFieldExists(field string) bool {
	// Split the field into the fieldset name and field name.
	parts := strings.Split(field, ".")

	if len(parts) != 2 {
		return false
	}

	fieldset, ok := st.fieldsets[parts[0]]
	if !ok {
		return false
	}

	for _, f := range fieldset.Fields {
		if f.Name == parts[1] {
			return true
		}
	}
	return false
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

// AllCategories returns all defined categories.
func (s *TOMLMetadataStore) AllCategories() ([]*model.Category, error) {
	cats := make([]*model.Category, 0, len(s.categories))
	for _, cat := range s.categories {
		cats = append(cats, cat)
	}
	return cats, nil
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

// AllFieldsets returns all defined fieldsets.
func (s *TOMLMetadataStore) AllFieldsets() ([]*model.Fieldset, error) {
	fsets := make([]*model.Fieldset, 0, len(s.fieldsets))
	for _, fset := range s.fieldsets {
		fsets = append(fsets, fset)
	}

	return fsets, nil
}

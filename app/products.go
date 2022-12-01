package app

import (
	"encoding/json"
	"fmt"

	"github.com/santhosh-tekuri/jsonschema/v5"
)

// ProductsService provides operations to retrieve, create and update products.
type ProductsService struct {
	meta    *MetadataService
	schemas map[string]*jsonschema.Schema
}

// NewProductsService returns a new ProductsService.
// The returned service will use the given MetadataService to retrieve category information.
func NewProductsService(meta *MetadataService) (*ProductsService, error) {
	s := &ProductsService{meta: meta}

	// Compile all fieldset schemas
	s.schemas = make(map[string]*jsonschema.Schema)

	fsets, err := meta.GetAllFieldsets()
	if err != nil {
		return nil, fmt.Errorf("error when retrieving fieldsets: %w", err)
	}

	for _, fset := range fsets {
		schema, err := meta.GetSchemaForFieldset(fset)
		if err != nil {
			return nil, fmt.Errorf("error when getting JSON schema for fieldset %s: %w", fset.Slug, err)
		}

		// the jsonschema package wants us to supply the schemas as raw JSON.
		encoded, err := json.Marshal(schema)
		if err != nil {
			return nil, fmt.Errorf("error when encoding schema for fieldset %s: %w", fset.Slug, err)
		}
		// Our schemas don't have a URL, but jsonschema requires one, so let's just make one up.
		url := "https://enably.me/schemas/" + fset.Slug

		compiled, err := jsonschema.CompileString(url, string(encoded))
		if err != nil {
			return nil, fmt.Errorf("error when compiling JSON schema for fieldset %s: %w", fset.Slug, err)
		}

		s.schemas[fset.Slug] = compiled
	}

	return s, nil
}

// CreateProduct creates a product in a given category.
//
// It accepts two arguments, the slug of the category to create the product in
// and a map whose keys are fieldset slugs and whose values are the decoded JSON representations satisfying the given fieldset's schema.
func (s *ProductsService) CreateProduct(categorySlug string, jsonData []byte) error {
	// Validate the received data:

	fieldsets := map[string]any{} // maps fieldset slugs to fields
	if err := json.Unmarshal(jsonData, &fieldsets); err != nil {
		// FIXME: This should cause a BadRequest, not an internalServerError.
		return fmt.Errorf("error when unmarshalling product JSON: %w", err)
	}

	fsets, err := s.meta.GetFieldsetsForCategory(categorySlug)
	if err != nil {
		// FIXME: This should cause a BadRequest, not an internalServerError.
		return fmt.Errorf("error when retrieving fieldsets for category %s: %w", categorySlug, err)
	}

	for _, fset := range fsets {
		fsetData, ok := fieldsets[fset.Slug]
		if !ok {
			// FIXME: This should cause a BadRequest, not an internalServerError.
			return fmt.Errorf("product in category %s doesn't contain fieldset %s", categorySlug, fset.Slug)
		}
		if err := s.schemas[fset.Slug].Validate(fsetData); err != nil {
			// FIXME: This should cause a BadRequest, not an internalServerError.
			return fmt.Errorf("error when validating schema for fieldset %s: %w", fset.Slug, err)
		}

	}

	// TODO: store the product in the DB.
	return nil
}

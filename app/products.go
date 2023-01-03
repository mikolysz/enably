package app

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mikolysz/enably/model"
	"github.com/santhosh-tekuri/jsonschema/v5"
)

// ProductsService provides operations to retrieve, create and update products.
type ProductsService struct {
	meta  *MetadataService
	store ProductsStore

	// compiled JSON schemas for all fieldsets, used for validation.
	schemas map[string]*jsonschema.Schema
}

// ProductsStore is an interface for a store that can retrieve, create and update products.
type ProductsStore interface {
	AddProduct(c context.Context, p model.Product) (model.Product, error)
	//GetProduct(c context.Context, id int) (Product, error)
	//UpdateProduct(c context.Context, p Product) (Product, error)
}

// NewProductsService returns a new ProductsService.
// The returned service will use the given MetadataService to retrieve category information.
func NewProductsService(meta *MetadataService, products ProductsStore) (*ProductsService, error) {
	s := &ProductsService{meta: meta, store: products}

	// Compile all  schemas
	s.schemas = make(map[string]*jsonschema.Schema)

	fsets, err := meta.GetAllFieldsets()
	if err != nil {
		return nil, fmt.Errorf("error when retrieving fieldsets: %w", err)
	}

	for _, fset := range fsets {
		schema, err := meta.getSchemaForFieldset(fset)
		if err != nil {
			return nil, fmt.Errorf("error when retrieving schema for fieldset %s: %w", fset.Slug, err)
		}

		// the jsonschema package wants us to supply the schemas as raw JSON.
		encoded, err := json.Marshal(schema)
		if err != nil {
			return nil, fmt.Errorf("error when encoding schema for fieldset %s: %w", fset.Slug, err)
		}
		// Our schemas don't have a URL, but jsonschema requires one, so let's just make one up.
		url := fmt.Sprintf("https://enably.me/schemas/%s", fset.Slug)

		compiled, err := jsonschema.CompileString(url, string(encoded))
		if err != nil {
			return nil, fmt.Errorf("error when compiling JSON schema for fieldset %s: %w", fset.Slug, err)
		}

		s.schemas[fset.Slug] = compiled
	}

	return s, nil
}

// CreateProduct creates a product in the specified category.
//
// Accepts the slug of the category to create the product in and a map of fieldset slugs to
// decoded JSON representations that satisfy the corresponding fieldset's schema.
func (s *ProductsService) CreateProduct(categorySlug string, jsonData []byte) (model.Product, error) {
	// Validate the received data:

	fieldsets := map[string]any{} // maps fieldset slugs to fields
	if err := json.Unmarshal(jsonData, &fieldsets); err != nil {
		// FIXME: This should cause a BadRequest, not an internalServerError.
		return model.Product{}, fmt.Errorf("error when unmarshalling product JSON: %w", err)
	}

	// Validate each fieldset's data against the corresponding schema.
	cat, err := s.meta.GetCategory(categorySlug)
	if err != nil {
		return model.Product{}, fmt.Errorf("error when retrieving category %s: %w", categorySlug, err)
	}

	for _, fset := range cat.Fieldsets {
		fsetData, ok := fieldsets[fset.Slug]
		if !ok {
			// FIXME: This should cause a BadRequest, not an internalServerError.
			return model.Product{}, fmt.Errorf("product in category %s doesn't contain fieldset %s", categorySlug, fset.Slug)
		}

		if err := s.schemas[fset.Slug].Validate(fsetData); err != nil {
			// FIXME: This should cause a BadRequest, not an internalServerError.
			return model.Product{}, fmt.Errorf("error when validating schema for fieldset %s: %w", fset.Slug, err)
		}
	}

	prod, err := s.store.AddProduct(context.Background(), model.Product{
		CategorySlug: categorySlug,
		Data:         jsonData,
	})
	if err != nil {
		return model.Product{}, fmt.Errorf("error when inserting product: %w", err)
	}

	return prod, nil
}

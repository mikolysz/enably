package app

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

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
	GetProductsByCategory(c context.Context, slug string) ([]model.Product, error)
	GetProductByID(c context.Context, id int) (model.Product, error)
	GetProductsNeedingApproval(c context.Context) ([]model.Product, error)
	ApproveProduct(c context.Context, id int) error
	RejectProduct(c context.Context, id int) error
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

	decoded := map[string]map[string]any{}
	if err := json.Unmarshal(jsonData, &decoded); err != nil {
		// FIXME: This should cause a BadRequest, not an internalServerError.
		return model.Product{}, fmt.Errorf("error when unmarshalling product JSON: %w", err)
	}

	// Validate each fieldset's data against the corresponding schema.
	cat, err := s.meta.GetCategory(categorySlug)
	if err != nil {
		return model.Product{}, fmt.Errorf("error when retrieving category %s: %w", categorySlug, err)
	}

	for _, fset := range cat.Fieldsets {
		fsetData, ok := decoded[fset.Slug]
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
		Data:         decoded,
	})
	if err != nil {
		return model.Product{}, fmt.Errorf("error when inserting product: %w", err)
	}

	return prod, nil
}

// GetProductsByCategory returns all products in the specified category.
func (s *ProductsService) GetProductsByCategory(categorySlug string) ([]model.Product, error) {
	prods, err := s.store.GetProductsByCategory(context.Background(), categorySlug)
	if err != nil {
		return nil, fmt.Errorf("error when retrieving products: %w", err)
	}

	for i := range prods {
		if err := s.SetDerivedFields(&prods[i]); err != nil {
			return nil, fmt.Errorf("error when setting derived fields for product %d: %w", prods[i].ID, err)
		}
	}
	return prods, nil
}

// GetProductByID returns the product with the specified ID.
func (s *ProductsService) GetProductByID(id int) (model.Product, error) {
	prod, err := s.store.GetProductByID(context.Background(), id)
	if err != nil {
		return model.Product{}, fmt.Errorf("error when retrieving product %d: %w", id, err)
	}

	if err := s.SetDerivedFields(&prod); err != nil {
		return model.Product{}, fmt.Errorf("error when setting derived fields for product %d: %w", prod.ID, err)
	}
	return prod, nil
}

// GetProductsNeedingApproval returns all products that need approval by the mod team.
func (s *ProductsService) GetProductsNeedingApproval() ([]model.Product, error) {
	prods, err := s.store.GetProductsNeedingApproval(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error when retrieving products: %w", err)
	}

	for i := range prods {
		if err := s.SetDerivedFields(&prods[i]); err != nil {
			return nil, fmt.Errorf("error when setting derived fields for product %d: %w", prods[i].ID, err)
		}
	}

	return prods, nil
}

// ApproveProduct approves the product with the specified ID.
func (s *ProductsService) ApproveProduct(id int) error {
	if err := s.store.ApproveProduct(context.Background(), id); err != nil {
		return fmt.Errorf("error when approving product %d: %w", id, err)
	}
	return nil
}

// RejectProduct rejects the product with the specified ID.
func (s *ProductsService) RejectProduct(id int) error {
	if err := s.store.RejectProduct(context.Background(), id); err != nil {
		return fmt.Errorf("error when rejecting product %d: %w", id, err)
	}
	return nil
}

// SetDerivedFields sets the name, description and featured fields of the given product.
// The schema determines which fields from the product's fieldsets are used here.
func (s *ProductsService) SetDerivedFields(p *model.Product) error {
	cat, err := s.meta.GetCategory(p.CategorySlug)
	if err != nil {
		return fmt.Errorf("error when retrieving category %s: %w", p.CategorySlug, err)
	}

	untypedName, err := s.getField(cat.NameField, p.Data)
	if err != nil {
		return fmt.Errorf("error when retrieving name field for product %d: %w", p.ID, err)
	}

	name, ok := untypedName.(string)
	if !ok {
		return fmt.Errorf("name field for product %d is not a string", p.ID)
	}
	p.Name = name

	untypedDescription, err := s.getField(cat.DescriptionField, p.Data)
	if err != nil {
		return fmt.Errorf("error when retrieving description field for product %d: %w", p.ID, err)
	}

	description, ok := untypedDescription.(string)
	if !ok {
		return fmt.Errorf("description field for product %d is not a string", p.ID)
	}

	p.Description = description

	p.FeaturedFields = map[string]any{}
	for _, field := range cat.FeaturedFields {
		value, err := s.getField(field, p.Data)
		if err != nil {
			return fmt.Errorf("error when retrieving featured field %s for product %d: %w", field, p.ID, err)
		}

		p.FeaturedFields[field] = value
	}

	return nil
}

// getField returns the value of the given field from the given product's data.
// The field name is expected to be in the form "fieldset_slug.field_name".
func (s *ProductsService) getField(field string, data map[string]map[string]any) (value any, err error) {
	parts := strings.Split(field, ".")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid field name %s", field)
	}

	fset, ok := data[parts[0]]
	if !ok {
		return nil, fmt.Errorf("product doesn't contain fieldset %s", parts[0])
	}

	value, ok = fset[parts[1]]
	if !ok {
		return nil, fmt.Errorf("fieldset %s doesn't contain field %s", parts[0], parts[1])
	}
	return value, nil
}

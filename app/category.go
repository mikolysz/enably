package app

import "github.com/mikolysz/enably/model"

// MetadataService provides information about categories and fieldsets.
type MetadataService struct {
	store MetadataStore
}

// MetadataStore lets you retrieve information about categories and fieldsets.
type MetadataStore interface {
	GetTopLevelCategories() []*model.SubcategoryInfo
	GetCategory(slug string) (*model.Category, error)
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
		Subcategories: s.store.GetTopLevelCategories(),
	}
}

// GetCategory returns the category with the given slug.
func (s *MetadataService) GetCategory(slug string) (*model.Category, error) {
	return s.store.GetCategory(slug)
}

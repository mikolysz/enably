package app

import "github.com/mikolysz/enably/model"

// MetadataService provides information about categories and fieldsets.
type MetadataService struct {
	store MetadataStore
}

// MetadataStore lets you retrieve information about categories and fieldsets.
type MetadataStore interface {
	GetTopLevelCategories() []*model.Category
	GetCategory(slug string) (*model.Category, error)
}

// NewMetadataService returns a MetadataService which 	uses the given MetadataStore.
func NewMetadataService(store MetadataStore) *MetadataService {
	return &MetadataService{store}
}

// GetTopLevelCategories returns all categories that have no parent.
func (s *MetadataService) GetTopLevelCategories() []*model.Category {
	return s.store.GetTopLevelCategories()
}

// GetCategory returns the category with the given slug.
func (s *MetadataService) GetCategory(slug string) (*model.Category, error) {
	return s.store.GetCategory(slug)
}

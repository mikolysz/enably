// For documentation, see the Go types in package model.

export interface Category {
  slug: string;
  name: string;
  parent: string;
  subcategories: SubcategoryInfo[];
}

export interface SubcategoryInfo {
  slug: string;
  name: string;
  is_leaf_category: boolean;
}

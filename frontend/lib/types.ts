// For documentation, see the Go types in package model.

import { RJSFSchema } from "@rjsf/utils";

export interface Category {
  slug: string;
  name: string;
  parent: string;
  subcategories: SubcategoryInfo[];
  fieldsets: Fieldset[];
}

export interface SubcategoryInfo {
  slug: string;
  name: string;
  is_leaf_category: boolean;
}

export interface Schemas {
  [key: string]: RJSFSchema;
}

export interface Fieldset {
  slug: string;
  name: string;
  fields: Field[];
}

export interface Field {
  name: string;
  label: string;
  type: string;
}

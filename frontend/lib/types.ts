// For documentation, see the Go types in package model.

import { RJSFSchema } from "@rjsf/utils";

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

export interface Fieldset {
  slug: string;
  name: string;
  fields: Field[];
  json_schema: RJSFSchema;
}

export interface Field {
  name: string;
  label: string;
  type: string;
}

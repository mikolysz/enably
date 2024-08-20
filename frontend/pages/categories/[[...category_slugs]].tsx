import { GetStaticPaths, GetStaticProps } from "next";
import Link from "next/link";
import { useApi } from "../../lib/api";
import { Category, SubcategoryInfo } from "../../lib/types";

import type { PageWithLayout } from "../../components/Layout";

export const getStaticPaths: GetStaticPaths = async () => {
  return {
    // FIXME: We should generate these statically, but recursively fetching all
    // categories would be annoying, and probably too slow in development.
    paths: [],

    fallback: "blocking",
  };
};

export const getStaticProps: GetStaticProps<Props> = async ({ params }) => {
  // iw we navigate to /categories/ (AKA the root category) the param is missing, not empty.
  // If it is not empty, Next guarantees it is an array of strings because of the path.
  const categorySlugs = (params?.category_slugs ?? []) as string[];

  const isRootCategory = categorySlugs.length === 0;

  const slugToFetch = isRootCategory
    ? ""
    : categorySlugs[categorySlugs.length - 1];

  const categoryResponse = await fetch(
    `${process.env.API_URL}/categories/${slugToFetch}`
  );
  const category = await categoryResponse.json();

  return {
    props: {
      category,
      categorySlugs,
      isRootCategory,
    },
  };
};

interface Props {
  categorySlugs: string[];
  category: Category;
  isRootCategory: boolean;
}

const Categories: PageWithLayout<Props> = ({
  categorySlugs,
  category,
  isRootCategory,
}) => {
  const subcategory_link = ({ slug, is_leaf_category }: SubcategoryInfo) => {
    if (is_leaf_category) {
      return `/products/by-category/${slug}`;
    } else {
      return `/categories/${categorySlugs.join("/")}/${slug}`;
    }
  };

  return (
    <>
      {isRootCategory ? (
        <h1>Categories</h1>
      ) : (
        <h1>{category.short_description}</h1>
      )}

      <p>
        This category contains {category.subcategories.length} subcategories.
      </p>

      <ul>
        {category.subcategories.map((subcategory) => (
          <li key={subcategory.slug}>
            <Link href={subcategory_link(subcategory)}>{subcategory.name}</Link>
          </li>
        ))}
      </ul>
    </>
  );
};

Categories.getTitle = ({ category, isRootCategory }) =>
  `${isRootCategory ? "Categories" : category.name} | Enably`;

export default Categories;

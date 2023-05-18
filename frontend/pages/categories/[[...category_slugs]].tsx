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
  const category_slugs = params?.category_slugs
    ? (params.category_slugs as string[])
    : [];

  const isRootCategory = category_slugs.length === 0;

  const slug_to_fetch = isRootCategory
    ? ""
    : category_slugs[category_slugs.length - 1];

  const categoryResponse = await fetch(
    `${process.env.API_URL}/categories/${slug_to_fetch}`
  );
  const category = await categoryResponse.json();

  return {
    props: {
      category,
      category_slugs,
      isRootCategory,
    },
  };
};

interface Props {
  category_slugs: string[];
  category: Category;
  isRootCategory: boolean;
}

const Categories: PageWithLayout<Props> = ({
  category_slugs,
  category,
  isRootCategory,
}) => {
  const subcategory_link = ({ slug, is_leaf_category }: SubcategoryInfo) => {
    if (is_leaf_category) {
      return `/products/by-category/${slug}`;
    } else {
      return `/categories/${category_slugs.join("/")}/${slug}`;
    }
  };

  return (
    <>
      {isRootCategory ? <h1>Categories</h1> : <h1>{category.name}</h1>}

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

import Link from "next/link";
import { useApi } from "../../lib/api";
import { GetStaticPaths, GetStaticProps } from "next";
import { Category, SubcategoryInfo } from "../../lib/types";

export const getStaticPaths: GetStaticPaths = async () => {
  return {
    // FIXME: We should generate these statically, but recursively fetching all
    // categories would be annoying, and probably too slow in development.
    paths: [],

    fallback: "blocking",
  };
};

export const getStaticProps: GetStaticProps = async ({ params }) => {
  return {
    props: {
      category_slugs: params?.category_slugs ? params.category_slugs : [],
    },
  };
};

interface Props {
  category_slugs: string[];
}

export default function Categories({ category_slugs }: Props) {
  const isRootCategory = category_slugs.length === 0;

  const slug_to_fetch = isRootCategory
    ? ""
    : category_slugs[category_slugs.length - 1];
  const path_to_fetch = `categories/${slug_to_fetch}`;

  const { data: category, error: categories_error } =
    useApi<Category>(path_to_fetch);

  if (categories_error) {
    console.log(categories_error);
    return <div>failed to load</div>;
  }

  if (!category) {
    return <div>loading...</div>;
  }

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

      {isRootCategory || (
        <p>
          This category contains {category.subcategories.length} subcategories.
        </p>
      )}

      <ul>
        {category.subcategories.map((subcategory) => (
          <li key={subcategory.slug}>
            <Link href={subcategory_link(subcategory)}>{subcategory.name}</Link>
          </li>
        ))}
      </ul>
    </>
  );
}

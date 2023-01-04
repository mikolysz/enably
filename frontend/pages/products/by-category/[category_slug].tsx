import { GetStaticPaths, GetStaticProps } from "next";
import { Category, Product } from "../../../lib/types";
import { useApi } from "../../../lib/api";
import Link from "next/link";

export const getStaticPaths: GetStaticPaths = async () => {
  return {
    // FIXME: generate statically
    paths: [],
    fallback: "blocking",
  };
};

export const getStaticProps: GetStaticProps = async ({ params }) => {
  return {
    props: {
      category_slug: params?.category_slug ? params.category_slug : "",
    },
  };
};

export default function CategoryPage({
  category_slug,
}: {
  category_slug: string;
}) {
  const { data: products, error: products_error } = useApi<Product[]>(
    `products/by-category/${category_slug}`
  );

  const { data: category, error: category_error } = useApi<Category>(
    `categories/${category_slug}`
  );

  if (products_error) {
    console.log(products_error);
    return <div>failed to load</div>;
  }

  if (category_error) {
    console.log(category_error);
    return <div>failed to load</div>;
  }

  if (!products || !category) {
    return <div>loading...</div>;
  }

  return (
    <>
      <h1>{category.name}</h1>
      <Link href={`/submit/${category.slug}`}>Submit a product</Link>
      <ul>
        {products.map((product) => (
          <li key={product.id}>
            <a href={`/products/${product.id}`}>{product.name}</a>
          </li>
        ))}
      </ul>
    </>
  );
}

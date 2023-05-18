import { GetStaticPaths, GetStaticProps } from "next";
import { Category, Product } from "../../../lib/types";
import { useApi } from "../../../lib/api";
import Link from "next/link";
import { PageWithLayout } from "../../../components/Layout";

export const getStaticPaths: GetStaticPaths = async () => {
  return {
    // FIXME: generate statically
    paths: [],
    fallback: "blocking",
  };
};

export const getStaticProps: GetStaticProps = async ({ params }) => {
  const category_slug = params?.category_slug as string;
  const productsResponse = await fetch(
    `${process.env.API_URL}/products/by-category/${category_slug}`
  );

  const products = await productsResponse.json();

  const categoryResponse = await fetch(
    `${process.env.API_URL}/categories/${category_slug}`
  );

  const category = await categoryResponse.json();

  return {
    props: {
      products,
      category,
    },
  };
};

interface Props {
  products: Product[];
  category: Category;
}

const CategoryPage: PageWithLayout<Props> = ({ products, category }) => {
  return (
    <>
      <h1>{category.name}</h1>
      <Link className="link-secondary" href={`/submit/${category.slug}`}>
        Submit a product
      </Link>
      <ul>
        {products.map((product) => (
          <li key={product.id}>
            <a href={`/products/${product.id}`}>{product.name}</a>
          </li>
        ))}
      </ul>
    </>
  );
};

CategoryPage.getTitle = ({ category }) => `${category.name} | Enably`;

export default CategoryPage;

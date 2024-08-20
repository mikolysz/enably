import { GetStaticPaths, GetStaticProps } from "next";
import { Category, Product } from "../../../lib/types";
import { useApi } from "../../../lib/api";
import Link from "next/link";
import { PageWithLayout } from "../../../components/Layout";
import { LoginLink } from "../../../components/LoginLink";

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

const AddProductLink = ({
  className,
  categorySlug,
  children,
}: {
  className?: string;
  categorySlug: string;
  children: React.ReactNode;
}) => <LoginLink href={`/submit/${categorySlug}`}>{children}</LoginLink>;

const CategoryPage: PageWithLayout<Props> = ({ products, category }) => {
  const productsList = (
    <ul>
      {products.map((product) => (
        <li key={product.id}>
          <a href={`/products/${product.id}`}>{product.name}</a>
        </li>
      ))}
    </ul>
  );

  return (
    <>
      <h1>{category.short_description}</h1>
      <AddProductLink categorySlug={category.slug} className="link-secondary">
        Add New Product
      </AddProductLink>
      {products.length > 0 ? (
        productsList
      ) : (
        <p>
          There are no products in this category yet. Maybe you'd like to{" "}
          <AddProductLink categorySlug={category.slug}>Add one</AddProductLink>?
        </p>
      )}
    </>
  );
};

CategoryPage.getTitle = ({ category }) => `${category.name} | Enably`;

export default CategoryPage;

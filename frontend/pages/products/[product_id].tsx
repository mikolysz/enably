import { GetServerSideProps } from "next";
import { PageWithLayout } from "../../components/Layout";
import { Product } from "../../lib/types";

export const getServerSideProps: GetServerSideProps = async ({ params }) => {
  const productID = params?.product_id as string;
  const productResponse = await fetch(
    `${process.env.API_URL}/products/${productID}`
  );

  const product = await productResponse.json();

  return {
    props: {
      product,
    },
  };
};

interface Props {
  product: Product;
}

const ProductPage: PageWithLayout<Props> = (props) => {
  console.log(props);
  const product = props.product;
  return (
    <>
      <h1>{product.name}</h1>
      {Object.keys(product.data).map((fieldsetName) => (
        <Fieldset
          key={fieldsetName}
          name={fieldsetName}
          fieldset={product.data[fieldsetName]}
        />
      ))}
    </>
  );
};

const Fieldset = ({ name, fieldset }: { name: string; fieldset: any }) => {
  return (
    <>
      <h2> {name} </h2>
      {Object.keys(fieldset).map((fieldName) => (
        <p>
          {fieldName}: {fieldset[fieldName]}
        </p>
      ))}
    </>
  );
};

ProductPage.getTitle = ({ product }) => `${product.name} | Enably`;

export default ProductPage;

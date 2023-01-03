import Form from "@rjsf/core";
import { RJSFSchema, UiSchema } from "@rjsf/utils";
import validator from "@rjsf/validator-ajv6";
import { useRouter } from "next/router";
import { GetStaticPaths, GetStaticProps } from "next";
import { useState, MouseEventHandler } from "react";

import { useApi } from "../../lib/api";
import { Fieldset, Category, Schemas } from "../../lib/types";

interface OnChange {
  formData: any;
}

export const getStaticPaths: GetStaticPaths = async () => {
  // FIXME: actually get the list of categories from the API.
  return {
    paths: [],
    fallback: "blocking",
  };
};

export const getStaticProps: GetStaticProps = async (context) => {
  return {
    props: {
      category_slug: context.params?.category_slug, //FIXME: What to do when category is undefined?
      // FIXME: We could probably fetch the fieldset here too.
    },
  };
};

const Submit = ({ category_slug }: { category_slug: string }) => {
  const { data: schemas, error: schemaError } = useApi<Schemas>(
    `categories/${category_slug}/schemas`
  );

  if (schemaError) {
    console.error(schemaError);
    return <div>Error loading schemas</div>;
  }

  const { data: category, error: categoryError } = useApi<Category>(
    `categories/${category_slug}`
  );

  if (categoryError) {
    console.error(categoryError);
    return <div>Error loading category</div>;
  }

  if (!schemas || !category) {
    return <div>Loading...</div>;
  }

  return <SubmitPage category={category} schemas={schemas} />;
};

const SubmitPage = ({
  schemas,
  category,
}: {
  schemas: Schemas;
  category: Category;
}) => {
  let initialData: any = {};
  for (let fieldset of category.fieldsets) {
    initialData[fieldset.slug] = {};
  }

  const [data, setData] = useState(initialData);

  const onSubmit: MouseEventHandler = (e) => {
    e.preventDefault();
    submit(category.slug, data);
  };

  return (
    <>
      <h1>Submit a product</h1>
      {category.fieldsets.map((fieldset) => {
        const onChange = (e: any) => {
          const newData = {
            ...data,
            [fieldset.slug]: e.formData,
          };
          setData(newData);
        };

        return (
          <Section
            key={fieldset.slug}
            fieldset={fieldset}
            schema={schemas[fieldset.slug]}
            onChange={onChange}
            formData={data[fieldset.slug]}
          />
        );
      })}
      <button onClick={onSubmit}>Submit</button>
    </>
  );
};

const Section = ({
  fieldset,
  onChange,
  formData,
  schema,
}: {
  fieldset: Fieldset;
  onChange: (data: OnChange) => void;
  formData: any;
  schema: RJSFSchema;
}) => {
  return (
    <>
      <h2>{fieldset.name}</h2>
      <Form
        schema={schema}
        uiSchema={ui_schema_for_fieldset(fieldset)}
        validator={validator}
        onChange={onChange}
        formData={formData}
      />
    </>
  );
};

const ui_schema_for_fieldset = ({ fields }: Fieldset): UiSchema => {
  const schema: UiSchema = {
    "ui:submitButtonOptions": {
      norender: true,
    },
    "ui:order": fields.map((field) => field.name),
  };

  return schema;
};

export default Submit;

// FIXME: URLencode the slug, potential vulnerability here.
const submit = (category_slug: string, data: any) =>
  fetch(`/api/v1/products/${category_slug}`, {
    method: "post",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(data),
  });

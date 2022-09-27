import Form from "@rjsf/core";
import { UiSchema } from "@rjsf/utils";
import validator from "@rjsf/validator-ajv6";
import { useRouter } from "next/router";
import { GetStaticPaths, GetStaticProps } from "next";
import { useState, MouseEventHandler } from "react";

import { useApi } from "../../lib/api";
import { Fieldset } from "../../lib/types";

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
  const { data: fieldsets, error } = useApi<Fieldset[]>(
    `categories/${category_slug}/fieldsets`
  );

  if (error) {
    console.error(error);
    return <div>Error loading fieldsets</div>;
  }

  if (!fieldsets) {
    return <div>Loading...</div>;
  }

  return <SubmitPage category_slug={category_slug} fieldsets={fieldsets} />;
};

const SubmitPage = ({
  category_slug,
  fieldsets,
}: {
  category_slug: string;
  fieldsets: Fieldset[];
}) => {
  let initialData: any = {};
  for (let fieldset of fieldsets) {
    initialData[fieldset.slug] = {};
  }

  const [data, setData] = useState(initialData);

  const onSubmit: MouseEventHandler = (e) => {
    e.preventDefault();
    console.log(data);
  };

  return (
    <>
      <h1>Submit a product</h1>
      {fieldsets.map((fieldset) => {
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
}: {
  fieldset: Fieldset;
  onChange: (data: OnChange) => void;
  formData: any;
}) => {
  return (
    <>
      <h2>{fieldset.name}</h2>
      <Form
        schema={fieldset.json_schema}
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

import Head from "next/head";

export default function CustomHead({ title }: { title: string }) {
  return (
    <Head>
      <title>{title}</title>
    </Head>
  );
}

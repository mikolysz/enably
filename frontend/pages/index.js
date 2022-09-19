import Head from "next/head";
import Link from "next/link";

export default function Home() {
  return (
    <>
      <h1>Enably</h1>
      <Link href="/categories">categories</Link>
    </>
  );
}

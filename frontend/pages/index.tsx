import Link from "next/link";

import { PageWithLayout } from "../components/Layout";

const Home: PageWithLayout = () => {
  return (
    <>
      <h1>Enably</h1>
      <Link href="/categories">categories</Link>
    </>
  );
};

export default Home;

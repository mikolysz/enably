import Layout, { AppPropsWithLayout } from "../components/Layout";
import "bootstrap/dist/css/bootstrap.min.css";

function MyApp(props: AppPropsWithLayout) {
  return <Layout {...props} />;
}

export default MyApp;

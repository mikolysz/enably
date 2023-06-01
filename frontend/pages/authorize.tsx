import { useEffect } from "react";
import { GetServerSideProps } from "next";
import { useRouter } from "next/router";

export const getServerSideProps: GetServerSideProps = async (context) => {
  return {
    props: {
      token: context.query.token,
      redirectURI: context.query.redirect_uri,
    },
  };
};

const Authorize = ({
  token,
  redirectURI,
}: {
  token: string;
  redirectURI: string;
}) => {
  const router = useRouter();
  useEffect(() => {
    localStorage.setItem("token", token);

    // set a timer for 1 second that redirects to redirect URI
    setTimeout(() => {
      router.push(redirectURI);
    }, 200);
  });

  return <div>Authorizing...</div>;
};

export default Authorize;

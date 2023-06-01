import { NextPage } from "next";
import { AppProps } from "next/app";
import Link from "next/link";
import { useState } from "react";
import {
  Col,
  Collapse,
  Container,
  Nav,
  Navbar,
  NavbarBrand,
  NavbarText,
  NavbarToggler,
  NavItem,
  NavLink,
  Row,
} from "reactstrap";

import CustomHead from "./CustomHead";
import { LoginModal } from "./LoginModal";

const defaultTitle =
  "Enably: Your go-to resource for accessible products and solutions.";

const getDefaultTitle = (props: {}) => defaultTitle;

export type PageWithLayout<P = {}, IP = P> = NextPage<P, IP> & {
  getTitle?: (props: P) => string;
};

export type AppPropsWithLayout = AppProps & {
  Component: PageWithLayout;
};

export default function Layout({ Component, pageProps }: AppPropsWithLayout) {
  const getTitle = Component.getTitle || getDefaultTitle;
  const title = getTitle(pageProps);
  return (
    <Container>
      <CustomHead title={title} />
      <Row>
        <Col xs="12">
          <OurNavbar />
        </Col>
      </Row>

      <Row>
        <Col md="8">
          <main>
            <Component {...pageProps} />
          </main>
        </Col>
      </Row>

      <Row>
        <Col>
          <footer>
            <p>
              <small>
                Copyright &copy; {new Date().getFullYear()} Mikołaj Hołysz, All
                Rights Reserved
              </small>
            </p>
          </footer>
        </Col>
      </Row>
    </Container>
  );
}

const NavbarLink = ({ href, children }: { href: string; children: string }) => (
  <NavItem>
    <NavLink tag={Link} href={href}>
      {children}
    </NavLink>
  </NavItem>
);

const LoginLink = () => {
  const [isLoginModalOpen, setIsLoginModalOpen] = useState(false);
  const toggleLoginModal = () => setIsLoginModalOpen(!isLoginModalOpen);
  return (
    <>
      <NavItem>
        <NavLink href="#" onClick={toggleLoginModal}>
          Log In
        </NavLink>
      </NavItem>
      <LoginModal
        isOpen={isLoginModalOpen}
        toggle={toggleLoginModal}
        redirectURI=""
      />
    </>
  );
};

const OurNavbar = () => {
  const [isOpen, setIsOpen] = useState(false);
  const toggle = () => setIsOpen(!isOpen);
  return (
    <Navbar color="light" light expand="md">
      <NavbarBrand href="/">Enably</NavbarBrand>
      <NavbarToggler onClick={toggle} />
      <Collapse isOpen={isOpen} navbar>
        <Nav className="mr-auto" navbar>
          <NavbarLink href="/about">About</NavbarLink>
          <NavbarLink href="/categories">Browse</NavbarLink>
          <NavbarText>Not signed in</NavbarText>
          <LoginLink />
        </Nav>
      </Collapse>
    </Navbar>
  );
};

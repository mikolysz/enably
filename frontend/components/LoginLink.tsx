import React, { useEffect, useState } from "react";
import { LoginModal } from "./LoginModal";

export const LoginLink = (props: React.ComponentPropsWithoutRef<"a">) => {
  const [isOpen, setIsOpen] = useState(false);
  const toggle = () => setIsOpen(!isOpen);
  const onClick = (e: React.MouseEvent<HTMLAnchorElement, MouseEvent>) => {
    e.preventDefault();
    toggle();
  };

  const newProps = { ...props };
  newProps.href ||= "#";

  const [isLoggedIn, setIsLoggedIn] = useState(false);
  useEffect(() => {
    setIsLoggedIn(!!localStorage.getItem("token"));
  });

  if (isLoggedIn) {
    return <a {...newProps}>{newProps.children}</a>;
  }

  return (
    <>
      <a onClick={onClick} {...newProps}>
        {props.children}
      </a>
      <LoginModal isOpen={isOpen} toggle={toggle} redirectURI={newProps.href} />
    </>
  );
};

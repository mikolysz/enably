import React, { useState } from "react";
import { LoginModal } from "./LoginModal";

const LoginLink = (props: React.ComponentPropsWithoutRef<"a">) => {
  const [isOpen, setIsOpen] = useState(false);
  const toggle = () => setIsOpen(!isOpen);
  const onClick = (e: React.MouseEvent<HTMLAnchorElement, MouseEvent>) => {
    e.preventDefault();
    toggle();
  };

  const href = props.href || "#";

  return (
    <>
      <a onClick={toggle} {...props} />
      <LoginModal isOpen={isOpen} toggle={toggle} redirectURI={href} />
    </>
  );
};

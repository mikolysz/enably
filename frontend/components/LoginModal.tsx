import React, { useState } from "react";
import {
  Alert,
  Button,
  Form,
  FormGroup,
  Input,
  Label,
  Modal,
  ModalBody,
  ModalFooter,
  ModalHeader,
} from "reactstrap";

import { getAPIResponse } from "../lib/api";
export const LoginModal = ({
  isOpen,
  toggle,
  redirectURI,
}: {
  isOpen: boolean;
  toggle: () => void;
  redirectURI: string;
}) => {
  const [email, setEmail] = useState("");
  const onEmailChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setEmail(e.target.value);
  };

  const [checkEmailAlert, setCheckEmailAlert] = useState(false);
  const dismissCheckEmailAlert = () => setCheckEmailAlert(false);

  const login = async () => {
    toggle();

    const loginRequest = {
      email,
      redirect_uri: `${process.env.NEXT_PUBLIC_ROOT}/${redirectURI}`,
    };

    const response = await getAPIResponse("auth/login", {
      method: "POST",
      body: JSON.stringify(loginRequest),
    });

    setCheckEmailAlert(true);
  };

  return (
    <>
      <Modal isOpen={isOpen} toggle={toggle}>
        <ModalHeader toggle={toggle}>Login Required</ModalHeader>
        <ModalBody>
          <p>Enter an email address to continue</p>
          <small>
            If you don't have an Enably account yet, we will create one for you.
          </small>
          <hr />
          <Label for="email" className="visually-hidden">
            Email
          </Label>
          <Input
            id="email"
            type="email"
            placeholder="john.smith@example.com"
            value={email}
            onChange={onEmailChange}
          />
        </ModalBody>
        <ModalFooter>
          <Button color="primary" onClick={login}>
            Log In
          </Button>
          <Button color="secondary" onClick={toggle}>
            Cancel
          </Button>
        </ModalFooter>
      </Modal>
      <Alert
        color="info"
        isOpen={checkEmailAlert}
        toggle={dismissCheckEmailAlert}
      >
        <h4 className="alert-heading">check your email</h4>
        <p>Please click the link we sent you to continue.{email}</p>
      </Alert>
    </>
  );
};

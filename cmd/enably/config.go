package main

import (
	"fmt"
	"os"
)

type config struct {
	dbConnectionString string
	senderEmail        string
	senderName         string
	sendgridAPIKey     string
}

func loadConfig() (config, error) {

	var c config

	if err := c.setStringValue("DB_CONNECTION_STRING", &c.dbConnectionString); err != nil {
		return config{}, err
	}

	if err := c.setStringValue("SENDER_EMAIL", &c.senderEmail); err != nil {
		return config{}, err
	}

	if err := c.setStringValue("SENDER_NAME", &c.senderName); err != nil {
		return config{}, err
	}

	if err := c.setStringValue("SENDGRID_API_KEY", &c.sendgridAPIKey); err != nil {
		return config{}, err
	}

	return c, nil
}

func (c *config) setStringValue(envVar string, field *string) error {
	value, ok := os.LookupEnv(envVar)
	if !ok {
		return fmt.Errorf("environment variable %s is not set", envVar)
	}
	*field = value
	return nil
}

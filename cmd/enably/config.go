package main

import (
	"fmt"
	"net/url"
	"os"
)

type config struct {
	dbConnectionString string
	senderEmail        string
	senderName         string
	sendgridAPIKey     string
	frontendURL        *url.URL
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

	urlStr, ok := os.LookupEnv("FRONTEND_URL")
	if !ok {
		return config{}, fmt.Errorf("environment variable FRONTEND_URL not found")
	}

	var err error
	c.frontendURL, err = url.Parse(urlStr)
	if err != nil {
		return config{}, fmt.Errorf("FRONTEND_URL %s is not a valid URL: %w", urlStr, err)
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

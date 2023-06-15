package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/mikolysz/enably/model"
	"github.com/urfave/cli/v2"
)

func main() {
	apiURL, ok := os.LookupEnv("ENABLY_API_URL")
	if !ok {
		must(fmt.Errorf("environment variable ENABLY_API_URL not found"))
	}

	apiURL += "/api/v1"

	apiKey, ok := os.LookupEnv("ENABLY_MODERATION_API_KEY")
	if !ok {
		must(fmt.Errorf("environment variable ENABLY_MODERATION_API_KEY not found"))
	}

	header := http.Header{}
	header.Set("X-Moderation-Api-Key", apiKey)

	app := &cli.App{
		Name:  "enctl",
		Usage: "Interact with Enably's moderation API",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "Get a list of products that need approval",
				Action: func(c *cli.Context) error {
					url := apiURL + "/moderation/pending"
					req, err := http.NewRequest(http.MethodGet, url, nil)
					must(err)
					req.Header = header
					resp, err := http.DefaultClient.Do(req)
					must(err)
					defer resp.Body.Close()

					var products []model.Product
					must(json.NewDecoder(resp.Body).Decode(&products))
					for _, p := range products {
						fmt.Printf("%d - %s (%s)\n", p.ID, p.Name, p.CategorySlug)
					}

					return nil
				},
			},
			{
				Name:  "info",
				Usage: "Get information about a product",
				Action: func(c *cli.Context) error {
					id := c.Args().First()
					url := apiURL + "/products/" + id
					resp, err := http.Get(url)
					must(err)
					defer resp.Body.Close()

					var product model.Product
					must(json.NewDecoder(resp.Body).Decode(&product))
					for name, fieldset := range product.Data {
						fmt.Printf("%s:\n", name)
						for fieldName, field := range fieldset {
							fmt.Printf("  %s: %s\n", fieldName, field)
						}
						fmt.Println()
					}

					return nil
				},
			},
			{
				Name:  "approve",
				Usage: "Approve a product",
				Action: func(c *cli.Context) error {
					id := c.Args().First()
					url := apiURL + "/moderation/products/" + id + "/approve"
					req, err := http.NewRequest(http.MethodPost, url, nil)
					must(err)
					req.Header = header
					resp, err := http.DefaultClient.Do(req)
					must(err)
					defer resp.Body.Close()

					return nil
				},
			},
			{
				Name:  "reject",
				Usage: "Reject a product",
				Action: func(c *cli.Context) error {
					id := c.Args().First()
					url := apiURL + "/moderation/products/" + id + "/reject"
					req, err := http.NewRequest(http.MethodPost, url, nil)
					must(err)
					req.Header = header
					resp, err := http.DefaultClient.Do(req)
					must(err)
					defer resp.Body.Close()

					return nil
				},
			},
		},
	}

	must(app.Run(os.Args))
}

func must(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

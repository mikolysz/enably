package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mikolysz/enably"
	"github.com/mikolysz/enably/api"
	"github.com/mikolysz/enably/app"
	"github.com/mikolysz/enably/store"
)

func main() {
	metaStore := store.NewTOMLMetadataStore(enably.Schema)
	meta := app.NewMetadataService(metaStore)

	deps := api.Dependencies{
		Metadata: meta,
	}

	a := api.New(deps)

	r := chi.NewMux()

	r.Mount("/api/v1", a)

	log.Fatal(http.ListenAndServe(":8080", r))
}

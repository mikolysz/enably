package main

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/mikolysz/enably"
	"github.com/mikolysz/enably/api"
	"github.com/mikolysz/enably/app"
	"github.com/mikolysz/enably/pkg/email/sendgrid"
	"github.com/mikolysz/enably/store"
)

func main() {
	// We don't care if this errors out, a missing .env is fine.
	godotenv.Load()

	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("Error when loading config: %s", err)
	}

	db, err := pgxpool.New(context.Background(), cfg.dbConnectionString)
	if err != nil {
		log.Fatalf("Error when connecting to database: %s", err)
	}
	defer db.Close()

	if err := db.Ping(context.Background()); err != nil {
		log.Fatalf("Error when pinging database: %s", err)
	}

	metaStore, err := store.NewTOMLMetadataStore(enably.Schema)
	if err != nil {
		log.Fatalf("Failed to create metadata store: %s", err)
	}

	meta := app.NewMetadataService(metaStore)

	productsStore := store.NewPostgresProductsStore(db)

	prod, err := app.NewProductsService(meta, productsStore)
	if err != nil {
		log.Fatalf("Error when creating products service: %s", err)
	}

	sendgridConfig := sendgrid.Config{
		APIKey:      cfg.sendgridAPIKey,
		SenderEmail: cfg.senderEmail,
		SenderName:  cfg.senderName,
	}

	emailSender := sendgrid.NewSender(sendgridConfig)

	authStore := store.PostgresTokenStore{DB: db}
	auth := app.NewAuthenticationService(authStore, emailSender, cfg.frontendURL)

	deps := api.Dependencies{
		Metadata: meta,
		Products: prod,
		Auth:     auth,
	}

	a := api.New(deps)

	r := chi.NewMux()

	r.Mount("/api/v1", a)

	log.Fatal(http.ListenAndServe(":8080", r))
}

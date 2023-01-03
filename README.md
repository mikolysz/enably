# Enably: Your go-to resource for accessible products and solutions.

## What is Enably?

Enably is a platform designed to help people with disabilities find products that meet their accessibility needs. Our focus is currently on assisting those who are blind, but we may expand to serve individuals with other disabilities in the future.

While there are often many options available in certain product categories, many of them can be difficult or impossible for blind people to use due to features such as touch controls. This is not just a problem with kitchen appliances, but also with mobile apps, websites, games, and even music and ham-radio equipment.

Enably aims to be a comprehensive resource for accessibility information, offering ratings for products and any necessary workarounds or additional software that may be required to make them more accessible for those with visual impairments. All of our content is created by our users, allowing them to add new products or update existing ones as needed. With Enably, we hope to make it easy for those with disabilities to find the products they need to live their lives to the fullest.

**Note**: This project is a very incomplete, early prototype. For more info, see the [Initial Specification](docs/initial_spec.md).

## How to contribute:

1. Install Node, Postgres, Go and [golang-migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate).
1. Create a postgres role called `enably` with password `enably`, and a database called `enably_dev`. Grant all privileges on that database to this role.
1. Clone the repo.
1. Run the migrations with `migrate -path=migrations -database='postgres://enably:enably@localhost/enably_dev?sslmode=disable' up`
1. Run the backend with `go run cmd/enably/main.go`. It will be available on localhost:8080.
1. In another terminal, `cd frontent`, `npm install`, `npm run dev`.

 The frontent will be running on localhost:3000.

If you need to create new migrations, use the `migrate create -dir migrations -ext sql <migration_name>` command.

As long as you work on things in the roadmap, you should be fine, but create an issue just in case.


## Architecture notes:

The frontend is a typical Next (and hence React JS) app. This allows us to easily get server-rendering, which will hopefully improve SEO, and that's something we care about. The backend is written in Go.

Here are the most important backend packages:
- `store`, deals purely with communicating with the database and retrieving data from static files.
- `app`, contains the business logic, validation etc. Uses the store when necessary.
- `api`, accepts requests and returns responses. Delegates most of the work to `app`.
- `cmd/enably`, the app entry point, sets up the database, manually injects all the required dependencies and starts the HTTP server.

Packages don't depend on eachother directly, instead exposing interfaces for the lower layers to implement. This will allow swapping implementations for testing in the future.

### Schemas:

The schema for the available product categories and their required fields is stored in a file called schema.toml. This schema uses the concept of a fieldset, which is a group of fields that are required by many categories. For example, the ios_games category will require the "basic_app_info", "game_info" and "app_store_link" fieldsets. The backend converts this toml file into JSON schemas, which are used to validate products. This makes it easier to create forms in React, as there are libraries that can do it automatically, and to validate them in Go. In the future, it will also be possible to get nice diffs as products change. This design allows for quick modification of the schema without the need for a full GUI. The disadvantage is that it may be more difficult to filter products based on certain criteria, as products are stored as JSON rather than in separate columns.

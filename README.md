# Recipe Service
Repository for Recipe Service

Product Journey:
https://varsentinel.atlassian.net/wiki/spaces/Inventiums/overview

## Using Gin Framework:

https://gin-gonic.com/en/docs/quickstart/

## SQLC:

https://docs.sqlc.dev/en/stable/index.html

## Project Structures:
```
    recipe-service/
    ├── api/                   # Gin Router Controller
    │     ├── server.go        #
    ├── config/                # Store Applications Configs
    │     ├── config.go        #
    ├── handlers/              # Handlers Controller for different API Methods
    │     └── recipe.go        #
    ├── middlewares/           # Middlewares to check foe authorized client
    |     └── authenticate.go  #
    |── models/                # Models for working with Postgresl
    |     |── migration        # DB Migration
    |     |── query.           # DB Query
    |     |── sqlc             # DB Connection
    |── routes/                # Stores Route
    |     └── routes.go        #
```
## API Routes

- List Recipes:   `GET /v1/recipe/list`
- Get Recipe:     `GET /v1/recipe/:id`
- Create Recipe:  `POST /v1/recipe/create`
- Update Recipe:  `PUT /v1/recipe/:id`
- Delete Recipe:  `DELETE /v1/recipe/:id`

## Usage

How to perform db migration:

Prerequisites:
- Set $DB_SOURCE to the PostgreSQL URL

Run the following DB Migration Steps:
- For DB Migration Up
```
    $ make migrateup
```
- For DB Migration Down
```
    $ make migratedown
```
Run this command to generate sqlc code
```
    $ make sqlc
```

TestOps:

Uses the correct /v1/recipe endpoint
Removed authentication (since your Recipe Service doesn't have auth)
Cycles through 4 operations: create, get, list, and update
Uses proper form data parameters: name, ingredients, cost
Generates realistic test data with varying names, ingredients, and costs
Defaults to 10 requests with 0.5s delay between them

You can run it with: 
```sh
stimualte.sh
```
or customize with 
```sh
./script/stimualte.sh 20 1 
```

for 20 requests with 1s delay.

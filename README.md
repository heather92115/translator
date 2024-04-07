# Verdure Admin

The Verdure Admin application is the backend system to Administrate the 
Verdure application written in Rust. It is used to manage the vocabulary 
stored in a PostgreSQL relational database. 

# Building
To build the server:
> go build -o server ./cmd/server

### DB Connectivity
The following environment variables are required to connect to Postgres. The DB_LINK in
conjunction with the aws cli to obtain db connection properties.
> export DB_LINK="dev/aws/secret"
> 
> export REGION="us-east-1"

You can find docs here: [AWS Secret Manager](https://docs.aws.amazon.com/secretsmanager/latest/userguide/managing-secrets.html?icmpid=docs_asm_help_panel)


### GraphQL is used to access the system.


### Generate 
To generate updated models and resolvers:
> go run github.com/99designs/gqlgen generate

Then rebuild the server
> go build -v -o server ./cmd/server

To build everything:
>go build -v ./...

# To run tests:
>go test -v ./...
> 

# To run Integration Tests

### Add these tags to the top of each integration test:
//go:build integration
// +build integration

### Set the env var to find the .env.test
> export ENV_TEST_PATH="/Users/dev/go/src/verdure-admin/.env.test"

### Setup different db for your integration tests
> export DB_LINK="test/aws/secret"
>
> export REGION="us-east-1"

### To Run integration tests
> go test -tags=integration ./...



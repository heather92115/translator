# Translator

The Translator application is the backend system to Administrate the 
Palabras system. It is used to manage the vocabulary stored in a 
PostgreSQL relational database. 

# Building
To build the server:
> go build -o server ./cmd/server

### DB Connectivity
The following environment variables are required to connect to Postgres
>export DATABASE_IP="0.0.0.0"
>
>export DATABASE_PORT="5432"
> 
>export DATABASE_USER="dbuser"
> 
>export DATABASE_PASSWORD="****"
> 
>export DATABASE_NAME="postgres"

### GraphQL is used to access the system.


### Generate 
To generate updated models and resolvers:
> go run github.com/99designs/gqlgen generate

Then rebuild the server
> go build -o server ./cmd/server


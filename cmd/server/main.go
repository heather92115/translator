package main

import (
	"fmt"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/heather92115/translator/graph"
	"github.com/heather92115/translator/internal/db"
	"log"
	"net/http"
	"os"
)

const defaultPort = "8080"

func main() {
	fmt.Println("Starting the gql server")

	err := db.CreatePool()
	if err != nil {
		fmt.Printf("Failed DB connections, %v\n", err)
		return
	}

	port := os.Getenv("GQL_PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

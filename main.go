package main

import (
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/gorilla/websocket"
	"github.com/subosito/gotenv"
	"graphql_chat/graph"
	"graphql_chat/package/common"
	"graphql_chat/package/model"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

const defaultPort = "8080"

func main() {

	if err := gotenv.Load(); err != nil {
		log.Fatalf("error with init env variables: %s", err.Error())
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	db, err := common.InitPostgres()
	if err != nil {
		log.Fatalf("error with postrges: %v \n", err)
	}

	customCtx := &common.CustomContext{Database: db}

	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		MessageChanals: map[string]chan *model.Message{},
		Mutex:          &sync.Mutex{},
	}}))

	srv.AddTransport(transport.SSE{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	})

	srv.Use(extension.Introspection{})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", common.CreateContext(customCtx, srv))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

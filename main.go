package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"github.com/racho8/product-service/datastore"
	"github.com/racho8/product-service/handlers"
)

func main() {
	ctx := context.Background()

	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		slog.Error("GOOGLE_CLOUD_PROJECT env var not set")
		os.Exit(1)
	}

	dsClient, err := datastoreclient.NewClient(ctx, projectID)
	if err != nil {
		slog.Error("Failed to create Datastore client", slog.Any("error", err))
		os.Exit(1)
	}

	handlers.Init(dsClient)

	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Welcome to Product Service"))
	}).Methods("GET")

	r.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	r.HandleFunc("/products", handlers.CreateProduct).Methods("POST")
	r.HandleFunc("/products/{id}", handlers.GetProduct).Methods("GET")
	r.HandleFunc("/products/{id}", handlers.UpdateProduct).Methods("PUT")
	r.HandleFunc("/products/{id}", handlers.DeleteProduct).Methods("DELETE")
	r.HandleFunc("/products", handlers.ListProducts).Methods("GET")

	slog.Info("Server started on :8080")
	http.ListenAndServe(":8080", r)
}

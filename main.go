package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/gorilla/mux"
	"github.com/racho8/product-service/handlers"
)

func main() {
	ctx := context.Background()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		slog.Error("GOOGLE_CLOUD_PROJECT env var not set")
		os.Exit(1)
	}

	dsClient, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		slog.Error("Failed to create Datastore client", slog.Any("error", err))
		os.Exit(1)
	}
	defer dsClient.Close()

	handlers.Init(dsClient)

	r := mux.NewRouter()
	SetupRoutes(r)

	address := "0.0.0.0:" + port
	srv := &http.Server{
		Handler:      r,
		Addr:         address,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	slog.Info("Server starting", "address", address)
	if err := srv.ListenAndServe(); err != nil {
		slog.Error("Server failed", "error", err)
		os.Exit(1)
	}
}

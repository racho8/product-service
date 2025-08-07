package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/racho8/product-service/models"

	"cloud.google.com/go/datastore"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var dsClient *datastore.Client

func Init(client *datastore.Client) {
	dsClient = client
}

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	var product models.Product
	body, _ := io.ReadAll(r.Body)
	json.Unmarshal(body, &product)

	product.ID = uuid.New().String()

	key := datastore.NameKey("Product", product.ID, nil)
	if _, err := dsClient.Put(ctx, key, &product); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(product)
}

func GetProduct(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	id := mux.Vars(r)["id"]

	key := datastore.NameKey("Product", id, nil)
	var product models.Product
	if err := dsClient.Get(ctx, key, &product); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	product.ID = id
	json.NewEncoder(w).Encode(product)
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	id := mux.Vars(r)["id"]

	key := datastore.NameKey("Product", id, nil)
	var product models.Product
	body, _ := io.ReadAll(r.Body)
	json.Unmarshal(body, &product)

	product.ID = id
	if _, err := dsClient.Put(ctx, key, &product); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(product)
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	id := mux.Vars(r)["id"]

	key := datastore.NameKey("Product", id, nil)
	if err := dsClient.Delete(ctx, key); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func ListProducts(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	query := datastore.NewQuery("Product")

	var products []models.Product
	keys, err := dsClient.GetAll(ctx, query, &products)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for i, key := range keys {
		parts := strings.Split(key.Name, "/")
		products[i].ID = parts[len(parts)-1]
	}

	json.NewEncoder(w).Encode(products)
}

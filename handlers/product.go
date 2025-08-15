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

func CreateMultipleProducts(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// Parse the request body to get the list of products
	var request struct {
		Products []models.Product `json:"products"`
	}
	body, _ := io.ReadAll(r.Body)
	if err := json.Unmarshal(body, &request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Iterate over the products and save them to the datastore
	for i := range request.Products {
		request.Products[i].ID = uuid.New().String()
		key := datastore.NameKey("Product", request.Products[i].ID, nil)
		if _, err := dsClient.Put(ctx, key, &request.Products[i]); err != nil {
			http.Error(w, "Failed to create product: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(request.Products)
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

	// Fetch the existing product
	key := datastore.NameKey("Product", id, nil)
	var existingProduct models.Product
	if err := dsClient.Get(ctx, key, &existingProduct); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Decode the payload and update only the provided fields
	var updatedProduct models.Product
	body, _ := io.ReadAll(r.Body)
	json.Unmarshal(body, &updatedProduct)

	if updatedProduct.Name != "" {
		existingProduct.Name = updatedProduct.Name
	}
	if updatedProduct.Category != "" {
		existingProduct.Category = updatedProduct.Category
	}
	if updatedProduct.Segment != "" {
		existingProduct.Segment = updatedProduct.Segment
	}
	if updatedProduct.Price != 0 {
		existingProduct.Price = updatedProduct.Price
	}

	// Save the updated product back to the datastore
	if _, err := dsClient.Put(ctx, key, &existingProduct); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(existingProduct)
}

func UpdateMultipleProducts(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// Parse the request body to get the list of products
	var request struct {
		Products []models.Product `json:"products"`
	}
	body, _ := io.ReadAll(r.Body)
	if err := json.Unmarshal(body, &request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Iterate over the products and update them
	for _, updatedProduct := range request.Products {
		// Fetch the existing product
		key := datastore.NameKey("Product", updatedProduct.ID, nil)
		var existingProduct models.Product
		if err := dsClient.Get(ctx, key, &existingProduct); err != nil {
			http.Error(w, "Product not found: "+updatedProduct.ID, http.StatusNotFound)
			return
		}

		// Update only the provided fields
		if updatedProduct.Name != "" {
			existingProduct.Name = updatedProduct.Name
		}
		if updatedProduct.Category != "" {
			existingProduct.Category = updatedProduct.Category
		}
		if updatedProduct.Segment != "" {
			existingProduct.Segment = updatedProduct.Segment
		}
		if updatedProduct.Price != 0 {
			existingProduct.Price = updatedProduct.Price
		}

		// Save the updated product back to the datastore
		if _, err := dsClient.Put(ctx, key, &existingProduct); err != nil {
			http.Error(w, "Failed to update product: "+updatedProduct.ID, http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
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

func DeleteMultipleProducts(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// Parse the request body to get the list of IDs
	var request struct {
		IDs []string `json:"ids"`
	}
	body, _ := io.ReadAll(r.Body)
	if err := json.Unmarshal(body, &request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Create keys for the IDs and delete them
	keys := make([]*datastore.Key, len(request.IDs))
	for i, id := range request.IDs {
		keys[i] = datastore.NameKey("Product", id, nil)
	}

	if err := dsClient.DeleteMulti(ctx, keys); err != nil {
		http.Error(w, "Failed to delete products: "+err.Error(), http.StatusInternalServerError)
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

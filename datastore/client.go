package datastoreclient

import (
	"context"

	"cloud.google.com/go/datastore"
)

func NewClient(ctx context.Context, projectID string) (*datastore.Client, error) {
	return datastore.NewClient(ctx, projectID)
}

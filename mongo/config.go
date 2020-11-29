package mongo

import (
	"context"
	"os"
	"time"

	"github.com/soap-parser/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DB holds connection with db colletions
type DB struct {
	TransactionCollection TransactionCollection
}

// NewConnection starts the connection with database
func NewConnection() (*DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("LOTS_API_MONGO_CONN")))

	defer cancel()

	if err != nil {
		return nil, errors.ErrorGettingDBConnection(err)
	}

	database := client.Database(os.Getenv("LOTS_API_MONGO_DB"))
	if database == nil {
		return nil, errors.ErrorDBNotFound(os.Getenv("KOHO_API_MONGO_DB"))
	}

	transactionCol, err := NewTransactionCollection(ctx, database)
	if err != nil {
		return nil, err
	}

	return &DB{
		TransactionCollection: *transactionCol,
	}, nil
}

// Drop deleted the whole database
func Drop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("LOTS_API_MONGO_CONN")))

	defer cancel()

	if err != nil {
		return errors.ErrorGettingDBConnection(err)
	}

	database := client.Database(os.Getenv("LOTS_API_MONGO_DB"))
	if database == nil {
		return errors.ErrorDBNotFound(os.Getenv("LOTS_API_MONGO_DB"))
	}

	return database.Drop(ctx)
}

package mongo

import (
	"context"
	"time"

	"github.com/soap-parser/errors"
	"github.com/soap-parser/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const transactionCollection = "transactions"

// TransactionCollection represents the transaction collection
type TransactionCollection struct {
	access *mongo.Collection
}

// NewTransactionCollection returns the transaction collection access
func NewTransactionCollection(ctx context.Context, database *mongo.Database) (*TransactionCollection, error) {
	transactionCol := database.Collection(transactionCollection)
	if transactionCol == nil {
		return nil, errors.ErrorCollectionNotFound(transactionCollection)
	}

	_, err := transactionCol.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys: bson.M{
				"_id": 1,
			},
		},
	})
	if err != nil {
		return nil, errors.ErrorCreatingIndexes(err)
	}

	return &TransactionCollection{access: transactionCol}, nil
}

// Create creates a new transaction
func (ac TransactionCollection) Create(ctx context.Context, transaction *model.Transaction) (*model.Transaction, error) {
	if transaction == nil {
		return nil, errors.ErrorModelCannotBeNil(transactionCollection)
	}

	err := transaction.Validate()
	if err != nil {
		return nil, errors.ErrorValidating(transactionCollection, err)
	}

	transaction.Version = 1
	transaction.Schema = transaction.SchemaVersion()
	transaction.CreatedAt = time.Now()

	result, err := ac.access.InsertOne(ctx, transaction)
	if err != nil {
		return nil, errors.ErrorInserting(transactionCollection, err)
	}

	transaction.ID = result.InsertedID.(primitive.ObjectID)

	return transaction, nil
}

// Update updates an transaction
func (ac TransactionCollection) Update(ctx context.Context, transaction *model.Transaction) (*model.Transaction, error) {
	if transaction == nil {
		return nil, errors.ErrorModelCannotBeNil(transactionCollection)
	}

	err := transaction.Validate()
	if err != nil {
		return nil, errors.ErrorValidating(transactionCollection, err)
	}

	now := time.Now()
	transaction.Version = transaction.Version + 1
	transaction.UpdatedAt = &now

	filterVersion := bson.M{
		"_id":        transaction.ID,
		"version":    transaction.Version - 1,
		"deleted_at": bson.M{"$exists": false},
	}

	result, err := ac.access.UpdateOne(ctx, filterVersion, bson.M{"$set": transaction})
	if err != nil {
		return nil, errors.ErrorUpdating(transactionCollection, err)
	}

	if result.ModifiedCount == 0 {
		return nil, errors.ErrorUpdating(transactionCollection, errors.ErrorDocumentMismatch(transactionCollection, transaction.ID.Hex()))
	}

	return transaction, nil
}

// GetByID gets an transaction by id
func (ac TransactionCollection) GetByID(ctx context.Context, id string) (*model.Transaction, error) {
	itemID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.ErrorGetting(transactionCollection, errors.ErrorParsingObjectID(id))
	}

	filter := bson.M{
		"_id":        itemID,
		"deleted_at": bson.M{"$exists": false},
	}

	found := ac.access.FindOne(ctx, filter)
	result := new(model.Transaction)

	err = found.Decode(result)

	if err != nil && err.Error() != errors.NoDocumentsInResult().Error() {
		return nil, errors.ErrorGetting(transactionCollection, err)
	}

	return result, nil
}

// GetByTicketAndMatricula gets an transaction by ticket and matricula
func (ac TransactionCollection) GetByTicketAndMatricula(ctx context.Context, ticket string, parking int64, matricula string) (*model.Transaction, error) {
	filter := bson.M{
		"sequence":        ticket,
		"matricula":       matricula,
		"parking_info.id": parking,
		"deleted_at":      bson.M{"$exists": false},
	}

	found := ac.access.FindOne(ctx, filter)
	result := new(model.Transaction)

	err := found.Decode(result)

	if err != nil && err.Error() != errors.NoDocumentsInResult().Error() {
		return nil, errors.ErrorGetting(transactionCollection, err)
	}

	return result, nil
}

// GetAllByMatricula gets all transaction by matricula
func (ac TransactionCollection) GetAllByMatricula(ctx context.Context, parking int64, matricula string) ([]model.Transaction, error) {
	filter := bson.M{
		"matricula":       matricula,
		"parking_info.id": parking,
		"deleted_at":      bson.M{"$exists": false},
	}

	cursor, err := ac.access.Find(ctx, filter)

	if err != nil {
		return nil, errors.ErrorListing(transactionCollection, err)
	}

	var transactions []model.Transaction
	err = cursor.All(ctx, &transactions)
	if err != nil {
		return nil, errors.ErrorListing(transactionCollection, err)
	}

	return transactions, nil
}

// Delete deleted an transaction (logically)
func (ac TransactionCollection) Delete(ctx context.Context, id string) error {
	found, err := ac.GetByID(ctx, id)
	if err != nil {
		return errors.ErrorDeleting(transactionCollection, err)
	}

	if found.ID == primitive.NilObjectID {
		return nil
	}

	now := time.Now()
	found.DeletedAt = &now
	_, err = ac.Update(ctx, found)
	if err != nil {
		return errors.ErrorDeleting(transactionCollection, err)
	}

	return nil
}

// List returns a list of transactions
func (ac TransactionCollection) List(ctx context.Context, page int64, size int64, filters map[string]interface{}) ([]model.Transaction, error) {
	filter := bson.M{
		"deleted_at": bson.M{"$exists": false},
	}

	skip := (page - 1) * size
	cursor, err := ac.access.Find(ctx, filter, &options.FindOptions{Skip: &skip, Limit: &size})

	if err != nil {
		return nil, errors.ErrorListing(transactionCollection, err)
	}

	var transactions []model.Transaction
	err = cursor.All(ctx, &transactions)
	if err != nil {
		return nil, errors.ErrorListing(transactionCollection, err)
	}

	return transactions, nil
}

// Count returns total
func (ac TransactionCollection) Count(ctx context.Context, filters map[string]interface{}) (int64, error) {
	filter := bson.M{
		"deleted_at": bson.M{"$exists": false},
	}

	counter, err := ac.access.CountDocuments(ctx, filter)

	if err != nil {
		return 0, errors.ErrorCounting(transactionCollection, err)
	}

	return counter, nil
}

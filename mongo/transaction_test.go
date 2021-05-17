package mongo

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/csv-processor/errors"
	"github.com/csv-processor/model"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestTransaction_Create(t *testing.T) {
	db, err := StartDB(nil, nil)
	if err != nil {
		log.Panic(err)
	}

	type TestRun struct {
		name           string
		item           model.Transaction
		expectedItem   func(run *TestRun, saved *model.Transaction) *model.Transaction
		expectedError  error
		expectedResult *model.Transaction
	}

	tt := []TestRun{
		{
			name: "success",
			item: model.Transaction{},
			expectedItem: func(run *TestRun, saved *model.Transaction) *model.Transaction {
				copy := run.item

				copy.Version = saved.Version
				copy.Schema = saved.Schema
				copy.CreatedAt = saved.CreatedAt

				return &copy
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result, err := db.TransactionCollection.Create(context.Background(), &tc.item)

			if err != nil {
				require.Equal(t, tc.expectedError, err)
			} else {
				require.Equal(t, tc.expectedItem(&tc, result), result)
			}
		})
	}

	require.Nil(t, DropDB(nil, nil))
}

func TestTransaction_Update(t *testing.T) {
	db, err := StartDB(nil, nil)
	if err != nil {
		log.Panic(err)
	}

	type TestRun struct {
		name          string
		item          model.Transaction
		create        func(run *TestRun) *model.Transaction
		update        func(run *TestRun, saved *model.Transaction) *model.Transaction
		expectedItem  func(run *TestRun, updated *model.Transaction) *model.Transaction
		expectedError func(run *TestRun, updated *model.Transaction) error
	}

	tt := []TestRun{
		{
			name: "document mismatch version",
			item: model.Transaction{},
			create: func(run *TestRun) *model.Transaction {
				item := run.item
				saved, err := db.TransactionCollection.Create(context.Background(), &item)
				if err != nil {
					log.Panic(err)
				}
				return saved
			},
			update: func(run *TestRun, saved *model.Transaction) *model.Transaction {
				saved.Version = 0
				return saved
			},
			expectedError: func(run *TestRun, updated *model.Transaction) error {
				return errors.ErrorUpdating(transactionCollection, errors.ErrorDocumentMismatch(transactionCollection, updated.ID.Hex()))
			},
		},
		{
			name: "document mismatch id",
			item: model.Transaction{},
			create: func(run *TestRun) *model.Transaction {
				item := run.item
				saved, err := db.TransactionCollection.Create(context.Background(), &item)
				if err != nil {
					log.Panic(err)
				}
				return saved
			},
			update: func(run *TestRun, saved *model.Transaction) *model.Transaction {
				saved.ID = primitive.NewObjectID()
				return saved
			},
			expectedError: func(run *TestRun, updated *model.Transaction) error {
				return errors.ErrorUpdating(transactionCollection, errors.ErrorDocumentMismatch(transactionCollection, updated.ID.Hex()))
			},
		},
		{
			name: "document mismatch deleted",
			item: model.Transaction{},
			create: func(run *TestRun) *model.Transaction {
				itemDeleted := run.item
				now := time.Now()
				itemDeleted.DeletedAt = &now
				saved, err := db.TransactionCollection.Create(context.Background(), &itemDeleted)
				if err != nil {
					log.Panic(err)
				}
				return saved
			},
			update: func(run *TestRun, saved *model.Transaction) *model.Transaction {
				return saved
			},
			expectedError: func(run *TestRun, updated *model.Transaction) error {
				return errors.ErrorUpdating(transactionCollection, errors.ErrorDocumentMismatch(transactionCollection, updated.ID.Hex()))
			},
		},
		{
			name: "success",
			item: model.Transaction{},
			create: func(run *TestRun) *model.Transaction {
				item := run.item
				saved, err := db.TransactionCollection.Create(context.Background(), &item)
				if err != nil {
					log.Panic(err)
				}
				return saved
			},
			update: func(run *TestRun, saved *model.Transaction) *model.Transaction {
				return saved
			},
			expectedItem: func(run *TestRun, updated *model.Transaction) *model.Transaction {
				copy := run.item

				copy.ID = updated.ID
				copy.CreatedAt = updated.CreatedAt
				copy.UpdatedAt = updated.UpdatedAt
				copy.Schema = updated.Schema
				copy.Version = updated.Version
				return &copy
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			saved := tc.create(&tc)
			saved = tc.update(&tc, saved)
			result, err := db.TransactionCollection.Update(context.Background(), saved)

			if err != nil {
				require.Equal(t, tc.expectedError(&tc, saved), err)
			} else {
				require.Equal(t, tc.expectedItem(&tc, result), result)
			}
		})
	}

	require.Nil(t, DropDB(nil, nil))
}

func TestTransaction_GetByID(t *testing.T) {
	db, err := StartDB(nil, nil)
	if err != nil {
		log.Panic(err)
	}

	type TestRun struct {
		name          string
		itemID        func(saved *model.Transaction) string
		create        func(run *TestRun) *model.Transaction
		expectedError func(run *TestRun, saved *model.Transaction) error
		found         bool
	}

	tt := []TestRun{
		{
			name: "error parsing hex value",
			itemID: func(saved *model.Transaction) string {
				return "error"
			},
			create: func(run *TestRun) *model.Transaction {
				return nil
			},
			expectedError: func(run *TestRun, saved *model.Transaction) error {
				return errors.ErrorGetting(transactionCollection, errors.ErrorParsingObjectID(run.itemID(nil)))
			},
		},
		{
			name: "deleted",
			itemID: func(saved *model.Transaction) string {
				return saved.ID.Hex()
			},
			create: func(run *TestRun) *model.Transaction {
				now := time.Now()
				item := &model.Transaction{
					DeletedAt: &now,
				}

				item, err := db.TransactionCollection.Create(context.Background(), item)
				if err != nil {
					log.Panic(err)
				}

				return item
			},
		},
		{
			name: "success",
			itemID: func(saved *model.Transaction) string {
				return saved.ID.Hex()
			},
			create: func(run *TestRun) *model.Transaction {
				item := &model.Transaction{}

				item, err := db.TransactionCollection.Create(context.Background(), item)
				if err != nil {
					log.Panic(err)
				}

				return item
			},
			found: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			saved := tc.create(&tc)
			result, err := db.TransactionCollection.GetByID(context.Background(), tc.itemID(saved))

			if err != nil {
				require.Equal(t, tc.expectedError(&tc, saved), err)
			} else {
				if !tc.found {
					require.Equal(t, primitive.NilObjectID, result.ID)
				} else {
					require.Equal(t, tc.itemID(saved), result.ID.Hex())
				}
			}
		})
	}

	require.Nil(t, DropDB(nil, nil))
}

func TestTransaction_Delete(t *testing.T) {
	db, err := StartDB(nil, nil)
	if err != nil {
		log.Panic(err)
	}

	type TestRun struct {
		name          string
		itemID        func(saved *model.Transaction) string
		create        func(run *TestRun) *model.Transaction
		expectedError func(run *TestRun, saved *model.Transaction) error
	}

	tt := []TestRun{
		{
			name: "error parsing hex value",
			itemID: func(saved *model.Transaction) string {
				return "error"
			},
			create: func(run *TestRun) *model.Transaction {
				return nil
			},
			expectedError: func(run *TestRun, saved *model.Transaction) error {
				return errors.ErrorDeleting(transactionCollection, errors.ErrorGetting(transactionCollection, errors.ErrorParsingObjectID(run.itemID(nil))))
			},
		},
		{
			name: "success",
			itemID: func(saved *model.Transaction) string {
				return saved.ID.Hex()
			},
			create: func(run *TestRun) *model.Transaction {
				item := &model.Transaction{}

				item, err := db.TransactionCollection.Create(context.Background(), item)
				if err != nil {
					log.Panic(err)
				}

				return item
			},
			expectedError: func(run *TestRun, saved *model.Transaction) error {
				return nil
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			saved := tc.create(&tc)
			err := db.TransactionCollection.Delete(context.Background(), tc.itemID(saved))

			require.Equal(t, tc.expectedError(&tc, saved), err)
		})
	}

	require.Nil(t, DropDB(nil, nil))
}

func TestTransaction_List(t *testing.T) {
	db, err := StartDB(nil, nil)
	if err != nil {
		log.Panic(err)
	}

	type TestRun struct {
		name   string
		create func(run *TestRun)
		total  int
	}

	tt := []TestRun{
		{
			name:  "found only one",
			total: 1,
			create: func(run *TestRun) {
				item1 := &model.Transaction{}
				_, err := db.TransactionCollection.Create(context.Background(), item1)
				if err != nil {
					log.Panic(err)
				}

				now := time.Now()
				item2 := &model.Transaction{
					DeletedAt: &now,
				}
				_, err = db.TransactionCollection.Create(context.Background(), item2)
				if err != nil {
					log.Panic(err)
				}
			},
		},
		{
			name:  "found only one with pagination",
			total: 1,
			create: func(run *TestRun) {
				item1 := &model.Transaction{}
				_, err := db.TransactionCollection.Create(context.Background(), item1)
				if err != nil {
					log.Panic(err)
				}

				now := time.Now()
				item2 := &model.Transaction{
					DeletedAt: &now,
				}
				_, err = db.TransactionCollection.Create(context.Background(), item2)
				if err != nil {
					log.Panic(err)
				}
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.create(&tc)
			result, err := db.TransactionCollection.List(context.Background(), 1, 1, nil)

			require.Equal(t, nil, err)
			require.Equal(t, tc.total, len(result))
		})
	}

	require.Nil(t, DropDB(nil, nil))
}

func TestTransaction_Count(t *testing.T) {
	db, err := StartDB(nil, nil)
	if err != nil {
		log.Panic(err)
	}

	type TestRun struct {
		name   string
		create func(run *TestRun)
		total  int64
	}

	tt := []TestRun{
		{
			name:  "found only one",
			total: 1,
			create: func(run *TestRun) {
				item1 := &model.Transaction{}
				_, err := db.TransactionCollection.Create(context.Background(), item1)
				if err != nil {
					log.Panic(err)
				}

				now := time.Now()
				item2 := &model.Transaction{
					DeletedAt: &now,
				}
				_, err = db.TransactionCollection.Create(context.Background(), item2)
				if err != nil {
					log.Panic(err)
				}
			},
		},
		{
			name:  "found two",
			total: 2,
			create: func(run *TestRun) {
				item1 := &model.Transaction{}
				_, err := db.TransactionCollection.Create(context.Background(), item1)
				if err != nil {
					log.Panic(err)
				}

				now := time.Now()
				item2 := &model.Transaction{
					DeletedAt: &now,
				}
				_, err = db.TransactionCollection.Create(context.Background(), item2)
				if err != nil {
					log.Panic(err)
				}
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.create(&tc)
			result, err := db.TransactionCollection.Count(context.Background(), nil)

			require.Equal(t, nil, err)
			require.Equal(t, tc.total, result)
		})
	}

	require.Nil(t, DropDB(nil, nil))
}

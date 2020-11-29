package mongo

import (
	"fmt"
	"os"
	"testing"

	"github.com/soap-parser/errors"
	"github.com/stretchr/testify/require"
)

const (
	baseMongoConn = "mongodb://localhost:27017"
	baseMongoDB   = "test"
)

// StartDB starts the database
// 	by default, we use env variables, but for testing purposes
// 	we are adding some checks to force fail
func StartDB(connectionStr *string, dbName *string) (*DB, error) {
	if connectionStr != nil {
		os.Setenv("LOTS_API_MONGO_CONN", *connectionStr)
	} else {
		os.Setenv("LOTS_API_MONGO_CONN", baseMongoConn)
	}

	if dbName != nil {
		os.Setenv("LOTS_API_MONGO_DB", *dbName)
	} else {
		os.Setenv("LOTS_API_MONGO_DB", baseMongoDB)
	}

	return NewConnection()
}

// DropDB drops test DB
func DropDB(connectionStr *string, dbName *string) error {
	if connectionStr != nil {
		os.Setenv("LOTS_API_MONGO_CONN", *connectionStr)
	} else {
		os.Setenv("LOTS_API_MONGO_CONN", baseMongoConn)
	}

	if dbName != nil {
		os.Setenv("LOTS_API_MONGO_DB", *dbName)
	} else {
		os.Setenv("LOTS_API_MONGO_DB", baseMongoDB)
	}

	return Drop()
}

func Test_NewConnection(t *testing.T) {
	type TestRun struct {
		name          string
		expectedError error
		connectionStr func() *string
		dbName        func() *string
	}

	tt := []TestRun{
		{
			name:          "fail to connect with mongo",
			expectedError: errors.ErrorGettingDBConnection(fmt.Errorf("error parsing uri: scheme must be \"mongodb\" or \"mongodb+srv\"")),
			connectionStr: func() *string {
				str := "error"
				return &str
			},
			dbName: func() *string {
				str := "error"
				return &str
			},
		},
		{
			name: "success",
			connectionStr: func() *string {
				str := "mongodb://127.0.0.1:27017"
				return &str
			},
			dbName: func() *string {
				str := "seeker"
				return &str
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			_, err := StartDB(tc.connectionStr(), tc.dbName())

			require.Equal(t, tc.expectedError, err)
		})
	}
}

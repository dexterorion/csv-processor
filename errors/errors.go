package errors

import "fmt"

const (
	serviceNumber            = 1000
	errorGettingDBConnection = 1
	errorDBNotFound          = 2
	errorCollectionNotFound  = 3
	errorCreatingIndexes     = 4
	errorModelNil            = 5
	errorInserting           = 6
	errorUpdating            = 7
	errorDeleting            = 8
	errorListing             = 9
	errorGetting             = 10
	errorValidating          = 11
	errorOverlap             = 12
	errorCounting            = 13
)

// errorBase is the error base structure
func errorBase(errorNumber int, err error) error {
	return fmt.Errorf("Service: %d | : Error Number: %d | Description: %v", serviceNumber, errorNumber, err)
}

// ErrorGettingDBConnection returns and error when there is a problem connecting to DB
func ErrorGettingDBConnection(err error) error {
	return errorBase(errorGettingDBConnection, err)
}

// ErrorDBNotFound returns an error when the db is not found
func ErrorDBNotFound(dbName string) error {
	return errorBase(errorDBNotFound, fmt.Errorf("Database with name %s was not found", dbName))
}

// ErrorCollectionNotFound returns an error when the collection is not found
func ErrorCollectionNotFound(collectionName string) error {
	return errorBase(errorCollectionNotFound, fmt.Errorf("Collection with name %s was not found", collectionName))
}

// ErrorCreatingIndexes returns an error when the creating indexes
func ErrorCreatingIndexes(err error) error {
	return errorBase(errorCreatingIndexes, err)
}

// ErrorModelCannotBeNil returns an error when a model is nil
func ErrorModelCannotBeNil(modelName string) error {
	return errorBase(errorModelNil, fmt.Errorf("Model [%s] is nil", modelName))
}

// ErrorInserting returns an error when trying to insert to DB
func ErrorInserting(modelName string, err error) error {
	return errorBase(errorInserting, fmt.Errorf("Model [%s] got error [%v] on inserting", modelName, err))
}

// ErrorUpdating returns an error when trying to update to DB
func ErrorUpdating(modelName string, err error) error {
	return errorBase(errorUpdating, fmt.Errorf("Model [%s] got error [%v] on updating", modelName, err))
}

// ErrorDeleting returns an error when trying to delete to DB
func ErrorDeleting(modelName string, err error) error {
	return errorBase(errorDeleting, fmt.Errorf("Model [%s] got error [%v] on deleting", modelName, err))
}

// ErrorGetting returns an error when trying to get from DB
func ErrorGetting(modelName string, err error) error {
	return errorBase(errorGetting, fmt.Errorf("Model [%s] got error [%v] on getting", modelName, err))
}

// ErrorListing returns an error when trying to list from DB
func ErrorListing(modelName string, err error) error {
	return errorBase(errorListing, fmt.Errorf("Model [%s] got error [%v] on listing", modelName, err))
}

// ErrorCounting returns an error when trying to count from DB
func ErrorCounting(modelName string, err error) error {
	return errorBase(errorCounting, fmt.Errorf("Model [%s] got error [%v] on counting", modelName, err))
}

// ErrorCheckingOverlap returns an error when trying to check a overlap DB
func ErrorCheckingOverlap(modelName string, err error) error {
	return errorBase(errorOverlap, fmt.Errorf("Model [%s] got error [%v] on checking overlap", modelName, err))
}

// ErrorValidating returns an error when trying to validate model
func ErrorValidating(modelName string, err error) error {
	return errorBase(errorValidating, fmt.Errorf("Model [%s] got error [%v] on validating", modelName, err))
}

// ErrorDocumentMismatch returns an error when we don't find a document to update
func ErrorDocumentMismatch(modelName string, itemID string) error {
	return fmt.Errorf("Error updating [%s] with ID [%s]. Document mismatch", modelName, itemID)
}

// ErrorParsingObjectID returns an error trying to parse object id
func ErrorParsingObjectID(itemID string) error {
	return fmt.Errorf("Error parsing [%s] to objectID", itemID)
}

// NoDocumentsInResult returns an error when no documents are found
func NoDocumentsInResult() error {
	return fmt.Errorf("mongo: no documents in result")
}

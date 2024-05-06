package database

import (
	"api/database/helpers/upsert"
	"reflect"
	"time"

	"github.com/jinzhu/gorm"
)

const CREATED_AT_COLUMN_NAME string = "created_at"
const UPDATED_AT_COLUMN_NAME string = "updated_at"

func getColumnNamePlaceholderAndConflictHandler(columnName string, primaryKeyName string, isLastColumn bool) (string, string, string) {
	name := columnName
	valuePlaceholder := "?"
	conflictHandler := ""

	// Add a conflict handler if we have a primary key name and the column is not the primary key.
	// Also, skip the created_at column because we never want to update it.
	if primaryKeyName != "" && columnName != primaryKeyName && columnName != CREATED_AT_COLUMN_NAME {
		conflictHandler = columnName + "=excluded." + columnName
	}

	// If the column is not the last column, add a comma to separate the
	// column names and values.
	if !isLastColumn {
		if name != "" {
			name += ","
		}
		if valuePlaceholder != "" {
			valuePlaceholder += ","
		}
		if conflictHandler != "" {
			conflictHandler += ","
		}
	}

	return name, valuePlaceholder, conflictHandler
}

// Returns the column names, value placeholders, and conflict handlers for the query in a single loop.
func getColumnNamesPlaceholdersAndConflictHandlers(exampleRow interface{}, primaryKeyName string) (string, string, string) {
	var allColumnNames string
	var allColumnValuePlaceholders string
	var allColumnConflictHandlers string

	// Extract the column names from the first instance.
	columnNames := upsert.GetModelColumnNames(exampleRow)

	// Add the created_at and updated_at columns to the column names if the model has those fields.
	if upsert.HasCreatedAtField(exampleRow) {
		columnNames = append([]string{CREATED_AT_COLUMN_NAME}, columnNames...)
	}
	if upsert.HasUpdatedAtField(exampleRow) {
		columnNames = append([]string{UPDATED_AT_COLUMN_NAME}, columnNames...)
	}

	// Loop through the column names and values to build the query.
	for index, columnName := range columnNames {

		// Get the column name, value placeholder, and conflict handler.
		isLastColumn := index == len(columnNames)-1
		columnName, valuePlaceholder, conflictHandler := getColumnNamePlaceholderAndConflictHandler(columnName, primaryKeyName, isLastColumn)

		// Append the column name, value placeholder, and conflict handler
		// to the query.
		allColumnNames += columnName
		allColumnValuePlaceholders += valuePlaceholder
		allColumnConflictHandlers += conflictHandler
	}

	return allColumnNames, allColumnValuePlaceholders, allColumnConflictHandlers
}

func repeatColumnValuePlaceholders(columnValuePlaceholders string, rowAmount int) string {
	repeatedColumnValuePlaceholders := ""

	// Repeat the column value placeholders for each row.
	for index := 0; index < rowAmount; index++ {
		repeatedColumnValuePlaceholders += columnValuePlaceholders

		// Add a comma to separate the column value placeholders if it is not the last row.
		if index != rowAmount-1 {
			repeatedColumnValuePlaceholders += ","
		}
	}

	return repeatedColumnValuePlaceholders
}

func getRowValues(instances []interface{}) []interface{} {
	// Get the current time to use as the created_at and updated_at values.
	now := time.Now()

	rows := []interface{}{}
	for _, instance := range instances {
		// Get the column values from the instance.
		columnValues := []interface{}{}

		// Add created_at and updated_at values if the model has those fields.
		if upsert.HasCreatedAtField(instance) {
			columnValues = append(columnValues, now)
		}
		if upsert.HasUpdatedAtField(instance) {
			columnValues = append(columnValues, now)
		}

		// Append the column values to the row values.
		columnValues = append(columnValues, upsert.GetModelColumnValues(instance)...)

		// Append the column values to the row values.
		rows = append(rows, columnValues...)
	}

	return rows
}

func getModelInterface(modelType reflect.Type) reflect.Value {
	// Create a slice of pointers to the model type.
	sliceType := reflect.SliceOf(reflect.PointerTo(modelType))
	slice := reflect.MakeSlice(sliceType, 0, 0)

	// Return a pointer to the slice.
	return reflect.New(slice.Type())
}

func convertResultsToInterfaces(results reflect.Value) []interface{} {
	// Convert the results back to a slice of interfaces.
	upsertedRows := make([]interface{}, results.Len())

	for index := 0; index < results.Len(); index++ {
		result := results.Index(index)
		upsertedRows[index] = result.Elem().Interface()
	}

	return upsertedRows
}

func Upsert(database *gorm.DB, instances []interface{}) ([]interface{}, error) {
	rowAmount := len(instances)
	if rowAmount == 0 {
		return []interface{}{}, nil
	}

	// Get the table name of the instance by reflection.
	firstInstance := instances[0]
	modelType := reflect.TypeOf(firstInstance).Elem()
	tableName := upsert.PascalToSnake(modelType.Name()) + "s"

	// Get the primary key name of the instance.
	primaryKeyName := upsert.GetPrimaryKeyFieldName(modelType)

	// Get the column names, value placeholders, and conflict handlers.
	allColumnNames, allColumnValuePlaceholders, allConflictHandlers := getColumnNamesPlaceholdersAndConflictHandlers(firstInstance, primaryKeyName)

	var query string
	query += "INSERT INTO " + tableName + " (" + allColumnNames + ") "
	query += "VALUES " + repeatColumnValuePlaceholders("("+allColumnValuePlaceholders+")", rowAmount) + " "
	if primaryKeyName == "" {
		query += "ON CONFLICT DO NOTHING "
	} else {
		query += "ON CONFLICT (" + primaryKeyName + ") DO UPDATE SET " + allConflictHandlers + " "
	}
	query += "RETURNING *;"

	// Get a pointer with the correct model type to scan the results into.
	resultsPointer := getModelInterface(modelType)

	// Execute the query.
	// TODO: Use prepared statements to prevent SQL injection.
	database.Raw(query, getRowValues(instances)...).Scan(resultsPointer.Interface())

	// Convert the results back to a slice of interfaces.
	results := resultsPointer.Elem()
	upsertedRows := convertResultsToInterfaces(results)

	return upsertedRows, nil
}

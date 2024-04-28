package database_helpers

import (
	"errors"
	"reflect"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

func isPrimaryKey(field reflect.StructField) bool {
	gormTag := field.Tag.Get("gorm")
	return strings.Contains(gormTag, "primaryKey")
}

func getPrimaryKeyFieldName(modelType reflect.Type) string {
	// Loop through the fields of the model.
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)

		// Check each field to see if it's the primary key.
		if isPrimaryKey(field) {

			// Get the column name by converting the field name to
			// snake case.
			return PascalToSnake(field.Name)
		}
	}

	return ""
}

func getModelTypeAndValue(instance interface{}) (reflect.Type, reflect.Value) {
	// Get the type and value of the instance by reflection.
	modelType := reflect.TypeOf(instance)
	modelValue := reflect.ValueOf(instance)

	// If the instance is a pointer, get the type and value of the
	// element.
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
		modelValue = modelValue.Elem()
	}

	return modelType, modelValue
}

func isDatabaseModel(field reflect.StructField) bool {
	if field.Type.Kind() != reflect.Struct {
		return false
	}

	// Check if the field is a Gorm model.
	for index := 0; index < field.Type.NumField(); index++ {
		if field.Type.Field(index).Type == reflect.TypeOf(gorm.Model{}) {
			return true
		}
	}

	return false
}

func getColumnValues(instance interface{}) []interface{} {
	modelType, modelValue := getModelTypeAndValue(instance)

	var columnValues []interface{}

	for index := 0; index < modelType.NumField(); index++ {
		field := modelType.Field(index)

		// Skip unexported fields and database models because they
		// should be associated with other tables later.
		if field.PkgPath != "" || isDatabaseModel(field) {
			continue
		}

		fieldValue := modelValue.Field(index).Interface()

		// Slices  and Structs are most often used for many2many associations, so
		// we skip them for now.
		if field.Type.Kind() == reflect.Slice || field.Type.Kind() == reflect.Struct {
			continue
		}

		columnValues = append(columnValues, fieldValue)
	}

	return columnValues
}

func getColumnNames(instance interface{}) []string {
	modelType, _ := getModelTypeAndValue(instance)

	var columnNames []string

	for index := 0; index < modelType.NumField(); index++ {
		field := modelType.Field(index)

		// Skip unexported fields and database models because they
		// should be associated with other tables later.
		if field.PkgPath != "" {
			continue
		}

		fieldName := field.Name

		// Slices  and Structs are most often used for many2many associations, so
		// we skip them for now.
		if field.Type.Kind() == reflect.Slice || field.Type.Kind() == reflect.Struct {
			continue
		}

		// If it is not a struct, convert the field name to a snake
		// case column name.
		columnName := PascalToSnake(fieldName)
		columnNames = append(columnNames, columnName)
	}

	return columnNames
}

func getUpsertColumnParts(columnName string, primaryKeyName string, isLastColumn bool) (string, string, string) {
	name := columnName
	valuePlaceholder := "?"
	conflictHandler := ""

	// Add a conflict handler if the column is not the primary key.
	if columnName != primaryKeyName {
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

func getUpsertQueryParts(exampleRow interface{}, primaryKeyName string) (string, string, string) {
	var allColumnNames string
	var allColumnValuePlaceholders string
	var allColumnConflictHandlers string

	// Extract the column names from the first instance.
	columnNames := getColumnNames(exampleRow)

	// Add the created_at and updated_at columns to the column names.
	columnNames = append([]string{"created_at", "updated_at"}, columnNames...)

	// Loop through the column names and values to build the query.
	for index, columnName := range columnNames {

		// Get the column name, value placeholder, and conflict handler.
		isLastColumn := index == len(columnNames)-1
		columnName, valuePlaceholder, conflictHandler := getUpsertColumnParts(columnName, primaryKeyName, isLastColumn)

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
	for index := 0; index < rowAmount; index++ {
		repeatedColumnValuePlaceholders += columnValuePlaceholders
		if index != rowAmount-1 {
			repeatedColumnValuePlaceholders += ","
		}
	}

	return repeatedColumnValuePlaceholders
}

func getRowValues(instances []interface{}) []interface{} {
	// Also add the created_at and updated_at values to the column values.
	now := time.Now()
	sharedValues := []interface{}{now, now}

	rows := []interface{}{}
	for _, instance := range instances {
		// Get the column values from the instance.
		columnValues := sharedValues
		columnValues = append(columnValues, getColumnValues(instance)...)

		// Append the column values to the row values.
		rows = append(rows, columnValues...)
	}

	return rows
}

func linkAssociatedModel(database *gorm.DB, upsertedRow interface{}) error {
	modelType := reflect.TypeOf(upsertedRow)

	for index := 0; index < modelType.NumField(); index++ {
		field := modelType.Field(index)
		gormTag := field.Tag.Get("gorm")

		// Check if the field is a many2many association.
		if !strings.Contains(gormTag, "many2many") {
			continue
		}

		// Use Gorm to add the new associations.
		associationError := database.Model(upsertedRow).Association(field.Name).Append().Error
		if associationError != nil {
			return associationError
		}
	}

	return nil
}

func linkAssociatedModels(database *gorm.DB, upsertedRows []interface{}) error {
	for _, upsertedRow := range upsertedRows {
		// Skip nil rows as it's possible that the row was not changed.
		if upsertedRow == nil {
			continue
		}

		associationError := linkAssociatedModel(database, upsertedRow)
		if associationError != nil {
			return associationError
		}
	}

	return nil
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
	tableName := strings.ToLower(modelType.Name()) + "s"

	// Get the primary key name of the instance.
	primaryKeyName := getPrimaryKeyFieldName(modelType)
	if primaryKeyName == "" {
		return []interface{}{}, errors.New("could not find primary key")
	}

	// Get the column names, value placeholders, and conflict handlers.
	allColumnNames, allColumnValuePlaceholders, allConflictHandlers := getUpsertQueryParts(firstInstance, primaryKeyName)

	var query string
	query += "INSERT INTO " + tableName + " (" + allColumnNames + ") "
	query += "VALUES " + repeatColumnValuePlaceholders("("+allColumnValuePlaceholders+")", rowAmount) + " "
	query += "ON CONFLICT (" + primaryKeyName + ") DO UPDATE SET " + allConflictHandlers + " "
	query += "RETURNING *;"

	// Get a pointer with the correct model type to scan the results into.
	resultsPointer := getModelInterface(modelType)

	// Execute the query.
	// TODO: Use prepared statements to prevent SQL injection.
	database.Raw(query, getRowValues(instances)...).Scan(resultsPointer.Interface())

	// Convert the results back to a slice of interfaces.
	results := resultsPointer.Elem()
	upsertedRows := convertResultsToInterfaces(results)

	// Handle associations.
	associationError := linkAssociatedModels(database, upsertedRows)
	if associationError != nil {
		return []interface{}{}, associationError
	}

	return upsertedRows, nil
}

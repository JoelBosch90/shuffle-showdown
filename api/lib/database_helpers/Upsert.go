package database_helpers

import (
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

func isPrimaryKey(field reflect.StructField) bool {
	gormTag := field.Tag.Get("gorm")
	return strings.Contains(gormTag, "primary_key")
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

func getColumnValues(instance interface{}) []interface{} {
	modelType, modelValue := getModelTypeAndValue(instance)

	var columnValues []interface{}

	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)

		// Skip unexported fields.
		if field.PkgPath != "" {
			continue
		}

		fieldValue := modelValue.Field(i).Interface()

		// Check if the field is a struct.
		if modelType.Field(i).Type.Kind() == reflect.Struct {
			structValues := getColumnValues(fieldValue)
			columnValues = append(columnValues, structValues...)
			continue
		}

		columnValues = append(columnValues, fieldValue)
	}

	return columnValues
}

func getColumnNames(instance interface{}) []string {
	modelType, modelValue := getModelTypeAndValue(instance)

	var columnNames []string

	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)

		// Skip unexported fields.
		if field.PkgPath != "" {
			continue
		}

		fieldName := field.Name
		fieldValue := modelValue.Field(i).Interface()

		// Check if the field is a struct.
		if modelType.Field(i).Type.Kind() == reflect.Struct {
			structColumns := getColumnNames(fieldValue)
			columnNames = append(columnNames, structColumns...)
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
	var allColumnNames strings.Builder
	var allColumnValuePlaceholders strings.Builder
	var allColumnConflictHandlers strings.Builder

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
		allColumnNames.WriteString(columnName)
		allColumnValuePlaceholders.WriteString(valuePlaceholder)
		allColumnConflictHandlers.WriteString(conflictHandler)
	}

	return allColumnNames.String(), allColumnValuePlaceholders.String(), allColumnConflictHandlers.String()
}

func getRowValues(instances []interface{}) []interface{} {
	// Also add the created_at and updated_at values to the column values.
	now := time.Now()
	sharedValues := []interface{}{now, now}

	log.Println("instances", instances)

	rows := []interface{}{}
	for _, instance := range instances {
		// Get the column values from the instance.
		columnValues := sharedValues
		columnValues = append(columnValues, getColumnValues(instance)...)

		// Append the column values to the row values.
		rows = append(rows, columnValues...)
	}

	log.Println("rows", rows)

	return rows
}

func Upsert(database *gorm.DB, instances []interface{}) error {
	if len(instances) == 0 {
		return nil
	}

	// Get the table name of the instance by reflection.
	firstInstance := instances[0]
	modelType := reflect.TypeOf(firstInstance).Elem()
	tableName := strings.ToLower(modelType.Name()) + "s"
	primaryKeyName := getPrimaryKeyFieldName(modelType)

	// Get the column names, value placeholders, and conflict handlers.
	allColumnNames, allColumnValuePlaceholders, allConflictHandlers := getUpsertQueryParts(firstInstance, primaryKeyName)

	var query strings.Builder
	query.WriteString("INSERT INTO " + tableName + " (" + allColumnNames + ") ")
	query.WriteString("VALUES (" + allColumnValuePlaceholders + ") ")
	query.WriteString("ON CONFLICT (" + primaryKeyName + ") DO UPDATE SET " + allConflictHandlers + ";")

	database.Exec(query.String(), getRowValues(instances)...)

	return nil
}

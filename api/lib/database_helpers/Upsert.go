package database_helpers

import (
	"log"
	"reflect"
	"strings"

	"github.com/jinzhu/gorm"
)

func getPrimaryKeyFieldName(modelType reflect.Type) string {
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		gormTag := field.Tag.Get("gorm")
		if strings.Contains(gormTag, "primary_key") {
			return field.Tag.Get("json")
		}
	}

	return ""
}

func getColumns(instance interface{}) ([]string, []interface{}) {
	modelType := reflect.TypeOf(instance).Elem()
	modelValue := reflect.ValueOf(instance).Elem()

	var columnNames []string
	var columnValues []interface{}

	// TODO: clean up
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		gormTag := field.Tag.Get("gorm")
		if gormTag != "" {
			tagParts := strings.Split(gormTag, ";")
			for _, part := range tagParts {
				if strings.HasPrefix(part, "column:") {
					columnNames = append(columnNames, strings.TrimPrefix(part, "column:"))
					columnValues = append(columnValues, modelValue.Field(i).Interface())
					break
				}
			}
		}
	}

	return columnNames, columnValues
}

func Upsert(database *gorm.DB, instance interface{}) error {
	// Get the table name of the instance by reflection.
	modelType := reflect.TypeOf(instance).Elem()
	tableName := strings.ToLower(modelType.Name()) + "s"
	columnNames, columnValues := getColumns(instance)
	primaryKeyName := getPrimaryKeyFieldName(modelType)

	var allColumnNames strings.Builder
	var allValuePlaceholders strings.Builder
	var allConflictHandlers strings.Builder
	for index, columnName := range columnNames {
		allColumnNames.WriteString(columnName)
		allValuePlaceholders.WriteString("?")
		if columnName != primaryKeyName {
			allConflictHandlers.WriteString(columnName + "=excluded." + columnName)
		}
		if index != len(columnNames)-1 {
			allColumnNames.WriteString(",")
			allValuePlaceholders.WriteString(",")
			if columnName != primaryKeyName {
				allConflictHandlers.WriteString(",")
			}
		}
	}

	var query strings.Builder
	query.WriteString("INSERT INTO " + tableName + " (" + allColumnNames.String() + ") ")
	query.WriteString("VALUES (" + allValuePlaceholders.String() + ") ")
	// TODO: also set created_at and updated_at
	query.WriteString("ON CONFLICT (" + primaryKeyName + ") DO UPDATE SET " + allConflictHandlers.String() + ";")

	log.Println(query.String())

	database.Exec(query.String(), columnValues...)

	return nil
}

package database

import "reflect"

func GetModelColumnNames(instance interface{}) []string {
	var columnNames []string

	// Get the column names from the instance.
	IterateModelFields(instance, func(field reflect.StructField, fieldValue reflect.Value) {

		// If it is not a struct, convert the field name to a snake
		// case column name.
		columnName := PascalToSnake(field.Name)
		columnNames = append(columnNames, columnName)
	})

	return columnNames
}

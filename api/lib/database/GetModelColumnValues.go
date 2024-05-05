package database

import "reflect"

func GetModelColumnValues(instance interface{}) []interface{} {
	var columnValues []interface{}

	// Get the column values from the instance.
	IterateModelFields(instance, func(field reflect.StructField, fieldValue reflect.Value) {
		columnValues = append(columnValues, fieldValue.Interface())
	})

	return columnValues
}

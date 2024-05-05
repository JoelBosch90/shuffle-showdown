package database

import "reflect"

func GetPrimaryKeyFieldName(modelType reflect.Type) string {
	// Loop through the fields of the model.
	for index := 0; index < modelType.NumField(); index++ {
		field := modelType.Field(index)

		// Check each field to see if it's the primary key.
		if IsPrimaryKey(field) {

			// Get the column name by converting the field name to
			// snake case.
			return PascalToSnake(field.Name)
		}
	}

	return ""
}

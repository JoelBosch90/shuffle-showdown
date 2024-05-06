package upsert

import "reflect"

func IterateModelFields(instance interface{}, callback func(field reflect.StructField, fieldValue reflect.Value)) {
	modelType := GetModelType(instance)
	modelValue := GetModelValue(instance)

	for index := 0; index < modelType.NumField(); index++ {
		field := modelType.Field(index)

		// Skip unexported fields and database models because they
		// should be associated with other tables later.
		if field.PkgPath != "" || IsDatabaseModel(field) {
			continue
		}

		// Slices and Structs are used for many2many associations, so
		// we skip them for now.
		if field.Type.Kind() == reflect.Slice || field.Type.Kind() == reflect.Struct {
			continue
		}

		callback(field, modelValue.Field(index))
	}
}

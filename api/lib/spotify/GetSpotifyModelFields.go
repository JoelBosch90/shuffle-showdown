package spotify

import (
	"reflect"
	"strings"
)

func GetSpotifyModelFields(model interface{}) string {
	modelType := reflect.TypeOf(model)
	modelFields := make([]string, 0, modelType.NumField())

	for index := 0; index < modelType.NumField(); index++ {
		modelField := modelType.Field(index)
		modelFieldName := modelField.Tag.Get("json")

		switch modelField.Type.Kind() {
		case reflect.Struct:
			modelFields = append(modelFields, handleStructType(modelField, modelFieldName))
		case reflect.Slice:
			modelFields = append(modelFields, handleSliceType(modelField, modelFieldName))
		default:
			modelFields = append(modelFields, modelFieldName)
		}
	}

	return strings.Join(modelFields, ",")
}

func handleStructType(field reflect.StructField, fieldName string) string {
	nestedModelFields := GetSpotifyModelFields(reflect.New(field.Type).Elem().Interface())

	return fieldName + "(" + nestedModelFields + ")"
}

func handleSliceType(field reflect.StructField, fieldName string) string {
	sliceType := field.Type.Elem()

	if sliceType.Kind() == reflect.Struct {
		nestedModelFields := GetSpotifyModelFields(reflect.New(sliceType).Elem().Interface())
		return fieldName + "(" + nestedModelFields + ")"
	}

	return fieldName
}

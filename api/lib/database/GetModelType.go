package database

import "reflect"

func GetModelType(instance interface{}) reflect.Type {
	// Get the type of the instance by reflection.
	modelType := reflect.TypeOf(instance)

	// If the instance is a pointer, get the type of the
	// element.
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	return modelType
}

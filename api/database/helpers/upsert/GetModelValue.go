package upsert

import "reflect"

func GetModelValue(instance interface{}) reflect.Value {
	// Get the type and value of the instance by reflection.
	modelType := reflect.TypeOf(instance)
	modelValue := reflect.ValueOf(instance)

	// If the instance is a pointer, get the value of the
	// element.
	if modelType.Kind() == reflect.Ptr {
		modelValue = modelValue.Elem()
	}

	return modelValue
}

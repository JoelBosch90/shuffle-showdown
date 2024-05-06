package upsert

import (
	"reflect"

	"github.com/jinzhu/gorm"
)

func IsDatabaseModel(field reflect.StructField) bool {
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

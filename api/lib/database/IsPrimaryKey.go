package database

import (
	"reflect"
	"strings"
)

func IsPrimaryKey(field reflect.StructField) bool {
	return strings.Contains(field.Tag.Get("gorm"), "primaryKey")
}

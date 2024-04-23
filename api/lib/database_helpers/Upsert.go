package database_helpers

import (
	"reflect"

	"github.com/jinzhu/gorm"
)

func count(database *gorm.DB, newInstance interface{}) (uint, error) {
	modelValue := reflect.ValueOf(newInstance).Elem()
	primaryKeyValue := modelValue.FieldByName("ID").Interface()

	var count uint
	countError := database.Model(newInstance).Where("id = ?", primaryKeyValue).Count(&count).Error
	if countError != nil {
		return 0, countError
	}

	return count, nil
}

func update(database *gorm.DB, newInstance interface{}) error {
	modelValue := reflect.ValueOf(newInstance).Elem()
	primaryKeyValue := modelValue.FieldByName("ID").Interface()

	return database.Model(newInstance).Where("id = ?", primaryKeyValue).Omit("ID").Updates(newInstance).Error
}

func create(database *gorm.DB, newInstance interface{}) error {
	return database.Create(newInstance).Error
}

func Upsert(database *gorm.DB, newInstance interface{}) error {
	count, countError := count(database, newInstance)
	if countError != nil {
		return countError
	}

	// If the record exists, return early
	if count > 0 {
		return update(database, newInstance)
	} else {
		return create(database, newInstance)
	}
}

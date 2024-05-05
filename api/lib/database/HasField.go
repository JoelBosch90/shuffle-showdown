package database

func HasField(model interface{}, fieldName string) bool {
	modelType := GetModelType(model)

	for index := 0; index < modelType.NumField(); index++ {
		field := modelType.Field(index)

		if field.Name == fieldName {
			return true
		}
	}

	return false
}

func HasCreatedAtField(model interface{}) bool {
	return HasField(model, "CreatedAt")
}

func HasUpdatedAtField(model interface{}) bool {
	return HasField(model, "UpdatedAt")
}

package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	uuid "github.com/satori/go.uuid"
)

type OwnedEntity struct {
	BaseEntity
	ID uuid.UUID `json:"id" gorm:"type:uuid;primary_key;"`
}

func (ownedEntity *OwnedEntity) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("ID", uuid.NewV4())
}

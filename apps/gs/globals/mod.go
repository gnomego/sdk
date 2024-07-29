package globals

import (
	"github.com/gnomego/apps/gs/models"
	"github.com/gnomego/apps/gs/validation"
	"gorm.io/gorm"
)

var db *gorm.DB
var validiator *validation.GsValidator

func InitDb(d *gorm.DB, migrate bool) {
	db = d

	if migrate {
		db.AutoMigrate(
			models.OrgTable{},
			models.OrgDomainTable{})
	}
}

func InitValidator(v *validation.GsValidator) {
	validiator = v
}

func GetDb() *gorm.DB {
	if db == nil {
		panic("db not initialized")
	}

	return db
}

func GetValidator() *validation.GsValidator {
	if validiator == nil {
		panic("validator not initialized")
	}

	return validiator
}

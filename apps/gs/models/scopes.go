package models

import "gorm.io/gorm"

func StatusActive(db *gorm.DB) *gorm.DB {
	return db.Where("active = ?", int16(STATUS_ACTIVE))
}

func StatusInactive(db *gorm.DB) *gorm.DB {
	return db.Where("active = ?", int16(STATUS_INACTIVE))
}

func StatusDeleted(db *gorm.DB) *gorm.DB {
	return db.Where("deleted = ?", int16(STATUS_DELETED))
}

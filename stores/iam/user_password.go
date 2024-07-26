package iam

import (
	"database/sql"
	"strings"
	"time"

	"github.com/google/uuid"
)

type UserPasswordTable struct {
	UserId        int32        `gorm:"column:user_id;index:ix_user_passwords_user_id,unique" json:"userId"`
	Uid           uuid.UUID    `gorm:"column:uid;type:uuid;index:ix_user_passwords_uid,unique" json:"uid"`
	Password      string       `gorm:"column:password;size:1024" json:"password"`
	IsLocked      bool         `gorm:"column:is_locked" json:"isLocked"`
	LockedAt      sql.NullTime `gorm:"column:locked_at" json:"lockedAt"`
	FailureCount  int32        `gorm:"column:failed_attempts" json:"failedAttempts"`
	LastFailureAt sql.NullTime `gorm:"column:last_failed_at" json:"lastFailedAt"`
	SecurityStamp string       `gorm:"column:security_stamp;size:128" json:"securityStamp"`
	CreatedAt     time.Time    `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     sql.NullTime `gorm:"column:updated_at" json:"updated_at"`
}

func (UserPasswordTable) TableName() string {
	return "user_passwords"
}

func (up *UserPasswordTable) SetPassword(password string) error {
	p := password
	p = strings.TrimSpace(p)
	hash, err := HashSecret(p)
	if err != nil {
		return err
	}
	up.Password = hash
	return nil
}

func (up *UserPasswordTable) ValidatePassword(password string) error {
	return ValidateSecret(password, up.Password)
}

package iam

import (
	"database/sql"
	"strings"
	"time"
	"unicode"

	"github.com/dchest/uniuri"
	"github.com/google/uuid"
)

type UserApiKeyTable struct {
	Id            int32        `gorm:"column:id;primary_key;auto_increment" json:"id"`
	Uid           uuid.UUID    `gorm:"column:uid;type:uuid;index:ix_user_api_keys_uid,unique" json:"uid"`
	UserId        int32        `gorm:"column:user_id;index:ix_user_api_keys_user_id" json:"userId"`
	Key           string       `gorm:"column:key;size:2048" json:"key"`
	IsLocked      bool         `gorm:"column:is_locked" json:"isLocked"`
	LockedAt      sql.NullTime `gorm:"column:locked_at" json:"lockedAt"`
	FailureCount  int32        `gorm:"column:failed_attempts" json:"failedAttempts"`
	LastFailureAt sql.NullTime `gorm:"column:last_failed_at" json:"lastFailedAt"`
	ExpiresAt     sql.NullTime `gorm:"column:expires_at" json:"expiresAt"`
	CreatedAt     time.Time    `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     sql.NullTime `gorm:"column:updated_at" json:"updated_at"`
}

func (UserApiKeyTable) TableName() string {
	return "user_api_keys"
}

func (apiKey *UserApiKeyTable) Generate() (string, error) {
	key := uniuri.NewLenChars(20, []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-=#@~[]|"))
	valid := false
	for !valid {
		upper := false
		lower := false
		digit := false
		special := false
		for _, r := range key {
			if unicode.IsDigit(r) {
				digit = true
				continue
			}

			if unicode.IsLetter(r) {
				if unicode.IsUpper(r) {
					upper = true
					continue
				}

				if unicode.IsLower(r) {
					lower = true
					continue
				}
			}

			special = true
		}

		valid = upper && lower && digit && special
		if valid {
			break
		} else {
			key = uniuri.NewLenChars(20, []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-=#@~[]|"))
		}
	}

	err := apiKey.SetKey(key)
	if err != nil {
		return "", err
	}

	return key, nil
}

func (apiKey *UserApiKeyTable) SetKey(key string) error {
	k := key
	k = strings.TrimSpace(k)
	hash, err := HashSecret(k)
	if err != nil {
		return err
	}
	apiKey.Key = hash
	return nil
}

func (apiKey *UserApiKeyTable) IsExpired() bool {
	if apiKey.ExpiresAt.Valid {
		return apiKey.ExpiresAt.Time.Before(time.Now())
	}
	return false
}

func (apiKey *UserApiKeyTable) Reset() {
	apiKey.LockedAt = sql.NullTime{}
	apiKey.FailureCount = 0
	apiKey.LastFailureAt = sql.NullTime{}
	apiKey.IsLocked = false
}

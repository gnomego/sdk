package iam

import (
	"database/sql"
	"time"
)

type UserLoginProviderTable struct {
	UserId            int32          `gorm:"column:user_id;index:ix_user_logins_user_id" json:"userId"`
	Provider          string         `gorm:"column:provider;size:64;index:ix_user_logins_provider" json:"provider"`
	ProviderFormatted sql.NullString `gorm:"column:provider_formatted;size:64" json:"providerFormatted"`
	Key               string         `gorm:"column:key;size:128;index:ix_user_logins_key" json:"key"`
	CreatedAt         time.Time      `gorm:"column:created_at" json:"created_at"`
}

func (UserLoginProviderTable) TableName() string {
	return "user_login_providers"
}

type UserLoginTokenTable struct {
	UserId   int32  `gorm:"column:user_id;index:ix_user_login_tokens_user_id,unique" json:"userId"`
	Provider string `gorm:"column:provider;size:64,index:ix_user_login_tokens_provider,unique" json:"provider"`
	Name     string `gorm:"column:name;size:64,index:ix_user_login_tokens_name,unique" json:"name"`
	Token    string `gorm:"column:token;size:128" json:"token"`
}

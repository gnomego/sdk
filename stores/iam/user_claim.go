package iam

import (
	"time"

	"github.com/google/uuid"
)

type UserClaimTable struct {
	Id        int32     `gorm:"column:id,primary_key;auto_increment" json:"id"`
	Uid       uuid.UUID `gorm:"column:uid;type:uuid;index:ix_user_claims_uid,unique" json:"uid"`
	UserId    int32     `gorm:"column:user_id;index:ix_user_claims_user_id,unique" json:"userId"`
	Name      string    `gorm:"column:name;size:64,index:ix_user_claims_name,unique" json:"name"`
	Value     string    `gorm:"column:value;size:128" json:"value"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
}

func (UserClaimTable) TableName() string {
	return "user_claims"
}

package iam

import (
	"database/sql"

	"github.com/google/uuid"
)

type RoleTable struct {
	Id          int32         `gorm:"column:id;primary_key;auto_increment" json:"id"`
	Uid         uuid.UUID     `gorm:"column:uid;type:uuid;;index:ix_roles_uid,unique" json:"uid"`
	OrgId       sql.NullInt32 `gorm:"column:org_id" json:"organizationId"`
	Name        string        `gorm:"column:name;size:64;index:ix_roles_name,unique" json:"name"`
	Description string        `gorm:"column:description;size:256" json:"description"`
	Org         *OrgTable     `gorm:"foreignKey:OrgId;references:Id" json:"organization"`
}

func (RoleTable) TableName() string {
	return "roles"
}

type RoleClaimTable struct {
	Id        int32        `gorm:"column:id;primary_key,auto_increment" json:"id"`
	Uid       uuid.UUID    `gorm:"column:uid;type:uuid;index:ix_role_claims_uid,unique" json:"uid"`
	RoleId    int32        `gorm:"column:role_id;index:ix_role_claims_role_id,unique" json:"roleId"`
	Name      string       `gorm:"column:name;size:64,index:ix_role_claims_name,unique" json:"name"`
	Value     string       `gorm:"column:value;size:128" json:"value"`
	CreatedAt sql.NullTime `gorm:"column:created_at" json:"created_at"`
}

func (RoleClaimTable) TableName() string {
	return "role_claims"
}

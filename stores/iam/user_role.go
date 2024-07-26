package iam

type UserRoleTable struct {
	UserId int32 `gorm:"column:user_id;index:ix_users_roles_user_id,unique" json:"userId"`
	RoleId int32 `gorm:"column:role_id;index:ix_users_roles_role_id,unique" json:"roleId"`
}

func (UserRoleTable) TableName() string {
	return "users_roles"
}

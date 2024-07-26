package iam

func (db IamDb) AutoMigrateIam() error {
	return db.AutoMigrate(&UserTable{},
		&UserClaimTable{},
		&UserLoginProviderTable{},
		&UserPasswordTable{},
		&UserLoginTokenTable{},
		&UserClaimTable{},
		&UserApiKeyTable{},
		&RoleTable{},
		&RoleClaimTable{},
		&UserRoleTable{})
}

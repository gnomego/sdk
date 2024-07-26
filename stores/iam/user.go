package iam

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type UserTable struct {
	Id                  int32              `gorm:"column:id;primary_key;auto_increment" json:"id"`
	Uid                 uuid.UUID          `gorm:"column:uid;type:uuid;index:ix_users_uid,unique" json:"uid"`
	OrgId               sql.NullInt32      `gorm:"column:organization_id," json:"organizationId"`
	Name                string             `gorm:"column:name;size:64;index:ix_users_name,unique" json:"name" validate:"required"`
	NameFormatted       sql.NullString     `gorm:"column:name_formatted;size:64" json:"name_formatted"`
	Email               string             `gorm:"column:email;size:128;index:ix_users_email,unique" json:"email" validate:"required,email"`
	EmailFormatted      sql.NullString     `gorm:"column:email_formatted;size:128" json:"emailFormatted"`
	EmailVerified       bool               `gorm:"column:email_verified" json:"emailVerified"`
	PhoneNumber         sql.NullString     `gorm:"column:phone_number;size:16" json:"phoneNumber"`
	PhoneNumberVerified bool               `gorm:"column:phone_number_verified" json:"phoneNumberVerified"`
	AvatarUrl           sql.NullString     `gorm:"column:avatar_url;size:1048" json:"avatarUrl"`
	Status              int8               `gorm:"column:status" json:"status"`
	LastLoginAt         sql.NullTime       `gorm:"column:last_login_at" json:"lastLoginAt"`
	LastLoginIp         sql.NullString     `gorm:"column:last_login_ip;size:45" json:"lastLoginIp"`
	ConcurrencyStamp    string             `gorm:"column:concurrency_stamp;size:128" json:"concurrencyStamp"`
	CreatedAt           time.Time          `gorm:"column:created_at" json:"created_at"`
	UpdatedAt           sql.NullTime       `gorm:"column:updated_at" json:"updated_at"`
	Organization        *OrgTable          `gorm:"foreignKey:OrganizationId" json:"organization"`
	Password            *UserPasswordTable `gorm:"foreignKey:UserId;references:Id" json:"password"`
	Claims              []UserClaimTable   `gorm:"foreignKey:UserId;references:Id" json:"claims"`
	ApiKeys             []UserApiKeyTable  `gorm:"foreignKey:UserId;references:Id" json:"apiKeys"`
}

func (UserTable) TableName() string {
	return "users"
}

func (db *IamDb) NewUserWithPassword(name string, email string, password string) (*UserTable, error) {
	user, err := db.newUser(name, email)
	if err != nil {
		return nil, err
	}

	err = user.Validate()
	if err != nil {
		return nil, err
	}

	user.Password = &UserPasswordTable{
		UserId: user.Id,
		Uid:    uuid.New(),
	}
	err = user.Password.SetPassword(password)
	if err != nil {
		return nil, err
	}

	err = db.Save(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (db *IamDb) newUser(name string, email string) (*UserTable, error) {
	e := strings.TrimSpace(email)
	e = strings.ToLower(e)
	n := strings.TrimSpace(name)
	n = strings.ToLower(n)

	var count int64

	db.DB.Model(&UserTable{}).Where("email = ?", e).Or("name = ?", n).Count(&count)
	if count > 0 {
		return nil, fmt.Errorf("a user with the email or name already exists")
	}

	user := UserTable{}
	user.SetEmail(email)
	user.SetName(name)
	user.Uid = uuid.New()
	user.Status = 1
	user.CreatedAt = time.Now()
	user.ConcurrencyStamp = uuid.NewString()

	return &user, nil
}

func (db *IamDb) NewUser(name string, email string) (*UserTable, error) {
	user, err := db.newUser(name, email)
	if err != nil {
		return nil, err
	}

	err = user.Validate()
	if err != nil {
		return nil, err
	}

	err = db.Save(&user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (user *UserTable) Save(db *IamDb) error {
	err := user.Validate()
	if err != nil {
		return err
	}

	err = db.DB.Save(&user).Error
	return err
}

func (user *UserTable) Validate() error {
	if user == nil {
		return nil
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(*user)
	return err
}

func (user *UserTable) LockPassword() *UserTable {
	if user.Status == 2 {
		return user
	}

	user.ConcurrencyStamp = uuid.NewString()

	if user.Password != nil {
		user.Password.IsLocked = true
	}

	user.Status = 2
	return user
}

func (user *UserTable) UnlockPassword() *UserTable {
	if user.Status == 1 {
		return user
	}

	user.ConcurrencyStamp = uuid.NewString()

	if user.Password != nil {
		user.Password.IsLocked = false
	}

	user.Status = 1
	return user
}

func (user *UserTable) SetLastLogin(ip *string) *UserTable {
	user.LastLoginAt = sql.NullTime{Time: time.Now(), Valid: true}
	if ip != nil {
		user.LastLoginIp = sql.NullString{String: *ip, Valid: true}
	}

	user.ConcurrencyStamp = uuid.NewString()
	return user
}

func (user *UserTable) SetActive() *UserTable {
	user.Status = 1
	return user
}

func (user *UserTable) SetInactive() *UserTable {
	user.Status = 0
	return user
}

func (user *UserTable) SetEmail(email string) *UserTable {
	e := email
	e = strings.TrimSpace(e)
	slug := strings.ToLower(e)

	if user.Email != slug {
		user.Email = slug
		println("setting email", slug)
		user.ConcurrencyStamp = uuid.NewString()
	}

	user.Email = slug
	if e != slug || (user.EmailFormatted.Valid && user.EmailFormatted.String != e) {
		user.EmailFormatted.String = e
		user.ConcurrencyStamp = uuid.NewString()
	}

	return user
}

func (user *UserTable) SetName(name string) *UserTable {
	n := name
	n = strings.TrimSpace(n)
	slug := strings.ToLower(n)

	if user.Name != slug {
		user.Name = slug
		println("setting name", slug)
		user.ConcurrencyStamp = uuid.NewString()
	}

	if n != slug || user.NameFormatted.Valid && user.NameFormatted.String != n {
		user.NameFormatted.String = n
		user.ConcurrencyStamp = uuid.NewString()
	}

	return user
}

func (db *IamDb) GetUserByUid(uid uuid.UUID) (*UserTable, error) {
	var user UserTable
	err := db.DB.Where("uid = ?", uid).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (db *IamDb) GetUserById(id int32) (*UserTable, error) {
	var user UserTable
	err := db.DB.Model(UserTable{}).Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

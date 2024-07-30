package models

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	Id                  int32             `gorm:"column:id;primaryKey;autoIncrement"`
	Uid                 uuid.UUID         `gorm:"column:uid;type:uuid;not null;index:idx_users_uid,unique"`
	OrgId               sql.NullInt32     `gorm:"column:org_id"`
	Email               string            `gorm:"column:email;type:string;not null;size:256;index:idx_users_email,unique"`
	EmailFormatted      sql.NullString    `gorm:"column:email_formatted;size:256;"`
	EmailVerified       bool              `gorm:"column:email_verified"`
	Name                string            `gorm:"column:name;size:64;not null"`
	NameFormatted       sql.NullString    `gorm:"column:name_formatted;size:64;"`
	DisplayName         sql.NullString    `gorm:"column:display_name;size:64;"`
	AvatarUrl           sql.NullString    `gorm:"column:avatar_url;size:1024;"`
	PhoneNumber         sql.NullString    `gorm:"column:phone_number;size:16;"`
	PhoneNumberVerified bool              `gorm:"column:phone_number_verified"`
	Status              int8              `gorm:"column:status" json:"status"`
	LastLoginAt         sql.NullTime      `gorm:"column:last_login_at" json:"lastLoginAt"`
	LastLoginIp         sql.NullString    `gorm:"column:last_login_ip;size:45" json:"lastLoginIp"`
	ConcurrencyStamp    string            `gorm:"column:concurrency_stamp;size:128" json:"concurrencyStamp"`
	CreatedAt           time.Time         `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt           time.Time         `gorm:"column:updated_at;autoUpdateTime"`
	Organization        *OrgTable         `gorm:"foreignKey:OrgId;references:Id"`
	PasswordAuth        *UserPasswordAuth `gorm:"foreignKey:UserId;references:Id"`
}

type UserPasswordAuth struct {
	UserId          int32        `gorm:"column:id;primaryKey"`
	Password        string       `gorm:"column:password;type:string;not null;size:256"`
	ExpiresAt       sql.NullTime `gorm:"column:expires_at"`
	FailedAttempts  int16        `gorm:"column:failed_attempts"`
	FailedAttemptAt sql.NullTime `gorm:"column:failed_attempt_at"`
	LockedAt        sql.NullTime `gorm:"column:locked_at"`
	IsLocked        bool         `gorm:"column:is_locked"`
}

type UserRepo struct {
	db      *gorm.DB
	Options *UserRepoOptions
}

type UserRepoOptions struct {
	LockWindow      time.Duration
	AllowedAttempts int16
}

func (u *User) TableName() string {
	return "users"
}

func (u *User) GetName() string {
	if u.NameFormatted.Valid && u.NameFormatted.String != "" {
		return u.NameFormatted.String
	}

	return u.Name
}

func (u *User) GetEmail() string {
	if u.EmailFormatted.Valid && u.EmailFormatted.String != "" {
		return u.EmailFormatted.String
	}

	return u.Email
}

func (u *User) SetEmail(email string) *User {
	email = strings.TrimSpace(email)
	e := strings.ToLower(email)

	u.Email = e

	if (!u.EmailFormatted.Valid && email != e) || (u.EmailFormatted.Valid && u.EmailFormatted.String != e) {
		u.EmailFormatted = sql.NullString{
			String: email,
			Valid:  true,
		}
	}

	return u
}

func (u *User) SetStatus(status int8) *User {
	u.Status = status
	return u
}

func (u *User) SetName(name string) *User {
	name = strings.TrimSpace(name)
	n := strings.ToLower(name)

	u.Name = n

	if (!u.NameFormatted.Valid && name != n) || (u.NameFormatted.Valid && u.NameFormatted.String != n) {
		u.NameFormatted = sql.NullString{
			String: name,
			Valid:  true,
		}
	}

	return u
}

func (upa *UserPasswordAuth) TableName() string {
	return "user_password_auth"
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.Uid = uuid.New()
	u.CreatedAt = time.Now().UTC()
	u.UpdatedAt = time.Now().UTC()
	return
}

func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	u.UpdatedAt = time.Now().UTC()
	return
}

func (upa *UserPasswordAuth) BeforeCreate(tx *gorm.DB) (err error) {
	hasher := DefaultSecretHasher()

	if upa.Password != "" {
		upa.Password, err = hasher.HashSecretString(upa.Password)
		if err != nil {
			return err
		}
	}

	return nil
}

func (upa *UserPasswordAuth) Authenticate(password string) bool {
	if upa.IsLocked {
		return false
	}

	hasher := DefaultSecretHasher()
	err := hasher.ValidateSecretString(password, upa.Password)
	return err == nil
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{
		db: db,
		Options: &UserRepoOptions{
			LockWindow:      30 * time.Minute,
			AllowedAttempts: 5,
		},
	}
}

func (r *UserRepo) FindById(id int32, expand ...string) (*User, error) {
	var user User

	query := r.preload(expand)
	err := query.First(&user).Error
	return &user, err
}

func (r *UserRepo) FindByUid(uid uuid.UUID, expand ...string) (*User, error) {
	var user User

	query := r.preload(expand)
	err := query.Where("uid = ?", uid).First(&user).Error
	return &user, err
}

func (r *UserRepo) FindByUidString(uid string, expand ...string) (*User, error) {
	id, err := uuid.Parse(uid)
	if err != nil {
		return nil, fmt.Errorf("invalid uid %s", uid)
	}

	return r.FindByUid(id, expand...)
}

func (r *UserRepo) FindByEmail(email string, expand ...string) (*User, error) {
	var user User
	e := strings.TrimSpace(email)
	e = strings.ToLower(e)

	query := r.preload(expand)
	err := query.Where("email = ?", e).First(&user).Error
	return &user, err
}

func (r *UserRepo) AuthenticatePassword(email, password string) (bool, *User, error) {
	var user User
	e := strings.TrimSpace(email)
	e = strings.ToLower(e)
	err := r.db.Where("email = ?", e).Preload("PasswordAuth").First(&user).Error
	if err != nil {
		return false, &user, err
	}

	if user.PasswordAuth == nil {
		return false, &user, fmt.Errorf("user does not exist or password is incorrect")
	}

	if user.PasswordAuth.IsLocked {
		if user.PasswordAuth.LockedAt.Time.Add(r.Options.LockWindow).Before(time.Now().UTC()) {
			user.PasswordAuth.LockedAt = sql.NullTime{
				Valid: false,
			}

			user.PasswordAuth.FailedAttempts = 0
			user.PasswordAuth.FailedAttemptAt = sql.NullTime{
				Valid: false,
			}

			user.PasswordAuth.IsLocked = false
			user.Status = STATUS_ACTIVE

			err := r.db.Save(&user).Error
			if err != nil {
				return false, nil, err
			}
		}
	}

	if !user.PasswordAuth.Authenticate(password) {
		user.PasswordAuth.FailedAttemptAt = sql.NullTime{
			Time:  time.Now().UTC(),
			Valid: true,
		}

		user.PasswordAuth.FailedAttempts++
		if user.PasswordAuth.FailedAttempts >= r.Options.AllowedAttempts {
			user.PasswordAuth.LockedAt = sql.NullTime{
				Time:  time.Now().UTC(),
				Valid: true,
			}

			user.PasswordAuth.IsLocked = true
			user.Status = STATUS_PASSWORD_LOCKED
		}

		err := r.db.Save(&user).Error
		if err != nil {
			return false, &user, err
		}

		return false, &user, fmt.Errorf("user account is locked. please try again later")
	}

	return false, &user, nil
}

func (r *UserRepo) Create(user *User) error {
	return r.db.Create(user).Error
}

func (r *UserRepo) preload(expand []string) *gorm.DB {
	query := r.db.Model(&User{})
	for _, e := range expand {
		switch e {
		case "organization":
		case "org":
			query = query.Preload("Organization")
		case "password":
		case "passwordauth":
			query = query.Preload("PasswordAuth")
		}
	}

	return query
}

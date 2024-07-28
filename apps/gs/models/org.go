package models

import (
	"database/sql"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/gobuffalo/flect"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Org struct {
	Id      string   `json:"id" validate:"required"`
	Name    string   `json:"name" validate:"required,max=64"`
	Slug    string   `json:"slug" validate:"required,max=64"`
	Status  *string  `json:"status,omitempty"`
	IsRoot  bool     `json:"root"`
	Domains []string `json:"domains,omitempty"`
}

type OrgSelect struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type OrgTable struct {
	Id            int32            `gorm:"column:id;primaryKey;autoIncrement"`
	Uid           uuid.UUID        `gorm:"column:uid;type:uuid;not null;index:idx_orgs_uid,unique"`
	Name          string           `gorm:"column:name;size:64;not null;index:idx_orgs_name,unique"`
	NameFormatted sql.NullString   `gorm:"column:name_formatted;size:64;"`
	Slug          string           `gorm:"column:slug;type:string;not null;size:64;index:idx_orgs_slug,unique"`
	Status        int16            `gorm:"column:status;index:idx_orgs_status"`
	IsRoot        bool             `gorm:"column:is_root"`
	CreatedAt     time.Time        `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt     time.Time        `gorm:"column:updated_at;autoUpdateTime"`
	Domains       []OrgDomainTable `gorm:"foreignKey:OrgId;references:Id"`
}

func (o *OrgTable) TableName() string {
	return "orgs"
}

func (ot *OrgTable) GetDomains() []string {
	domains := []string{}
	for _, d := range ot.Domains {
		domains = append(domains, d.Domain)
	}

	return domains
}

func (ot *OrgTable) BeforeCreate(tx *gorm.DB) (err error) {
	ot.Uid = uuid.New()
	ot.Status = 1
	ot.CreatedAt = time.Now().UTC()
	ot.UpdatedAt = time.Now().UTC()

	return
}

func (ot *OrgTable) SetName(name string) *OrgTable {
	n := strings.TrimSpace(name)
	n = strings.ToLower(n)

	if n != ot.Name {
		ot.Name = n
	}

	if (!ot.NameFormatted.Valid && n != name) || (ot.NameFormatted.Valid && ot.NameFormatted.String != name) {
		ot.NameFormatted = sql.NullString{
			String: name,
			Valid:  true,
		}
	}

	if ot.Slug == "" {
		ot.Slug = flect.Dasherize(n)
	}

	return ot
}

func (ot *OrgTable) SetSlug(slug string) *OrgTable {
	ot.Slug = slug

	ot.Slug = slug

	return ot
}

type OrgRepo struct {
	db *gorm.DB
}

func NewOrgRepo(db *gorm.DB) *OrgRepo {
	return &OrgRepo{
		db: db,
	}
}

func (s *OrgRepo) Create(org *OrgTable) error {
	return s.db.Create(org).Error
}

func (s *OrgRepo) Count() (int64, error) {
	var count int64
	err := s.db.Model(&OrgTable{}).Count(&count).Error
	return count, err
}

func (o *OrgRepo) DeleteById(id int32) error {
	return o.db.Delete(OrgTable{}, id).Error
}

func (o *OrgRepo) DeleteByUid(uid uuid.UUID) error {
	return o.db.Delete(OrgTable{}, "uid = ?", uid).Error
}

func (os *OrgRepo) FindByUid(uid uuid.UUID, expand ...string) (*OrgTable, error) {
	org := &OrgTable{}

	tx := os.db.Model(org)

	if len(expand) > 0 && slices.Contains(expand, "domains") {
		tx = tx.Preload("Domains")
	}

	err := tx.Where("uid = ?", uid).First(org).Error
	return org, err
}

func (os *OrgRepo) FindBySlug(slug string, expand ...string) (*OrgTable, error) {
	org := &OrgTable{}

	tx := os.db.Model(org)

	if len(expand) > 0 && slices.Contains(expand, "domains") {
		tx = tx.Preload("Domains")
	}

	err := tx.Where("slug = ?", slug).First(org).Error
	return org, err
}

func (os *OrgRepo) FindByName(name string, expand ...string) (*OrgTable, error) {
	org := &OrgTable{}

	tx := os.db.Model(org)

	name = strings.TrimSpace(name)
	name = strings.ToLower(name)

	if len(expand) > 0 && slices.Contains(expand, "domains") {
		tx = tx.Preload("Domains")
	}

	r := tx.Where("name = ?", name).First(org)
	if r.Error != nil {
		return nil, r.Error
	}

	if r.RowsAffected == 0 {
		return nil, nil
	}

	return org, nil
}

func (os *OrgRepo) UpdateName(uid string, name string) error {
	id, err := uuid.Parse(uid)
	if err != nil {
		return fmt.Errorf("invalid uid: %s", uid)
	}

	n := strings.TrimSpace(name)
	n = strings.ToLower(n)

	values := map[string]interface{}{
		"name": n,
	}

	if n != name {
		values["name_formatted"] = name
	}

	return os.db.Model(OrgTable{}).Where("uid = ?", id).Updates(values).Error
}

func (os *OrgRepo) UpdateSlug(uid string, slug string) error {
	id, err := uuid.Parse(uid)
	if err != nil {
		return fmt.Errorf("invalid uid: %s", uid)
	}

	return os.db.Model(OrgTable{}).Where("uid = ?", id).Update("slug", slug).Error
}

func (os *OrgRepo) Update(org *OrgTable) error {
	return os.db.Save(org).Error
}

func (os *OrgRepo) Delete(org *OrgTable) error {
	return os.db.Delete(org).Error
}

func (os *OrgRepo) All() ([]OrgTable, error) {
	orgs := []OrgTable{}
	err := os.db.Find(&orgs).Error
	return orgs, err
}

func (os *OrgRepo) AllByStatus(status int16) ([]OrgTable, error) {
	orgs := []OrgTable{}
	err := os.db.Where("status = ?", status).Find(&orgs).Error
	return orgs, err
}

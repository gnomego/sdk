package iam

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrgPair struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Org struct {
	Id      string   `json:"id"`
	Name    string   `json:"name" validate:"required,max=64"`
	Slug    string   `json:"slug" validate:"max=64"`
	Domains []string `json:"domains"`
}

func (org *Org) Validate(v *validator.Validate) error {
	if v == nil {
		return validator.New(validator.WithRequiredStructEnabled()).Struct(org)
	}

	return v.Struct(org)
}

type OrgStore struct {
	db        *gorm.DB
	Validator *validator.Validate
}

func NewOrgStore(db *gorm.DB, v *validator.Validate) *OrgStore {
	if v == nil {
		v = validator.New()
	}

	return &OrgStore{
		db:        db,
		Validator: v,
	}
}

func (store *OrgStore) DeleteByUid(id string) (int64, error) {

	err := uuid.Validate(id)

	if err != nil {
		return 0, fmt.Errorf("id must be a valid UUID")
	}

	res := store.db.Delete(&OrgTable{}, "uid = ?", id)
	if res.Error != nil {
		return 0, res.Error
	}

	return res.RowsAffected, nil
}

func (store *OrgStore) DeleteById(id int32) (int64, error) {

	res := store.db.Delete(&OrgTable{}, "id = ?", id)
	if res.Error != nil {
		return 0, res.Error
	}

	return res.RowsAffected, nil
}

func (store *OrgStore) Create(org *Org) error {

	table := OrgTable{}
	_, err := org.ToOrgTable(&table)
	if err != nil {
		return err
	}

	println(table.Name)
	println(table.Slug)
	err = store.db.Create(&table).Error
	if err != nil {
		return err
	}

	org.Id = table.Uid.String()

	return nil
}

func (store *OrgStore) CountAll() (int64, error) {
	var count int64
	res := store.db.Model(&OrgTable{}).Count(&count)
	if res.Error != nil {
		return 0, res.Error
	}

	return count, nil
}

func (store *OrgStore) UpdateName(uid string, name string) error {

	id, err := uuid.Parse(uid)
	if err != nil {
		return fmt.Errorf("id must be a valid UUID")
	}

	n := strings.TrimSpace(name)
	n = strings.ToLower(n)

	values := map[string]interface{}{
		"name": n,
	}

	if n != name {
		values["name_formatted"] = name
	}

	res := store.db.Model(&OrgTable{}).Where("uid = ?", id).Updates(values)

	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (store *OrgStore) UpdateSlug(uid string, slug string) error {

	id, err := uuid.Parse(uid)
	if err != nil {
		return fmt.Errorf("id must be a valid UUID")
	}

	res := store.db.Model(&OrgTable{}).Where("uid = ?", id).Update("slug", slug)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (store *OrgStore) AddDomain(uid string, domain string) error {
	id, err := uuid.Parse(uid)
	if err != nil {
		return fmt.Errorf("id must be a valid UUID")
	}

	table := OrgTable{}
	r := store.db.Model(&OrgTable{}).Preload("Domains").Where("uid = ?", id).First(&table)
	if r.Error != nil {
		return r.Error
	}

	if r.RowsAffected == 0 {
		return errors.New("org not found")
	}

	d := strings.TrimSpace(domain)
	d = strings.ToLower(d)

	if len(table.Domains) > 0 {
		for _, domain := range table.Domains {
			if domain.Domain == d {
				return nil
			}
		}
	}

	table.Domains = append(table.Domains, OrgDomain{Domain: d})
	err = store.db.Save(&table).Error
	if err != nil {
		return err
	}

	return nil
}

func (store *OrgStore) RemoveDomain(uid string, domain string) error {
	id, err := uuid.Parse(uid)
	if err != nil {
		return fmt.Errorf("id must be a valid UUID")
	}

	table := OrgTable{}
	r := store.db.Model(&OrgTable{}).Preload("Domains").Where("uid = ?", id).First(&table)
	if r.Error != nil {
		return r.Error
	}

	if r.RowsAffected == 0 {
		return errors.New("org not found")
	}

	d := strings.TrimSpace(domain)
	d = strings.ToLower(d)

	found := false
	if len(table.Domains) > 0 {
		for i, domain := range table.Domains {
			if domain.Domain == d {
				table.Domains = append(table.Domains[:i], table.Domains[i+1:]...)
				found = true
				break
			}
		}
	}

	if found {
		err = store.db.Save(&table).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func (store *OrgStore) Save(org *Org) error {
	table := OrgTable{}
	res := store.db.Find(&table, "uid = ?", org.Id)
	if res.Error != nil {
		return res.Error
	}

	_, err := org.ToOrgTable(&table)
	if err != nil {
		return err
	}

	err = store.db.Save(&table).Error
	if err != nil {
		return err
	}

	org.Id = table.Uid.String()

	return nil
}

func (store *OrgStore) AllPairs() ([]OrgPair, error) {
	tables := []OrgTable{}
	res := store.db.Find(&tables)
	if res.Error != nil {
		return nil, res.Error
	}

	orgs := make([]OrgPair, 0)
	for _, table := range tables {
		orgs = append(orgs, table.ToOrgPair())
	}

	return orgs, nil
}

func (store *OrgStore) PagePairs(page int, size int) ([]OrgPair, error) {
	tables := []OrgTable{}
	res := store.db.Offset(page * size).Limit(size).Find(&tables)
	if res.Error != nil {
		return nil, res.Error
	}

	orgs := make([]OrgPair, 0)
	for _, table := range tables {
		orgs = append(orgs, table.ToOrgPair())
	}

	return orgs, nil
}

func (store *OrgStore) All(expand ...string) ([]Org, error) {
	tables := []OrgTable{}

	if len(expand) > 0 && slices.Contains(expand, "domains") {
		res := store.db.Joins("Domain").Find(&tables)
		if res.Error != nil {
			return nil, res.Error
		}
	} else {
		res := store.db.Find(&tables)
		if res.Error != nil {
			return nil, res.Error
		}
	}

	orgs := make([]Org, 0)
	for _, table := range tables {
		orgs = append(orgs, table.ToOrg())
	}

	return orgs, nil
}

func (store *OrgStore) Page(page int, size int, expand bool) ([]Org, error) {
	tables := []OrgTable{}

	if expand {
		res := store.db.Joins("Domain").Offset(page * size).Limit(size).Find(&tables)
		if res.Error != nil {
			return nil, res.Error
		}
	} else {
		res := store.db.Offset(page * size).Limit(size).Find(&tables)
		if res.Error != nil {
			return nil, res.Error
		}
	}

	orgs := make([]Org, 0)
	for _, table := range tables {
		orgs = append(orgs, table.ToOrg())
	}

	return orgs, nil
}

func (store *OrgStore) FindByUid(id string, expand ...string) (*Org, error) {

	err := uuid.Validate(id)

	if err != nil {
		return nil, fmt.Errorf("id must be a valid UUID")
	}

	table := OrgTable{}
	if len(expand) > 0 && slices.Contains(expand, "domains") {
		res := store.db.Model(&OrgTable{}).Preload("Domains").First(&table, "uid = ?", id)
		if res.Error != nil {
			return nil, res.Error
		}

		if res.RowsAffected == 0 {
			return nil, errors.New("org not found")
		}

	} else {
		res := store.db.First(&table, "uid = ?", id)
		if res.Error != nil {
			return nil, res.Error
		}

		if res.RowsAffected == 0 {
			return nil, errors.New("org not found")
		}
	}

	org := table.ToOrg()

	return &org, nil
}

func (store *OrgStore) FindById(id int32, expand ...string) (*Org, error) {

	table := OrgTable{}
	if len(expand) > 0 && slices.Contains(expand, "domains") {
		res := store.db.Model(&OrgTable{}).Preload("Domains").First(&table, "id = ?", id)
		if res.Error != nil {
			return nil, res.Error
		}

		if res.RowsAffected == 0 {
			return nil, errors.New("org not found")
		}

	} else {
		res := store.db.First(&table, "id = ?", id)
		if res.Error != nil {
			return nil, res.Error
		}

		if res.RowsAffected == 0 {
			return nil, errors.New("org not found")
		}
	}

	org := table.ToOrg()

	return &org, nil
}

func (store *OrgStore) FindBySlug(slug string, expand ...string) (*Org, error) {

	table := OrgTable{}
	if len(expand) > 0 && slices.Contains(expand, "domains") {
		res := store.db.Model(&OrgTable{}).Preload("Domains").First(&table, "slug = ?", slug)
		if res.Error != nil {
			return nil, res.Error
		}

		if res.RowsAffected == 0 {
			return nil, errors.New("org not found")
		}

	} else {
		res := store.db.First(&table, "slug = ?", slug)
		if res.Error != nil {
			return nil, res.Error
		}

		if res.RowsAffected == 0 {
			return nil, errors.New("org not found")
		}
	}

	org := table.ToOrg()

	return &org, nil
}

func (store *OrgStore) FindByName(name string, expand ...string) (*Org, error) {

	name = strings.TrimSpace(name)
	name = strings.ToLower(name)
	table := OrgTable{}
	if len(expand) > 0 && slices.Contains(expand, "domains") {
		res := store.db.Model(&OrgTable{}).Preload("Domains").First(&table, "name = ?", name)
		if res.Error != nil {
			return nil, res.Error
		}

		if res.RowsAffected == 0 {
			return nil, errors.New("org not found")
		}

	} else {
		res := store.db.First(&table, "name = ?", name)
		if res.Error != nil {
			return nil, res.Error
		}

		if res.RowsAffected == 0 {
			return nil, errors.New("org not found")
		}
	}

	org := table.ToOrg()

	return &org, nil
}

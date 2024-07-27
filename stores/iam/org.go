package iam

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/gobuffalo/flect"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrgTable struct {
	Id            int32          `gorm:"column:id;primaryKey,autoIncrement"`
	Uid           uuid.UUID      `gorm:"column:uid;type:uuid;index:ix_orgs_uid,unique"`
	Name          string         `gorm:"column:name;size:64,index:ix_orgs_name,unique"`
	NameFormatted sql.NullString `gorm:"column:name_formatted;size:64"`
	Slug          string         `gorm:"column:slug;size:64,index:ix_orgs_slug,unique"`
	Domains       []OrgDomain    `gorm:"foreignKey:OrgId;references:Id"`
}

func (OrgTable) TableName() string {
	return "orgs"
}

func (org *OrgTable) SetName(name string) *OrgTable {
	n := strings.TrimSpace(name)
	n = strings.ToLower(n)

	org.Name = n

	if (!org.NameFormatted.Valid && name != n) || org.NameFormatted.String != name {
		println("setting name formatted", name, n)
		org.NameFormatted = sql.NullString{
			String: name,
			Valid:  true,
		}
	}

	if org.Slug == "" {
		s := flect.Dasherize(n)
		org.Slug = s
	}

	return org
}

func (org *OrgTable) SetSlug(slug string) *OrgTable {
	s := strings.TrimSpace(slug)
	s = flect.Dasherize(s)

	if org.Slug != s {
		org.Slug = s
	}

	return org
}

func (org *OrgTable) AddDomain(domain string) *OrgTable {
	d := strings.TrimSpace(domain)
	d = strings.ToLower(d)

	for _, dom := range org.Domains {
		if dom.Domain == d {
			return org
		}
	}

	org.Domains = append(org.Domains, OrgDomain{
		OrgId:  org.Id,
		Domain: d,
	})

	return org
}

func (org *OrgTable) RemoveDomain(domain string) *OrgTable {
	d := strings.TrimSpace(domain)
	d = strings.ToLower(d)

	for i, dom := range org.Domains {
		if dom.Domain == d {
			org.Domains = append(org.Domains[:i], org.Domains[i+1:]...)
			break
		}
	}

	return org
}

func (org *OrgTable) SetDomains(domains []string) *OrgTable {
	org.Domains = []OrgDomain{}

	for _, domain := range domains {
		org.AddDomain(domain)
	}

	return org
}

func (org *OrgTable) ToOrgPair() OrgPair {
	return OrgPair{
		Id:   org.Uid.String(),
		Name: org.Name,
	}
}

func (org *OrgTable) ToOrg() Org {
	domains := []string{}
	for _, domain := range org.Domains {
		domains = append(domains, domain.Domain)
	}

	n := org.Name
	if org.NameFormatted.Valid {
		n = org.NameFormatted.String
	}

	return Org{
		Id:      org.Uid.String(),
		Name:    n,
		Slug:    org.Slug,
		Domains: domains,
	}
}

func (org *Org) ToOrgTable(table *OrgTable) (*OrgTable, error) {

	id := uuid.New()

	if org.Id != "" {
		id2, err := uuid.Parse(org.Id)
		if err != nil {
			return nil, fmt.Errorf("id must be a valid UUID")
		}

		id = id2
	}

	if table == nil {
		table = &OrgTable{}
	}

	table.Uid = id
	table.SetName(org.Name)
	if org.Slug != "" {
		table.SetSlug(org.Slug)
	}

	if org.Domains != nil && len(org.Domains) > 0 {
		for _, domain := range org.Domains {
			d := strings.TrimSpace(domain)
			d = strings.ToLower(d)
			table.Domains = append(table.Domains, OrgDomain{Domain: d})
		}
	}

	return table, nil
}

func (org *OrgTable) BeforeCreate(db *gorm.DB) error {
	org.Uid = uuid.New()

	return nil
}

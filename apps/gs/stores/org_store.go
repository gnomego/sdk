package stores

import (
	"github.com/gnomego/apps/gs/einfo"
	"github.com/gnomego/apps/gs/globals"
	"github.com/gnomego/apps/gs/log"
	"github.com/gnomego/apps/gs/models"
)

type OrgStore struct {
	repo *models.OrgRepo
}

type NewOrg struct {
	Name    string   `json:"name" validate:"required,max=64"`
	Slug    *string  `json:"slug" validate:"max=64"`
	Domains []string `json:"domains" validate:"dive,domain,max=256"`
}

type Org struct {
	Id      string   `json:"id"`
	Name    string   `json:"name"`
	Slug    string   `json:"slug"`
	Status  *string  `json:"status,omitempty"`
	IsRoot  bool     `json:"root"`
	Domains []string `json:"domains,omitempty"`
}

func (o *NewOrg) Validate() *einfo.ErrorInfo {
	v := globals.GetValidator()
	err := v.Struct(o)
	if err == nil {
		return nil
	}

	return v.TranslateStruct(o, err)
}

func NewOrgStore() *OrgStore {
	db := globals.GetDb()
	repo := models.NewOrgRepo(db)
	return &OrgStore{
		repo: repo,
	}
}

func (s *OrgStore) Create(org *NewOrg) *Response[*Org] {
	response := &Response[*Org]{
		Ok: true,
	}

	e := org.Validate()
	if e != nil {
		response.Ok = false
		response.Error = e
		return response
	}

	table, err := s.repo.FindByName(org.Name)

	if err != nil {
		log.Error(err, "error finding org by name %s", org.Name)
		response.Ok = false
		response.Error = einfo.Sprintf("failed to find org by name: %v", org.Name)
		return response
	}

	if table != nil {
		response.Ok = false
		response.Error = einfo.Sprintf("org with name %s already exists", org.Name)
		return response
	}

	table = &models.OrgTable{}
	if org.Slug != nil {
		table.SetSlug(*org.Slug)
	}

	table.SetName(org.Name)

	for _, domain := range org.Domains {
		table.Domains = append(table.Domains, models.OrgDomainTable{
			Domain: domain,
		})
	}

	err = s.repo.Create(table)
	if err != nil {
		log.Error(err, "error creating org %v", org.Name)
		response.Ok = false
		response.Error = einfo.Sprintf("failed to create org: %v", org.Name)
		return response
	}

	status := mapToStatus(table.Status)

	response.Value = &Org{
		Id:      table.Uid.String(),
		Name:    table.Name,
		Slug:    table.Slug,
		Status:  &status,
		IsRoot:  table.IsRoot,
		Domains: table.GetDomains(),
	}

	return response
}

func mapToStatus(s int16) string {
	switch s {
	case 1:
		return "active"
	case 2:
		return "inactive"
	case 4:
		return "deleted"
	default:
		return "active"
	}
}

package stores

import (
	"github.com/gnomego/apps/gs/einfo"
	"github.com/gnomego/apps/gs/globals"
	"github.com/gnomego/apps/gs/models"
	"github.com/gnomego/apps/gs/xgin"
	"github.com/google/uuid"
)

type OrgStore struct {
	repo *models.OrgRepo
}

type NewOrg struct {
	Name    string    `json:"name" validate:"required,max=64"`
	Slug    *string   `json:"slug" validate:"max=64"`
	Domains *[]string `json:"domains" validate:"dive,domain,max=256"`
}

type Org struct {
	Id      string   `json:"id"`
	Name    string   `json:"name"`
	Slug    string   `json:"slug"`
	Status  *string  `json:"status,omitempty"`
	IsRoot  bool     `json:"root"`
	Domains []string `json:"domains,omitempty"`
}

func (o *Org) Validate() *einfo.ErrorInfo {
	v := globals.GetValidator()
	err := v.Struct(o)
	if err == nil {
		return nil
	}

	return v.TranslateStructWithName(o, "Org", err)
}

func (o *NewOrg) Validate() *einfo.ErrorInfo {
	v := globals.GetValidator()
	err := v.Struct(o)
	if err == nil {
		return nil
	}

	return v.TranslateStructWithName(o, "NewOrg", err)
}

func NewOrgStore() *OrgStore {
	db := globals.GetDb()
	repo := models.NewOrgRepo(db)
	return &OrgStore{
		repo: repo,
	}
}

// All godoc
// @Summary gets all orgs
// @Schemes
// @Description gets all orgs
// @Tags orgs
// @Accept json
// @Produce json
// @Success 200 {object} xgin.Response[[]Org] ok
// @Failure 500 {object} xgin.Response[[]Org] error
// @Router / [get]
func (s *OrgStore) All() *xgin.Response[*[]Org] {
	response := &xgin.Response[*[]Org]{
		Ok: true,
	}

	tables, err := s.repo.All()
	if err != nil {
		return response.SetError(err, "error getting all orgs")
	}

	orgs := []Org{}
	for _, table := range tables {
		status := mapToStatus(table.Status)
		org := &Org{
			Id:     table.Uid.String(),
			Name:   table.Name,
			Slug:   table.Slug,
			Status: &status,
			IsRoot: table.IsRoot,
		}

		orgs = append(orgs, *org)
	}

	response.Value = &orgs

	return response
}

func (s *OrgStore) Create(org *NewOrg) *xgin.Response[*Org] {
	response := &xgin.Response[*Org]{
		Value: nil,
		Ok:    true,
	}

	e := org.Validate()
	if e != nil {
		return response.Invalid(e)
	}

	exists, err := s.repo.ExistsByName(org.Name)
	if err != nil {
		return response.SetError(err, "ExistsByName failed: %v", org.Name)
	}

	if exists {
		return response.SetErrorMessage("org_exists", "org with name already exists")
	}

	table := &models.OrgTable{}
	if org.Slug != nil {
		table.SetSlug(*org.Slug)
	}

	table.SetName(org.Name)

	if org.Domains != nil {
		for _, domain := range *org.Domains {
			table.Domains = append(table.Domains, models.OrgDomainTable{
				Domain: domain,
			})
		}
	}

	err = s.repo.Create(table)
	if err != nil {
		return response.SetError(err, "error creating org %s", org.Name)
	}

	response.Value = mapOrg(table)

	return response
}

func (s *OrgStore) Save(org *Org) *xgin.Response[*Org] {
	response := &xgin.Response[*Org]{
		Value: nil,
		Ok:    true,
	}

	e := org.Validate()
	if e != nil {
		return response.Invalid(e)
	}

	uid, err := uuid.Parse(org.Id)
	if err != nil {
		return response.SetError(err, "invalid org id")
	}

	table, err := s.repo.FindByUid(uid, "domains")
	if err != nil {
		return response.SetError(err, "error finding org %s", org.Id)
	}

	name := table.Name
	if table.NameFormatted.Valid {
		name = table.NameFormatted.String
	}
	if org.Name != name {
		table.SetName(org.Name)
	}

	if org.Slug != table.Slug {
		table.SetSlug(org.Slug)
	}

	if (org.Status != nil) && (*org.Status != mapToStatus(table.Status)) {
		table.Status = mapFromStatus(*org.Status)
	}

	err = s.repo.Update(table)
	if err != nil {
		return response.SetError(err, "error updating org %s", org.Id)
	}

	response.Value = mapOrg(table)
	return response
}

func mapOrg(table *models.OrgTable) *Org {
	status := mapToStatus(table.Status)
	name := table.Name
	if table.NameFormatted.Valid {
		name = table.NameFormatted.String
	}
	org := &Org{
		Id:      table.Uid.String(),
		Name:    name,
		Slug:    table.Slug,
		Status:  &status,
		IsRoot:  table.IsRoot,
		Domains: table.GetDomains(),
	}

	return org
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

func mapFromStatus(s string) int16 {
	switch s {
	case "active":
		return 1
	case "inactive":
		return 2
	case "deleted":
		return 4
	default:
		return 1
	}
}

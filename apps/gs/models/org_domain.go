package models

type OrgDomainTable struct {
	Id     int32  `gorm:"column:id;primaryKey;autoIncrement"`
	OrgId  int32  `gorm:"column:org_id;index:idx_org_domains_org_id"`
	Domain string `gorm:"column:domain;size:256;not null;index:idx_org_domains_domain,unique"`
}

func (od *OrgDomainTable) TableName() string {
	return "org_domains"
}

func (od *OrgDomainTable) SetDomain(domain string) *OrgDomainTable {
	od.Domain = domain
	return od
}

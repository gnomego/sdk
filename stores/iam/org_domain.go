package iam

type OrgDomain struct {
	Id     int32  `gorm:"column:id;primaryKey;autoIncrement"`
	OrgId  int32  `gorm:"column:org_id"`
	Domain string `gorm:"column:domain;size:128;index:ix_org_domains_domain,unique"`
}

func (OrgDomain) TableName() string {
	return "org_domains"
}

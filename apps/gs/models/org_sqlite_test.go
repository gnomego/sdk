package models_test

import (
	"testing"

	"github.com/gnomego/apps/gs/models"
	assert2 "github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	orgDb *gorm.DB
)

func init() {
	if orgDb == nil {
		db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
		if err != nil {
			panic(err)
		}

		err = db.AutoMigrate(&models.OrgTable{}, &models.OrgDomainTable{})
		if err != nil {
			panic(err)
		}

		orgDb = db
	}
}

func TestOrgStore(t *testing.T) {
	assert := assert2.New(t)

	repo := models.NewOrgRepo(orgDb)

	org1 := &models.OrgTable{
		Slug: "test-org",
		Domains: []models.OrgDomainTable{
			{
				Domain: "test.org",
			},
		},
	}

	org1.SetName("Test Org")

	org2 := &models.OrgTable{
		Slug: "test-org-2",
		Domains: []models.OrgDomainTable{
			{
				Domain: "test2.org",
			},
		},
	}
	org2.SetName("Test Org 2")

	err := repo.Create(org1)
	if err != nil {
		t.Fatalf("failed to create org: %v", err)
	}

	err = repo.Create(org2)
	if err != nil {
		t.Fatalf("failed to create org: %v", err)
	}

	one, err := repo.FindByUid(org1.Uid)
	if err != nil {
		t.Fatalf("failed to find org: %v", err)
	}

	assert.Equal(org1.Name, one.Name)
	assert.Equal(org1.Slug, one.Slug)
	assert.Empty(one.Domains)

	one, err = repo.FindByUid(org1.Uid, "domains")
	if err != nil {
		t.Fatalf("failed to find org: %v", err)
	}

	assert.Equal(org1.Name, one.Name)
	assert.Equal(org1.Slug, one.Slug)
	assert.NotEmpty(one.Domains)

	one, err = repo.FindBySlug("test-org")
	if err != nil {
		t.Fatalf("failed to find org: %v", err)
	}

	assert.Equal(org1.Name, one.Name)
	assert.Equal(org1.Slug, one.Slug)
	assert.Empty(one.Domains)

	one, err = repo.FindBySlug("test-org", "domains")

	if err != nil {
		t.Fatalf("failed to find org: %v", err)
	}

	assert.Equal(org1.Name, one.Name)
	assert.Equal(org1.Slug, one.Slug)
	assert.NotEmpty(one.Domains)

	err = repo.UpdateName(org1.Uid.String(), "WWF")

	if err != nil {
		t.Fatalf("failed to update org: %v", err)
	}

	one, err = repo.FindByUid(org1.Uid)
	if err != nil {
		t.Fatalf("failed to find org: %v", err)
	}

	assert.Equal("WWF", one.NameFormatted.String)
	assert.Equal("wwf", one.Name)

	one, err = repo.FindByName("WWF")
	if err != nil {
		t.Fatalf("failed to find org: %v", err)
	}

	assert.Equal("wwf", one.Name)
	assert.Equal("WWF", one.NameFormatted.String)

	count, err := repo.Count()
	if err != nil {
		t.Fatalf("failed to count orgs: %v", err)
	}

	assert.Equal(int64(2), count)

	all, err := repo.All()
	if err != nil {
		t.Fatalf("failed to get all orgs: %v", err)
	}

	assert.Len(all, 2)

	err = repo.DeleteById(1)
	if err != nil {
		t.Fatalf("failed to delete org: %v", err)
	}

	count, err = repo.Count()
	if err != nil {
		t.Fatalf("failed to count orgs: %v", err)
	}

	assert.Equal(int64(1), count)
}

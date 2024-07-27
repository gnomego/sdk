package iam_test

import (
	"testing"

	"github.com/gnomego/sdk/stores/iam"
	assert2 "github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	orgDb *gorm.DB
)

func init() {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	orgDb = db

	err = orgDb.AutoMigrate(&iam.OrgTable{}, &iam.OrgDomain{})
	if err != nil {
		panic(err)
	}
}

func TestOrgStore(t *testing.T) {
	assert := assert2.New(t)

	orgStore := iam.NewOrgStore(orgDb, nil)

	org := &iam.Org{
		Name:    "Test Org",
		Slug:    "test-org",
		Domains: []string{"test.org"},
	}

	err := orgStore.Create(org)
	if err != nil {
		t.Fatalf("failed to create org: %v", err)
	}

	err = orgStore.Create(&iam.Org{
		Name: "Test Org 2",
		Slug: "test-org-2",
	})

	if err != nil {
		t.Fatalf("failed to create org: %v", err)
	}

	id := org.Id
	one, err := orgStore.FindByUid(id)
	if err != nil {
		t.Fatalf("failed to find org: %v", err)
	}

	assert.Equal(org.Name, one.Name)
	assert.Equal(org.Slug, one.Slug)
	assert.Empty(one.Domains)

	one, err = orgStore.FindByUid(id, "domains")
	if err != nil {
		t.Fatalf("failed to find org: %v", err)
	}

	assert.Equal(org.Name, one.Name)
	assert.Equal(org.Slug, one.Slug)
	assert.NotEmpty(one.Domains)

	one, err = orgStore.FindBySlug("test-org")
	if err != nil {
		t.Fatalf("failed to find org: %v", err)
	}

	assert.Equal(org.Name, one.Name)
	assert.Equal(org.Slug, one.Slug)
	assert.Empty(one.Domains)

	one, err = orgStore.FindBySlug("test-org", "domains")

	if err != nil {
		t.Fatalf("failed to find org: %v", err)
	}

	assert.Equal(org.Name, one.Name)
	assert.Equal(org.Slug, one.Slug)
	assert.NotEmpty(one.Domains)

	err = orgStore.UpdateName(id, "WWF")

	if err != nil {
		t.Fatalf("failed to update org: %v", err)
	}

	one, err = orgStore.FindByUid(id)
	if err != nil {
		t.Fatalf("failed to find org: %v", err)
	}

	assert.Equal("WWF", one.Name)

	one, err = orgStore.FindByName("WWF")
	if err != nil {
		t.Fatalf("failed to find org: %v", err)
	}

	assert.Equal("WWF", one.Name)

	count, err := orgStore.CountAll()
	if err != nil {
		t.Fatalf("failed to count orgs: %v", err)
	}

	assert.Equal(int64(2), count)

	all, err := orgStore.All()

	if err != nil {
		t.Fatalf("failed to get all orgs: %v", err)
	}

	assert.Len(all, 2)

	_, err = orgStore.DeleteById(1)

	if err != nil {
		t.Fatalf("failed to delete org: %v", err)
	}

	count, err = orgStore.CountAll()
	if err != nil {
		t.Fatalf("failed to count orgs: %v", err)
	}

	assert.Equal(int64(1), count)
}

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

	// assert := assert2.New(t)
	// db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	// if err != nil {
	// 	t.Fatalf("failed to open database: %v", err)
	// }

	// iamDb := iam.IamDb{db}

	// err = iamDb.AutoMigrateIam()
	// if err != nil {
	// 	t.Fatalf("failed to migrate database: %v", err)
	// }

	// user, err := iamDb.NewUser("bob", "
}

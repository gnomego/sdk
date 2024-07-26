package iam_test

import (
	"testing"

	"github.com/gnomego/sdk/stores/iam"
	"github.com/go-playground/validator/v10"
	assert2 "github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestIam(t *testing.T) {
	assert := assert2.New(t)
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	iamDb := iam.IamDb{db}

	err = iamDb.AutoMigrateIam()
	if err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	user, err := iamDb.NewUser("bob", "bob@test.org")

	if err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if ok {
			for _, e := range errs {
				v := e.Value().(string)
				println(e.Field(), v, e.StructField())
			}
		}

		t.Fatalf("failed to create user: %v", err)
	}

	println("before print")
	println(&user.Name)
	println("after print")

	assert.Equal(user.Name, "bob")
	assert.Equal(user.Email, "bob@test.org")
	assert.NotEmpty(user.Uid)
}

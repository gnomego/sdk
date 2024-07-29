package main

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/gnomego/apps/gs/api/v1"
	"github.com/gnomego/apps/gs/globals"
	"github.com/gnomego/apps/gs/validation"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	r := gin.Default()

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	globals.InitDb(db, true)

	globals.InitValidator(validation.NewGsValidator())

	v1.Register(r)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

package v1

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/gnomego/apps/gs/stores"
	"github.com/gnomego/apps/gs/xgin"
)

func Create(c *gin.Context) {
	store := stores.NewOrgStore()
	newOrg := &stores.NewOrg{}

	xgin.BindTargetJson(c, newOrg, "newOrg")
	slog.Debug("new org", slog.Any("newOrg", newOrg))
	response := store.Create(newOrg)
	xgin.Respond(c, response)
}

func Update(c *gin.Context) {
	store := stores.NewOrgStore()
	org := &stores.Org{}

	xgin.BindTargetJson(c, org, "org")
	slog.Debug("new org", slog.Any("newOrg", org))
	response := store.Save(org)
	xgin.Respond(c, response)
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
func All(c *gin.Context) {
	store := stores.NewOrgStore()
	response := store.All()
	xgin.Respond(c, response)
}

package admin

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/melsincostan/borders/db/country"
	"github.com/melsincostan/borders/db/models"
	"github.com/melsincostan/borders/db/transport"
	"github.com/melsincostan/borders/utils"
	"gorm.io/gorm"
)

func create(db *gorm.DB, base, modelType string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var edp editDTO

		if err := ctx.ShouldBind(&edp); err != nil {
			ctx.Error(err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrJSON("could not get creation information"))
			return
		}

		var msg uint
		if modelType == "country" {
			msg = 2
			if _, err := country.Create(db, models.CountryBase{
				Name: edp.Name,
			}); err != nil {
				ctx.Error(err)
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, utils.ErrJSON("could not create country"))
				return
			}
		} else if modelType == "transport" {
			msg = 4
			if _, err := transport.Create(db, models.TransportBase{
				Name: edp.Name,
			}); err != nil {
				ctx.Error(err)
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, utils.ErrJSON("could not create transport"))
				return
			}
		} else {
			ctx.Error(fmt.Errorf("unknown model type: '%s' (expected one of 'country', 'transport')", modelType))
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, utils.ErrJSON("could not create model"))
			return
		}

		ctx.Redirect(http.StatusFound, fmt.Sprintf("%s?msg=%d", base, msg))
	}
}

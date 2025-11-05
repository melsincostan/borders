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

type itemDTO struct {
	ID uint `uri:"id" binding:"required"`
}

type editDTO struct {
	Name string `form:"name" binding:"required"`
}

func edit(db *gorm.DB, base, modelType string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var idp itemDTO
		var edp editDTO

		if err := ctx.ShouldBindUri(&idp); err != nil {
			ctx.Error(err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrJSON("could not get ID parameter from the URL"))
			return
		}

		if err := ctx.ShouldBind(&edp); err != nil {
			ctx.Error(err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrJSON("could not get edit information"))
			return
		}

		var msg uint
		if modelType == "country" {
			msg = 2
			if err := country.Update(db, idp.ID, models.CountryBase{
				Name: edp.Name,
			}); err != nil {
				ctx.Error(err)
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, utils.ErrJSON("could not edit country"))
				return
			}
		} else if modelType == "transport" {
			msg = 4
			if err := transport.Update(db, idp.ID, models.TransportBase{
				Name: edp.Name,
			}); err != nil {
				ctx.Error(err)
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, utils.ErrJSON("could not edit transport"))
				return
			}
		} else {
			ctx.Error(fmt.Errorf("unknown model type: '%s' (expected one of 'country', 'transport')", modelType))
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, utils.ErrJSON("could not edit model"))
			return
		}

		ctx.Redirect(http.StatusFound, fmt.Sprintf("%s?msg=%d", base, msg))
	}
}

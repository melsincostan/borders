package admin

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/melsincostan/borders/db/country"
	"github.com/melsincostan/borders/db/transport"
	"github.com/melsincostan/borders/utils"
	"gorm.io/gorm"
)

func delete(db *gorm.DB, base, modelType string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var idp itemDTO

		if err := ctx.ShouldBindUri(&idp); err != nil {
			ctx.Error(err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrJSON("could not get ID parameter from the URL"))
			return
		}

		var msg uint
		if modelType == "country" {
			msg = 3
			if err := country.Delete(db, idp.ID); err != nil {
				ctx.Error(err)
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, utils.ErrJSON("could not delete country"))
				return
			}
		} else if modelType == "transport" {
			msg = 6
			if err := transport.Delete(db, idp.ID); err != nil {
				ctx.Error(err)
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, utils.ErrJSON("could not delete transport"))
				return
			}
		} else {
			ctx.Error(fmt.Errorf("unknown model type: '%s' (expected one of 'country', 'transport')", modelType))
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, utils.ErrJSON("could not delete model"))
			return
		}

		ctx.Redirect(http.StatusFound, fmt.Sprintf("%s?msg=%d", base, msg))
	}
}

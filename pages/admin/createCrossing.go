package admin

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/melsincostan/borders/auth"
	"github.com/melsincostan/borders/db/crossing"
	"github.com/melsincostan/borders/db/models"
	"github.com/melsincostan/borders/pages/stats"
	"github.com/melsincostan/borders/utils"
	"gorm.io/gorm"
)

type createCrossingDTO struct {
	Country   uint     `form:"country" binding:"required"`
	Transport uint     `form:"transport" binding:"required"`
	RawTime   string   `form:"time" binding:"required"`
	Check     []string `form:"check[]"`
}

const (
	timeFormat  = "2006-01-02T15:04"
	borderValue = "border"
	papersValue = "papers"
)

func createCrossing(db *gorm.DB, base string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		rusr, ok := ctx.Get(auth.UIDContextKey)
		if !ok {
			ctx.Error(fmt.Errorf("could not find subject in context (key: '%s')", auth.UIDContextKey))
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, utils.ErrJSON("could not find user ID in context"))
			return
		}

		usr, ok := rusr.(uint)
		if !ok {
			if !ok {
				ctx.Error(errors.New("could not coax user id into an integer"))
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, utils.ErrJSON("could not parse user ID"))
				return
			}
		}

		var params createCrossingDTO
		if err := ctx.ShouldBind(&params); err != nil {
			ctx.Error(err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrJSON("couldn't parse request"))
			return
		}

		actualTime, err := time.ParseInLocation(timeFormat, params.RawTime, time.Local)
		if err != nil {
			ctx.Error(err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrJSON("couldn't parse time"))
			return
		}

		borderCheck := lookup(params.Check, borderValue)
		papersCheck := lookup(params.Check, papersValue)

		if _, err := crossing.Create(db, models.CrossingBase{
			When:        actualTime,
			BorderCheck: borderCheck,
			PapersCheck: papersCheck,
		}, params.Country, params.Transport, usr); err != nil {
			ctx.Error(err)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, utils.ErrJSON("couldn't create crossing"))
			return
		}
		stats.ExpiryLock.Lock()
		stats.Expiry = time.Now().Add(-5 * time.Second)
		stats.ExpiryLock.Unlock()
		ctx.Redirect(http.StatusFound, fmt.Sprintf("%s?msg=1", base))
	}
}

func lookup(arr []string, val string) (ok bool) {
	for _, arrval := range arr {
		if arrval == val {
			return true
		}
	}
	return false
}

package admin

import (
	"crypto/rand"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/melsincostan/borders/auth"
	"github.com/melsincostan/borders/db/models"
	"github.com/melsincostan/borders/db/user"
	"github.com/melsincostan/borders/utils"
	"gorm.io/gorm"
)

type changePWDTO struct {
	Old     string `form:"old" binding:"required"`
	New     string `form:"password" binding:"required"`
	Confirm string `form:"password-confirm" binding:"required"`
}

type createDTO struct {
	Name string `form:"name" binding:"required"`
}

func createUser(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var params createDTO

		if err := ctx.ShouldBind(&params); err != nil {
			ctx.Error(err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrJSON("could not get parameters"))
			return
		}

		tmpPW := fmt.Sprintf("change-me-%s", rand.Text())

		base := models.UserBase{
			Name: params.Name,
		}

		if _, err := user.Create(db, base, tmpPW); err != nil {
			ctx.Error(err)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, utils.ErrJSON("could not create user"))
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message":  "The new user was created. User details follow. The password will not be shown again and will be lost if not copied now. Please make sure this password is changed as soon as convenient.",
			"username": params.Name,
			"password": tmpPW,
		})
	}
}

func changePW(db *gorm.DB, base string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var params changePWDTO

		if err := ctx.ShouldBind(&params); err != nil {
			ctx.Error(err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrJSON("could not get parameters"))
			return
		}

		ruid, ok := ctx.Get(auth.UIDContextKey)
		if !ok {
			ctx.Error(fmt.Errorf("could not get uid from context (key: '%s')", auth.UIDContextKey))
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, utils.ErrJSON("could not get user ID"))
			return
		}

		uid, ok := ruid.(uint)
		if !ok {
			ctx.Error(fmt.Errorf("could not coax uid key into uint (value: '%#v')", ruid))
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, utils.ErrJSON("could not parse user ID"))
			return
		}

		if params.New != params.Confirm {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrJSON("password and confirmation do not match"))
			return
		}

		if err := user.ChangePassword(db, uid, params.Old, params.New); err != nil {
			ctx.Error(err)
			if errors.Is(err, user.ErrWrongPassword) {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, utils.ErrJSON("old password could not be verified"))
			} else {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, utils.ErrJSON("could not change password"))
			}
			return
		}
		ctx.Redirect(http.StatusFound, fmt.Sprintf("%s?msg=%d", base, 5))
	}
}

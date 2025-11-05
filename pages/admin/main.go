package admin

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/melsincostan/borders/db/country"
	"github.com/melsincostan/borders/db/models"
	"github.com/melsincostan/borders/db/transport"
	"github.com/melsincostan/borders/utils"
	"gorm.io/gorm"
)

type manageable interface {
	models.Country | models.Transport
}

type adminObj struct {
	CreateCrossing   string
	ChangePassword   string
	CreateUser       string
	BorderValue      string
	PapersValue      string
	Logout           string
	ManageCountries  manageObj[models.Country]
	ManageTransports manageObj[models.Transport]
	Message          string
}

type manageObj[T manageable] struct {
	Items []T
	Base  string
	Type  string
}

type msg struct {
	MsgNumber uint `form:"msg"`
}

func main(db *gorm.DB, base string) gin.HandlerFunc {
	createCrossing := fmt.Sprintf("%s%s", strings.TrimSuffix(base, "/"), crossingSubmit)
	logout := fmt.Sprintf("%s/logout", strings.TrimSuffix(base, "/"))
	changePassword := fmt.Sprintf("%s/change-password", strings.TrimSuffix(base, "/"))
	cUser := fmt.Sprintf("%s/create-user", strings.TrimSuffix(base, "/"))
	return func(ctx *gin.Context) {
		countries, err := country.Read(db)
		if err != nil {
			ctx.Error(err)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, utils.ErrJSON("could not load countries"))
			return
		}

		transports, err := transport.Read(db)
		if err != nil {
			ctx.Error(err)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, utils.ErrJSON("could not load transports"))
		}

		var msgp msg
		if err := ctx.ShouldBindQuery(&msgp); err != nil {
			log.Printf("Error parsing message number param on admin page: %s", err.Error())
		}

		ctx.HTML(http.StatusOK, "admin.html", adminObj{
			ManageCountries: manageObj[models.Country]{
				Items: countries,
				Base:  base,
				Type:  "country",
			},
			ManageTransports: manageObj[models.Transport]{
				Items: transports,
				Base:  base,
				Type:  "transport",
			},
			CreateCrossing: createCrossing,
			ChangePassword: changePassword,
			CreateUser:     cUser,
			BorderValue:    borderValue,
			PapersValue:    papersValue,
			Logout:         logout,
			Message:        messages[msgp.MsgNumber],
		})
	}
}

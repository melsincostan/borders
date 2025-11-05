package stats

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/melsincostan/borders/db/country"
	"github.com/melsincostan/borders/db/crossing"
	"github.com/melsincostan/borders/db/models"
	"github.com/melsincostan/borders/db/transport"
	"github.com/melsincostan/borders/utils"
	"gorm.io/gorm"
)

type statsObj struct {
	crossing.Stats
	OGTitle              string
	CrossingCheckPer     float32
	CheckIDPer           float32
	CrossingIdPer        float32
	WorstCountryPer      float32
	WorstTransportPer    float32
	CountryLeaderboard   []models.LeaderboardEntry[models.Country]
	TransportLeaderboard []models.LeaderboardEntry[models.Transport]
}

var cache *statsObj
var Expiry time.Time

const expiryDelay = 15 * time.Minute

func show(db *gorm.DB, ogtitle string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if cache == nil || time.Now().After(Expiry) {

			stats, err := crossing.GetStats(db)
			if err != nil {
				ctx.Error(err)
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, utils.ErrJSON("there was an issue retrieving the overview statistics"))
				return
			}

			countryLeaderboard, err := country.Leaderboard(db)
			if err != nil {
				ctx.Error(err)
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, utils.ErrJSON("there was an issue retrieving the country leaderboard"))
				return
			}

			transportLeaderboard, err := transport.Leaderboard(db)
			if err != nil {
				ctx.Error(err)
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, utils.ErrJSON("there was an issue retrieving the transport leaderboard"))
				return
			}

			Expiry = time.Now().Add(expiryDelay)
			cache = &statsObj{
				Stats:                *stats,
				CrossingCheckPer:     float32(stats.TotalBorderChecks) / float32(stats.TotalCrossings) * (100),
				CheckIDPer:           float32(stats.TotalIdChecks) / float32(stats.TotalBorderChecks) * 100,
				CrossingIdPer:        float32(stats.TotalIdChecks) / float32(stats.TotalCrossings) * 100,
				WorstCountryPer:      stats.WorstCountry.Ratio * 100,
				WorstTransportPer:    stats.WorstTransport.Ratio * 100,
				CountryLeaderboard:   countryLeaderboard,
				TransportLeaderboard: transportLeaderboard,
				OGTitle:              ogtitle,
			}
		}

		ctx.HTML(http.StatusOK, "stats.html", *cache)
	}
}

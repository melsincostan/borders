package country

import (
	"github.com/melsincostan/borders/db/models"
	"gorm.io/gorm"
)

const (
	leaderboardQuery = "SELECT countries.*, CAST(COUNT(DISTINCT(tbc.id)) AS REAL) / COUNT(DISTINCT(tc.id)) AS ratio, COUNT(DISTINCT(tc.id)) AS crossings, COUNT(DISTINCT(tbc.id)) AS checks, COUNT(DISTINCT(tpc.id)) AS papers_checks FROM countries LEFT JOIN crossings AS tc ON tc.country_id = countries.id LEFT JOIN crossings AS tbc ON tbc.country_id = countries.id AND tbc.border_check = 1 LEFT JOIN crossings AS tpc ON tpc.country_id = countries.id AND tpc.papers_check = 1 GROUP BY countries.id HAVING crossings > 0 ORDER BY ratio DESC, crossings DESC"
)

func Read(db *gorm.DB) (res []models.Country, err error) {
	res = []models.Country{}
	if err := db.Model(&models.Country{}).Order("name asc").Scan(&res).Error; err != nil {
		return nil, err
	}
	return
}

func Leaderboard(db *gorm.DB) (res []models.LeaderboardEntry[models.Country], err error) {
	res = []models.LeaderboardEntry[models.Country]{}
	if err := db.Raw(leaderboardQuery).Scan(&res).Error; err != nil {
		return nil, err
	}
	return
}

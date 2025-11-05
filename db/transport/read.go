package transport

import (
	"github.com/melsincostan/borders/db/models"
	"gorm.io/gorm"
)

const (
	leaderboardQuery = "SELECT transports.*, CAST(COUNT(DISTINCT(tbc.id)) AS REAL) / COUNT(DISTINCT(tc.id)) AS ratio, COUNT(DISTINCT(tc.id)) AS crossings, COUNT(DISTINCT(tbc.id)) AS checks, COUNT(DISTINCT(tpc.id)) AS papers_checks FROM transports LEFT JOIN crossings AS tc ON tc.transport_id = transports.id LEFT JOIN crossings AS tbc ON tbc.transport_id = transports.id AND tbc.border_check = 1 LEFT JOIN crossings AS tpc ON tpc.transport_id = transports.id AND tpc.papers_check = 1 GROUP BY transports.id HAVING crossings > 0 ORDER BY ratio DESC, crossings DESC"
)

func Read(db *gorm.DB) (res []models.Transport, err error) {
	if err := db.Model(&models.Transport{}).Order("name ASC").Scan(&res).Error; err != nil {
		return nil, err
	}
	return
}

func Leaderboard(db *gorm.DB) (res []models.LeaderboardEntry[models.Transport], err error) {
	res = []models.LeaderboardEntry[models.Transport]{}
	if err := db.Raw(leaderboardQuery).Scan(&res).Error; err != nil {
		return nil, err
	}
	return
}

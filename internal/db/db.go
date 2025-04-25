package db

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/sjimenezl/phishrivals/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("phish.db"), &gorm.Config{})
	if err != nil {
		panic("cant connect to db")
	}
}

func SaveEnrichment(entity *models.DomainInfo) error {
	result := DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&entity)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

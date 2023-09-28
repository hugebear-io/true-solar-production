package main

import (
	"database/sql"
	"log"

	"github.com/hugebear-io/true-solar-production/infra"
	"github.com/hugebear-io/true-solar-production/model"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/gorm"
)

type LegacyData struct {
	Code string  `json:"code" binding:"required"`
	Name string  `json:"name" binding:"required"`
	Area *string `json:"area"`
}

func main() {
	legacy, err := sql.Open("sqlite3", "legacy.db")
	if err != nil {
		log.Fatal(err)
	}

	row, err := legacy.Query("SELECT code, name, area FROM site_region_mapping;")
	if err != nil {
		log.Fatal(err)
	}

	db, _ := infra.NewGormDB()
	err = db.Transaction(func(tx *gorm.DB) error {
		count := 1
		for row.Next() {
			var tmp LegacyData
			err := row.Scan(&tmp.Code, &tmp.Name, &tmp.Area)
			if err != nil {
				return err
			}

			site := model.SiteRegionMapping{
				ID:   int64(count),
				Code: tmp.Code,
				Name: tmp.Name,
				Area: tmp.Area,
			}

			if err := tx.Create(&site).Error; err != nil {
				return err
			}

			count++
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}

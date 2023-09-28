package model

import (
	"time"

	"gorm.io/gorm"
)

type SiteRegionMapping struct {
	ID        int64      `gorm:"column:id"`
	Code      string     `gorm:"column:code"`
	Name      string     `gorm:"column:name"`
	Area      *string    `gorm:"column:area"`
	CreatedAt *time.Time `gorm:"created_at"`
	UpdatedAt *time.Time `gorm:"updated_at"`
}

func (*SiteRegionMapping) TableName() string {
	return "tbl_site_region_mapping"
}

func (r *SiteRegionMapping) BeforeCreate(tx *gorm.DB) error {
	var count int64
	tx.Model(&SiteRegionMapping{}).Count(&count)
	r.ID = count + 1
	return nil
}

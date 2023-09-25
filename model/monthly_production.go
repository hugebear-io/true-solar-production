package model

import (
	"time"

	"github.com/hugebear-io/true-solar-production/util"
	"go.openly.dev/pointy"
)

type MonthlyProductionDocument struct {
	Date               *time.Time `json:"date"`
	VendorType         *string    `json:"vendor_type"`
	Area               *string    `json:"area"`
	SiteName           *string    `json:"site_name"`
	InstalledCapacity  *float64   `json:"installed_capacity"`
	MonthlyProduction  *float64   `json:"monthly_production"`
	Latitude           *float64   `json:"lat"`
	Longitude          *float64   `json:"lng"`
	Target             *float64   `json:"target"`
	ProductionToTarget *float64   `json:"production_to_target"`
	Criteria           *string    `json:"criteria"`
}

func (d *MonthlyProductionDocument) parseString(data string) *string {
	if util.EmptyString(data) {
		return nil
	}

	return pointy.String(data)
}

func (d *MonthlyProductionDocument) SetDate(data *time.Time) {
	d.Date = data
}

func (d *MonthlyProductionDocument) SetVendorType(data string) {
	d.VendorType = d.parseString(data)
}

func (d *MonthlyProductionDocument) SetArea(data string) {
	d.Area = d.parseString(data)
}

func (d *MonthlyProductionDocument) SetSiteName(data string) {
	d.SiteName = d.parseString(data)
}

func (d *MonthlyProductionDocument) SetInstalledCapacity(data *float64) {
	d.InstalledCapacity = data
}

func (d *MonthlyProductionDocument) SetMonthlyProduction(data *float64) {
	d.MonthlyProduction = data
}

func (d *MonthlyProductionDocument) SetLatitude(data *float64) {
	d.Latitude = data
}

func (d *MonthlyProductionDocument) SetLongitude(data *float64) {
	d.Longitude = data
}

func (d *MonthlyProductionDocument) SetProductionToTarget(data *float64) {
	d.ProductionToTarget = data
}

func (d *MonthlyProductionDocument) SetTarget(data *float64) {
	d.Target = data
}

func (d *MonthlyProductionDocument) SetCriteria(data *float64) {
	if data == nil {
		d.Criteria = nil
	}

	value := pointy.Float64Value(data, 0)
	var criteria string
	if value >= 100 {
		criteria = ">=100%"
	} else if value >= 80 {
		criteria = ">=80%"
	} else if value >= 60 {
		criteria = ">=60%"
	} else if value >= 50 {
		criteria = ">=50%"
	} else {
		criteria = "<50%"
	}

	d.Criteria = &criteria
}

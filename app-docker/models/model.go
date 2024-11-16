package model

import (
	"gorm.io/gorm"
)

type Visit struct {
    gorm.Model
    StoreID  string `json:"store_id" gorm:"index"` 
    ImageURL string `json:"image_url" gorm:"type:json"` // Use JSON to store array of strings
    VisitTime string `json:"visit_time"`
}

type Store struct {
	gorm.Model
	StoreID      string           `json:"store_id" gorm:"index"`
	StoreName    string           `json:"store_name"`       // New field for Store Name
	AreaCode     string           `json:"area_code"`        // New field for Area Code
	VisitorCount int              `json:"visitor_count"`
	Images       []Image          `json:"images" gorm:"foreignKey:StoreID;references:StoreID"` 
}

type VisitsResponse struct {
	Count  int     `json:"count"`
	Visits []Visit `json:"visits"`
}

type Job struct {
    gorm.Model
    JobID  string  `json:"job_id" gorm:"index"` // Add index to JobID
    Images []Image `json:"images" gorm:"foreignKey:JobID;references:JobID"` 
    Status string  `json:"status"`
}


type Image struct {
	gorm.Model
	StoreID   string  `json:"-" gorm:"index"` // Add index for faster lookups, "-" to exclude from JSON
	JobID     string  `json:"-" gorm:"index"` // Add index for faster lookups, "-" to exclude from JSON
	ImageURL  string  `json:"image_url"`
	Status    string  `json:"status"`
	Perimeter float64 `json:"perimeter"`
	ErrMessage string  `json:"err"`
}
package model

// Visit model with foreign key to Store
type Visit struct {
	ID        uint   `gorm:"primaryKey"`
	StoreID   string `json:"store_id"`
	ImageURL  []string `json:"image_url"`
	VisitTime string `json:"visit_time"`
	Store     Store `gorm:"foreignKey:StoreID;constraint:OnDelete:CASCADE;"` // Foreign key reference to Store
}

// Store model with foreign key reference for Job and Image
type Store struct {
	ID            uint              `gorm:"primaryKey"`
	StoreID       string            `json:"store_id"`
	VisitorCount  int               `json:"visitor_count"`
	Images        []Image           `json:"images"`   // One-to-many relation with Image
	Visits        []Visit           `json:"visits"`   // One-to-many relation with Visit
	Jobs          []Job             `json:"jobs"`     // One-to-many relation with Job
}

// VisitsResponse model
type VisitsResponse struct {
	Count  int     `json:"count"`
	Visits []Visit `json:"visits"`
}

// Job model with foreign key to Store
type Job struct {
	ID      uint              `gorm:"primaryKey"`
	JobID   string            `json:"job_id"`
	StoreID string            `json:"store_id"`
	Images  []Image           `json:"images"`
	Status  string            `json:"status"`
	Store   Store             `gorm:"foreignKey:StoreID;constraint:OnDelete:CASCADE;"` // Foreign key reference to Store
}

// Image model with foreign key to Job
type Image struct {
	ID        uint    `gorm:"primaryKey"`
	ImageURL  string  `json:"image_url"`
	Status    string  `json:"status"`
	Perimeter float64 `json:"perimeter"`
	ErrMessage string `json:"err"`
	JobID     uint    `json:"job_id"`
	Job       Job     `gorm:"foreignKey:JobID;constraint:OnDelete:CASCADE;"` // Foreign key reference to Job
}

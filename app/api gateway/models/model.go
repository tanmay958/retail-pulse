package model

type Visit struct {
	StoreID   string   `json:"store_id"`
	ImageURL  []string `json:"image_url"`
	VisitTime string   `json:"visit_time"`
}
type Store struct {
	StoreID      string           `json:"store_id"`
	VisitorCount int              `json:"visitor_count"`
	Images       map[string]Image `json:"images"` // Changed to map[string]Image
}

type VisitsResponse struct {
	Count  int     `json:"count"`
	Visits []Visit `json:"visits"`
}

type Job struct {
	JobID  string           `json:"job_id"`
	Images map[string]Image `json:"images"`
	Status string           `json:"status"`
}

type Image struct {
	ImageURL   string  `json:"image_url"`
	Status     string  `json:"status"` // Added status field
	Perimeter  float64 `json:"perimeter"`
	ErrMessage string  `json:"err"`
}
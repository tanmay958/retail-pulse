package model

type Visit struct {
	StoreID   string   `json:"store_id"`
	ImageURL  []string `json:"image_url"`
	VisitTime string   `json:"visit_time"`
}

type VisitsResponse struct {
	Count  int     `json:"count"`
	Visits []Visit `json:"visits"`
}

type Job struct {
	JobID    int      `json:"job_id"`
	StoreID  string   `json:"store_id"`
	ImageURL []string `json:"image_url"`
}

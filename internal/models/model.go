package models

type Event struct {
	ID              string      `json:"id"`
	Camera          string      `json:"camera"`
	FrameTime       float64     `json:"frame_time"`
	SnapshotTime    float64     `json:"snapshot_time"`
	Label           string      `json:"label"`
	TopScore        float64     `json:"top_score"`
	FalsePositive   bool        `json:"false_positive"`
	StartTime       float64     `json:"start_time"`
	EndTime         interface{} `json:"end_time"`
	Score           float64     `json:"score"`
	Box             []int       `json:"box"`
	Area            int         `json:"area"`
	Region          []int       `json:"region"`
	CurrentZones    []string    `json:"current_zones"`
	EnteredZones    []string    `json:"entered_zones"`
	Thumbnail       interface{} `json:"thumbnail"`
	HasSnapshot     bool        `json:"has_snapshot"`
	HasClip         bool        `json:"has_clip"`
	Stationary      bool        `json:"stationary"`
	MotionlessCount int         `json:"motionless_count"`
	PositionChanges int         `json:"position_changes"`
}

type Events struct {
	Type   string `json:"type"`
	Before Event  `json:"before"`
	After  Event  `json:"after"`
}

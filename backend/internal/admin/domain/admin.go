package domain

// DashboardMetrics holds admin dashboard aggregates.
type DashboardMetrics struct {
	Users    MetricGroup `json:"users"`
	Tasks    MetricGroup `json:"tasks"`
	Revenue  MetricGroup `json:"revenue"`
	Disputes MetricGroup `json:"disputes"`
}

type MetricGroup struct {
	Total    int64   `json:"total,omitempty"`
	Active   int64   `json:"active,omitempty"`
	Open     int64   `json:"open,omitempty"`
	Growth7d float64 `json:"growth7d,omitempty"`
	Delta7d  int64   `json:"delta7d,omitempty"`
	MonthToDate int64 `json:"monthToDate,omitempty"`
}

// UserAdminUpdate holds admin user moderation fields.
type UserAdminUpdate struct {
	IsActive *bool
	Role     *string
}

// DisputeResolution holds admin dispute outcome.
type DisputeResolution struct {
	Outcome string // refund_requester, release_agent, split
	Note    string
}

// RevenueReport holds revenue analytics.
type RevenueReport struct {
	Range   string       `json:"range"`
	Total   int64        `json:"total"`
	Points  []DataPoint  `json:"points"`
}

type DataPoint struct {
	Date  string `json:"date"`
	Value int64  `json:"value"`
}

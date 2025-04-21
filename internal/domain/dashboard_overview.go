package domain

type DashboardOverviewKPIs struct {
	Logs       DashboardOverviewKPI
	Errors     DashboardOverviewKPI
	Warnings   DashboardOverviewKPI
	Info       DashboardOverviewKPI
	LogsPerApp map[string]struct {
		AppName string
		Total   int64
	}
	LogsByPeriod map[Period]struct {
		Total int64
	}
}

type DashboardOverviewKPI struct {
	Total      int64
	Percentage float64
}

type Period struct {
	Year  int
	Month int
}

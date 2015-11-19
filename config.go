package main

type config struct {
	StatsAccount string `json:"stats_account"`
	StatsFormat  string `json:"stats_format"`
	Uchiwa       struct {
		Host      string `json:"host"`
		Interval  int    `json:"interval"`
		MaxErrors int    `json:"max_errors"`
	} `json:"uchiwa"`
	Alerts []string `json:"alerts"`
}

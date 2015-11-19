package main

type Event struct {
	Acknowledged bool   `json:"acknowledged"`
	Action       string `json:"action"`
	Check        struct {
		Command  string  `json:"command"`
		Duration float64 `json:"duration"`
		//Executed         uint64   `json:"executed"`
		Handlers []string `json:"handlers"`
		History  []string `json:"history"`
		Interval uint64   `json:"interval"`
		//Issued           uint64   `json:"issued"`
		Name string `json:"name"`
		//Occurrences      uint64 `json:"occurrences"`
		Output     string `json:"output"`
		Standalone bool   `json:"standalone"`
		//Status           uint64 `json:"status"`
		//TotalStateChange uint64 `json:"total_state_change"`
	} `json:"check"`
	Client struct {
		Acknowledged bool   `json:"acknowledged"`
		Address      string `json:"address"`
		Keepalive    struct {
			Thresholds struct {
				//Critical uint64 `json:"critical"`
				//Warning  uint64 `json:"warning"`
			} `json:"thresholds"`
		} `json:"keepalive"`
		Name          string   `json:"name"`
		Subscriptions []string `json:"subscriptions"`
		//Timestamp     uint64   `json:"timestamp"`
		Version string `json:"version"`
	} `json:"client"`
	Dc          string `json:"dc"`
	ID          string `json:"id"`
	Occurrences uint64 `json:"occurrences"`
}

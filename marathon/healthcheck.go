package marathon

type HealthCheckResult struct {
	Alive                 bool   `json:"alive"`
	ConsucutiveFailures   int    `json:"consecutiveFailures"`
	FirstSuccessTimestamp string `json:"firstSuccess"`
	LastFailureTimestamp  string `json:"lastFailure"`
	LastSuccessTimestamp  string `json:"lastSuccess"`
}

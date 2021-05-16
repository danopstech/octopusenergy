package octopusenergy

// errorResponse is the returned body when API error accrues
type errorResponse struct {
	Detail string `json:"detail"`
}

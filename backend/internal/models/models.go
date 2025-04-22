package models

type AccessCheckResponse struct {
	VPN            bool     `json:"vpn"`
	Production     bool     `json:"production"`
	ConfigTool     bool     `json:"configTool"`
	CurrentProfile string   `json:"currentProfile"`
	MissingGroups  []string `json:"missingGroups"`
	ProfileARN     string   `json:"profileArn,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

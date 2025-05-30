package model

type APIHealth struct {
	Version    string           `json:"version"`
	Downstream []DownstreamItem `json:"downstream"`
}

type DownstreamItem struct {
	ServiceName string `json:"serviceName"`
	Version     string `json:"version,omitempty"`
	Status      Status `json:"status"`
}

type Status string

const (
	Up   = "Up"
	Down = "Down"
)

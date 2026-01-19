package event

type Envelope struct {
	Event      string      `json:"event"`
	Version    string      `json:"version"`
	OccurredAt string      `json:"occurred_at"`
	Source     string      `json:"source"`
	Data       interface{} `json:"data"`
}

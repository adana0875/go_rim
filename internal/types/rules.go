package types

type CommunityRules struct {
	Timestamp int64           `json:"timestamp"`
	Rules     map[string]Rule `json:"rules"`
}

type Rule struct {
	LoadBefore map[string]any `json:"loadBefore,omitempty"`
	LoadAfter  map[string]any `json:"loadAfter,omitempty"`
}

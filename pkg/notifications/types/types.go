package types

type EventMeta struct {
	Message string         `yaml:"message" json:"message"`
	Meta    map[string]any `yaml:"meta" json:"meta"`
}

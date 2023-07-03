package model

type Item struct {
	PlayID int    `json:"id" mapstructure:"id"`
	Key    string `json:"key" mapstructure:"key"`
	Value  any    `json:"value" mapstructure:"value"`
}

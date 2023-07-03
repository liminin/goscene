package model

type Play struct {
	ID           int    `json:"id" mapstructure:"id"`
	UserID       int    `json:"user_id" mapstructure:"user_id"`
	SceneKey     string `json:"scene_key" mapstructure:"scene_key"`
	CurrentIndex int    `json:"current" mapstructure:"current"`
	FirstTime    bool   `json:"first_time" mapstructure:"first_time"`
}

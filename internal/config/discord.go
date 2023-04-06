package config

// Discord is the configuration for the Discord bot.
type Discord struct {
	Production bool                         `json:"production" yaml:"production" default:"false"`
	Token      string                       `json:"token" yaml:"token" default:""`
	Locales    map[string]map[string]string `json:"locales" yaml:"locales" default:"{}"`
	Servers    []struct {
		ID           string `json:"id" yaml:"id"`
		Language     string `json:"language" yaml:"language" default:"enUS"`
		ChatChannels []struct {
			ID         string `json:"id" yaml:"id"`
			TokenLimit int    `json:"token_limit" yaml:"token_limit"`
		} `json:"chat_channels" yaml:"chat_channels" default:"[]"`
		Commands struct {
			ClearContext struct {
				Enable  bool     `json:"enable" yaml:"enable" default:"false"`
				Aliases []string `json:"aliases" yaml:"aliases" default:"[]"`
			} `json:"clear_context" yaml:"clear_context"`
		} `json:"commands" yaml:"commands"`
	} `json:"servers" yaml:"servers" default:"[]"`
}

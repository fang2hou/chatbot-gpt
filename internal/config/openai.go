package config

// OpenAI is the configuration for the OpenAI API.
type OpenAI struct {
	Token                  string `json:"token" yaml:"token" default:""`
	ModelID                string `json:"model_id" yaml:"model_id" default:"gpt-3.5-turbo-0301"`
	TokenPredictionModelID string `json:"token_prediction_model_id" yaml:"token_prediction_model_id" default:"gpt-3.5-turbo"`
}

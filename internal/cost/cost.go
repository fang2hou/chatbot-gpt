package cost

import "github.com/sashabaranov/go-openai"

// Calculator is a calculator for calculating the cost of a completion.
type Calculator struct {
	*openai.Model
}

// NewCalculator creates a new calculator.
func NewCalculator(model *openai.Model) *Calculator {
	return &Calculator{
		Model: model,
	}
}

// GetPromptCost returns the cost of a prompt.
func (c *Calculator) GetPromptCost(numTokens int) float64 {
	modelID := shortModelID(c.ID)
	if costTable, ok := modelCosts[modelID]; ok {
		return float64(numTokens) * costTable.promptCost
	}

	return 0
}

// GetSampledCost returns the cost of a sampled completion.
func (c *Calculator) GetSampledCost(numTokens int) float64 {
	modelID := shortModelID(c.ID)
	if costTable, ok := modelCosts[modelID]; ok {
		return float64(numTokens) * costTable.sampledCost
	}

	return 0
}

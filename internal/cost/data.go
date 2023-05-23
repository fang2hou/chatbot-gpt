package cost

type modelCostData struct {
	promptCost  float64
	sampledCost float64
}

var modelCosts = map[string]modelCostData{
	"gpt-3.5-turbo": {
		promptCost:  0.002 / 1000,
		sampledCost: 0.002 / 1000,
	},
	"gpt-4": {
		promptCost:  0.03 / 1000,
		sampledCost: 0.06 / 1000,
	},
	"gpt-4-32k": {
		promptCost:  0.06 / 1000,
		sampledCost: 0.12 / 1000,
	},
}

func shortModelID(modelID string) string {
	switch modelID {
	case "gpt-3.5-turbo-0301":
		return "gpt-3.5-turbo"
	case "gpt-4-0314":
		return "gpt-4"
	case "gpt-4-32k-0314":
		return "gpt-4-32k"
	}

	return modelID
}

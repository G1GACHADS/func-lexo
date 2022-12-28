package clients

import "github.com/G1GACHADS/func-lexo/internal/logger"

type ClientsConfig struct {
	AzureOCPSubscriptionKey string
	AzureOCPEndpoint        string
}

type Clients struct {
	AzureComputerVisionClient *AzureComputerVisionClient
}

func New(cfg ClientsConfig) (Clients, error) {
	c := Clients{}

	logger.M.Debug("Creating Azure Computer Vision client")

	c.AzureComputerVisionClient = NewAzureComputerVisionClient(cfg.AzureOCPEndpoint, cfg.AzureOCPSubscriptionKey)

	logger.M.Debug("Azure Computer Vision client created")

	return c, nil
}

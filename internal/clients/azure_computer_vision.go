package clients

import (
	"context"
	"mime/multipart"

	"github.com/Azure/azure-sdk-for-go/services/cognitiveservices/v2.0/computervision"
	"github.com/Azure/go-autorest/autorest"
)

// AzureComputerVisionClient is a wrapper around the Azure Computer Vision client
type AzureComputerVisionClient struct {
	client computervision.BaseClient
}

// NewAzureComputerVisionClient creates a new AzureComputerVisionClient
func NewAzureComputerVisionClient(endpoint, subscriptionKey string) *AzureComputerVisionClient {
	client := computervision.New(endpoint)
	client.Authorizer = autorest.NewCognitiveServicesAuthorizer(subscriptionKey)

	return &AzureComputerVisionClient{
		client: client,
	}
}

// RecognizePrinterOCR performs OCR on a given image file and returns the result
func (acv AzureComputerVisionClient) RecognizePrintedOCR(ctx context.Context, img multipart.File) (computervision.OcrResult, error) {
	// TODO: dynamic language param
	ocrResult, err := acv.client.RecognizePrintedTextInStream(ctx, true, img, computervision.En)
	if err != nil {
		return computervision.OcrResult{}, err
	}

	return ocrResult, nil
}

type MergedOCRResultRegion struct {
	BoundingBox string `json:"bounding_box"`
	Content     string `json:"content"`
}

type MergedOCRResultLinesResult struct {
	Regions       []MergedOCRResultRegion `json:"regions"`
	MergedContent string                  `json:"merged_content"`
}

func (acv AzureComputerVisionClient) MergeOCRResultLines(result computervision.OcrResult) MergedOCRResultLinesResult {
	var mergedResult MergedOCRResultLinesResult

	mergedResult.Regions = make([]MergedOCRResultRegion, len(*result.Regions))

	for regionIDX, region := range *result.Regions {
		mergedResult.Regions[regionIDX].BoundingBox = *region.BoundingBox
		for _, line := range *region.Lines {
			mergedResult.Regions[regionIDX].Content = ""
			for _, word := range *line.Words {
				mergedResult.Regions[regionIDX].Content += *word.Text + " "
				mergedResult.MergedContent += *word.Text + " "
			}
		}
	}

	return mergedResult
}

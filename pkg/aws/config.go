package aws

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
)

// LoadAWSConfig loads the AWS SDK configuration
func LoadAWSConfig(region string) aws.Config {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		fmt.Println("Failed to load AWS SDK configuration, " + err.Error())
		os.Exit(1)
	}
	defaultRegion := os.Getenv("AWS_DEFAULT_REGION")
	if region != "" {
		cfg.Region = region
	} else if defaultRegion != "" {
		cfg.Region = defaultRegion
	}
	return cfg
}

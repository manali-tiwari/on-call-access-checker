package auth

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type AWSAuthenticator interface {
	GetProfileInfo() (*AWSProfile, error)
}

type AWSAuth struct {
	region string
	mock   bool
}

type AWSProfile struct {
	Name      string
	ARN       string
	AccountID string
}

func NewAWSAuth() (*AWSAuth, error) {
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "us-east-1"
	}

	// Enable mock mode if no AWS credentials are set
	if os.Getenv("AWS_ACCESS_KEY_ID") == "" &&
		os.Getenv("AWS_SECRET_ACCESS_KEY") == "" &&
		os.Getenv("AWS_PROFILE") == "" {
		fmt.Println("WARNING: Running in AWS mock mode - no credentials provided")
		return &AWSAuth{region: region, mock: true}, nil
	}

	return &AWSAuth{region: region}, nil
}

func (a *AWSAuth) GetProfileInfo() (*AWSProfile, error) {
	if a.mock {
		// Return mock AWS profile
		return &AWSProfile{
			Name:      "mock-profile",
			ARN:       "arn:aws:iam::123456789012:user/mock-user",
			AccountID: "123456789012",
		}, nil
	}

	profile := os.Getenv("AWS_PROFILE")
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithSharedConfigProfile(profile),
		config.WithRegion(a.region),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	stsClient := sts.NewFromConfig(cfg)
	result, err := stsClient.GetCallerIdentity(context.TODO(), &sts.GetCallerIdentityInput{})
	if err != nil {
		return nil, fmt.Errorf("failed to get AWS caller identity: %w", err)
	}

	return &AWSProfile{
		Name:      profile,
		ARN:       *result.Arn,
		AccountID: *result.Account,
	}, nil
}

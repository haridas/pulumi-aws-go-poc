package main

import (
	"os"

	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func GetDeployRegion(ctx *pulumi.Context) string {
	var rolloutRegion string
	rolloutRegion, found := os.LookupEnv("ROLLOUT_REGION")
	if !found {
		rolloutRegion, _ = ctx.GetConfig("deploy:region")
	}
	return rolloutRegion
}

func GetRegionProvider(ctx *pulumi.Context, region string) *aws.Provider {
	provider, _ := aws.NewProvider(ctx,
		"aws-"+region+"-provider",
		&aws.ProviderArgs{
			Region: pulumi.String(region),
		})
	return provider
}

func GetAWSProvider(ctx *pulumi.Context) (*aws.Provider, string, error) {
	//region, err := aws.GetRegion(ctx, nil)
	// if err != nil {
	// 	return nil, "", errors.New("Failed to get the region" + err.Error())
	// }

	region := GetDeployRegion(ctx)
	provider := GetRegionProvider(ctx, region)
	return provider, region, nil
}

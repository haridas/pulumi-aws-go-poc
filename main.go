package main

import (
	"fmt"
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/s3"
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/sqs"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

type Infra struct {
	Cpu int
	Ram int
}

func GetInfra(ctx *pulumi.Context) *Infra {
	var infra Infra
	cfg := config.New(ctx, "us-west-1")
	cfg.RequireObject("infra", &infra)
	return &infra
}

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Create an AWS resource (S3 Bucket)
		regions := []string{"us-west-1"}

		in := GetInfra(ctx)
		fmt.Println("===> ", in.Cpu)

		//provider, region, err := GetAWSProvider(ctx)
		// if err != nil {
		// 	return errors.New("failed to create the provider")
		// }
		for _, region := range regions {
			provider := GetRegionProvider(ctx, region)
			bucket, err := s3.NewBucket(ctx, "my-bucket"+region, nil, pulumi.Provider(provider))

			if err != nil {
				return err
			}

			sqsQueue, err := sqs.NewQueue(ctx,
				"hari-pulumi-testqueue-"+region,
				&sqs.QueueArgs{
					ContentBasedDeduplication: pulumi.Bool(true),
					FifoQueue:                 pulumi.Bool(true),
				},
				pulumi.Provider(provider),
				pulumi.Aliases([]pulumi.Alias{{Name: pulumi.String("hari-pulumi-test-queue-" + region)}}),
				// pulumi.Protect(true),
			)
			if err != nil {
				return err
			}

			// Export the name of the bucket
			ctx.Export("bucketName-"+region, bucket.ID())
			ctx.Export("queueArn-"+region, sqsQueue.Arn)
		}

		return nil
	})
}

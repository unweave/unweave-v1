package awsprov

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate  . StsAPI
//counterfeiter:generate  . Ec2API
//counterfeiter:generate  . IamAPI

type StsAPI interface {
	GetCallerIdentity(ctx context.Context,
		params *sts.GetCallerIdentityInput,
		optFns ...func(*sts.Options)) (*sts.GetCallerIdentityOutput, error)
}

type IamAPI interface {
	CreatePolicy(ctx context.Context,
		params *iam.CreatePolicyInput,
		optFns ...func(*iam.Options)) (*iam.CreatePolicyOutput, error)

	CreateRole(ctx context.Context,
		params *iam.CreateRoleInput,
		optFns ...func(*iam.Options)) (*iam.CreateRoleOutput, error)

	GetInstanceProfile(ctx context.Context,
		params *iam.GetInstanceProfileInput,
		optFns ...func(*iam.Options)) (*iam.GetInstanceProfileOutput, error)

	AttachRolePolicy(ctx context.Context,
		params *iam.AttachRolePolicyInput,
		optFns ...func(*iam.Options)) (*iam.AttachRolePolicyOutput, error)

	CreateInstanceProfile(ctx context.Context,
		params *iam.CreateInstanceProfileInput,
		optFns ...func(*iam.Options)) (*iam.CreateInstanceProfileOutput, error)

	AddRoleToInstanceProfile(ctx context.Context,
		params *iam.AddRoleToInstanceProfileInput,
		optFns ...func(*iam.Options)) (*iam.AddRoleToInstanceProfileOutput, error)
}

type Ec2API interface {
	RunInstances(ctx context.Context,
		params *ec2.RunInstancesInput,
		optFns ...func(*ec2.Options)) (*ec2.RunInstancesOutput, error)

	DescribeInstanceStatus(ctx context.Context,
		params *ec2.DescribeInstanceStatusInput,
		optFns ...func(*ec2.Options)) (*ec2.DescribeInstanceStatusOutput, error)

	DescribeInstances(ctx context.Context,
		params *ec2.DescribeInstancesInput,
		optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error)

	DescribeInstanceTypes(ctx context.Context,
		params *ec2.DescribeInstanceTypesInput,
		optFns ...func(*ec2.Options)) (*ec2.DescribeInstanceTypesOutput, error)

	TerminateInstances(ctx context.Context,
		params *ec2.TerminateInstancesInput,
		optFns ...func(*ec2.Options)) (*ec2.TerminateInstancesOutput, error)

	CreateTags(ctx context.Context,
		params *ec2.CreateTagsInput,
		optFns ...func(*ec2.Options)) (*ec2.CreateTagsOutput, error)

	CreateVolume(ctx context.Context,
		params *ec2.CreateVolumeInput,
		optFns ...func(*ec2.Options)) (*ec2.CreateVolumeOutput, error)

	DeleteVolume(ctx context.Context,
		params *ec2.DeleteVolumeInput,
		optFns ...func(*ec2.Options)) (*ec2.DeleteVolumeOutput, error)

	ModifyVolume(ctx context.Context,
		params *ec2.ModifyVolumeInput,
		optFns ...func(*ec2.Options)) (*ec2.ModifyVolumeOutput, error)
}

func NewAwsApis(region, accessKey, secretKey string) (Ec2API, StsAPI, IamAPI, error) {
	creds := aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""))

	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion(region),
		config.WithCredentialsProvider(creds),
	)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("load aws config: %w", err)
	}

	ec2API := ec2.NewFromConfig(cfg)
	stsAPI := sts.NewFromConfig(cfg)
	iamAPI := iam.NewFromConfig(cfg)

	return ec2API, stsAPI, iamAPI, nil
}

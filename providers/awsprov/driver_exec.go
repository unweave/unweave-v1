package awsprov

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	iamtypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/rs/zerolog/log"
	"github.com/unweave/unweave/api/types"
	"github.com/unweave/unweave/providers/awsprov/internal/nodes"
	"github.com/unweave/unweave/services/execsrv"
	"github.com/unweave/unweave/tools/random"
)

type ExecDriver struct {
	userID string
	ec2API Ec2API
	stsAPI StsAPI
	iamAPI IamAPI
	region string
}

func NewExecDriverAPI(region, userID string, ec2API Ec2API, stsAPI StsAPI, iamAPI IamAPI) *ExecDriver {
	return &ExecDriver{
		region: region,
		userID: userID,
		ec2API: ec2API,
		stsAPI: stsAPI,
		iamAPI: iamAPI,
	}
}

func (d *ExecDriver) ExecCreate(
	ctx context.Context,
	project string,
	image string,
	spec types.HardwareSpec,
	network types.ExecNetwork,
	volumes []types.ExecVolume,
	pubKeys []string,
	region *string,
) (string, error) {
	log.Warn().Msgf("Ignoring hardware spec in aws driver")

	if network.HTTPService != nil {
		return "", errors.New("exposing an http service is not supported in this aws provider")
	}

	instanceType := nodes.NodeType(spec)
	if instanceType == "" {
		return "", fmt.Errorf("could not find node matching spec")
	}

	execID, err := newExecID()
	if err != nil {
		return "", fmt.Errorf("generate exec ID: %w", err)
	}

	uData, err := UserData(*region, pubKeys, volumes)
	if err != nil {
		return "", fmt.Errorf("failed to build user data: %w", err)
	}

	arn, err := d.setupIamPermissions(ctx)
	if err != nil {
		return "", fmt.Errorf("setup iam permissions: %w", err)
	}

	minMaxCount := int32(1)
	input := &ec2.RunInstancesInput{
		ImageId:           &image,
		InstanceType:      instanceType,
		MinCount:          &minMaxCount,
		MaxCount:          &minMaxCount,
		UserData:          &uData,
		TagSpecifications: d.tags(project, execID),
		Placement: &ec2types.Placement{
			AvailabilityZone: aws.String(d.region + "a"),
		},
		IamInstanceProfile: &ec2types.IamInstanceProfileSpecification{
			Arn: &arn,
		},
	}

	_, err = d.ec2API.RunInstances(ctx, input)
	if err != nil {
		return "", fmt.Errorf("failed to start instance: %w", err)
	}

	return execID, nil
}

const trustPolicy = `{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "sts:AssumeRole"
            ],
            "Principal": {
                "Service": [
                    "ec2.amazonaws.com"
                ]
            }
        }
    ]
}`

const rolePolicy = `{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "UnweaveEc2AttachVol",
            "Effect": "Allow",
            "Action": [
                "ec2:AttachVolume",
                "ec2:DescribeVolumes"
            ],
            "Resource": "*"
        }
    ]
}`

func (d *ExecDriver) setupIamPermissions(ctx context.Context) (string, error) {
	getipOut, err := d.iamAPI.GetInstanceProfile(
		ctx,
		&iam.GetInstanceProfileInput{InstanceProfileName: aws.String("UnweaveEc2ExecInstanceProfile")},
	)
	if err == nil {
		log.Debug().Msg("Instance profile already exists, skipping")

		return *getipOut.InstanceProfile.Arn, nil
	}

	var nse *iamtypes.NoSuchEntityException

	hasError := err != nil
	roleMissing := hasError && errors.As(err, &nse)

	if hasError && !roleMissing {
		return "", fmt.Errorf("get instance profile: %w", err)
	}

	createRoleInput := iam.CreateRoleInput{
		AssumeRolePolicyDocument: aws.String(trustPolicy),
		RoleName:                 aws.String("UnweaveEc2ExecRole"),
		Description:              aws.String("Role for Ec2Execs to assume"),
	}

	crOut, err := d.iamAPI.CreateRole(ctx, &createRoleInput)
	if err != nil {
		return "", fmt.Errorf("create role: %w", err)
	}

	cpOut, err := d.iamAPI.CreatePolicy(
		ctx,
		&iam.CreatePolicyInput{
			PolicyDocument: aws.String(rolePolicy),
			PolicyName:     aws.String("UnweaveEc2AttachVolume"),
		},
	)
	if err != nil {
		return "", fmt.Errorf("create policy: %w", err)
	}

	_, err = d.iamAPI.AttachRolePolicy(
		ctx,
		&iam.AttachRolePolicyInput{
			PolicyArn: cpOut.Policy.Arn,
			RoleName:  crOut.Role.RoleName,
		},
	)
	if err != nil {
		return "", fmt.Errorf("attach policy to role: %w", err)
	}

	cipOut, err := d.iamAPI.CreateInstanceProfile(
		ctx,
		&iam.CreateInstanceProfileInput{
			InstanceProfileName: aws.String("UnweaveEc2ExecInstanceProfile"),
		},
	)
	if err != nil {
		return "", fmt.Errorf("create instance profile: %w", err)
	}

	_, err = d.iamAPI.AddRoleToInstanceProfile(
		ctx,
		&iam.AddRoleToInstanceProfileInput{
			InstanceProfileName: cipOut.InstanceProfile.InstanceProfileName,
			RoleName:            crOut.Role.RoleName,
		},
	)
	if err != nil {
		return "", fmt.Errorf("add role to instance profile: %w", err)
	}

	return "", nil
}

func (d *ExecDriver) tags(project, execID string) []ec2types.TagSpecification {
	return []ec2types.TagSpecification{
		{
			ResourceType: ec2types.ResourceTypeInstance,
			Tags: []ec2types.Tag{
				{
					Key:   aws.String("unweave.io/project"),
					Value: &project,
				},
				{
					Key:   aws.String("unweave.io/user"),
					Value: &d.userID,
				},
				{
					Key:   aws.String("unweave.io/exec"),
					Value: &execID,
				},
			},
		},
	}
}

func (d *ExecDriver) ExecDriverName() string {
	return "aws"
}

func (d *ExecDriver) ExecGetStatus(ctx context.Context, execID string) (types.Status, error) {
	_, status, err := d.instanceState(ctx, execID)
	if err != nil {
		return types.StatusUnknown, fmt.Errorf("failed to get state: %w", err)
	}

	return status, nil
}

func (d *ExecDriver) instanceState(ctx context.Context, execID string) (string, types.Status, error) {
	instance, err := d.instance(ctx, execID)
	if err != nil {
		return "", types.StatusUnknown, fmt.Errorf("get instance: %w", err)
	}

	instanceState := instance.State.Name

	var status types.Status

	switch instanceState {
	case ec2types.InstanceStateNamePending:
		status = types.StatusInitializing
	case ec2types.InstanceStateNameRunning:
		// Check the health checks to make sure
		// the instance is both running and contactable
		status, err = d.instanceSummaryStatus(ctx, *instance.InstanceId)
		if err != nil {
			return "", types.StatusError, fmt.Errorf("summary status: %w", err)
		}
	case ec2types.InstanceStateNameShuttingDown,
		ec2types.InstanceStateNameTerminated,
		ec2types.InstanceStateNameStopping,
		ec2types.InstanceStateNameStopped:
		status = types.StatusTerminated
	default:
		status = types.StatusUnknown
	}

	return *instance.InstanceId, status, nil
}

func (d *ExecDriver) instanceSummaryStatus(ctx context.Context, instanceID string) (types.Status, error) {
	out, err := d.ec2API.DescribeInstanceStatus(
		ctx,
		&ec2.DescribeInstanceStatusInput{
			InstanceIds:         []string{instanceID},
			IncludeAllInstances: aws.Bool(true),
		},
	)
	if err != nil {
		return "", fmt.Errorf("describe status: %w", err)
	}

	if len(out.InstanceStatuses) == 0 {
		return "", errors.New("not found")
	}

	summaryStatus := out.InstanceStatuses[0].InstanceStatus.Status

	switch summaryStatus {
	case ec2types.SummaryStatusInitializing, ec2types.SummaryStatusInsufficientData:
		return types.StatusInitializing, nil
	case ec2types.SummaryStatusOk, ec2types.SummaryStatusNotApplicable:
		return types.StatusRunning, nil
	case ec2types.SummaryStatusImpaired:
		return types.StatusError, nil
	}

	return types.StatusUnknown, nil
}

func (d *ExecDriver) instance(ctx context.Context, execID string) (ec2types.Instance, error) {
	input := &ec2.DescribeInstancesInput{
		Filters: []ec2types.Filter{
			{
				Name:   aws.String("tag:unweave.io/exec"),
				Values: []string{execID},
			},
		},
	}

	output, err := d.ec2API.DescribeInstances(ctx, input)
	if err != nil {
		return ec2types.Instance{}, fmt.Errorf("failed to get exec instance: %w", err)
	}

	if len(output.Reservations) == 0 {
		return ec2types.Instance{}, fmt.Errorf("exec not found: %w", err)
	}

	instance := output.Reservations[0].Instances[0]

	return instance, nil
}

func (d *ExecDriver) ExecProvider() types.Provider {
	return types.AWSProvider
}

func (d *ExecDriver) ExecTerminate(ctx context.Context, execID string) error {
	instanceID, _, err := d.instanceState(ctx, execID)
	if err != nil {
		return fmt.Errorf("failed to get instance: %w", err)
	}

	_, err = d.ec2API.TerminateInstances(ctx, &ec2.TerminateInstancesInput{InstanceIds: []string{instanceID}})
	if err != nil {
		return fmt.Errorf("failed to terminate: %w", err)
	}

	return nil
}

func (d *ExecDriver) ExecSpec(_ context.Context, _ string) (types.HardwareSpec, error) {
	panic("not implemented")
}

func (d *ExecDriver) ExecStats(_ context.Context, _ string) (execsrv.Stats, error) {
	panic("not implemented")
}

// ExecPing pings the driver availability on behalf of a user. This can be used to
// check if the driver is configured correctly and healthy.
func (d *ExecDriver) ExecPing(ctx context.Context, _ *string) error {
	_, err := d.stsAPI.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		return fmt.Errorf("failed to ping: %w", err)
	}

	return nil
}

func (d *ExecDriver) ExecConnectionInfo(ctx context.Context, execID string) (types.ConnectionInfo, error) {
	instance, err := d.instance(ctx, execID)
	if err != nil {
		return types.ConnectionInfo{}, fmt.Errorf("connection info: %w", err)
	}

	return types.ConnectionInfo{
		Host: *instance.PublicIpAddress,
		User: "unweave", // we make this user in the userData on startup
		Port: 22,
	}, nil
}

func newExecID() (string, error) {
	str, err := random.GenerateRandomString(11)
	if err != nil {
		return "", fmt.Errorf("failed to generate random string, %w", err)
	}

	execID := "exc_" + strings.ToLower(str)

	return execID, nil
}

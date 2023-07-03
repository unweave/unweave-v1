package awsprov

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/unweave/unweave/api/types"
)

type VolumeDriver struct {
	userID string
	region string
	ec2Api Ec2API
}

func NewVolumeDriverAPI(region, userID string, ec2Api Ec2API) *VolumeDriver {
	return &VolumeDriver{
		region: region,
		userID: userID,
		ec2Api: ec2Api,
	}
}

func (v *VolumeDriver) VolumeCreate(ctx context.Context, projectID, name string, size int) (string, error) {
	input := &ec2.CreateVolumeInput{
		AvailabilityZone:  aws.String(v.region + "a"),
		Size:              aws.Int32(int32(size)),
		TagSpecifications: v.tags(projectID, name),
		VolumeType:        ec2types.VolumeTypeStandard,
	}

	out, err := v.ec2Api.CreateVolume(ctx, input)
	if err != nil {
		return "", fmt.Errorf("failed to create volume: %w", err)
	}

	return *out.VolumeId, nil
}

func (v *VolumeDriver) tags(project, name string) []ec2types.TagSpecification {
	return []ec2types.TagSpecification{
		{
			ResourceType: ec2types.ResourceTypeVolume,
			Tags: []ec2types.Tag{
				{
					Key:   aws.String("Name"),
					Value: &name,
				},
				{
					Key:   aws.String("unweave.io/project"),
					Value: &project,
				},
				{
					Key:   aws.String("unweave.io/user"),
					Value: &v.userID,
				},
			},
		},
	}
}

func (v *VolumeDriver) VolumeDelete(ctx context.Context, id string) error {
	_, err := v.ec2Api.DeleteVolume(ctx, &ec2.DeleteVolumeInput{VolumeId: aws.String(id)})
	if err != nil {
		// handle still attached errors?
		return fmt.Errorf("failed to delete volume: %w", err)
	}

	return nil
}

func (v *VolumeDriver) VolumeResize(ctx context.Context, id string, size int) error {
	_, err := v.ec2Api.ModifyVolume(ctx, &ec2.ModifyVolumeInput{VolumeId: &id, Size: aws.Int32(int32(size))})
	if err != nil {
		return fmt.Errorf("failed to resize: %w", err)
	}

	return nil
}

func (v *VolumeDriver) VolumeProvider() types.Provider        { return types.AWSProvider }
func (v *VolumeDriver) VolumeDriver(_ context.Context) string { return "aws" }

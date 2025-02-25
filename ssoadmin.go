package setlist

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssoadmin"
	"github.com/aws/aws-sdk-go-v2/service/ssoadmin/types"
)

// SsoInstance retrieves the SSO instance metadata.
// AWS SSO allows only a single instance per organization.
func SsoInstance(ctx context.Context, client *ssoadmin.Client) (types.InstanceMetadata, error) {
	resp, err := client.ListInstances(ctx, nil)
	if err != nil {
		return types.InstanceMetadata{}, err
	}

	if len(resp.Instances) == 0 {
		return types.InstanceMetadata{}, errors.New("SSO is not enabled.  No SSO instances exist.")
	}

	instance := resp.Instances[0]

	return instance, nil
}

// PermissionSets retrieves the list of permission sets provisioned to an
// account.
func PermissionSets(ctx context.Context, client *ssoadmin.Client, instanceArn string, accountId string) ([]types.PermissionSet, error) {
	permissionSets := []types.PermissionSet{}

	var permissionSetArns []string
	var token *string
	for {
		params := &ssoadmin.ListPermissionSetsProvisionedToAccountInput{
			InstanceArn: aws.String(instanceArn),
			AccountId:   aws.String(accountId),
			NextToken:   token,
		}
		resp, err := client.ListPermissionSetsProvisionedToAccount(ctx, params)

		if err != nil {
			return permissionSets, err
		}

		for _, i := range resp.PermissionSets {
			permissionSetArns = append(permissionSetArns, i)
		}

		if resp.NextToken == nil {
			break
		}

		token = resp.NextToken
	}

	// Retrieve detailed information for each permission set.
	for _, arn := range permissionSetArns {
		params := &ssoadmin.DescribePermissionSetInput{
			InstanceArn:      aws.String(instanceArn),
			PermissionSetArn: aws.String(arn),
		}
		resp, err := client.DescribePermissionSet(ctx, params)
		if err != nil {
			return permissionSets, err
		}

		permissionSets = append(permissionSets, *resp.PermissionSet)
	}

	return permissionSets, nil
}

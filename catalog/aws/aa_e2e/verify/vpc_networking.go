package verify

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/smithy-go"
	pkgerrors "github.com/pkg/errors"
)

// --- AwsVpcPeering ------------------------------------------------------------

// vpcPeeringVerifier verifies an AwsVpcPeering via
// DescribeVpcPeeringConnections, keyed on peering_connection_id.
//
// Deleted peerings linger describable in "deleted"/"deleting" states for
// hours, so absence treats those states as gone. NOTE the accept arm's
// destroy is a NO-OP at AWS (management is abandoned, the peering
// persists ACTIVE) - accept-arm lanes cannot assert absence and record
// that contract in their profile instead.
type vpcPeeringVerifier struct{}

func (*vpcPeeringVerifier) IDOutputKey() string { return "peering_connection_id" }

func (*vpcPeeringVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := vpcPeeringExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsvpcpeering verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awsvpcpeering %q not found after deploy", id)
	}
	return nil
}

func (*vpcPeeringVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := vpcPeeringExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsvpcpeering verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awsvpcpeering %q still exists after destroy", id)
	}
	return nil
}

func vpcPeeringExists(ctx context.Context, cfg aws.Config, id, region string) (bool, error) {
	client := ec2.NewFromConfig(cfg, func(o *ec2.Options) {
		if region != "" {
			o.Region = region
		}
	})
	out, err := client.DescribeVpcPeeringConnections(ctx, &ec2.DescribeVpcPeeringConnectionsInput{
		VpcPeeringConnectionIds: []string{id},
	})
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) && apiErr.ErrorCode() == "InvalidVpcPeeringConnectionID.NotFound" {
			return false, nil
		}
		return false, err
	}
	for _, connection := range out.VpcPeeringConnections {
		if connection.Status == nil {
			continue
		}
		switch connection.Status.Code {
		case ec2types.VpcPeeringConnectionStateReasonCodeDeleted,
			ec2types.VpcPeeringConnectionStateReasonCodeDeleting:
			// A deleted peering lingers describable - treat as absent.
			continue
		default:
			return true, nil
		}
	}
	return false, nil
}

// --- AwsNetworkAcl --------------------------------------------------------------

// networkAclVerifier verifies an AwsNetworkAcl via DescribeNetworkAcls,
// keyed on network_acl_id. The in-line rules and subnet associations are
// children of the ACL and cannot outlive it.
type networkAclVerifier struct{}

func (*networkAclVerifier) IDOutputKey() string { return "network_acl_id" }

func (*networkAclVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := networkAclExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsnetworkacl verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awsnetworkacl %q not found after deploy", id)
	}
	return nil
}

func (*networkAclVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := networkAclExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsnetworkacl verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awsnetworkacl %q still exists after destroy", id)
	}
	return nil
}

func networkAclExists(ctx context.Context, cfg aws.Config, id, region string) (bool, error) {
	client := ec2.NewFromConfig(cfg, func(o *ec2.Options) {
		if region != "" {
			o.Region = region
		}
	})
	out, err := client.DescribeNetworkAcls(ctx, &ec2.DescribeNetworkAclsInput{
		NetworkAclIds: []string{id},
	})
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) && apiErr.ErrorCode() == "InvalidNetworkAclID.NotFound" {
			return false, nil
		}
		return false, err
	}
	return len(out.NetworkAcls) > 0, nil
}

// --- AwsManagedPrefixList --------------------------------------------------------

// managedPrefixListVerifier verifies an AwsManagedPrefixList via
// DescribeManagedPrefixLists, keyed on prefix_list_id. Deletes drain
// through "delete-in-progress"/"delete-complete" states - both count as
// absent.
type managedPrefixListVerifier struct{}

func (*managedPrefixListVerifier) IDOutputKey() string { return "prefix_list_id" }

func (*managedPrefixListVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := managedPrefixListExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsmanagedprefixlist verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awsmanagedprefixlist %q not found after deploy", id)
	}
	return nil
}

func (*managedPrefixListVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := managedPrefixListExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsmanagedprefixlist verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awsmanagedprefixlist %q still exists after destroy", id)
	}
	return nil
}

func managedPrefixListExists(ctx context.Context, cfg aws.Config, id, region string) (bool, error) {
	client := ec2.NewFromConfig(cfg, func(o *ec2.Options) {
		if region != "" {
			o.Region = region
		}
	})
	out, err := client.DescribeManagedPrefixLists(ctx, &ec2.DescribeManagedPrefixListsInput{
		PrefixListIds: []string{id},
	})
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) && apiErr.ErrorCode() == "InvalidPrefixListID.NotFound" {
			return false, nil
		}
		return false, err
	}
	for _, prefixList := range out.PrefixLists {
		switch prefixList.State {
		case ec2types.PrefixListStateDeleteComplete, ec2types.PrefixListStateDeleteInProgress:
			continue
		default:
			return true, nil
		}
	}
	return false, nil
}

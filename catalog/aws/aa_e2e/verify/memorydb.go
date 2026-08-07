package verify

import (
	"context"
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/memorydb"
	memorydbtypes "github.com/aws/aws-sdk-go-v2/service/memorydb/types"
	pkgerrors "github.com/pkg/errors"
)

// memorydbClusterVerifier verifies an AwsMemorydbCluster via DescribeClusters.
// MemoryDB deletes are slow (the provider polls only after a 5-minute initial
// wait), and a destroyed cluster lingers in the "deleting" state before the
// API returns the typed ClusterNotFoundFault -- both mean "absent" for
// verification purposes (the NAT-gateway lifecycle class).
type memorydbClusterVerifier struct{}

func (*memorydbClusterVerifier) IDOutputKey() string { return "cluster_name" }

func (*memorydbClusterVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := memorydbClusterExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsmemorydbcluster verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awsmemorydbcluster %q not found after deploy", id)
	}
	return nil
}

func (*memorydbClusterVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := memorydbClusterExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsmemorydbcluster verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awsmemorydbcluster %q still exists after destroy", id)
	}
	return nil
}

func memorydbClusterExists(ctx context.Context, cfg aws.Config, name, region string) (bool, error) {
	client := memorydbClient(cfg, region)
	out, err := client.DescribeClusters(ctx, &memorydb.DescribeClustersInput{ClusterName: &name})
	if err != nil {
		var notFound *memorydbtypes.ClusterNotFoundFault
		if errors.As(err, &notFound) {
			return false, nil
		}
		return false, err
	}
	for _, c := range out.Clusters {
		if aws.ToString(c.Name) != name {
			continue
		}
		// A cluster on its way out counts as absent.
		if strings.EqualFold(aws.ToString(c.Status), "deleting") {
			return false, nil
		}
		return true, nil
	}
	return false, nil
}

// memorydbUserVerifier verifies an AwsMemorydbUser via DescribeUsers. User
// deletion is quick; a missing user returns the typed UserNotFoundFault.
type memorydbUserVerifier struct{}

func (*memorydbUserVerifier) IDOutputKey() string { return "user_name" }

func (*memorydbUserVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := memorydbUserExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsmemorydbuser verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awsmemorydbuser %q not found after deploy", id)
	}
	return nil
}

func (*memorydbUserVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := memorydbUserExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsmemorydbuser verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awsmemorydbuser %q still exists after destroy", id)
	}
	return nil
}

func memorydbUserExists(ctx context.Context, cfg aws.Config, name, region string) (bool, error) {
	client := memorydbClient(cfg, region)
	out, err := client.DescribeUsers(ctx, &memorydb.DescribeUsersInput{UserName: &name})
	if err != nil {
		var notFound *memorydbtypes.UserNotFoundFault
		if errors.As(err, &notFound) {
			return false, nil
		}
		return false, err
	}
	for _, u := range out.Users {
		if aws.ToString(u.Name) != name {
			continue
		}
		// A user pending deletion counts as absent.
		if strings.EqualFold(aws.ToString(u.Status), "deleting") {
			return false, nil
		}
		return true, nil
	}
	return false, nil
}

// memorydbAclVerifier verifies an AwsMemorydbAcl via DescribeACLs. ACL
// deletion is quick; a missing ACL returns the typed ACLNotFoundFault.
type memorydbAclVerifier struct{}

func (*memorydbAclVerifier) IDOutputKey() string { return "acl_name" }

func (*memorydbAclVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := memorydbAclExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsmemorydbacl verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awsmemorydbacl %q not found after deploy", id)
	}
	return nil
}

func (*memorydbAclVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := memorydbAclExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsmemorydbacl verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awsmemorydbacl %q still exists after destroy", id)
	}
	return nil
}

func memorydbAclExists(ctx context.Context, cfg aws.Config, name, region string) (bool, error) {
	client := memorydbClient(cfg, region)
	out, err := client.DescribeACLs(ctx, &memorydb.DescribeACLsInput{ACLName: &name})
	if err != nil {
		var notFound *memorydbtypes.ACLNotFoundFault
		if errors.As(err, &notFound) {
			return false, nil
		}
		return false, err
	}
	for _, a := range out.ACLs {
		if aws.ToString(a.Name) != name {
			continue
		}
		if strings.EqualFold(aws.ToString(a.Status), "deleting") {
			return false, nil
		}
		return true, nil
	}
	return false, nil
}

func memorydbClient(cfg aws.Config, region string) *memorydb.Client {
	return memorydb.NewFromConfig(cfg, func(o *memorydb.Options) {
		if region != "" {
			o.Region = region
		}
	})
}

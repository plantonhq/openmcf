package verify

import (
	"context"
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticache"
	"github.com/aws/aws-sdk-go-v2/service/elasticache/types"
	pkgerrors "github.com/pkg/errors"
)

// elasticacheReplicationGroupVerifier verifies AwsRedisElasticache via
// DescribeReplicationGroups. A group mid-deletion stays describable with a
// "deleting" status before ElastiCache stops returning it — the NAT-gateway
// lifecycle class — so existence is "described AND not deleting/deleted".
type elasticacheReplicationGroupVerifier struct{}

func (*elasticacheReplicationGroupVerifier) IDOutputKey() string { return "replication_group_id" }

func (*elasticacheReplicationGroupVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := replicationGroupExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsrediselasticache verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awsrediselasticache replication group %q not found after deploy", id)
	}
	return nil
}

func (*elasticacheReplicationGroupVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := replicationGroupExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsrediselasticache verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awsrediselasticache replication group %q still exists after destroy", id)
	}
	return nil
}

func replicationGroupExists(ctx context.Context, cfg aws.Config, id, region string) (bool, error) {
	client := elasticacheClient(cfg, region)
	out, err := client.DescribeReplicationGroups(ctx, &elasticache.DescribeReplicationGroupsInput{
		ReplicationGroupId: &id,
	})
	if err != nil {
		var notFound *types.ReplicationGroupNotFoundFault
		if errors.As(err, &notFound) {
			return false, nil
		}
		return false, err
	}
	for _, rg := range out.ReplicationGroups {
		if rg.Status == nil {
			continue
		}
		switch strings.ToLower(*rg.Status) {
		case "deleting", "deleted":
			continue
		default:
			return true, nil
		}
	}
	return false, nil
}

// elasticacheClusterVerifier verifies AwsMemcachedElasticache via
// DescribeCacheClusters on the cluster id output.
type elasticacheClusterVerifier struct{}

func (*elasticacheClusterVerifier) IDOutputKey() string { return "cluster_id" }

func (*elasticacheClusterVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := cacheClusterExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsmemcachedelasticache verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awsmemcachedelasticache cluster %q not found after deploy", id)
	}
	return nil
}

func (*elasticacheClusterVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := cacheClusterExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsmemcachedelasticache verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awsmemcachedelasticache cluster %q still exists after destroy", id)
	}
	return nil
}

func cacheClusterExists(ctx context.Context, cfg aws.Config, id, region string) (bool, error) {
	client := elasticacheClient(cfg, region)
	out, err := client.DescribeCacheClusters(ctx, &elasticache.DescribeCacheClustersInput{
		CacheClusterId: &id,
	})
	if err != nil {
		var notFound *types.CacheClusterNotFoundFault
		if errors.As(err, &notFound) {
			return false, nil
		}
		return false, err
	}
	for _, cluster := range out.CacheClusters {
		if cluster.CacheClusterStatus == nil {
			continue
		}
		switch strings.ToLower(*cluster.CacheClusterStatus) {
		case "deleting", "deleted":
			continue
		default:
			return true, nil
		}
	}
	return false, nil
}

// elasticacheServerlessCacheVerifier verifies AwsServerlessElasticache via
// DescribeServerlessCaches on the cache name output.
type elasticacheServerlessCacheVerifier struct{}

func (*elasticacheServerlessCacheVerifier) IDOutputKey() string { return "name" }

func (*elasticacheServerlessCacheVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := serverlessCacheExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsserverlesselasticache verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awsserverlesselasticache cache %q not found after deploy", id)
	}
	return nil
}

func (*elasticacheServerlessCacheVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := serverlessCacheExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsserverlesselasticache verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awsserverlesselasticache cache %q still exists after destroy", id)
	}
	return nil
}

func serverlessCacheExists(ctx context.Context, cfg aws.Config, id, region string) (bool, error) {
	client := elasticacheClient(cfg, region)
	out, err := client.DescribeServerlessCaches(ctx, &elasticache.DescribeServerlessCachesInput{
		ServerlessCacheName: &id,
	})
	if err != nil {
		var notFound *types.ServerlessCacheNotFoundFault
		if errors.As(err, &notFound) {
			return false, nil
		}
		return false, err
	}
	for _, cache := range out.ServerlessCaches {
		if cache.Status == nil {
			continue
		}
		switch strings.ToLower(*cache.Status) {
		case "deleting", "deleted":
			continue
		default:
			return true, nil
		}
	}
	return false, nil
}

// elasticacheUserVerifier verifies AwsElasticacheUser via DescribeUsers.
type elasticacheUserVerifier struct{}

func (*elasticacheUserVerifier) IDOutputKey() string { return "user_id" }

func (*elasticacheUserVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := elasticacheUserExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awselasticacheuser verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awselasticacheuser %q not found after deploy", id)
	}
	return nil
}

func (*elasticacheUserVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := elasticacheUserExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awselasticacheuser verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awselasticacheuser %q still exists after destroy", id)
	}
	return nil
}

func elasticacheUserExists(ctx context.Context, cfg aws.Config, id, region string) (bool, error) {
	client := elasticacheClient(cfg, region)
	out, err := client.DescribeUsers(ctx, &elasticache.DescribeUsersInput{
		UserId: &id,
	})
	if err != nil {
		var notFound *types.UserNotFoundFault
		if errors.As(err, &notFound) {
			return false, nil
		}
		return false, err
	}
	for _, user := range out.Users {
		if user.Status == nil {
			continue
		}
		switch strings.ToLower(*user.Status) {
		case "deleting", "deleted":
			continue
		default:
			return true, nil
		}
	}
	return false, nil
}

// elasticacheUserGroupVerifier verifies AwsElasticacheUserGroup via
// DescribeUserGroups.
type elasticacheUserGroupVerifier struct{}

func (*elasticacheUserGroupVerifier) IDOutputKey() string { return "user_group_id" }

func (*elasticacheUserGroupVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := elasticacheUserGroupExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awselasticacheusergroup verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awselasticacheusergroup %q not found after deploy", id)
	}
	return nil
}

func (*elasticacheUserGroupVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := elasticacheUserGroupExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awselasticacheusergroup verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awselasticacheusergroup %q still exists after destroy", id)
	}
	return nil
}

func elasticacheUserGroupExists(ctx context.Context, cfg aws.Config, id, region string) (bool, error) {
	client := elasticacheClient(cfg, region)
	out, err := client.DescribeUserGroups(ctx, &elasticache.DescribeUserGroupsInput{
		UserGroupId: &id,
	})
	if err != nil {
		var notFound *types.UserGroupNotFoundFault
		if errors.As(err, &notFound) {
			return false, nil
		}
		return false, err
	}
	for _, group := range out.UserGroups {
		if group.Status == nil {
			continue
		}
		switch strings.ToLower(*group.Status) {
		case "deleting", "deleted":
			continue
		default:
			return true, nil
		}
	}
	return false, nil
}

func elasticacheClient(cfg aws.Config, region string) *elasticache.Client {
	return elasticache.NewFromConfig(cfg, func(o *elasticache.Options) {
		if region != "" {
			o.Region = region
		}
	})
}

package pulumiawsnativeprovider

import (
	"context"
	"testing"

	awsprovider "github.com/plantonhq/planton/catalog/aws"
	awsnative "github.com/pulumi/pulumi-aws-native/sdk/go/aws"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

// pulumi-aws-native's provider-block surface is narrower than the classic
// provider's; these tests pin the single-hop mapping and every hard-error
// degradation (a silent one would deploy as the wrong identity or with
// silently different semantics).

func TestBuildProviderArgs_ProviderBlock_SingleHop(t *testing.T) {
	cfg := &awsprovider.AwsProviderConfig{
		AccountId: "123456789012",
		Region:    "us-east-1",
		AssumeRoleChain: []*awsprovider.AwsAssumeRole{
			{
				RoleArn:     "arn:aws:iam::222222222222:role/deploy",
				SessionName: "planton",
				ExternalId:  "expected-external-id",
				Duration:    "30m",
			},
		},
		DefaultTags: &awsprovider.AwsDefaultTags{Tags: map[string]string{"CostCenter": "eng"}},
		MaxRetries:  proto.Int32(5),
	}

	args, err := buildProviderArgs(context.Background(), cfg, "us-east-1", failingResolver(t))
	require.NoError(t, err)

	hop, ok := args.AssumeRole.(*awsnative.ProviderAssumeRoleArgs)
	require.True(t, ok, "AssumeRole must be *ProviderAssumeRoleArgs")
	assert.Equal(t, pulumi.String("arn:aws:iam::222222222222:role/deploy"), hop.RoleArn)
	assert.Equal(t, pulumi.String("planton"), hop.SessionName)
	assert.Equal(t, pulumi.String("expected-external-id"), hop.ExternalId)
	// Native takes seconds where Terraform/classic take a duration string.
	assert.Equal(t, pulumi.Int(1800), hop.DurationSeconds)

	defaultTags, ok := args.DefaultTags.(awsnative.ProviderDefaultTagsArgs)
	require.True(t, ok)
	assert.Equal(t, pulumi.ToStringMap(map[string]string{"CostCenter": "eng"}), defaultTags.Tags)
	assert.Equal(t, pulumi.Int(5), args.MaxRetries)
}

func TestBuildProviderArgs_ProviderBlock_ChainLongerThanOne_Errors(t *testing.T) {
	cfg := &awsprovider.AwsProviderConfig{
		AccountId: "123456789012",
		Region:    "us-east-1",
		AssumeRoleChain: []*awsprovider.AwsAssumeRole{
			{RoleArn: "arn:aws:iam::111111111111:role/intermediate"},
			{RoleArn: "arn:aws:iam::222222222222:role/deploy"},
		},
	}

	_, err := buildProviderArgs(context.Background(), cfg, "us-east-1", failingResolver(t))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "single assume_role hop")
	assert.Contains(t, err.Error(), "chain of 2")
}

func TestBuildProviderArgs_ProviderBlock_SourceIdentity_Errors(t *testing.T) {
	cfg := &awsprovider.AwsProviderConfig{
		AccountId: "123456789012",
		Region:    "us-east-1",
		AssumeRoleChain: []*awsprovider.AwsAssumeRole{
			{RoleArn: "arn:aws:iam::222222222222:role/deploy", SourceIdentity: "engineer"},
		},
	}

	_, err := buildProviderArgs(context.Background(), cfg, "us-east-1", failingResolver(t))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "source_identity")
}

func TestBuildProviderArgs_ProviderBlock_RetryMode_Errors(t *testing.T) {
	cfg := &awsprovider.AwsProviderConfig{
		AccountId: "123456789012",
		Region:    "us-east-1",
		RetryMode: "adaptive",
	}

	_, err := buildProviderArgs(context.Background(), cfg, "us-east-1", failingResolver(t))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no retry_mode support")
}

func TestBuildProviderArgs_ProviderBlock_Endpoints(t *testing.T) {
	// cloudformation is in native's small endpoint vocabulary; s3 is not.
	cfg := &awsprovider.AwsProviderConfig{
		AccountId: "123456789012",
		Region:    "us-east-1",
		Endpoints: map[string]string{"cloudformation": "https://cfn.internal.example.com"},
	}
	args, err := buildProviderArgs(context.Background(), cfg, "us-east-1", failingResolver(t))
	require.NoError(t, err)
	endpointsArray, ok := args.Endpoints.(awsnative.ProviderEndpointArray)
	require.True(t, ok)
	require.Len(t, endpointsArray, 1)
	endpoint := endpointsArray[0].(awsnative.ProviderEndpointArgs)
	assert.Equal(t, pulumi.StringPtrInput(pulumi.String("https://cfn.internal.example.com")), endpoint.Cloudformation)

	cfg.Endpoints = map[string]string{"s3": "https://s3.internal.example.com"}
	_, err = buildProviderArgs(context.Background(), cfg, "us-east-1", failingResolver(t))
	require.Error(t, err)
	assert.Contains(t, err.Error(), `"s3"`)
	assert.Contains(t, err.Error(), "small subset")
}

func TestBuildProviderArgs_ProviderBlock_InvalidDuration_Errors(t *testing.T) {
	cfg := &awsprovider.AwsProviderConfig{
		AccountId: "123456789012",
		Region:    "us-east-1",
		AssumeRoleChain: []*awsprovider.AwsAssumeRole{
			{RoleArn: "arn:aws:iam::222222222222:role/deploy", Duration: "one-hour"},
		},
	}

	_, err := buildProviderArgs(context.Background(), cfg, "us-east-1", failingResolver(t))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid assume_role duration")
}

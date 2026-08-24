package pulumiawsprovider

import (
	"context"
	"testing"

	awsprovider "github.com/plantonhq/planton/catalog/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

// The provider-block surface (assume-role chain, default tags, endpoints,
// retries) composes with every base-credential mode; these tests pin the
// classic builder's full-surface mapping and its loud unknown-endpoint refusal.

func TestBuildProviderArgs_ProviderBlock_FullSurface(t *testing.T) {
	cfg := &awsprovider.AwsProviderConfig{
		AccountId:       "123456789012",
		AccessKeyId:     "AKIAIOSFODNN7EXAMPLE",
		SecretAccessKey: "wJalrXUtnFEMIK7MDENGbPxRfiCYEXAMPLEKEY123",
		Region:          "us-east-1",
		AssumeRoleChain: []*awsprovider.AwsAssumeRole{
			{RoleArn: "arn:aws:iam::111111111111:role/intermediate", SessionName: "hop-1"},
			{
				RoleArn:           "arn:aws:iam::222222222222:role/deploy",
				ExternalId:        "expected-external-id",
				Duration:          "1h",
				Policy:            `{"Version":"2012-10-17"}`,
				PolicyArns:        []string{"arn:aws:iam::aws:policy/ReadOnlyAccess"},
				Tags:              map[string]string{"Team": "platform"},
				TransitiveTagKeys: []string{"Team"},
				SourceIdentity:    "platform-engineer",
			},
		},
		DefaultTags: &awsprovider.AwsDefaultTags{
			Tags: map[string]string{"CostCenter": "eng"},
		},
		Endpoints:  map[string]string{"sts": "https://sts.internal.example.com"},
		MaxRetries: proto.Int32(10),
		RetryMode:  "adaptive",
	}

	args, err := buildProviderArgs(context.Background(), cfg, "us-east-1", failingResolver(t))
	require.NoError(t, err)

	// The chain maps hop-for-hop, in order (chained evaluation).
	chain, ok := args.AssumeRoles.(aws.ProviderAssumeRoleArray)
	require.True(t, ok, "AssumeRoles must be a ProviderAssumeRoleArray")
	require.Len(t, chain, 2)
	first := chain[0].(aws.ProviderAssumeRoleArgs)
	assert.Equal(t, pulumi.String("arn:aws:iam::111111111111:role/intermediate"), first.RoleArn)
	assert.Equal(t, pulumi.String("hop-1"), first.SessionName)
	assert.Nil(t, first.ExternalId)
	second := chain[1].(aws.ProviderAssumeRoleArgs)
	assert.Equal(t, pulumi.String("arn:aws:iam::222222222222:role/deploy"), second.RoleArn)
	assert.Equal(t, pulumi.String("expected-external-id"), second.ExternalId)
	assert.Equal(t, pulumi.String("1h"), second.Duration)
	assert.Equal(t, pulumi.String(`{"Version":"2012-10-17"}`), second.Policy)
	assert.Equal(t, pulumi.ToStringArray([]string{"arn:aws:iam::aws:policy/ReadOnlyAccess"}), second.PolicyArns)
	assert.Equal(t, pulumi.ToStringMap(map[string]string{"Team": "platform"}), second.Tags)
	assert.Equal(t, pulumi.ToStringArray([]string{"Team"}), second.TransitiveTagKeys)
	assert.Equal(t, pulumi.String("platform-engineer"), second.SourceIdentity)

	// The web-identity slot stays untouched: the chain rides AssumeRoles, so a
	// keyless base credential and a chain can coexist.
	assert.Nil(t, args.AssumeRoleWithWebIdentity)

	defaultTags, ok := args.DefaultTags.(aws.ProviderDefaultTagsArgs)
	require.True(t, ok)
	assert.Equal(t, pulumi.ToStringMap(map[string]string{"CostCenter": "eng"}), defaultTags.Tags)

	endpointsArray, ok := args.Endpoints.(aws.ProviderEndpointArray)
	require.True(t, ok)
	require.Len(t, endpointsArray, 1)
	endpoint := endpointsArray[0].(aws.ProviderEndpointArgs)
	assert.Equal(t, pulumi.StringPtrInput(pulumi.String("https://sts.internal.example.com")), endpoint.Sts)

	assert.Equal(t, pulumi.Int(10), args.MaxRetries)
	assert.Equal(t, pulumi.String("adaptive"), args.RetryMode)

	// The base credentials still landed (the surface composes, never replaces).
	assert.Equal(t, pulumi.String("AKIAIOSFODNN7EXAMPLE"), args.AccessKey)
}

func TestBuildProviderArgs_ProviderBlock_UnknownEndpointService_Errors(t *testing.T) {
	cfg := &awsprovider.AwsProviderConfig{
		AccountId: "123456789012",
		Region:    "us-east-1",
		Endpoints: map[string]string{"not-a-service": "https://example.com"},
	}

	_, err := buildProviderArgs(context.Background(), cfg, "us-east-1", failingResolver(t))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not-a-service")
	assert.Contains(t, err.Error(), "no such endpoint attribute")
}

func TestBuildProviderArgs_ProviderBlock_MaxRetriesZero_IsExplicit(t *testing.T) {
	// optional int32 presence: zero means "no retries", distinct from unset.
	cfg := &awsprovider.AwsProviderConfig{
		AccountId:  "123456789012",
		Region:     "us-east-1",
		MaxRetries: proto.Int32(0),
	}

	args, err := buildProviderArgs(context.Background(), cfg, "us-east-1", failingResolver(t))
	require.NoError(t, err)
	assert.Equal(t, pulumi.Int(0), args.MaxRetries)

	// And unset stays nil.
	cfg.MaxRetries = nil
	args, err = buildProviderArgs(context.Background(), cfg, "us-east-1", failingResolver(t))
	require.NoError(t, err)
	assert.Nil(t, args.MaxRetries)
}

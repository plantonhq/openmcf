// EKS bearer-token minting: a presigned STS GetCallerIdentity URL bound to one cluster.
package kubetoken

import (
	"context"
	"encoding/base64"
	"time"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	smithyhttp "github.com/aws/smithy-go/transport/http"
	"github.com/pkg/errors"
)

const (
	// eksTokenPrefix is the encoding EKS mandates: the prefix followed by the
	// base64url-encoded presigned URL (aws-iam-authenticator token format v1).
	eksTokenPrefix = "k8s-aws-v1."

	// clusterIDHeader binds the token to one cluster: the header is part of the
	// signature, and the EKS API server rejects tokens signed for a different cluster.
	clusterIDHeader = "X-K8s-Aws-Id"

	// presignExpiresSeconds is signed into the URL as X-Amz-Expires. EKS enforces its
	// own ~15-minute acceptance window regardless; 60 is the value aws-iam-authenticator
	// signs and AWS documents for client implementations.
	presignExpiresSeconds = "60"

	// eksTokenLifetime is the honest validity reported to callers: the EKS server-side
	// 15-minute window minus a 1-minute refresh margin (aws-iam-authenticator convention),
	// so a consumer refreshing on expiry never presents a token the server just rejected.
	eksTokenLifetime = 14 * time.Minute
)

// EksTokenOptions identifies the cluster a token is minted for and, optionally, the
// static AWS credentials that sign it. When the static fields are empty the ambient
// AWS credential chain of the process signs instead (env, shared profile, instance role).
type EksTokenOptions struct {
	ClusterName string
	Region      string

	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string
}

// MintEksToken presigns an STS GetCallerIdentity request with the cluster-binding
// header and encodes it as an EKS bearer token. Presigning is pure local signing --
// no network call is made -- which is also what makes this fully unit-testable.
func MintEksToken(ctx context.Context, opts EksTokenOptions) (Token, error) {
	if opts.ClusterName == "" || opts.Region == "" {
		return Token{}, errors.New("cluster name and region are required to mint an EKS token")
	}

	loadOpts := []func(*awsconfig.LoadOptions) error{awsconfig.WithRegion(opts.Region)}
	if opts.AccessKeyID != "" {
		loadOpts = append(loadOpts, awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(opts.AccessKeyID, opts.SecretAccessKey, opts.SessionToken)))
	}
	cfg, err := awsconfig.LoadDefaultConfig(ctx, loadOpts...)
	if err != nil {
		return Token{}, errors.Wrap(err, "loading AWS config for EKS token minting")
	}

	mintedAt := time.Now()
	presigned, err := sts.NewPresignClient(sts.NewFromConfig(cfg)).PresignGetCallerIdentity(ctx,
		&sts.GetCallerIdentityInput{},
		func(po *sts.PresignOptions) {
			po.ClientOptions = append(po.ClientOptions, func(o *sts.Options) {
				o.APIOptions = append(o.APIOptions,
					smithyhttp.SetHeaderValue(clusterIDHeader, opts.ClusterName),
					smithyhttp.SetHeaderValue("X-Amz-Expires", presignExpiresSeconds),
				)
			})
		})
	if err != nil {
		return Token{}, errors.Wrap(err, "presigning STS GetCallerIdentity for EKS token")
	}

	return Token{
		Value:     eksTokenPrefix + base64.RawURLEncoding.EncodeToString([]byte(presigned.URL)),
		ExpiresAt: mintedAt.Add(eksTokenLifetime),
	}, nil
}

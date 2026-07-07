package verify

import (
	"context"
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/wafv2"
	wafv2types "github.com/aws/aws-sdk-go-v2/service/wafv2/types"
	pkgerrors "github.com/pkg/errors"
)

// WAFv2 resources are addressed by a three-part key (id + name + scope), all
// of which the resource ARN encodes:
//
//	arn:aws:wafv2:<region>:<account>:<regional|global>/<type>/<name>/<id>
//
// so every WAF verifier keys on the ARN output and derives the Get* call's
// parameters from it. The "global" ARN segment is the CLOUDFRONT scope (and
// pins the API endpoint to us-east-1); "regional" resources verify in their
// own region. Deletion is synchronous — a deleted resource returns the typed
// WAFNonexistentItemException immediately.

// wafv2ArnParts extracts (scope, name, id, region) from a WAFv2 resource ARN.
func wafv2ArnParts(arn string) (scope wafv2types.Scope, name, id, region string, err error) {
	// arn:aws:wafv2:us-west-2:123456789012:regional/webacl/<name>/<id>
	arnFields := strings.SplitN(arn, ":", 6)
	if len(arnFields) != 6 {
		return "", "", "", "", pkgerrors.Errorf("unexpected WAFv2 ARN shape %q", arn)
	}
	region = arnFields[3]
	resourceFields := strings.Split(arnFields[5], "/")
	if len(resourceFields) != 4 {
		return "", "", "", "", pkgerrors.Errorf("unexpected WAFv2 ARN resource shape %q", arn)
	}
	scope = wafv2types.ScopeRegional
	if resourceFields[0] == "global" {
		scope = wafv2types.ScopeCloudfront
		// CLOUDFRONT-scope calls must go to the WAF global endpoint.
		region = "us-east-1"
	}
	return scope, resourceFields[2], resourceFields[3], region, nil
}

// wafv2Client builds a client pinned to the region the ARN dictates (the
// harness's ambient region may differ for CLOUDFRONT-scope resources).
func wafv2Client(cfg aws.Config, region string) *wafv2.Client {
	cfg.Region = region
	return wafv2.NewFromConfig(cfg)
}

// isWafNotFound reports the typed absent signal.
func isWafNotFound(err error) bool {
	var notFound *wafv2types.WAFNonexistentItemException
	return errors.As(err, &notFound)
}

// wafWebAclVerifier verifies an AwsWafWebAcl via GetWebACL, keyed on the
// web_acl_arn output.
type wafWebAclVerifier struct{}

func (*wafWebAclVerifier) IDOutputKey() string { return "web_acl_arn" }

func (*wafWebAclVerifier) VerifyExists(ctx context.Context, cfg aws.Config, arn, region string) error {
	exists, err := wafWebAclExists(ctx, cfg, arn)
	if err != nil {
		return pkgerrors.Wrapf(err, "awswafwebacl verify-exists failed for %q", arn)
	}
	if !exists {
		return pkgerrors.Errorf("awswafwebacl %q not found after deploy", arn)
	}
	return nil
}

func (*wafWebAclVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, arn, region string) error {
	exists, err := wafWebAclExists(ctx, cfg, arn)
	if err != nil {
		return pkgerrors.Wrapf(err, "awswafwebacl verify-absent failed for %q", arn)
	}
	if exists {
		return pkgerrors.Errorf("awswafwebacl %q still exists after destroy", arn)
	}
	return nil
}

func wafWebAclExists(ctx context.Context, cfg aws.Config, arn string) (bool, error) {
	scope, name, id, region, err := wafv2ArnParts(arn)
	if err != nil {
		return false, err
	}
	_, err = wafv2Client(cfg, region).GetWebACL(ctx, &wafv2.GetWebACLInput{
		Id:    aws.String(id),
		Name:  aws.String(name),
		Scope: scope,
	})
	if err != nil {
		if isWafNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// wafIpSetVerifier verifies an AwsWafIpSet via GetIPSet, keyed on the
// ip_set_arn output.
type wafIpSetVerifier struct{}

func (*wafIpSetVerifier) IDOutputKey() string { return "ip_set_arn" }

func (*wafIpSetVerifier) VerifyExists(ctx context.Context, cfg aws.Config, arn, region string) error {
	exists, err := wafIpSetExists(ctx, cfg, arn)
	if err != nil {
		return pkgerrors.Wrapf(err, "awswafipset verify-exists failed for %q", arn)
	}
	if !exists {
		return pkgerrors.Errorf("awswafipset %q not found after deploy", arn)
	}
	return nil
}

func (*wafIpSetVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, arn, region string) error {
	exists, err := wafIpSetExists(ctx, cfg, arn)
	if err != nil {
		return pkgerrors.Wrapf(err, "awswafipset verify-absent failed for %q", arn)
	}
	if exists {
		return pkgerrors.Errorf("awswafipset %q still exists after destroy", arn)
	}
	return nil
}

func wafIpSetExists(ctx context.Context, cfg aws.Config, arn string) (bool, error) {
	scope, name, id, region, err := wafv2ArnParts(arn)
	if err != nil {
		return false, err
	}
	_, err = wafv2Client(cfg, region).GetIPSet(ctx, &wafv2.GetIPSetInput{
		Id:    aws.String(id),
		Name:  aws.String(name),
		Scope: scope,
	})
	if err != nil {
		if isWafNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// wafRegexPatternSetVerifier verifies an AwsWafRegexPatternSet via
// GetRegexPatternSet, keyed on the regex_pattern_set_arn output.
type wafRegexPatternSetVerifier struct{}

func (*wafRegexPatternSetVerifier) IDOutputKey() string { return "regex_pattern_set_arn" }

func (*wafRegexPatternSetVerifier) VerifyExists(ctx context.Context, cfg aws.Config, arn, region string) error {
	exists, err := wafRegexPatternSetExists(ctx, cfg, arn)
	if err != nil {
		return pkgerrors.Wrapf(err, "awswafregexpatternset verify-exists failed for %q", arn)
	}
	if !exists {
		return pkgerrors.Errorf("awswafregexpatternset %q not found after deploy", arn)
	}
	return nil
}

func (*wafRegexPatternSetVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, arn, region string) error {
	exists, err := wafRegexPatternSetExists(ctx, cfg, arn)
	if err != nil {
		return pkgerrors.Wrapf(err, "awswafregexpatternset verify-absent failed for %q", arn)
	}
	if exists {
		return pkgerrors.Errorf("awswafregexpatternset %q still exists after destroy", arn)
	}
	return nil
}

func wafRegexPatternSetExists(ctx context.Context, cfg aws.Config, arn string) (bool, error) {
	scope, name, id, region, err := wafv2ArnParts(arn)
	if err != nil {
		return false, err
	}
	_, err = wafv2Client(cfg, region).GetRegexPatternSet(ctx, &wafv2.GetRegexPatternSetInput{
		Id:    aws.String(id),
		Name:  aws.String(name),
		Scope: scope,
	})
	if err != nil {
		if isWafNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

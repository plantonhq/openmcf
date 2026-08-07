package verify

import (
	"context"
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
	pkgerrors "github.com/pkg/errors"
)

// eventBridgeRuleVerifier verifies an AwsEventBridgeRule via DescribeRule.
// Custom-bus rules need the bus name as well as the rule name; rule_name alone
// is insufficient, so verification reads rule_arn to derive the bus when the
// ARN carries a bus segment (arn:...:rule/<bus>/<rule>).
type eventBridgeRuleVerifier struct{}

func (*eventBridgeRuleVerifier) IDOutputKey() string { return "rule_name" }

func (*eventBridgeRuleVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	return pkgerrors.Errorf("awseventbridgerule verify-exists requires full outputs (rule_arn); use OutputsVerifier path")
}

func (*eventBridgeRuleVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	return pkgerrors.Errorf("awseventbridgerule verify-absent requires full outputs (rule_arn); use OutputsVerifier path")
}

func (*eventBridgeRuleVerifier) VerifyExistsFromOutputs(ctx context.Context, cfg aws.Config, outputs map[string]interface{}, region string) error {
	busName, ruleName, err := eventBridgeRuleLookup(outputs)
	if err != nil {
		return pkgerrors.Wrap(err, "awseventbridgerule verify-exists")
	}
	exists, err := eventBridgeRuleExists(ctx, cfg, busName, ruleName, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awseventbridgerule verify-exists failed for bus %q rule %q", busName, ruleName)
	}
	if !exists {
		return pkgerrors.Errorf("awseventbridgerule %q on bus %q not found after deploy", ruleName, busName)
	}
	return nil
}

func (*eventBridgeRuleVerifier) VerifyAbsentFromOutputs(ctx context.Context, cfg aws.Config, outputs map[string]interface{}, region string) error {
	busName, ruleName, err := eventBridgeRuleLookup(outputs)
	if err != nil {
		return pkgerrors.Wrap(err, "awseventbridgerule verify-absent")
	}
	exists, err := eventBridgeRuleExists(ctx, cfg, busName, ruleName, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awseventbridgerule verify-absent failed for bus %q rule %q", busName, ruleName)
	}
	if exists {
		return pkgerrors.Errorf("awseventbridgerule %q on bus %q still exists after destroy", ruleName, busName)
	}
	return nil
}

func eventBridgeRuleLookup(outputs map[string]interface{}) (busName, ruleName string, err error) {
	ruleName = stringOutputMap(outputs, "rule_name")
	if ruleName == "" {
		return "", "", pkgerrors.New("no rule_name in outputs -- cannot verify")
	}
	ruleARN := stringOutputMap(outputs, "rule_arn")
	busName = "default"
	if ruleARN != "" {
		if bus, name, ok := parseEventBusFromRuleARN(ruleARN); ok {
			busName = bus
			if ruleName == "" {
				ruleName = name
			}
		}
	}
	return busName, ruleName, nil
}

// parseEventBusFromRuleARN extracts the event bus and rule name from an
// EventBridge rule ARN. Default-bus rules use arn:...:rule/<rule-name>;
// custom-bus rules use arn:...:rule/<bus-name>/<rule-name>.
func parseEventBusFromRuleARN(ruleARN string) (busName, ruleName string, ok bool) {
	const marker = ":rule/"
	idx := strings.Index(ruleARN, marker)
	if idx < 0 {
		return "", "", false
	}
	rest := ruleARN[idx+len(marker):]
	parts := strings.SplitN(rest, "/", 2)
	if len(parts) == 2 {
		return parts[0], parts[1], true
	}
	return "default", parts[0], true
}

func eventBridgeRuleExists(ctx context.Context, cfg aws.Config, busName, ruleName, region string) (bool, error) {
	client := eventbridge.NewFromConfig(cfg, func(o *eventbridge.Options) {
		if region != "" {
			o.Region = region
		}
	})
	in := &eventbridge.DescribeRuleInput{Name: aws.String(ruleName)}
	if busName != "" && busName != "default" {
		in.EventBusName = aws.String(busName)
	}
	_, err := client.DescribeRule(ctx, in)
	if err == nil {
		return true, nil
	}
	var notFound *types.ResourceNotFoundException
	if errors.As(err, &notFound) {
		return false, nil
	}
	return false, err
}
